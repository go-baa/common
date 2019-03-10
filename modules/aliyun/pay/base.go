package pay

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

const (
	// Gateway 接口地址
	Gateway = "https://openapi.alipay.com/gateway.do"
	// GatewaySandbox 沙箱接口地址
	GatewaySandbox = "https://openapi.alipaydev.com/gateway.do"
)

const (
	// TimeFormat 时间格式
	TimeFormat = "2006-01-02 15:04:05"
	// DataFormat 数据格式
	DataFormat = "JSON"
	// Charset 数据编码
	Charset = "utf-8"
	// SuccessCode 成功状态码
	SuccessCode = "10000"
	// ResponseSuffix 前缀
	ResponseSuffix = "_response"
	// ErrorResponse ...
	ErrorResponse = "error_response"
	// SignNodeName ...
	SignNodeName = "sign"
	// APIVersion api版本
	APIVersion = "1.0"
	// SignTypeRSA2 rsa2签名
	SignTypeRSA2 = "RSA2"
	// SignTypeRSA rsa签名
	SignTypeRSA = "RSA"
)

// AliPay 微信支付
type AliPay struct {
	appID      string
	sellerID   string
	publicKey  []byte
	privateKey []byte
	aliPubkey  []byte
	sandbox    bool
}

// New 实例化支付宝
func New(appid, sellerid string, pubKey, privKey, aliPubkey []byte) (*AliPay, error) {
	ins := new(AliPay)
	if appid == "" {
		return nil, fmt.Errorf("Invalid appid")
	}
	ins.appID = appid

	if sellerid == "" {
		return nil, fmt.Errorf("Invalid sellerid")
	}
	ins.sellerID = sellerid

	if len(pubKey) == 0 {
		return nil, fmt.Errorf("Invalid pubKey")
	}
	ins.publicKey = pubKey

	if len(privKey) == 0 {
		return nil, fmt.Errorf("Invalid privKey")
	}
	ins.privateKey = privKey

	if len(aliPubkey) == 0 {
		return nil, fmt.Errorf("Invalid aliPubKey")
	}
	ins.aliPubkey = aliPubkey

	return ins, nil
}

// GetAppID 获取appid
func (t *AliPay) GetAppID() string {
	return t.appID
}

// GetSellerID 获取sellerid
func (t *AliPay) GetSellerID() string {
	return t.sellerID
}

// SetSandbox 设置沙箱开启状态
func (t *AliPay) SetSandbox(status bool) {
	t.sandbox = status
}

func (t *AliPay) doRequest(param AliPayParam, results interface{}) error {
	gateway := Gateway
	if t.sandbox {
		gateway = GatewaySandbox
	}

	reqBody := t.buildParams(param)
	if setting.Debug {
		log.Printf("AliPay api:%s reqbody:%s\n", param.APIName(), reqBody)
	}

	fmt.Printf("req: %s\n", reqBody)

	res, err := t.request(gateway, "", []byte(reqBody))
	if err != nil {
		return fmt.Errorf("请求错误:%v", err)
	}

	fmt.Printf("res: %s\n", string(res))

	if setting.Debug {
		log.Printf("AliPay api:%s resbody:%s\n", param.APIName(), string(res))
	}

	// 响应签名验证
	if len(t.aliPubkey) > 0 {
		var dataStr = string(res)

		var rootNodeName = strings.Replace(param.APIName(), ".", "_", -1) + ResponseSuffix

		var rootIndex = strings.LastIndex(dataStr, rootNodeName)
		var errorIndex = strings.LastIndex(dataStr, ErrorResponse)

		var content string
		var sign string

		if rootIndex > 0 {
			content, sign = parserJSONSource(dataStr, rootNodeName, rootIndex)
		} else if errorIndex > 0 {
			content, sign = parserJSONSource(dataStr, ErrorResponse, errorIndex)
		} else {
			return nil
		}

		if ok, err := t.VerifyResponseData([]byte(content), sign); ok == false {
			return fmt.Errorf("响应签名验证错误:%v", err)
		}
	}

	err = json.Unmarshal(res, results)
	if err != nil {
		return fmt.Errorf("响应JSON解码错误:%v", err)
	}

	return nil
}

// AliPayParam 支付宝参数
type AliPayParam interface {
	// 用于提供访问的 method
	APIName() string

	// 返回参数列表
	Params() map[string]string

	// 返回扩展 JSON 参数的字段名称
	ExtJSONParamName() string

	// 返回扩展 JSON 参数的字段值
	ExtJSONParamValue() string
}

func jsonMarshal(obj interface{}) string {
	var bytes, err = json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// buildParams 组装请求参数
func (t *AliPay) buildParams(param AliPayParam) string {
	var p = url.Values{}
	// 公共参数
	p.Add("app_id", t.appID)
	p.Add("method", param.APIName())
	p.Add("format", DataFormat)
	p.Add("charset", Charset)
	p.Add("sign_type", SignTypeRSA2)
	p.Add("timestamp", time.Now().Format(TimeFormat))
	p.Add("version", APIVersion)

	// 补充可选公共参数
	var ps = param.Params()
	for key, value := range ps {
		p.Add(key, value)
	}

	if len(param.ExtJSONParamName()) > 0 {
		p.Add(param.ExtJSONParamName(), param.ExtJSONParamValue())
	}

	var keys = make([]string, 0, 0)
	for key := range p {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	sign := t.SignRSA2(keys, p)
	p.Add("sign", sign)

	return p.Encode()
}

// request http请求
func (t *AliPay) request(url string, query string, reqBody []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	// 超时设置
	client := new(http.Client)
	client.Timeout = time.Second * 60

	// https 支持
	if strings.HasPrefix(url, "https") {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	// 执行请求
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 处理响应
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func parserJSONSource(rawData string, nodeName string, nodeIndex int) (content string, sign string) {
	var dataStartIndex = nodeIndex + len(nodeName) + 2
	var signIndex = strings.LastIndex(rawData, "\""+SignNodeName+"\"")
	var dataEndIndex = signIndex - 1

	var indexLen = dataEndIndex - dataStartIndex
	if indexLen < 0 {
		return "", ""
	}
	content = rawData[dataStartIndex:dataEndIndex]

	var signStartIndex = signIndex + len(SignNodeName) + 4
	sign = rawData[signStartIndex:]
	var signEndIndex = strings.LastIndex(sign, "\"}")
	sign = sign[:signEndIndex]

	return content, sign
}
