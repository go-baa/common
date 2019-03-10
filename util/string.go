package util

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	STR_PAD_LEFT int = iota
	STR_PAD_RIGHT
	STR_PAD_BOTH
)

// ToString 将任意一个类型转换为字符串
func ToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// StringToInt 将字符串转为int
func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// StringToInt64 将字符串转为int64
func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// IntToString 将数字转换为字符串
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// CamelCase 将一个字符串转为大驼峰命名
func CamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	ss := strings.Split(s, " ")
	for k, v := range ss {
		ss[k] = strings.Title(v)
	}
	return strings.Join(ss, "")
}

// Link: https://github.com/golang/lint/blob/master/lint.go
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

// CamelCaseInitialism 将一个字符串转为大驼峰命名，强制首字母缩写命名规范
func CamelCaseInitialism(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	ss := strings.Split(s, " ")
	var sm bool
	var uv string
	for k, v := range ss {
		sm = false
		uv = strings.ToUpper(v)
		for m := range commonInitialisms {
			if uv == m {
				ss[k] = m
				sm = true
				break
			}
		}
		if !sm {
			ss[k] = strings.Title(v)
		}
	}
	return strings.Join(ss, "")
}

// IsNumeric 判断是否是纯数字，空字符串不是数字，0才是
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		if v < 48 || v > 57 {
			return false
		}
	}
	return true
}

// StrPad 使用另一个字符串填充字符串为指定长度
func StrPad(v interface{}, length int, pad string, padType int) string {
	s := fmt.Sprintf("%v", v)
	if len(s) >= length {
		return s
	}
	var pos int
	switch padType {
	case STR_PAD_BOTH:
		left := true
		for len(s) < length {
			// first left, then right
			if left {
				s = pad + s
				if len(s) > length {
					pos = len(s) - length
					return s[pos:]
				}
				left = false
			} else {
				s = s + pad
				if len(s) > length {
					return s[:length]
				}
				left = true
			}
		}
	case STR_PAD_RIGHT:
		for len(s) < length {
			s = s + pad
			if len(s) > length {
				return s[:length]
			}
		}
	default:
		for len(s) < length {
			s = pad + s
			if len(s) > length {
				pos = len(s) - length
				return s[pos:]
			}
		}

	}

	return s
}

// StrNatCut 字符串截取，中文算一个 英文算两个
func StrNatCut(s string, length int, dots ...string) string {
	source := []rune(s)
	n := len(source)
	if n <= length {
		return s
	}

	if dots == nil {
		dots = []string{"..."}
	}

	dst := make([]rune, length)
	copy(dst, source)
	return fmt.Sprintf("%s%s", string(dst), dots[0])
}

// PurgeText 纯净的文本，不带html标签，没有换行，没有制表符
func PurgeText(s string) string {
	s = string(StripTags([]byte(s), ""))
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "&nbsp;", "", -1)
	s = strings.Replace(s, "&quot;", "", -1)
	s = strings.Replace(s, "&ldquo;", "", -1)
	s = strings.Replace(s, "&rdquo;", "", -1)
	s = strings.Replace(s, "&bdquo;", "", -1)
	s = strings.Replace(s, "&lsquo;", "", -1)
	s = strings.Replace(s, "&middot;", "", -1)
	s = strings.Replace(s, "&hellip;", "", -1)
	s = strings.Replace(s, "&tilde;", "", -1)
	s = strings.Replace(s, "&mdash;", "", -1)
	s = strings.Replace(s, "&ndash;", "", -1)
	return s
}

// Concat 连接字符串
func Concat(args ...string) string {
	var buf bytes.Buffer
	for i := range args {
		buf.WriteString(args[i])
	}
	return buf.String()
}

// MaskEmail 使用指定字符遮罩邮箱
func MaskEmail(email string, masks ...string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return MaskString(email, masks...)
	}
	return MaskString(parts[0], masks...) + "@" + parts[1]
}

// MaskMobile 使用指定字符遮罩手机号
func MaskMobile(mobile string, masks ...string) string {
	if len(mobile) != 11 {
		return MaskString(mobile, masks...)
	}
	mask := "*"
	if len(masks) > 0 {
		mask = masks[0]
	}
	return mobile[:3] + strings.Repeat(mask, 4) + mobile[7:]
}

// MaskString 使用指定字符遮罩字符串
func MaskString(str string, masks ...string) string {
	parts := []rune(str)

	length := len(parts)
	if length <= 2 {
		return str
	}

	n := (length - length%3) / 3
	m := length - n*2

	mask := "*"
	if len(masks) > 0 {
		mask = masks[0]
	}

	return string(parts[:n]) + strings.Repeat(mask, m) + string(parts[n+m:])
}

// SplitStringToSlice 将字符串分割成数组，并去除空行
func SplitStringToSlice(s string, sep string) []string {
	if len(s) == 0 || len(sep) == 0 {
		return nil
	}

	ss := strings.Split(s, sep)
	if len(ss) == 0 {
		return nil
	}

	items := make([]string, 0, len(ss))
	for i := range ss {
		ss[i] = strings.TrimSpace(ss[i])
		if len(ss[i]) > 0 {
			items = append(items, ss[i])
		}
	}

	return items
}

// SplitStringToIntSlice 将字符串分割成数字数组，并去除空行
func SplitStringToIntSlice(s string, sep string) []int {
	if len(s) == 0 || len(sep) == 0 {
		return nil
	}

	ss := strings.Split(s, sep)
	if len(ss) == 0 {
		return nil
	}

	items := make([]int, 0, len(ss))
	for i := range ss {
		ss[i] = strings.TrimSpace(ss[i])
		if len(ss[i]) > 0 {
			if n, err := strconv.Atoi(ss[i]); err == nil {
				items = append(items, n)
			}
		}
	}

	return items
}
