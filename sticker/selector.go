package sticker

import (
	"crypto/rand"
	"math/big"
)

type StickerSelector struct {
	stickers []*Sticker
}

type Sticker struct {
	PackageID string
	StickerID string
}

func NewStickerSelector() *StickerSelector {
	sticker := []*Sticker{
		{PackageID: "6136", StickerID: "10551394"}, // 謝罪のプロ！LINEキャラクターズ
		{PackageID: "6136", StickerID: "10551378"},
		{PackageID: "6325", StickerID: "10979914"},  //ちっちゃいブラコニ
		{PackageID: "11537", StickerID: "52002742"}, //動くブラウン＆コニー＆サリー スペシャル
		{PackageID: "11537", StickerID: "52002745"},
		{PackageID: "11537", StickerID: "52002752"},
		{PackageID: "11538", StickerID: "51626498"}, //動くチョコ＆LINEキャラ スペシャル
		{PackageID: "11538", StickerID: "51626521"},
		{PackageID: "11538", StickerID: "51626508"},
		{PackageID: "11539", StickerID: "52114115"}, //ユニバースター BT21 スペシャルVer.
		{PackageID: "1070", StickerID: "17878"},     //ムーン スペシャル
		{PackageID: "1070", StickerID: "17844"},
		{PackageID: "1070", StickerID: "17842"},
		{PackageID: "1070", StickerID: "17840"},
	}
	return &StickerSelector{stickers: sticker}
}

func (s *StickerSelector) GetRandomSticker() *Sticker {
	size := len(s.stickers)
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(size)))
	return s.stickers[n.Int64()]
}
