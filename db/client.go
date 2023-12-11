package db

import (
	"github.com/xiaoxue1272/club-5fw-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDatabase(dbConfig *config.DataBaseConfiguration) {
	var err error
	db, err = gorm.Open(mysql.Open(dbConfig.Dns))
	if err != nil {
		panic(err)
	}
}
