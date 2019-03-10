package alipay

import "fmt"

// 文档
// 0. https://cshall.alipay.com/support/help_detail.htm?help_id=486063&enctraceid=4_9UawjFSD2vJAfO2SyrmntJ763ax0f2uqEwpf_9c0o,
// 1. https://doc.open.alipay.com/doc2/detail.htm?spm=a219a.7629140.0.0.ukVaju&treeId=62&articleId=104743&docType=1#s7

type CreateDirectPayByUserOption struct {
	OutTradeNO       string  // 订单号
	Subject          string  // 商品名
	TotalFee         float64 // 交易金额
	SellerID         string  // 卖家支付宝用户号
	Price            float64 // 单价
	ShowURL          string  // 商品展示网址
	ExtraCommonParam string  // 公用回传参数
	ItBPay           string  // 超时时间
	QrPayMode        string  // 扫码支付方式
	QrcodeWidth      string  // 商户自定二维码宽度
	GoodsType        string  // 商品类型, 1表示实物类商品, 0表示虚拟类商品
	NotifyURL        string
	ReturnURL        string
}

func CreateDirectPayByUser(config *Config, option *CreateDirectPayByUserOption) string {
	params := map[string]string{}

	params["service"] = "create_direct_pay_by_user"
	params["partner"] = config.Partner
	params["_input_charset"] = "utf-8"
	params["sign_type"] = "MD5"
	params["notify_url"] = option.NotifyURL
	params["return_url"] = option.ReturnURL

	params["out_trade_no"] = option.OutTradeNO
	params["subject"] = option.Subject
	params["payment_type"] = "1"
	params["total_fee"] = fmt.Sprintf("%.2f", option.TotalFee)
	if len(option.SellerID) > 0 {
		params["seller_id"] = option.SellerID
	} else {
		params["seller_id"] = config.Partner
	}
	if option.Price > 0 {
		params["price"] = fmt.Sprintf("%.2f", option.Price)
	}
	params["show_url"] = option.ShowURL
	params["extra_common_param"] = option.ExtraCommonParam
	params["it_b_pay"] = option.ItBPay
	params["qr_pay_mode"] = option.QrPayMode
	params["qrcode_width"] = option.QrcodeWidth
	params["goods_type"] = option.GoodsType

	// 构建签名
	params["sign"] = BuildSign(params, config.MD5Key)

	return Gateway + "?" + BuildQuery(params)
}
