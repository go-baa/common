package pay

import (
	"fmt"
)

const (
	// TradeAppPayMethod APP交易接口
	TradeAppPayMethod = "alipay.trade.app.pay"
	// TradeQueryMethod 订单查询接口
	TradeQueryMethod = "alipay.trade.query"
)

// DefaultTradeExpiredHours 订单默认有效期
const DefaultTradeExpiredHours = 24

const (
	// TradeStatusWaitBuyerPay 交易创建，等待买家付款
	TradeStatusWaitBuyerPay = "WAIT_BUYER_PAY"
	// TradeStatusClosed 未付款交易超时关闭，或支付完成后全额退款
	TradeStatusClosed = "TRADE_CLOSED"
	// TradeStatusSuccess 交易支付成功
	TradeStatusSuccess = "TRADE_SUCCESS"
	// TradeStatusFinished // 交易结束，不可退款
	TradeStatusFinished = "TRADE_FINISHED"
)

// TradeAppPay App创建订单参数
type TradeAppPay struct {
	OutTradeNo     string `json:"out_trade_no"`
	SellerID       string `json:"seller_id"`
	TotalAmount    string `json:"total_amount"`
	Subject        string `json:"subject"`
	TimeoutExpress string `json:"timeout_express,omitempty"`
	NotifyURL      string `json:"-"`
}

// APIName api名称
func (t TradeAppPay) APIName() string {
	return TradeAppPayMethod
}

// Params 扩展公共参数
func (t TradeAppPay) Params() map[string]string {
	var m = make(map[string]string)
	m["notify_url"] = t.NotifyURL
	return m
}

// ExtJSONParamName 扩展json字段键名
func (t TradeAppPay) ExtJSONParamName() string {
	return "biz_content"
}

// ExtJSONParamValue 扩展json字段值
func (t TradeAppPay) ExtJSONParamValue() string {
	return jsonMarshal(t)
}

// GetTradeAppPayParams 获取app支付信息参数
type GetTradeAppPayParams struct {
	TradeNo      string
	Amount       int // 单位：分
	Title        string
	ExpiredHours int
	NotifyURL    string
}

// GetTradeAppPay 组装客户端发起交易参数
func (t *AliPay) GetTradeAppPay(params *GetTradeAppPayParams) string {
	trade := new(TradeAppPay)
	trade.OutTradeNo = params.TradeNo
	trade.Subject = params.Title
	trade.SellerID = t.sellerID
	trade.TotalAmount = fmt.Sprintf("%0.2f", float64(params.Amount)/100)
	if params.ExpiredHours == 0 {
		params.ExpiredHours = DefaultTradeExpiredHours
	}
	trade.TimeoutExpress = fmt.Sprintf("%dh", params.ExpiredHours)
	trade.NotifyURL = params.NotifyURL

	return t.buildParams(trade)
}

// TradeQueryRequest 交易查询请求参数
type TradeQueryRequest struct {
	OutTradeNo string `json:"out_trade_no,omitempty"`
	TradeNo    string `json:"trade_no,omitempty"`
}

// APIName api名称
func (t TradeQueryRequest) APIName() string {
	return TradeQueryMethod
}

// Params 扩展公共参数
func (t TradeQueryRequest) Params() map[string]string {
	var m = make(map[string]string)
	return m
}

// ExtJSONParamName 扩展json字段键名
func (t TradeQueryRequest) ExtJSONParamName() string {
	return "biz_content"
}

// ExtJSONParamValue 扩展json字段值
func (t TradeQueryRequest) ExtJSONParamValue() string {
	return jsonMarshal(t)
}

// TradeQueryResponse 交易查询响应
type TradeQueryResponse struct {
	TradeQueryResponse TradeQuery `json:"alipay_trade_query_response"`
	Sign               string     `json:"sign"`
}

