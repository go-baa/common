package im

import (
	"encoding/json"
	"fmt"

	"time"

	"git.code.tencent.com/xinhuameiyu/common/util"
)

// ServiceOpenIM 消息服务名
const ServiceOpenIM = "openim"

// AccountRequest 通用用户操作请求
type AccountRequest struct {
	ToAccount []string `json:"To_Account"`
}

// 单聊消息

// msgSync 消息同步选项
var msgSync = map[bool]int{true: 1, false: 2}

// msgImportRealtime 消息导入选项
var msgImportRealtime = map[bool]int{true: 1, false: 2}

// MsgRequest 发送消息请求
type MsgRequest struct {
	SyncOtherMachine int              `json:",omitempty"`             // 1 消息同步至发送方(默认), 2 消息不同步至发送方
	FromAccount      string           `json:"From_Account,omitempty"` // 消息发送方账号, 不指定则为管理员账号
	ToAccount        interface{}      `json:"To_Account"`             // 消息接收方账号, string | []string
	MsgRandom        int              // 消息随机数
	MsgTimeStamp     int64            // 消息时间戳
	MsgBody          []*MsgBodyItem   // 消息内容
	OfflinePushInfo  *OfflinePushInfo `json:",omitempty"` // 离线推送信息配置
}

// ImportMsgRequest 消息导入请求
type ImportMsgRequest struct {
	SyncFromOldSystem int            // 1，实时消息导入，消息加入未读计数；2，历史消息导入，消息不计入未读
	FromAccount       string         `json:"From_Account"` // 消息发送方账号, 不指定则为管理员账号
	ToAccount         string         `json:"To_Account"`   // 消息接收方账号
	MsgRandom         int            // 消息随机数
	MsgTimeStamp      int64          // 消息时间戳
	MsgBody           []*MsgBodyItem // 消息内容
}

