package core

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"taogin/config/global"
)

type GormMysql struct {
}

func NewGormMysql() *GormMysql {
	return &GormMysql{}
}

//连接
func (this *GormMysql) Conn(uri string, maxIdleConns int, maxOpenConns int) *gorm.DB {
	if db, err := gorm.Open("mysql", uri); err != nil {
		panic(err)
	} else {
		sqlDB := db.DB()
		sqlDB.SetMaxIdleConns(maxIdleConns) //空闲中的最大连接数
		sqlDB.SetMaxOpenConns(maxOpenConns) //打开到数据库的最大连接数
		return db
	}
}

//连接多个mysql
func (this *GormMysql) ConnList() map[string]*gorm.DB {
	dbMap := make(map[string]*gorm.DB)
	for _, info := range global.CONFIG.Mysql {
		dbMap[info.AliasName] = this.Conn(info.MysqlConf.Uri, info.MysqlConf.MaxIdleConns, info.MysqlConf.MaxOpenConns)
	}
	return dbMap
}
