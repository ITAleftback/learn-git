package user

import (
	"github.com/gin-gonic/gin"
	"redrock2/middleware"
)

func Setuprouter(router *gin.Engine){
	router.POST("/register",Register)
	router.POST("/login",Login)
	router.Use(middleware.User)
	{
		router.POST("/joinrace", Joinrace)//参赛
		router.POST("/quitrace", Quitrace)//退赛
		router.POST("/like",  Like)// 投票系统
		router.POST("/list", List)//  排行榜
	}
}
