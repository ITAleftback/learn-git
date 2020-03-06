package str

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
}
type Instruct struct {
	gorm.Model
	Username string
	Instruct string
}