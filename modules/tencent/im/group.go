package im

import (
	"encoding/json"
	"fmt"

	"git.code.tencent.com/xinhuameiyu/common/util"
)

// ServiceGroupOpen 群组服务名
const ServiceGroupOpen = "group_open_http_svc"

const (
	// GroupTypePublic 公开群
	GroupTypePublic = "Public"
	// GroupTypePrivate 私有群
	GroupTypePrivate = "Private"
	// GroupTypeChatRoom 聊天室
	GroupTypeChatRoom = "ChatRoom"
	// GroupTypeAVChatRoom 互动直播聊天室
	GroupTypeAVChatRoom = "AVChatRoom"
	// GroupTypeBChatRoom 在线成员广播大群
	GroupTypeBChatRoom = "BChatRoom"
)

const (
	// GroupRoleOwner 群主
	GroupRoleOwner = "Owner"
	// GroupRoleAdmin 管理员
	GroupRoleAdmin = "Admin"
	// GroupRoleMember 群成员
	GroupRoleMember = "Member"
	// GroupRoleNotMember 非群成员
	GroupRoleNotMember = "NotMember"
)

const (
	// GroupMsgFlagNotify 接收并提示
	GroupMsgFlagNotify = "AcceptAndNotify"
	// GroupMsgFlagNotNotify 接收不提示
	GroupMsgFlagNotNotify = "AcceptNotNotify"
	// GroupMsgFlagDiscard 屏蔽群消息
	GroupMsgFlagDiscard = "Discard"
)

const (
	// GroupApplyOptionDisable 禁止任何人申请加入
	GroupApplyOptionDisable = "DisableApply"
	// GroupApplyOptionNeedPermission 需要群主或管理员审批
	GroupApplyOptionNeedPermission = "NeedPermission"
	// GroupApplyOptionFreeAccess 允许无需审批自由加入群组
	GroupApplyOptionFreeAccess = "FreeAccess"
)

const (
	// AppDefinedKeyGroupLevel 自定义群级别字段
	AppDefinedKeyGroupLevel = "CustomLevel"
	// AppDefinedKeyAllowInvite 自定义是否允许邀请字段
	AppDefinedKeyAllowInvite = "AllowInvite"
)

const (
	// GroupImportMemberFailed 导入或添加群成员失败
	GroupImportMemberFailed = iota
	// GroupImportMemberSuccess 导入或添加群成员成功
	GroupImportMemberSuccess
	// GroupImportMemberExist 导入或添加的群成员已存在
	GroupImportMemberExist
)

// silenceAction 静默处理选项
var silenceAction = map[bool]int{true: 1, false: 0}

// GroupID 群组ID
type GroupID struct {
	GroupID string `json:"GroupId"` // 群组ID
}

// GroupIDList 群组ID列表
type GroupIDList struct {
	GroupIDList []string `json:"GroupIdList"`
}

// GroupInfo 群组基本资料
type GroupInfo struct {
	Type            string             `json:",omitempty"`              // 群组形态
	Name            string             `json:",omitempty"`              // 群组名称（最长30字节）
	GroupID         string             `json:"GroupId,omitempty"`       // 群组ID
	Introduction    string             `json:",omitempty"`              // 群组简介（最长120字节）
	Notification    string             `json:",omitempty"`              // 群组公告（最长150字节）
	FaceURL         string             `json:"FaceUrl,omitempty"`       // 群组头像URL（最长100字节）
	OwnerAccount    string             `json:"Owner_Account,omitempty"` // 群主ID
	CreateTime      int64              `json:",omitempty"`              // 群组的创建时间
	InfoSeq         int                `json:",omitempty"`              // 群资料的每次变都会增加该值
	LastInfoTime    int64              `json:",omitempty"`              // 群组最后一次信息变更时间
	LastMsgTime     int                `json:",omitempty"`              // 群组内最后发消息的时间
	NextMsgSeq      int64              `json:",omitempty"`              // 群内下一条消息的Seq
	MemberNum       int                `json:",omitempty"`              // 当前成员数量
	MaxMemberNum    int                `json:",omitempty"`              // 最大成员数量
	ApplyJoinOption string             `json:",omitempty"`              // 申请加群选项
	AppDefinedData  []*AppDefinedData  `json:",omitempty"`              // 群资料自定义字段
	MemberList      []*GroupMemberInfo `json:",omitempty"`              // 群成员列表
}

