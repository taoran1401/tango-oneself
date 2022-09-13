package load

type Email struct {
	Host     string `mapstructure:"Host" json:"Host" yaml:"Host"`             // smtp host
	Port     int    `mapstructure:"Port" json:"Port" yaml:"Port"`             // smtp 端口
	IsSSL    bool   `mapstructure:"IsSSL" json:"IsSSL" yaml:"IsSSL"`          // smtp 是否开启ssl
	Username string `mapstructure:"Username" json:"Username" yaml:"Username"` // 账号，用于发邮件验证的账号，一般和发件人相同
	Password string `mapstructure:"Password" json:"Password" yaml:"Password"` // smtp授权码、密码（非邮箱密码）
	From     string `mapstructure:"From" json:"From" yaml:"From"`             // 发件人  你自己要发邮件的邮箱
	Nickname string `mapstructure:"Nickname" json:"Nickname" yaml:"Nickname"` // 昵称    发件人昵称 通常为自己的邮箱
}

//邮件模板
type EmailTemplate struct {
}

func NewEmailTemplate() *EmailTemplate {
	return &EmailTemplate{}
}

// verify模板
func (this *EmailTemplate) VerifyTemplate(code string) []byte {
	//验证
	content := []byte("<h1>当前验证码:" + code + "</h1>")
	return content
}

// notify模板
func (this *EmailTemplate) NotifyTemplate(msg string) []byte {
	//验证
	content := []byte("<h1>通知：" + msg + "</h1>")
	return content
}
