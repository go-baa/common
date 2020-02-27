package base

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-baa/baa"
	"github.com/go-baa/cache"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
	"github.com/jinzhu/gorm"

	// 导入mysql驱动
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// MapParams 声明一个通用的参数结构
type MapParams map[string]interface{}

// DbConfig database config struct
type DbConfig struct {
	Type, Host, Name, User, Passwd, Path, SSLMode string
}

// Errorf 对fmt.Errorf()的一个包装
func Errorf(format string, a ...interface{}) error {
	if len(a) > 0 {
		return fmt.Errorf(format, a...)
	}
	return fmt.Errorf(format)
}

// LoadConfigs 加载数据库配置
func LoadConfigs(name string) *DbConfig {
	config := new(DbConfig)
	config.Host = setting.Config.MustString("db."+name+".host", "")
	config.Name = setting.Config.MustString("db."+name+".name", "")
	config.User = setting.Config.MustString("db."+name+".user", "")
	config.Passwd = setting.Config.MustString("db."+name+".pass", "")
	return config
}

// setLogger 切换日志
func setLogger(db *gorm.DB, date string) {
	logpath := strings.TrimRight(setting.Config.MustString("orm.logpath", "data/log"), "/")
	logfile := logpath + "/orm-" + date + ".log"

	if logpath == "os.Stderr" || logpath == "os.Stdout" {
		db.SetLogger(gorm.Logger{LogWriter: log.New(os.Stdout, "[orm]", 0)})
	} else {
		os.MkdirAll(path.Dir(logfile), os.ModePerm)
		f, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err == nil {
			db.SetLogger(log.New(f, "[orm]", 0))
		}
	}
}

func getEngine(config *DbConfig) (*gorm.DB, error) {
	cnnstr := ""
	if config.Host[0] == '/' { // looks like a unix socket
		cnnstr = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8mb4&timeout=3s&parseTime=true&loc=Local",
			config.User, config.Passwd, config.Host, config.Name)
	} else {
		cnnstr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&timeout=3s&parseTime=true&loc=Local",
			config.User, config.Passwd, config.Host, config.Name)
	}
	return gorm.Open("mysql", cnnstr)
}

// NewEngine ...
func NewEngine(config *DbConfig) (*gorm.DB, error) {
	db, err := getEngine(config)
	if err != nil {
		return nil, fmt.Errorf("Fail to connect to database: %v", err)
	}

	// 关闭tableName自动复数
	db.SingularTable(true)

	// 默认不打印日志
	db.LogMode(false)

	// 设置日志
	if baa.Env != baa.PROD {
		// 设置日志
		db.LogMode(true)
		date := time.Now().Format("2006-01-02")
		setLogger(db, date)
		// 设置日志轮询
		go func(x *gorm.DB, date string) {
			var dateNew string
			for {
				dateNew = time.Now().Format("2006-01-02")
				if dateNew == date {
					time.Sleep(time.Second * 1)
					continue
				}
				date = dateNew
				setLogger(db, date)
			}
		}(db, date)
	}

	return db, nil
}

// DB 在gorm.DB基础上封装了嵌套事务的支持，实际上是一次性提交
type DB struct {
	*gorm.DB
	ox               *gorm.DB
	tx               bool
	transactionLevel int
}

// NewDB 创建新的数据库连接
func NewDB(db *gorm.DB, tx bool) *DB {
	if tx {
		return &DB{db, nil, true, 0}
	}
	return &DB{db, nil, false, 0}
}

// Begin 开启事务
func (t *DB) Begin() *DB {
	if !t.tx {
		log.Panic("[orm] db.Begin error: current connection not support transaction\n")
	}
	if t.transactionLevel == 0 {
		t.ox = t.DB
		t.DB = t.DB.Begin()
	}
	t.transactionLevel++
	return t
}

