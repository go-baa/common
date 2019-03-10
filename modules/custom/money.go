package custom

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Money 金额自定义结构体
type Money struct {
	MoneyInt64 int64
	MoneyFloat float64
	MoneyStr   string
	Valide     bool
}

// formatMoney
func (t Money) formatMoney(money int64) string {
	switch money {
	case 0:
		return "0.00"
	default:
		return fmt.Sprintf("%.2f", float64(money)/100)
	}
}

// NewMoney 创建
func NewMoney(money int64) Money {
	return Money{
		MoneyInt64: money,
		MoneyFloat: float64(money) / 100,
		MoneyStr:   Money{}.formatMoney(money),
		Valide:     true,
	}
}

// IsNull 判断是否为空
func (t Money) IsNull() bool {
	return !t.Valide
}

// MarshalJSON implements json.Marshaler.
func (t Money) MarshalJSON() ([]byte, error) {
	if t.IsNull() {
		return []byte("\"0.00\""), nil
	}
	return []byte("\"" + t.MoneyStr + "\""), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Money) UnmarshalJSON(data []byte) (err error) {
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
func (t *Money) SetInt64(v int64) {
	t.MoneyInt64 = v
	t.MoneyFloat = float64(v) / 100
	t.MoneyStr = Money{}.formatMoney(v)
	t.Valide = true
}

// SetString 设置金额（元）
func (t *Money) SetString(v string) (err error) {
	t.Valide = true
	if v == "" {
		t.MoneyStr = "0.00"
		t.MoneyInt64 = 0
		t.MoneyFloat = 0.0
		return
	}
	t.MoneyFloat, err = strconv.ParseFloat(v, 64)
	t.MoneyInt64 = int64(math.Floor(t.MoneyFloat*100 + 0.5))
	t.MoneyStr = Money{}.formatMoney(t.MoneyInt64)
	if strings.Contains(t.MoneyStr, ".") {
		nums := strings.Split(t.MoneyStr, ".")
		if len(nums[1]) > 2 {
			return ErrMoneyFmt
		}
	}
	return
}

// SetFloat64 设置金额（元）
func (t *Money) SetFloat64(v float64) (err error) {
	t.MoneyFloat = v
	t.MoneyInt64 = int64(v * 100)
	t.MoneyStr = Money{}.formatMoney(t.MoneyInt64)
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
func (t *Money) AddMoney(v Money) {
	t.MoneyInt64 += v.MoneyInt64
	t.MoneyFloat = float64(t.MoneyInt64) / 100
	t.MoneyStr = Money{}.formatMoney(t.MoneyInt64)
}

// SubMoney 金额相减
func (t *Money) SubMoney(v Money) {
	t.MoneyInt64 -= v.MoneyInt64
	t.MoneyFloat = float64(t.MoneyInt64) / 100
	t.MoneyStr = Money{}.formatMoney(t.MoneyInt64)
}

// Scan implements the Scanner interface.
func (t *Money) Scan(src interface{}) error {
	var err error
	switch x := src.(type) {
	case int64:
		t.MoneyInt64 = x
		t.MoneyFloat = float64(t.MoneyInt64) / 100
		t.MoneyStr = Money{}.formatMoney(x)
	case nil:
		t.MoneyInt64 = 0
		t.MoneyFloat = 0.0
		t.MoneyStr = "0.00"
	case []byte:
		val := string(x)
		moneyInt64, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			err = fmt.Errorf("ParseInt error: %v", src)
		}
		t.MoneyInt64 = moneyInt64
		t.MoneyFloat = float64(t.MoneyInt64) / 100
		t.MoneyStr = Money{}.formatMoney(moneyInt64)
	default:
		err = fmt.Errorf("null: cannot scan type %T into Money: %v", src, src)
	}
	t.Valide = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (t Money) Value() (driver.Value, error) {
	if t.IsNull() {
		return int64(0), nil
	}
	return int64(t.MoneyInt64), nil
}

// String toString方法
func (t Money) String() string {
	if t.IsNull() {
		return "0.00"
	}
	return t.MoneyStr
}

// CnyString 大写输出
func (t Money) CnyString() string {
	num := t.MoneyFloat
	strnum := strconv.FormatFloat(num*100, 'f', 0, 64)
	sliceUnit := []string{"仟", "佰", "拾", "亿", "仟", "佰", "拾", "万", "仟", "佰", "拾", "元", "角", "分"}
	// log.Println(sliceUnit[:len(sliceUnit)-2])
	s := sliceUnit[len(sliceUnit)-len(strnum) : len(sliceUnit)]
	upperDigitUnit := map[string]string{"0": "零", "1": "壹", "2": "贰", "3": "叁", "4": "肆", "5": "伍", "6": "陆", "7": "柒", "8": "捌", "9": "玖"}
	str := ""
	for k, v := range strnum[:] {
		str = str + upperDigitUnit[string(v)] + s[k]
	}
	reg, err := regexp.Compile(`零角零分$`)
	str = reg.ReplaceAllString(str, "整")

	reg, err = regexp.Compile(`零角`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零分$`)
	str = reg.ReplaceAllString(str, "整")

	reg, err = regexp.Compile(`零[仟佰拾]`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零{2,}`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零亿`)
	str = reg.ReplaceAllString(str, "亿")

	reg, err = regexp.Compile(`零万`)
	str = reg.ReplaceAllString(str, "万")

	reg, err = regexp.Compile(`零*元`)
	str = reg.ReplaceAllString(str, "元")

	reg, err = regexp.Compile(`亿零{0, 3}万`)
	str = reg.ReplaceAllString(str, "^元")

	reg, err = regexp.Compile(`零元`)
	str = reg.ReplaceAllString(str, "零")
	if err != nil {
		log.Fatal(err)
	}
	return str
}

// FormatMoney 格式化为两位小数字符串，若无小数则舍去小数点后表示
func FormatMoney(money int) string {
	return Money{}.formatMoney(int64(money))
}
