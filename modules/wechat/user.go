package wechat

import (
	"encoding/json"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/go-baa/log"
)

const (
	// GatewayGetUserInfo ...
	GatewayGetUserInfo = "https://api.weixin.qq.com/sns/userinfo"
)

// UserInfo 用户信息
type UserInfo struct {
	OpenID       string   `json:"openid"`
	Nickname     string   `json:"nickname"`
	Sex          int      `json:"sex"`
	Province     string   `json:"provinc"`
	City         string   `json:"city"`
	Country      string   `json:"country"`
	HeadImageURL string   `json:"headimgurl"`
	Privilege    []string `json:"privileg"`
	UnionID      string   `json:"unionid"`
}

// GetUserInfo 获取用户信息
func GetUserInfo(accessToken, openID, lang string) (*UserInfo, *ErrorResult) {
	data, err := util.HTTPGet(buildAPIRequestURL(GatewayGetUserInfo, map[string]string{
		"access_token": accessToken,
		"openid":       openID,
		"lang":         lang,
	}), APIRequestTimeout)

	if err != nil {
		log.Errorf("获取微信用户信息失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取微信用户信息失败"}
	}

	if err := checkAPIResultError(data); err != nil {
		log.Errorf("微信用户信息接口错误：%s", string(data))
		return nil, err
	}

	ret := new(UserInfo)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析微信用户信息响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析微信用户信息响应失败"}
	}

	return ret, nil
}
