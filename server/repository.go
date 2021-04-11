package server

import "github.com/kiddikn/poicord/poicwater"

type PoicWaterRepository interface {
	Create(p *poicwater.PoicWater) error
	Get() ([]*poicwater.PoicWater, error)
}
