package main

import (
	"9Kicks/router"
)

func main() {
	r := router.InitRouter()
	r.Run(":8080")
}
