package nsq

import (
	"os"

	"github.com/go-baa/log"
	"github.com/go-baa/setting"
	nsq "github.com/nsqio/go-nsq"
)

const (
	// DefaultMaxInFlight 默认MaxInFlight配置，根据nsqd服务数可以动态调整
	DefaultMaxInFlight = 3
)

// Handler consumer handler
type Handler struct {
	Topic      string
	Channel    string
	Config     *nsq.Config
	MsgHandler nsq.Handler
}

// Consumer ...
type Consumer struct {
	nsqLookupdAddrs []string
	handlers        []*Handler
}

// NewConsumer creates a new instance of Consumer
func NewConsumer(nsqLookupdAddrs ...string) *Consumer {
	return &Consumer{
		nsqLookupdAddrs: nsqLookupdAddrs,
		handlers:        make([]*Handler, 0),
	}
}

// DefaultConfig 默认配置
func DefaultConfig() *nsq.Config {
	conf := nsq.NewConfig()
	conf.MaxInFlight = DefaultMaxInFlight
	return conf
}

// AddHandler add handler
func (t *Consumer) AddHandler(handler *Handler) {
	t.handlers = append(t.handlers, handler)
}

// Run run consumer
func (t *Consumer) Run() {
	log.Println("[consumer] starting ...")
	for _, handler := range t.handlers {
		go t.subscribe(handler)
	}
	select {}
}

// subscribe 启动topic订阅
func (t *Consumer) subscribe(h *Handler) {
	nsqPrefix := setting.Config.MustString("nsq.prefix", "")
	if nsqPrefix != "" {
		h.Topic = nsqPrefix + "_" + h.Topic
	}
	consumer, err := nsq.NewConsumer(h.Topic, h.Channel, h.Config)
	if err != nil {
		log.Fatalf("[%s:%s] init error: %v", h.Topic, h.Channel, err)
	} else {
		log.Printf("[%s:%s] init ok", h.Topic, h.Channel)
	}
	consumer.SetLogger(log.New(os.Stderr, "NSQ", log.Flags()), nsq.LogLevelError)

	consumer.AddHandler(h.MsgHandler)
	err = consumer.ConnectToNSQLookupds(t.nsqLookupdAddrs)
	if err != nil {
		log.Fatalf("[%s:%s] connect lookup error: %v", h.Topic, h.Channel, err)
	}

	<-consumer.StopChan
	log.Printf("[%s:%s] stopped subscribe", h.Topic, h.Channel)
}
