package custom

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ErrMoneyFmt 金额格式化错误
var ErrMoneyFmt = errors.New("金钱相关字段只支持两位小数，最小单位精确到分")

// SimpleMoney 简写金额自定义结构体，18.00元显示为18元
type SimpleMoney struct {
	MoneyInt64 int64
	MoneyFloat float64
	MoneyStr   string
	Valide     bool
}

// formatMoney
func (t SimpleMoney) formatMoney(money int64) string {
	switch money {
	case 0:
		return "0"
	case (money / 100) * 100:
		return fmt.Sprintf("%.0f", float64(money)/100)
	case (money / 10) * 10:
		return fmt.Sprintf("%.1f", float64(money)/100)
	default:
		return fmt.Sprintf("%.2f", float64(money)/100)
	}
}

// NewSimpleMoney 创建
func NewSimpleMoney(money int64) SimpleMoney {
	return SimpleMoney{
		MoneyInt64: money,
		MoneyFloat: float64(money) / 100,
		MoneyStr:   SimpleMoney{}.formatMoney(money),
		Valide:     true,
	}
}

// IsNull 判断是否为空
func (t SimpleMoney) IsNull() bool {
	return !t.Valide
}

// MarshalJSON implements json.Marshaler.
func (t SimpleMoney) MarshalJSON() ([]byte, error) {
	if t.IsNull() {
		return []byte("\"0\""), nil
	}
	return []byte("\"" + t.MoneyStr + "\""), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *SimpleMoney) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		t.Valide = false
		return
	}
	if bytes.Equal(data, nil) {
		t.Valide = false
		return
	}

	t.MoneyStr = string(data)
	if strings.Contains(t.MoneyStr, `"`) {
		t.MoneyStr = strings.Replace(t.MoneyStr, `"`, "", -1)
	}
	t.MoneyFloat, err = strconv.ParseFloat(t.MoneyStr, 64)
	// 此处采用四舍五入，避免精度问题
	t.MoneyInt64 = int64(math.Floor(t.MoneyFloat*100 + 0.5))
	if strings.Contains(t.MoneyStr, ".") {
		nums := strings.Split(t.MoneyStr, ".")
		if len(nums[1]) > 2 {
			return ErrMoneyFmt
		}
	}
	t.Valide = true
	return
}

// SetInt64 设置金额（分）
func (t *SimpleMoney) SetInt64(v int64) {
	t.MoneyInt64 = v
	t.MoneyFloat = float64(v) / 100
	t.MoneyStr = SimpleMoney{}.formatMoney(v)
	t.Valide = true
}

// SetString 设置金额（元）
func (t *SimpleMoney) SetString(v string) (err error) {
	t.Valide = true
	if v == "" {
		t.MoneyStr = "0"
		t.MoneyInt64 = 0
		t.MoneyFloat = 0.0
		return
	}
	t.MoneyFloat, err = strconv.ParseFloat(v, 64)
	t.MoneyInt64 = int64(math.Floor(t.MoneyFloat*100 + 0.5))
	t.MoneyStr = SimpleMoney{}.formatMoney(t.MoneyInt64)
	if strings.Contains(t.MoneyStr, ".") {
		nums := strings.Split(t.MoneyStr, ".")
		if len(nums[1]) > 2 {
			return ErrMoneyFmt
		}
	}
	return
}

// SetFloat64 设置金额（元）
func (t *SimpleMoney) SetFloat64(v float64) (err error) {
	t.MoneyFloat = v
	t.MoneyInt64 = int64(v * 100)
	t.MoneyStr = SimpleMoney{}.formatMoney(t.MoneyInt64)
	if strings.Contains(t.MoneyStr, ".") {
		nums := strings.Split(t.MoneyStr, ".")
		if len(nums[1]) > 2 {
			return ErrMoneyFmt
		}
	}
	t.Valide = true
	return
}

// AddMoney 金额相加
func (t *SimpleMoney) AddMoney(v SimpleMoney) {
	t.MoneyInt64 += v.MoneyInt64
	t.MoneyFloat = float64(t.MoneyInt64) / 100
	t.MoneyStr = SimpleMoney{}.formatMoney(t.MoneyInt64)
}

// SubMoney 金额相减
func (t *SimpleMoney) SubMoney(v SimpleMoney) {
	t.MoneyInt64 -= v.MoneyInt64
	t.MoneyFloat = float64(t.MoneyInt64) / 100
	t.MoneyStr = SimpleMoney{}.formatMoney(t.MoneyInt64)
}

// Scan implements the Scanner interface.
func (t *SimpleMoney) Scan(src interface{}) error {
	var err error
	switch x := src.(type) {
	case int64:
		t.MoneyInt64 = x
		t.MoneyFloat = float64(t.MoneyInt64) / 100
		t.MoneyStr = SimpleMoney{}.formatMoney(x)
	case nil:
		t.MoneyInt64 = 0
		t.MoneyFloat = 0.0
		t.MoneyStr = "0"
	case []byte:
		val := string(x)
		moneyInt64, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			err = fmt.Errorf("ParseInt error: %v", src)
		}
		t.MoneyInt64 = moneyInt64
		t.MoneyFloat = float64(t.MoneyInt64) / 100
		t.MoneyStr = SimpleMoney{}.formatMoney(moneyInt64)
	default:
		err = fmt.Errorf("null: cannot scan type %T into SimpleMoney: %v", src, src)
	}
	t.Valide = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (t SimpleMoney) Value() (driver.Value, error) {
	if t.IsNull() {
		return int64(0), nil
	}
	return int64(t.MoneyInt64), nil
}

// String toString方法
func (t SimpleMoney) String() string {
	if t.IsNull() {
		return "0"
	}
	return t.MoneyStr
}
