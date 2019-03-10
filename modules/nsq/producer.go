package nsq

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
	nsq "github.com/nsqio/go-nsq"
)

const (
	// LookupdPollInterval ...
	LookupdPollInterval = 15 * time.Second
)

// Producer nsq producer manager
type Producer struct {
	nsqLookupdAddrs   []string
	lookupdQueryIndex int
	locker            sync.RWMutex
	nsqdAddrs         []string
	nsqdAddrIndex     int
	conn              map[string]*nsq.Producer
}

// NewProducer creates a new instance of Producer
func NewProducer(nsqLookupdAddrs ...string) *Producer {
	p := &Producer{
		conn: make(map[string]*nsq.Producer),
	}
	go func() {
		p.connectToNSQLookupds(nsqLookupdAddrs)
	}()
	return p
}

// Publish publish a msg to a topic
func (t *Producer) Publish(topic string, body []byte) error {
	var conn *nsq.Producer
	t.locker.RLock()
	if len(t.conn) == 0 {
		t.locker.RUnlock()
		return fmt.Errorf("no alive nsqd")
	}

	nsqPrefix := setting.Config.MustString("nsq.prefix", "")
	if nsqPrefix != "" {
		topic = nsqPrefix + "_" + topic
	}
	addr := t.nextNSQDAddr()
	conn = t.conn[addr]
	t.locker.RUnlock()
	if setting.Debug {
		log.Debugf("pub:\n    conn: %v\n    topic: %s\n    body: %s\n", conn, topic, body)
	}
	return conn.Publish(topic, body)
}

// Count return producer count
func (t *Producer) Count() int {
	t.locker.RLock()
	num := len(t.conn)
	t.locker.RUnlock()
	return num
}

// connectToNSQLookupds adds multiple nsqlookupd address to the list for this Producer instance.
func (t *Producer) connectToNSQLookupds(addresses []string) error {
	for _, addr := range addresses {
		err := t.connectToNSQLookupd(addr)
		if err != nil {
			return err
		}
	}
	return nil
}

// connectToNSQLookupd add an nsqlookupd address to the list for this Producer instance.
func (t *Producer) connectToNSQLookupd(addr string) error {
	if err := validatedLookupAddr(addr); err != nil {
		return err
	}
	t.locker.Lock()
	for _, v := range t.nsqLookupdAddrs {
		if v == addr {
			t.locker.Unlock()
			return nil
		}
	}
	t.nsqLookupdAddrs = append(t.nsqLookupdAddrs, addr)
	numLookupdAddrs := len(t.nsqLookupdAddrs)
	t.locker.Unlock()

	if numLookupdAddrs == 1 {
		go t.lookupdLoop()
	}

	return nil
}

// poll all known lookup servers every LookupdPollInterval
func (t *Producer) lookupdLoop() {
	var ticker *time.Ticker

	ticker = time.NewTicker(LookupdPollInterval)
	t.lookupAddr()
	for {
		select {
		case <-ticker.C:
			t.lookupAddr()
		}
	}
}

// LookupAddr query alive nsqd services addr
func (t *Producer) lookupAddr() {
	data := new(lookupResp)
	addr := t.nextLookupdEndpoint()
	err := apiRequestNegotiateV1("GET", addr, nil, data)
	if err != nil {
		log.Errorf("[nsq] err: querying nsqlookupd (%s) - %s\n", addr, err)
		return
	}

	// 当前存活的所有节点
	var addrs []string
	for _, producer := range data.Producers {
		broadcastAddress := producer.BroadcastAddress
		port := producer.TCPPort
		addr := net.JoinHostPort(broadcastAddress, strconv.Itoa(port))
		addrs = append(addrs, addr)
	}
	t.updateNSQDConn(addrs)
}

// updateNSQDConn update conn list of this Producer instance
func (t *Producer) updateNSQDConn(addrs []string) {
	t.locker.Lock()
	delAddrList := util.SliceStringDiff(t.nsqdAddrs, addrs)
	for _, addr := range delAddrList {
		if t.nsqdAddrIndex > 0 {
			t.nsqdAddrIndex--
		}
		idx := indexOf(addr, t.nsqdAddrs)
		t.nsqdAddrs = append(t.nsqdAddrs[:idx], t.nsqdAddrs[idx+1:]...)
		delete(t.conn, addr)
	}

	addAddrList := util.SliceStringDiff(addrs, t.nsqdAddrs)
	for _, addr := range addAddrList {
		if producer, err := nsq.NewProducer(addr, nsq.NewConfig()); err == nil {
			t.nsqdAddrs = append(t.nsqdAddrs, addr)
			t.conn[addr] = producer
		}
	}

	if len(t.nsqdAddrs) == 0 {
		log.Println("[nsq] info: nsqd service list empty")
	}

	t.locker.Unlock()
}

type lookupResp struct {
	Producers []*peerInfo `json:"producers"`
}

type peerInfo struct {
	RemoteAddress    string `json:"remote_address"`
	Hostname         string `json:"hostname"`
	BroadcastAddress string `json:"broadcast_address"`
	TCPPort          int    `json:"tcp_port"`
	HTTPPort         int    `json:"http_port"`
	Version          string `json:"version"`
}

// return the next nsqd addr to publish msg
// keep track of which one was last used
func (t *Producer) nextNSQDAddr() string {
	if t.nsqdAddrIndex >= len(t.nsqdAddrs) {
		t.nsqdAddrIndex = 0
	}
	addr := t.nsqdAddrs[t.nsqdAddrIndex]
	num := len(t.nsqLookupdAddrs)
	t.nsqdAddrIndex = (t.nsqdAddrIndex + 1) % num

	return addr
}

// return the next lookupd endpoint to query
// keeping track of which one was last used
func (t *Producer) nextLookupdEndpoint() string {
	t.locker.RLock()
	if t.lookupdQueryIndex >= len(t.nsqLookupdAddrs) {
		t.lookupdQueryIndex = 0
	}
	addr := t.nsqLookupdAddrs[t.lookupdQueryIndex]
	num := len(t.nsqLookupdAddrs)
	t.locker.RUnlock()
	t.lookupdQueryIndex = (t.lookupdQueryIndex + 1) % num

	urlString := addr
	if !strings.Contains(urlString, "://") {
		urlString = "http://" + addr
	}

	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	if u.Path == "/" || u.Path == "" {
		u.Path = "/nodes"
	}

	v, err := url.ParseQuery(u.RawQuery)
	u.RawQuery = v.Encode()
	return u.String()
}

// validate lookupd address format
func validatedLookupAddr(addr string) error {
	if strings.Contains(addr, "/") {
		_, err := url.Parse(addr)
		if err != nil {
			return err
		}
		return nil
	}
	if !strings.Contains(addr, ":") {
		return errors.New("missing port")
	}
	return nil
}

// find index of string n in slice h
// if not found, return -1
func indexOf(n string, h []string) int {
	for i, a := range h {
		if n == a {
			return i
		}
	}
	return -1
}
