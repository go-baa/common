package cps

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-baa/common/modules/aliyun"
	"github.com/go-baa/common/util"
)

const (
	// CPSGateway 阿里云推送服务接入地址
	CPSGateway = "https://cloudpush.aliyuncs.com/"
)

// CPS 阿里云移动推送服务
type CPS struct {
	appKey          string
	accessKeyID     string
	accessKeySecret string
}

// Response 响应结果
type Response struct {
	RequestID string `json:"RequestId"`
	HostID    string `json:"HostId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
}

// New 推送服务
func New(appKey, accessKeyID, accessKeySecret string) (*CPS, error) {
	ins := new(CPS)

	if appKey == "" {
		return nil, fmt.Errorf("Invalid CPS appKey")
	}
	ins.appKey = appKey

	if accessKeyID == "" {
		return nil, fmt.Errorf("Invalid CPS accessKeyID")
	}
	ins.accessKeyID = accessKeyID

	if accessKeySecret == "" {
		return nil, fmt.Errorf("Invalid CPS accessKeySecret")
	}
	ins.accessKeySecret = accessKeySecret

	return ins, nil
}

func (t *CPS) request(url string, query url.Values) ([]byte, error) {
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

func (t *CPS) sign(method string, params map[string]string) string {
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
		pairs = append(pairs, encode(key)+"="+encode(params[key]))
	}
	str := strings.ToUpper(method) + "&" + encodePercent("/") + "&" + encode(strings.Join(pairs, "&"))

	// 加密
	mac := hmac.New(sha1.New, []byte(t.accessKeySecret+"&"))
	mac.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func encode(str string) string {
	enc := url.QueryEscape(str)
	enc = strings.Replace(enc, "+", "%20", -1)
	enc = strings.Replace(enc, "*", "%2A", -1)
	enc = strings.Replace(enc, "%7E", "~", -1)
	return enc
}

func encodePercent(str string) string {
	return strings.Replace(str, "/", "%2F", -1)
}

func (t *CPS) getCommonParams() map[string]string {
	return map[string]string{
		"Format":           "JSON",
		"RegionId":         aliyun.LocationCnHangzhou,
		"Version":          "2016-08-01",
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   string(util.RandStr(16, util.KC_RAND_KIND_ALL)),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"AppKey":           t.appKey,
		"AccessKeyId":      t.accessKeyID,
	}
}
