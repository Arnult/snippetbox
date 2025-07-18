package database

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"snippetbox.lhsort.top/internal/models"
)

func NewDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", viper.GetString("database.user"), viper.GetString("database.password"), viper.GetString("database.host"), viper.GetString("database.port"), viper.GetString("database.name"))
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn, // DSN data source name
		DefaultStringSize: 256, // string 类型字段的默认长度
	}), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.Snippets{}, &models.Users{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
