package load

type Zap struct {
	Level         string `mapstructure:"Level" json:"Level" yaml:"Level"`                         // 级别
	Prefix        string `mapstructure:"Prefix" json:"Prefix" yaml:"Prefix"`                      // 日志前缀
	Format        string `mapstructure:"Format" json:"Format" yaml:"Format"`                      // 输出
	Director      string `mapstructure:"Director" json:"Director"  yaml:"Director"`               // 日志文件夹
	EncodeLevel   string `mapstructure:"EncodeLevel" json:"EncodeLevel" yaml:"EncodeLevel"`       // 编码级
	StacktraceKey string `mapstructure:"StacktraceKey" json:"StacktraceKey" yaml:"StacktraceKey"` // 栈名

	MaxAge       int  `mapstructure:"MaxAge" json:"MaxAge" yaml:"MaxAge"`                   // 日志留存时间
	ShowLine     bool `mapstructure:"ShowLine" json:"ShowLine" yaml:"ShowLine"`             // 显示行
	LogInConsole bool `mapstructure:"LogInConsole" json:"LogInConsole" yaml:"LogInConsole"` // 输出控制台
}
