package im

import (
	"encoding/json"
	"fmt"
)

// ServiceConfig 全局配置服务名
const ServiceConfig = "openconfigsvr"

// SetNoSpeakingRequest 设置全局禁言请求
type SetNoSpeakingRequest struct {
	SetAccount             string `json:"Set_Account"`
	C2CmsgNospeakingTime   int64  `json:",omitempty"` // 单聊消息禁言时间，秒为单位，非负整数。等于0代表没有被设置禁言；
	GroupmsgNospeakingTime int64  `json:",omitempty"` // 群组消息禁言时间，秒为单位，非负整数。等着0代表没有被设置禁言；
}

// GetNoSpeakingRequest 查询全局禁言配置请求
type GetNoSpeakingRequest struct {
	GetAccount string `json:"Get_Account"`
}

// NoSpeakingConfig 全局禁言配置项
type NoSpeakingConfig struct {
	C2CmsgNospeakingTime   int64
	GroupmsgNospeakingTime int64
}

// NoSpeakingResponse 查询全局禁言响应
type NoSpeakingResponse struct {
	Response
	NoSpeakingConfig
}

// SetNoSpeaking 设置全局禁言
func (t *IM) SetNoSpeaking(account string, c2cMsgLimit, groupMsgLimit int64) error {
	req := &SetNoSpeakingRequest{
		SetAccount:             account,
		C2CmsgNospeakingTime:   c2cMsgLimit,
		GroupmsgNospeakingTime: groupMsgLimit,
	}
	res, err := t.api(ServiceConfig, "setnospeaking", req)
	if err != nil {
		return err
	}

	response := new(Response)
	err = json.Unmarshal(res, response)
	if err != nil {
		return fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return nil
}

// GetNoSpeaking 查询全局禁言
func (t *IM) GetNoSpeaking(account string) (*NoSpeakingConfig, error) {
	req := &GetNoSpeakingRequest{GetAccount: account}
	res, err := t.api(ServiceConfig, "getnospeaking", req)
	if err != nil {
		return nil, err
	}

	response := new(NoSpeakingResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return &response.NoSpeakingConfig, nil
}