// SendMsg 单发单聊消息
// pushInfo为可选值，不需要推送时使用nil
func (t *IM) SendMsg(from, to string, msg []*MsgBodyItem, sync bool, pushInfo *OfflinePushInfo) error {
	req := &MsgRequest{
		SyncOtherMachine: msgSync[sync],
		FromAccount:      from,
		ToAccount:        to,
		MsgRandom:        util.StringToInt(string(util.RandStr(6, util.KC_RAND_KIND_NUM))),
		MsgTimeStamp:     time.Now().Unix(),
		MsgBody:          msg,
		OfflinePushInfo:  pushInfo,
	}

	res, err := t.api(ServiceOpenIM, "sendmsg", req)
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

// BatchSendMsg 批量发单聊消息
func (t *IM) BatchSendMsg(from string, to []string, msg []*MsgBodyItem, sync bool, pushInfo *OfflinePushInfo) error {
	if len(to) > 500 {
		return fmt.Errorf("一次最多给500个用户进行单发消息")
	}

	req := &MsgRequest{
		SyncOtherMachine: msgSync[sync],
		FromAccount:      from,
		ToAccount:        to,
		MsgRandom:        util.StringToInt(string(util.RandStr(6, util.KC_RAND_KIND_NUM))),
		MsgTimeStamp:     time.Now().Unix(),
		MsgBody:          msg,
	}

	if pushInfo != nil {
		req.OfflinePushInfo = pushInfo
	}

	res, err := t.api(ServiceOpenIM, "batchsendmsg", req)
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

// ImportMsg 导入单聊消息
func (t *IM) ImportMsg(from, to string, msg []*MsgBodyItem, msgTimeStamp int64, realtime bool) error {
	req := &ImportMsgRequest{
		SyncFromOldSystem: msgImportRealtime[realtime],
		FromAccount:       from,
		ToAccount:         to,
		MsgRandom:         util.StringToInt(string(util.RandStr(6, util.KC_RAND_KIND_NUM))),
		MsgTimeStamp:      msgTimeStamp,
		MsgBody:           msg,
	}

	res, err := t.api(ServiceOpenIM, "importmsg", req)
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

// 消息推送 增值服务，需要申请开通才能使用

// Push 推送 TODO
func (t *IM) Push() {

}

// PushReport 推送任务状态
type PushReport struct {
	TaskID    string `json:"TaskId"`
	Status    int    // 0(未处理) / 1（推送中) / 2（推送完成）
	StartTime string `json:",omitempty"` // 推送开始时间，Status不为0，才有这个字段
	Finished  int    `json:",omitempty"` // 已完成推送人数
	Total     int    `json:",omitempty"` // 需推送总人数
}

// QueryPushReportRequest 推送任务查询请求
type QueryPushReportRequest struct {
	TaskIds []string
}

// QueryPushReportResponse 推送任务查询响应
type QueryPushReportResponse struct {
	Response
	Reports []*PushReport
}

// GetPushReport 获取推送报告
func (t *IM) GetPushReport(taskids []string) ([]*PushReport, error) {
	if len(taskids) > 500 {
		return nil, fmt.Errorf("每次最多只能查询500个任务")
	}

	req := new(QueryPushReportRequest)
	req.TaskIds = taskids
	res, err := t.api(ServiceOpenIM, "im_get_push_report", req)
	if err != nil {
		return nil, err
	}

	response := new(QueryPushReportResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.Reports, nil
}

// AppAttr 应用推送属性
type AppAttr struct {
	AttrNames map[string]string // 数字键，表示第几个属性，（"0"到"9"之间），属性名最长不超过50字节。应用最多可以有10个推送属性
}

// QueryAppAttrResponse 查询应用推送属性请求
type QueryAppAttrResponse struct {
	Response
	AppAttr
}

// SetPushAttr 设置应用属性名称
func (t *IM) SetPushAttr(attrs []string) error {
	if len(attrs) > 10 {
		return fmt.Errorf("最多只能设置10个属性")
	}

	req := new(AppAttr)
	req.AttrNames = make(map[string]string)
	for k, v := range attrs {
		req.AttrNames[util.IntToString(k)] = v
	}
	res, err := t.api(ServiceOpenIM, "im_set_attr_name", req)
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

// GetAppAttr 获取应用属性名称
func (t *IM) GetAppAttr() (map[string]string, error) {
	req := new(EmptyRequest)
	res, err := t.api(ServiceOpenIM, "im_get_attr_name", req)
	if err != nil {
		return nil, err
	}

	response := new(QueryAppAttrResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.AttrNames, nil
}

// AccountAttr 用户推送属性
type AccountAttr struct {
	ToAccount string `json:"To_Account"`
	Attrs     map[string]string
}

// SetAccountAttrRequest 设置用户属性请求
type SetAccountAttrRequest struct {
	UserAttrs []*AccountAttr
}

// RemoveAccountAttrItem 删除用户属性
type RemoveAccountAttrItem struct {
	ToAccount string `json:"To_Account"`
	Attrs     []string
}

// RemoveAccountAttrRequest 删除用户属性请求
type RemoveAccountAttrRequest struct {
	UserAttrs []*RemoveAccountAttrItem
}

// QueryAccountAttrResponse 查询用户属性响应
type QueryAccountAttrResponse struct {
	Response
	UserAttrs []*AccountAttr
}

// SetAccountAttr 设置用户属性
func (t *IM) SetAccountAttr(accountAttrs []*AccountAttr) error {
	if len(accountAttrs) > 500 {
		return fmt.Errorf("每次最多只能给500个用户设置属性")
	}

	req := new(SetAccountAttrRequest)
	req.UserAttrs = accountAttrs
	res, err := t.api(ServiceOpenIM, "im_set_attr", req)
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

// RemoveAccountAttr 删除用户属性
func (t *IM) RemoveAccountAttr(accountAttrs []*RemoveAccountAttrItem) error {
	if len(accountAttrs) > 500 {
		return fmt.Errorf("每次最多只能给500个用户删除属性")
	}

	req := new(RemoveAccountAttrRequest)
	req.UserAttrs = accountAttrs
	res, err := t.api(ServiceOpenIM, "im_remove_attr", req)
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

// GetAccountAttr 获取用户属性
func (t *IM) GetAccountAttr(accounts []string) ([]*AccountAttr, error) {
	if len(accounts) > 500 {
		return nil, fmt.Errorf("每次最多只能获取500个用户的属性")
	}

	req := new(AccountRequest)
	req.ToAccount = accounts
	res, err := t.api(ServiceOpenIM, "im_get_attr", req)
	if err != nil {
		return nil, err
	}

	response := new(QueryAccountAttrResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.UserAttrs, nil
}

// AccountTags 用户标签
type AccountTags struct {
	ToAccount string `json:"To_Account"`
	Tags      []string
}

// SetAccountTagsRequest 用户标签操作请求
type SetAccountTagsRequest struct {
	UserTags []*AccountTags
}

// AddTag 添加用户标签
func (t *IM) AddTag(accountTags []*AccountTags) error {
	if len(accountTags) > 500 {
		return fmt.Errorf("每次最多只能给500个用户添加标签")
	}

	req := new(SetAccountTagsRequest)
	req.UserTags = accountTags
	res, err := t.api(ServiceOpenIM, "im_add_tag", req)
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

// RemoveTag 删除用户标签
func (t *IM) RemoveTag(accountTags []*AccountTags) error {
	if len(accountTags) > 500 {
		return fmt.Errorf("每次最多只能给500个用户删除标签")
	}

	req := new(SetAccountTagsRequest)
	req.UserTags = accountTags
	res, err := t.api(ServiceOpenIM, "im_remove_tag", req)
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

// RemoveAllTags 删除用户所有标签
func (t *IM) RemoveAllTags(accounts []string) error {
	if len(accounts) > 500 {
		return fmt.Errorf("每次最多只能删除500个用户的所有标签")
	}
	req := new(AccountRequest)
	req.ToAccount = accounts
	res, err := t.api(ServiceOpenIM, "im_remove_all_tags", req)
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

// 在线状态

// AccountState 用户在线状态
type AccountState struct {
	ToAccount string `json:"To_Account"`
	State     string
}

// QueryStateRequest 用户在线状态查询请求
type QueryStateRequest struct {
	ToAccount []string `json:"To_Account"`
}

// QueryStateResponse 用户在线状态查询响应
type QueryStateResponse struct {
	Response
	QueryResult []*AccountState
}

// QueryState 用户在线状态查询
func (t *IM) QueryState(accounts []string) ([]*AccountState, error) {
	if len(accounts) > 500 {
		return nil, fmt.Errorf("一次最多查询500个账号状态")
	}

	req := new(QueryStateRequest)
	req.ToAccount = accounts
	res, err := t.api(ServiceOpenIM, "querystate", req)
	if err != nil {
		return nil, err
	}

	response := new(QueryStateResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.QueryResult, nil
}
