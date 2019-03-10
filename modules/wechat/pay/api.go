package pay

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

const (
	// TradeTypeJS 公众号支付
	TradeTypeJS = "JSAPI"
	// TradeTypeNative 原生扫码支付
	TradeTypeNative = "NATIVE"
	// TradeTypeApp app支付
	TradeTypeApp = "APP"
)

const (
	// SignTypeMD5 md5签名算法
	SignTypeMD5 = "MD5"
)

const (
	// FeeTypeCNY 人民币
	FeeTypeCNY = "CNY"
)

const (
	// ReturnCodeSuccess 请求成功
	ReturnCodeSuccess = "SUCCESS"
	// ReturnCodeFail 请求失败
	ReturnCodeFail = "FAIL"
)

const (
	// TradeStatSuccess 支付成功
	TradeStatSuccess = "SUCCESS"
	// TradeStatRefund 转入退款
	TradeStatRefund = "REFUND"
	// TradeStatNotPay 未支付
	TradeStatNotPay = "NOTPAY"
	// TradeStatClosed 已关闭
	TradeStatClosed = "CLOSED"
	// TradeStatRevoked 已撤销（刷卡支付）
	TradeStatRevoked = "REVOKED"
	// TradeStatUserPaying 用户支付中
	TradeStatUserPaying = "USERPAYING"
	// TradeStatPayError 支付失败(其他原因，如银行返回失败)
	TradeStatPayError = "PAYERROR"
)

// CommonResponse 通用响应结构
type CommonResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
	ResultCode string   `xml:"result_code"`
	ErrCode    string   `xml:"err_code"`
	ErrCodeDes string   `xml:"err_code_des"`
}

// UnifiedOrderRequest 统一下单请求参数
type UnifiedOrderRequest struct {
	TradeNo    string    // 订单号，必填
	Amount     int       // 单位：分， 必填
	Desc       string    // 商品描述， 必填
	IP         string    // 用户端IP， 必填
	NotifyURL  string    // 微信支付异步通知回调地址，必填
	TradeType  string    // 交易类型, 必填
	TimeStart  time.Time // 订单生成时间，格式为yyyyMMddHHmmss
	TimeExpire time.Time // 订单失效时间，格式为yyyyMMddHHmmss
	Attach     string    // 附加参数
}

