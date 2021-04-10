package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	help      = "ヘルプ"
	poicStart = "本日のポイックウォーター開始"
	poicEnd   = "ポイックウォーター終了しました！がんばりました！"
)

type Server struct {
	bot *linebot.Client
}

func NewServer(channelSecret, channelToken string) (*Server, error) {
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}
	return &Server{
		bot: bot,
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

		if event.Type != linebot.EventTypeMessage {
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
	// case *linebot.StickerMessage:
	// 	replyMessage := fmt.Sprintf(
	// 		"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
	// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
	// 		log.Print(err)
	// 	}

	// メッセージによって処理を変更する
	if msg == help {
		if _, err := s.bot.ReplyMessage(
			e.ReplyToken,
			linebot.NewTextMessage(fmt.Sprintf("以下のメッセージから始まる文章に反応します。\n・"+poicStart)),
		).WithContext(ctx).Do(); err != nil {
			log.Print(err)
		}
	} else if strings.HasPrefix(msg, poicStart) {
		// 開始に対応するメッセージは打たなくても良いようにボタンテンプレートを返す
		t := linebot.NewButtonsTemplate(
			"",
			"ポイックウォーター終了",
			"ポイックウォーターが終わったら押してね",
			linebot.NewPostbackAction("終わったよ", "poicend", "", "ポイックウォーター終了しました！がんばりました！"),
		)
		if _, err := s.bot.ReplyMessage(
			e.ReplyToken,
			linebot.NewTextMessage("ええやん！！がんばれ！！"),
			linebot.NewTemplateMessage("ポイックウォーター終了", t)).WithContext(ctx).Do(); err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) postback(ctx context.Context, e *linebot.Event) {
	const (
		packageID = "6136" // 謝罪のプロ！LINEキャラクターズ
		stickerID = "10551378"
	)

	if _, err := s.bot.ReplyMessage(
		e.ReplyToken,
		linebot.NewTextMessage("お疲れ様ーー"),
		&linebot.StickerMessage{
			PackageID: packageID,
			StickerID: stickerID,
		},
	).WithContext(ctx).Do(); err != nil {
		log.Print(err)
	}
}
