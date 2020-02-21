package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strings"
)

//主函数
//////0.3
//127.0.0.1:8080
func main() {

	router:=gin.Default()
	router.Use(Cors())
	router.POST("/registe",Registe)//注册 http://localhost:8080/registe?username=?&password=?
	router.POST("/login",Login)//登录   http://localhost:8080/login?username=?&password=?
	router.POST("/logout",Logout)//注销      http://localhost:8080/logout?username=?
	router.POST("/recommend",Recommend)//推荐  http://localhost:8080/recommend?
	router.POST("/search",Search)//搜索  http://localhost:8080/search?keyword=?
	router.POST("/wrtessay",Wrtessay)//写文章 http://localhost:8080/wrtessay?essay=?
	router.POST("/askquestion",Askquestion)  //提问 http://localhost:8080/askquestion?question=?
	router.POST("/follow",Follow)//关注   http://localhost:8080/follow?usernames=?
    router.POST("/notfollow",Notfollow)//取消关注  http://localhost:8080/notfollow?usernames=?
	router.POST("/anser",Anser)//回答   http://localhost:8080/anser?anser=?
    router.POST("/agree",Agree)//点赞   http://localhost:8080/agree?id=?
    router.POST("/disagree",Disagree)//踩  http://localhost:8080/disagree?id=?
    router.POST("/comment",Comment)//评论回复 http://localhost:8080/comment=?
    router.POST("/collect",Collect)//收藏 http://localhost:8080/collect?question=?
    router.POST("/selectfollow",Selectfollow)//查看我关注的人  http://localhost:8080/selectfollow?
    router.POST("/selectanser",Selectanser)//查看回答  http://localhost:8080/selectanser?
	router.POST("/selectessay",Selectessay)//查看文章  http://localhost:8080/selectessay?
	router.POST("/selectquestion",Selectquestion)//查看提问  http://localhost:8080/selectquestion?
	router.POST("/selectcollect",Selectcollect)// 查看收藏 http://localhost:8080/selectcollect?
    router.POST("/updateusername",Updateusername)// 修改昵称 http://localhost:8080/updateusername?newusername=?
    router.POST("/instruct",Instruct)//个人简介  http://localhost:8080/instruct?instruct=?
    router.POST("/delquestion",Delquestion)//删除提问 http://localhost:8080/delquestion?id=?
    router.POST("/delanser",Delanser)//删除回答  http://localhost:8080/delanser?id=?
    router.POST("/delessay",Delessay)//删除文章 http://localhost:8080/delessay?id=?


	router.Run(":8080")
}
/////跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method      //请求方法
		origin := c.Request.Header.Get("Origin")        //请求头部
		var headerKeys []string                             // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")        // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")      //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")      // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")        // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")       //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")       // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next()        //  处理请求
	}
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
	c.SetCookie("username", username, 1000, "/", "127.0.0.1", false, true)
	//第一个参数为 cookie 名；第二个参数为 cookie 值；第三个参数为 cookie 有效时长；第四个参数为 cookie 所在的目录；第五个为所在域，表示我们的 cookie 作用范围；第六个表示是否只能通过 https 访问；第七个表示 cookie 是否支持HttpOnly属性。
	if checkUserSignin(username, password) {
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
//改变cookie 存在时长就完事了

func Logout(c *gin.Context)  {
	username,err:=c.Cookie("username")
	if err != nil {
		log.Fatal(err)
	}
	c.SetCookie("username", username, -1, "/", "127.0.0.1", false, true)
	//第一个参数为 cookie 名；第二个参数为 cookie 值；第三个参数为 cookie 有效时长；第四个参数为 cookie 所在的目录；第五个为所在域，表示我们的 cookie 作用范围；第六个表示是否只能通过 https 访问；第七个表示 cookie 是否支持HttpOnly属性。

	c.JSON(200,gin.H{"Status":http.StatusOK,"message":"注销成功"})
}

//========================================主页部分========================================================
//推荐 essay里的前三条
func Recommend(c*gin.Context){
	rows, err := db.Query("SELECT essay FROM essay where id<=3")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
	var essay string
	err:=rows.Scan(&essay)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"essay":essay})
	}


}
// 搜索文章
////
func Search(c*gin.Context){
	_,err:=c.Cookie("username")
	if err!=nil {
		log.Fatal(err)
	}
	keyword:=c.PostForm("keyword")
	rows, err := db.Query("select essay from essay where essay like ?;","%"+keyword)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var essay string
		err:=rows.Scan(&essay)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"essay":essay})
	}
}
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
		"insert into question(username,question,collectnum)values(?,?,?)")//建一个名为question的表，把question放里面
	if err!=nil{
		log.Fatal(err)
		return false
	}
	stmt.Exec(username,question,0)
	return true
}
//====================================================================================================
//关注  我想的可能是和添加好友一样 我加一个取消关注
//这里要仔细注意 username是自己 usernames是想要关注的人
//
func Follow(c *gin.Context) {
	username,err:=c.Cookie("username")
	usernames:=c.PostForm("usernames")
	if err!=nil{
		log.Fatal(err)
	}

	stmt, err := DBConn().Prepare(
		`insert into friends(username,usernames,follownum)values(?,?,?)`)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(username,usernames, 1)
	c.JSON(200,gin.H{"message":"关注成功"})
}

