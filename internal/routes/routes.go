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

// NewHttpServer creates a new HTTP server using the provided HttpServerParam.
// The server is configured with gin as the underlying engine and returns a *gin.Engine.
// The server is started on a separate goroutine and listens on the specified port.
// It also registers the routes provided in the HttpServerParam.
func NewHttpServer(lc fx.Lifecycle, v HttpServerParam) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.DebugMode)

	// Create a new gin engine with default middleware
	r := gin.Default()

	// Register HTML templates and static files
	r.LoadHTMLGlob("ui/html/**/*.tmpl")
	r.Static("/static", "ui/static")

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
					err := endless.ListenAndServe(":"+viper.GetString("server.port"), r)
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
