package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"snippetbox.lhsort.top/internal/repositories"
)

type Snippet struct {
	repo *repositories.SnippetRepository
}

func NewSnippet(repo *repositories.SnippetRepository) *Snippet {
	s := &Snippet{repo: repo}
	return s
}

// Home is a function that handles the home route of the application.
//
// It takes a pointer to a gin.Context type parameter, which represents the HTTP request and response context.
// It does not return any value.
func (s *Snippet) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home", gin.H{"title": "Home"})
}

type ViewParams struct {
	Id int `form:"id"`
}

func (s *Snippet) View(c *gin.Context) {
	var p ViewParams
	if c.BindQuery(&p) == nil {
		sp, err := s.repo.Get(p.Id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(err)
				return
			}
		}

		c.HTML(http.StatusOK, "view", gin.H{
			"title":   fmt.Sprintf("Snippet #%d", p.Id),
			"snippet": sp,
		})
	}
}

type SnippetCreateParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Expires int    `json:"expires"`
}

func (s *Snippet) Create(c *gin.Context) {
	var p SnippetCreateParams
	if c.BindJSON(&p) == nil {
		id, err := s.repo.Insert(p.Title, p.Content, p.Expires)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/snippet/view?id=%d", id))
	}
}

func (s *Snippet) RouteRegister(r *gin.Engine) {
	r.GET("/", s.Home)
	st := r.Group("/snippet")
	st.GET("/view", s.View)
	st.POST("/create", s.Create)
}
