package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"snippetbox.lhsort.top/internal/models"
	"snippetbox.lhsort.top/internal/repositories"
	"snippetbox.lhsort.top/internal/routes"
	"strconv"
)

type Snippet struct {
	repo *repositories.SnippetsRepository
}

func NewSnippet(repo *repositories.SnippetsRepository) *Snippet {
	s := &Snippet{repo: repo}
	return s
}

type snippetsData struct {
	Snippets []*models.Snippets
	Snippet  *models.Snippets
	PageInfo pageInfo
}

type pageInfo struct {
	CurrentPage int
	TotalPage   int
	I           []int
	CPMI        []int
	J           []int
	CPAJ        []int
	TPM5        int
}

// Home is a function that handles the home route of the server.
//
// It takes a gin.Context parameter and does the following:
// - Retrieves the latest snippets from the repository.
// - Creates a snippetsData struct with the retrieved snippets.
// - Renders the "home" template with the snippetsData.
// - Returns an HTTP 200 status code along with the rendered HTML.
func (s *Snippet) Home(c *gin.Context) {
	session := sessions.Default(c)
	flash := session.Flashes("normalMsg")
	err := session.Save()
	if err != nil {
		_ = c.Error(err)
		return
	}
	userId := session.Get("authenticatedUserID")
	data := snippetsData{}
	if userId != nil {
		ms, err := s.repo.Latest(userId.(int), c.Request)
		totalPage := s.repo.GetLatestPageInfo(userId.(int))
		if err != nil {
			_ = c.Error(err)
			return
		}
		cup := 1
		cp := c.Query("page")
		if cp != "" {
			cup, err = strconv.Atoi(cp)
			if err != nil {
				_ = c.Error(err)
				return
			}
		}
		info := getPageInfo(cup, totalPage)
		data.Snippets = ms
		data.PageInfo = info
	}
	td := newTemplateData(c, data, "Home", flash)
	c.HTML(http.StatusOK, "home", td)
}

func getPageInfo(cup int, totalPage int) pageInfo {
	i := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	currentPageMinusI := []int{cup - 10, cup - 9, cup - 8, cup - 7, cup - 6, cup - 5, cup - 4, cup - 3, cup - 2, cup - 1}
	j := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	currentPageAddJ := []int{cup + 1, cup + 2, cup + 3, cup + 4, cup + 5, cup + 6, cup + 7, cup + 8, cup + 9, cup + 10}
	totalPageMinus10 := totalPage - 10
	info := pageInfo{TotalPage: totalPage, CurrentPage: cup, I: i, CPMI: currentPageMinusI, J: j, CPAJ: currentPageAddJ, TPM5: totalPageMinus10}
	return info
}

type ViewParams struct {
	Id int `uri:"id"`
}

func (s *Snippet) View(c *gin.Context) {
	var p ViewParams
	if c.ShouldBindUri(&p) == nil {
		session := sessions.Default(c)
		userId := session.Get("authenticatedUserID")
		sp, err := s.repo.Get(p.Id, userId.(int))
		if err != nil {
			_ = c.Error(err)
			return
		}
		if len(sp.Content) == 0 {
			_ = c.Error(gorm.ErrRecordNotFound)
			return
		}

		flash := session.Flashes("createKey")
		err = session.Save()
		if err != nil {
			_ = c.Error(err)
			return
		}
		td := newTemplateData(c, sp, fmt.Sprintf("Snippets #%d", p.Id), flash)
		c.HTML(http.StatusOK, "view", td)
	}
}

type snippetCreateParams struct {
	Title       string `form:"title" binding:"required,min=3,max=100"`
	Content     string `form:"content" binding:"required"`
	Expires     int    `form:"expires" binding:"required,oneof=1 7 365"`
	FieldErrors map[string]string
}

func (s *Snippet) CreatePost(c *gin.Context) {
	var p snippetCreateParams
	err := c.ShouldBind(&p)
	if err == nil {
		session := sessions.Default(c)
		userId := session.Get("authenticatedUserID")
		id, err := s.repo.Insert(p.Title, p.Content, p.Expires, userId.(int))
		if err != nil {
			_ = c.Error(err)
			return
		}

		session.AddFlash("Snippets successfully created!", "createKey")
		err = session.Save()
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/snippet/view/%d", id))
	} else {
		var verr validator.ValidationErrors
		ok := errors.As(err, &verr)
		if ok {
			errm := routes.FormatErr(verr)
			params := snippetCreateParams{
				Title:       p.Title,
				Content:     p.Content,
				Expires:     p.Expires,
				FieldErrors: errm,
			}
			td := newTemplateData(c, params, "Create a new snippet", nil)
			c.HTML(http.StatusUnprocessableEntity, "create", td)
		}
	}

}

func (s *Snippet) Create(c *gin.Context) {
	td := newTemplateData(c, snippetCreateParams{}, "Create a new snippet", nil)
	c.HTML(http.StatusOK, "create", td)
}

func (s *Snippet) RouteRegister(r *gin.Engine) {
	r.GET("/", s.Home)
	st := r.Group("/snippet")
	st.Use(requireAuthentication)
	st.GET("/view/:id", s.View)
	st.GET("/create", s.Create)
	st.POST("/create", s.CreatePost)
}
