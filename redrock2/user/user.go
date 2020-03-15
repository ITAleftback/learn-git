package user

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	. "redrock2/init"
	"redrock2/jwt"
	"redrock2/str"
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
	//注册的同时 赋予这个账户3张票
	p:=str.Userpoll{
		Username: username,
		Poll:3,
	}
	DB.Create(&p)
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
//参赛===================================
func Joinrace(c*gin.Context){
	username:=c.GetString("username")
	_, err:= Conn.Do(
		"zadd",
		"poll",
		0,
		username,
	)
	if nil != err {
		return
	}else{
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"参赛成功"})
	}
}
/// 退赛===========================
func Quitrace(c*gin.Context){
	username:=c.GetString("username")
	_,err:=Conn.Do(
		"zrem",
		"poll",
		username,
		)
	if err != nil {
		return
	}else{
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"退赛成功"})
	}
}
// 投票======================================
func Like(c*gin.Context){
	username:=c.GetString("username")
	var p str.Userpoll
	DB.Where("username=?",username).Find(&p)
	poll:=p.Poll
	if poll==0{
		c.JSON(500,gin.H{"Status":http.StatusInternalServerError,"message":"你没有票了"})
		return
	}else {
		name:=c.PostForm("name")//  你想投给谁
		_,err:=Conn.Do(
			"zincrby",
			"poll",
			1,
			name,
		)
		if err != nil {
			return
		}else {
			c.JSON(200,gin.H{"status":http.StatusOK,"message":"投票成功"})
		}
		//票数-1
		DB.Model(&str.Userpoll{}).Where("username = ?",username).Update("poll",poll-1)
		return
	}
}
//每隔24h 更新票数为3  因为24h太久 为了加快效果我设置为每隔一分钟就更新一次票数
func Update(){
	DB.Model(&str.Userpoll{}).Update("poll",3)
	log.Println("更新成功吖~")
}
//排行榜==============================
// 根据票数从大到小
func List(c*gin.Context){
	list,err:=redis.Strings(Conn.Do(
		"zrevrange",
		"poll",
		0,
		-1,
		"withscores",
		))
	if err != nil {
		return
	}else{
		c.JSON(200,gin.H{"status":http.StatusOK,"message":list})
	}
}