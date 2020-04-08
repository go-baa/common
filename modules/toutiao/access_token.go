package toutiao

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
)

const (
	// GatewayGetAccessToken ...
	GatewayGetAccessToken = "https://developer.toutiao.com/api/apps/token"
)

// AccessTokenResult ...
type AccessTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// AccessTokenErrorResult ...
type AccessTokenErrorResult struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

// GetAccessToken 获取 access_token
func GetAccessToken(appID, secret, grantType string) (*AccessTokenResult, *AccessTokenErrorResult) {
	data, err := HTTPGet(fmt.Sprintf("%s?%s", GatewayGetAccessToken, util.HTTPBuildQuery(map[string]interface{}{
		"appid":      appID,
		"secret":     secret,
		"grant_type": grantType,
	}, 0)), 3)

	if err != nil {
		log.Errorf("toutiao.GetAccessToken error: %s", err.Error())
		return nil, &AccessTokenErrorResult{Error: 1, Message: err.Error()}
	}

	var apiError = new(AccessTokenErrorResult)
	if err := json.Unmarshal(data, &apiError); err == nil {
		if apiError.Error > 0 {
			return nil, apiError
		}
	}

	var ret = new(AccessTokenResult)
	if err := json.Unmarshal(data, &ret); err != nil {
		log.Errorf("toutiao.GetAccessToken error: parse reponse %s", err.Error())
		return nil, &AccessTokenErrorResult{Error: 1, Message: "parse reponse " + err.Error()}
	}

	return ret, nil
}

// HTTPGet 带超时设置的请求一个url，单位: 秒
func HTTPGet(uri string, timeout int) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request err status code is %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
