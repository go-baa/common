package pay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"net/url"
	"sort"
	"strings"

	"github.com/go-baa/log"
)

// SignRSA2 rsa2签名
func (t *AliPay) SignRSA2(keys []string, param url.Values) (s string) {
	if param == nil {
		param = make(url.Values, 0)
	}

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var src = strings.Join(pList, "&")
	var sig, err = signPKCS1v15([]byte(src), t.privateKey, crypto.SHA256)
	if err != nil {
		log.Error(err)
		return ""
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s
}

func signPKCS1v15(src, key []byte, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)

	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error")
	}

	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, pri, hash, hashed)
}

func verifyPKCS1v15(src, sig, key []byte, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)

	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return errors.New("public key error")
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	var pub = pubInterface.(*rsa.PublicKey)

	return rsa.VerifyPKCS1v15(pub, hash, hashed, sig)
}

// VerifyResponseData 同步请求响应数据签名验证
func (t *AliPay) VerifyResponseData(data []byte, sign string) (ok bool, err error) {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	err = verifyPKCS1v15(data, signBytes, t.aliPubkey, crypto.SHA256)
	if err != nil {
		return false, err
	}

	return true, nil
}

// VerifyNotifySign 验证通知签名
func (t *AliPay) VerifyNotifySign(reqParams map[string]interface{}) (ok bool, err error) {
	sign, err := base64.StdEncoding.DecodeString(reqParams["sign"].(string))
	if err != nil {
		return false, err
	}
	signType := reqParams["sign_type"]

	var keys = make([]string, 0, 0)
	for k, v := range reqParams {
		if k == "sign" || k == "sign_type" {
			continue
		}

		if len(v.(string)) > 0 {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(reqParams[key].(string))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var s = strings.Join(pList, "&")
	if signType == SignTypeRSA {
		err = verifyPKCS1v15([]byte(s), sign, t.aliPubkey, crypto.SHA1)
	} else {
		err = verifyPKCS1v15([]byte(s), sign, t.aliPubkey, crypto.SHA256)
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
