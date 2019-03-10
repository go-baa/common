package sms

import (
	"fmt"
	"time"

	_ "github.com/go-baa/common/modules/sms/aliyun"
	"github.com/go-baa/common/modules/sms/base"
	_ "github.com/go-baa/common/modules/sms/juhe"
	_ "github.com/go-baa/common/modules/sms/qcloud"
	"github.com/go-baa/log"
)

// SendSMSCode 发送短信验证码
func SendSMSCode(mobile string, code string) (ret string, err error) {
	var success bool
	var count int
	var provider base.SMSProvider

	// 首先尝试阿里云接口
	if provider = base.GetSMSProvider("aliyun"); provider != nil {
		for i := 0; i < 2; i++ {
			count++
			ret, err = provider.SendSMSCode(mobile, code)
			if err == nil {
				log.Debugf("短信发送成功：阿里云 %s\n", ret)
				success = true
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
	}

	if success {
		return
	}

	// 再尝试腾讯云接口
	if provider = base.GetSMSProvider("qcloud"); provider != nil {
		for i := 0; i < 2; i++ {
			count++
			ret, err = provider.SendSMSCode(mobile, code)
			if err == nil {
				log.Debugf("短信发送成功：腾讯云 %s\n", ret)
				success = true
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
	}

	if count == 0 {
		log.Errorf("短信发送失败：没有可用的短信配置\n")
		err = fmt.Errorf("没有可用的短信配置")
	}

	return
}

// SendVoiceCode 发送语音短信验证码
func SendVoiceCode(mobile string, code string) (ret string, err error) {
	var success bool
	var count int
	var provider base.VoiceProvider

	// 再尝试腾讯云接口
	if provider = base.GetVoiceProvider("qcloud"); provider != nil {
		for i := 0; i < 2; i++ {
			count++
			ret, err = provider.SendVoiceCode(mobile, code)
			if err == nil {
				log.Debugf("短信发送成功：腾讯云 %s\n", ret)
				success = true
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
	}

	if success {
		return
	}

	if count == 0 {
		log.Errorf("语音短信发送失败：没有可用的语音短信配置\n")
		err = fmt.Errorf("没有可用的语音短信配置")
	}

	return
}
