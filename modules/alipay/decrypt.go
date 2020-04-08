package alipay

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"net/url"
	"sort"
	"strings"
)

func sign(unsigned, key string) string {
	pk, err := parsePKCS8PrivateKey(formatPKCS8PrivateKey(key))
	if err != nil {
		return ""
	}
	bs, _ := rsaSignWithKey([]byte(unsigned), pk, crypto.SHA256)
	return base64.StdEncoding.EncodeToString(bs)
}

func rsaSignWithKey(src []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

func parsePKCS1PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("")
	}

	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, err
}

func parsePKCS8PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("")
	}

	rawKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := rawKey.(*rsa.PrivateKey)
	if ok == false {
		return nil, errors.New("")
	}

	return key, err
}

const (
	publicKeyPrefix = "-----BEGIN PUBLIC KEY-----"
	publicKeySuffix = "-----END PUBLIC KEY-----"

	pKCS1Prefix = "-----BEGIN RSA PRIVATE KEY-----"
	pKCS1Suffix = "-----END RSA PRIVATE KEY-----"

	pKCS8Prefix = "-----BEGIN PRIVATE KEY-----"
	pKCS8Suffix = "-----END PRIVATE KEY-----"
)

func formatPKCS1PrivateKey(raw string) []byte {
	raw = strings.Replace(raw, pKCS8Prefix, "", 1)
	raw = strings.Replace(raw, pKCS8Suffix, "", 1)
	return formatKey(raw, pKCS1Prefix, pKCS1Suffix, 64)
}

func formatPKCS8PrivateKey(raw string) []byte {
	raw = strings.Replace(raw, pKCS1Prefix, "", 1)
	raw = strings.Replace(raw, pKCS1Suffix, "", 1)
	return formatKey(raw, pKCS8Prefix, pKCS8Suffix, 64)
}

func formatKey(raw, prefix, suffix string, lineCount int) []byte {
	if raw == "" {
		return nil
	}
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.Replace(raw, " ", "", -1)
	raw = strings.Replace(raw, "\n", "", -1)
	raw = strings.Replace(raw, "\r", "", -1)
	raw = strings.Replace(raw, "\t", "", -1)

	var sl = len(raw)
	var c = sl / lineCount
	if sl%lineCount > 0 {
		c = c + 1
	}

	var buf bytes.Buffer
	buf.WriteString(prefix + "\n")
	for i := 0; i < c; i++ {
		var b = i * lineCount
		var e = b + lineCount
		if e > sl {
			buf.WriteString(raw[b:])
		} else {
			buf.WriteString(raw[b:e])
		}
		buf.WriteString("\n")
	}
	buf.WriteString(suffix)
	return buf.Bytes()
}

func sortKeys(param url.Values) string {
	if param == nil {
		param = make(url.Values, 0)
	}
	pList := make([]string, 0, 0)
	for key := range param {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	sort.Strings(pList)
	return strings.Join(pList, "&")
}

// AESDecrypt 支付宝解密
func AESDecrypt(ociphertext, okey string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ociphertext)
	if err != nil {
		return nil, err
	}
	key, err := base64.StdEncoding.DecodeString(okey)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, block.BlockSize())
	for index := 0; index < len(iv); index++ {
		iv[index] = 0x0
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext, err = PKCS7UnPadding(ciphertext)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// PKCS7Padding PKCS7填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding PKCS7去除填充
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	if unpadding > length {
		return nil, errors.New("padding out of range")
	}
	return origData[:(length - unpadding)], nil
}
