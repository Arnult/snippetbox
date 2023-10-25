package routes

import (
	"context"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"snippetbox.lhsort.top/log"
)

var Module = fx.Module("routes", fx.Provide(
	log.NewLog,
	NewHttpServer,
))

type HttpServerParam struct {
	fx.In

	L      *zap.Logger
	Routes []Handler `group:"routes"`
}

func NewHttpServer(lc fx.Lifecycle, v HttpServerParam) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	l := v.L
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				r.LoadHTMLGlob("ui/html/**/*.tmpl")
				r.Static("/static", "ui/static")
				for _, route := range v.Routes {
					route.RouteRegister(r)
				}
				l.Info("load html templates")
				go func() {
					err := endless.ListenAndServe(":"+viper.GetString("server.port"), r)
					if err != nil {
						l.Error("failed to start http server", zap.Error(err))
					}
					l.Info("starting http server", zap.String("address", ":4000"))
				}()
				return nil
			},
		},
	)
	return r
}

type Handler interface {
	RouteRegister(r *gin.Engine)
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Handler)),
		fx.ResultTags(`group:"routes"`),
	)
}
