package database

import (
	"fmt"
	"log"

	"github.com/ronnrein/eps/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB gorm.DB

func init() {
	db, err := getDatabase(utils.Config)
	if err != nil {
		log.Panicln(err)
	}
	DB = db
}

func getDatabase(config utils.Conf) (gorm.DB, error) {
	connectionStr := fmt.Sprintf(
		"%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True", config.DBUser,
		config.DBPassword, config.DBProtocol, config.DBHost, config.DBPort, config.DBName,
	)

	db, err := gorm.Open("mysql", connectionStr)
	if err != nil {
		return gorm.DB{}, err
	}
	db.LogMode(false)

	return *db, nil
}
