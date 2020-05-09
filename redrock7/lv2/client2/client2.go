package main

import (
	"context"
	"github.com/robfig/cron"

	"google.golang.org/grpc"
	"log"
	"os"
	"redrock7/lv0/proto"
)
/*
为了方便测试多客户端  所以我就把第二个客户端不写成路由的形式，直接自行取值
然后我就只测试注册
 */
func main()  {
	crontab := cron.New()
	task := func() {
		username:=`jjj`
		password:=`dada`
		RPCRegister(username,password)
	}
	// 添加定时任务, * * * * * 是 crontab,表示每分钟执行一次
	crontab.AddFunc("0 */1 * * * ?", task)
	// 启动定时器
	crontab.Start()
	// 定时任务是另起协程执行的,这里使用 select 简答阻塞.实际开发中需要
	// 根据实际情况进行控制
	select {}
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