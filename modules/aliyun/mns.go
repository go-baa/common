package aliyun

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"gopkg.in/baa.v1"
)

type MNSNotification struct {
	TopicOwner       string
	TopicName        string
	Subscriber       string
	SubscriptionName string
	MessageId        string
	Message          string
	MessageMD5       string
	MessageTag       string
	PublishTime      int
}

type MNS struct{}

// ValidateRequest 验证请求
func (t MNS) ValidateRequest(c *baa.Context) (*MNSNotification, error) {
	// 获取请求体
	content, err := c.Body().Bytes()
	if err != nil {
		return nil, err
	}

	// TODO 临时禁用校验
	if false {
		// 获取公钥
		cert, err := t.getCertificate(c)
		if err != nil {
			return nil, err
		}

		// 获取以 x-mns- 开头的请求头
		keys := make([]string, 0)
		for k := range c.Req.Header {
			if strings.HasPrefix(c.Req.Header.Get(k), "x-mns-") {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)

		// 获取请求方式和 URL
		method := c.Req.Method
		requestURI := c.Req.RequestURI
		contentMD5 := c.Req.Header.Get("Content-MD5")

		// 拼装为待加密的字符串
		items := make([]string, 0)
		items = append(items, strings.ToUpper(method))
		items = append(items, contentMD5)
		items = append(items, c.Req.Header.Get("Content-Type"))
		items = append(items, c.Req.Header.Get("Date"))
		for i := range keys {
			items = append(items, keys[i]+":"+c.Req.Header.Get(keys[i]))
		}
		items = append(items, requestURI)
		strToSign := strings.Join(items, "\n")

		// 获取待验证的签名
		authorization, err := base64.StdEncoding.DecodeString(c.Req.Header.Get("Authorization"))
		if err != nil {
			return nil, err
		}

		// 验证签名
		err = cert.CheckSignature(x509.SHA1WithRSA, []byte(strToSign), authorization)
		if err != nil {
			return nil, err
		}

		// 检查内容校验值
		hexStr := md5.Sum(content)
		md5Str := hex.EncodeToString(hexStr[:])
		if md5Str != contentMD5 {
			return nil, fmt.Errorf("Header field Content-MD5 verify failed, expect %s, got %s", md5Str, contentMD5)
		}
	}

	notify := new(MNSNotification)
	if err := json.Unmarshal(content, notify); err != nil {
		return nil, err
	}

	return notify, nil
}

func (t MNS) getCertificate(c *baa.Context) (*x509.Certificate, error) {
	certStr := c.Req.Header.Get("x-mns-signing-cert-url")
	if certStr == "" {
		return nil, fmt.Errorf("Header field %s required", "x-mns-signing-cert-url")
	}

	certURL, err := base64.StdEncoding.DecodeString(certStr)
	if err != nil {
		return nil, fmt.Errorf("Header field %s decode failed with error: %s", "x-mns-signing-cert-url", err.Error())
	}

	content, err := util.HTTPGet(string(certURL), 10)
	if err != nil {
		return nil, fmt.Errorf("Get public key cert failed with error: %s", err.Error())
	}

	cert, err := x509.ParseCertificate(content)
	if err != nil {
		return nil, fmt.Errorf("Parse public key cert failed with error: %s", err.Error())
	}

	return cert, nil
}

func NewMNS() *MNS {
	ins := new(MNS)
	return ins
}
