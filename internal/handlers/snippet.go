package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Snippet struct {
}

func NewSnippet() *Snippet {
	s := &Snippet{}
	return s
}

func (s *Snippet) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "base", nil)
}

type ViewParams struct {
	Id int `form:"id"`
}

func (s *Snippet) View(c *gin.Context) {
	var p ViewParams
	if c.BindQuery(&p) == nil {
		c.JSON(http.StatusOK, gin.H{
			"id": p.Id,
		})
	}
}

func (s *Snippet) RouteRegister(r *gin.Engine) {
	r.GET("/", s.Home)
	st := r.Group("/snippet")
	st.GET("/view", s.View)
}
