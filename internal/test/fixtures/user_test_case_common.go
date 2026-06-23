package fixtures_test

import (
	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"gorm.io/gorm"
)

type UserCommonTestDB struct {
	base
}

func (u *UserCommonTestDB) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

func (u *UserCommonTestDB) GenerateData() error {
	db := u.db.Session(&gorm.Session{SkipHooks: true})

	users := []*model.User{
		{
			Base: model.Base{
				ID: "3f1c2e9a-8b7a-4c1f-9a2d-1e5f6a7b8c9d",
			},
			Username:    "hoang01",
			Email:       "hoang01@gmail.com",
			Password:    "123456",
			Displayname: "hoang01",
		},
		{
			Base: model.Base{
				ID: "7a9d2c4b-1f6e-4b3a-8c2d-5e7f9a1b2c3d",
			},
			Username:    "hoang02",
			Email:       "hoang02@gmail.com",
			Password:    "123456",
			Displayname: "hoang02",
		},
	}
	return db.CreateInBatches(users, 10).Error
}
