package repositories

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"snippetbox.lhsort.top/internal/models"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (u *UsersRepository) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user := &models.Users{
		Name:           name,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}
	return u.db.Create(&user).Error
}

func (u *UsersRepository) Authenticate(email, password string) (int, error) {
	var user models.Users
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return 0, err
	}
	return int(user.ID), nil
}
