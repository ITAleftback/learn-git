package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"redrock2/user"
)

func main() {
	//每隔24h 更新票数为3  因为24h太久 为了加快效果我设置为每隔一分钟就更新一次票数
	go func() {
		crontab := cron.New()
		crontab.AddFunc("0 */1 * * * ?", user.Update)
		crontab.Start()
	}()

	router := gin.Default()
	user.Setuprouter(router)
	router.Run(":8080")
}
