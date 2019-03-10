package unionpay

import (
	"strings"
	"time"
)

var (
	merID      string
	frontURL   string
	encoding   = "utf-8"
	signMethod = "01"
	version    = "5.1.0"
	baseURL    = "https://gateway.test.95516.com/"
)

// Config 初始使用的配置
type Config struct {
	// 版本号 默认5.1.0
	Version string

	// 请求银联的地址
	URL string

	// 商户代码
	MerID string

	// 验签私钥证书地址，传入pfx此路径可不传
	// openssl pkcs12 -in xxxx.pfx -nodes -out server.pem 生成为原生格式pem 私钥
	// openssl rsa -in server.pem -out server.key  生成为rsa格式私钥文件
	PrivateKey []byte

	// 验签证书地址,传入pfx此路径可以不传
	// openssl pkcs12 -in xxxx.pfx -clcerts -nokeys -out key.cert
	CertKey []byte

	// 加密证书地址
	EncryptCertKey []byte

	// 前端回调地址
	FrontURL string

	// 后端回调地址
	BackURL string
}

// Init init
func Init(config *Config) error {
	if err := LoadCert(config); err != nil {
		return err
	}
	SetConfig(config)
	return nil
}

// SetConfig 设置用户配置
func SetConfig(config *Config) {
	merID = config.MerID
	frontURL = config.FrontURL
	if config.Version != "" {
		version = config.Version
	}
	if config.URL != "" {
		baseURL = config.URL
	}
}

func sysParams(c APIConfig, data *RequestParams) map[string]string {
	request := map[string]string{
		"version":       version,
		"encoding":      encoding,
		"certId":        certData.CertID,
		"signMethod":    signMethod,
		"encryptCertId": certData.EncryptID,
		"accessType":    c.accessType,
		"channelType":   c.channelType,
		"bizType":       c.bizType,
		"merId":         merID,
		"backUrl":       "http://wwww.badiu.com",
	}

	request["txnTime"] = time.Now().Format("20060102150405")
	request["orderId"] = data.OrderID
	request["accNo"] = getaccNo(data.AccNo)
	request["customerInfo"] = getCustomerInfo(data.Customer)
	if data.Extend != "" {
		request["reqReserved"] = data.Extend
	}
	if data.Reserved != nil {
		list := []string{}
		for k, v := range data.Reserved {
			list = append(list, k+"&"+v)
		}
		if len(list) > 0 {
			request["reserved"] = "{" + strings.Join(list, "&") + "}"
		}
	}

	return request
}

func getTxnTime() string {
	return sec2Str("20060102150405", getNowSec())
}
func getaccNo(no string) string {
	str, _ := EncryptData(no)
	return str
}
func getCustomerInfo(customer *CustomerInfo) string {
	enmap := map[string]string{}
	other := map[string]string{}
	m := obj2Map(*customer)
	for k, v := range m {
		if v.(string) != "" {
			if k == "phoneNo" || k == "cvn2" || k == "expired" {
				enmap[k] = v.(string)
			} else {
				other[k] = v.(string)
			}
		}
	}
	if len(enmap) > 0 {
		tmp := []string{}
		for k, v := range enmap {
			tmp = append(tmp, k+"="+v)
		}
		str := strings.Join(tmp, "&")
		enc, _ := EncryptData(str)
		other["encryptedInfo"] = enc
	}
	tmp := []string{}
	for k, v := range other {
		tmp = append(tmp, k+"="+v)
	}
	return base64Encode([]byte("{" + strings.Join(tmp, "&") + "}"))
}
