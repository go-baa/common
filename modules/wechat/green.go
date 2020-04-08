package wechat

import (
	"encoding/json"
	"fmt"

	"github.com/go-baa/common/util"
)

const (
	host   = "https://api.weixin.qq.com/wxa/msg_sec_check?access_token=%s"
	method = "POST"
)

// GreenTextScanner .
type GreenTextScanner struct{}

type req struct {
	Tasks []task `json:"tasks"`
}
type task struct {
	Content string `json:"content"`
}

// GreenScanResponse .
type GreenScanResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// GreenTextScran 垃圾文本检测
func (t *GreenTextScanner) GreenTextScran(content, accessToken string) (*GreenScanResponse, error) {
	tk := task{Content: content}
	respBs, err := util.HTTPPostJSON(fmt.Sprintf(host, accessToken), &tk, 10)
	if err != nil {
		return nil, err
	}
	resp := GreenScanResponse{}
	err = json.Unmarshal(respBs, &resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode > 0 && resp.ErrCode != 87014 {
		return nil, fmt.Errorf("%d %s", resp.ErrCode, resp.ErrMsg)
	}
	return &resp, nil
}

// NewGreenTextScanner .
func NewGreenTextScanner() (*GreenTextScanner, error) {
	return &GreenTextScanner{}, nil
}
