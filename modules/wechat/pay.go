package wechat

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
)

const (
	// payOrderURL ...
	payOrderURL    = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	getPayOrderURL = "https://api.mch.weixin.qq.com/pay/orderquery"
)

const (
	// tradeType 支付类型
	tradeType = "JSAPI"
	// signType 签名方式
	signType = "MD5"
	// timeExpireNum 过期时间 单位（分钟）
	timeExpireNum = "5"
)

// PayOrderParam 获取同一订单号的参数
type PayOrderParam struct {
	MchID          string // 商户号
	SpbillCreateIP string // 终端IP
	TotalFee       string // 标价金额
	Body           string // 商品描述
	OutTradeNo     string // 商户订单号
	NotifyURL      string // 回调URL
	OpenID         string // 用户OPENID
}

// payOrderResult 获取统一订单号返回结果
type payOrderResult struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"` // 返回状态码 SUCCESS 成功/FAIL 失败
	ReturnMsg  string   `xml:"return_msg"`  // 返回信息
	AppID      string   `xml:"appid"`       // 小程序ID
	MchID      string   `xml:"mch_id"`      //  商户号
	NonceStr   string   `xml:"nonce_str"`   // 加密随机字符串
	Sign       string   `xml:"sign"`        // 签名验证sign
	ResultCode string   `xml:"result_code"` // 业务结果 SUCCESS 成功/FAIL 失败
	PrepayID   string   `xml:"prepay_id"`   // 预支付交易会话标识
	TradeType  string   `xml:"trade_type"`  // 交易类型 SAPI，NATIVE，APP
}

// GetPayOrderParam 查询订单的 参数
type GetPayOrderParam struct {
	MchID         string
	TransactionID string
	OutTradeNo    string
}

/*
 * PayOrder 获取统一的订单号
 *	appID 小程序ID
 *	appSecret 小程序KEY
 */
