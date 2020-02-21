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


//最后总结一下我的思路：
//注册登录方面很简单课上教过没必要说，注销的话我是改变cooike存在时长为-1这样就能达到注销的效果
//推荐  因为没有找到随机出的方式  所以我干脆就直接推荐数据库里essay的前三条
//搜索  我百度到一个没学过的sql语句 like  能够模糊搜索出相关数据
//写文章   很简单 就把essay当value写进数据库
//展示问题   就查询展示呗
//提问 跟写文章一样 把question当value写进数据库
//关注  关注我的用法与点赞 踩之类的很像  改变follownum状态1为关注  0为没有关注
//取消关注 就跟关注一样的嘛
//查看问题相应的回答   我的思路是输入一个pid  pid指的是想回答问题的id  这样数据库代表你回答的对应的问题可以辨明  按理来说应该由前端保存id数据然后直接传给pid 的 但我菜鸡没办法
//回答   需要输入pid和anser  anser就是回答  pid是想回答的问题   而前面我写的展示问题自然会告诉id  用户输入pid  就是对应id 想回答的问题
//点赞  同关注
//踩   同点赞
//展示回答的评论  这个与上面查看问题相应的回答相似    我记得要用到以前教过的留言板套娃   我记得用了啥递归来着  但是我没听懂  所以用的输入pid2这种菜鸡想法
//评论回复  就跟写文章也是一样的
//收藏   需要输入你收藏东西的id
//查看我关注的人    将friends表里  查看follownum=1的人  同时检查表里的username是cookie登录状态中的username 然后筛选出来的usernames就是关注的人
//查看我的粉丝     就跟上面很像   但是呢  筛选的是username和follownum=1
//查看自己的回答   查看表里 对应cookie中登录状态中的username  筛选出自己的回答
//查看我的文章 同上
//查看我的提问 同上
//查看我的收藏  同查看我关注的人
//修改我的昵称   从cookie中读取现在的username是啥   去表里找对应的username  然后用sql语句将username改成newusername
//个人简介    跟写文章很像
//删除我的提问   就在上面查看自己的提问中会告诉每个提问的id  你像删除谁  就输入对应的id就可以了
//删除我的回答  同上
//删除我的文章  同上
/////
///
//

//
func main() {

	router:=gin.Default()
	router.Use(Cors())
	router.POST("/registe",Registe)//注册 http://localhost:8080/registe?username=?&password=?
	router.POST("/login",Login)//登录   http://localhost:8080/login?username=?&password=?
	router.POST("/logout",Logout)//注销      http://localhost:8080/logout?username=?
	router.POST("/recommend",Recommend)//推荐  http://localhost:8080/recommend?
	router.POST("/search",Search)//搜索  http://localhost:8080/search?keyword=?
	router.POST("/wrtessay",Wrtessay)//写文章 http://localhost:8080/wrtessay?essay=?
	router.POST("/showquestion",Showquestion)//展示问题  http://locallhost:8080/showquestion?
	router.POST("/askquestion",Askquestion)  //提问 http://localhost:8080/askquestion?question=?
	router.POST("/follow",Follow)//关注   http://localhost:8080/follow?usernames=?
    router.POST("/notfollow",Notfollow)//取消关注  http://localhost:8080/notfollow?usernames=?
	router.POST("/showanser",Showanser)//查看问题相应的回答  http://localhost:8080/showanser?pid=?
    router.POST("/anser",Anser)//回答   http://localhost:8080/anser?anser=?&pid=?   pid 指回答的问题的id是多少   方便对应
    router.POST("/agree",Agree)//点赞   http://localhost:8080/agree?id=?
    router.POST("/disagree",Disagree)//踩  http://localhost:8080/disagree?id=?
	router.POST("/showcomment",Showcomment)//展示回答的评论 与展示回答相似  http://localhost:8080/showcomment?pid2=?
	router.POST("/comment",Comment)//评论回复 http://localhost:8080/comment?comment=?
    router.POST("/collect",Collect)//收藏 http://localhost:8080/collect?id=?
    router.POST("/selectfollow",Selectfollow)//查看我关注的人  http://localhost:8080/selectfollow?
    router.POST("/selectfans",Selectfans)//查看我的粉丝  http://localhost:8080/selectfans?
    router.POST("/selectanser",Selectanser)//查看自己的回答  http://localhost:8080/selectanser?
	router.POST("/selectessay",Selectessay)//查看我的文章  http://localhost:8080/selectessay?
	router.POST("/selectquestion",Selectquestion)//查看我的提问  http://localhost:8080/selectquestion?
	router.POST("/selectcollect",Selectcollect)// 查看我的收藏 http://localhost:8080/selectcollect?
    router.POST("/updateusername",Updateusername)// 修改我的昵称 http://localhost:8080/updateusername?newusername=?
    router.POST("/instruct",Instruct)//个人简介  http://localhost:8080/instruct?instruct=?
    router.POST("/delquestion",Delquestion)//删除我的提问 http://localhost:8080/delquestion?id=?
    router.POST("/delanser",Delanser)//删除我的回答  http://localhost:8080/delanser?id=?
    router.POST("/delessay",Delessay)//删除我的文章 http://localhost:8080/delessay?id=?


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
	rows, err := db.Query("SELECT username,essay,id FROM essay where id<=3")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
	var essay string
	var username string
	var id int
	err:=rows.Scan(&username,&essay,&id)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"username":username,"essay":essay,"id":id})
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
	rows, err := db.Query("select id,essay,username from essay where essay like ?;","%"+keyword)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var essay string
		var username string
		var id int
		err:=rows.Scan(&username,&id,&essay)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"username":username,"id":id,"essay":essay})
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
/////展示问题的页面  返回数据库中的问题数据  如推荐
func Showquestion(c*gin.Context){
	rows, err := db.Query("SELECT username,id,question FROM question where id<=4")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var question string
		var username string
		var id int
		err:=rows.Scan(&username,&id,&question)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"username":username,"id":id,"question":question})
	}
}
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
//  一个问题的回答   这个接口放在问题的旁边 点进去就是这个问题的回答
func Showanser(c*gin.Context){
	_,err:=c.Cookie("username")
	if err != nil {
		log.Println(err)
	}
	pid:=c.PostForm("pid")///  这个pid  你看看想要查看问题的id  id是多少输入pid就是多少 返回的就是相应答案
	rows, err :=db.Query("select username,id,anser from anser where pid=?",pid)
	defer rows.Close()
	for rows.Next(){
		var anser string
		var username string
		var id int
		err:=rows.Scan(&username,&id,&anser)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"username":username,"id":id,"anser":anser})
	}


}

