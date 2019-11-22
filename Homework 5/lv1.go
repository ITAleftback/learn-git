package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var accounts map[string]Users = make(map[string]Users)
func main() {
	router := gin.Default()
	router.POST("",func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"code":200,
			"message":"Hello guest",
		})
	})
	router.POST("/login", func(c *gin.Context) {
		var user Users
		err := c.Bind(&user)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		if v, ok := accounts[user.Username]; ok && v.Password == user.Password {
			message:="Hello "+v.Username
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"message": message,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message" : "账号或者密码有误",
			})
		}
	})

	router.POST("/register", func(c *gin.Context) {
		var user Users
		err := c.Bind(&user)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		username := user.Username
		if _, ok := accounts[username]; ok {
			message := "用户名" + username + "已存在"
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"message": message,
			})
		} else {
			accounts[username] = user
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"message": "注册成功",
			})
		}
	})
	router.Run(":8000")
}
type Users struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}