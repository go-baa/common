package wechat

import (
	"encoding/json"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/go-baa/log"
)

const (
	// GatewayGetAccessToken ...
	GatewayGetAccessToken = "https://api.weixin.qq.com/sns/oauth2/access_token"
	// GatewayRefreshAccessToken ...
	GatewayRefreshAccessToken = "https://api.weixin.qq.com/sns/oauth2/refresh_token"
	// GetwapAccessTokenURL 获取网页使用的token路径
	GetwapAccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"
)

// AccessTokenResult ...
type AccessTokenResult struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}

// GetAccessToken 获取 access_token
func GetAccessToken(appID, secret, code, grantType string) (*AccessTokenResult, *ErrorResult) {
	data, err := util.HTTPGet(buildAPIRequestURL(GatewayGetAccessToken, map[string]string{
		"appid":      appID,
		"secret":     secret,
		"code":       code,
		"grant_type": grantType,
	}), APIRequestTimeout)

	if err != nil {
		log.Errorf("获取微信 access_token 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取微信 access_token 失败"}
	}

	if err := checkAPIResultError(data); err != nil {
		log.Errorf("获取微信 access_token 接口错误：%s", string(data))
		return nil, err
	}

	ret := new(AccessTokenResult)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析微信 access_token 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析微信 access_token 响应失败"}
	}

	return ret, nil
}

// RefreshAccessToken 刷新 access_token
func RefreshAccessToken(appID, grantType, refreshToken string) (*AccessTokenResult, *ErrorResult) {
	data, err := util.HTTPGet(buildAPIRequestURL(GatewayRefreshAccessToken, map[string]string{
		"appid":         appID,
		"grant_type":    grantType,
		"refresh_token": refreshToken,
	}), APIRequestTimeout)

	if err != nil {
		log.Errorf("刷新微信 access_token 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "刷新微信 access_token 失败"}
	}

	if err := checkAPIResultError(data); err != nil {
		log.Errorf("刷新微信 access_token 接口错误：%s", string(data))
		return nil, err
	}

	ret := new(AccessTokenResult)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析微信 access_token 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析微信 access_token 响应失败"}
	}

	return ret, nil
}

// WapAccessTokenResult 返回值
type WapAccessTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetWapAccessToken 获取网页使用的access_token
func GetWapAccessToken(appID, appSecret string) (*WapAccessTokenResult, *ErrorResult) {
	data, err := util.HTTPGet(buildAPIRequestURL(GetwapAccessTokenURL, map[string]string{
		"appid":      appID,
		"secret":     appSecret,
		"grant_type": "client_credential",
	}), APIRequestTimeout)

	if err != nil {
		log.Errorf("获取微信 网页 access_token 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取微信 网页 access_token 失败"}
	}

	if err := checkAPIResultError(data); err != nil {
		log.Errorf("获取微信 网页 access_token 接口错误：%s", string(data))
		return nil, err
	}

	ret := new(WapAccessTokenResult)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析微信 网页 access_token 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析微信 网页 access_token 响应失败"}
	}

	return ret, nil
}
