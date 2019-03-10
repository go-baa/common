package sms

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"git.code.tencent.com/xinhuameiyu/common/util"
)

// Sms 短信通知接口
type Sms struct {
	gateway string
	appid   string
	appkey  string
}

type tel struct {
	Nationcode string `json:"nationcode"` // 国家代码
	Mobile     string `json:"mobile"`     // 手机号码
}

type smsTplCode struct {
	Tel    tel      `json:"tel"`
	Sign   string   `json:"sign"`   // 短信签名，如果使用默认签名，该字段可缺省
	TplID  int      `json:"tpl_id"` // 业务在控制台审核通过的模板ID
	Params []string `json:"params"` // 假定这个模板为：您的{1}是{2}，请于{3}分钟内填写。如非本人操作，请忽略本短信。
	Sig    string   `json:"sig"`    // 签名
	Time   int64    `json:"time"`   // unix时间戳，请求发起时间，如果和系统时间相差超过10分钟则会返回失败
	Extend string   `json:"extend"` // 通道扩展码，可选字段，默认没有开通(需要填空)。
	Ext    string   `json:"ext"`    // 用户的session内容，腾讯server回包中会原样返回，可选字段，不需要就填空。
	random string   // 随机数
}

type voiceCode struct {
	Tel       tel    `json:"tel"`
	Msg       string `json:"msg"`       // 验证码
	Playtimes int    `json:"playtimes"` // 语音播放次数，默认为2，不超过3
	Sig       string `json:"sig"`       // 签名
	Time      int64  `json:"time"`      // unix时间戳，请求发起时间，如果和系统时间相差超过10分钟则会返回失败
	Ext       string `json:"ext"`       // 用户的session内容，腾讯server回包中会原样返回，可选字段，不需要就填空。
	random    string // 随机数
}

type resultSms struct {
	Result int    `json:"result"` // 0表示成功(计费依据)，非0表示失败
	ErrMsg string `json:"errmsg"` // result非0时的具体错误信息
	Ext    string `json:"ext"`    // 用户的session内容，腾讯server回包中会原样返回
	Sid    string `json:"sid"`    // 标识本次发送id，标识一次短信下发记录
	Fee    int    `json:"fee"`    // 短信计费的条数
}

type resultVoice struct {
	Result int    `json:"result"` // 0表示成功(计费依据)，非0表示失败
	ErrMsg string `json:"errmsg"` // result非0时的具体错误信息
	Ext    string `json:"ext"`    // 用户的session内容，腾讯server回包中会原样返回
	Sid    string `json:"callid"` // //标识本次发送id，标识一次下发记录
}

// New 创建新的短信/语音对象
func New(appid, appkey string) *Sms {
	t := new(Sms)
	t.gateway = "https://yun.tim.qq.com/v5"
	t.appid = appid
	t.appkey = appkey
	return t
}

// SendSmsTplCode 发送短信验证码(使用短信模板)
// country 国家代码，sign 短信签名，tplID 短信模板
func (t *Sms) SendSmsTplCode(country, mobile, code string, sign string, tplID int) (string, error) {
	data := new(smsTplCode)
	if country == "" {
		country = "86" // 中国
	}
	if mobile == "" {
		return "", fmt.Errorf("tencent.sms.SendSmsTplCode: mobile is empty")
	}
	if code == "" {
		return "", fmt.Errorf("tencent.sms.SendSmsTplCode: code is empty")
	}
	if sign == "" {
		sign = "药视通"
	}
	data.Tel.Nationcode = country
	data.Tel.Mobile = mobile
	data.Sign = sign
	data.TplID = tplID
	data.Params = []string{code}
	data.Time = time.Now().Unix()
	data.random = string(util.RandStr(4, util.KC_RAND_KIND_NUM))
	data.Sig = t.makeSig(data.Tel.Mobile, data.Time, data.random)
	body, err := util.HTTPPostJSON(fmt.Sprintf("%s/tlssmssvr/sendsms?sdkappid=%s&random=%s", t.gateway, t.appid, data.random), data, 3)
	if err != nil {
		return "", fmt.Errorf("tencent.sms.SendSmsTplCode: request api error %v", err)
	}
	ret := new(resultSms)
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return "", fmt.Errorf("tencent.sms.SendSmsTplCode: json.Marshal response error %v", err)
	}
	if ret.Result != 0 {
		return "", fmt.Errorf("tencent.sms.SendSmsTplCode: response error [%d]%s", ret.Result, ret.ErrMsg)
	}
	return ret.Sid, nil
}

// SendVoiceCode 发送语音验证码
func (t *Sms) SendVoiceCode(country, mobile, code string) (string, error) {
	data := new(voiceCode)
	if country == "" {
		country = "86" // 中国
	}
	if mobile == "" {
		return "", fmt.Errorf("tencent.sms.SendVoiceCode: mobile is empty")
	}
	if code == "" {
		return "", fmt.Errorf("tencent.sms.SendVoiceCode: code is empty")
	}
	data.Tel.Nationcode = country
	data.Tel.Mobile = mobile
	data.Msg = code
	data.Playtimes = 2
	data.Time = time.Now().Unix()
	data.random = string(util.RandStr(4, util.KC_RAND_KIND_NUM))
	data.Sig = t.makeSig(data.Tel.Mobile, data.Time, data.random)
	body, err := util.HTTPPostJSON(fmt.Sprintf("%s/tlsvoicesvr/sendvoice?sdkappid=%s&random=%s", t.gateway, t.appid, data.random), data, 3)
	if err != nil {
		return "", fmt.Errorf("tencent.sms.SendVoiceCode: request api error %v", err)
	}
	ret := new(resultVoice)
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return "", fmt.Errorf("tencent.sms.SendVoiceCode: json.Marshal response error %v", err)
	}
	if ret.Result != 0 {
		return "", fmt.Errorf("tencent.sms.SendVoiceCode: response error [%d]%s", ret.Result, ret.ErrMsg)
	}
	return ret.Sid, nil
}

// makeSig 生成签名
func (t *Sms) makeSig(mobile string, time int64, random string) string {
	str := fmt.Sprintf("appkey=%s&random=%s&time=%d&mobile=%s", t.appkey, random, time, mobile)
	hash := sha256.New()
	hash.Write([]byte(str))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
