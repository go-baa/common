// Package consul provider consul service discovery and health check
package consul

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/hashicorp/consul/api"
)

// Consul consul api client
type Consul struct {
	client    *api.Client
	nodeAddr  string // 当前的consul节点服务器
	localAddr string // 应用的本地IP地址
}

// KV consul KV client
type KV struct {
	client *api.KV
}

// New returns a consul api client
func New(addrs ...string) (*Consul, error) {
	config := api.DefaultConfig()
	addr := chooseAddr(addrs...)
	if addr == "" {
		panic("consul.New: addr is empty or addrs cannot connect")
	}
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	// 得到去除端口号的地址
	var nodeAddr string
	splitPos := strings.IndexByte(addr, ':')
	if splitPos > 0 {
		nodeAddr = addr[:splitPos]
	} else {
		nodeAddr = addr
	}

	return &Consul{client, nodeAddr, ""}, nil
}

// ServiceRegisger register service and health check
func (t *Consul) ServiceRegisger(id, name, addr, port string, tags []string) error {
	if id == "" {
		if id = os.Getenv("CSPHERE_CONTAINER_NAME"); id == "" {
			id, _ = os.Hostname()
			id = id + "." + name
		}
	}

	if addr == "" {
		addr = getIPAddress()
	}

	// 尝试查找上次注册的节点
	client := t.client
	nodeAddr := t.ServiceQueryNode(name, "", addr)
	if nodeAddr != "" {
		config := api.DefaultConfig()
		config.Address = nodeAddr + ":8500"
		newClient, err := api.NewClient(config)
		if err == nil {
			client = newClient
		}
	}

	check := &api.AgentServiceCheck{
		Status: "passing",
		TTL:    "30s",
		DeregisterCriticalServiceAfter: "2h",
	}
	service := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Port:    util.StringToInt(port),
		Address: addr,
		Check:   check,
	}
	err := client.Agent().ServiceRegister(service)
	if err != nil {
		return err
	}

	// begin health check notice
	go func(client *api.Client, id string) {
		for {
			time.Sleep(time.Second * 10)
			client.Agent().UpdateTTL(id, time.Now().Format(time.RFC1123Z), "pass")
		}
	}(client, "service:"+id)

	return nil
}

// ServiceQuery query a service and return a normal service addr
func (t *Consul) ServiceQuery(name, tag string) (string, error) {
	entris, _, err := t.client.Health().Service(name, tag, true, nil)
	if err != nil {
		return "", err
	}
	if len(entris) == 0 {
		return "", fmt.Errorf("consul.ServiceQuery: %v:%v none health node is available", name, tag)
	}
	rand.Seed(time.Now().Unix())
	n := rand.Intn(len(entris))
	return entris[n].Service.Address + ":" + util.IntToString(entris[n].Service.Port), nil
}

// ServiceQueryNode query node address of a local service
func (t *Consul) ServiceQueryNode(name, tag, localAddr string) string {
	// localAddr := getIPAddress()
	catalog, _, err := t.client.Catalog().Service(name, tag, nil)
	if err != nil {
		return ""
	}
	for _, agent := range catalog {
		if localAddr == agent.ServiceAddress {
			return agent.Address
		}
	}

	return ""
}

// KV returns consul KV handle
func (t *Consul) KV() *KV {
	return &KV{t.client.KV()}
}

// Get get a value from kv
func (t *KV) Get(key string) ([]byte, error) {
	pair, _, err := t.client.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, fmt.Errorf("consul.KV error: key [%s] is not exist", key)
	}
	return pair.Value, nil
}

// Put set a value to kv
func (t *KV) Put(key string, val []byte) error {
	if len(key) == 0 {
		return fmt.Errorf("consul.KV error: key [%s] is empty", key)
	}
	if key[0] == '/' {
		key = key[1:]
	}
	_, err := t.client.Put(&api.KVPair{Key: key, Value: val}, nil)
	return err
}

// Delete delete a key from kv
func (t *KV) Delete(key string) error {
	if len(key) == 0 {
		return fmt.Errorf("consul.KV error: key [%s] is empty", key)
	}
	if key[0] == '/' {
		key = key[1:]
	}
	_, err := t.client.Delete(key, nil)
	return err
}

// chooseAddr rand check service addr and returns
func chooseAddr(addrs ...string) string {
	used := make(map[int]struct{})
	l := len(addrs)
	rand.Seed(time.Now().Unix())
	for {
		// none can be used
		if len(used) == l {
			break
		}
		// rand a addr
		n := rand.Intn(l)
		if _, ok := used[n]; ok {
			continue
		}
		used[n] = struct{}{}
		conn, err := net.Dial("tcp4", addrs[n])
		if err != nil {
			continue
		}
		conn.Close()
		return addrs[n]
	}
	return ""
}

func getIPAddress() string {
	addr := os.Getenv("CSPHERE_CONTAINER_IP")
	if addr == "" {
		if localAddrs := util.LocalIPAddrs(); localAddrs != nil {
			addr = localAddrs[0]
		} else {
			addr = "127.0.0.1"
		}
	}
	return addr
}
