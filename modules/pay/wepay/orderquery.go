package wepay

import (
	"encoding/xml"
	"fmt"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/go-baa/log"
)

const (
	OrderQueryGateway = "https://api.mch.weixin.qq.com/pay/orderquery"
)

type OrderQueryOption struct {
	TransactionID string `xml:"transaction_id"`
	OutTradeNO    string `xml:"out_trade_no"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
}

type OrderQueryRequest struct {
	XMLName xml.Name `xml:"xml"`
	AppID   string   `xml:"appid"`
	MchID   string   `xml:"mch_id"`
	OrderQueryOption
}

func OrderQuery(config *Config, option *OrderQueryOption) (*OrderMessage, error) {
	params := map[string]string{}
	request := new(OrderQueryRequest)

	params["appid"] = config.AppID
	request.AppID = config.AppID

	params["mch_id"] = config.MchID
	request.MchID = config.MchID

	params["transaction_id"] = option.TransactionID
	request.TransactionID = option.TransactionID

	params["out_trade_no"] = option.OutTradeNO
	request.OutTradeNO = option.OutTradeNO

	params["nonce_str"] = string(util.RandStr(32, util.KC_RAND_KIND_ALL))
	request.NonceStr = params["nonce_str"]

	// 构建签名
	sign := BuildSign(params, config.MD5Key)
	request.Sign = sign

	post, err := xml.MarshalIndent(request, "", "  ")
	if err != nil {
		log.Errorf("生成微信支付查询订单数据失败: %s\n", err)
		return nil, fmt.Errorf("生成微信支付查询订单数据失败")
	}

	requestXML := string(post)
	log.Printf("\n微信支付查询订单请求 XML: \n%s\n", requestXML)

	body, err := Request(OrderQueryGateway, requestXML, 5)
	if err != nil {
		log.Errorf("请求微信支付查询订单接口失败: %s\n", err)
		return nil, fmt.Errorf("请求微信支付查询订单接口失败")
	}

	responseXML := string(body)
	log.Printf("\n微信支付查询订单响应 XML: \n%s\n", responseXML)

	response := new(OrderMessage)
	if err := xml.Unmarshal(body, &response); err != nil {
		log.Errorf("解析微信支付查询订单响应失败: %s\n", err)
		return nil, fmt.Errorf("解析微信支付查询订单响应失败")
	}

	if response.ReturnCode != "SUCCESS" {
		log.Errorf("请求微信支付查询订单通信失败: %s\n", response.ReturnMsg)
		return nil, fmt.Errorf("请求微信支付查询订单通信失败")
	}

	if response.ResultCode != "SUCCESS" {
		log.Errorf(
			"调用微信支付查询订单接口失败: [%s] %s\n",
			response.ErrCode, response.ErrCodeDes,
		)
		return nil, fmt.Errorf("调用微信支付查询订单接口失败")
	}

	return response, nil
}
