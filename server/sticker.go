package server

import "github.com/kiddikn/poicord/sticker"

type Sticker interface {
	GetRandomSticker() *sticker.Sticker
}
