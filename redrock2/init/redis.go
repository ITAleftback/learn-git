package init

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)
var Conn redis.Conn
// 链接redis数据库==========================================
func init(){
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	} else {
		fmt.Println("Connect to redis ok.")
	}
	Conn=conn
}
