package main

import (
	"context"

	"fmt"


	"final/proto"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

func main()  {

	//注册路由
	http.HandleFunc("/Register", Reg)//注册
	http.HandleFunc("/Login", Login)//登录
	http.HandleFunc("/Ready",Ready)//准备  一开始默认为为准备状态
	http.HandleFunc("/Playchess",Playchess)
	http.HandleFunc("/Surrender",Surrender)

	//http.HandleFunc("/Surrender",Surrenderp)
	err := http.ListenAndServe("localhost:8001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//注册的路由
func Reg(w http.ResponseWriter, r *http.Request){
	username := r.FormValue("username")
	password := r.FormValue("password")
	RPCRegister(username,password)
}
func Login(w http.ResponseWriter, r *http.Request){
	username := r.FormValue("username")
	password := r.FormValue("password")
	RPCLogin(username,password)
}
func Ready(w http.ResponseWriter,r *http.Request){

	RPCReady()
}
func Playchess(w http.ResponseWriter,r *http.Request)  {
	var i,j int32
	//从键盘输入值,  i代表行  j代表列
	fmt.Println("请输入想下的位置")
	_, _ = fmt.Scanln(&i, &j)
	RPCplaychess(i,j)
}
func Surrender(w http.ResponseWriter,r *http.Request)  {
	RPCSurrender()
}
// 建立RPC 注册服务
func RPCRegister(username string,password string){
	// 建立连接到gRPC服务
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Register服务的客户端
	c :=proto.NewRegisterClient(conn)
	//注册
	if len(os.Args)>1{
		username = os.Args[1]
		password = os.Args[1]
	}
	//调用服务
	r,err:=c.Register(context.Background(),&proto.Register{Username:username,Password:password})
	if err!=nil{
		log.Fatalln(err)
	}
	log.Println(r.GetMessage())
}

// 建立RPC 登录服务
func RPCLogin(username string,password string){
	// 建立连接到gRPC服务
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Login服务的客户端
	c :=proto.NewLoginClient(conn)
	//注册
	if len(os.Args)>1{
		username = os.Args[1]
		password = os.Args[1]
	}
	//调用服务
	r,err:=c.Login(context.Background(),&proto.Login{Username:username,Password:password})
	if err!=nil{
		log.Fatalln(err)
	}
	log.Println(r.GetMessage())
}

func RPCReady(){
	// 建立连接到gRPC服务
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Ready服务的客户端
	c :=proto.NewReadyClient(conn)

	//调用服务
	r,err:=c.Ready(context.Background(),&proto.Ready{})
	if err!=nil{
		log.Fatalln(err)
	}
	log.Println(r.GetMessage())
}
func RPCplaychess(i int32,j int32){
	// 建立连接到gRPC服务
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Login服务的客户端
	c :=proto.NewPlaychessClient(conn)
	//调用服务
	r,err:=c.Playchess(context.Background(),&proto.Playchess{I:i,J:j})
	if err!=nil{
		log.Fatalln(err)
	}
	log.Println(r.GetMessage())
}
//认输
func RPCSurrender(){
	// 建立连接到gRPC服务
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Peace服务的客户端
	c :=proto.NewSurrenderClient(conn)
	//调用服务
	r,err:=c.Surrender(context.Background(),&proto.Surrender{})
	if err!=nil{
		log.Fatalln(err)
	}
	log.Println(r.GetMessage())
}
