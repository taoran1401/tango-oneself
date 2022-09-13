package load

type JWT struct {
	SigningKey  string `mapstructure:"SigningKey" json:"SigningKey" yaml:"SigningKey"`    // jwt签名
	ExpiresTime int64  `mapstructure:"ExpiresTime" json:"ExpiresTime" yaml:"ExpiresTime"` // 过期时间
	BufferTime  int64  `mapstructure:"BufferTime" json:"BufferTime" yaml:"BufferTime"`    // 缓冲时间
	Issuer      string `mapstructure:"Issuer" json:"Issuer" yaml:"Issuer"`                // 签发者
}
