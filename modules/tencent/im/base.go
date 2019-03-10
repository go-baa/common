package im

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-baa/common/util"
	// 导入cache的Redis适配器
	_ "github.com/go-baa/cache/redis"
)

// IMGateway REST API地址
const IMGateway = "https://console.tim.qq.com/v4/"

const (
	// ActionStatusOK 请求处理结果成功
	ActionStatusOK = "OK"
	// ActionStatusFail 请求处理结果失败
	ActionStatusFail = "FAIL"
)

// IM 腾讯云通信
type IM struct {
	sdkappid   string
	identifier string
	usersig    string
}

// EmptyRequest 空请求
type EmptyRequest struct{}

// Response 基本响应
type Response struct {
	ActionStatus string // 请求处理的结果，OK表示处理成功，FAIL表示失败
	ErrorCode    int    // 错误码
	ErrorInfo    string // 错误信息
}

// IMApiFrequent Api调用频率
var IMApiFrequent = map[string]int{
	ServiceConfig + "_setnospeaking":                     100,
	ServiceConfig + "_getnospeaking":                     100,
	ServiceDirtyWords + "_get":                           100,
	ServiceDirtyWords + "_add":                           100,
	ServiceDirtyWords + "_delete":                        100,
	ServiceGroupOpen + "_get_appid_group_list":           100,
	ServiceGroupOpen + "_create_group":                   100,
	ServiceGroupOpen + "_get_group_info":                 100,
	ServiceGroupOpen + "_get_group_member_info":          100,
	ServiceGroupOpen + "_modify_group_base_info":         100,
	ServiceGroupOpen + "_add_group_member":               100,
	ServiceGroupOpen + "_delete_group_member":            100,
	ServiceGroupOpen + "_modify_group_member_info":       100,
	ServiceGroupOpen + "_destroy_group":                  100,
	ServiceGroupOpen + "_get_joined_group_list":          100,
	ServiceGroupOpen + "_get_role_in_group":              100,
	ServiceGroupOpen + "_forbid_send_msg":                100,
	ServiceGroupOpen + "_get_group_shutted_uin":          100,
	ServiceGroupOpen + "_send_group_msg":                 100,
	ServiceGroupOpen + "_send_group_system_notification": 100,
	ServiceGroupOpen + "_change_group_owner":             100,
	ServiceGroupOpen + "_import_group":                   100,
	ServiceGroupOpen + "_import_group_msg":               100,
	ServiceGroupOpen + "_import_group_member":            100,
	ServiceGroupOpen + "_set_unread_msg_num":             100,
	ServiceGroupOpen + "_delete_group_msg_by_sender":     100,
	ServiceGroupOpen + "_search_group":                   100,
	ServiceGroupOpen + "_group_msg_get_simple":           100,
	ServiceOpenLogin + "_account_import":                 1000,
	ServiceOpenLogin + "_multiaccount_import":            10,
	ServiceOpenLogin + "_kick":                           1000,
	ServiceRegistration + "_register_account_v1":         1000,
	ServiceOpenMSG + "_get_history":                      100,
	ServiceOpenIM + "_sendmsg":                           100,
	ServiceOpenIM + "_batchsendmsg":                      100,
	ServiceOpenIM + "_importmsg":                         100,
	ServiceOpenIM + "_im_get_push_report":                100,
	ServiceOpenIM + "_im_set_attr_name":                  100,
	ServiceOpenIM + "_im_get_attr_name":                  100,
	ServiceOpenIM + "_im_set_attr":                       100,
	ServiceOpenIM + "_im_remove_attr":                    100,
	ServiceOpenIM + "_im_get_attr":                       100,
	ServiceOpenIM + "_im_add_tag":                        100,
	ServiceOpenIM + "_im_remove_tag":                     100,
	ServiceOpenIM + "_im_remove_all_tags":                100,
	ServiceOpenIM + "_querystate":                        100,
	ServiceOpenConfig + "_getappinfo":                    100,
	ServiceProfile + "_portrait_set":                     100,
	ServiceProfile + "_portrait_get":                     100,
	ServiceSNS + "_friend_add":                           100,
	ServiceSNS + "_friend_get_all":                       100,
	ServiceSNS + "_friend_delete":                        100,
	ServiceSNS + "_friend_delete_all":                    100,
	ServiceSNS + "_friend_check":                         100,
	ServiceSNS + "_black_list_check":                     100,
}

// New 获取IM实例
func New(appid, identifier, usersig string) (*IM, error) {
	ins := new(IM)
	if appid == "" {
		return nil, fmt.Errorf("Invalid appid")
	}
	ins.sdkappid = appid

	if identifier == "" {
		return nil, fmt.Errorf("Invalid identifier")
	}
	ins.identifier = identifier

	if usersig == "" {
		return nil, fmt.Errorf("Invalid usersig")
	}
	ins.usersig = usersig

	return ins, nil
}

// getQueryParams 组装查询参数，注意不要编码
func (t *IM) getQueryParams() string {
	params := map[string]string{
		"usersig":     t.usersig,
		"identifier":  t.identifier,
		"sdkappid":    t.sdkappid,
		"random":      string(util.RandStr(32, util.KC_RAND_KIND_NUM)),
		"contenttype": "json",
	}
	s := make([]string, 0, len(params))
	for k, v := range params {
		item := k + "=" + v
		s = append(s, item)
	}

	return strings.Join(s, "&")
}

// api 调用api
func (t *IM) api(service, command string, reqdata interface{}) ([]byte, error) {
	url := IMGateway + service + "/" + command
	query := t.getQueryParams()
	reqBodyJSON, err := json.Marshal(reqdata)

	if err != nil {
		return nil, fmt.Errorf("请求数据json编码错误:%v", err)
	}

	res, err := t.request(url, query, reqBodyJSON)
	if err != nil {
		time.Sleep(time.Second)
		return nil, fmt.Errorf("请求失败:%v", err)
	}

	return res, nil
}

// request http请求
func (t *IM) request(url string, query string, reqBody []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url+"?"+query, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// 超时设置
	client := new(http.Client)
	client.Timeout = time.Second * 60

	// https 支持
	if strings.HasPrefix(url, "https") {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	// 执行请求
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 处理响应
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}
