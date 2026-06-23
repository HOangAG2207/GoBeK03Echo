package model

type User struct {
	Base               //having ID, UpdatedAt, CreatedAt
	Username    string `gorm:"unique;not null;column:username" json:"username"`
	Email       string `gorm:"unique;not null;column:email" json:"email"`
	Password    string `gorm:"not null;column:password" json:"-"`
	DisplayName string `gorm:"column:display_name" json:"display_name"`
}

func (User) TableName() string {
	return "users"
}
