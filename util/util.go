package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// GetDir get pwd
func GetDir() string {
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		return ""
	}
	return filepath.Dir(path)
}

// MD5 checksum for str
func MD5(str string) string {
	hexStr := md5.Sum([]byte(str))
	return hex.EncodeToString(hexStr[:])
}

// MD5File checksum for file path
func MD5File(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		return ""
	}

	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return ""
	}

	hexStr := md5hash.Sum(nil)
	return hex.EncodeToString(hexStr[:])
}

// SHA1 sha1 Encrypted data
func SHA1(str string) string {
	sha1 := sha1.New()
	io.WriteString(sha1, string(str))
	return fmt.Sprintf("%x", sha1.Sum(nil))
}
