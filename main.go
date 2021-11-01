package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const dsn = "root:12345@tcp(127.0.0.1:3306)/wager?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	fmt.Println("init server")
	r := gin.Default()
	r.GET("/todo", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// list wagers
	r.GET("/wagers", ListWager)

	// place wager
	r.POST("/wagers", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//buy wager
	r.POST("buy/:wager_id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	err := r.Run()
	if err != nil {
		log.Fatal("err when init server", err)
	}
}

//	HANDLER LAYERS
func ListWager(c *gin.Context) {
	pageInfo := GetPagingInfo(c)

	if pageInfo.OrderBy != "asc" && pageInfo.OrderBy != "desc" {
		pageInfo.OrderBy = "asc"
	}

	wagers, err := ListWagerSrv(pageInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, wagers)
}
func ListWagerSrv(pageInfo PageInfo) ([]Wager, error) {
	return ListWagerRepo(pageInfo)
}

func ListWagerRepo(pageInfo PageInfo) ([]Wager, error) {
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}
	if pageInfo.OrderBy == "" {
		pageInfo.OrderBy = "asc"
	}

	if pageInfo.SortBy == "" {
		pageInfo.SortBy = "id"
	}
	//TODO move this to repo ctor
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	//end todo

	if err != nil {
		return []Wager{}, err
	}
	var models []Wager

	offset := pageInfo.Page*pageInfo.Limit - pageInfo.Limit
	order := fmt.Sprintf("%s %s", strings.ToLower(pageInfo.SortBy), strings.ToLower(pageInfo.OrderBy))

	query := db.Table("wagers")
	query = query.Limit(pageInfo.Limit).Offset(offset).Order(order)

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func GetPagingInfo(c *gin.Context) PageInfo {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Printf("err when getting page from context, set page to 1, err: %+v", err)
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		log.Printf("err when getting limit from context, set limit to 10, err: %+v", err)
		limit = 10
	}

	sortBy := strings.ToLower(c.DefaultQuery("sort_by", "id"))
	orderBy := strings.ToLower(c.DefaultQuery("order_by", "asc"))

	return PageInfo{Page: page, Limit: limit, SortBy: sortBy, OrderBy: orderBy}
}

type PageInfo struct {
	Page    int
	Limit   int
	SortBy  string
	OrderBy string
}

// Wager : wager information
type Wager struct {
	ID                  uint      `gorm:"primary_key;column:id" json:"id"`
	TotalWagerValue     uint      `gorm:"column:total_wager_value" json:"total_wager_value"`
	Odds                uint      `gorm:"column:odds" json:"odds"`
	SellingPercentage   uint      `gorm:"column:selling_percentage" json:"selling_percentage"`
	SellingPrice        float64   `gorm:"column:selling_price" json:"selling_price"`
	CurrentSellingPrice float64   `gorm:"column:current_selling_price" json:"current_selling_price"`
	PercentageSold      uint      `gorm:"default:null;column:percentage_sold" json:"percentage_sold"`
	AmountSold          float64   `gorm:"default:null;column:amount_sold" json:"amount_sold"`
	PlacedAt            time.Time `gorm:"column:placed_at" json:"placed_at"`
}

// BuyWager : record when wager is placed
type BuyWager struct {
	ID          uint      `gorm:"primary_key;column:id" json:"id"`
	WagerID     uint      `gorm:"index;column:wager_id" json:"wager_id"`
	BuyingPrice float64   `gorm:"column:buying_price" json:"buying_price"`
	BoughtAt    time.Time `gorm:"column:bought_at" json:"bought_at"`
}
