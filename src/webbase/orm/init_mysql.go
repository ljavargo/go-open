package orm

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql.20160224"
	"github.com/go-xorm/xorm"

)

type MysqlConfig struct {
	DBDriver     string `json:"db_driver"`
	DBHost       string `json:"db_host"`
	DBName       string `json:"db_name"`
	DBUser       string `json:"db_user"`
	DBPassword   string `json:"db_pwd"`
	PrintSql     bool   `json:"print_sql"`
	DefaultEngin *xorm.Engine
}

func Init(cfg MysqlConfig) {
	if cfg.DefaultEngin != nil {
		OrmEngine = cfg.DefaultEngin
		return
	}
	var err error
	OrmEngine, err = xorm.NewEngine(cfg.DBDriver, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName))
	if err != nil {
		log.Fatal("new mysql OrmEngine failed", err)
	}

	if cfg.PrintSql {
		OrmEngine.ShowSQL(true)
	}
}
