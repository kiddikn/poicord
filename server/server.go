package server

import (
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
	poicEnd   = "終了"
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
	r := c.Request

	events, err := s.bot.ParseRequest(r)
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

	for _, event := range events {
		if event.Type != linebot.EventTypeMessage {
			continue
		}

		var msg string
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			msg = message.Text
		default:
			continue
		}
		// case *linebot.StickerMessage:
		// 	replyMessage := fmt.Sprintf(
		// 		"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
		// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		// 		log.Print(err)
		// 	}

		if msg == help {
			s.replyMessage(
				event.ReplyToken,
				fmt.Sprintf("以下のメッセージから始まる文章に反応します。\n・"+poicStart+"\n・"+poicEnd),
			)
		} else if strings.HasPrefix(msg, poicStart) {
			s.replyMessage(event.ReplyToken, "頑張って行こうね")
		} else if strings.HasPrefix(msg, poicEnd) {
			s.replyMessageWithSticker(event.ReplyToken, "お疲れ様ーー")
		}

	}
	c.JSON(http.StatusOK, nil)
}

func (s *Server) replyMessage(replyToken, reply string) {
	if _, err := s.bot.ReplyMessage(replyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		log.Print(err)
	}
}

func (s *Server) replyMessageWithSticker(replyToken, reply string) {
	if _, err := s.bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(reply),
		&linebot.StickerMessage{
			PackageID: "6136",
			StickerID: "10551378",
		},
	).Do(); err != nil {
		log.Print(err)
	}
}
