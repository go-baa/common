package qcloud

import (
	"fmt"

	"git.code.tencent.com/xinhuameiyu/common/modules/sms/base"
	"git.code.tencent.com/xinhuameiyu/common/modules/tencent/sms"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

// SMS 腾讯云短信
type SMS struct {
}

// Config 配置
type Config struct {
	Name     string
	AppID    string
	AppKey   string
	Sign     string
	SMSTplID int
}

func getConfig(name string) *Config {
	c := new(Config)

	c.Name = name

	c.AppID = setting.Config.MustString("sms."+name+".tencent.appid", "")
	if c.AppID == "" {
		c.AppID = setting.Config.MustString("sms.tencent.appid", "")
	}
	if c.AppID == "" {
		log.Errorf("短信发送失败：腾讯 短信模板配置缺少 appid %s\n", name)
		return nil
	}

	c.AppKey = setting.Config.MustString("sms."+name+".tencent.appkey", "")
	if c.AppKey == "" {
		c.AppKey = setting.Config.MustString("sms.tencent.appkey", "")
	}
	if c.AppKey == "" {
		log.Errorf("短信发送失败：腾讯云 短信模板配置缺少 appkey %s\n", name)
		return nil
	}

	c.Sign = setting.Config.MustString("sms."+name+".tencent.sign_name", "")
	if c.Sign == "" {
		c.Sign = setting.Config.MustString("sms.tencent.sign_name", "")
	}
	if c.Sign == "" {
		log.Errorf("短信发送失败：腾讯云 短信模板配置缺少 sign_name %s\n", name)
		return nil
	}

	c.SMSTplID = setting.Config.MustInt("sms."+name+".tencent.sms_tplid", 0)
	if c.SMSTplID == 0 {
		c.SMSTplID = setting.Config.MustInt("sms.tencent.sms_tplid", 0)
	}
	if c.SMSTplID == 0 {
		log.Errorf("短信发送失败：腾讯云 短信模板配置缺少 sms_tplid %s\n", name)
		return nil
	}

	return c
}

// SendSMSCode 发送短信验证码
func (t *SMS) SendSMSCode(mobile string, code string) (string, error) {
	c := getConfig("code")
	if c == nil {
		return "", fmt.Errorf("短信发送失败：腾讯云 %s 配置初始化失败", "code")
	}
	client := sms.New(c.AppID, c.AppKey)
	sid, err := client.SendSmsTplCode("86", mobile, code, c.Sign, c.SMSTplID)
	if err != nil {
		log.Errorf("短信发送失败：腾讯云 %s Error：%s\n", c.Name, err.Error())
		return "", err
	}

	return sid, nil
}

// SendVoiceCode 发送语音短信验证码
func (t *SMS) SendVoiceCode(mobile string, code string) (string, error) {
	c := getConfig("code")
	if c == nil {
		return "", fmt.Errorf("短信发送失败：腾讯云 %s 配置初始化失败", "code")
	}
	client := sms.New(c.AppID, c.AppKey)
	sid, err := client.SendVoiceCode("86", mobile, code)
	if err != nil {
		log.Errorf("短信发送失败：腾讯云 %s Error：%s\n", c.Name, err.Error())
		return "", nil
	}

	return sid, nil
}

func init() {
	base.Register("qcloud", new(SMS))
}
