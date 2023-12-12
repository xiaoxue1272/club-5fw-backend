package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	name     string
	password string
}
