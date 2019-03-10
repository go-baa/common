package apple

import (
	"encoding/json"
	"fmt"

	"github.com/go-baa/common/util"
)

const (
	// PayVerifyURL 正式验证地址
	PayVerifyURL = "https://buy.itunes.apple.com/verifyReceipt"
	// PaySendboxVerifyURL 沙箱验证地址
	PaySendboxVerifyURL = "https://sandbox.itunes.apple.com/verifyReceipt"

	// PayVerifyStatusOK 交易状态正常
	PayVerifyStatusOK = 0
)

// PayVerifyStatusDesc 验证结果状态描述
var PayVerifyStatusDesc = map[int]string{
	21000: "App Store无法读取提供的JSON数据",
	21002: "收据数据不符合格式",
	21003: "收据无法被验证",
	21004: "提供的共享密钥和账户的共享密钥不一致",
	21005: "收据服务器当前不可用",
	21006: "收据是有效的，但订阅服务已经过期",
	21007: "收据信息是测试用（sandbox），但却被发送到产品环境中验证",
	21008: "收据信息是产品环境中使用，但却被发送到测试环境中验证",
}

// PayVerifyRequest 验证请求数据
type PayVerifyRequest struct {
	ReceiptData string `json:"receipt-data"`
}

// PayVerifyReceipt 验证详情
type PayVerifyReceipt struct {
	AdamID                     int                          `json:"adam_id"`
	AppItemID                  int                          `json:"app_item_id"`
	ApplicationVersion         string                       `json:"application_version"`
	BundleID                   string                       `json:"bundle_id"`
	DownloadID                 int                          `json:"download_id"`
	OriginalApplicationAersion string                       `json:"Original_application_aersion"`
	OriginalPurchaseDate       string                       `json:"original_purchase_date"`
	OriginalPurchaseDateMs     string                       `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePst    string                       `json:"original_purchase_date_pst"`
	ReceiptCreationDate        string                       `json:"receipt_creation_date"`
	ReceiptCreationDateMs      string                       `json:"receipt_creation_date_ms"`
	ReceiptCreationDatePst     string                       `json:"receipt_creation_date_pst"`
	ReceiptType                string                       `json:"receipt_type"`
	RequestDate                string                       `json:"request_date"`
	RequestDateMs              string                       `json:"request_date_ms"`
	RequestDatePst             string                       `json:"request_date_pst"`
	VersionExternalIdentifier  int                          `json:"version_external_identifier"`
	InApp                      []*PayVerifyReceiptInappItem `json:"in_app"`
}

// PayVerifyReceiptInappItem 交易项目
type PayVerifyReceiptInappItem struct {
	IsTrialPeriod           string `json:"is_trial_period"`
	OriginalPurchaseDate    string `json:"original_purchase_date"`
	OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
	OriginalTransactionID   string `json:"original_transaction_id"`
	ProductID               string `json:"product_id"`
	PurchaseDate            string `json:"purchase_date"`
	PurchaseDateMs          string `json:"purchase_date_ms"`
	PurchaseDatePst         string `json:"purchase_date_pst"`
	Quantity                string `json:"quantity"`
	TransactionID           string `json:"transaction_id"`
}

// PayVerifyResponse 支付验证响应
type PayVerifyResponse struct {
	Environment string            `json:"environment,omitempty"`
	Receipt     *PayVerifyReceipt `json:"receipt"`
	Status      int               `json:"status"`
}

// PayVerify 支付验证
func PayVerify(transactionID, receipt string, sandbox bool) (*PayVerifyReceiptInappItem, int, error) {
	uri := PayVerifyURL
	if sandbox {
		uri = PaySendboxVerifyURL
	}

	reqData := &PayVerifyRequest{
		ReceiptData: receipt,
	}
	body, err := util.HTTPPostJSON(uri, reqData, 30)
	if err != nil {
		return nil, 0, err
	}

	response := new(PayVerifyResponse)
	if err = json.Unmarshal(body, response); err != nil {
		return nil, 0, fmt.Errorf("JSON解码错误:%v", err)
	}

	if response.Status != PayVerifyStatusOK {
		if desc, ok := PayVerifyStatusDesc[response.Status]; ok {
			return nil, util.MustInt(response.Status), fmt.Errorf("%s", desc)
		}
		return nil, response.Status, fmt.Errorf("未知错误")
	}

	for _, v := range response.Receipt.InApp {
		if v.TransactionID == transactionID {
			return v, 0, nil
		}
	}

	return nil, 0, fmt.Errorf("未找到交易信息")
}
