package aliyun

import (
	"fmt"

	"git.code.tencent.com/xinhuameiyu/common/modules/sms/base"
	"github.com/denverdino/aliyungo/sms"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

// SMS 阿里云短信
type SMS struct {
}

// Config 配置
type Config struct {
	Name            string
	AppKey          string
	Secret          string
	SMSFreeSignName string
	SMSTemplateCode string
	Timeout         int
}

func getConfig(name string) *Config {
	c := new(Config)

	c.Name = name

	c.AppKey = setting.Config.MustString("sms."+name+".ali.app_key", "")
	if c.AppKey == "" {
		c.AppKey = setting.Config.MustString("sms.ali.app_key", "")
	}
	if c.AppKey == "" {
		log.Errorf("短信发送失败：阿里云 短信模板配置缺少 app_key %s\n", name)
		return nil
	}

	c.Secret = setting.Config.MustString("sms."+name+".ali.secret", "")
	if c.Secret == "" {
		c.Secret = setting.Config.MustString("sms.ali.secret", "")
	}
	if c.Secret == "" {
		log.Errorf("短信发送失败：阿里云 短信模板配置缺少 secret %s\n", name)
		return nil
	}

	c.SMSFreeSignName = setting.Config.MustString("sms."+name+".ali.sms_free_sign_name", "")
	if c.SMSFreeSignName == "" {
		c.SMSFreeSignName = setting.Config.MustString("sms.ali.sms_free_sign_name", "")
	}
	if c.SMSFreeSignName == "" {
		log.Errorf("短信发送失败：阿里云 短信模板配置缺少 sms_free_sign_name %s\n", name)
		return nil
	}

	c.SMSTemplateCode = setting.Config.MustString("sms."+name+".ali.sms_template_code", "")
	if c.SMSTemplateCode == "" {
		log.Errorf("短信发送失败：阿里云 短信模板配置缺少 sms_template_code %s\n", name)
		return nil
	}

	c.Timeout = setting.Config.MustInt("sms."+name+".ali.timeout", 0)
	if c.Timeout == 0 {
		c.Timeout = setting.Config.MustInt("sms.ali.timeout", 0)
	}
	if c.Timeout == 0 {
		c.Timeout = 3
	}

	return c
}

// SendSMSCode 发送短信验证码
func (t *SMS) SendSMSCode(mobile string, code string) (string, error) {
	c := getConfig("code")
	if c == nil {
		return "", fmt.Errorf("短信发送失败：阿里云 %s 配置初始化失败", "code")
	}
	client := sms.NewDYSmsClient(c.AppKey, c.Secret)

	response, err := client.SendSms(&sms.SendSmsArgs{
		PhoneNumbers:  mobile,
		SignName:      c.SMSFreeSignName,
		TemplateCode:  c.SMSTemplateCode,
		TemplateParam: `{"code": ` + code + `}`,
	})
	if err != nil {
		log.Errorf("短信发送失败：阿里云 %s Error：%s\n", c.Name, err.Error())
		return "", err
	}

	return response.RequestId, nil
}

func init() {
	base.Register("aliyun", new(SMS))
}
