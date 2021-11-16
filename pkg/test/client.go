package test

import (
	"fmt"
	"log"
	"os"
	"test-wage/wager"
	"testing"
	"time"

	"gorm.io/gorm/logger"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	Db     *gorm.DB
	Router *gin.Engine
}

const testDb = "wager_test"

func NewClient(t *testing.T) *Client {
	r := gin.Default()
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// setup test db
	dsn := fmt.Sprintf("root:12345@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", testDb)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal("err when create test db connection", err)
	}

	wagerRepo := wager.NewWagerRepo(db)
	wagerService := wager.NewWagerService(wagerRepo)
	wager.NewWagerHandler(r, wagerService)

	return &Client{
		Db:     db,
		Router: r,
	}
}
