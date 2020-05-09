package service

import (
	"context"
	"fmt"

	"log"

	. "redrock7/lv0/init"

	"redrock7/lv0/proto"
)

type Server struct {}

func (s *Server)Register(ctx context.Context, in *proto.Register) (*proto.Reply, error)  {

	log.Printf("Received:%v\n",in.GetUsername(),in.GetPassword())
	username:=in.GetUsername()
	password:=in.GetPassword()
	message:=Reg(username,password)

	return &proto.Reply{Message: message},nil
}
func (s *Server) Login(ctx context.Context, in *proto.Login) (*proto.Reply, error){
	username:=in.GetUsername()
	password:=in.GetPassword()
	message:=Login(username,password)
	return &proto.Reply{Message: message},nil
}

func (s *Server) Update(ctx context.Context,in *proto.Update) (*proto.Reply, error){
	username:=in.GetUsername()
	newusername:=in.GetNewusername()
	message:=Update(newusername,username)
	return &proto.Reply{Message: message},nil
}
//这是 有关注册的数据库操作
func Reg(username string,password string)(message string){
	//注册sql语句
	stmt,err:= Db.Prepare(
		"insert into user(username,password)values(?,?)")
	if err!=nil{
		log.Fatal(err)
	}
	stmt.Exec(username,password)
	message=`注册成功`+username
	return
}
//有关登录的数据库操作
func Login(username string,password string)(message string){
	rows, err :=Db.Query("select id from user where username=? and password=?",username,password)

	if err != nil {
		fmt.Println("db.query is error login",err)
	}

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println("failed to rows.Scan",err)
		}
		if id>0{
			message=`登录成功`
		}else{
			message=`用户名或密码不正确`
		}
	}
	return
}
func Update(newusername string,username string)(message string){
	stmt,err:=Db.Prepare(
		`update user set username=? where username=?`)
	if err != nil {
		log.Fatal(err)

	}
	stmt.Exec(newusername,username)
	message=`更改名字成功`
	return
}