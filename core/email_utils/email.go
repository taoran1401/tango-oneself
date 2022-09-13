package email_utils

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

type Email struct {
	Host     string //smtp服务器
	Port     int    //smtp服务器端口
	Username string //邮箱账号
	Password string //注意 这里是写smtp的授权码
	From     string //发件人邮箱； 格式：xx <xxx@xxx.com>
	IsSSL    bool   //是否启用ssl
}

func NewEmail(Host string, Port int, Username string, Password string, Nickname string, From string, IsSSL bool) *Email {
	return &Email{
		Host:     Host,
		Port:     Port,
		Username: Username,
		Password: Password,
		From:     fmt.Sprintf("%s <%s>", Nickname, From),
		IsSSL:    IsSSL,
	}
}

//发送邮件
//subject： 主题
//to： 收件人
//attachFile： 附件
func (this *Email) Send(subject string, to string, body []byte, attachFile string) error {
	e := email.NewEmail()
	//发件人
	e.From = this.From
	//收件人
	e.To = []string{to}
	//e.Bcc = []string{bcc} //私密发送邮件
	//e.Cc = []string{cc}   //转发邮件
	//主题
	e.Subject = subject
	//邮件模板
	e.HTML = body
	//附件
	if len(attachFile) > 0 {
		e.AttachFile(attachFile)
	}
	//smtp信息
	hostAddr := fmt.Sprintf("%s:%d", this.Host, this.Port)
	//认证信息
	auth := smtp.PlainAuth("", this.Username, this.Password, this.Host)
	//发送
	var err error
	if this.IsSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: this.Host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}
