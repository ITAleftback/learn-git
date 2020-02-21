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

//主函数
//下面的每一个应该就是所谓的接口，我一般都在Postman本地运行
//127.0.0.1:8080
func main() {
	router:=gin.Default()
	router.POST("/registe",Registe)//注册 http://localhost:8080/registe?username=?&password=?
	router.POST("/login",Login)//登录   http://localhost:8080/login?username=?&password=?
	router.GET("/quit",Quit)//注销      http://localhost:8080/quit?username=?
	router.POST("/wrtessay",Wrtessay)//写文章 http://localhost:8080/wrtessay?username=?&essay=?
	router.POST("/askquestion",Askquestion)  //提问 http://localhost:8080/askquestion?username=?&question=?
    router.POST("/follow",Follow)//关注   http://localhost:8080/follow?username=?&follow=?
    router.POST("/anser",Anser)//回答   http://localhost:8080/anser?username=?&anser=?
    router.POST("/agree",Agree)//点赞   http://localhost:8080/agree?agree=?&anser=?
    router.POST("/disagree",Disagree)//踩  http://localhost:8080/disagree=?disagree=?&anser=?
    router.POST("/collect",Collect)//收藏 http://localhost:8080/collect?question=?&collect=?


	router.Run(":8080")
}


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
	password := c.PostForm("password")
	if checkUserSignin(username, password) {
		c.SetCookie("username", username, 100, "/", "localhost", false, true)
		//第一个参数为 cookie 名；第二个参数为 cookie 值；第三个参数为 cookie 有效时长；第四个参数为 cookie 所在的目录；第五个为所在域，表示我们的 cookie 作用范围；第六个表示是否只能通过 https 访问；第七个表示 cookie 是否支持HttpOnly属性。

		c.JSON(200, gin.H{"status": http.StatusOK, "message": "登录成功"})
	} else {
		c.JSON(403, gin.H{"status": http.StatusForbidden, "message": "登录失败，用户名或密码错误"})
	}
}
//判断密码是否一致
func checkUserSignin(username string,password string)bool {
	rows, err :=db.Query("select id from user where username=? and password=?",username,password)
	if err != nil {
		fmt.Println("db.query is error login",err)
		return false
	}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		if id>=0{
			return true
		}else{
			return false
		}
	}
	return false
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
}

