package handlers

import (
	"go.uber.org/fx"
	"snippetbox.lhsort.top/internal/routes"
)

var Module = fx.Module("handlers", fx.Provide(
	routes.AsRoute(NewSnippet),
))
