package alipay

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/go-baa/log"
)

// UserInfo 小程序登录后 信息
type UserInfo struct {
	UserInfoResponse userInfoResponse `json:"alipay_user_info_share_response"`
	ErrorResponse    errorResponse    `json:"error_response"`
	Sign             string           `json:"sign"`
}

// SnsUserInfo .
func SnsUserInfo(appID, pk, authToken string) (*UserInfo, error) {
	params := url.Values{}
	params.Add("app_id", appID)
	params.Add("auth_token", authToken)
	params.Add("charset", "utf-8")
	params.Add("method", userInfoAPI)
	params.Add("sign_type", "RSA2")
	params.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Add("version", "1.0")
	sign := sign(sortKeys(params), pk)
	params.Add("sign", sign)
	client := &http.Client{
		Timeout: time.Second * time.Duration(APIRequestTimeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.PostForm(alipayGateway+"?"+params.Encode(), params)
	if err != nil {
		log.Errorf("请求支付宝小程序接口 失败：%s", err.Error())
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("支付宝小程序登录 失败：%s", err.Error())
		return nil, err
	}
	ret := new(UserInfo)
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	if ret.ErrorResponse.Code != "" {
		return nil, fmt.Errorf("支付宝小程序登录 支付宝API响应错误: %s %s", ret.ErrorResponse.Code, ret.ErrorResponse.Msg)
	}

	if ret.UserInfoResponse.Code != "" {
		return nil, fmt.Errorf("支付宝小程序登录 支付宝API响应错误: %s %s %s %s", ret.UserInfoResponse.Code, ret.UserInfoResponse.Msg,
			ret.UserInfoResponse.SubCode, ret.UserInfoResponse.SubMsg)
	}

	return ret, nil
}
