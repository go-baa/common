package custom

import (
	"database/sql"
	"encoding/json"
	"strconv"
)

// NullString 可为nil的string类型
type NullString struct {
	sql.NullString
}

// IsNull 是否为空
func (v NullString) IsNull() bool {
	return !v.Valid
}

// SetValide 设值
func (v NullString) SetValide(s string) {
	v.String = s
	v.Valid = true
}

// NewNullString 初始化string类型
func NewNullString(s string, valid bool) NullString {
	var n NullString
	n.Valid = valid
	n.String = s
	return n
}

// MarshalJSON 转换为json类型
func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON json转换为NullString类型
func (v *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.String = *s
	} else {
		v.Valid = false
	}
	return nil
}

// NullInt64 可为nil的数字类型
type NullInt64 struct {
	sql.NullInt64
}

// NewNullInt64 初始化int64类型
func NewNullInt64(i int64, valid bool) NullInt64 {
	var n NullInt64
	n.Valid = valid
	n.Int64 = i
	return n
}

// IsNull 是否为空
func (v NullInt64) IsNull() bool {
	return !v.Valid
}

// SetValide 设值
func (v *NullInt64) SetValide(i int64) {
	v.Int64 = i
	v.Valid = true
}

// SetValideByString 设值
func (v NullInt64) SetValideByString(s string) {
	result, _ := strconv.Atoi(s)
	v.Int64 = int64(result)
	v.Valid = true
}

// MarshalJSON 转换为json类型
func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON json转换为NullString类型
func (v *NullInt64) UnmarshalJSON(data []byte) error {
	var s *int64
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Int64 = *s
	} else {
		v.Valid = false
	}
	return nil
}
