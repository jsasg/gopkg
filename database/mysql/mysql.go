package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlCfg struct {
	Host     string
	Port     int
	Username string
	Password string
	Dbname   string
}

func New(cfg MysqlCfg) gorm.Dialector {
	var (
		dsn string
	)
	dsn = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
	)
	return mysql.Open(dsn)
}
