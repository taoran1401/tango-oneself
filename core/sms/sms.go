package sms

import (
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Sms struct {
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
}

func NewSms(AccessKeyId string, AccessKeySecret string, SignName string) *Sms {
	return &Sms{
		AccessKeyId:     AccessKeyId,
		AccessKeySecret: AccessKeySecret,
		SignName:        SignName,
	}
}

//phone: 手机号
//TemplateCode: 模板id
//TemplateParam: 模板内容
func (this *Sms) Send(Phone string, TemplateCode string, TemplateParam string) (err error) {
	client, err := this.CreateClient()
	if err != nil {
		return err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(this.SignName),
		PhoneNumbers:  tea.String(Phone),
		TemplateCode:  tea.String(TemplateCode),
		TemplateParam: tea.String(TemplateParam),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return err
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		err = errors.New(error.Error())
	}
	return err
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func (this *Sms) CreateClient() (result *dysmsapi20170525.Client, err error) {
	config := &openapi.Config{
		AccessKeyId:     &this.AccessKeyId,
		AccessKeySecret: &this.AccessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	result = &dysmsapi20170525.Client{}
	result, err = dysmsapi20170525.NewClient(config)
	return result, err
}
