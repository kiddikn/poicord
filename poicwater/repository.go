package poicwater

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PoicWaterRepository struct {
	db *gorm.DB
}

func NewPoicWaterRepository(dsn string) (*PoicWaterRepository, error) {
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PoicWaterRepository{
		db: db,
	}, nil
}
