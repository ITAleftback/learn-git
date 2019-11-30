package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)
var db *sql.DB
func initDB()(err error){
	//数据库信息
	dsn:="root:@tcp(127.0.0.1:3306)/KAMERIDER"
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
//查询
func query(n int){
	sqlStr:=`select id,name,year from kamerider where id>?;`
	rows,err:=db.Query(sqlStr,n)
	if err!=nil{
		fmt.Printf("exec %s query failed,err:%v\n",sqlStr,err)
		return
	}
	defer rows.Close()
	for rows.Next(){
		var k1 kamerider
		err:=rows.Scan(&k1.id,&k1.name,&k1.year)
		if err!=nil{
			fmt.Printf("scan failed,err:%v\n",err)
		}
		fmt.Printf("k1:%v\n",k1)
	}
}
//更改
func updataDB(){
	stmt, err := db.Prepare("UPDATE kamerider SET name = '空我' WHERE id = 1")
	if err != nil{
		log.Fatal(err)
	}
	stmt.Exec();
}
//删除
func deleteDB(){
	stmt, err := db.Prepare("delete from kamerider where year = '2002'")
	if err != nil{
		log.Fatal(err);
	}
	stmt.Exec();
}


//插入
func insertDB(db *sql.DB)  {
	stmt, err := db.Prepare("insert into kamerider(name,year) values (?,?)")
	if err != nil{
		log.Fatal(err)
	}
	stmt.Exec("龙骑",2002)

}
func main(){
	err:=initDB()
	if err!=nil{
		fmt.Printf("initDB failed,err:%v\n",err)
	}
	fmt.Println("连接数据库成功")
	insertDB(db)
	//query(0)
	//updataDB()
	//deleteDB()
}
type kamerider struct {
	id   int
	name string
	year int
}
