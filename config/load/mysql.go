package load

type Mysql struct {
	AliasName string `mapstructure:"AliasName" json:"AliasName" yaml:"AliasName"`
	MysqlConf `yaml:",inline" mapstructure:",squash"`
}

type MysqlConf struct {
	Uri          string `mapstructure:"Uri" json:"Uri" yaml:"Uri"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
}
