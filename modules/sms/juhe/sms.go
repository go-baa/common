package juhe

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"git.code.tencent.com/xinhuameiyu/common/modules/sms/base"
	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

// SMS 聚合数据短信
type SMS struct {
}

// Config 配置
type Config struct {
	Name        string
	Key         string
	TplID       string
	TplValueRaw string
	Timeout     int
}

// juheSMSAPIResponse 聚合API的返回数据结构
type response struct {
	ErrorCode int            `json:"error_code"` // 返回码：0 表示成功
	Reason    string         `json:"reason"`     // 返回码说明
	Result    responseResult `json:"result"`
}

// responseResult 聚合API的处理结果结构
type responseResult struct {
	Count int    `json:"count"` // 发送数量
	Fee   int    `json:"fee"`   // 扣除条数
	Sid   string `json:"sid"`   // 短信ID
}

func getConfig(name string) *Config {
	c := new(Config)

	c.Name = name

	c.Key = setting.Config.MustString("sms."+name+".juhe.key", "")
	if c.Key == "" {
		c.Key = setting.Config.MustString("sms.juhe.key", "")
	}
	if c.Key == "" {
		log.Errorf("短信发送失败：聚合数据 短信模板配置缺少 key %s\n", name)
		return nil
	}

	c.TplID = setting.Config.MustString("sms."+name+".juhe.tpl_id", "")
	if c.TplID == "" {
		log.Errorf("短信发送失败：聚合数据 短信模板配置缺少 tpl_id %s\n", name)
		return nil
	}

	c.TplValueRaw = setting.Config.MustString("sms."+name+".juhe.tpl_value_raw", "")
	if c.TplValueRaw == "" {
		log.Errorf("短信发送失败：聚合数据 短信模板配置缺少 tpl_value_raw %s\n", name)
		return nil
	}

	c.Timeout = setting.Config.MustInt("sms."+name+".juhe.timeout", 0)
	if c.Timeout == 0 {
		c.Timeout = setting.Config.MustInt("sms.juhe.timeout", 0)
	}
	if c.Timeout == 0 {
		c.Timeout = 3
	}

	return c
}

func renderTpl(c *Config, params map[string]string) (string, error) {
	tpl := c.TplValueRaw

	for k, v := range params {
		tpl = strings.Replace(tpl, "{"+k+"}", v, -1)
	}

	return url.QueryEscape(tpl), nil
}

// SendSMSCode 发送短信验证码
func (t *SMS) SendSMSCode(mobile string, code string) (string, error) {
	c := getConfig("code")
	if c == nil {
		return "", fmt.Errorf("短信发送失败：聚合数据 %s 配置初始化失败", "code")
	}
	tplValue, err := renderTpl(c, map[string]string{
		"code": code,
	})
	if err != nil {
		return "", err
	}

	body, err := util.HTTPGet("http://v.juhe.cn/sms/send?mobile="+mobile+"&tpl_id="+c.TplID+"&tpl_value="+tplValue+"&key="+c.Key, c.Timeout)
	if err != nil {
		log.Errorf("短信发送失败：聚合数据 %s %s\n", mobile, err)
		return "", fmt.Errorf("网络异常")
	}

	// DEBUG
	if setting.Debug {
		log.Debugln(string(body))
	}

	var r *response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}

	if r.ErrorCode == 0 {
		return r.Result.Sid, nil
	}

	log.Errorf("短信发送失败：聚合数据 %s [%d] %s", mobile, r.ErrorCode, r.Reason)
	return "", fmt.Errorf("短信发送失败：[%d] %s", r.ErrorCode, r.Reason)
}

func init() {
	base.Register("juhe", new(SMS))
}
