package routes

import (
	"context"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"snippetbox.lhsort.top/internal/snippet"
)

func NewHttpServer(lc fx.Lifecycle, log *zap.Logger) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				r.LoadHTMLGlob("ui/html/**/*.tmpl")
				r.Static("/static", "ui/static")
				log.Info("load html templates")
				go func() {
					err := endless.ListenAndServe(":4000", r)
					if err != nil {
						log.Error("failed to start http server", zap.Error(err))
					}
					log.Info("starting http server", zap.String("address", ":4000"))
				}()
				return nil
			},
		},
	)
	return r
}

func AddRoutes(r *gin.Engine) {
	r.GET("/", snippet.Home)
	st := r.Group("/snippet")
	st.GET("/view", snippet.View)
	st.POST("/create", func(c *gin.Context) {
		c.String(200, "Create a new s")
	})
}
