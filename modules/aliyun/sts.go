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
	STSGateway = "https://sts.aliyuncs.com/"
)

type STSCredentials struct {
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken"`
	Expiration      string `json:"Expiration"`
}

type STSAssumedRoleUser struct {
	Arn               string `json:"arn"`
	AssumedRoleUserId string `json:"AssumedRoleUserId"`
}

type STSResult struct {
	RequestId string `json:"RequestId"`
	HostId    string `json:"HostId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
}

type STSAssumeRoleResult struct {
	STSResult
	Credentials     *STSCredentials     `json:"Credentials"`
	AssumedRoleUser *STSAssumedRoleUser `json:"AssumedRoleUser"`
}

type STS struct {
	accessKeyId     string
	accessKeySecret string
}

func (s *STS) AssumeRole(roleArn, roleSessionName string, durationSeconds int) (*STSCredentials, *STSAssumedRoleUser, error) {
	return s.AssumeRoleWithPolicy(roleArn, roleSessionName, durationSeconds, "")
}

func (s *STS) AssumeRoleWithPolicy(roleArn, roleSessionName string, durationSeconds int, policy string) (*STSCredentials, *STSAssumedRoleUser, error) {
	params := s.getCommonParams()
	params["Action"] = "AssumeRole"
	params["RoleArn"] = roleArn
	params["RoleSessionName"] = roleSessionName
	params["DurationSeconds"] = fmt.Sprintf("%d", durationSeconds)

	// 指定策略
	if policy != "" {
		params["Policy"] = policy
	}

	// 签名
	sign := s.sign(http.MethodGet, params)
	params["Signature"] = sign

	// 获取响应
	query := url.Values{}
	for k := range params {
		query.Add(k, params[k])
	}
	res, err := s.request(STSGateway, query)
	if err != nil {
		return nil, nil, err
	}

	// 解析响应
	ret := new(STSAssumeRoleResult)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, nil, err
	}

	// 判断返回状态
	if ret.Code != "" {
		return nil, nil, fmt.Errorf(ret.Message)
	}

	return ret.Credentials, ret.AssumedRoleUser, nil
}

func (s *STS) request(url string, query url.Values) ([]byte, error) {
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

func (s *STS) getCommonParams() map[string]string {
	return map[string]string{
		"Format":           "JSON",
		"Version":          "2015-04-01",
		"AccessKeyId":      s.accessKeyId,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   string(util.RandStr(16, util.KC_RAND_KIND_ALL)),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}
}

func (s *STS) sign(method string, params map[string]string) string {
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
		pairs = append(pairs, s.encode(key)+"="+s.encode(params[key]))
	}
	str := strings.ToUpper(method) + "&" + s.encodePercent("/") + "&" + s.encode(strings.Join(pairs, "&"))

	// 加密
	mac := hmac.New(sha1.New, []byte(s.accessKeySecret+"&"))
	mac.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (s *STS) encode(str string) string {
	enc := url.QueryEscape(str)
	enc = strings.Replace(enc, "+", "%20", -1)
	enc = strings.Replace(enc, "*", "%2A", -1)
	enc = strings.Replace(enc, "%7E", "~", -1)
	return enc
}

func (s *STS) encodePercent(str string) string {
	return strings.Replace(str, "/", "%2F", -1)
}

func NewSTS(accessKeyId, accessKeySecret string) (*STS, error) {
	ins := new(STS)

	if accessKeyId == "" {
		return nil, fmt.Errorf("Invalid STS accessKeyId: %s", accessKeyId)
	}
	ins.accessKeyId = accessKeyId

	if accessKeySecret == "" {
		return nil, fmt.Errorf("Invalid STS accessKeySecret: %s", accessKeySecret)
	}
	ins.accessKeySecret = accessKeySecret

	return ins, nil
}
