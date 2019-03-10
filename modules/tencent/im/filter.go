package im

// ResponseFilter 过滤器
type ResponseFilter struct {
	GroupBasePublicInfoFilter       []string `json:",omitempty"`                                 // 群组公开基础信息过滤器
	GroupBaseInfoFilter             []string `json:",omitempty"`                                 // 群组基础信息字段过滤器，和上面没什么区别，不同接口用的名字不一样
	MemberInfoFilter                []string `json:",omitempty"`                                 // 成员信息字段过滤器
	AppDefinedDataFilterGroup       []string `json:"AppDefinedDataFilter_Group,omitempty"`       // 群组维度的自定义字段过滤器
	AppDefinedDataFilterGroupMember []string `json:"AppDefinedDataFilter_GroupMember,omitempty"` // 群成员维度自定义字段过滤器
	MemberRoleFilter                []string `json:",omitempty"`                                 // 群成员身份过滤器
	SelfInfoFilter                  []string `json:",omitempty"`                                 // 用户在群组中的个人资料
}

var (
	// GroupBasePublicInfoFilter 群组公开基础信息过滤器
	GroupBasePublicInfoFilter = []string{
		"Type",
		"Name",
		"Introduction",
		"Notification",
		"FaceUrl",
		"Owner_Account",
		"CreateTime",
		"InfoSeq",
		"LastInfoTime",
		"LastMsgTime",
		"NextMsgSeq",
		"MemberNum",
		"MaxMemberNum",
		"ApplyJoinOption",
	}
	// GroupBaseInfoFilter 群组基础信息字段过滤器
	GroupBaseInfoFilter = []string{
		"Type",
		"Name",
		"Introduction",
		"Notification",
		"FaceUrl",
		"Owner_Account",
		"CreateTime",
		"InfoSeq",
		"LastInfoTime",
		"LastMsgTime",
		"NextMsgSeq",
		"MemberNum",
		"MaxMemberNum",
		"ApplyJoinOption",
	}
	// MemberInfoFilter 成员信息字段过滤器
	MemberInfoFilter = []string{
		"Role",
		"JoinTime",
		"MsgSeq",
		"MsgFlag",
		"LastSendMsgTime",
		"ShutUpUntil",
		"NameCard",
	}
	// SelfInfoFilter 用户在群组中的个人资料
	SelfInfoFilter = []string{
		"Role",
		"JoinTime",
		"MsgSeq",
		"MsgFlag",
		"LastSendMsgTime",
		"ShutUpUntil",
		"NameCard",
	}
	// AppDefinedDataFilterGroup 群组维度的自定义字段过滤器
	AppDefinedDataFilterGroup = []string{}
	// AppDefinedDataFilterGroupMember 群成员维度自定义字段过滤器
	AppDefinedDataFilterGroupMember = []string{}
	// MemberRoleFilter 群成员身份过滤器
	MemberRoleFilter = []string{GroupRoleOwner, GroupRoleAdmin, GroupRoleMember}
)
