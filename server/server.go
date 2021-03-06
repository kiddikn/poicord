package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hako/durafmt"
	"github.com/kiddikn/poicord/poicwater"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
)

const (
	help      = "ヘルプ"
	poicStart = "いでよ"
	poicEnd   = "歯磨き終了しました！がんばりました！"
)

type Server struct {
	bot   *linebot.Client
	r     PoicWaterRepository
	s     Sticker
	units durafmt.Units
}

func NewServer(channelSecret, channelToken string, r PoicWaterRepository, s Sticker) (*Server, error) {
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}

	units, err := durafmt.DefaultUnitsCoder.Decode("年:年,週:週,日:日,時間:時間,分:分,秒:秒,ミリ秒:ミリ秒,マイクロ秒:マイクロ秒")
	if err != nil {
		return nil, err
	}

	return &Server{
		bot:   bot,
		r:     r,
		s:     s,
		units: units,
	}, nil
}

func (s *Server) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "poicoder is running.",
	})
}

func (s *Server) LineHandler(c *gin.Context) {
	events, err := s.bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalidate signature",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
		}
		return
	}

	ctx := context.Background()

	for _, event := range events {
		if event.Type == linebot.EventTypePostback {
			s.postback(ctx, event)
			continue
		}

		if event.Type == linebot.EventTypeMessage {
			s.message(ctx, event)
			continue
		}
	}
	c.JSON(http.StatusOK, nil)
}

func (s *Server) message(ctx context.Context, e *linebot.Event) {
	var msg string
	switch message := e.Message.(type) {
	case *linebot.TextMessage:
		msg = message.Text
	default:
		return
	}

	// メッセージによって処理を変更する
	if msg == help {
		if _, err := s.bot.ReplyMessage(
			e.ReplyToken,
			linebot.NewTextMessage(fmt.Sprintf("以下のメッセージから始まる文章に反応します。\n・"+poicStart)),
		).WithContext(ctx).Do(); err != nil {
			log.Print(err)
		}
	} else if strings.HasPrefix(msg, poicStart) {
		// 既存データで終了していないデータはrevoke
		if err := s.r.RevokeEver(e.Source.UserID); err != nil {
			log.Print("レコードのrevoke大失敗")
			log.Print(err)
		}

		// 最新の開始のみデータを残す
		p := poicwater.NewPoicWater(e.Source.UserID)
		if err := s.r.Create(p); err != nil {
			log.Print("レコードの作成大失敗")
			log.Print(err)
		}

		// 開始に対応するメッセージは打たなくても良いようにボタンテンプレートを返す
		t := linebot.NewButtonsTemplate(
			"",
			"歯磨き終了",
			"歯磨きが終わったら押してね",
			linebot.NewPostbackAction("終わったよ", "poicend", "", poicEnd),
		)
		if _, err := s.bot.ReplyMessage(
			e.ReplyToken,
			linebot.NewTextMessage("いっぱい磨こうな^_^デンタルフロスは歯間を上下左右に丁寧にやるねんで！！"),
			linebot.NewTemplateMessage("歯磨き終了", t)).WithContext(ctx).Do(); err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) postback(ctx context.Context, e *linebot.Event) {
	// 終了処理を実施
	id, err := s.r.Finish(e.Source.UserID)
	if err != nil {
		if errors.Is(err, poicwater.ErrRecordNotFound) {
			log.Print("開始レコードがありません")
			if _, err := s.bot.ReplyMessage(
				e.ReplyToken,
				linebot.NewTextMessage("開始してますか？計算が開始できませんでした..."),
			).WithContext(ctx).Do(); err != nil {
				log.Print(err)
			}
			return
		}
		log.Print("レコードの終了処理の大失敗")
		log.Print(err)
		return
	}

	if id == 0 {
		log.Print("歯磨きの終了失敗")
		return
	}

	// 終了したポイックウォーターを取得する
	p, err := s.r.GetByID(id)
	if err != nil {
		log.Print("歯磨き開始の取得失敗:" + strconv.FormatUint(uint64(id), 10))
		log.Print(err)
		return
	}

	// 時間差を計算
	diff := p.FinishedAt.Time.Sub(p.StartedAt)
	duration := durafmt.Parse(diff).LimitFirstN(2).Format(s.units)

	// LINE通知
	st := s.s.GetRandomSticker()
	msg := "お疲れ様でした。\n所要時間は" + duration + "です。"
	if diff < 10*time.Minute {
		msg += "\nサボり？？これが続くと怒っちゃうよ"
	} else if diff < 15*time.Minute {
		msg += "\nもう少し頑張って汚れを落とそうぜ！"
	} else if diff >= 30*time.Minute {
		msg += "\nやば！！まじ天才かよ！！"
	} else if diff >= 25*time.Minute {
		msg += "\nExcellent！！とってもいけてる！！"
	} else if diff >= 20*time.Minute {
		msg += "\nよく頑張ってるね！いい子だいい子だ！"
	}

	if _, err := s.bot.ReplyMessage(
		e.ReplyToken,
		linebot.NewTextMessage(msg),
		&linebot.StickerMessage{
			PackageID: st.PackageID,
			StickerID: st.StickerID,
		},
	).WithContext(ctx).Do(); err != nil {
		log.Print(err)
	}
}

func (s *Server) GetHandler(c *gin.Context) {
	defer c.JSON(http.StatusOK, nil)

	p, err := s.r.BatchGet()
	if err != nil {
		log.Print("レコードの取得大失敗")
		return
	}
	log.Print(p)
}

func (s *Server) CreateHandler(c *gin.Context) {
	defer c.JSON(http.StatusOK, nil)

	p := poicwater.NewPoicWater("test")
	if err := s.r.Create(p); err != nil {
		log.Print("レコードの作成大失敗")
		log.Print(err)
		return
	}
}

func (s *Server) RevokeEverHandler(c *gin.Context) {
	defer c.JSON(http.StatusOK, nil)

	if err := s.r.RevokeEver("test"); err != nil {
		log.Print("レコードのrevoke大失敗")
		log.Print(err)
		return
	}
}

func (s *Server) FinishHandler(c *gin.Context) {
	defer c.JSON(http.StatusOK, nil)

	id, err := s.r.Finish("test")
	if err != nil {
		if errors.Is(err, poicwater.ErrRecordNotFound) {
			log.Print("開始レコードがありませんでした")
			return
		}
		log.Print("レコードの終了処理の大失敗")
		log.Print(err)
		return
	}

	if id == 0 {
		return
	}

	p, err := s.r.GetByID(id)
	if err != nil {
		log.Print("ポイックウォーターの取得失敗:" + strconv.FormatUint(uint64(id), 10))
		log.Print(err)
		return
	}

	diff := p.FinishedAt.Time.Sub(p.StartedAt)
	duration := durafmt.Parse(diff).String()
	log.Print("お疲れ様でした。\n所要時間は" + duration + "です。")
}
