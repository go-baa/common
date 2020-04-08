// Package accesslog 提供了一个baa的访问日志中间件
package accesslog

import (
	"fmt"
	"time"

	"github.com/go-baa/baa"
)

// DefaultFlushTTL 日志缓存时间，单位：秒
const DefaultFlushTTL = 10

// Logger 访问日志接口
type Logger interface {
	Log(string)
	Flush() error
	Config(*Options) error
}

// Options 访问日志配置
type Options struct {
	Open     bool                   // 日志是否开启
	Format   string                 // 日志格式
	FlushTTL time.Duration          // 日志缓存刷新时间，单位：秒
	Adapter  string                 // 适配器, file, flume
	Config   map[string]interface{} // 适配器配置
}

// Get 从配置中获取一个值
func (o Options) Get(key string, defaultValue interface{}) interface{} {
	v, exists := o.Config[key]
	if !exists {
		return defaultValue
	}
	return v
}

type instanceFunc func() Logger

var adapters = make(map[string]instanceFunc)

// New 创建一个新的访问记录中间件
func New(o Options) baa.Middleware {
	if !o.Open {
		return func(c *baa.Context) { c.Next() }
	}

	if len(o.Adapter) == 0 {
		panic("accesslog.New: cannot use empty adapter")
	}

	// adapter
	f, ok := adapters[o.Adapter]
	if !ok {
		panic("accesslog.New: unknown adapter (forgot to import?)")
	}
	adapter := f()

	// Set default logger format
	if o.Format == "" {
		o.Format = DefaultFormatter
	}

	date := time.Now().Format("2006-01-02")
	lopath, ok := o.Config["path"].(string)
	if ok && lopath != "" {
		o.Config["file"] = lopath + "/access-" + date + ".log"
	}
	if err := adapter.Config(&o); err != nil {
		panic(fmt.Sprintf("accesslog.New: %s incorrect configuration, %s", o.Adapter, err.Error()))
	}

	// 设置日志轮询
	go func(date string) {
		var dateNew string
		for {
			dateNew = time.Now().Format("2006-01-02")
			if dateNew == date {
				time.Sleep(time.Second * 1)
				continue
			}
			date = dateNew
			lopath, ok := o.Config["path"].(string)
			if ok && lopath != "" {
				o.Config["file"] = lopath + "/access-" + date + ".log"
			}
			if err := adapter.Config(&o); err != nil {
				panic(fmt.Sprintf("accesslog.New: %s incorrect configuration, %s", o.Adapter, err.Error()))
			}
		}
	}(date)

	// 刷新日志缓冲
	go func() {
		ttl := o.FlushTTL
		if ttl == 0 {
			ttl = DefaultFlushTTL
		}
		ticker := time.NewTicker(time.Second * ttl)
		for _ = range ticker.C {
			adapter.Flush()
		}
	}()

	// new formater
	formater := newFormatter(o.Format)

	return func(c *baa.Context) {
		start := time.Now()

		c.Next()

		// 异步日志
		go func() {
			line := formater.build(c, start)
			adapter.Log(line)
		}()
	}
}

// Register 注册新的适配器
func Register(name string, adapter instanceFunc) {
	if adapter == nil {
		panic("accesslog.Register: cannot register adapter with nil func")
	}
	if _, ok := adapters[name]; ok {
		panic(fmt.Errorf("accesslog.Register: cannot register adapter '%s' twice", name))
	}
	adapters[name] = adapter
}
