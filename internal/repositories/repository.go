package repositories

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
	"net/http"
	"snippetbox.lhsort.top/internal/database"
	"strconv"
)

var Module = fx.Module("repositories", fx.Provide(
	database.NewDatabase,
	NewSnippetsRepository,
	NewUsersRepository,
),
)

func paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page <= 0 {
			page = 1
		}
		pageSize := 10

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
