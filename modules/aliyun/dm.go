package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
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
	DMGateway = "https://dm.aliyuncs.com/"
)

type DM struct {
	accessKeyId     string
	accessKeySecret string
}

// DMResponse 邮件发送响应结果
type DMResponse struct {
	RequestId string
	HostId    string
	Code      string
	Message   string
}

func (t *DM) getCommonParams() map[string]string {
	return map[string]string{
		"Format":           "JSON",
		"Version":          "2015-11-23",
		"AccessKeyId":      t.accessKeyId,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   string(util.RandStr(16, util.KC_RAND_KIND_ALL)),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}
}

// DMConfig 发信参数
type DMConfig struct {
	AccountName    string
	FromAlias      string
	ReplyToAddress bool
	ToAddress      []string
	Subject        string
	HtmlBody       string
}

// SingleSendMail 发送邮件
func (t *DM) SingleSendMail(config DMConfig) (*DMResponse, error) {
	params := t.getCommonParams()
	params["Action"] = "SingleSendMail"

	// 检查必填参数
	if config.AccountName == "" {
		return nil, errors.New("发信账号不能为空：AccountName")
	}
	params["AccountName"] = config.AccountName

	if config.FromAlias != "" {
		params["FromAlias"] = config.FromAlias
	}

	if config.ReplyToAddress {
		params["ReplyToAddress"] = "true"
	} else {
		params["ReplyToAddress"] = "false"
	}

	params["AddressType"] = "0"

	if len(config.ToAddress) == 0 {
		return nil, errors.New("收信人不能为空：ToAddress")
	}
	params["ToAddress"] = strings.Join(config.ToAddress, ",")

	if config.Subject == "" {
		return nil, errors.New("主题不能为空：Subject")
	}
	params["Subject"] = config.Subject

	if config.HtmlBody == "" {
		return nil, errors.New("内容不能为空：HtmlBody")
	}
	params["HtmlBody"] = config.HtmlBody

	// 签名
	sign := t.sign(http.MethodGet, params)
	params["Signature"] = sign

	// 获取响应
	query := url.Values{}
	for k := range params {
		query.Add(k, params[k])
	}
	res, err := t.request(DMGateway, query)
	if err != nil {
		return nil, err
	}

	// 解析响应
	ret := new(DMResponse)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, err
	}

	if ret.Code != "" {
		return nil, errors.New(ret.Message)
	}

	return ret, nil
}

func (t *DM) request(url string, query url.Values) ([]byte, error) {
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

func (t *DM) sign(method string, params map[string]string) string {
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

func (t *DM) encode(str string) string {
	enc := url.QueryEscape(str)
	enc = strings.Replace(enc, "+", "%20", -1)
	enc = strings.Replace(enc, "*", "%2A", -1)
	enc = strings.Replace(enc, "%7E", "~", -1)
	return enc
}

func (t *DM) encodePercent(str string) string {
	return strings.Replace(str, "/", "%2F", -1)
}

func NewDM(accessKeyId, accessKeySecret string) (*DM, error) {
	ins := new(DM)

	if accessKeyId == "" {
		return nil, fmt.Errorf("Invalid DM accessKeyId: %s", accessKeyId)
	}
	ins.accessKeyId = accessKeyId

	if accessKeySecret == "" {
		return nil, fmt.Errorf("Invalid DM accessKeySecret: %s", accessKeySecret)
	}
	ins.accessKeySecret = accessKeySecret

	return ins, nil
}
