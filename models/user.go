package models

type User struct {
	GormModel
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Level    int    `json:"level" gorm:"not null;default:1"`
}

func (User) TableName() string {
	return "users"
}
