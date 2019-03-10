package unionpay

import (
	"errors"
	"fmt"
	"time"
)

const (
	// UnionpayTestURL 测试地址
	UnionpayTestURL = "https://gateway.test.95516.com"
	// UnionpayProURL 正式生产地址
	UnionpayProURL = "https://gateway.95516.com"
)

// APIConfig Config
type APIConfig struct {
	URL         string
	bizType     string
	accessType  string
	channelType string
}

// CustomerInfo 用户数据
type CustomerInfo struct {
	// 证件类型 01：身份证 02：军官证 03：护照 04：回乡证 05：台胞证 06：警官证 07：士兵证 99：其它证件
	CertifTp string `json:"certifTp"`

	// 证件ID
	CertifID string `json:"certifId"`

	// 名称
	CustomerNm string `json:"customerNm"`

	// 短信验证码
	SmsCode string `json:"smsCode"`

	//使用敏感信息加密证书对 ANSI X9.8 格式的 PIN 加密，并做 Base64 编码
	Pin string `json:"pin"`

	// 三位长度的cvn 信用卡反面后三位
	Cvn2 string `json:"cvn2"`

	// YYMM四位长度的过期时间
	Expired string `json:"expired"`

	// 开卡时预留的手机号
	PhoneNo string `json:"phoneNo"`
}

// RequestParams 入参
type RequestParams struct {
	// 订单ID
	OrderID string
	// 订单日期 format=20060102150405 默认当前时间
	TnxTime string

	// 账户（银行卡号）
	AccNo string

	// 用户验证信息
	Customer *CustomerInfo

	// 扩展参数 交易应答或callback时原样返回 非必填
	Extend string

	// 保留域 非必填
	Reserved map[string]string
}

// PayForAnother 代付类API
type PayForAnother struct {
	c APIConfig
}

// NewPayForAnother 初始化一个代付类
func NewPayForAnother(c APIConfig) (o PayForAnother, err error) {
	if certData.CertID == "" || certData.EncryptID == "" {
		err = errors.New("请先配置证书信息")
		return
	}
	if c.URL == "" {
		c.URL = baseURL
	}
	c.bizType = "000201"
	c.channelType = "07"
	c.accessType = "0"
	return PayForAnother{c}, nil
}

// RealNameAuth 实名认证接口
func (n *PayForAnother) RealNameAuth(bindid string, data *RequestParams) (result interface{}, err error) {
	request := sysParams(n.c, data)
	request["bindId"] = bindid
	request["txnType"] = "72"
	request["txnSubType"] = "01"
	return post(n.c.URL+"/gateway/api/backTransReq.do", request)
}

// Pay 支付接口
func (n *PayForAnother) Pay(orderTradeNO string, amount int, frontURL string, backURL string, data *RequestParams) (string, error) {
	request := make(map[string]string)

	request["certId"] = certData.CertID //证书id
	request["merId"] = merID            //商户代码，请改自己的测试商户号

	request["frontUrl"] = frontURL                           //前台通知地址
	request["backUrl"] = backURL                             //后台通知地址
	request["orderId"] = orderTradeNO                        //商户订单号
	request["txnTime"] = time.Now().Format("20060102150405") //订单发送时间
	request["txnAmt"] = fmt.Sprintf("%d", amount)            //交易金额，单位分
	request["signMethod"] = "01"                             //签名方法

	request["version"] = "5.0.0"       //版本号
	request["encoding"] = "utf-8"      //编码方式
	request["txnType"] = "01"          //交易类型
	request["txnSubType"] = "01"       //交易子类
	request["bizType"] = "000201"      //业务类型
	request["channelType"] = "08"      //渠道类型，07-PC，08-手机
	request["accessType"] = "0"        //接入类型
	request["currencyCode"] = "156"    //交易币种
	request["defaultPayType"] = "0001" //默认支付方式

	kvs, err := GenKVpairs(frontConsumeParamMap, request, "signature")
	if err != nil {
		return "", err
	}

	request["signature"], _ = signature(certData.Private, kvs)
	return post(n.c.URL+"/gateway/api/frontTransReq.do", request)
}

