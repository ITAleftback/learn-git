package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"redrock2/jwt"
	"redrock2/resps"
)
func User(c *gin.Context) {
	auth:= c.GetHeader("Authorization")
	fmt.Println(auth)
	if len(auth)<7 {
		resps.Error(c, 10011, "token error")
		c.Abort()
		return
	}
	token := auth[7:]
	uid, username, err := jwt.CheckToken(token)
	fmt.Println(err)
	if err != nil {
		resps.Error(c, 10011, "token error")
		c.Abort()
		return
	}
	c.Set("uid", uid)
	c.Set("username", username)
	c.Next()
	return
}