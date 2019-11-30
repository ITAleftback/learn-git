package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	//"strconv"
	"sync"
	//""
)
var db *sql.DB
var j,i int64
var WEB string
var lock sync.Mutex
func main()  {
	err:=initDB()
	if err!=nil{
		fmt.Printf("initDB failed,err:%v\n",err)
	}
	fmt.Println("连接数据库成功")
	file, err1 := os.Create("list.txt");
	if err1 != nil {
		fmt.Println(err1);
	}
	client := &http.Client{}

	for i = 0;i<=27;i++ {
		gogo()
		resp, err := client.Get(WEB)
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		a := regexp.MustCompile(`201921[\d]{4}[\p{Han}]+`)
		c := regexp.MustCompile(`[\p{Han}]+`)
		name := strings.Join(c.FindStringSubmatch(strings.Join(a.FindStringSubmatch(string(body)),``)),``)
		if err != nil {
			fmt.Println(err)
		}
		file.Write([]byte(name));
		file.Write([]byte(" "))
		insertDB(db,name)
	}
	file.Close()
}
func gogo()  {
	WEB = "http://jwzx.cqupt.edu.cn/kebiao/kb_stu.php?xh="
	lock.Lock()
	j = 2019213796 + i
	s :=strconv.FormatInt(j, 10)
	WEB += s
	lock.Unlock()
	return
}
func insertDB(db *sql.DB,n string)  {
	stmt, err := db.Prepare("insert into student(name) values (?)")
	if err != nil{
		log.Fatal(err)
	}
	stmt.Exec(n)

}
func initDB()(err error){
	//数据库信息
	dsn:="root:@tcp(127.0.0.1:3306)/student"
	//连接数据库
	db, err=sql.Open("mysql",dsn)
	if err != nil {
		return
	}
	err=db.Ping()
	if err!=nil{
		return
	}
	return
}
type student struct {
	name string
}