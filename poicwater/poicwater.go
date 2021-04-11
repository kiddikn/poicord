package poicwater

import "time"

type PoicWater struct {
	ID         uint `gorm:"primary_key"`
	UserID     string
	StartedAt  time.Time
	FinishedAt time.Time
	RevokedAt  time.Time
	CreatedAt  time.Time `gorm:"-"`
}

func NewPoicWater(userID string) *PoicWater {
	return &PoicWater{
		UserID:    userID,
		StartedAt: time.Now(),
	}
}
