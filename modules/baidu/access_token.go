package baidu

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
)

// OpenAPITokenURL OPEN API token获取地址
const OpenAPITokenURL = "https://openapi.baidu.com/oauth/2.0/token"

// APIRequestTimeout 请求超时时间
var APIRequestTimeout = 120

// baiduErrorResult 百度返回的 错误结果
type baiduErrorResult struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// ErrorResult 错误信息
type ErrorResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// GetAccessToken 获取token
func GetAccessToken(code, appID, secret, redirectURL string) (*AccessToken, *ErrorResult) {
	if code == "" || appID == "" || secret == "" || redirectURL == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}
	data, err := util.HTTPGet(buildAPIRequestURL(OpenAPITokenURL, map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     appID,
		"client_secret": secret,
		"redirect_uri":  redirectURL,
	}), APIRequestTimeout)

	if err != nil {
		log.Errorf("获取百度 access_token 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取百度 access_token 失败"}
	}

	if err := checkAPIResultError(data); err != nil {
		log.Errorf("获取百度 access_token 接口错误：%s", string(data))
		return nil, err
	}

	ret := new(AccessToken)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析百度 access_token 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析百度 access_token 响应失败"}
	}

	return ret, nil
}

// APIAccessToken 百度api返回结果
type APIAccessToken struct {
	AccessToken string `json:"access_token"` // 获取到的网页授权接口调用凭证
	ExpiresIn   int    `json:"expires_in"`   // 凭证有效时间，单位：秒
	Scope       string `json:"scope"`        // 用户授权权限，如：snsapi_userinfo，多个空格分隔
}

// GetAPIAccessToken 直接获取token
func GetAPIAccessToken(appID, secret, scope string) (*APIAccessToken, *ErrorResult) {
	if appID == "" || secret == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}
	data, err := util.HTTPGet(buildAPIRequestURL("https://openapi.baidu.com/oauth/2.0/token", map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     appID,
		"client_secret": secret,
		"scope":         scope,
	}), APIRequestTimeout)

	if err != nil {
		log.Errorf("获取百度 access_token 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取百度 access_token 失败"}
	}

	if err := checkAPIResultError(data); err != nil {
		log.Errorf("获取百度 access_token 接口错误：%s", string(data))
		return nil, err
	}

	ret := new(APIAccessToken)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析百度 access_token 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析百度 access_token 响应失败"}
	}

	return ret, nil
}

// RefreshAccessToken 刷新token
func RefreshAccessToken(refreshToken, appID, secret string) (*AccessToken, *ErrorResult) {
	if refreshToken == "" || appID == "" || secret == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}

	data, err := util.HTTPGet(buildAPIRequestURL(OpenAPITokenURL, map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     appID,
		"client_secret": secret,
	}), APIRequestTimeout)

	if err != nil {
		log.Errorf("刷新百度 access_token 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "刷新百度 access_token 失败"}
	}

	if err := checkAPIResultError(data); err != nil {
		log.Errorf("刷新百度 access_token 接口错误：%s", string(data))
		return nil, err
	}

	ret := new(AccessToken)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("刷新百度 access_token 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "刷新百度 access_token 响应失败"}
	}

	return ret, nil
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

// checkAPIResultError 检查接口是否返回了错误结果
func checkAPIResultError(data []byte) *ErrorResult {
	ret := new(baiduErrorResult)

	if err := json.Unmarshal(data, ret); err != nil {
		return &ErrorResult{Message: ret.ErrorDescription}
	}

	if ret.Error != "" {
		return &ErrorResult{Message: ret.ErrorDescription}
	}

	return nil
}
