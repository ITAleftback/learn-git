package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)
//连接数据库
func init() {
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/food?charset=utf8")
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
//查询(今天吃小于15块钱的美食）
func Query(n int){
	sqlStr:=`select name,price from food where price<?;`
	rows,err:=db.Query(sqlStr,n)
	if err!=nil{
		fmt.Printf("exec %s query failed,err:%v\n",sqlStr,err)
		return
	}
	defer rows.Close()
	for rows.Next(){
		var k1 food
		_=rows.Scan(&k1.name,&k1.price)

		fmt.Printf("k1:%v\n",k1)
	}

}
//传入今天想吃少于多少钱的美食
func main(){
	Query(15)
}
func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

type food struct {
	id   int
	name string
	price string
}

