package repositories

import (
	"go.uber.org/fx"
	"snippetbox.lhsort.top/internal/database"
)

var Module = fx.Module("repositories", fx.Provide(
	database.NewDatabase,
	NewSnippetRepository,
),
)
