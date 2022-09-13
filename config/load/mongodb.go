package load

type MongoDB struct {
	AliasName   string `mapstructure:"AliasName" json:"AliasName" yaml:"AliasName"`
	MongoDBConf `yaml:",inline" mapstructure:",squash"`
}

type MongoDBConf struct {
	Uri string `mapstructure:"Uri" json:"Uri" yaml:"Uri"`
}
