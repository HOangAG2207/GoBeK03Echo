package model

type User struct {
	Base

	Username string `json:"username" gorm:"column:username;uniqueIndex;not null"`
	Email    string `json:"email" gorm:"column:email;uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null;column:password"`

	Displayname string `json:"display_name" gorm:"column:display_name"`
}

type UserRegisterRequest struct {
	Username    string `json:"username" validate:"required,gt=0"`
	Displayname string `json:"display_name" validate:"required,gt=0"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,gte=8"`
}
type RegisterUserSwaggerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	User    User   `json: "data"`
}

func (User) TableName() string {
	return "users"
}
