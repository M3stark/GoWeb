package models

// User 自定义用户表
type User struct {
	ID       uint   `gorm:"primaryKey; AUTO_INCREMENT"` // ID - 自增主键
	Username string `gorm:"not null; unique;"`
	Password string `gorm:"not null"`
	Email    string `gorm:"not null; unique;"`
	Mobile   string `gorm:"not null; unique;"`
}

// LoginUser 登陆信息
type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
