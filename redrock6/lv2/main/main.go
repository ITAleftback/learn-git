package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"log"
	."redrock6/lv2/init"
	"redrock6/lv2/str"
)




func main()  {
	router:=gin.Default()
	router.POST("/selectlesson",Selectlesson)
	router.Run(":8080")
}



func Selectlesson(c* gin.Context)  {
	//输入想选择的课
	//我给三个选择  A 语文 8-10点 B 数学8-10点 C 英语14到-16点
	lesson:=c.PostForm("lesson")
	var Time  string

	//判断选的哪个课
	switch {
	case lesson == "A" :
		lesson="Chinese"
		Time ="8-10"
	case lesson == "B" :
		lesson="Math"
		Time ="8-10"
	case lesson == "C" :
		lesson="English"
		Time ="14-16"
	default:
		fmt.Println("无效的选择")
		c.JSON(500,gin.H{"Status":http.StatusInternalServerError,"message":"无效的选择"})
		return
	}


	// 开启事务
	tx := DB.Begin()
	var s str.Selectlesson
	DB.Where("Time=?",Time).Find(&s)
	if s.ID>0{
		log.Println("课程时间冲突")
		c.JSON(500,gin.H{"Status":http.StatusInternalServerError,"message":"课程时间冲突"})
		return
	}

	if err := tx.Create(&str.Selectlesson{Lesson:lesson,Time:Time}).Error; err != nil {
		tx.Rollback()
		log.Println(err)
	}

	//提交事务
	tx.Commit()
}
