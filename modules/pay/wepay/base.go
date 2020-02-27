package wepay

import (
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-baa/baa"
	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

type Config struct {
	AppID  string
	MchID  string
	MD5Key string
}

type CommonMessage struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

type OrderMessage struct {
	CommonMessage
	AppID              string     `xml:"appid"`
	MchID              string     `xml:"mch_id"`
	DeviceInfo         string     `xml:"device_info"`
	NonceStr           string     `xml:"nonce_str"`
	Sign               string     `xml:"sign"`
	ResultCode         string     `xml:"result_code"`
	ErrCode            string     `xml:"err_code"`
	ErrCodeDes         string     `xml:"err_code_des"`
	OpenID             string     `xml:"openid"`
	IsSubscribe        string     `xml:"is_subscribe"`
	TradeType          string     `xml:"trade_type"`
	BankType           string     `xml:"bank_type"`
	TotalFee           int        `xml:"total_fee"`
	SettlementTotalFee int        `xml:"settlement_total_fee"`
	FeeType            string     `xml:"fee_type"`
	CashFee            int        `xml:"cash_fee"`
	CashFeeType        string     `xml:"cash_fee_type"`
	CouponFee          int        `xml:"coupon_fee"`
	CouponCount        int        `xml:"coupon_count"`
	TransactionID      string     `xml:"transaction_id"`
	OutTradeNO         string     `xml:"out_trade_no"`
	Attach             string     `xml:"attach"`
	TimeEndStr         string     `xml:"time_end"`
	TimeEnd            *time.Time `xml:"-"`
}

func GetConfig() *Config {
	appID := setting.Config.MustString("pay.wepay.appid", "")
	if len(appID) == 0 {
		log.Errorln("获取微信支付支付配置失败：缺少 appid")
		return nil
	}

	mchID := setting.Config.MustString("pay.wepay.mch_id", "")
	if len(mchID) == 0 {
		log.Errorln("获取微信支付支付配置失败：缺少 mch_id")
		return nil
	}

	md5Key := setting.Config.MustString("pay.wepay.md5key", "")
	if len(md5Key) == 0 {
		log.Errorln("获取微信支付支付配置失败：缺少 md5key")
		return nil
	}

	return &Config{
		AppID:  appID,
		MchID:  mchID,
		MD5Key: md5Key,
	}
}

// BuildSign https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=4_3
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
		if keys[i] != "sign" {
			str += keys[i] + "=" + params[keys[i]] + "&"
		}
	}

	return strings.ToUpper(util.MD5(str + "key=" + key))
}

func GetRawMessage(c *baa.Context) ([]byte, error) {
	body, err := ioutil.ReadAll(c.Req.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Request(url, data string, timeout int) ([]byte, error) {
	// 创建请求对象
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	// 超时设置
	client := new(http.Client)
	client.Timeout = time.Duration(timeout) * time.Second

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