// GroupMemberAccount 群组成员账号
type GroupMemberAccount struct {
	MemberAccount string `json:"Member_Account"` // 群成员帐号
}

// GroupMemberInfo 群组成员资料
type GroupMemberInfo struct {
	MemberAccount        string            `json:"Member_Account"` // 群成员帐号
	Role                 string            `json:",omitempty"`     // 群内身份
	JoinTime             int64             `json:",omitempty"`     // 入群时间
	MsgSeq               int64             `json:",omitempty"`     // 该成员当前已读消息Seq
	MsgFlag              string            `json:",omitempty"`     // 消息接收选项
	LastSendMsgTime      int64             `json:",omitempty"`     // 最后发送消息的时间
	NameCard             string            `json:",omitempty"`     // 群名片
	AppMemberDefinedData []*AppDefinedData `json:",omitempty"`     // 群成员维度自定义字段
}

// AppDefinedData 自定义资料
type AppDefinedData struct {
	Key   string
	Value string
}

// QueryGroupListRequest 群组列表查询请求
type QueryGroupListRequest struct {
	GroupType string `json:",omitempty"` // 特定群组形态的群组
	Limit     int    `json:",omitempty"` // 本次获取的群组ID数量的上限，不得超过10000
	Next      int    `json:",omitempty"` // 群太多时分页拉取标志，第一次填0
}

// GroupListResponse 群组列表查询响应
type GroupListResponse struct {
	Response
	TotalCount  int        // APP当前的群组总数，可以通过GroupType进行过滤
	GroupIDList []*GroupID `json:"GroupIdList"` // 获取到的群组ID的集合
	Next        int        // 分页拉取的标志
}

