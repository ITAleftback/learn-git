package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reredrock5/lv2/par"
	"time"
)

func main() {
	router:=gin.Default()
	router.POST("/query",Query)
	router.Run(":8080")
}
//我就懒得在开一个包放post了
func Query(c *gin.Context)  {
	t1:=time.Now()
	for i:=1;i<5204 ;i++  {
		 value:=par.Parse(i)
		 for _,v:=range value{
		 	fmt.Println(v)
			 c.JSON(200,gin.H{"status":http.StatusOK,"message":v})
		 }
	}
	elapsed := time.Since(t1)
	fmt.Println("爬虫结束,总共耗时: ", elapsed)
}
