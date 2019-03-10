package cps

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// TagsResponse 标签查询响应
type TagsResponse struct {
	Response
	TagInfos tagInfo
}

type tagInfo struct {
	TagInfo []*tagInfoItem
}

type tagInfoItem struct {
	TagName string
}

// TagItem 标签
type TagItem struct {
	Tag string `json:"tag"`
}

// MultiTag 多标签交集
type MultiTag struct {
	And []*TagItem `json:"and"`
}

const (
	// KeyTypeDevice 设备
	KeyTypeDevice = "DEVICE"
	// KeyTypeAccount 账号
	KeyTypeAccount = "ACCOUNT"
	// KeyTypeAlias 别名
	KeyTypeAlias = "ALIAS"
)

// BindTag 绑定标签
func (t *CPS) BindTag(keyType string, clientKey []string, tags []string) (*Response, error) {
	if len(clientKey) > 1000 {
		return nil, fmt.Errorf("每次绑定操作不能超过1000个账号")
	}
	if len(tags) > 10 {
		return nil, fmt.Errorf("每次操作标签不能超过10个")
	}
	params := t.getCommonParams()
	params["Action"] = "BindTag"
	params["KeyType"] = keyType
	params["ClientKey"] = strings.Join(clientKey, ",")
	params["TagName"] = strings.Join(tags, ",")

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

// UnbindTag 解绑定标签
func (t *CPS) UnbindTag(keyType string, clientKey []string, tags []string) (*Response, error) {
	if len(clientKey) > 1000 {
		return nil, fmt.Errorf("每次解绑操作不能超过1000个账号")
	}
	if len(tags) > 10 {
		return nil, fmt.Errorf("每次操作标签不能超过10个")
	}
	params := t.getCommonParams()
	params["Action"] = "UnbindTag"
	params["KeyType"] = keyType
	params["ClientKey"] = strings.Join(clientKey, ",")
	params["TagName"] = strings.Join(tags, ",")

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

// ListTags tag列表
func (t *CPS) ListTags() ([]string, error) {
	params := t.getCommonParams()
	params["Action"] = "ListTags"

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
	ret := new(TagsResponse)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, err
	}

	var tags = make([]string, 0)
	for _, v := range ret.TagInfos.TagInfo {
		tags = append(tags, v.TagName)
	}

	return tags, nil
}

// QueryTags 标签查询
func (t *CPS) QueryTags(keyType, clientKey string) ([]string, error) {
	params := t.getCommonParams()
	params["Action"] = "QueryTags"
	params["KeyType"] = keyType
	params["ClientKey"] = clientKey

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
	ret := new(TagsResponse)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, err
	}

	var tags = make([]string, 0)
	for _, v := range ret.TagInfos.TagInfo {
		tags = append(tags, v.TagName)
	}

	return tags, nil
}

// MergeTag 合并多个标签
func (t *CPS) MergeTag(tags []string) (string, error) {
	ret := new(MultiTag)
	for _, v := range tags {
		t := new(TagItem)
		t.Tag = v
		ret.And = append(ret.And, t)
	}

	retJSON, err := json.Marshal(ret)
	if err != nil {
		return "", nil
	}

	return string(retJSON), nil
}
