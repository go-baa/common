package im

import (
	"encoding/json"
	"fmt"
)

// ServiceDirtyWords 脏字管理服务名
const ServiceDirtyWords = "openim_dirty_words"

// DirtyWordsList 脏字列表请求
type DirtyWordsList struct {
	DirtyWordsList []string
}

// DirtyWordsQueryResponse 脏字查询响应
type DirtyWordsQueryResponse struct {
	Response
	DirtyWordsList
}

// GetDirtyWords 查询自定义脏字
func (t *IM) GetDirtyWords() ([]string, error) {
	req := new(EmptyRequest)
	res, err := t.api(ServiceDirtyWords, "get", req)
	if err != nil {
		return nil, err
	}

	response := new(DirtyWordsQueryResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("code:%d, info: %s", response.ErrorCode, response.ErrorInfo)
	}

	return response.DirtyWordsList.DirtyWordsList, nil
}

// AddDirtyWords 添加自定义脏字
func (t *IM) AddDirtyWords(wordsList []string) error {
	if len(wordsList) == 0 {
		return fmt.Errorf("脏字列表不能为空")
	}

	req := new(DirtyWordsList)
	req.DirtyWordsList = wordsList
	res, err := t.api(ServiceDirtyWords, "add", req)
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

// DeleteDirtyWords 删除自定义脏字
func (t *IM) DeleteDirtyWords(wordsList []string) error {
	if len(wordsList) == 0 {
		return fmt.Errorf("脏字列表不能为空")
	}

	req := new(DirtyWordsList)
	req.DirtyWordsList = wordsList
	res, err := t.api(ServiceDirtyWords, "delete", req)
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
