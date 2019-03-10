package cps

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PushResponse 推送请求响应
type PushResponse struct {
	Response
	MessageID string `json:"MessageId"`
}

// PushConfig 推送配置
type PushConfig struct {
	Target          string
	TargetValue     []string
	DeviceType      string
	PushType        string
	ExtParameters   map[string]string
	PushTime        time.Time
	StoreOffline    bool
	Env             string
	AndroidOpenType string
}

const (
	// CancelAction 取消定时任务接口名称
	CancelAction = "CancelPush"
	// PushAction 推送接口名称
	PushAction = "Push"
)

// 设备类型
const (
	// DeviceTypeAll 全部设备类型
	DeviceTypeAll = "ALL"
	// DeviceTypeiOS iOS设备
	DeviceTypeiOS = "iOS"
	// DeviceTypeAndroid Android设备
	DeviceTypeAndroid = "ANDROID"
)

// 推送目标
const (
	// TargetAll 全部设备
	TargetAll = "ALL"
	// TargetTag 标签
	TargetTag = "TAG"
	// TargetAlias 别名
	TargetAlias = "ALIAS"
	// TargetAccount 账号
	TargetAccount = "ACCOUNT"
	// TargetDevice 设备
	TargetDevice = "DEVICE"
)

const (
	// TargetValueAll 全部设备目标取值
	TargetValueAll = "all"
)

const (
	// PushTypeMessage 消息类型
	PushTypeMessage = "MESSAGE"
	// PushTypeNotice 通知类型
	PushTypeNotice = "NOTICE"
)

const (
	// IOSApnsEnvDev iOS通知开发环境
	IOSApnsEnvDev = "DEV"
	// IOSApnsEnvProduct iOS通知生产环境
	IOSApnsEnvProduct = "PRODUCT"
)

const (
	// AndroidOpenTypeApplication 打开应用
	AndroidOpenTypeApplication = "APPLICATION"
	// AndroidOpenTypeActivity 打开地址
	AndroidOpenTypeActivity = "ACTIVITY"
)

// CancelPush 取消定时推送任务
func (t *CPS) CancelPush(messageID string) (*Response, error) {
	params := t.getCommonParams()
	params["Action"] = CancelAction
	params["MessageId"] = messageID

	// 签名
	sign := t.sign(http.MethodGet, params)
	params["Signature"] = sign

	// 获取响应
	query := url.Values{}
	for k := range params {
		query.Add(k, params[k])
	}
	res, err := t.request(CPSGateway, query)
	if err != nil {
		return nil, err
	}

	// 解析响应
	ret := new(Response)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, err
	}

	if ret.Code != "" {
		return nil, errors.New(ret.Message)
	}

	return ret, nil
}

// Push 高级推送
func (t *CPS) Push(title, body string, config *PushConfig) (*PushResponse, error) {
	if body == "" {
		return nil, fmt.Errorf("内容不能为空")
	}
	params, err := t.getPushParams(config)
	if err != nil {
		return nil, err
	}
	params["Title"] = title
	params["Body"] = body

	// 签名
	sign := t.sign(http.MethodGet, params)
	params["Signature"] = sign

	// 获取响应
	query := url.Values{}
	for k := range params {
		query.Add(k, params[k])
	}
	res, err := t.request(CPSGateway, query)
	if err != nil {
		return nil, err
	}

	// 解析响应
	ret := new(PushResponse)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, err
	}

	if ret.Code != "" {
		return nil, errors.New(ret.Message)
	}

	return ret, nil
}

// DefaultPushConfig 默认配置
func (t *CPS) DefaultPushConfig() *PushConfig {
	return &PushConfig{
		Target:       TargetAll,
		TargetValue:  []string{TargetValueAll},
		DeviceType:   DeviceTypeAll,
		PushType:     PushTypeNotice,
		Env:          IOSApnsEnvProduct,
		StoreOffline: true,
	}
}

func (t *CPS) getPushParams(config *PushConfig) (map[string]string, error) {
	params := t.getCommonParams()
	params["Action"] = PushAction

	// 推送目标
	if config.Target == "" || config.Target == TargetAll {
		params["Target"] = TargetAll
		params["TargetValue"] = TargetValueAll
	} else {
		params["Target"] = config.Target
		if len(config.TargetValue) == 1 {
			params["TargetValue"] = config.TargetValue[0]
		} else {
			params["TargetValue"] = strings.Join(config.TargetValue, ",")
			if config.Target == TargetTag {
				tags, err := t.MergeTag(config.TargetValue)
				if err != nil {
					return nil, fmt.Errorf("格式化标签错误: %v", err)
				}
				params["TargetValue"] = tags
			}
		}
	}

	// 设备类型
	if config.DeviceType == "" {
		params["DeviceType"] = DeviceTypeAll
	} else {
		params["DeviceType"] = config.DeviceType
	}

	// 通知类型
	if config.PushType == "" {
		params["PushType"] = PushTypeNotice
	} else {
		params["PushType"] = config.PushType
	}

	// 定时发送
	if !config.PushTime.IsZero() {
		params["PushTime"] = config.PushTime.UTC().Format("2006-01-02T15:04:05Z")
	}

	if config.Env != "" {
		params["iOSApnsEnv"] = config.Env
	} else {
		params["iOSApnsEnv"] = IOSApnsEnvProduct
	}

	// 扩展属性
	if len(config.ExtParameters) > 0 {
		ext, err := json.Marshal(config.ExtParameters)
		if err != nil {
			return nil, fmt.Errorf("格式化扩展参数错误: %v", err)
		}
		params["iOSExtParameters"] = string(ext)
		params["AndroidExtParameters"] = string(ext)
	}

	// 离线消息
	if config.StoreOffline {
		params["StoreOffline"] = "true"
		params["ExpireTime"] = time.Now().Add(time.Duration(86400) * time.Second).UTC().Format("2006-01-02T15:04:05Z")
	}

	if config.AndroidOpenType != "" {
		params["AndroidOpenType"] = config.AndroidOpenType
	}

	return params, nil
}
