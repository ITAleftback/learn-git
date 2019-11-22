package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
)

var accountsss map[string]student = make(map[string]student)
func main() {
	router := gin.Default()
	router.POST("/student?stuId=2016214229", func(c *gin.Context) {
		var user student
		err := c.Bind(&user)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		if v, ok := accountsss[user.stuID]; ok {

			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data":[
				{
					"studentID":"2016214229",
					"name":"匡俊嘉",
					"gender":  "男",
					"classID": "13001609",
					"major":   "软件工程",
					"college": "软件工程学院",
				}
			]
			})

		}

	})
	router.Run("8888")
}
type student struct {
	stuID int `form:"stuID" json:"stuID" binding:"required"`
}