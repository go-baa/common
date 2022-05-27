package base

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// 时间格式
const (
	// TimeFormatDefault 时间默认格式
	TimeFormatDefault = "2006-01-02 15:04:05"
	// TimeFormatDate 日期格式
	TimeFormatDate = "2006-01-02"
	// TimeFormatMonth 日期格式
	TimeFormatMonth = "2006-01"
)

// DateTime 格式：2006-01-02 15:04:05 再定义time时间类型
type DateTime struct {
	Time  time.Time
	Valid bool
}

// NewDateTime 创建
func NewDateTime(t time.Time, valid bool) DateTime {
	return DateTime{
		Time:  t,
		Valid: valid,
	}
}

// DateTimeFromPtr 初始化指针到自定义time类型
func DateTimeFromPtr(t *time.Time) DateTime {
	if t == nil {
		return NewDateTime(time.Time{}, false)
	}
	return NewDateTime(*t, true)
}

// MarshalJSON implements json.Marshaler.
func (t DateTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("\"\""), nil
	}

	b := make([]byte, 0, len(TimeFormatDefault)+2)
	b = append(b, '"')
	b = t.Time.AppendFormat(b, TimeFormatDefault)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *DateTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "\"\"" {
		t.Valid = false
		t.Time = time.Time{}
		return nil
	}

	if bytes.Equal(data, nil) {
		t.Valid = false
		t.Time = time.Time{}
		return nil
	}
	timeFmt := TimeFormatDefault
	dataStr := string(data)
	if strings.HasPrefix(strings.Trim(dataStr, "\""), "0001") {
		t.Valid = false
		t.Time = time.Time{}
		return nil
	}
	if strings.Contains(dataStr, "T") && strings.Contains(dataStr, "+") {
		timeFmt = time.RFC3339
	}
	now, err := time.ParseInLocation(`"`+timeFmt+`"`, dataStr, time.Local)
	if err != nil {
		return err
	}
	t.Time = now
	t.Valid = true
	return nil
}

// SetValid 设置时间
func (t *DateTime) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// SetValidByString 设置string时间
func (t *DateTime) SetValidByString(v string) (err error) {
	if v == "" {
		t.Valid = false
		return
	}
	timeFmt := TimeFormatDefault
	if strings.Contains(v, "T") && strings.Contains(v, "+") {
		timeFmt = time.RFC3339
	}
	t.Time, err = time.ParseInLocation(timeFmt, v, time.Local)
	if err == nil {
		t.Valid = true
	}
	return
}

// Ptr 当前结构体转换为指针结构
func (t DateTime) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsNull 判断是否为空
func (t DateTime) IsNull() bool {
	return !t.Valid
}

// Scan implements the Scanner interface.
func (t *DateTime) Scan(src interface{}) error {
	var err error
	switch x := src.(type) {
	case time.Time:
		if x.IsZero() || strings.HasPrefix(x.String(), "0001-01-01") {
			t.Valid = false
			return nil
		}
		t.Time = x
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", src, src)
	}
	t.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (t DateTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// String toString方法
func (t DateTime) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(TimeFormatDefault)
}

// AgoFormat 时间格式化
func (t DateTime) AgoFormat() string {
	now := time.Now()
	interval := now.Unix() - t.Time.Unix()
	if interval <= 60 {
		return "1分钟前"
	} else if interval > 60 && interval <= 3600 {
		return fmt.Sprintf("%d分钟前", int(interval/60))
	} else if interval > 3600 && interval <= 86400 {
		return fmt.Sprintf("%d小时前", int(interval/3600))
	} else if t.Time.Unix() > now.Add(-time.Second*86400).Unix() {
		return "昨天"
	} else {
		return t.Time.Format("2006-01-02")
	}
}

// Date 格式：2006-01-02 再定义date日期类型
type Date struct {
	Date  time.Time
	Valid bool
}

// NewDate 创建
func NewDate(t time.Time, valid bool) Date {
	return Date{
		Date:  t,
		Valid: valid,
	}
}

// DateFromPtr 初始化指针到自定义date类型
func DateFromPtr(t *time.Time) Date {
	if t == nil {
		return NewDate(time.Time{}, false)
	}
	return NewDate(*t, true)
}

// MarshalJSON implements json.Marshaler.
func (t Date) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("\"\""), nil
	}

	b := make([]byte, 0, len(TimeFormatDate)+2)
	b = append(b, '"')
	b = t.Date.AppendFormat(b, TimeFormatDate)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Date) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "\"\"" {
		t.Valid = false
		t.Date = time.Time{}
		return nil
	}

	if bytes.Equal(data, nil) {
		t.Valid = false
		t.Date = time.Time{}
		return nil
	}
	timeFmt := TimeFormatDate
	dataStr := string(data)
	if strings.HasPrefix(strings.Trim(dataStr, "\""), "0001") {
		t.Valid = false
		t.Date = time.Time{}
		return nil
	}
	if strings.Contains(dataStr, "T") && strings.Contains(dataStr, "+") {
		timeFmt = time.RFC3339
	}
	now, err := time.ParseInLocation(`"`+timeFmt+`"`, dataStr, time.Local)
	if err != nil {
		return err
	}
	t.Date = now
	t.Valid = true
	return nil
}

// SetValid 设置时间
func (t *Date) SetValid(v time.Time) {
	t.Date = v
	t.Valid = true
}

// SetValidByString 设置string时间
func (t *Date) SetValidByString(v string) (err error) {
	if v == "" {
		t.Valid = false
		return
	}
	timeFmt := TimeFormatDate
	if strings.Contains(v, "T") && strings.Contains(v, "+") {
		timeFmt = time.RFC3339
	}
	t.Date, err = time.ParseInLocation(timeFmt, v, time.Local)
	if err == nil {
		t.Valid = true
	}
	return
}

// Ptr 当前结构体转换为指针结构
func (t Date) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Date
}

// IsNull 判断是否为空
func (t Date) IsNull() bool {
	return !t.Valid
}

// Scan implements the Scanner interface.
func (t *Date) Scan(src interface{}) error {
	var err error
	switch x := src.(type) {
	case time.Time:
		t.Date = x
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", src, src)
	}
	t.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (t Date) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Date, nil
}

// String toString方法
func (t Date) String() string {
	if !t.Valid {
		return ""
	}
	return t.Date.Format(TimeFormatDate)
}
