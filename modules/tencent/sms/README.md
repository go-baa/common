package sms

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var client *Sms

func init() {
	client = New("1400024347", "xxxx")
}

func TestSendSmsTplCode(t *testing.T) {
	Convey("测试发送短信验证码", t, func() {
		_, err := client.SendSmsTplCode("86", "18600362795", "8876", "药视通", 18371)
		So(err, ShouldBeNil)
	})
}

func TestSendVoiceCode(t *testing.T) {
	Convey("测试发送语音验证码", t, func() {
		_, err := client.SendVoiceCode("86", "18610367118", "8876")
		So(err, ShouldBeNil)
	})
}
