package ali

import (
	"strings"

	"github.com/go-baa/common/modules/aliyun"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

// Send 发送
func Send(to, subject, body string) error {
	accessKeyID := setting.Config.MustString("aliyun.accessKeyId", "")
	accessKeySecret := setting.Config.MustString("aliyun.accessKeySecret", "")

	dm, err := aliyun.NewDM(accessKeyID, accessKeySecret)
	if err != nil {
		log.Errorf("初始化DM错误: %v", err)
		return err
	}

	account := setting.Config.MustString("ali.mail.account", "")
	fromAlias := setting.Config.MustString("ali.mail.alias", "")
	_, err = dm.SingleSendMail(aliyun.DMConfig{
		AccountName: account,
		FromAlias:   fromAlias,
		ToAddress:   strings.Split(to, ","),
		Subject:     subject,
		HtmlBody:    body,
	})

	return err
}
