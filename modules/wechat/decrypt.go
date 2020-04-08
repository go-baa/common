package wechat

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// AESDecrypt 解密微信加密接口
func AESDecrypt(owords, okey, oiv string) ([]byte, error) {
	words, err := base64.StdEncoding.DecodeString(owords)
	if err != nil {
		return nil, err
	}
	key, err := base64.StdEncoding.DecodeString(okey)
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(oiv)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(words, words)
	return words, nil
}
