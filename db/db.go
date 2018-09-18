package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hohice/gin-web/db/mapper"
	"github.com/hohice/gin-web/pkg/setting"
	"github.com/hohice/gin-web/pkg/util/once"
)

var (
	db        *sql.DB
	dbConf    setting.Configs
	openOnce  once.Once
	closeOnce once.Once
)

func init() {
	configChan := make(chan struct{})
	setting.RegNotifyChannel(configChan)
	go func() {
		for {
			select {
			case _, ok := <-configChan:
				{
					if !ok {
						return
					} else {
						dbConf = setting.Config
						if dbConf.Database.Enable {
							openOnce.Do(Open)
						}
					}
				}
			}
		}
	}()
}

func Open() (err error) {
	if db, err = sql.Open(dbConf.Database.Dirver, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf.Database.Username,
		dbConf.Database.Password,
		dbConf.Database.Host,
		dbConf.Database.Dbname)); err != nil {
		return err
	}

	//db.SingularTable(true)
	db.SetMaxIdleConns(dbConf.Database.MaxIdleConn)
	db.SetMaxOpenConns(dbConf.Database.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(dbConf.Database.MaxLifeTime) * time.Second)

	return
}

func Close() {
	if db != nil {
		closeOnce.Do(db.Close)
	}
}

func GetUserMapper() *mapper.UserMapper {
	return &mapper.UserMapper{
		DB: db,
	}
}
