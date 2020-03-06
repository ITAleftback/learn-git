package user

import (
	"github.com/gin-gonic/gin"
	"redrock1/middleware"
)

func Setuprouter(router *gin.Engine){
	router.POST("/register",Register)
	router.POST("/login",Login)
	router.POST("/selectuser",Selectuser)
	router.POST("/updateusername",middleware.User,Updateusername)
}
