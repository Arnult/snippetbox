package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"snippetbox.lhsort.top/internal/routes"
)

func main() {
	fx.New(
		fx.Provide(zap.NewExample, routes.NewHttpServer),
		fx.Invoke(
			func(r *gin.Engine) {
				routes.AddRoutes(r)
			},
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
