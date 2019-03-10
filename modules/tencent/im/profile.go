package im

import (
	"encoding/json"
	"fmt"
)

// ServiceProfile 资料管理服务名
const ServiceProfile = "profile"

// 标准资料Tag字段
const (
	// ProfileNick 昵称
	ProfileNick = "Tag_Profile_IM_Nick"
	// ProfileGender 性别
	ProfileGender = "Tag_Profile_IM_Gender"
	// ProfileBirthDay 生日
	ProfileBirthDay = "Tag_Profile_IM_BirthDay"
	// ProfileLocation 所在地
	ProfileLocation = "Tag_Profile_IM_Location"
	// ProfileSelfSignature 个性签名
	ProfileSelfSignature = "Tag_Profile_IM_SelfSignature"
	// ProfileAllowType 加好友验证方式
	ProfileAllowType = "Tag_Profile_IM_AllowType"
	// ProfileLanguage 语言
	ProfileLanguage = "Tag_Profile_IM_Language"
	// ProfileImage 头像URL
	ProfileImage = "Tag_Profile_IM_Image"
	// ProfileMsgSettings 消息设置: Bit0：置0表示接收消息，置1则不接收消息
	ProfileMsgSettings = "Tag_Profile_IM_MsgSettings"
)

// 自定义资料Tag字段
const (
	// CustomProfileCompanyID 企业ID
	CustomProfileCompanyID = "Tag_Profile_Custom_Comid"
	// CustomProfileDepartmentID 门店ID
	CustomProfileDepartmentID = "Tag_Profile_Custom_Depid"
	// CustomProfileAccountRole 账号类型
	CustomProfileAccountRole = "Tag_Profile_Custom_Role"
	// CustomProfileUID 用户ID
	CustomProfileUID = "Tag_Profile_Custom_UID"
)

const (
	// GenderUnknown 没设置性别
	GenderUnknown = "Gender_Type_Unknown"
	// GenderFemale 女性
	GenderFemale = "Gender_Type_Female"
	// GenderMale 男性
	GenderMale = "Gender_Type_Male"
)

const (
	// AllowTypeNeedConfirm 需要经过自己确认才能添加自己为好友
	AllowTypeNeedConfirm = "AllowType_Type_NeedConfirm"
	// AllowTypeAllowAny 允许任何人添加自己为好友
	AllowTypeAllowAny = "AllowType_Type_AllowAny"
	// AllowTypeDenyAny 不允许任何人添加自己为好友
	AllowTypeDenyAny = "AllowType_Type_DenyAny"
)

// ProfileItem 资料单条项目
type ProfileItem struct {
	Tag   string // 资料字段名称
	Value string // 资料字段值
}

// UserProfileItem 用户资料结构化信息
type UserProfileItem struct {
	ToAccount   string         `json:"To_Account"`
	ProfileItem []*ProfileItem // 用户的资料对象列表
	ResultCode  int            // 单个用户的结果，0表示正确，非0表示错误
	ResultInfo  string         // 单个用户的结果详细信息
}

// SetProfileRequest 设置资料请求
type SetProfileRequest struct {
	FromAccount
	ProfileItem []*ProfileItem
}

// SetProfile 设置资料
func (t *IM) SetProfile(account string, profile map[string]string) error {
	var profileList = make([]*ProfileItem, 0)
	for k, v := range profile {
		item := new(ProfileItem)
		item.Tag = k
		item.Value = v
		profileList = append(profileList, item)
	}
	req := &SetProfileRequest{
		FromAccount: FromAccount{FromAccount: account},
		ProfileItem: profileList,
	}
	res, err := t.api(ServiceProfile, "portrait_set", req)
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

// GetProfileRequest 拉取资料请求
type GetProfileRequest struct {
	ToAccount
	TagList []string // 要拉取的资料对象的名称
}

// GetProfileResponse 拉取资料响应
type GetProfileResponse struct {
	Response
	UserProfileItem []*UserProfileItem // 用户对象列表
	FailAccount     []string           `json:"Fail_Account,omitempty"`
	InvalidAccount  []string           `json:"Invalid_Account,omitempty"`
}

// GetProfile 拉取资料
func (t *IM) GetProfile(accounts []string, tags []string) ([]*UserProfileItem, []string, []string, error) {
	req := &GetProfileRequest{
		ToAccount: ToAccount{ToAccount: accounts},
		TagList:   tags,
	}
	res, err := t.api(ServiceProfile, "portrait_get", req)
	if err != nil {
		return nil, nil, nil, err
	}

	response := new(GetProfileResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, nil, nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.UserProfileItem, response.FailAccount, response.InvalidAccount, nil
}