// GetAppGroupList 获取APP中的群组ID列表
// 返回值：列表，总数，下页偏移值
func (t *IM) GetAppGroupList(groupType string, next, limit int) ([]string, int, int, error) {
	req := &QueryGroupListRequest{
		GroupType: groupType,
		Limit:     limit,
		Next:      next,
	}

	res, err := t.api(ServiceGroupOpen, "get_appid_group_list", req)
	if err != nil {
		return nil, 0, 0, err
	}

	response := new(GroupListResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, 0, 0, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	var groupID []string
	for _, v := range response.GroupIDList {
		groupID = append(groupID, v.GroupID)
	}

	return groupID, response.TotalCount, response.Next, nil
}

// CreateGroupResponse 创建群组响应
type CreateGroupResponse struct {
	Response
	GroupID
}

// CreateGroup 创建群组
func (t *IM) CreateGroup(req *GroupInfo) (string, error) {
	if req.Name == "" || req.Type == "" {
		return "", fmt.Errorf("群组名称与类型为必填项")
	}

	res, err := t.api(ServiceGroupOpen, "create_group", req)
	if err != nil {
		return "", err
	}

	response := new(CreateGroupResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return "", fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		if response.ErrorCode == 10021 {
			return req.GroupID, nil
		}
		return "", fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.GroupID.GroupID, nil
}

// GetGroupInfoRequest 获取群组资料请求
type GetGroupInfoRequest struct {
	GroupIDList    []string        `json:"GroupIdList"` // 群组ID列表
	ResponseFilter *ResponseFilter `json:",omitempty"`  // 过滤器，选填
}

// GetGroupInfoResponse 获取群组资料响应
type GetGroupInfoResponse struct {
	Response
	GroupInfo []*GroupInfo
}

// GetGroupInfo 获取群组详细资料
func (t *IM) GetGroupInfo(groupIDs []string, filter *ResponseFilter) ([]*GroupInfo, error) {
	if len(groupIDs) > 50 {
		return nil, fmt.Errorf("一次请求最多获取50个群组资料")
	}

	req := &GetGroupInfoRequest{
		GroupIDList: groupIDs,
	}
	if filter != nil {
		req.ResponseFilter = filter
	}
	res, err := t.api(ServiceGroupOpen, "get_group_info", req)
	if err != nil {
		return nil, err
	}

	response := new(GetGroupInfoResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.GroupInfo, nil
}

// GetGroupMemberRequest 获取群成员资料请求
type GetGroupMemberRequest struct {
	GroupID string `json:"GroupId"`    // 群组ID
	Limit   int    `json:",omitempty"` // 最多获取多少个成员的资料
	Offset  int    `json:",omitempty"` // 从第多少个成员开始获取
	ResponseFilter
}

// GetGroupMemberResponse 获取群成员资料响应
type GetGroupMemberResponse struct {
	Response
	MemberNum  int // 本群组的群成员总数
	MemberList []*GroupMemberInfo
}

// GetGroupMemberInfo 获取群成员详细资料
// filter支持：MemberInfoFilter，MemberRoleFilter，AppDefinedDataFilterGroupMember
func (t *IM) GetGroupMemberInfo(groupID string, offset, limit int, filter *ResponseFilter) ([]*GroupMemberInfo, int, error) {
	req := &GetGroupMemberRequest{
		GroupID: groupID,
		Offset:  offset,
		Limit:   limit,
	}
	if filter != nil {
		req.ResponseFilter = *filter
	}

	res, err := t.api(ServiceGroupOpen, "get_group_member_info", req)
	if err != nil {
		return nil, 0, err
	}

	response := new(GetGroupMemberResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, 0, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, 0, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.MemberList, response.MemberNum, nil
}

// ModifyGroupBaseInfo 修改群组基础资料
func (t *IM) ModifyGroupBaseInfo(req *GroupInfo) error {
	if req.GroupID == "" {
		return fmt.Errorf("必须制定群组ID")
	}

	res, err := t.api(ServiceGroupOpen, "modify_group_base_info", req)
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

// AddGroupMemberRequest 增加群组成员请求
type AddGroupMemberRequest struct {
	GroupID    string                `json:"GroupId"`    // 操作的群ID
	Silence    int                   `json:",omitempty"` // 是否静默加人。0：非静默加人；1：静默加人。不填该字段默认为0
	MemberList []*GroupMemberAccount // 待添加的群成员数组
}

// AddGroupMember 增加群组成员
func (t *IM) AddGroupMember(groupID string, accounts []string, silence bool) ([]*ImportGroupMemberResponseItem, error) {
	if len(accounts) > 500 {
		return nil, fmt.Errorf("一次最多支持添加500个成员")
	}

	var memberList []*GroupMemberAccount
	for _, v := range accounts {
		member := new(GroupMemberAccount)
		member.MemberAccount = v
		memberList = append(memberList, member)
	}
	req := &AddGroupMemberRequest{
		GroupID:    groupID,
		Silence:    silenceAction[silence],
		MemberList: memberList,
	}
	res, err := t.api(ServiceGroupOpen, "add_group_member", req)
	if err != nil {
		return nil, err
	}

	response := new(ImportGroupMemberResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.MemberList, nil
}

// DeleteGroupMemberRequest 删除群组成员请求
type DeleteGroupMemberRequest struct {
	GroupID     string   `json:"GroupId"`             // 操作的群ID
	Silence     int      `json:",omitempty"`          // 是否静默删人。0：非静默删人；1：静默删人。不填该字段默认为0
	Reason      string   `json:",omitempty"`          // 踢出用户原因
	MemberToDel []string `json:"MemberToDel_Account"` // 待删除的群成员
}

// DeleteGroupMember 删除群组成员
func (t *IM) DeleteGroupMember(groupID string, accounts []string, reason string, silence bool) error {
	if len(accounts) > 500 {
		return fmt.Errorf("一次请求最多删除500个成员")
	}

	req := &DeleteGroupMemberRequest{
		GroupID:     groupID,
		Silence:     silenceAction[silence],
		Reason:      reason,
		MemberToDel: accounts,
	}
	res, err := t.api(ServiceGroupOpen, "delete_group_member", req)
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

// GroupMemberWritableInfo 群组成员可修改信息
type GroupMemberWritableInfo struct {
	Role                 string            `json:",omitempty"` // 群内身份
	MsgFlag              string            `json:",omitempty"` // 消息接收选项
	NameCard             string            `json:",omitempty"` // 群名片
	AppMemberDefinedData []*AppDefinedData `json:",omitempty"` // 群成员维度自定义字段
}

// ModifyMemberInfoRequest 修改群组成员资料请求
type ModifyMemberInfoRequest struct {
	GroupID       string `json:"GroupId"`        // 群组ID
	MemberAccount string `json:"Member_Account"` // 群成员帐号
	GroupMemberWritableInfo
}

// ModifyGroupMemberInfo 修改群组成员资料
func (t *IM) ModifyGroupMemberInfo(groupID, account string, data *GroupMemberWritableInfo) error {
	req := &ModifyMemberInfoRequest{
		GroupID:                 groupID,
		MemberAccount:           account,
		GroupMemberWritableInfo: *data,
	}
	res, err := t.api(ServiceGroupOpen, "modify_group_member_info", req)
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

// DestroyGroup 解散群组
func (t *IM) DestroyGroup(groupID string) error {
	req := &GroupID{GroupID: groupID}
	res, err := t.api(ServiceGroupOpen, "destroy_group", req)
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

// GetJoinedGroupRequest 获取用户所加入的群组请求
type GetJoinedGroupRequest struct {
	MemberAccount string `json:"Member_Account"` // 群成员帐号
	Limit         int    `json:",omitempty"`     // 拉取多少个，不填标识拉取全部
	Offset        int    `json:",omitempty"`     // 从第多少个开始拉取
	GroupType     string // 群组形态，不填为拉取所有
}

// GetJoinedGroupResponse 获取用户所加入的群组响应
type GetJoinedGroupResponse struct {
	Response
	GroupIDList []*GroupID
	TotalCount  int
}

// GetJoinedGroupList 获取用户所加入的群组, 暂不支持高级信息过滤
func (t *IM) GetJoinedGroupList(account, groupType string, offset, limit int) ([]string, int, error) {
	req := &GetJoinedGroupRequest{
		MemberAccount: account,
		Limit:         limit,
		Offset:        offset,
		GroupType:     groupType,
	}
	res, err := t.api(ServiceGroupOpen, "get_joined_group_list", req)
	if err != nil {
		return nil, 0, err
	}

	response := new(GetJoinedGroupResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, 0, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, 0, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	var groupIDList []string
	for _, v := range response.GroupIDList {
		groupIDList = append(groupIDList, v.GroupID)
	}

	return groupIDList, response.TotalCount, nil
}

// QueryRoleInGroupRequest 查询用户身份请求
type QueryRoleInGroupRequest struct {
	GroupID     string   `json:"GroupId"`      // 需要查询的群组ID
	UserAccount []string `json:"User_Account"` // 需要查询的用户账号，最多支持500个账号
}

// GroupMemberRole 用户在群组内身份
type GroupMemberRole struct {
	MemberAccount string `json:"Member_Account"` //
	Role          string // 成员在群内的身份信息，可能的身份包括Owner/Admin/Member/NotMember
}

// RoleInGroupResponse 用户身份查询响应
type RoleInGroupResponse struct {
	Response
	UserIDList []*GroupMemberRole `json:"UserIdList"`
}

// GetRoleInGroup 查询用户在群组中的身份
func (t *IM) GetRoleInGroup(groupID string, accounts []string) ([]*GroupMemberRole, error) {
	if len(accounts) > 500 {
		return nil, fmt.Errorf("一次请求最多查询500个成员")
	}

	req := &QueryRoleInGroupRequest{
		GroupID:     groupID,
		UserAccount: accounts,
	}
	res, err := t.api(ServiceGroupOpen, "get_role_in_group", req)
	if err != nil {
		return nil, err
	}

	response := new(RoleInGroupResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.UserIDList, nil
}

// ForbidSendMsgRequest 禁言设置与查询请求
type ForbidSendMsgRequest struct {
	GroupID       string   `json:"GroupId"`        // 群组ID
	MemberAccount []string `json:"Member_Account"` // 需要禁言的用户账号，最多支持500个账号
	ShutUpTime    int      // 需禁言时间，单位为秒，为0时表示取消禁言
}

// ForbidSendMsg 批量禁言和取消禁言
func (t *IM) ForbidSendMsg(groupID string, accounts []string, forbidTime int) error {
	if len(accounts) > 500 {
		return fmt.Errorf("一次请求最多操作500个成员")
	}

	req := &ForbidSendMsgRequest{
		GroupID:       groupID,
		MemberAccount: accounts,
		ShutUpTime:    forbidTime,
	}
	res, err := t.api(ServiceGroupOpen, "forbid_send_msg", req)
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

// ShuttedAccount 被禁言用户信息
type ShuttedAccount struct {
	MemberAccount string `json:"Member_Account"` // 用户ID
	ShuttedUntil  int64  // 禁言到的时间（使用UTC时间，即世界协调时间）
}

// GroupShuttedUinResponse 群组禁言用户查询响应
type GroupShuttedUinResponse struct {
	Response
	GroupID        string `json:"GroupId"`
	ShuttedUinList []*ShuttedAccount
}

// GetGroupShuttedUin 获取群组被禁言用户列表
func (t *IM) GetGroupShuttedUin(groupID string) ([]*ShuttedAccount, error) {
	req := &GroupID{GroupID: groupID}
	res, err := t.api(ServiceGroupOpen, "get_group_shutted_uin", req)
	if err != nil {
		return nil, err
	}

	response := new(GroupShuttedUinResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.ShuttedUinList, nil
}

// GroupMsgRequest 群组消息请求
type GroupMsgRequest struct {
	GroupID         string           `json:"GroupId"`
	FromAccount     string           `json:"From_Account,omitempty"` // 指定消息发送者（选填）
	Random          int64            // 32位随机数字，五分钟数字相同认为是重复消息
	MsgBody         []*MsgBodyItem   // 消息内容
	OfflinePushInfo *OfflinePushInfo `json:",omitempty"` // 离线推送信息配置
}

// SendGroupMsg 在群组中发送普通消息
func (t *IM) SendGroupMsg(groupID, from string, msg []*MsgBodyItem, pushInfo *OfflinePushInfo) error {
	req := &GroupMsgRequest{
		GroupID:         groupID,
		FromAccount:     from,
		Random:          util.StringToInt64(string(util.RandStr(6, util.KC_RAND_KIND_NUM))),
		MsgBody:         msg,
		OfflinePushInfo: pushInfo,
	}

	res, err := t.api(ServiceGroupOpen, "send_group_msg", req)
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

// SystemNotificationRequest 系统通知请求
type SystemNotificationRequest struct {
	GroupID         string   `json:"GroupId"` // 群组ID
	Content         string   // 系统通知内容
	ToMembersAccount []string `json:"ToMembers_Account,omitempty"` // 接收者群成员列表，不填或为空表示全员下发
}

// SendGroupSystemNotification 在群组中发送系统通知
func (t *IM) SendGroupSystemNotification(groupID, content string, toAccounts []string) error {
	req := &SystemNotificationRequest{
		GroupID:         groupID,
		Content:         content,
	}
	if len(toAccounts) > 0 {
		req.ToMembersAccount = toAccounts
	}
	res, err := t.api(ServiceGroupOpen, "send_group_system_notification", req)
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

// ChangeGroupOwnerRequest 转让群组请求
type ChangeGroupOwnerRequest struct {
	GroupID         string `json:"GroupId"`          // 要被转移的群组ID
	NewOwnerAccount string `json:"NewOwner_Account"` // 新群主ID
}

// ChangeGroupOwner 转让群组
func (t *IM) ChangeGroupOwner(groupID, newOwner string) error {
	req := &ChangeGroupOwnerRequest{
		GroupID:         groupID,
		NewOwnerAccount: newOwner,
	}
	res, err := t.api(ServiceGroupOpen, "change_group_owner", req)
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

// ImportGroup 导入群基础资料
func (t *IM) ImportGroup(req *GroupInfo) (string, error) {
	if req.Name == "" || req.Type == "" {
		return "", fmt.Errorf("群组名称与类型为必填项")
	}

	res, err := t.api(ServiceGroupOpen, "import_group", req)
	if err != nil {
		return "", err
	}

	response := new(CreateGroupResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return "", fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		if response.ErrorCode == 10021 {
			return req.GroupID, nil
		}
		return "", fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.GroupID.GroupID, nil
}

// ImportGroupMsgItem 导入群消息条目
type ImportGroupMsgItem struct {
	FromAccount string `json:"From_Account"`
	SendTime    int64
	Random      int64 `json:",omitempty"`
	MsgBody     []*MsgBodyItem
}

// ImportGroupMsgRequest 导入群消息请求
type ImportGroupMsgRequest struct {
	GroupID string                `json:"GroupId"` // 群组ID
	MsgList []*ImportGroupMsgItem // 消息列表，一次最多导入20条
}

// ImportGroupMsgResultItem 导入群消息结果条目
type ImportGroupMsgResultItem struct {
	Result  int
	MsgSeq  int
	MsgTime int64
}

// ImportGroupMsgResponse 导入群消息响应
type ImportGroupMsgResponse struct {
	Response
	ImportMsgResult []*ImportGroupMsgResultItem
}

// ImportGroupMsg 导入群消息
func (t *IM) ImportGroupMsg(groupID string, msg []*ImportGroupMsgItem) ([]*ImportGroupMsgResultItem, error) {
	req := &ImportGroupMsgRequest{
		GroupID: groupID,
		MsgList: msg,
	}

	res, err := t.api(ServiceGroupOpen, "import_group_msg", req)
	if err != nil {
		return nil, err
	}

	response := new(ImportGroupMsgResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.ImportMsgResult, nil
}

// ImportGroupMemberInfo 导入群成员信息
type ImportGroupMemberInfo struct {
	MemberAccount string `json:"Member_Account"` // 群成员帐号
	Role          string `json:",omitempty"`     // 导入成员的角色，目前只有Admin(可选)
	JoinTime      int64  `json:",omitempty"`     // 导入的成员入群时间（选填）
	UnreadMsgNum  int    `json:",omitempty"`     // 该成员的未读消息数（选填）
}

// ImportGroupMemberRequest 导入群成员请求
type ImportGroupMemberRequest struct {
	GroupID    string `json:"GroupId"` // 群组ID
	MemberList []*ImportGroupMemberInfo
}

// ImportGroupMemberResponseItem 导入群成员响应条目
type ImportGroupMemberResponseItem struct {
	MemberAccount string `json:"Member_Account"` // 群成员帐号
	Result        int    // 导入结果：0为失败；1为成功；2表示已经是群成员
}

// ImportGroupMemberResponse 导入群成员响应
type ImportGroupMemberResponse struct {
	Response
	MemberList []*ImportGroupMemberResponseItem
}

// ImportGroupMember 导入群成员
func (t *IM) ImportGroupMember(groupID string, members []*ImportGroupMemberInfo) ([]*ImportGroupMemberResponseItem, error) {
	req := &ImportGroupMemberRequest{
		GroupID:    groupID,
		MemberList: members,
	}
	res, err := t.api(ServiceGroupOpen, "import_group_member", req)
	if err != nil {
		return nil, err
	}

	response := new(ImportGroupMemberResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.MemberList, nil
}

// SetUnreadMsgNumRequest 设置成员未读消息计数请求
type SetUnreadMsgNumRequest struct {
	GroupID       string `json:"GroupId"`        // 要操作的群组ID
	MemberAccount string `json:"Member_Account"` // 要操作的群成员ID
	UnreadMsgNum  int    // 未读消息数
}

// SetUnreadMsgNum 设置成员未读消息计数
func (t *IM) SetUnreadMsgNum(groupID, account string, unread int) error {
	req := &SetUnreadMsgNumRequest{
		GroupID:       groupID,
		MemberAccount: account,
		UnreadMsgNum:  unread,
	}
	res, err := t.api(ServiceGroupOpen, "set_unread_msg_num", req)
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

// DeleteGroupMsgBySenderRequest 删除指定用户发送的消息请求
type DeleteGroupMsgBySenderRequest struct {
	GroupID       string `json:"GroupId"`        // 要删除消息的群组ID
	SenderAccount string `json:"Sender_Account"` // 被删除消息的发送者ID
}

// DeleteGroupMsgBySender 删除指定用户发送的消息
func (t *IM) DeleteGroupMsgBySender(groupID, account string) error {
	req := &DeleteGroupMsgBySenderRequest{
		GroupID:       groupID,
		SenderAccount: account,
	}
	res, err := t.api(ServiceGroupOpen, "delete_group_msg_by_sender", req)
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

// SearchGroupRequest 搜索群组请求
type SearchGroupRequest struct {
	Content        string          // 群组名关键字
	PageNum        int             `json:",omitempty"` // 页码
	GroupPerPage   int             `json:",omitempty"` // 每页数量
	ResponseFilter *ResponseFilter `json:",omitempty"` // 基础公开信息字段过滤器，指定需要获取的基础信息字段
}

// SearchGroupResponse 搜索群组响应
type SearchGroupResponse struct {
	Response
	GroupInfo   []*GroupInfo
	TotalRecord int
}

// SearchGroup 搜索群组
func (t *IM) SearchGroup(keyword string, page, pagesize int, filter []string) ([]*GroupInfo, int, error) {
	req := &SearchGroupRequest{
		Content:      keyword,
		PageNum:      page,
		GroupPerPage: pagesize,
	}

	if len(filter) > 0 {
		req.ResponseFilter = &ResponseFilter{GroupBasePublicInfoFilter: filter}
	}

	res, err := t.api(ServiceGroupOpen, "search_group", req)
	if err != nil {
		return nil, 0, err
	}

	response := new(SearchGroupResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, 0, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, 0, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.GroupInfo, response.TotalRecord, nil
}

// GroupMsgGetRequest 拉取群漫游消息请求
type GroupMsgGetRequest struct {
	GroupID      string `json:"GroupId"` // 要拉取漫游消息的群组ID
	ReqMsgNumber int    // 拉取的漫游消息的条数，目前一次请求最多返回20条漫游消息
	ReqMsgSeq    int64  `json:",omitempty"` // 拉取消息的最大seq
}

// GroupMsgItem 群组漫游消息单个条目
type GroupMsgItem struct {
	FromAccount  string         `json:"From_Account"` // 消息的发送者
	IsPlaceMsg   int            // 是否是空洞消息，当消息被删除或者消息过期后，MsgBody为空，这个字段为1
	MsgBody      []*MsgBodyItem // 消息内容
	MsgRandom    int            // 消息随机值
	MsgSeq       int64          // 消息seq，用来标识唯一消息，值越小发送的越早
	MsgTimeStamp int64          // 消息被发送的时间戳
}

// GroupMsgGetResponse 群漫游消息响应
type GroupMsgGetResponse struct {
	Response
	GroupID    string          `json:"GroupId"`
	IsFinished int             // 是否返回了请求区间的全部消息
	RspMsgList []*GroupMsgItem // 返回的消息列表
}

// GroupMsgGetSimple 拉取群漫游消息
func (t *IM) GroupMsgGetSimple(groupID string, limit int, msgseq int64) ([]*GroupMsgItem, bool, error) {
	if limit > 20 {
		return nil, false, fmt.Errorf("一次最多拉取20条漫游消息")
	}

	req := &GroupMsgGetRequest{
		GroupID:      groupID,
		ReqMsgNumber: limit,
		ReqMsgSeq:    msgseq,
	}
	res, err := t.api(ServiceGroupOpen, "group_msg_get_simple", req)
	if err != nil {
		return nil, false, err
	}

	response := new(GroupMsgGetResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, false, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, false, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	var finished bool
	if response.IsFinished == 1 {
		finished = true
	}

	return response.RspMsgList, finished, nil
}
