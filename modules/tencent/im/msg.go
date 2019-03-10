package im

import (
	"encoding/json"
	"fmt"
)

// ServiceOpenMSG 消息数据管理服务名
const ServiceOpenMSG = "open_msg_svc"

const (
	// ChatTypeC2C 单发消息类型
	ChatTypeC2C string = "C2C"
	// ChatTypeGroup 群组消息类型
	ChatTypeGroup string = "Group"
)

// GetMsgHistoryRequest 消息记录下载请求
type GetMsgHistoryRequest struct {
	ChatType string // 消息类型，C2C:单发消息 Group:群组消息
	MsgTime  string // 需要下载的时间段，2015120121表示获取2015年12月1日21:00~21:59的消息的下载地址
}

// MsgHistory 历史消息
type MsgHistory struct {
	URL        string // 消息记录文件下载地址
	ExpireTime string // 下载地址过期时间
	FileSize   int64  // GZip压缩前的文件大小（单位Byte）
	FileMD5    string // GZip压缩前的文件MD5
	GzipSize   int64  // GZip压缩后的文件大小（单位Byte）
	GzipMD5    string // GZip压缩后的文件MD5
}

// MsgHistoryResponse 历史消息查询响应
type MsgHistoryResponse struct {
	Response
	File []*MsgHistory
}

// GetHistory 消息记录下载
// 获取指定某天某小时的消息记录下载地址
// 消息类型，C2C:单发消息 Group:群组消息
// 时间，2015120121表示获取2015年12月1日21:00~21:59的消息的下载地址
func (t *IM) GetHistory(chatType string, msgTime string) (*MsgHistory, error) {
	req := &GetMsgHistoryRequest{ChatType: chatType, MsgTime: msgTime}
	res, err := t.api(ServiceOpenMSG, "get_history", req)
	if err != nil {
		return nil, err
	}

	response := new(MsgHistoryResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.File[0], nil
}
