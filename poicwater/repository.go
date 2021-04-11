package poicwater

import (
	"fmt"

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

func (r *PoicWaterRepository) Create(p *PoicWater) error {
	db := r.db.Create(p)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (r *PoicWaterRepository) Get() ([]PoicWater, error) {
	var p []PoicWater
	db := r.db.Find(&p)

	fmt.Println("結果")
	fmt.Println(p)

	if db.Error != nil {
		return nil, db.Error
	}
	return p, nil
}
