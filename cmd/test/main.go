package main

import (
	"context"
	"fmt"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	userRepo "github.com/HOangAG2207/GoBeK03Echo/internal/repository/user"
	pkgdb "github.com/HOangAG2207/GoBeK03Echo/pkg/sqldb"
)

func main() {
	dbClient, err := pkgdb.NewClient("")
	if err != nil {
		panic(err)
	}

	dbClient.AutoMigrate(&model.User{})
	user := model.User{
		Username:    "hoàng",
		Email:       "hoang@gmail.com",
		Password:    "juwdhuwhd",
		DisplayName: "HOàng",
	}
	repo := userRepo.NewRepository(dbClient)
	create_user, err := repo.CreateUser(context.Background(), &user)
	if err != nil {
		panic(err)
	}
	fmt.Print(create_user)
}
