package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
	"go.uber.org/fx"
	"net/http"
	"snippetbox.lhsort.top/internal/routes"
)

var Module = fx.Module("handlers", fx.Provide(
	routes.AsRoute(NewSnippet),
	routes.AsRoute(NewUser),
))

type templateData struct {
	Title           string
	Flash           []interface{}
	IsAuthenticated bool
	Data            any
	CsrfToken       string
	UserCsrfToken   any
}

func newTemplateData(c *gin.Context, data any, title string, flash []interface{}) templateData {
	session := sessions.Default(c)
	token := nosurf.Token(c.Request)
	id := session.Get("authenticatedUserID")
	userCsrfToken := session.Get("user_csrf_token")
	return templateData{
		Title:           title,
		Flash:           flash,
		IsAuthenticated: id != nil,
		Data:            data,
		CsrfToken:       token,
		UserCsrfToken:   userCsrfToken,
	}
}

func requireAuthentication(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("authenticatedUserID")
	if id == nil {
		c.Redirect(http.StatusSeeOther, "/user/login")
	}
	c.Header("Cache-Control", "no-store")
	c.Next()
}
