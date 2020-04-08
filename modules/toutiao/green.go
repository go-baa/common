package toutiao

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	host           = "https://developer.toutiao.com/api/v2/tags/text/antidirt"
	method         = "POST"
	tokenGrantType = "client_credential"
)

// GreenTextScanner .
type GreenTextScanner struct {
	appID     string
	appSecret string
}
type req struct {
	Tasks []task `json:"tasks"`
}
type task struct {
	Content string `json:"content"`
}

// GreenScanResponse .
type GreenScanResponse struct {
	LogID string `json:"log_id"`
	Data  []struct {
		Code     int    `json:"code"`
		TaskID   string `json:"task_id"`
		Predicts []struct {
			Prob float32 `json:"prob"`
		} `json:"predicts"`
	} `json:"data"`
}

// GreenTextScran 垃圾文本检测
func (t *GreenTextScanner) GreenTextScran(content string) (*GreenScanResponse, error) {
	reqJSON := req{Tasks: []task{task{Content: content}}}
	reqBody, err := json.Marshal(&reqJSON)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, host, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	// token
	tk, terr := GetAccessToken(t.appID, t.appSecret, tokenGrantType)
	if err != nil {
		return nil, fmt.Errorf("get toutiao access token failed, code:%d message:%s", terr.Error, terr.Message)
	}

	req.Header = map[string][]string{"X-Token": []string{tk.AccessToken}}
	// 超时设置
	client := new(http.Client)
	client.Timeout = time.Second * 10

	// https 支持
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// 执行请求
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 处理响应
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	resp := GreenScanResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// NewGreenTextScanner .
func NewGreenTextScanner(appID, appSecret string) (*GreenTextScanner, error) {
	if appID == "" {
		return nil, fmt.Errorf("Invalid GreenTextScanner appID: %s", appID)
	}

	if appSecret == "" {
		return nil, fmt.Errorf("Invalid GreenTextScanner appSecret: %s", appSecret)
	}
	scanner := &GreenTextScanner{
		appID:     appID,
		appSecret: appSecret,
	}
	return scanner, nil
}
