package baidu

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

// OpenAPIOathGateway Open API 授权地址
const OpenAPIOathGateway = "https://openapi.baidu.com/oauth/2.0/token"

// Baidu 百度OpenApi实例
type Baidu struct {
	GrantType    string
	ClientID     string
	ClientSecret string
}

// New 实例化百度OpenAPI
func New(clientID, clientSecret string) (*Baidu, error) {
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("client_id or client_secret is empty")
	}

	return &Baidu{
		GrantType:    "client_credentials",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, nil
}

// AccessToken 授权秘钥
type AccessToken struct {
	AccessToken      string `json:"access_token"`
	SessionKey       string `json:"session_key"`
	Scope            string `json:"scope"`
	RefreshToken     string `json:"refresh_token"`
	SessionSecret    string `json:"session_secret"`
	ExpiresIn        int64  `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GetAccessToken 获取授权秘钥
func (t *Baidu) GetAccessToken() (*AccessToken, error) {
	params := map[string]interface{}{
		"grant_type":    t.GrantType,
		"client_id":     t.ClientID,
		"client_secret": t.ClientSecret,
	}

	query := util.HTTPBuildQuery(params, 0)
	res, err := t.request(OpenAPIOathGateway, query, []byte{}, map[string]string{})
	if err != nil {
		return nil, err
	}

	token := new(AccessToken)
	err = json.Unmarshal(res, token)
	if err != nil {
		return nil, fmt.Errorf("JSON解码错误: %v", err)
	}

	if token.AccessToken == "" {
		return nil, fmt.Errorf("获取access_token错误:%v", token.ErrorDescription)
	}

	return token, nil
}

// request http请求
func (t *Baidu) request(url string, query string, reqBody []byte, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url+"?"+query, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	// 超时设置
	client := new(http.Client)
	client.Timeout = time.Second * 60

	// https 支持
	if strings.HasPrefix(url, "https") {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	// 执行请求
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 处理响应
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	if setting.Debug {
		log.Debugf("Baidu debug: url: %v, res: %v", url+"?"+query, string(body))
	}

	return body, nil
}
