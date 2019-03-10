package alipay

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

const (
	Gateway = "https://mapi.alipay.com/gateway.do"
)

type Config struct {
	Partner string
	MD5Key  string
}

func GetConfig() *Config {
	partner := setting.Config.MustString("pay.alipay.partner", "")
	if len(partner) == 0 {
		log.Errorln("获取支付宝支付配置失败：缺少 partner")
		return nil
	}

	md5Key := setting.Config.MustString("pay.alipay.md5key", "")
	if len(md5Key) == 0 {
		log.Errorln("获取支付宝支付配置失败：缺少 md5key")
		return nil
	}

	return &Config{
		Partner: partner,
		MD5Key:  md5Key,
	}
}

func BuildQuery(params map[string]string) string {
	query := url.Values{}
	for key, value := range params {
		if len(value) > 0 {
			query.Add(key, value)
		}
	}
	return query.Encode()
}

func BuildSign(params map[string]string, key string) string {
	var keys []string
	for key, value := range params {
		if len(value) > 0 {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var str string
	for i := range keys {
		if keys[i] != "sign" && keys[i] != "sign_type" {
			str += keys[i] + "=" + params[keys[i]] + "&"
		}
	}

	return util.MD5(strings.TrimSuffix(str, "&") + key)
}

func ValidateRequest(tradeNO string, querys map[string]interface{}) error {
	var sign string
	var requestTradeNO string

	if v, ok := querys["sign"]; ok {
		sign = v.(string)
	}

	if v, ok := querys["out_trade_no"]; ok {
		requestTradeNO = v.(string)
	}

	if len(sign) == 0 || tradeNO != requestTradeNO {
		return fmt.Errorf("Invalid sign or out_trade_no")
	}

	config := GetConfig()
	if config == nil {
		return fmt.Errorf("Invalid alipay config")
	}

	params := map[string]string{}
	for i, v := range querys {
		if s, ok := v.(string); ok {
			params[i] = s
		}
	}

	if BuildSign(params, config.MD5Key) != sign {
		return fmt.Errorf("Invalid sign")
	}

	return nil
}
