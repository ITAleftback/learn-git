package init

import (
	"github.com/jinzhu/gorm"
	"redrock1/str"
)
//连接数据库============================================================================
var (
	DB *gorm.DB
)
func init(){
	mysql, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/user?charset=utf8")
	mysql.SingularTable(true)
	if err != nil {
		panic(err)
	}

	DB=mysql
	DB.AutoMigrate(&str.User{})
	DB.AutoMigrate(&str.Instruct{})
}
