package repositories

import (
	"gorm.io/gorm"
	"snippetbox.lhsort.top/internal/models"
	"time"
)

type SnippetRepository struct {
	db *gorm.DB
}

func NewSnippetRepository(db *gorm.DB) *SnippetRepository {
	return &SnippetRepository{
		db: db,
	}
}

func (s *SnippetRepository) Insert(title string, content string, expires int) (int, error) {
	ex := time.Now().Add(time.Duration(expires) * 24 * 60 * time.Minute)
	snippet := models.Snippet{Title: title, Content: content, Expires: ex}
	rs := s.db.Create(&snippet)
	return int(snippet.ID), rs.Error
}

func (s *SnippetRepository) Get(id int) (ms *models.Snippet, err error) {
	err = s.db.First(&ms, id).Error
	return
}

func (s *SnippetRepository) Latest() (ms []*models.Snippet) {
	s.db.Where("expires > ?", time.Now()).Order("created_at desc").Order("created_at desc").Find(&ms)
	return
}
