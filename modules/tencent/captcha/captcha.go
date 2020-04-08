package captcha

import (
	"errors"
	"fmt"

	captcha "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

// const (
// 	RequestDomain               = "https://captcha.tencentcloudapi.com"
// 	ActionDescribeCaptchaResult = "DescribeCaptchaResult"
// 	Version                     = "2019-07-22"
// 	CaptchaTypeSlid             = 9
// )

type ValidTicketParams struct {
	SecretId     string
	SecretKey    string
	CaptchaAppId uint64
	AppSecretKey string
	UserIp       string
	Ticket       string
	Randstr      string
}

func GetDescribeCaptchaResult(params ValidTicketParams) error {
	// credential := common.NewCredential(params.CaptchaSecretId, params.AppSecretKey)
	credential := common.NewCredential(params.SecretId, params.SecretKey)
	cli, err := captcha.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return err
	}
	req := captcha.NewDescribeCaptchaResultRequest()
	req.Ticket = &params.Ticket
	req.UserIp = &params.UserIp
	req.Randstr = &params.Randstr
	req.AppSecretKey = &params.AppSecretKey
	req.CaptchaAppId = &params.CaptchaAppId
	req.CaptchaType = common.Uint64Ptr(9)
	// http.CompleteCommonParams(req, regions.Beijing)
	resp, err := cli.DescribeCaptchaResult(req)
	if err != nil {
		return err
	}
	if *resp.Response.CaptchaCode != 1 {
		return errors.New(fmt.Sprintf("%d : %s", *resp.Response.CaptchaCode, *resp.Response.CaptchaMsg))
	}
	fmt.Println(resp.ToJsonString())
	return nil
}
