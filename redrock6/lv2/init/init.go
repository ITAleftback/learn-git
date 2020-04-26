package init

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"redrock6/lv2/str"
)

var (
	DB *gorm.DB
)
//连接数据库==========================================
func init(){
	mysql, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/selectlesson?charset=utf8")
	mysql.SingularTable(true)
	if err != nil {
		panic(err)
	}

	DB=mysql
	DB.AutoMigrate(&str.Selectlesson{})
}
