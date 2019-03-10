package im

import (
	"encoding/json"
	"fmt"
)

// ServiceOpenConfig 运营数据服务名
const ServiceOpenConfig = "openconfigsvr"

// AppInfoRequestField 全部运营字段
var AppInfoRequestField = []string{
	"APNSMsgNum",
	"ActiveUserNum",
	"AppId",
	"AppName",
	"C2CAPNSMsgNum",
	"C2CDownMsgNum",
	"C2CSendMsgUserNum",
	"C2CUpMsgNum",
	"CallBackReq",
	"CallBackRsp",
	"ChainDecrease",
	"ChainIncrease",
	"Company",
	"Date",
	"DownMsgNum",
	"GroupAPNSMsgNum",
	"GroupAllGroupNum",
	"GroupDestroyGroupNum",
	"GroupDownMsgNum",
	"GroupJoinGroupTimes",
	"GroupNewGroupNum",
	"GroupQuitGroupTimes",
	"GroupSendMsgGroupNum",
	"GroupSendMsgUserNum",
	"GroupUpMsgNum",
	"LoginTimes",
	"LoginUserNum",
	"MaxOnlineNum",
	"QQ",
	"RegistUserNumOneDay",
	"RegistUserNumTotal",
	"SendMsgUserNum",
	"TextMsgInNum",
	"TextMsgOutNum",
	"UpMsgNum",
	"VoiceMsgInNum",
	"VoiceMsgOutNum",
}

// AppInfoDefaultRequestField 默认请求字段
var AppInfoDefaultRequestField = []string{
	"ActiveUserNum",
	"UpMsgNum",
	"DownMsgNum",
	"SendMsgUserNum",
	"MaxOnlineNum",
}

// AppInfoItem 统计数据
type AppInfoItem struct {
	APNSMsgNum           string
	ActiveUserNum        string
	AppID                string `json:"AppId"`
	AppName              string
	C2CAPNSMsgNum        string
	C2CDownMsgNum        string
	C2CSendMsgUserNum    string
	C2CUpMsgNum          string
	CallBackReq          string
	CallBackRsp          string
	ChainDecrease        string
	ChainIncrease        string
	Company              string
	Date                 string
	DownMsgNum           string
	GroupAPNSMsgNum      string
	GroupAllGroupNum     string
	GroupDestroyGroupNum string
	GroupDownMsgNum      string
	GroupJoinGroupTimes  string
	GroupNewGroupNum     string
	GroupQuitGroupTimes  string
	GroupSendMsgGroupNum string
	GroupSendMsgUserNum  string
	GroupUpMsgNum        string
	LoginTimes           string
	LoginUserNum         string
	MaxOnlineNum         string
	QQ                   string
	RegistUserNumOneDay  string
	RegistUserNumTotal   string
	SendMsgUserNum       string
	TextMsgInNum         string
	TextMsgOutNum        string
	UpMsgNum             string
	VoiceMsgInNum        string
	VoiceMsgOutNum       string
}

// AppInfoRequest 运营数据请求结构
type AppInfoRequest struct {
	RequestField []string
}

// AppInfoResponse 运营数据响应
type AppInfoResponse struct {
	Response
	Result []*AppInfoItem
}

// GetAppInfo 拉取运营数据
func (t *IM) GetAppInfo(fields []string) ([]*AppInfoItem, error) {
	var req interface{}
	if len(fields) > 0 {
		req = &AppInfoRequest{RequestField: fields}
	} else {
		req = &EmptyRequest{}
	}
	res, err := t.api(ServiceOpenConfig, "getappinfo", req)
	if err != nil {
		return nil, err
	}

	response := new(AppInfoResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.Result, nil
}
