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

	ctx := context.Background()

	for _, event := range events {
		if event.Type != linebot.EventTypeMessage && event.Type != linebot.EventTypePostback {
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

		// メッセージによって処理を変更する
		if msg == help {
			s.replyMessage(
				ctx,
				event.ReplyToken,
				fmt.Sprintf("以下のメッセージから始まる文章に反応します。\n・"+poicStart),
			)
		} else if strings.HasPrefix(msg, poicStart) {
			tem := linebot.NewButtonsTemplate(
				"",
				"ポイックウォータ終了",
				"ポイックウォーターが終わったら押してね",
				linebot.NewPostbackAction("終わったよ", "poicend", "", "ポイックウォーター終了しました！がんばりました！"),
			)
			if _, err := s.bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("ええやん！！がんばれ！！"),
				linebot.NewTemplateMessage("ポイックウォーター終了", tem)).WithContext(ctx).Do(); err != nil {
				log.Print(err)
			}
		} else if strings.HasPrefix(msg, poicEnd) {
			s.replyMessageWithSticker(ctx, event.ReplyToken, "お疲れ様ーー", "6136", "10551378")
		}
	}
	c.JSON(http.StatusOK, nil)
}

func (s *Server) replyMessage(ctx context.Context, replyToken, reply string) {
	if _, err := s.bot.ReplyMessage(replyToken, linebot.NewTextMessage(reply)).WithContext(ctx).Do(); err != nil {
		log.Print(err)
	}
}

func (s *Server) replyMessageWithSticker(ctx context.Context, replyToken, reply, packageID, stickerID string) {
	if _, err := s.bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(reply),
		&linebot.StickerMessage{
			PackageID: packageID,
			StickerID: stickerID,
		},
	).WithContext(ctx).Do(); err != nil {
		log.Print(err)
	}
}
