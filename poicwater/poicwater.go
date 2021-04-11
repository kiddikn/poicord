package poicwater

import "time"

type PoicWater struct {
	ID          uint `gorm:"primary_key"`
	UserID      string
	StartedAt   time.Time
	FinishsedAt time.Time
	RevokedAt   time.Time
	CreatedAt   time.Time `gorm:"-"`
}

// type User struct {
// 	ID           uint
// 	Name         string
// 	Email        *string
// 	Age          uint8
// 	Birthday     *time.Time
// 	MemberNumber sql.NullString
// 	ActivatedAt  sql.NullTime
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time
//   }
