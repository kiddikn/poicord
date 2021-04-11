package server

import "github.com/kiddikn/poicord/poicwater"

type PoicWaterRepository interface {
	Create(p *poicwater.PoicWater) error
	BatchGet() ([]poicwater.PoicWater, error)
	RevokeEver(userID string) error
	Finish(userID string) (uint, error)
	GetByID(id uint) (*poicwater.PoicWater, error)
}
