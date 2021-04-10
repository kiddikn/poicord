package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port must be set")
	}

	lcs := os.Getenv("LINE_CHANNEL_SECRET")
	if lcs == "" {
		log.Fatal("line channel secret must be set")
	}

	lat := os.Getenv("LINE_ACCESS_TOKEN")
	if lat == "" {
		log.Fatal("line access token must be set")
	}

	server, err := NewServer(lcs, lat)
	if err != nil {
		log.Fatal("initialize new server is failed")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/health", func(c *gin.Context) {
		server.healthHandler(c)
	})
	router.POST("/v1/callback", func(c *gin.Context) {
		server.lineHandler(c)
	})

	log.Print("http://localhost:" + port)
	router.Run(":" + port)
}

type server struct {
	bot *linebot.Client
}

func NewServer(channelSecret, channelToken string) (*server, error) {
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}
	return &server{
		bot: bot,
	}, nil
}

func (s *server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "poicoder is running.",
	})
}

func (s *server) lineHandler(c *gin.Context) {
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

		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			replyMessage := message.Text
			if replyMessage == "ヘルプ" {
				replyMessage = fmt.Sprintf("まずは返信できるかなを確認するね")
			}
			if _, err = s.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
				log.Print(err)
			}
			// case *linebot.StickerMessage:
			// 	replyMessage := fmt.Sprintf(
			// 		"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
			// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
			// 		log.Print(err)
			// 	}
		}
	}
	c.JSON(http.StatusOK, nil)
}
