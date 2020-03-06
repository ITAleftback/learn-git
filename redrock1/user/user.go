package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	. "redrock1/init"
	"redrock1/jwt"
	"redrock1/str"
)

//注册===================================================================
func Register(c *gin.Context)  {
	username := c.PostForm("username")
	password := c.PostForm("password")
	u:=str.User{
		Username: username,
		Password: password,
	}
	var user str.User
	DB.Where("username=?",username).Find(&user)
	if user.ID != 0 {
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"用户名已存在"})
		return
	}
	DB.Create(&u)
	c.JSON(200,gin.H{"status":http.StatusOK,"message":"注册成功"})

}
//登录=======================================
func Login(c *gin.Context) {
	username:=c.PostForm("username")
	password:=c.PostForm("password")
	var u str.User
	DB.Where("username=? AND password=?",username,password).Find(&u)
	if u.ID>0{
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"登陆成功"})
		//登录时创建这个用户名的token
		signtoken:=jwt.Create(username,u.ID)
		fmt.Println(signtoken)

		return
	}else {
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"密码错误"})
		return
	}
}
//查询用户=================  我这里就查询 寒假作业里的用户的个人介绍
func Selectuser(c*gin.Context){
	username:=c.PostForm("username")//输入你想查询的用户
	var i str.Instruct
	DB.Where("username=?",username).Find(&i)
	c.JSON(200,gin.H{"status":http.StatusOK,"message":i.Instruct})

}
// 修改个人信息  我这里就直接用寒假作业里的  修改个人昵称
func Updateusername(c*gin.Context){
	newusername:=c.PostForm("newusername")
	username:=c.MustGet("username")
	DB.Model(&str.User{}).Where("username = ?",username).Update("username",newusername)
	c.JSON(200,gin.H{"status":http.StatusOK,"message":"修改成功"})
}