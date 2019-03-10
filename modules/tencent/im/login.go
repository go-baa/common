package im

import (
	"encoding/json"
	"fmt"
)

const (
	// ServiceOpenLogin 独立账号服务名
	ServiceOpenLogin = "im_open_login_svc"
	// ServiceRegistration 托管账号服务名
	ServiceRegistration = "registration_service"
)

const (
	// IdentifierTypeMobile 手机类型用户名
	IdentifierTypeMobile = iota + 1
	// IdentifierTypeEmail 邮箱类型用户名
	IdentifierTypeEmail
	// IdentifierTypeString 字符串类型用户名
	IdentifierTypeString
)

// Account 账号基本信息
type Account struct {
	Identifier string // 用户名，长度不超过 32 字节
	Nick       string // 用户昵称
	FaceURL    string `json:"FaceUrl"` // 用户头像URL
}

// MultiImportRequest 批量导入请求
type MultiImportRequest struct {
	Accounts []string
}

// MultiImportResponse 批量导入响应
type MultiImportResponse struct {
	Response
	FailAccounts []string // 导入失败的帐号列表
}

// RegisterAccount 托管模式导入账号
type RegisterAccount struct {
	Identifier     string // 为用户申请同步的帐号，长度为4-24个字符
	IdentifierType int    // Identifier的类型，1:手机号(国家码-手机号) 2:邮箱 3:字符串帐号
	Password       string // Identifier的密码，长度为8-16个字符
}

// AccountImport 独立模式帐号导入
// nickname, faceURL 选填
func (t *IM) AccountImport(identifier, nickname, faceURL string) error {
	req := &Account{
		Identifier: identifier,
		Nick:       nickname,
		FaceURL:    faceURL,
	}
	res, err := t.api(ServiceOpenLogin, "account_import", req)
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

// MultiAccountImport 独立模式帐号批量导入, 返回导入失败的账号名
func (t *IM) MultiAccountImport(accounts []string) ([]string, error) {
	if len(accounts) > 100 {
		return nil, fmt.Errorf("单次最多导入100个用户名")
	}

	req := new(MultiImportRequest)
	req.Accounts = accounts
	res, err := t.api(ServiceOpenLogin, "multiaccount_import", req)
	if err != nil {
		return nil, err
	}

	response := new(MultiImportResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.FailAccounts, nil
}

// RegisterAccountV1 托管模式帐号导入
func (t *IM) RegisterAccountV1(identifier string, identifierType int, password string) error {
	req := &RegisterAccount{
		Identifier:     identifier,
		IdentifierType: identifierType,
		Password:       password,
	}
	res, err := t.api(ServiceRegistration, "register_account_v1", req)
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

// Kick 失效帐号登录态
func (t *IM) Kick(identifier string) error {
	req := new(Account)
	req.Identifier = identifier
	res, err := t.api(ServiceOpenLogin, "kick", req)
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
