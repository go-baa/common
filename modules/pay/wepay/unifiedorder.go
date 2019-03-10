package wepay

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
)

const (
	UnifiedOrderGateway = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

type UnifiedOrderOptionGoods struct {
	GoodsID       string `json:"goods_id"`
	WXPayGoodsID  string `json:"wxpay_goods_id"`
	GoodsName     string `json:"goods_name"`
	GoodsNum      int    `json:"goods_num"`
	Price         int    `json:"price"`
	GoodsCategory string `json:"goods_category"`
	Body          string `json:"body"`
}

type UnifiedOrderOptionDetail struct {
	GoodsDetail []*UnifiedOrderOptionGoods `json:"goods_detail"`
}

type UnifiedOrderOption struct {
	DeviceInfo     string                    `xml:"device_info,omitempty"` // 设备号
	NonceStr       string                    `xml:"nonce_str"`             // 必填, 随机字符串，不长于32位
	Sign           string                    `xml:"sign"`                  // 必填, 签名
	Body           string                    `xml:"body"`                  // 必填, 商品描述
	Detail         *UnifiedOrderOptionDetail `xml:"-"`                     // 商品详情
	DetailJSON     string                    `xml:"detail,omitempty"`      // 商品详情
	Attach         string                    `xml:"attach,omitempty"`      // 附加数据
	OutTradeNO     string                    `xml:"out_trade_no"`          // 必填, 商户订单号
	FeeType        string                    `xml:"fee_type,omitempty"`    // 货币类型
	TotalFee       int                       `xml:"total_fee"`             // 必填, 订单总金额，单位为分
	SPBillCreateIP string                    `xml:"spbill_create_ip"`      // 必填, 终端IP
	TimeStart      *time.Time                `xml:"-"`                     // 订单生成时间，格式为yyyyMMddHHmmss
	TimeStartStr   string                    `xml:"time_start,omitempty"`  //
	TimeExpire     *time.Time                `xml:"-"`                     // 订单失效时间，格式为yyyyMMddHHmmss
	TimeExpireStr  string                    `xml:"time_expire,omitempty"` //
	GoodsTag       string                    `xml:"goods_tag,omitempty"`   // 商品标记
	NotifyURL      string                    `xml:"notify_url"`            // 必填, 接收微信支付异步通知回调地址
	TradeType      string                    `xml:"trade_type"`            // 必填, 取值如下：JSAPI，NATIVE，APP
	ProductID      string                    `xml:"product_id,omitempty"`  // trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义
	LimitPay       string                    `xml:"limit_pay,omitempty"`   // 指定支付方式, no_credit--指定不能使用信用卡支付
	OpenID         string                    `xml:"openid,omitempty"`      // 用户标识
}

type UnifiedOrderRequest struct {
	XMLName xml.Name `xml:"xml"`
	AppID   string   `xml:"appid"`
	MchID   string   `xml:"mch_id"`
	UnifiedOrderOption
}

type UnifiedOrderResponse struct {
	CommonMessage
	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	TradeType  string `xml:"trade_type"`
	PrePayID   string `xml:"prepay_id"`
	CodeURL    string `xml:"code_url"`
}

func UnifiedOrder(config *Config, option *UnifiedOrderOption) (*UnifiedOrderResponse, error) {
	params := map[string]string{}
	request := new(UnifiedOrderRequest)

	params["appid"] = config.AppID
	request.AppID = config.AppID

	params["mch_id"] = config.MchID
	request.MchID = config.MchID

	params["device_info"] = option.DeviceInfo
	request.DeviceInfo = option.DeviceInfo

	params["nonce_str"] = string(util.RandStr(32, util.KC_RAND_KIND_ALL))
	request.NonceStr = params["nonce_str"]

	params["body"] = option.Body
	request.Body = option.Body

	if option.Detail != nil {
		detail, err := json.Marshal(option.Detail)
		if err == nil {
			params["detail"] = string(detail)
			request.DetailJSON = params["detail"]
		}
	}

	params["attach"] = option.Attach
	request.Attach = option.Attach

	params["out_trade_no"] = option.OutTradeNO
	request.OutTradeNO = option.OutTradeNO

	params["fee_type"] = option.FeeType
	request.FeeType = option.FeeType

	params["total_fee"] = strconv.Itoa(option.TotalFee)
	request.TotalFee = option.TotalFee

	params["spbill_create_ip"] = option.SPBillCreateIP
	request.SPBillCreateIP = option.SPBillCreateIP

	if option.TimeStart != nil {
		params["time_start"] = option.TimeStart.Format("20060102150405")
		request.TimeStartStr = params["time_start"]
	}
	if option.TimeExpire != nil {
		params["time_expire"] = option.TimeExpire.Format("20060102150405")
		request.TimeExpireStr = params["time_expire"]
	}

	params["goods_tag"] = option.GoodsTag
	request.GoodsTag = option.GoodsTag

	params["notify_url"] = option.NotifyURL
	request.NotifyURL = option.NotifyURL

	if option.TradeType == "" {
		params["trade_type"] = "NATIVE"
	} else {
		params["trade_type"] = option.TradeType
	}
	request.TradeType = params["trade_type"]

	params["product_id"] = option.ProductID
	request.ProductID = option.ProductID

	params["limit_pay"] = option.LimitPay
	request.LimitPay = option.LimitPay

	params["openid"] = option.OpenID
	request.OpenID = option.OpenID

	// 构建签名
	sign := BuildSign(params, config.MD5Key)
	request.Sign = sign

	post, err := xml.MarshalIndent(request, "", "  ")
	if err != nil {
		log.Errorf("生成微信支付统一下单数据失败: %s\n", err)
		return nil, fmt.Errorf("生成微信支付统一下单数据失败")
	}

	requestXML := string(post)
	log.Printf("\n微信支付统一下单请求 XML: \n%s\n", requestXML)

	body, err := Request(UnifiedOrderGateway, requestXML, 5)
	if err != nil {
		log.Errorf("请求微信支付统一下单接口失败: %s\n", err)
		return nil, fmt.Errorf("请求微信支付统一下单接口失败")
	}

	responseXML := string(body)
	log.Printf("\n微信支付统一下单响应 XML: \n%s\n", responseXML)

	response := new(UnifiedOrderResponse)
	if err := xml.Unmarshal(body, &response); err != nil {
		log.Errorf("解析微信支付统一下单响应失败: %s\n", err)
		return nil, fmt.Errorf("解析微信支付统一下单响应失败")
	}

	if response.ReturnCode != "SUCCESS" {
		log.Errorf("请求微信支付统一下单通信失败: %s\n", response.ReturnMsg)
		return nil, fmt.Errorf("请求微信支付统一下单通信失败")
	}

	if response.ResultCode != "SUCCESS" {
		log.Errorf(
			"调用微信支付统一下单接口失败: [%s] %s\n",
			response.ErrCode, response.ErrCodeDes,
		)
		return nil, fmt.Errorf("调用微信支付统一下单接口失败")
	}

	return response, nil
}

func UnifiedOrderParseNotify(body []byte) (*OrderMessage, error) {
	notify := new(OrderMessage)

	if err := xml.Unmarshal(body, &notify); err != nil {
		log.Errorf("解析微信支付支付结果通用通知失败: %s\n", err)
		return nil, fmt.Errorf("解析微信支付支付结果通用通知失败")
	}

	if notify.ReturnCode != "SUCCESS" {
		log.Errorf("微信支付支付结果通用通知通信失败: %s\n", notify.ReturnMsg)
		return nil, fmt.Errorf("微信支付支付结果通用通知通信失败")
	}

	if notify.ResultCode != "SUCCESS" {
		log.Errorf(
			"微信支付支付结果通用通知业务失败: [%s] %s\n",
			notify.ErrCode, notify.ErrCodeDes,
		)
		return nil, fmt.Errorf("微信支付支付结果通用通知业务失败")
	}

	return notify, nil
}
