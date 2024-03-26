package main

import (
	"9Kicks/config"
	"9Kicks/router"
	"fmt"
)

func main() {
	r := router.InitRouter()
	secrets := config.GetJWTSecrets()
	fmt.Println(secrets.JWTAdminSecret)
	r.Run(":8080")
}