// Rollback 回滚事务
func (t *DB) Rollback() *gorm.DB {
	t.transactionLevel--
	if t.transactionLevel == 0 {
		tx := t.DB.Rollback()
		t.DB = t.ox
		return tx
	} else if t.transactionLevel < 0 {
		log.Panic("[orm] db.Rollback error: over transaction level\n")
	}
	return t.DB
}

// Commit 提交事务
func (t *DB) Commit() *gorm.DB {
	t.transactionLevel--
	if t.transactionLevel == 0 {
		tx := t.DB.Commit()
		t.DB = t.ox
		return tx
	} else if t.transactionLevel < 0 {
		log.Panic("[orm] db.Commit error: over transaction level\n")
	}
	return t.DB
}

// MustCommit 强制提交所有事务，跳过层级检查
func (t *DB) MustCommit() *gorm.DB {
	tx := t.DB.Commit()
	t.DB = t.ox
	t.transactionLevel = 0
	return tx
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (t *DB) Save(value interface{}) *gorm.DB {
	if !t.tx {
		log.Panic("[orm] db.Save error: current connection not support modify data\n")
	}
	return t.DB.Save(value)
}

// FirstOrCreate find first matched record or create a new one with given conditions (only works with struct, map conditions)
func (t *DB) FirstOrCreate(out interface{}, where ...interface{}) *gorm.DB {
	if !t.tx {
		log.Panic("[orm] db.Save error: current connection not support modify data\n")
	}
	return t.DB.FirstOrCreate(out, where...)
}

// Update update attributes with callbacks
func (t *DB) Update(attrs ...interface{}) *gorm.DB {
	if !t.tx {
		log.Panic("[orm] db.Save error: current connection not support modify data\n")
	}
	return t.DB.Update(attrs...)
}

// Updates update attributes with callbacks
func (t *DB) Updates(values interface{}, ignoreProtectedAttrs ...bool) *gorm.DB {
	if !t.tx {
		log.Panic("[orm] db.Save error: current connection not support modify data\n")
	}
	return t.DB.Updates(values, ignoreProtectedAttrs...)
}

// UpdateColumn update attributes without callbacks
func (t *DB) UpdateColumn(attrs ...interface{}) *gorm.DB {
	if !t.tx {
		log.Panic("[orm] db.Save error: current connection not support modify data\n")
	}
	return t.DB.UpdateColumn(attrs...)
}

// UpdateColumns update attributes without callbacks
func (t *DB) UpdateColumns(values interface{}) *gorm.DB {
	if !t.tx {
		log.Panic("[orm] db.Save error: current connection not support modify data\n")
	}
	return t.DB.UpdateColumns(values)
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
func (t *DB) Delete(value interface{}, where ...interface{}) *gorm.DB {
	if !t.tx {
		log.Panic("[orm] db.Save error: current connection not support modify data\n")
	}
	return t.DB.Delete(value, where...)
}

// Create insert the value into database
func (t *DB) Create(value interface{}) *gorm.DB {
	if !t.tx {
		log.Panic("[orm] db.Save error: current connection not support create data\n")
	}
	return t.DB.Create(value)
}

// Raw use raw sql as conditions, won't run it unless invoked by other methods
//    db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
// func (t *DB) Raw(sql string, values ...interface{}) *gorm.DB {
// 	if !t.tx {
// 		log.Panic("[orm] db.Save error: current connection not support raw sql\n")
// 	}
// 	return t.DB.Raw(sql, values...)
// }

// Exec execute raw sql
func (t *DB) Exec(sql string, values ...interface{}) *gorm.DB {
	if !t.tx {
		log.Panic("[orm] db.Save error: current connection not support raw sql\n")
	}
	return t.DB.Exec(sql, values...)
}

// Cacher 获取缓存控制
func Cacher() cache.Cacher {
	if c := baa.Default().GetDI("cache"); c != nil {
		return c.(cache.Cacher)
	}
	return nil
}
