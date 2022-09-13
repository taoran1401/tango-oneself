package load

type Sms struct {
	AccessKeyId     string `mapstructure:"AccessKeyId" json:"AccessKeyId" yaml:"AccessKeyId"`
	AccessKeySecret string `mapstructure:"AccessKeySecret" json:"AccessKeySecret" yaml:"AccessKeySecret"`
	SignName        string `mapstructure:"SignName" json:"SignName" yaml:"SignName"`
}
