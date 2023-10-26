package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"snippetbox.lhsort.top/configs"
	"snippetbox.lhsort.top/internal/handlers"
	"snippetbox.lhsort.top/internal/repositories"
	"snippetbox.lhsort.top/internal/routes"
)

func main() {
	configs.NewConfig()
	fx.New(
		routes.Module,
		repositories.Module,
		handlers.Module,
		fx.Invoke(
			func(r *gin.Engine) {},
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
