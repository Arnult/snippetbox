package handlers

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/justinas/nosurf"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"snippetbox.lhsort.top/internal/repositories"
	"snippetbox.lhsort.top/internal/routes"
)

type User struct {
	logger *zap.Logger
	repo   *repositories.UsersRepository
}

func NewUser(repo *repositories.UsersRepository, logger *zap.Logger) *User {
	return &User{repo: repo, logger: logger}
}

func (u *User) SignUp(c *gin.Context) {
	td := newTemplateData(c, signupDataParam{}, "SignUp", nil)
	c.HTML(http.StatusOK, "signup", td)
}

type signupDataParam struct {
	Name        string `form:"name" binding:"required,min=3,max=100"`
	Email       string `form:"email" binding:"required,email"`
	Password    string `form:"password" binding:"required,min=8"`
	FieldErrors map[string]string
}

func (u *User) SignUpPost(c *gin.Context) {
	var p signupDataParam
	err := c.ShouldBind(&p)
	if err != nil {
		var verr validator.ValidationErrors
		ok := errors.As(err, &verr)
		if ok {
			errm := routes.FormatErr(verr)
			p.FieldErrors = errm
			td := newTemplateData(c, p, "SignUp", nil)
			c.HTML(http.StatusUnprocessableEntity, "signup", td)
		}
	}
	err = u.repo.Insert(p.Name, p.Email, p.Password)
	if err != nil {
		fe := make(map[string]string)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			fe["Email"] = "Email is already in use"
		}
		p.FieldErrors = fe
		td := newTemplateData(c, p, "SignUp", nil)
		c.HTML(http.StatusUnprocessableEntity, "signup", td)
	} else {
		c.Redirect(http.StatusSeeOther, "/user/login")
	}

}

type loginParams struct {
	Email          string `form:"email" binding:"required,email"`
	Password       string `form:"password" binding:"required,min=8"`
	FieldErrors    map[string]string
	NonFieldErrors []string
}

func (u *User) Login(c *gin.Context) {
	td := newTemplateData(c, loginParams{}, "Login", nil)
	c.HTML(http.StatusOK, "login", td)
}

func (u *User) LoginPost(c *gin.Context) {
	var p loginParams
	err := c.ShouldBind(&p)
	if err != nil {
		var verr validator.ValidationErrors
		ok := errors.As(err, &verr)
		if ok {
			errm := routes.FormatErr(verr)
			p.FieldErrors = errm
			td := newTemplateData(c, p, "Login", nil)
			c.HTML(http.StatusUnprocessableEntity, "login", td)
		}
	}
	id, err := u.repo.Authenticate(p.Email, p.Password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			u.logger.Warn("incorrect password", zap.String("email", p.Email))
		}
		p.NonFieldErrors = append(p.NonFieldErrors, "Email or password is incorrect")
		td := newTemplateData(c, p, "Login", nil)
		c.HTML(http.StatusUnprocessableEntity, "login", td)
	} else {
		session := sessions.Default(c)
		session.Clear()
		session.Set("authenticatedUserID", id)
		token := nosurf.Token(c.Request)
		session.Set("user_csrf_token", token)
		err := session.Save()
		if err != nil {
			p.NonFieldErrors = append(p.NonFieldErrors, "Login failed")
			td := newTemplateData(c, p, "Login", nil)
			c.HTML(http.StatusUnprocessableEntity, "login", td)
		} else {
			c.Redirect(http.StatusSeeOther, "/snippet/create")
		}
	}
}

func (u *User) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Delete("authenticatedUserID")
	session.Delete("user_csrf_token")
	session.AddFlash("Logout successful!", "normalMsg")
	err := session.Save()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func (u *User) RouteRegister(r *gin.Engine) {
	g := r.Group("/user")
	g.GET("/signup", u.SignUp)
	g.POST("/signup", u.SignUpPost)
	g.GET("/login", u.Login)
	g.POST("/login", u.LoginPost)
	g.POST("/logout", requireAuthentication, u.Logout)
}
