package util

import "strconv"

// MustInt 转换值为 int 类型
func MustInt(v interface{}) int {
	switch v.(type) {
	case int:
		return v.(int)
	case int8:
		return int(v.(int8))
	case int16:
		return int(v.(int16))
	case int32:
		return int(v.(int32))
	case int64:
		return int(v.(int64))
	case float32:
		return int(v.(float32))
	case float64:
		return int(v.(float64))
	case string:
		val, _ := strconv.Atoi(v.(string))
		return val
	}
	return 0
}

// MustInt64 转换值为 int64 类型
func MustInt64(v interface{}) int64 {
	switch v.(type) {
	case int:
		return int64(v.(int))
	case int8:
		return int64(v.(int8))
	case int16:
		return int64(v.(int16))
	case int32:
		return int64(v.(int32))
	case int64:
		return v.(int64)
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	case string:
		val, _ := strconv.ParseInt(v.(string), 10, 64)
		return val
	}
	return 0
}

// MustFloat64 转换值为 float64 类型
func MustFloat64(v interface{}) float64 {
	switch v.(type) {
	case int:
		return float64(v.(int))
	case int8:
		return float64(v.(int8))
	case int16:
		return float64(v.(int16))
	case int32:
		return float64(v.(int32))
	case int64:
		return float64(v.(int64))
	case float32:
		return float64(v.(float32))
	case float64:
		return v.(float64)
	case string:
		val, _ := strconv.ParseFloat(v.(string), 64)
		return val
	}
	return 0
}