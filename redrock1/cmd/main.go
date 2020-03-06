package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"redrock1/user"
)

func main() {

	router := gin.Default()
	user.Setuprouter(router)
	router.Run(":8080")
}


