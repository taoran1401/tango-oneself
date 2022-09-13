package load

type App struct {
	AppName string `mapstructure:"AppName" json:"AppName" yaml:"AppName"`
	Debug   string `mapstructure:"Debug" json:"Debug" yaml:"Debug"`
	Port    string `mapstructure:"Port" json:"Port" yaml:"Port"`
}
