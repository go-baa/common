package toutiao

import (
	"encoding/json"
	"fmt"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
)

const (
	// GatSnsGetUserInfo ...
	GatSnsGetUserInfo = "https://developer.toutiao.com/api/apps/jscode2session?appid=%s&secret=%s&code=%s&anonymous_code=%s"
)

// APIRequestTimeout 请求超时时间
var APIRequestTimeout = 10

// SnsLoginInfo 小程序登录后 信息
type SnsLoginInfo struct {
	OpenID          string `json:"openid"`
	SessionKey      string `json:"session_key"`
	AnonymousOpenID string `json:"anonymous_openid"`
}

// ErrorResult 错误信息
type ErrorResult struct {
	Error   int    `json:"error"`
	Message string `json:"errmsg"`
}

// SnsLogin 获取用户信息 code anonymousCode 必须有一个
func SnsLogin(appID, appSecret, code, anonymousCode string) (*SnsLoginInfo, *ErrorResult) {
	fmt.Println(fmt.Sprintf(GatSnsGetUserInfo, appID, appSecret, code, anonymousCode))
	data, err := util.HTTPGet(fmt.Sprintf(GatSnsGetUserInfo, appID, appSecret, code, anonymousCode), APIRequestTimeout)
	if err != nil {
		log.Errorf("头条小程序登录 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "头条小程序登录 失败"}
	}

	if err := checkAPIResultError(data); err != nil {
		log.Errorf("头条小程序登录 错误：%s", string(data))
		return nil, err
	}

	ret := new(SnsLoginInfo)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析头条小程序登录信息 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析头条小程序登录信息 响应失败"}
	}
	return ret, nil
}

// checkAPIResultError 检查接口是否返回了错误结果
func checkAPIResultError(data []byte) *ErrorResult {
	ret := new(ErrorResult)

	if err := json.Unmarshal(data, ret); err != nil {
		return ret
	}

	if ret.Error > 0 {
		return ret
	}

	return nil
}
