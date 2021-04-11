package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/kiddikn/poicord/poicwater"
	"github.com/kiddikn/poicord/server"
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

	dbu := os.Getenv("DATABASE_URL")
	if dbu == "" {
		log.Fatal("database url must be set")
	}

	r, err := poicwater.NewPoicWaterRepository("host=ec2-54-145-102-149.compute-1.amazonaws.com:5432 user=yzhligxpsnczcx dbname=d81hku2triopmn sslmode=disable password=11f7bb65e9c1ef1629acfefdc27745502aeca6494273eaa2f5527799550b8d2e")
	if err != nil {
		log.Fatal(err)
	}

	server, err := server.NewServer(lcs, lat, r)
	if err != nil {
		log.Fatal("initialize new server is failed")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/health", func(c *gin.Context) {
		server.HealthHandler(c)
	})
	router.POST("/v1/callback", func(c *gin.Context) {
		server.LineHandler(c)
	})

	log.Print("http://localhost:" + port)
	router.Run(":" + port)
}
