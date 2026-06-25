package main

import (
	pkgjwt "github.com/HOangAG2207/GoBeK03Echo/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	jwtGen, err := pkgjwt.NewGeneratorJWT("./private_key.pem")
	if err != nil {
		panic(err)
	}

	tokenStr, err := jwtGen.GenerateJWTToken(jwt.MapClaims{
		"sub":  "1234",
		"name": "user1234",
	})
	if err != nil {
		panic(err)
	}
	println(tokenStr)
}
