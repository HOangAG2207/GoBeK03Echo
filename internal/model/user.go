package model

type User struct {
	Base

	Username string `json:"username" gorm:"column:username;uniqueIndex;not null"`
	Email    string `json:"email" gorm:"column:email;uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null;column:password"`

	DisplayName string `json:"display_name" gorm:"column:display_name"`
}

func (User) TableName() string {
	return "users"
}