// UnifiedOrderResponse 统一下单请求响应
type UnifiedOrderResponse struct {
	XMLName     xml.Name `xml:"xml"`
	ReturnCode  string   `xml:"return_code"`
	ReturnMsg   string   `xml:"return_msg"`
	ResultCode  string   `xml:"result_code"`
	ErrCode     string   `xml:"err_code"`
	ErrCodeDesc string   `xml:"err_code_des"`
	AppID       string   `xml:"appid"`
	MchID       string   `xml:"mch_id"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	TradeType   string   `xml:"trade_type"`
	PrepayID    string   `xml:"prepay_id"`
}

// UnifiedOrder 统一下单
func (t *WXPay) UnifiedOrder(req *UnifiedOrderRequest) (string, error) {
	err := req.validate()
	if err != nil {
		return "", err
	}

	params := map[string]string{
		"out_trade_no":     req.TradeNo,
		"total_fee":        fmt.Sprintf("%d", req.Amount),
		"body":             req.Desc,
		"spbill_create_ip": req.IP,
		"notify_url":       req.NotifyURL,
		"trade_type":       req.TradeType,
	}
	if req.TradeType == "" {
		params["trade_type"] = TradeTypeApp
	}
	if !req.TimeStart.IsZero() {
		params["time_start"] = req.TimeStart.Format(TimeLayout)
	}
	if !req.TimeExpire.IsZero() {
		params["time_expire"] = req.TimeExpire.Format(TimeLayout)
	}

	res, err := t.api("unifiedorder", params)
	if err != nil {
		return "", err
	}

	response := new(UnifiedOrderResponse)
	err = xml.Unmarshal(res, response)
	if err != nil {
		return "", fmt.Errorf("xml解码错误:%v", err)
	}

	if response.ReturnCode == ReturnCodeFail {
		return "", fmt.Errorf("returnCodeFail, err:%s", response.ReturnMsg)
	}

	if err = t.CheckSign(response, response.Sign); err != nil {
		return "", err
	}

	if response.ResultCode == ReturnCodeFail {
		return "", fmt.Errorf("resultCodeFail, errcode:%s, errmsg:%s", response.ErrCode, response.ErrCodeDesc)
	}

	return response.PrepayID, nil
}

func (t *UnifiedOrderRequest) validate() error {
	if t.TradeNo == "" {
		return fmt.Errorf("订单号不能为空")
	}
	if t.Amount <= 0 {
		return fmt.Errorf("金额要求大于0")
	}
	if t.Desc == "" {
		return fmt.Errorf("商品描述不能为空")
	}
	if t.IP == "" {
		return fmt.Errorf("客户端IP不能为空")
	}
	if t.NotifyURL == "" {
		return fmt.Errorf("回调通知地址不能为空")
	}

	return nil
}

// QueryOrderResponse 订单查询/支付回调通知响应结构
type QueryOrderResponse struct {
	XMLName        xml.Name `xml:"xml" json:"-"`
	ReturnCode     string   `xml:"return_code" json:"return_code"`
	ReturnMsg      string   `xml:"return_msg" json:"return_msg"`
	ResultCode     string   `xml:"result_code" json:"result_code"`
	ErrCode        string   `xml:"err_code" json:"err_code"`
	ErrCodeDesc    string   `xml:"err_code_des" json:"err_code_desc"`
	AppID          string   `xml:"appid" json:"app_id"`
	MchID          string   `xml:"mch_id" json:"mch_id"`
	NonceStr       string   `xml:"nonce_str" json:"nonce_str"`
	Sign           string   `xml:"sign" json:"sign"`
	DeviceInfo     string   `xml:"device_info" json:"device_info"`
	OpenID         string   `xml:"openid" json:"open_id"`
	TradeType      string   `xml:"trade_type" json:"trade_type"`
	TradeState     string   `xml:"trade_state" json:"trade_state"`
	TradeStateDesc string   `xml:"trade_state_desc" json:"trade_state_desc"`
	BankType       string   `xml:"bank_type" json:"bank_type"`
	TotalFee       string   `xml:"total_fee" json:"total_fee"`
	FeeType        string   `xml:"fee_type" json:"fee_type"`
	CashFee        string   `xml:"cash_fee" json:"cash_fee"`
	CashFeeType    string   `xml:"cash_fee_type" json:"cash_fee_type"`
	CouponFee      string   `xml:"coupon_fee" json:"coupon_fee"`
	CouponCount    string   `xml:"coupon_count" json:"coupon_count"`
	TransactionID  string   `xml:"transaction_id" json:"transaction_id"`
	TradeNo        string   `xml:"out_trade_no" json:"trade_no"`
	Attach         string   `xml:"attach" json:"attach"`
	TimeEnd        string   `xml:"time_end" json:"time_end"`
	IsSubscribe    string   `xml:"is_subscribe" json:"is_subscribe"`
}

// OrderQuery 订单查询
func (t *WXPay) OrderQuery(tradeNo string) (*QueryOrderResponse, error) {
	params := map[string]string{
		"out_trade_no": tradeNo,
	}

	res, err := t.api("orderquery", params)
	if err != nil {
		return nil, err
	}

	response := new(QueryOrderResponse)
	err = xml.Unmarshal(res, response)
	if err != nil {
		return nil, fmt.Errorf("xml解码错误:%v", err)
	}

	if response.ReturnCode == ReturnCodeFail {
		return nil, fmt.Errorf("returnCodeFail, err:%s", response.ReturnMsg)
	}

	if err := t.CheckSign(response, response.Sign); err != nil {
		return nil, err
	}

	if response.ResultCode == ReturnCodeFail {
		return nil, fmt.Errorf("resultCodeFail, errcode:%s, errmsg:%s", response.ErrCode, response.ErrCodeDesc)
	}

	return response, nil
}

// NotifyResponse 回调通知响应
type NotifyResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

// CloseOrder 关闭订单
func (t *WXPay) CloseOrder() {

}

// Refund 申请退款
func (t *WXPay) Refund() {

}

// RefundQuery 退款查询
func (t *WXPay) RefundQuery() {

}

// DownloadBill 下载对账单
func (t *WXPay) DownloadBill() {

}

// TransferRequest 付款请求结构
type TransferRequest struct {
	TradeNo string
	Openid  string
	Amount  int
	Desc    string
	IP      string
}

const (
	// TransferCheckNameNO 不校验姓名
	TransferCheckNameNO = "NO_CHECK"
	// TransferCheckNameForce 强校验姓名
	TransferCheckNameForce = "FORCE_CHECK"
	// MinimumTransferAmount 最小转账金额
	MinimumTransferAmount = 100
	// TransferAPIGateway 微信付款接口地址
	TransferAPIGateway = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
)

// TransferResponse 付款响应结构
type TransferResponse struct {
	CommonResponse
	ResultCode     string `xml:"result_code" json:"result_code"`
	ErrCode        string `xml:"err_code" json:"err_code"`
	ErrCodeDesc    string `xml:"err_code_des" json:"err_code_desc"`
	PartnerTradeNo string `xml:"partner_trade_no" json:"partner_trade_no"`
	PaymentNo      string `xml:"payment_no" json:"payment_no"`
	PaymentTime    string `xml:"payment_time" json:"payment_time"`
}

// Transfer 付款到个人微信, 返回微信流水号
func (t *WXPay) Transfer(req *TransferRequest) (string, error) {
	if req.Amount < MinimumTransferAmount {
		return "", fmt.Errorf("付款金额最小1元")
	}

	params := map[string]string{
		"mch_appid":        t.appID,
		"mchid":            t.mchID,
		"partner_trade_no": req.TradeNo,
		"openid":           req.Openid,
		"check_name":       TransferCheckNameNO,
		"amount":           fmt.Sprintf("%d", req.Amount),
		"desc":             req.Desc,
		"spbill_create_ip": req.IP,
		"nonce_str":        string(util.RandStr(32, util.KC_RAND_KIND_ALL)),
	}
	params["sign"] = t.Sign(params)
	reqBody := MapToXMLString(params)

	if setting.Debug {
		log.Printf("WXPay transfer reqbody:%s\n", reqBody)
	}

	res, err := t.request(TransferAPIGateway, "", []byte(reqBody), true)
	if err != nil {
		return "", fmt.Errorf("请求错误:%v", err)
	}

	if setting.Debug {
		log.Printf("WXPay transfer resbody:%s\n", string(res))
	}

	response := new(TransferResponse)
	err = xml.Unmarshal(res, response)
	if err != nil {
		return "", fmt.Errorf("xml解码错误:%v", err)
	}

	if response.ReturnCode == ReturnCodeFail {
		return "", fmt.Errorf("returnCodeFail, err:%s", response.ReturnMsg)
	}

	if response.ResultCode == ReturnCodeFail {
		return "", fmt.Errorf("resultCodeFail, errcode:%s, errmsg:%s", response.ErrCode, response.ErrCodeDesc)
	}

	return response.PaymentNo, nil
}

// PaymentRequest APP支付数结构
type PaymentRequest struct {
	AppID     string `xml:"appid" json:"app_id"`
	PartnerID string `xml:"partnerid" json:"partner_id"`
	PrepayID  string `xml:"prepayid" json:"prepay_id"`
	Package   string `xml:"package" json:"package_value"`
	NonceStr  string `xml:"noncestr" json:"nonce_str"`
	Timestamp string `xml:"timestamp" json:"timestamp"`
	Sign      string `xml:"sign" json:"sign"`
}

// GetAppPaymentRequest 组装APP支付结构
func (t *WXPay) GetAppPaymentRequest(prepayID string) (*PaymentRequest, error) {
	req := new(PaymentRequest)
	req.AppID = t.appID
	req.PartnerID = t.mchID
	req.PrepayID = prepayID
	req.Package = "Sign=WXPay"
	req.NonceStr = string(util.RandStr(32, util.KC_RAND_KIND_ALL))
	req.Timestamp = fmt.Sprintf("%d", time.Now().Unix())

	reqMap, err := StructToMap(req)
	if err != nil {
		return nil, err
	}
	req.Sign = t.Sign(reqMap)

	return req, nil
}

// SandboxSignKeyResponse 沙箱秘钥响应
type SandboxSignKeyResponse struct {
	CommonResponse
	MchID          string `xml:"mch_id"`
	SandboxSignkey string `xml:"sandbox_signkey"`
}

// GetSanboxSignKey 获取沙箱秘钥
func (t *WXPay) GetSanboxSignKey() (string, error) {
	url := GatewaySandbox + "getsignkey"
	params := map[string]string{
		"mch_id":    t.mchID,
		"nonce_str": string(util.RandStr(32, util.KC_RAND_KIND_ALL)),
	}

	reqBody := t.buildXMLParams(params)
	res, err := t.request(url, "", []byte(reqBody), false)
	if err != nil {
		return "", fmt.Errorf("请求错误:%v", err)
	}

	response := new(SandboxSignKeyResponse)
	err = xml.Unmarshal(res, response)
	if err != nil {
		return "", fmt.Errorf("xml解码错误:%v", err)
	}

	if response.ReturnCode == ReturnCodeFail {
		return "", fmt.Errorf("returnCodeFail, err:%s", response.ReturnMsg)
	}

	return response.SandboxSignkey, nil
}
