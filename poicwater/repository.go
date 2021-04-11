package poicwater

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
)

var ErrRecordNotFound = errors.New("record not found")

type PoicWaterRepository struct {
	db *gorm.DB
}

func NewPoicWaterRepository(db *gorm.DB) *PoicWaterRepository {
	return &PoicWaterRepository{
		db: db,
	}
}

func (r *PoicWaterRepository) Create(p *PoicWater) error {
	db := r.db.Create(p)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (r *PoicWaterRepository) BatchGet() ([]PoicWater, error) {
	p := []PoicWater{}
	db := r.db.Find(&p)

	if db.Error != nil {
		return nil, db.Error
	}
	return p, nil
}

func (r *PoicWaterRepository) RevokeEver(userID string) error {
	// TODO:トランザクション

	p := []PoicWater{}
	db := r.db.Find(&p, "user_id=? and finished_at is null and revoked_at is null", userID)

	if db.Error != nil {
		return db.Error
	}

	var now sql.NullTime
	now.Scan(time.Now())

	for _, poic := range p {
		poic.RevokedAt = now
		// TODO:エラー処理
		r.db.Save(poic)
	}
	return nil
}

// Finish finishes last poicwater record, and return ID
func (r *PoicWaterRepository) Finish(userID string) (uint, error) {
	p := PoicWater{}
	db := r.db.First(&p, "user_id=? and finished_at is null and revoked_at is null", userID)

	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return 0, ErrRecordNotFound
		}
		return 0, db.Error
	}

	var now sql.NullTime
	now.Scan(time.Now())

	p.FinishedAt = now
	r.db.Save(p)

	return p.ID, nil
}

func (r *PoicWaterRepository) GetByID(id uint) (*PoicWater, error) {
	p := PoicWater{}
	db := r.db.First(&p, "id=?", id)

	if db.Error != nil {
		return nil, db.Error
	}
	return NewPoicWaterWithFinished(p.UserID, p.StartedAt, p.FinishedAt), nil
}
