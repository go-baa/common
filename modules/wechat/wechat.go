package wechat

import (
	"encoding/json"
	"net/url"
	"strings"
)

const (
	// GrantTypeAuthorizationCode ...
	GrantTypeAuthorizationCode = "authorization_code"
)

const (
	LangCN = "zh_CN"
	LangTW = "zh_TW"
	LangEN = "en"
)

// APIRequestTimeout 请求超时时间
var APIRequestTimeout = 10

// ErrorResult 错误结果
type ErrorResult struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

// checkAPIResultError 检查接口是否返回了错误结果
func checkAPIResultError(data []byte) *ErrorResult {
	ret := new(ErrorResult)

	if err := json.Unmarshal(data, ret); err != nil {
		return ret
	}

	if ret.Code > 0 {
		return ret
	}

	return nil
}

// buildAPIRequestURL 构建请求参数
func buildAPIRequestURL(gateway string, params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	query := values.Encode()
	if len(query) > 0 {
		if strings.Contains(gateway, "?") {
			return gateway + "&" + query
		}
		return gateway + "?" + query
	}

	return gateway
}
