package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"git.code.tencent.com/xinhuameiyu/common/util"
)

const (
	SMSGateway = "https://dysmsapi.aliyuncs.com/"
)

// SMS 短信服务
type SMS struct {
	accessKeyId     string
	accessKeySecret string
}

func (t *SMS) getCommonParams() map[string]string {
	return map[string]string{
		"Format":           "JSON",
		"Version":          "2016-09-27",
		"AccessKeyId":      t.accessKeyId,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   string(util.RandStr(16, util.KC_RAND_KIND_ALL)),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}
}

// SMSResponse 短信响应
type SMSResponse struct {
	RequestId string
	Message   string
	Code      string
}

// SingleSendSms 发送短信
func (t *SMS) SingleSendSms(mobiles []string, signName string, templateCode string, templateParams map[string]string) (string, error) {
	params := t.getCommonParams()
	params["Action"] = "SendSms"
	params["Version"] = "2017-05-25"
	params["RegionId"] = "cn-hangzhou"
	params["PhoneNumbers"] = strings.Join(mobiles, ",")
	params["SignName"] = signName
	params["TemplateCode"] = templateCode

	enc, err := json.Marshal(templateParams)
	if err != nil {
		return "", err
	}
	params["TemplateParam"] = string(enc)

	// 签名
	sign := t.sign(http.MethodGet, params)
	params["Signature"] = sign

	// 获取响应
	query := url.Values{}
	for k := range params {
		query.Add(k, params[k])
	}
	res, err := t.request(SMSGateway, query)
	if err != nil {
		return "", err
	}

	fmt.Printf("res: %#v\n", string(res))

	// 解析响应
	ret := new(SMSResponse)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return "", err
	}
	if ret.Code != "OK" {
		return "", fmt.Errorf(ret.Message)
	}

	return ret.RequestId, nil
}

func (t *SMS) request(url string, query url.Values) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	// 超时设置
	client := new(http.Client)
	client.Timeout = time.Second * 10

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

func (t *SMS) sign(method string, params map[string]string) string {
	// 排序
	keys := make([]string, 0)
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接
	pairs := make([]string, 0)
	for i := range keys {
		key := keys[i]
		pairs = append(pairs, t.encode(key)+"="+t.encode(params[key]))
	}
	str := strings.ToUpper(method) + "&" + t.encodePercent("/") + "&" + t.encode(strings.Join(pairs, "&"))

	// 加密
	mac := hmac.New(sha1.New, []byte(t.accessKeySecret+"&"))
	mac.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (t *SMS) encode(str string) string {
	enc := url.QueryEscape(str)
	enc = strings.Replace(enc, "+", "%20", -1)
	enc = strings.Replace(enc, "*", "%2A", -1)
	enc = strings.Replace(enc, "%7E", "~", -1)
	return enc
}

func (t *SMS) encodePercent(str string) string {
	return strings.Replace(str, "/", "%2F", -1)
}

func NewSMS(accessKeyId, accessKeySecret string) (*SMS, error) {
	ins := new(SMS)

	if accessKeyId == "" {
		return nil, fmt.Errorf("Invalid SMS accessKeyId: %s", accessKeyId)
	}
	ins.accessKeyId = accessKeyId

	if accessKeySecret == "" {
		return nil, fmt.Errorf("Invalid SMS accessKeySecret: %s", accessKeySecret)
	}
	ins.accessKeySecret = accessKeySecret

	return ins, nil
}