// TradeQuery 交易查询响应信息主体
type TradeQuery struct {
	Code           string      `json:"code"`
	Msg            string      `json:"msg"`
	SubCode        string      `json:"sub_code"`
	SubMsg         string      `json:"sub_msg"`
	TradeNo        string      `json:"trade_no"`                 // 支付宝交易号
	OutTradeNo     string      `json:"out_trade_no"`             // 商家订单号
	BuyerLogonID   string      `json:"buyer_logon_id"`           // 买家支付宝账号
	TradeStatus    string      `json:"trade_status"`             // 交易状态
	TotalAmount    float64     `json:"total_amount,string"`      // 交易的订单金额
	ReceiptAmount  float64     `json:"receipt_amount,string"`    // 实收金额，单位为元，两位小数
	BuyerPayAmount float64     `json:"buyer_pay_amount,string"`  // 买家实付金额，单位为元，两位小数。
	PointAmount    float64     `json:"point_amount,string"`      // 积分支付的金额，单位为元，两位小数。
	InvoiceAmount  float64     `json:"invoice_amount,string"`    // 交易中用户支付的可开具发票的金额，单位为元，两位小数。
	SendPayDate    string      `json:"send_pay_date"`            // 本次交易打款给卖家的时间
	StoreID        string      `json:"store_id"`                 // 商户门店编号
	TerminalID     string      `json:"terminal_id"`              // 商户机具终端编号
	FundBillList   []*FundBill `json:"fund_bill_list,omitempty"` // 交易支付使用的资金渠道
	StoreName      string      `json:"store_name"`               // 请求交易支付中的商户店铺的名称
	BuyerUserID    string      `json:"buyer_user_id"`            // 买家在支付宝的用户id
}

// FundBill 资金渠道
type FundBill struct {
	FundChannel string  `json:"fund_channel"`       // 交易使用的资金渠道，详见 支付渠道列表
	Amount      string  `json:"amount"`             // 该支付工具类型所使用的金额
	RealAmount  float64 `json:"real_amount,string"` // 渠道实际付款金额
}

// QueryTrade 交易查询
func (t *AliPay) QueryTrade(tradeNo string) (result *TradeQueryResponse, err error) {
	param := &TradeQueryRequest{
		OutTradeNo: tradeNo,
	}

	err = t.doRequest(param, &result)
	return result, err
}

// TradeNotification 交易通知
type TradeNotification struct {
	NotifyTime        string `json:"notify_time"`                   // 通知时间
	NotifyID          string `json:"notify_id"`                     // 通知校验ID
	NotifyType        string `json:"notify_type"`                   // 通知类型
	AppID             string `json:"app_id"`                        // 开发者的app_id
	Charset           string `json:"charset"`                       // 编码格式
	Version           string `json:"version"`                       // 接口版本
	SignType          string `json:"sign_type"`                     // 签名类型
	Sign              string `json:"sign"`                          // 签名
	TradeNo           string `json:"trade_no"`                      // 支付宝交易号
	OutTradeNo        string `json:"out_trade_no"`                  // 商户订单号
	OutBizNo          string `json:"out_biz_no,omitempty"`          // 商户业务号
	BuyerID           string `json:"buyer_id,omitempty"`            // 买家支付宝用户号
	BuyerLogonID      string `json:"buyer_logon_id,omitempty"`      // 买家支付宝账号
	SellerID          string `json:"seller_id,omitempty"`           // 卖家支付宝用户号
	SellerEmail       string `json:"seller_email,omitempty"`        // 卖家支付宝账号
	TradeStatus       string `json:"trade_status"`                  // 交易状态
	TotalAmount       string `json:"total_amount"`                  // 订单金额
	ReceiptAmount     string `json:"receipt_amount"`                // 实收金额
	InvoiceAmount     string `json:"invoice_amount"`                // 开票金额
	BuyerPayAmount    string `json:"buyer_pay_amount"`              // 付款金额
	PointAmount       string `json:"point_amount"`                  // 集分宝金额
	RefundFee         string `json:"refund_fee"`                    // 总退款金额
	Subject           string `json:"subject"`                       // 总退款金额
	Body              string `json:"body"`                          // 商品描述
	GmtCreate         string `json:"gmt_create"`                    // 交易创建时间
	GmtPayment        string `json:"gmt_payment"`                   // 交易付款时间
	GmtRefund         string `json:"gmt_refund"`                    // 交易退款时间
	GmtClose          string `json:"gmt_close"`                     // 交易结束时间
	FundBillList      string `json:"fund_bill_list,omitempty"`      // 支付金额信息
	PassbackParams    string `json:"passback_params,omitempty"`     // 回传参数
	VoucherDetailList string `json:"voucher_detail_list,omitempty"` // 优惠券信息
}
