package routes

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/justinas/nosurf"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"reflect"
	"snippetbox.lhsort.top/log"
	"strings"
	"time"
)

var Module = fx.Module("routes", fx.Provide(
	log.NewLog,
	NewHttpServer,
))

var trans ut.Translator

type HttpServerParam struct {
	fx.In

	L      *zap.Logger
	Routes []Handler `group:"routes"`
}

// NewHttpServer creates a new HTTP server using the provided HttpServerParam.
// The server is configured with gin as the underlying engine and returns a *gin.Engine.
// The server is started on a separate goroutine and listens on the specified port.
// It also registers the routes provided in the HttpServerParam.
func NewHttpServer(lc fx.Lifecycle, v HttpServerParam) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.DebugMode)

	// CreatePost a new gin engine with default middleware
	r := gin.Default()

	authKey, err := base64.StdEncoding.DecodeString(viper.GetString("redis.authKey"))
	if err != nil {
		panic(err)
	}
	encryptionKey, err := base64.StdEncoding.DecodeString(viper.GetString("redis.encryptionKey"))
	if err != nil {
		panic(err)
	}
	store, err := redis.NewStore(10,
		"tcp",
		fmt.Sprint(viper.GetString("redis.host"), ":", viper.GetString("redis.port")),
		viper.GetString("redis.password"),
		authKey, encryptionKey)
	if err != nil {
		panic(err)
	}

	r.Use(sessions.Sessions("snippet_session", store))

	err = initTrans(viper.GetString("i18n.lang"))
	if err != nil {
		return nil
	}

	// Set up template functions
	r.SetFuncMap(template.FuncMap{
		"currentYear": func() int {
			return time.Now().Year()
		},
	})

	// Set up secure headers
	r.Use(secureHeaders)
	r.Use(notFound)

	// Register HTML templates and static files
	r.LoadHTMLGlob("ui/html/**/*.gohtml")
	r.Static("/static", "ui/static")
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "page_404", gin.H{
			"Title": "Page not found",
		})
	})
	// Register routes
	for _, route := range v.Routes {
		route.RouteRegister(r)
	}

	// Append the startup hook to the lifecycle
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				// Start the server on a separate goroutine
				go func() {
					csrf := nosurf.New(r)
					csrf.SetBaseCookie(http.Cookie{
						HttpOnly: true,
						Path:     "/",
						Secure:   true,
					})
					csrf.SetFailureHandler(http.HandlerFunc(csrfFailHandler))
					err := endless.ListenAndServe(":"+viper.GetString("server.port"), csrf)
					if err != nil {
						v.L.Error("failed to start http server", zap.Error(err))
					}
					v.L.Info("starting http server", zap.String("address", ":4000"))
				}()
				return nil
			},
		},
	)

	// Return the gin engine
	return r
}

type Handler interface {
	// RouteRegister registers the provided Gin engine with the appropriate routes.
	//
	// r: A pointer to the Gin engine instance.
	//
	// No return value.
	RouteRegister(r *gin.Engine)
}

// AsRoute returns an annotated function with the provided input function and additional configurations.
//
// Parameters:
// - f: the input function to be annotated.
//
// Returns:
// - any: the annotated function.
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Handler)),
		fx.ResultTags(`group:"routes"`),
	)
}

func secureHeaders(c *gin.Context) {
	c.Header("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
	c.Header("Referrer-Policy", "origin-when-cross-origin")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "deny")
	c.Header("X-XSS-Protection", "0")
	c.Next()
}

func FormatErr(err validator.ValidationErrors) map[string]string {
	fields := err.Translate(trans)
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func initTrans(locale string) error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT)
		trans, _ = uni.GetTranslator(locale)

		switch locale {
		case "en":
			_ = entranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			_ = zhtranslations.RegisterDefaultTranslations(v, trans)
		default:
			_ = zhtranslations.RegisterDefaultTranslations(v, trans)
		}
	}
	return nil

}

func csrfFailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", nosurf.Reason(r))
}

func notFound(c *gin.Context) {
	c.Next()
	le := c.Errors.Last()
	if le == nil {
		return
	}
	c.HTML(http.StatusNotFound, "page_404", gin.H{
		"Title": "record not found",
	})
}
