package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hohice/gin-web/db/model"
	"github.com/jinzhu/gorm"

	"github.com/hohice/gin-web/pkg/setting"
	"github.com/hohice/gin-web/pkg/util/logger"
)

func createNoneExistDB(conf *setting.Configs) error {
	dbconf := fmt.Sprintf("%s:%s@tcp(%s)/", dbConf.Database.Username, dbConf.Database.Password, dbConf.Database.Host)
	if datebase, err := sql.Open(dbConf.Database.Dirver, dbconf); err != nil {
		return err
	} else {
		dbstr := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci", dbConf.Database.Dbname)
		if _, err := datebase.Exec(dbstr); err != nil {
			return err
		}
		datebase.Close()
	}
	return nil
}

func AutoMigrate() error {

	if !dbConf.Database.Enable {
		return nil
	}

	if err := createNoneExistDB(&dbConf); err != nil {
		return errors.New("create database failed!")
	}

	if db, err := gorm.Open(dbConf.Database.Dirver, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf.Database.Username,
		dbConf.Database.Password,
		dbConf.Database.Host,
		dbConf.Database.Dbname)); err != nil {
		return errors.New("connect to database failed!")
	} else {
		db.LogMode(dbConf.Debug)
		db.SetLogger(logger.DefaultLogger)

		db.SingularTable(true)
		db.DB().SetMaxIdleConns(dbConf.Database.MaxIdleConn)
		db.DB().SetMaxOpenConns(dbConf.Database.MaxOpenConn)

		//auto migrate User
		db.Set("gorm:table_options", "ENGINE=InnoDB").Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&model.User{})

	}
	Close()
	return nil
}

/*
func CreateProductTable() bool {
	db.DropTableIfExists("user","")

	db.Set("gorm:table_options", "ENGINE=InnoDB").Table("User").CreateTable(&model.User{})

	if len(db.GetErrors()) > 0 {
		return false
	}

	return true
}
*/
