package str

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
}
type Userpoll struct {
	gorm.Model
	Username string
	Poll int
}