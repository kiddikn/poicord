package poicwater

import (
	"time"

	"github.com/jinzhu/gorm"
)

type PoicWater struct {
	gorm.Model

	UserID     string
	StartedAt  time.Time
	FinishedAt time.Time
	RevokedAt  time.Time
}

func NewPoicWater(userID string) *PoicWater {
	return &PoicWater{
		UserID:    userID,
		StartedAt: time.Now(),
	}
}
