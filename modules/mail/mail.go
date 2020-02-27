package mail

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/go-baa/baa"
	"github.com/go-baa/common/modules/mail/ali"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

// Config 邮件配置
type Config struct {
	Addr  string // smtp.163.com:25
	User  string // robot@163.com
	Pass  string // 12345
	From  string // 发件人
	To    string // 收件人, 以逗号隔开的多个
	Title string // 邮件标题
	Body  string // 邮件正文
	Type  string // 内容类型，纯文本plain或网页html
}

// NewConfig 根据配置邮件获取发送配置
func NewConfig() *Config {
	return &Config{
		Addr:  setting.Config.MustString("mail.host", ""),
		User:  setting.Config.MustString("mail.user", ""),
		Pass:  setting.Config.MustString("mail.pass", ""),
		From:  setting.Config.MustString("mail.from", ""),
		To:    "",
		Title: "",
		Body:  "",
		Type:  "html",
	}
}

// SendMail 发送邮件
func SendMail(conf *Config) error {
	// 尝试阿里云发送
	err := ali.Send(conf.To, conf.Title, conf.Body)
	if err != nil {
		log.Errorf("阿里云发送邮件错误: %v\n", err)
		err = smtpSend(conf)
	}

	return err
}

// smtpSend SMTP发送
func smtpSend(conf *Config) error {
	var (
		contentType string
		vs          string
		message     string
		toaddr      mail.Address
	)
	encode := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	host := strings.Split(conf.Addr, ":")
	auth := smtp.PlainAuth("", conf.User, conf.Pass, host[0])
	if conf.Type == "html" {
		contentType = "text/html; charset=UTF-8"
	} else {
		contentType = "text/plain; charset=UTF-8"
	}
	from := mail.Address{Name: setting.Config.MustString("mail.name", "XinMa System"), Address: conf.From}
	tolist := strings.Split(conf.To, ",")
	var to []string
	for i, addr := range tolist {
		addr = strings.TrimSpace(addr)
		tolist[i] = addr
		toaddr = mail.Address{Name: "", Address: tolist[i]}
		to = append(to, toaddr.String())
	}

	header := make(mail.Header)
	header["From"] = []string{from.String()}
	header["To"] = to
	header["Subject"] = []string{conf.Title}
	header["MIME-Version"] = []string{"1.0"}
	header["Content-Type"] = []string{contentType}
	header["Content-Transfer-Encoding"] = []string{"base64"}

	for k, v := range header {
		vs = strings.Join(v, ", ")
		message += fmt.Sprintf("%s: %s\r\n", k, vs)
	}
	message += "\r\n" + encode.EncodeToString([]byte(conf.Body))

	err := smtp.SendMail(
		conf.Addr,
		auth,
		from.Address,
		tolist,
		[]byte(message),
	)
	return err
}

// SendTemplateMail 发送模板邮件
func SendTemplateMail(conf *Config, templatePath string, data map[string]interface{}) error {
	if templatePath == "" {
		return errors.New("模板路径不能为空")
	}
	if len(data) == 0 {
		return errors.New("数据不能为空")
	}

	b := baa.Default()
	buf := new(bytes.Buffer)
	if err := b.Render().Render(buf, templatePath, data); err != nil {
		log.Errorf("邮件渲染模板错误: %v\n", err)
		return err
	}

	conf.Body = string(buf.Bytes())
	return SendMail(conf)
}
