
package resps

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 10000, "message": "ok!"})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": code, "message": msg})
}