//取消关注
func Notfollow(c*gin.Context){
	_,err:=c.Cookie("username")
	usernames:=c.PostForm("usernames")
	if err !=nil{
		log.Fatal(err)
	}
	stmt,err:=DBConn().Prepare(
		`update friends set follownum=0 where usernames=?`)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(usernames)
	c.JSON(200,gin.H{"message":"取消关注成功"})
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
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"回答成功"})
	}else{
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
func ansersql(username string,anser string)bool{
	stmt,err:=DBConn().Prepare(
		`insert into anser(username,anser,agreenum,disagreenum)values(?,?,?,?)`)
	if err!=nil{
		log.Fatal(err)
		return false
	}
		stmt.Exec(username,anser,0,0)
	return true
}

//==============================================================================================
//点赞&取消点赞  (第一次点赞，点两次取消点赞）
//1为点赞，0为取消点赞,在anser表里加value：agreenum
//这个不需要cookie把  不需要知道是谁点的赞
func Agree(c *gin.Context) {
	_,err:=c.Cookie("username")
	if err!=nil{
		log.Fatal(err)
	}
	id:=c.PostForm("id")
	rows, err :=db.Query("select agreenum from anser where id=?",id)
	if err != nil {
		fmt.Println("db.query is error agree",err)
	}
	for rows.Next() {
		var agreenum int
		err := rows.Scan(&agreenum)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		if agreenum==1{
			stmt,err:=DBConn().Prepare(
				`update anser set agreenum=0 where id=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(id)
			c.JSON(200,gin.H{"message":"点赞成功"})

		}else{
			stmt,err:=DBConn().Prepare(
				`update anser set agreenum=1 where id=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(id)
			c.JSON(200,gin.H{"message":"取消点赞成功"})

		}

	}
}

//===============================================================================================
//踩跟点赞类似直接ctrl+c
//踩&取消踩  (第一次踩，点两次取消踩）
////1为踩，0为取消踩,在anser表里加value：disagreenum
func Disagree(c *gin.Context) {
	_,err:=c.Cookie("username")
	id:=c.PostForm("id")
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	rows, err :=db.Query("select disagreenum from anser where id=?",id)
	if err != nil {
		fmt.Println("db.query is error agree",err)
	}
	for rows.Next() {
		var disagreenum int
		err := rows.Scan(&disagreenum)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		if disagreenum==1{
			stmt,err:=DBConn().Prepare(
				`update anser set disagreenum=0 where id=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(id)
			c.JSON(200,gin.H{"message":"取消踩人成功"})

		}else{
			stmt,err:=DBConn().Prepare(
				`update anser set disagreenum=1 where id=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(id)
			c.JSON(200,gin.H{"message":"踩人成功"})

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
//收藏 跟点赞一个道理
func Collect(c *gin.Context) {
	_,err:=c.Cookie("username")
	id:=c.PostForm("id")
	if err != nil {
		log.Fatal(err)
	}
	rows, err :=db.Query("select collectnum from question where id=?",id)
	if err != nil {
		fmt.Println("db.query is error collect",err)
	}
	for rows.Next() {
		var collectnum int
		err := rows.Scan(&collectnum)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		if collectnum==1{
			stmt,err:=DBConn().Prepare(
				`update question set collectnum=0 where id=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(id)
			c.JSON(200,gin.H{"message":"取消收藏成功"})

		}else{
			stmt,err:=DBConn().Prepare(
				`update question set collectnum=1 where id=?`)
			if err!=nil{
				fmt.Println(err)
			}
			stmt.Query(id)
			c.JSON(200,gin.H{"message":"收藏成功"})

		}
	}
}

//=======================================个人信息页===============================================
//查看关注了，关注者
func Selectfollow(c*gin.Context){
	username,err:=c.Cookie("username")
	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select usernames from friends where follownum=1 and username=?;",username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var usernames string
		err:=rows.Scan(&usernames)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"usernames":usernames})
	}
}
//查看自己的回答
func Selectanser(c*gin.Context){
	username,err:=c.Cookie("username")

	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select anser from anser where username=?;",username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var anser string
		err:=rows.Scan(&anser)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"anser":anser})
	}
}

//查看提问
func Selectquestion(c*gin.Context){
	username,err:=c.Cookie("username")

	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select question from question where username=?;",username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var question string
		err:=rows.Scan(&question)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"question":question})
	}
}

//查看文章
func Selectessay(c*gin.Context){
	username,err:=c.Cookie("username")

	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select essay from essay where username=?;",username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var essay string
		err:=rows.Scan(&essay)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"essay":essay})

	}

}

//查看收藏的问题  同关注 输入0是未收藏  输入1是收藏
func Selectcollect(c*gin.Context){
	username,err:=c.Cookie("username")
	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select question from question where collectnum=1 and username=?;",username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var question string
		err:=rows.Scan(&question)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"question":question})
	}
}

//修改昵称   username是想要更改的名字   newusername是新的名字
func Updateusername(c*gin.Context){
	username,err:=c.Cookie("username")
	if err != nil {
		log.Fatal(err)
	}

	newusername:=c.PostForm("newusername") //新的名字
	if updateusernamesql(newusername,username){
		c.JSON(200,gin.H{"Status":http.StatusOK,"message":"更改成功"})
	}else{
		c.JSON(500,gin.H{"Status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
func updateusernamesql(newusername string,username string)bool{
	stmt,err:=DBConn().Prepare(
		`update user set username=? where username=?`)
	if err != nil {
		log.Fatal(err)
		return false
	}
	stmt.Exec(newusername,username)
	return true
}
//个人简介
func Instruct(c*gin.Context){
	username,err:=c.Cookie("username")
	if err != nil {
		log.Fatal()
	}

	instruct:=c.PostForm("instruct")
	if instructsql(username,instruct){
		c.JSON(200,gin.H{"Status":http.StatusOK,"message":"上传成功"})
	}else{
		c.JSON(500,gin.H{"Status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
func instructsql(username string,instruct string)bool{
	stmt,err:=DBConn().Prepare(
		`insert into instruct(username,instruct)values(?,?)`)
	if err != nil {
		log.Fatal(err)
		return false
	}
	stmt.Exec(username,instruct)
	return true
}
//自己的提问回答文章=================================================================
//前面已经写了查看了  再写一点删除吧
//删除提问
func Delquestion(c*gin.Context){
	username,err:=c.Cookie("username")
	if err != nil {
		log.Fatal(err)
	}
	id:=c.PostForm("id")//想删除提问的ID
	if delquestionsql(id,username){
		c.JSON(200,gin.H{"Status":http.StatusOK,"message":"删除成功"})
	}else {
		c.JSON(500,gin.H{"Status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
func delquestionsql(id string,username string)bool{

	stmt,err:=DBConn().Prepare(
		`delete from question where id=? and username=?`)
	if err != nil {
		log.Fatal(err)
		return false
	}
	stmt.Exec(id,username)
	return true
}
//删除回答
func Delanser(c*gin.Context){
	username,err:=c.Cookie("username")
	if err != nil {
		log.Fatal(err)
	}
	id:=c.PostForm("id")
	if delansersql(id,username){
		c.JSON(200,gin.H{"Status":http.StatusOK,"message":"删除成功"})
	}else {
		c.JSON(500,gin.H{"Status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
func delansersql(id string,username string)bool{
	stmt,err:=DBConn().Prepare(
		`delete from anser where id=? and username=?`)
	if err != nil {
		log.Fatal(err)
		return false
	}
	stmt.Exec(id,username)
	return true
}
//删除文章
func Delessay(c*gin.Context){
	username,err:=c.Cookie("username")
	if err != nil {
		log.Fatal(err)
	}
	id:=c.PostForm("id")
	if delessaysql(id,username){
		c.JSON(200,gin.H{"Status":http.StatusOK,"message":"删除成功"})
	}else {
		c.JSON(500,gin.H{"Status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
func delessaysql(id string,username string)bool{
	stmt,err:=DBConn().Prepare(
		`delete from essay where id=? and username=?`)
	if err != nil {
		log.Fatal(err)
		return false
	}
	stmt.Exec(id,username)
	return true
}
//=====================================================================================================
//后面都是建立的结构体了
//又写了个comment  拿来放评论  可能不需要 只需要anser结构体？

type instruct struct {
	username string
	instruct string
}
type comment struct {
	username string
	comment string
	id int
}
//建一个anser表， 回答放着 点赞 踩  评论回复也放这
type anser struct{
	anser string
	agreenum int
	disagreenum int
	username string
	id int
}
//关注列表， friends表里面放的就是关注的人
type friends struct {
	follownum int
	username string
	id int
}
//这里建一个question表  提问全部放这
type question struct{
	username string
	question string
	collectnum int
	id int
}
//文章全放这
type essay struct{
	id int
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
