package main

import (
	"log"
	"os"
	"test-wage/wager"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gin-gonic/gin"
)

//const dsn = "wager:12345@tcp(127.0.0.1:3306)/wager?charset=utf8mb4&parseTime=True&loc=Local"
//TODO  move to env, change localhost -> mysql viper, docker read from .env too
const dsn = "wager:12345@tcp(mysql:3306)/wager?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	r := gin.Default()
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal("err when create db connection", err)
	}

	wagerRepo := wager.NewWagerRepo(db)
	wagerService := wager.NewWagerService(wagerRepo)
	wager.NewWagerHandler(r, wagerService)

	err = r.Run()
	if err != nil {
		log.Fatal("err when init server", err)
	}
}
