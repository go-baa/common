package baidu

import (
	"encoding/json"
	"fmt"

	"github.com/go-baa/common/util"
)

const sessionKeyURL = "https://spapi.baidu.com/oauth/jscode2sessionkey"

// GetSessionKey 获取session key
func GetSessionKey(code, appKey, secret string) (*SessionkeyInfo, error) {
	param := make(map[string]string)
	param["code"] = code
	param["client_id"] = appKey
	param["sk"] = secret
	bs, err := util.HTTPPost(sessionKeyURL, param, 10, nil)
	if err != nil {
		return nil, err
	}
	baiduResp := &struct {
		OpenID           string `json:"openid"`
		Sessionkey       string `json:"session_key"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}{}
	err = json.Unmarshal(bs, baiduResp)
	if err != nil {
		return nil, err
	}
	if baiduResp.Error != "" {
		return nil, fmt.Errorf("code:%s ,error:%s", baiduResp.Error, baiduResp.ErrorDescription)

	}
	return &SessionkeyInfo{
		OpenID:     baiduResp.OpenID,
		SessionKey: baiduResp.Sessionkey,
	}, nil
}

// SessionkeyInfo .
type SessionkeyInfo struct {
	OpenID     string
	SessionKey string
}
