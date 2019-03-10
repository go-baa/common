package im

import (
	"encoding/json"
	"fmt"
)

// ServiceSNS 关系链服务名
const ServiceSNS = "sns"

const (
	// DeleteTypeSingle 单向删除
	DeleteTypeSingle = "Delete_Type_Single"
	// DeleteTypeBoth 双向删除
	DeleteTypeBoth = "Delete_Type_Both"
)

const (
	// CheckTypeSingle 单向校验
	CheckTypeSingle = "CheckResult_Type_Singal"
	// CheckTypeBoth 双向校验
	CheckTypeBoth = "CheckResult_Type_Both"
	// CheckResultBoth 双向好友
	CheckResultBoth = "CheckResult_Type_BothWay"
	// CheckResultAWithB from有to，to没有from
	CheckResultAWithB = "CheckResult_Type_AWithB"
	// CheckResultBWithA from没有to，to有from
	CheckResultBWithA = "CheckResult_Type_BWithA"
	// CheckResultNoRelation 无关系
	CheckResultNoRelation = "CheckResult_Type_NoRelation"
)

const (
	// BlackCheckTypeSingle 单向校验黑名单
	BlackCheckTypeSingle = "BlackCheckResult_Type_Singal"
	// BlackCheckTypeBoth 双向校验黑名单
	BlackCheckTypeBoth = "BlackCheckResult_Type_Both"
	// BlackCheckResultBoth 双向黑名单
	BlackCheckResultBoth = "BlackCheckResult_Type_BothWay"
	// BlackCheckResultAWithB from有to，to没有from
	BlackCheckResultAWithB = "BlackCheckResult_Type_AWithB"
	// BlackCheckResultBWithA from没有to，to有from
	BlackCheckResultBWithA = "BlackCheckResult_Type_BWithA"
	// BlackCheckResultNoRelation 无关系
	BlackCheckResultNoRelation = "BlackCheckResult_Type_NoRelation"
)

const (
	// GetAllTypeYes 全量更新
	GetAllTypeYes = "GetAll_Type_YES"
	// GetAllTypeNo 非全量更新
	GetAllTypeNo = "GetAll_Type_NO"
)

const (
	// AddFriendTypeBoth 双向加好友
	AddFriendTypeBoth = "Add_Type_Both"
	// AddFriendTypeSingle 单向加好友
	AddFriendTypeSingle = "Add_Type_Single"
)

// 加好友来源，自定义
const (
	// AddSourceTypeSystem 系统
	AddSourceTypeSystem = "AddSource_Type_System"
	// AddSourceTypeAndroid 安卓
	AddSourceTypeAndroid = "AddSource_Type_Android"
	// AddSourceTypeIOS ios
	AddSourceTypeIOS = "AddSource_Type_IOS"
)

// addFriendForce 强制加好友标记
var addFriendForce = map[bool]int{true: 1, false: 0}

// FromAccount ...
type FromAccount struct {
	FromAccount string `json:"From_Account"`
}

// ToAccount ...
type ToAccount struct {
	ToAccount []string `json:"To_Account"`
}

// SNSResponseItem 关系链响应单条项目
type SNSResponseItem struct {
	ToAccount  string `json:"To_Account"`
	ResultCode int
	ResultInfo string
	Relation   string `json:",omitempty"` // 好友、黑名单校验结果
}

// SNSRelationResponse 好友关系响应
type SNSRelationResponse struct {
	Response
	FailAccount    []string           `json:"Fail_Account"`    // 处理失败的用户列表
	InvalidAccount []string           `json:"Invalid_Account"` // 请求包中的非法用户列表
	InfoItem       []*SNSResponseItem `json:",omitempty"`      // 好友关系校验结果
	ResultItem     []*SNSResponseItem `json:",omitempty"`      // 黑名单校验结果
}

// SNSInfoItem 好友列表信息单条项目
type SNSInfoItem struct {
	InfoAccount    string         `json:"Info_Account"` // 好友账号
	SnsProfileItem []*ProfileItem // 好友资料列表
}

// AddFriendItem 添加好友单条项目
type AddFriendItem struct {
	ToAccount  string   `json:"To_Account"`
	AddSource  string   // 加好友来源字段
	Remark     string   `json:",omitempty"` // From_Account对To_Account的好友备注
	GroupName  []string `json:",omitempty"` // From_Account对To_Account的分组信息
	AddWording string   `json:",omitempty"` // From_Account和To_Account形成好友关系时的附言信息
}

// AddFriendRequest 添加好友请求
type AddFriendRequest struct {
	FromAccount   string           `json:"From_Account"`
	AddFriendItem []*AddFriendItem // 添加好友列表
	AddType       string           `json:",omitempty"` // 加好友方式（默认双向加好友方式）："Add_Type_Single"表示单向加好友；"Add_Type_Both"表示双向加好友。
	ForceAddFlags int              // 管理员强制加好友标记：1表示强制加好友；0表示常规加好友方式
}

