package dao

import (
	"blog/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Manager interface {
	AddUser(user *model.User)
}

type manager struct {
	db *gorm.DB
}

var Mgr Manager

func init() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/db_go?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to init db", err)
	}
	Mgr = &manager{db: db}
	//db.AutoMigrate(&model.User{})
}

func (mgr *manager) AddUser(user *model.User) {

	mgr.db.Create(user)
}
