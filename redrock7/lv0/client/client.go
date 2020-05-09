package main

import (
	"context"

	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"redrock7/lv0/proto"
)

func main()  {

	//注册路由
	http.HandleFunc("/Register", Reg)
	http.HandleFunc("/Login", Login)
	http.HandleFunc("/Update",Update)
	err := http.ListenAndServe("localhost:8080", nil)
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
func Update(w http.ResponseWriter, r *http.Request){
	//  这里为了方便我就懒得用cookiet token session之类的了
	username := r.FormValue("username")
	newusername := r.FormValue("newusername")
	RPCUpdate(newusername,username)
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
// 建立RPC 更改服务
func RPCUpdate(newusername string,username string){
	//	// 建立连接到gRPC服务
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Update服务的客户端
	c :=proto.NewUpdateClient(conn)
	//注册
	if len(os.Args)>1{
		newusername = os.Args[1]
		username = os.Args[1]
	}
	//调用服务
	r,err:=c.Update(context.Background(),&proto.Update{Newusername:newusername,Username:username})
	if err!=nil{
		log.Fatalln(err)
	}
	log.Println(r.GetMessage())
}