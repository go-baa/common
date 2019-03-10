package pay

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

const (
	// Gateway 接口地址
	Gateway = "https://api.mch.weixin.qq.com/pay/"
	// GatewaySandbox 沙箱接口地址
	GatewaySandbox = "https://api.mch.weixin.qq.com/sandboxnew/pay/"
)

// TimeLayout 时间格式
const TimeLayout = "20060102150405"

// WXPay 微信支付
type WXPay struct {
	appID      string
	appKey     string
	mchID      string
	sandbox    bool
	sandboxKey string
	apicert    []byte
	apikey     []byte
	rootca     []byte
}

// New 实例化WXpay
func New(appid, appkey, mchid string) (*WXPay, error) {
	ins := new(WXPay)
	if appid == "" {
		return nil, fmt.Errorf("Invalid appid")
	}
	ins.appID = appid

	if appkey == "" {
		return nil, fmt.Errorf("Invalid appkey")
	}
	ins.appKey = appkey

	if mchid == "" {
		return nil, fmt.Errorf("Invalid mchid")
	}
	ins.mchID = mchid

	return ins, nil
}

// SetAPICert 设置apiclient_cert
func (t *WXPay) SetAPICert(data []byte) {
	t.apicert = data
}

// SetAPIKey 设置apiclient_key
func (t *WXPay) SetAPIKey(data []byte) {
	t.apikey = data
}

// SetRootCA 设置rootca
func (t *WXPay) SetRootCA(data []byte) {
	t.rootca = data
}

// SetSandbox 设置沙箱状态
func (t *WXPay) SetSandbox(status bool, key string) {
	t.sandbox = status
	t.sandboxKey = key
}

// api 调用api
func (t *WXPay) api(service string, params map[string]string) ([]byte, error) {
	gateway := Gateway
	if t.sandbox {
		gateway = GatewaySandbox
	}
	url := gateway + service
	reqBody := t.buildXMLParams(params)

	if setting.Debug {
		log.Printf("WXPay api:%s reqbody:%s\n", service, reqBody)
	}

	res, err := t.request(url, "", []byte(reqBody), false)
	if err != nil {
		return nil, fmt.Errorf("请求错误:%v", err)
	}

	if setting.Debug {
		log.Printf("WXPay api:%s resbody:%s\n", service, string(res))
	}

	fmt.Printf("res: %v\n", string(res))

	return res, nil
}

// request http请求
func (t *WXPay) request(url string, query string, reqBody []byte, checkCert bool) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url+"?"+query, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// 超时设置
	client := new(http.Client)
	client.Timeout = time.Second * 60

	// https 支持
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	if checkCert {
		tlsCert, err := tls.X509KeyPair(t.apicert, t.apikey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}

		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM(t.rootca)
		if !ok {
			return nil, errors.New("failed to parse root certificate")
		}
		tlsConfig.RootCAs = pool
		tlsConfig.InsecureSkipVerify = false
	}

	client.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
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

// buildXMLParams 组装xml参数
func (t *WXPay) buildXMLParams(params map[string]string) string {
	params["appid"] = t.appID
	params["mch_id"] = t.mchID
	params["nonce_str"] = string(util.RandStr(32, util.KC_RAND_KIND_ALL))
	params["sign_type"] = "MD5"
	params["sign"] = t.Sign(params)

	return MapToXMLString(params)
}

// Sign 签名
func (t *WXPay) Sign(params map[string]string) string {
	newMap := make(map[string]string)
	for k, v := range params {
		if k == "sign" {
			continue
		}

		if v == "" {
			continue
		}
		newMap[k] = v
	}
	preSignStr := SortAndConcat(newMap)
	key := t.appKey
	if t.sandbox {
		key = t.sandboxKey
	}
	preSignWithKey := preSignStr + "&key=" + key

	md := md5.New()
	md.Write([]byte(preSignWithKey))
	sign := fmt.Sprintf("%X", md.Sum(nil))
	return strings.ToUpper(sign)
}

// CheckSign 响应签名验证
func (t *WXPay) CheckSign(data interface{}, sign string) error {
	resMap, err := StructToMap(data)
	if err != nil {
		return fmt.Errorf("签名验证 数据map转换错误:%v", err)
	}

	newsign := t.Sign(resMap)
	if sign != newsign {
		return fmt.Errorf("签名验证错误:got:%s,want:%s", newsign, sign)
	}

	return nil
}

// SortAndConcat map按键名排序并拼接
func SortAndConcat(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}

	var sortedParams []string
	sort.Strings(keys)
	for _, k := range keys {
		sortedParams = append(sortedParams, k+"="+params[k])
	}

	return strings.Join(sortedParams, "&")
}

// MapToXMLString map[string]string转xml
func MapToXMLString(params map[string]string) string {
	xml := "<xml>"
	for k, v := range params {
		xml = xml + fmt.Sprintf("<%s>%s</%s>", k, v, k)
	}
	xml = xml + "</xml>"

	return xml
}

// Map map[string]string
type Map map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// UnmarshalXML xml解码
func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = Map{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

// StructToMap struct转map[string]string
func StructToMap(in interface{}) (map[string]string, error) {
	out := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("StructToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagv := fi.Tag.Get("xml"); tagv != "" && tagv != "xml" {
			out[tagv] = v.Field(i).String()
		}
	}
	return out, nil
}

// TimeParse 时间转换
func TimeParse(timeStr string) (time.Time, error) {
	t, err := time.ParseInLocation(TimeLayout, timeStr, time.Local)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

// // 附着商户证书
// func (c *Client) WithCert(certFile, keyFile, rootcaFile string) error {
// 	cert, err := ioutil.ReadFile(certFile)
// 	if err != nil {
// 		return err
// 	}
// 	key, err := ioutil.ReadFile(keyFile)
// 	if err != nil {
// 		return err
// 	}
// 	rootca, err := ioutil.ReadFile(rootcaFile)
// 	if err != nil {
// 		return err
// 	}
// 	return c.WithCertBytes(cert, key, rootca)
// }

// func (c *Client) WithCertBytes(cert, key, rootca []byte) error {
// 	tlsCert, err := tls.X509KeyPair(cert, key)
// 	if err != nil {
// 		return err
// 	}
// 	pool := x509.NewCertPool()
// 	ok := pool.AppendCertsFromPEM(rootca)
// 	if !ok {
// 		return errors.New("failed to parse root certificate")
// 	}
// 	conf := &tls.Config{
// 		Certificates: []tls.Certificate{tlsCert},
// 		RootCAs:      pool,
// 	}
// 	trans := &http.Transport{
// 		TLSClientConfig: conf,
// 	}
// 	c.tlsClient = &http.Client{
// 		Transport: trans,
// 	}
// 	return nil
// }