var frontConsumeParamMap = map[string]bool{
	"version":         true,  // 版本号 固定填写5.0.0
	"encoding":        true,  // 编码方式 默认值 UTF-8
	"certId":          true,  // 证书id
	"signature":       true,  // 签名 填写对报文摘要的签名
	"signMethod":      true,  // 签名方式 取值：01 表示采用的是RSA
	"txnType":         true,  // 交易类型 取值：01
	"txnSubType":      true,  // 交易子类 01:自主消费，通过地址的方式区分前台消费和后台消费（含无跳转支付） 03:分期付款
	"bizType":         true,  // 产品类型 000201
	"channelType":     true,  // 渠道类型
	"frontUrl":        false, // 前台通知地址 前台返回商户结果时使用，前台类交易需上送
	"backUrl":         true,  // 后台通知地址 后台返回商户结果时使用，如上送，则发送商户后台交易结果通知
	"accessType":      true,  // 接入类型 0:普通商户直接接入 2:平台类商户接入
	"merId":           true,  // 商户代码
	"subMerId":        false, // 二级商户代码 商户类型为平台商户接入时必须上送
	"subMerName":      false, // 二级商户全称 商户类型为平台商户接入时必须上送
	"subMerAbbr":      false, // 二级商户简称 商户类型为平台商户接入时必须上送
	"orderId":         true,  // 商户订单号 商户端生成
	"txnTime":         true,  // 订单发送时间 商户发送交易时间
	"accType":         false, // 账号类型 后台类交易且卡号上送; 跨行收单且收单机构收集银行卡 信息时上送 01: 02: 03:IC  默认取值: 取值“03”表示以 IC 终端发起的 IC 卡交易,IC 作为普通银行卡进行支 付时,此域填写为“01”
	"accNo":           false, // 账号 1、 后台类消费交易时上送全卡号 2、 跨行收单且收单机构收集银行 卡信息时上送、 3、前台类交易可通过配置后返回, 卡号可选上送
	"txnAmt":          true,  // 交易金额 单位为分
	"currencyCode":    true,  // 交易币种 默认为156
	"customerInfo":    false, // 银行卡验证信息及身法信息 1、后台类消费交易时上送 2、认证支付 2.0,后台交易时可选 Key=value 格式
	"orderTimeout":    false, // 账号接受超时时间（防钓鱼使用）1、前台类消费交易时上送 2、认证支付 2.0,后台交易时可选
	"payTimeout":      false, // 订单支付超时时间 超过此时间用户支付成功的交易, 不通知商户,系统自动退款,大约 5 个工作日金额返还到用户账户
	"termId":          false, // 终端号
	"reqReserved":     false, // 请求方保留域 商户自定义保留域，交易应答时会原样返回
	"reserved":        false, // 保留域
	"riskRateInfo":    false, // 风险信息域
	"encryptCertId":   false, // 加密证书
	"frontFailUrl":    false, // 失败交易前台跳转地址 前台消费交易弱商户上送此字段，则在支付失败时，页面跳转至商户该URL（不带交易信息，仅跳转）
	"instalTransInfo": false, // 分期付款信息域 分期付款交易，商户端选择分期信息时，需上送组合域，填法见数据元说明
	"defaultPayType":  false, // 默认支付方式 取值参考数据字典
	"issInsCode":      false, // 发卡机构代码 1、当账号类型为 02-存折时需填写 2、在前台类交易时填写默认银行 代码,支持直接跳转到网银。银行简码列表参考附录：C.1,C.2，其中C.2银行列表仅支持借记卡
	"supPayType":      false, // 支持支付方式 仅仅 pc 使用,使用哪种支付方式 由收单机构填写,取值为以下内容 的一种或多种,通过逗号(,)分 割。取值参考数据字典
	"userMac":         false, // 终端信息域 移动支付业务需要上送
	"customerIp":      false, // 持卡人IP 前台交易，有IP防钓鱼要求的商户上送
	"cardTransData":   false, // 有卡交易信息域 有卡交易必填
	"orderDesc":       false, // 订单描述 移动支付上送
}
