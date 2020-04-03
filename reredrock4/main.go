package main

import "redrock4/gin"

func main(){
	router:=gin.Default()
	router.GET("/book",QueryBook)
	router.GET("/cake",QueryCake)

	router.Run(8080)
}

func QueryBook(c *gin.Context){
	bid:=c.Query("id")
	c.String("your book id is "+bid)
	c.JSON(200,gin.H{"message":bid})
}

func QueryCake(c *gin.Context){
	c.String("cake cake cake")
}