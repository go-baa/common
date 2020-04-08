package wechat

import (
	"encoding/json"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
)

const (
	// WxacodeUnlimitedURL 无限制小程序码请求地址
	WxacodeUnlimitedURL = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
)

// 获取无限制小程序码请求参数
type WxacodeUnlimitedRequest struct {
	Scene     string `json:"scene"`      // 场景值 *必须
	Page      string `json:"page"`       // 扫码后跳转的页面地址（默认主页）
	Width     int    `json:"width"`      // 二维码宽度（默认430）
	AutoColor bool   `json:"auto_color"` // 自动配置线条颜色（默认false）
	LineColor struct {
		R string `json:"r"`
		G string `json:"g"`
		B string `json:"b"`
	} `json:"line_color"` // 线条颜色RGB（auto_color为false是有效）
	IsHyaline bool `json:"is_hyaline"` // 是否透明底色（默认false）
}

// GetUnlimitedWxacode 获取无限制小程序码
func GetUnlimitedWxacode(accessToken string, req WxacodeUnlimitedRequest) ([]byte, *ErrorResult) {
	uri := buildAPIRequestURL(WxacodeUnlimitedURL, map[string]string{
		"access_token": accessToken,
	})
	data, err := util.HTTPPostJSON(uri, req, APIRequestTimeout)
	if err != nil {
		log.Errorf("获取微信小程序码失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取微信小程序码失败"}
	}
	errRet := &ErrorResult{}
	_ = json.Unmarshal(data, errRet)
	if errRet.Code > 0 {
		return nil, errRet
	}

	return data, nil
}
