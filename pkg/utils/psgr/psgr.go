package psgr

import (
	"fmt"
	"gorm.io/gorm"
	"gpixivImageDownload/dao/dao"
	"gpixivImageDownload/dao/sql"
)

var Authdb *gorm.DB

func GetAuthTable() *gorm.DB {
	db := dao.GetClient()
	neu := sql.InitAuth()
	err := db.DB.AutoMigrate(&sql.Auth{})
	Authdb := db.DB.Create(neu)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(dv)
	return Authdb
}
