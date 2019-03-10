package wechat

import (
	"encoding/json"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
)

const (
	// GatSnsGetUserInfo ...
	GatSnsGetUserInfo = "https://api.weixin.qq.com/sns/jscode2session"
)

// SnsLoginInfo 小程序登录后 信息
type SnsLoginInfo struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
}

// SnsLogin 获取用户信息
func SnsLogin(appID, appSecret, code string) (*SnsLoginInfo, *ErrorResult) {
	data, err := util.HTTPGet(buildAPIRequestURL(GatSnsGetUserInfo, map[string]string{
		"appid":      appID,
		"secret":     appSecret,
		"js_code":    code,
		"grant_type": "authorization_code",
	}), APIRequestTimeout)

	if err != nil {
		log.Errorf("微信小程序登录 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "微信小程序登录 失败"}
	}

	if err := checkAPIResultError(data); err != nil {
		log.Errorf("微信小程序登录 错误：%s", string(data))
		return nil, err
	}

	ret := new(SnsLoginInfo)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析微信小程序登录信息 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析微信小程序登录信息 响应失败"}
	}

	return ret, nil
}
