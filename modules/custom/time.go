package custom

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"
)

// 时间格式
const (
	// TimeFormatDefault 时间默认格式
	TimeFormatDefault = "2006-01-02 15:04:05"
	// TimeFormatDate 日期格式
	TimeFormatDate = "2006-01-02"
)

// Time 格式：2006-01-02 15:04:05 再定义time时间类型
type Time struct {
	Time  time.Time
	Valid bool
}

// NewTime 创建
func NewTime(t time.Time, valid bool) Time {
	return Time{
		Time:  t,
		Valid: valid,
	}
}

// TimeFromPtr 初始化指针到自定义time类型
func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, false)
	}
	return NewTime(*t, true)
}

// MarshalJSON implements json.Marshaler.
func (t Time) MarshalJSON() ([]byte, error) {
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
func (t *Time) UnmarshalJSON(data []byte) error {
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
	now, err := time.ParseInLocation(`"`+TimeFormatDefault+`"`, string(data), time.Local)
	if err != nil {
		return err
	}
	t.Time = now
	t.Valid = true
	return nil
}

// SetValid 设置时间
func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// SetValidByString 设置string时间
func (t *Time) SetValidByString(v string) (err error) {
	if v == "" {
		t.Valid = false
		return
	}
	t.Time, err = time.ParseInLocation(TimeFormatDefault, v, time.Local)
	t.Valid = true
	return
}

// Ptr 当前结构体转换为指针结构
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsNull 判断是否为空
func (t Time) IsNull() bool {
	return !t.Valid
}

// Scan implements the Scanner interface.
func (t *Time) Scan(src interface{}) error {
	var err error
	switch x := src.(type) {
	case time.Time:
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
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// String toString方法
func (t Time) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(TimeFormatDefault)
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
	now, err := time.ParseInLocation(`"`+TimeFormatDate+`"`, string(data), time.Local)
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
	t.Date, err = time.ParseInLocation(TimeFormatDate, v, time.Local)
	t.Valid = true
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
