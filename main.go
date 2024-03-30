package main

import (
	"9Kicks/router"
	"fmt"
)

func main() {
	route := router.SetupRouter()
	if err := route.Run(":8080"); err != nil {
		fmt.Printf("startup server failed,err: %v", err)
	}
}
