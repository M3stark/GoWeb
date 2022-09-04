package dao // Data Access Object

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全局DB变量
var DB *gorm.DB

// InitMySQL 初始化MySQL
func InitMySQL() (err error) {
	dsn := "root:mysql_password@tcp(127.0.0.1:3306)/myWeb?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

// Close defer 关闭MySQL
func Close() {
	// https://gorm.io/docs/generic_interface.html
	sqlDB, _ := DB.DB()
	err := sqlDB.Close()
	if err != nil {
		return
	}
}
