package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

//=======================================有关数据库===================================================
//连接数据库
func init() {
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/user?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("fail to connect to db")
	}else {
		fmt.Println("连接数据库成功！")
	}
}
//数据库
var db *sql.DB

func DBConn()*sql.DB{
	return  db
}
//====================================用户==========================================
//签到
func qiandao(c *gin.Context) {
	username := c.PostForm("username")
	if Qiandao(username) {
		c.SetCookie("username", username, 10, "/", "localhost", false, true)
		//第一个参数为 cookie 名；第二个参数为 cookie 值；第三个参数为 cookie 有效时长；第四个参数为 cookie 所在的目录；第五个为所在域，表示我们的 cookie 作用范围；第六个表示是否只能通过 https 访问；第七个表示 cookie 是否支持HttpOnly属性。

		c.JSON(200, gin.H{"status": http.StatusOK, "message": "签到成功，积分+1"})
	} else {
		c.JSON(403, gin.H{"status": http.StatusForbidden, "message": "签到失败"})
	}
}
func Qiandao(username string)bool {

		stmt, err := db.Prepare("UPDATE qiandao SET jifen = 1 WHERE username=?")
		if err != nil {
			log.Fatal(err)
			return false
		}
		stmt.Exec(username);
		return true

}
//兑换奖品
func duihuan(c*gin.Context){
	prize:=c.PostForm("prize")
	jifen:=c.PostForm("jifen")
	jifenn,err:=strconv.Atoi(jifen)
	if err!=nil {
		fmt.Println("failed to strconv")
	}
	if change(prize,jifenn) {
		c.JSON(200, gin.H{"status": http.StatusOK, "message": "兑换成功"})
	} else {
		c.JSON(403, gin.H{"status": http.StatusForbidden, "message": "兑换失败"})
	}
}
func change(prize string,jifenn int)bool{
	rows, err :=db.Query("select id from prize where prize=? and jifen=?",prize,jifenn)
	if err != nil {
		fmt.Println("db.query is error login",err)
		return false
	}
	for rows.Next() {
		var jifen int
		err := rows.Scan(&jifen)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		if jifen>=0{
			return true
		}else{
			return false
		}
	}
	return false

}
//=======================================管理员===================================================
//查看用户积分
func check(c *gin.Context) {
	username := c.PostForm("username")
	if chakan(username) {
		c.SetCookie("username", username, 10, "/", "localhost", false, true)
		//第一个参数为 cookie 名；第二个参数为 cookie 值；第三个参数为 cookie 有效时长；第四个参数为 cookie 所在的目录；第五个为所在域，表示我们的 cookie 作用范围；第六个表示是否只能通过 https 访问；第七个表示 cookie 是否支持HttpOnly属性。

		c.JSON(200, gin.H{"status": http.StatusOK, "message": "查看成功"})
	} else {
		c.JSON(403, gin.H{"status": http.StatusForbidden, "message": "查看失败"})
	}
}
func chakan(username string)bool{
	stmt, err := db.Prepare("select  qiandao from qiandao where username=?")
	if err != nil{
		log.Fatal(err)
		return false
	}
	stmt.Exec(username);
	return true
}
//设置奖品
func setprize(c*gin.Context){
	prize:=c.PostForm("prize")
	jifen:=c.PostForm("jifen")
	if set(prize,jifen) {
		c.JSON(200, gin.H{"status": http.StatusOK, "message": "设置成功"})
	} else {
		c.JSON(403, gin.H{"status": http.StatusForbidden, "message": "设置失败"})
	}
}
func set(prize string,jifen string)bool{
	stmt,err:=DBConn().Prepare(
		"insert into prize(prize,jifen)values(?,?)")
	if err!=nil{
		fmt.Println("fail to insert")
		return false
	}
	_,err=stmt.Exec(prize,jifen)
	if err!=nil{
		fmt.Println("fail to insert")
		return false
	}
	return true
}
//===========================================================================================
func main() {
	router:=gin.Default()
	router.POST("/qiandao",qiandao)
	router.POST("/check",check)
	router.POST("/duihuan",duihuan)
	router.POST("/setprize",setprize)
	router.Run(":8020")
}
//结构体
type user struct {
	username string
	qiandao string
	id int
}
type guanliyuan struct {
	username string
}
type prize struct {
	prize string
	jifen int
	id int
}