func Quit(c *gin.Context)  {
	username,err := c.Cookie("username")
	if err!=nil{
		fmt.Println("fail to quit")
	}
	c.SetCookie("username", username, -1, "/", "localhost", false, true)
	if quitSql(username){
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"注销成功"})
	}else{
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
//========================================主页部分========================================================
//写文章
func Wrtessay(c *gin.Context){
	username,err:=c.Cookie("username")
	if err!=nil{
		log.Fatal(err)
	}
	essay:=c.PostForm("essay")
	if wrtessay(username,essay){
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"上传成功"})
	}else{
		c.JSON(403,gin.H{"status":http.StatusInternalServerError,"message":"上传失败"})
}
}
func wrtessay(username string,essay string)bool{
	stmt,err:=DBConn().Prepare(
		"insert into essay(username,essay)values(?,?)")//建一个名为essay的表，把文章放里面
	if err!=nil{
		fmt.Println("fail to insert")
		return false
	}
	_,err=stmt.Exec(username,essay)
	if err!=nil{
		fmt.Println("fail to insert")
		return false
	}
	return true
}
//===================================================================================================
//我写个提问，你看我写出来的东西都是上面写文章ctrl+c，真的感觉没多大区，但我还是把question与essay分别建一个表
func Askquestion(c *gin.Context){
	username,err:=c.Cookie("username")
	if err!=nil{
		log.Fatal(err)
	}
	question:=c.PostForm("question")
	if askquestion(username,question){
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"上传成功"})
	}else{
		c.JSON(403,gin.H{"status":http.StatusInternalServerError,"message":"上传失败"})
	}

}
func askquestion(username string,question string)bool{
	stmt,err:=DBConn().Prepare(
		"insert into question(username,question)values(?,?)")//建一个名为question的表，把question放里面
	if err!=nil{
		fmt.Println("fail to insert")
		return false
	}
	_,err=stmt.Exec(username,question)
	if err!=nil{
		fmt.Println("fail to insert")
		return false
	}
	return true
}
//====================================================================================================
//关注  我想的可能是和添加好友一样 我加一个取消关注
//这里要仔细注意 username是想要关注的人
//跟点赞一样   点两下也是取消
func Follow(c *gin.Context) {
	username:=c.PostForm("username")
	follow:= c.PostForm("follow")
	follownum, err := strconv.Atoi(follow)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	followSql(username,follownum)
}
func followSql(username string,follownum int) {
	rows, err :=db.Query("select follownum from friends where username=?",username)
	if err != nil {
		fmt.Println("db.query is error follow",err)
	}
	for rows.Next() {
		var follownum int
		err := rows.Scan(&follownum)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		if follownum==0{
			stmt,err:=DBConn().Prepare(
				`update friends set follownum=1 where username=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(username)

		}else{
			stmt,err:=DBConn().Prepare(
				`update friends set follownum=0 where username=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(username)

		}
	}
}
//=======================================问题详情页===================================================
//写个回答，  跟提问文章那些差不多 放在anser表里
func Anser(c*gin.Context){
	anser:=c.PostForm("anser")
	username,err:=c.Cookie("username")
	if err!=nil{
		log.Fatal(err)
	}
	if ansersql(username,anser){
		c.SetCookie("anser", anser, 100, "/", "localhost", false, true)
		//第一个参数为 cookie 名；第二个参数为 cookie 值；第三个参数为 cookie 有效时长；第四个参数为 cookie 所在的目录；第五个为所在域，表示我们的 cookie 作用范围；第六个表示是否只能通过 https 访问；第七个表示 cookie 是否支持HttpOnly属性。
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"回答成功"})
	}else{
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
func ansersql(username string,anser string)bool{
	stmt,err:=DBConn().Prepare(
		`insert into anser(username,anser)values=(?,?)`)
	if err!=nil{
		log.Fatal(err)
		return false
	}
	stmt.Exec(username,anser)
	return true
}

//==============================================================================================
//点赞&取消点赞  (第一次点赞，点两次取消点赞）
//1为点赞，0为取消点赞,在anser表里加value：agreenum
//这个不需要cookie把  不需要知道是谁点的赞
func Agree(c *gin.Context) {
	agree := c.PostForm("agree")
	anser:=c.PostForm("anser")
	agreenum, err := strconv.Atoi(agree)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	AgreeSql(anser, agreenum)
}
func AgreeSql(anser string,agreenum int) {
	rows, err :=db.Query("select agreenum from anser where anser=?",anser)
	if err != nil {
		fmt.Println("db.query is error agree",err)
	}
	for rows.Next() {
		var agreenum int
		err := rows.Scan(&agreenum)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		if agreenum==0{
			stmt,err:=DBConn().Prepare(
				`update anser set agreenum=1 where anser=?`)
				if err!=nil{
					fmt.Println(err)
				}
				stmt.Query(anser)

		}else{
			stmt,err:=DBConn().Prepare(
			`update anser set agreenum=0 where anser=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(anser)

		}
	}
}
//===============================================================================================
//踩跟点赞类似直接ctrl+c
//踩&取消踩  (第一次踩，点两次取消踩）
////1为踩，0为取消踩,在anser表里加value：disagreenum
func Disagree(c *gin.Context) {
	disagree := c.PostForm("disagree")
	anser:=c.PostForm("anser")
	disagreenum, err := strconv.Atoi(disagree)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	disagreeSql(anser, disagreenum)
}
func disagreeSql(anser string,disagreenum int) {
	rows, err :=db.Query("select disagreenum from anser where anser=?",anser)
	if err != nil {
		fmt.Println("db.query is error agree",err)
	}
	for rows.Next() {
		var disagreenum int
		err := rows.Scan(&disagreenum)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		if disagreenum==0{
			stmt,err:=DBConn().Prepare(
				`update anser set disagreenum=1 where anser=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(anser)

		}else{
			stmt,err:=DBConn().Prepare(
				`update anser set disagreenum=0 where anser=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(anser)

		}
	}
}

//==================================================================================================
//评论回复
//用到username的cookie  username就代表是谁评论的
func Comment(c*gin.Context){
	comment:=c.PostForm("comment")
	username,err:=c.Cookie("username")//用之前的cookie
	if err!=nil{
		log.Fatal(err)
	}
	if commentsql(username,comment){
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"评论成功"})
	}else{
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
func commentsql(username string,comment string)bool{
	stmt,err:=DBConn().Prepare(
		`insert into comment(username,comment)values(?,?)`)//又建一个comment表.... 拿来放评论  可是我有点疑惑 是不是只需把评论放anser表就行了？
	if err!=nil{
		log.Fatal(err)
		return false
	}
	stmt.Exec(username,comment)
	return true
}
//===============================================================================================
//收藏 跟关注一个道理
func Collect(c *gin.Context) {
	question:=c.PostForm("question")
	collect:= c.PostForm("collect")
	collectnum, err := strconv.Atoi(collect)
	if err != nil {
		log.Fatal(err)
	}
	followSql(question,collectnum)
}
func collectSql(question string,collectnum int) {
	rows, err :=db.Query("select collectnum from question where question=?",question)
	if err != nil {
		fmt.Println("db.query is error follow",err)
	}
	for rows.Next() {
		var collectnum int
		err := rows.Scan(&collectnum)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		if collectnum==0{
			stmt,err:=DBConn().Prepare(
				`update friends set collectnum=1 where question=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(question)

		}else{
			stmt,err:=DBConn().Prepare(
				`update friends set collectnum=0 where question=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(question)

		}
	}
}

//=====================================================================================================
//后面都是建立的结构体了
//又写了个coment  拿来放评论  可能不需要 只需要anser结构体？
type comment struct {
	username string
	comment string
}
//建一个anser表， 回答放着 点赞 踩  评论回复也放这
type anser struct{
	anser string
	agreenum int
	disagreenum int
	username string
}
//关注列表， friends表里面放的就是关注的人
type friends struct {
	follownum int
	username string
}
//这里建一个question表  提问全部放这
type question struct{
	username string
	question string
	collectnum int
}
//文章全放这
type essay struct{
	essay string
	username string
}
//存放用户的数据  用户名 密码 登录状态
type user struct {
	username string
	password string
	id int
	situation int
}