func PayOrder(appID, appSecret string, param PayOrderParam) (*payOrderResult, *ErrorResult) {
	xmlParam := payOrderXML(appID, appSecret, param)
	//return nil, &ErrorResult{Message: "微信支付 获取统一订单号 失败"}
	// data, err := util.HTTPPostJSON(payOrderURL, xml, APIRequestTimeout)
	resp, err := http.Post(payOrderURL, "application/x-www-form-urlencoded", strings.NewReader(xmlParam))
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("微信支付 获取统一订单号 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "微信支付 获取统一订单号 失败"}
	}

	row := new(payOrderResult)
	err = xml.Unmarshal([]byte(data), &row)
	if err != nil {
		log.Errorf("微信支付 获取统一订单号 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "微信支付 获取统一订单号 失败"}
	}
	if row.ReturnCode != "SUCCESS" {
		log.Errorf("微信支付 获取统一订单号 失败：%s", row.ReturnMsg)
		return nil, &ErrorResult{Message: "微信支付 获取统一订单号 失败"}
	}

	return row, nil
}

// payOrderXML 获取xml结构的数据
func payOrderXML(appID, appSecret string, param PayOrderParam) string {
	nonceStr := getNonceStr()
	m, _ := time.ParseDuration(timeExpireNum + "m")
	timeExpire := time.Now().Add(m).Format("20060102150405")
	xmlParam := make(map[string]string)
	xmlParam["mch_id"] = param.MchID
	xmlParam["nonce_str"] = nonceStr
	xmlParam["body"] = param.Body
	xmlParam["out_trade_no"] = param.OutTradeNo
	xmlParam["total_fee"] = param.TotalFee
	xmlParam["spbill_create_ip"] = param.SpbillCreateIP
	xmlParam["notify_url"] = param.NotifyURL
	xmlParam["trade_type"] = tradeType
	xmlParam["sign_type"] = signType
	xmlParam["openid"] = param.OpenID
	xmlParam["time_expire"] = timeExpire

	return getXML(appID, appSecret, xmlParam)
}

// getOrderResult 查询订单结果
type getOrderResult struct {
	XMLName        xml.Name `xml:"xml"`
	ReturnCode     string   `xml:"return_code"`      // 返回状态 SUCCESS 成功/FAIL 失败
	ReturnMsg      string   `xml:"return_msg"`       // 返回信息 成功或失败的原因
	AppID          string   `xml:"appid"`            // 小程序ID
	MchID          string   `xml:"mch_id"`           // 商户号
	NonceStr       string   `xml:"nonce_str"`        // 随机加密字符串
	Sign           string   `xml:"sign"`             // 加密验证sign
	ResultCode     string   `xml:"result_code"`      // 业务结果 SUCCESS 成功/FAIL 失败
	ErrCodeDes     string   `xml:"err_code_des"`     // 业务结果为失败时 错误代码描述
	OpenID         string   `xml:"openid"`           // 用户OPENID
	IsSubscribe    string   `xml:"is_subscribe"`     // 是否关注公众账号 Y-关注，N-未关注
	TradeType      string   `xml:"trade_type"`       // 交易类型 JSAPI，NATIVE，APP，MICROPAY
	BankType       string   `xml:"bank_type"`        // 付款银行
	TotalFee       string   `xml:"total_fee"`        // 标价金额
	FeeType        string   `xml:"fee_type"`         // 标价币种
	TransactionID  string   `xml:"transaction_id"`   // 微信订单号
	OutTradeNo     string   `xml:"out_trade_no"`     // 商户订单号
	Attach         string   `xml:"attach"`           // 附加数据
	TimeEnd        string   `xml:"time_end"`         // 支付完成时间
	TradeState     string   `xml:"trade_state"`      // 交易状态
	CashFee        string   `xml:"cash_fee"`         // 现金支付金额
	TradeStateDesc string   `xml:"trade_state_desc"` // 交易状态描述
}

/*
 * GetPayOrder 查询订单数据是否已支付
 *	appID 小程序ID
 *	appSecret 小程序KEY
 */
func GetPayOrder(appID, appSecret string, param GetPayOrderParam) (*getOrderResult, *ErrorResult) {
	xmlParam := getOrderXML(appID, appSecret, param)

	resp, err := http.Post(getPayOrderURL, "application/x-www-form-urlencoded", strings.NewReader(xmlParam))
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("微信支付 查询订单数据 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "微信支付 查询订单数据 失败"}
	}

	row := new(getOrderResult)
	err = xml.Unmarshal([]byte(data), &row)
	if err != nil {
		log.Errorf("微信支付 查询订单数据 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "微信支付 查询订单数据 失败"}
	}
	if row.ReturnCode != "SUCCESS" {
		log.Errorf("微信支付 查询订单数据 失败：%s", row.ReturnMsg)
		return nil, &ErrorResult{Message: "微信支付 查询订单数据 失败"}
	}
	if row.ResultCode != "SUCCESS" {
		log.Errorf("微信支付 查询订单数据 失败：%s", row.ErrCodeDes)
		return nil, &ErrorResult{Message: "微信支付 查询订单数据 失败"}
	}

	return row, nil
}

// getOrderXML 获取查询 xml结构的数据
func getOrderXML(appID, appSecret string, param GetPayOrderParam) string {
	nonceStr := getNonceStr()
	xmlParam := make(map[string]string)
	xmlParam["mch_id"] = param.MchID
	// xmlParam["transaction_id"] = param.TransactionID
	xmlParam["out_trade_no"] = param.OutTradeNo
	xmlParam["nonce_str"] = nonceStr
	xmlParam["sign_type"] = signType

	return getXML(appID, appSecret, xmlParam)
}

// getXML 获取XML结构
func getXML(appID, appSecret string, param map[string]string) string {
	param["appid"] = appID
	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
	buffer.WriteString("<xml>\n")
	for k, v := range param {
		buffer.WriteString("<" + k + ">" + v + "</" + k + ">\n")
	}
	sign := getSign(param, appSecret)
	buffer.WriteString("<sign>" + sign + "</sign>\n")
	buffer.WriteString("</xml>")

	return buffer.String()
}

// getSign 获取签名数据
func getSign(param map[string]string, appSecret string) string {
	httpParam := make(map[string]interface{})
	for k, v := range param {
		httpParam[k] = v
	}
	// 排序数据
	paramString := util.HTTPSortQuery(httpParam)
	//paramString = "appid=" + param["appid"] + "&body=" + param["body"] + "&mch_id=" + param["mch_id"] + "&nonce_str=" + param["nonce_str"] + "&notify_url=" + param["notify_url"] + "&openid=" + param["openid"] + "&out_trade_no=" + param["out_trade_no"] + "&sign_type=" + param["sign_type"] + "&spbill_create_ip=" + param["spbill_create_ip"] + "&time_expire=" + param["time_expire"] + "&total_fee=" + param["total_fee"] + "&trade_type=" + param["trade_type"]
	// 数据先进行MD5加密
	md5String := util.MD5(paramString + "&key=" + appSecret)
	return strings.ToUpper(md5String)
	// 数据进行sha256加密
	// hashSha256 := hmac.New(sha256.New, []byte(appSecret))
	// io.WriteString(hashSha256, paramString+"&key="+appSecret)
	// hashString := fmt.Sprintf("%x", hashSha256.Sum(nil))
	// return strings.ToUpper(hashString)
}

// getNonceStr 获取随机字符串
func getNonceStr() string {
	return util.MD5(util.IntToString(int(time.Now().Unix())) + util.IntToString(rand.Intn(10000)) + util.IntToString(rand.Intn(10000)*rand.Intn(10000)))
}
