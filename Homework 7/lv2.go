package main


import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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
//================================================注册==================================================
func Registe(c *gin.Context){
	username:=c.PostForm("username")
	password:=c.PostForm("password")
	fmt.Println("user:"+username+password)
	if UserSignup(username,password){
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"数据库Insert报错"})
	}else {
		c.JSON(200, gin.H{"status": http.StatusOK, "message": "注册成功"})
	}
}
//通过用户名和密码完成user表中注册操作
func UserSignup(username string,password string)bool {
	stmt,err:=DBConn().Prepare(
		"insert into user(username,password)values(?,?)")
	if err!=nil{
		fmt.Println("fail to insert")
		return false
	}
	_,err=stmt.Exec(username,password)
	if err!=nil{
		fmt.Println("fail to insert")
		return false
	}
	return false
}

//===============================================登录===================================================
func Login(c *gin.Context) {
	username := c.PostForm("username")
	Password := c.PostForm("Password")
	if UserSignin(username, Password) {
		c.SetCookie("username", username, 10, "/", "localhost", false, true)
		//第一个参数为 cookie 名；第二个参数为 cookie 值；第三个参数为 cookie 有效时长；第四个参数为 cookie 所在的目录；第五个为所在域，表示我们的 cookie 作用范围；第六个表示是否只能通过 https 访问；第七个表示 cookie 是否支持HttpOnly属性。

		c.JSON(200, gin.H{"status": http.StatusOK, "message": "登录成功"})
	} else {
		c.JSON(403, gin.H{"status": http.StatusForbidden, "message": "登录失败，用户名或密码错误"})
	}
}


//判断密码是否一致
func UserSignin(username string,Password string)bool{
	var password  string
	stmt,err:=DBConn().Prepare(
		"select *from user where username=? AND password=?  ")
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}
	rows,err:=stmt.Query(username,password)
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}
	rows.Next()
	rows.Scan(&password,&username)
	if Password ==password{
		stmt,err:=DBConn().Prepare(
			"update user set situation  = 1 where username = ? ")
		if err!=nil {
			return false
		}else {
			stmt.Query(username)
			return true
		}
	}else {
		return false
	}

}
//=======================================留言评论回复====================================================
//发表留言
func SendMessageSql(username string,message string,pid int) bool {
	stmt,err:=DBConn().Prepare(
		"insert into messages(username,message,pid)values(?,?,?) ")
	if err!=nil{
		fmt.Println("fail to insert")
		return false
	}
	defer stmt.Close()

	_,err=stmt.Exec(username,message,pid)
	if err!=nil{
		fmt.Println("fail to insert")
		return false
	}
	return true
}

//发送信息--net--从cookie里读取用户名//区分了pid是否为空（回帖与否）//有pid的key就考虑是在回帖，没有的视为pid=0
func SendMsg(c *gin.Context){
	username,err:=c.Cookie("username")
	fmt.Println("username"+username)
	if err != nil{
		c.JSON(500,gin.H{"status": http.StatusForbidden,"message":"cookie读取失败"})
		return
	}
	message:=c.PostForm("message")
	var pid_input string
	pid_input=c.PostForm("pid")
	if pid_input == "" {
		//空pid非回复
		pid := 0
		if SendMessageSql(username,message,pid){
			c.JSON(200, gin.H{"内容":message,"用户名":username})
		}else {
			c.JSON(403, gin.H{"status": http.StatusForbidden, "message": "发送失败","": pid})
		}
	}else {
		//有pid//回复
		pid,err := strconv.Atoi(pid_input)
		if err!=nil{
			fmt.Println(err)
		}
		if SendMessageSql(username,message,pid){
			c.JSON(200, gin.H{"内容":message,"用户名":username,"回复":pid})
		}else {
			c.JSON(403, gin.H{"status": http.StatusForbidden, "message": "发送失败"})
		}
	}

}
//===============================================注销==========================================
//situation为0时为没登陆，为1时是登录状态
func quitSql(username string)bool{
	stmt,err:=DBConn().Prepare(
		"update user set situation  = 0 where username = ? ")
	if err!=nil {
		return false
	}else {
		stmt.Query(username)
		return true
	}
	return true
}

func Quit(c *gin.Context)  {
	username,err := c.Cookie("username")
	if err!=nil{
		fmt.Println("fail to quit")
	}
	c.SetCookie("username", username, -1, "/", "localhost", false, true)
	if quitSql(username){
		c.JSON(200,gin.H{"status":http.StatusInternalServerError,"message":"注销成功"})
	}else{
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
//看信息
//=========================================点赞================================================
//1为点赞，0为取消点赞
func AgreeSql(username string,id int)  {
	{
		stmt, err := DBConn().Prepare(
			`update message set agreenum=1 where username=?`)
		if err != nil {
			fmt.Println(err)
		}
		stmt.Query(username)

	}
}

func agree(c *gin.Context){
	agree_id_input := c.PostForm("agree")
	agree_id,err := strconv.Atoi(agree_id_input)
	if err!=nil {
		fmt.Println(err)
	}
	username,err := c.Cookie("username")
	AgreeSql(username,agree_id)

}
//==============================================权限========================================
func power(c*gin.Context){

}

//==================================================================================================
//主函数
func main() {
	router:=gin.Default()
	router.POST("/registe",Registe)
	router.POST("/login",Login)
	router.POST("/sendmsg",SendMsg)
	router.GET("/quit",Quit)
	router.POST("/agree",agree)
	router.Run(":8000")
}
//结构体
type user struct {
	username string
	password string
	Situation int
}
type message struct {
	username string
	message  string
	pid int
	agreenum int
}

