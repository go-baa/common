package qihu360

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
)

const (
	LangCN = "zh_CN"
	LangTW = "zh_TW"
	LangEN = "en"
)

// APIURL 请求地址
//const APIURL = "http://demo.api.360kan.com/"
//const APIHost = "demo.api.360kan.com"
const APIHost = ""

const APIURL = "http://openapi.k.360kan.com/"

// APIRequestTimeout 请求超时时间
var APIRequestTimeout = 10

// ErrorResult 错误结果
type ErrorResult struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

// QihuTokenInfo 360的token数据
type QihuTokenInfo struct {
	Status   int          `json:"status"`
	Token    string       `json:"token"`
	Userinfo QihuUserInfo `json:"userinfo"`
}

// QihuUserInfo 供应商账号信息
type QihuUserInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

// GetToken 获取360的TOKEN
func GetToken(username, passwd string) (*QihuTokenInfo, *ErrorResult) {
	if username == "" || passwd == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}
	params := make(map[string]interface{})
	params["username"] = username
	params["password"] = util.MD5(passwd)
	data, err := util.HTTPPostJSON(APIURL+"mgrvideo/jwt/token", params, APIRequestTimeout)
	if err != nil {
		log.Errorf("获取360 token 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取360 token 失败"}
	}

	ret := new(QihuTokenInfo)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析360 token 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析360 token 响应失败"}
	}
	if ret.Status != 0 {
		return nil, &ErrorResult{Message: "获取TOKEN错误"}
	}

	return ret, nil
}

// RefreshToken 获取360的TOKEN
func RefreshToken(token string) (*QihuTokenInfo, *ErrorResult) {
	if token == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}
	client := &http.Client{}
	//提交请求
	reqest, err := http.NewRequest("GET", APIURL+"mgrvideo/jwt/token/refresh", nil)
	//增加header选项
	reqest.Header.Add("Authorization", token)
	if err != nil {
		log.Errorf("刷新360 token 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "刷新360 token 失败"}
	}
	//处理返回结果
	resp, err := client.Do(reqest)
	if err != nil {
		log.Errorf("刷新360 token 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "刷新360 token 失败"}
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("刷新360 token 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "刷新360 token 失败"}
	}

	ret := new(QihuTokenInfo)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("刷新360 token 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "刷新360 token 响应失败"}
	}
	if ret.Status != 0 {
		return nil, &ErrorResult{Message: "刷新TOKEN错误"}
	}

	return ret, nil
}
