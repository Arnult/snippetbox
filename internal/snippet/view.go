package snippet

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ViewParams struct {
	Id int `form:"id"`
}

func View(c *gin.Context) {
	var p ViewParams
	if c.BindQuery(&p) == nil {
		c.JSON(http.StatusOK, gin.H{
			"id": p.Id,
		})
	}
}
