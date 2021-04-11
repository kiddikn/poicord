package poicwater

import (
	"database/sql"
	"time"
)

type PoicWater struct {
	ID         uint
	UserID     string
	StartedAt  time.Time
	FinishedAt sql.NullTime
	RevokedAt  sql.NullTime
	CreatedAt  time.Time
}

func NewPoicWater(userID string) *PoicWater {
	return &PoicWater{
		UserID:    userID,
		StartedAt: time.Now(),
	}
}

func NewPoicWaterWithFinished(userID string, start time.Time, finish sql.NullTime) *PoicWater {
	return &PoicWater{
		UserID:     userID,
		StartedAt:  start,
		FinishedAt: finish,
	}
}

func (PoicWater) TableName() string {
	return "poic_water"
}
