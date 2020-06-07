package service

import (
	"context"
	."final/chessboard"
	"final/judge"
	"fmt"
	"log"

	. "final/init"

	"final/proto"
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
func (s *Server)Ready(ctx context.Context, in *proto.Ready) (*proto.Reply, error)  {
	//准备好了就发一个棋盘
	for i := 0; i < 14; i++ {
		for j := 0; j < 14; j++ {
			fmt.Printf("%d\t", Board[i][j])

		}
		fmt.Println("\n")
	}
	message:=Rea()
	return &proto.Reply{Message: message},nil
}
func (s *Server)Playchess(ctx context.Context, in *proto.Playchess) (*proto.Reply, error)  {
	i:=in.GetI()
	j:=in.GetJ()
	message:=Play(i,j)
	return &proto.Reply{Message: message},nil
}
func (s *Server)Surrender(ctx context.Context, in *proto.Surrender) (*proto.Reply, error)  {

	message:=Sua()
	return &proto.Reply{Message: message},nil
}


////////////////////////////=============================
func Sua()(message string){
	if Flag%2==0 {
		message=`游戏结束，玩家一认输`
	}else {
		message=`游戏结束，玩家二认输`
	}
	return
}
func Play(i int32,j int32)(message string){
	if Board[i-1][j-1]!=0 {
		fmt.Println("这里已经有棋子了")
	}else {
		//玩家1执黑棋，第一个下
		if Flag%2 == 0 {
			Board[i-1][j-1] = 1
		} else {
			Board[i-1][j-1] = 2
		}

		for i := 0; i < 14; i++ {
			for j := 0; j < 14; j++ {
				fmt.Printf("%d\t", Board[i][j])
			}
			fmt.Println("\n")
		}
		//每次下完都要flag % 2 == 0进行判定是否输赢，若不分胜负则将次数++
		Flag++

		if judge.Judge(int(i-1), int(j-1)) {
			if Flag%2 == 0 {
				message=`玩家二胜利，游戏结束`
			} else {
				message=`玩家一胜利，游戏结束`
			}
		}
	}
	return
}
//有关准备的操作
func Rea()(message string){
	for i := 0; i < 14; i++ {
		for j := 0; j < 14; j++ {
			fmt.Printf("%d\t", Board[i][j])

		}
		fmt.Println("\n")
	}
	message=`已准备`
	return
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