//写个回答，  跟提问文章那些差不多 放在anser表里
func Anser(c*gin.Context){
	anser:=c.PostForm("anser")
	username,err:=c.Cookie("username")
	pid:=c.PostForm("pid") //麻烦用户在回答的时候把这个问题的pid写上  这个问题id是多少就填pid多少  id的值问题前面会返回

	if err!=nil{
		log.Fatal(err)
	}
	if ansersql(username,anser,pid){
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"回答成功"})
	}else{
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
func ansersql(username string,anser string,pid string)bool{
	stmt,err:=DBConn().Prepare(
		`insert into anser(username,anser,agreenum,disagreenum,pid)values(?,?,?,?,?)`)
	if err!=nil{
		log.Fatal(err)
		return false
	}
		stmt.Exec(username,anser,0,0,pid)
	return true
}

//==============================================================================================
//点赞&取消点赞  (第一次点赞，点两次取消点赞）
//1为点赞，0为取消点赞,在anser表里加value：agreenum
//这个不需要cookie把  不需要知道是谁点的赞
//  返回的数据中会告诉id

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
//  评论是回复的回答  跟回答回复提问一样 也要用到pid  麻烦用户输入你评论回复的回答是哪个
//  一会应该还要写一个查看回答相应的评论回复接口

func Comment(c*gin.Context){
	comment:=c.PostForm("comment")
	pid2:=c.PostForm("pid2")
	username,err:=c.Cookie("username")//用之前的cookie
	if err!=nil{
		log.Fatal(err)
	}

	if commentsql(username,comment,pid2){
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"评论成功"})
	}else{
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"未知错误"})
	}
}
func commentsql(username string,comment string,pid2 string)bool{
	stmt,err:=DBConn().Prepare(
		`insert into comment(username,comment,pid2)values(?,?,?)`)//又建一个comment表.... 拿来放评论  可是我有点疑惑 是不是只需把评论放anser表就行了？
	if err!=nil{
		log.Fatal(err)
		return false
	}
	stmt.Exec(username,comment,pid2)
	return true
}
////  展示回答的评论来了
func Showcomment(c*gin.Context){
	_,err:=c.Cookie("username")
	if err != nil {
		log.Println(err)
	}
	pid2:=c.PostForm("pid2")///  这个pid  你看看想要查看问题的id  id是多少输入pid就是多少 返回的就是相应答案
	rows, err :=db.Query("select username,id,comment from comment where pid2=?",pid2)
	defer rows.Close()
	for rows.Next(){
		var comment string
		var username string
		var id int
		err:=rows.Scan(&username,&id,&comment)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"username":username,"id":id,"comment":comment})
	}


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
//查看我关注的人
func Selectfollow(c*gin.Context){
	username,err:=c.Cookie("username")
	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select id,usernames from friends where follownum=1 and username=?;",username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var usernames string
		var id int
		err:=rows.Scan(&id,&usernames)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"id":id,"usernames":usernames})
	}
}
//查看我的粉丝
func Selectfans(c*gin.Context){
	usernames,err:=c.Cookie("username")
	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select id,username from friends where follownum=1 and usernames=?;",usernames)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var username string
		var id int
		err:=rows.Scan(&id,&username)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"id":id,"username":username})
	}
}
//查看自己的回答
func Selectanser(c*gin.Context){
	username,err:=c.Cookie("username")

	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select username,id,anser from anser where username=?;",username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var anser string
		var username string
		var id int
		err:=rows.Scan(&username,&id,&anser)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"username":username,"id":id,"anser":anser})
	}
}

//查看提问
func Selectquestion(c*gin.Context){
	username,err:=c.Cookie("username")

	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select username,id,question from question where username=?;",username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var question string
		var username string
		var id int
		err:=rows.Scan(&username,&id,&question)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"username":username,"id":id,"question":question})
	}
}

//查看文章
func Selectessay(c*gin.Context){
	username,err:=c.Cookie("username")

	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select username,id,essay from essay where username=?;",username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var essay string
		var username string
		var id int
		err:=rows.Scan(&username,&id,&essay)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"username":username,"id":id,"essay":essay})

	}

}

//查看收藏的问题  同关注 输入0是未收藏  输入1是收藏
func Selectcollect(c*gin.Context){
	username,err:=c.Cookie("username")
	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query("select username,id,question from question where collectnum=1 and username=?;",username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var question string
		var username string
		var id int
		err:=rows.Scan(&username,&id,&question)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		c.JSON(200,gin.H{"username":username,"id":id,"question":question})
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
//删除提问   前面的查看  会告诉id
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