// FriendAdd 添加好友
func (t *IM) FriendAdd(from string, to []*AddFriendItem, addType string, force bool) ([]*SNSResponseItem, []string, []string, error) {
	req := &AddFriendRequest{
		FromAccount:   from,
		AddFriendItem: to,
		AddType:       addType,
		ForceAddFlags: addFriendForce[force],
	}
	res, err := t.api(ServiceSNS, "friend_add", req)
	if err != nil {
		return nil, nil, nil, err
	}

	response := new(SNSRelationResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, nil, nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.ResultItem, response.FailAccount, response.InvalidAccount, nil
}

// GetFriendRequest 拉取好友请求
type GetFriendRequest struct {
	FromAccount          string   `json:"From_Account"` // 需要拉取该Identifier的好友
	TimeStamp            int64    `json:",omitempty"`   // 上次拉取的时间戳，不填或为0时表示全量拉取
	StartIndex           int      // 拉取的起始位置
	TagList              []string `json:",omitempty"` // 指定要拉取的资料字段及好友字段
	LastStandardSequence int      `json:",omitempty"` // 上次拉取标配关系链的Sequence，仅在只拉取标配关系链字段时有用
	GetCount             int      `json:",omitempty"` // 每页需要拉取的数量，默认每页拉去100个好友
}

// GetFriendResponse 拉取好友响应
type GetFriendResponse struct {
	Response
	NeedUpdateAll           string         // 是否需要全量更新："GetAll_Type_YES"表示需要全量更新，"GetAll_Type_NO"表示不需要全量更新
	TimeStampNow            int64          // 本次拉取的时间戳，客户端需要保存该时间
	StartIndex              int            // 下页拉取的起始位置
	InfoItem                []*SNSInfoItem // 好友对象数组
	CurrentStandardSequence int            // 本次拉取标配关系链的Sequence，客户端需要保存该Sequence
	FriendNum               int            // 好友总数
}

// FriendGetAll 拉取好友
func (t *IM) FriendGetAll(from string, timestamp int64, startIndex, lastSeq, limit int, tags []string) (*GetFriendResponse, error) {
	req := &GetFriendRequest{
		FromAccount:          from,
		TimeStamp:            timestamp,
		StartIndex:           startIndex,
		LastStandardSequence: lastSeq,
		GetCount:             limit,
	}

	if len(tags) > 0 {
		req.TagList = tags
	}

	res, err := t.api(ServiceSNS, "friend_get_all", req)
	if err != nil {
		return nil, err
	}

	response := new(GetFriendResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response, nil
}

// SNSDeleteRequest 删除关系请求
type SNSDeleteRequest struct {
	FromAccount string   `json:"From_Account"`
	ToAccount   []string `json:"To_Account"`
	DeleteType  string   // 删除模式
}

// FriendDelete 删除好友
func (t *IM) FriendDelete(from string, to []string, deleteType string) ([]*SNSResponseItem, []string, []string, error) {
	req := &SNSDeleteRequest{
		FromAccount: from,
		ToAccount:   to,
		DeleteType:  deleteType,
	}
	res, err := t.api(ServiceSNS, "friend_delete", req)
	if err != nil {
		return nil, nil, nil, err
	}

	response := new(SNSRelationResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, nil, nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.ResultItem, response.FailAccount, response.InvalidAccount, nil
}

// FriendDeleteAll 删除所有好友
func (t *IM) FriendDeleteAll(from string) error {
	req := &FromAccount{FromAccount: from}
	res, err := t.api(ServiceSNS, "friend_delete_all", req)
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

// SNSCheckRequest 关系校验请求
type SNSCheckRequest struct {
	FromAccount string   `json:"From_Account"`
	ToAccount   []string `json:"To_Account"`
	CheckType   string   // 校验模式
}

// FriendCheck 校验好友
// 返回校验结果列表、失败用户列表、非法用户列表
func (t *IM) FriendCheck(from string, to []string, checkType string) ([]*SNSResponseItem, []string, []string, error) {
	req := &SNSCheckRequest{
		FromAccount: from,
		ToAccount:   to,
		CheckType:   checkType,
	}
	res, err := t.api(ServiceSNS, "friend_check", req)
	if err != nil {
		return nil, nil, nil, err
	}

	response := new(SNSRelationResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, nil, nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.InfoItem, response.FailAccount, response.InvalidAccount, nil
}

// BlackListCheck 校验黑名单
func (t *IM) BlackListCheck(from string, to []string, checkType string) ([]*SNSResponseItem, []string, []string, error) {
	req := &SNSCheckRequest{
		FromAccount: from,
		ToAccount:   to,
		CheckType:   checkType,
	}
	res, err := t.api(ServiceSNS, "black_list_check", req)
	if err != nil {
		return nil, nil, nil, err
	}

	response := new(SNSRelationResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, nil, nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.ResultItem, response.FailAccount, response.InvalidAccount, nil
}
