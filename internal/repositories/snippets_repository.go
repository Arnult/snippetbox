package repositories

import (
	"gorm.io/gorm"
	"net/http"
	"snippetbox.lhsort.top/internal/models"
	"time"
)

type SnippetsRepository struct {
	db *gorm.DB
}

func NewSnippetsRepository(db *gorm.DB) *SnippetsRepository {
	return &SnippetsRepository{
		db: db,
	}
}

func (s *SnippetsRepository) Insert(title string, content string, expires int, userId int) (int, error) {
	ex := time.Now().Add(time.Duration(expires) * 24 * 60 * time.Minute)
	snippet := models.Snippets{Title: title, Content: content, Expires: ex}
	users := models.Users{}
	users.ID = uint(userId)
	err := s.db.Model(&users).Association("Snippets").Append(&snippet)
	return int(snippet.ID), err
}

func (s *SnippetsRepository) Get(id int, userId int) (ms *models.Snippets, err error) {
	users := models.Users{}
	users.ID = uint(userId)
	ms = &models.Snippets{}
	ms.ID = uint(id)
	err = s.db.Model(&users).Association("Snippets").Find(&ms)
	return
}

func (s *SnippetsRepository) Latest(userId int, r *http.Request) (ms []*models.Snippets, err error) {
	users := models.Users{}
	users.ID = uint(userId)
	err = s.db.Scopes(paginate(r)).Model(&users).
		Select("id", "title", "created_at").
		Where("expires > ?", time.Now()).
		Order("created_at desc").
		Association("Snippets").Find(&ms)
	return
}

func (s *SnippetsRepository) GetLatestPageInfo(userId int) int {
	users := models.Users{}
	users.ID = uint(userId)
	total := s.db.Model(&users).
		Select("id", "title", "created_at").
		Where("expires > ?", time.Now()).
		Order("created_at desc").
		Association("Snippets").Count()
	totalPages := total / 10
	if total%10 != 0 {
		totalPages++
	}
	return int(totalPages)
}
