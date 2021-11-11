package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"test-wage/wager"
	"time"

	"github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.user", "wager")
	viper.SetDefault("database.pass", "12345")
	viper.SetDefault("database.name", "wager")
	viper.SetDefault("server.address", "8080")

}

func main() {
	err := setupEnv()
	if err != nil {
		panic(err)
	}

	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	fmt.Println(dbHost, dbPort, dbUser, dbPass, dbName)
	conStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("charset", "utf8mb4")
	val.Add("parseTime", "True")
	val.Add("loc", "Local")
	dsn := fmt.Sprintf("%s?%s", conStr, val.Encode())

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
	log.Println("done init db")

	wagerRepo := wager.NewWagerRepo(db)
	wagerService := wager.NewWagerService(wagerRepo)
	wager.NewWagerHandler(r, wagerService)

	err = r.Run(fmt.Sprintf(":%s", viper.GetString("server.address")))
	if err != nil {
		log.Fatal("err when init server", err)
	}

}

// setupEnv Setup viper to read config from env first, otherwise read from config file
func setupEnv() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	err := viper.BindEnv("database.host", "DATABASE_HOST")
	if err != nil {
		return err
	}

	err = viper.BindEnv("database.port", "DATABASE_PORT")
	if err != nil {
		return err
	}

	err = viper.BindEnv("database.user", "DATABASE_USER")
	if err != nil {
		return err
	}

	err = viper.BindEnv("database.pass", "DATABASE_PASS")
	if err != nil {
		return err
	}

	err = viper.BindEnv("database.name", "DATABASE_NAME")
	if err != nil {
		return err
	}

	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}
