package unionpay

import (
	"errors"

	"github.com/go-baa/log"
)

// 返回参数校验
var notifyParamMap = map[string]bool{
	"queryId":            true, // 消费交易流水号，供后续查询使用
	"currencyCode":       true, // 交易币种 默认为156
	"traceTime":          true, // 交易传输时间
	"signature":          true, // 签名 填写对报文摘要的签名
	"signMethod":         true, // 签名方式 取值：01 表示采用的是RSA
	"settleCurrencyCode": true, // 清算币种
	"settleAmt":          true, // 清算金额
	"settleDate":         true, // 清算日期
	"traceNo":            true, // 系统跟踪号
	"respCode":           true, // 应答码
	"respMsg":            true, // 应答信息

	"exchangeDate":     false, // 兑换日期：交易成功，交易币种和清算币种不一致的时候返回
	"signPubKeyCert":   false, // 签名公钥证书：使用RSA签名方式时必选，此域填写银联签名公钥证书。
	"exchangeRate":     false, // 清算汇率：交易成功，交易币种和清算币种不一致的时候返回
	"accNo":            false, // 账号：1、 后台类消费交易时上送全卡号 2、 跨行收单且收单机构收集银行 卡信息时上送、 3、前台类交易可通过配置后返回, 卡号可选上送
	"payType":          false, // 支付方式：根据商户配置返回
	"payCardNo":        false, // 支付卡标识：移动支付交易时，根据商户配置返回
	"payCardType":      false, // 支付卡类型：根据商户配置返回
	"payCardIssueName": false, // 支付卡名称：	移动支付交易时，根据商户配置返回
	"version":          false, // 版本号
	"bindId":           false, // 绑定标识号：绑定支付时，根据商户配置返回
	"encoding":         false, // 编码方式
	"bizType":          false, // 产品类型：000201
	"txnTime":          false, // 订单发送时间 商户发送交易时间
	"txnAmt":           false, // 交易金额 交易金额 单位为分
	"txnType":          false, // 交易类型
	"txnSubType":       false, // 交易子类：01:自主消费，通过地址的方式区分前台消费和后台消费（含无跳转支付） 03:分期付款
	"accessType":       false, // 接入类型：0:普通商户直接接入 2:平台类商户接入
	"reqReserved":      false, // 接入类型
	"merId":            false, // 商户代码
	"orderId":          false, // 用户订单
	"reserved":         false, // 保留域
	"accSplitData":     false, // 分账域
}

// ValidateSignature 校验签名
func ValidateSignature(params map[string]string) error {
	kvs, err := GenKVpairs(notifyParamMap, params, "signature")
	if err != nil {
		log.Errorf("签名验证转化失败:%v", err)
		return errors.New("签名验证转化kvs错误")
	}
	sign, err := signature(certData.Private, kvs)
	if err != nil {
		return err
	}
	if params["signature"] == sign {
		return errors.New("签名验证失败")
	}
	return nil
}
