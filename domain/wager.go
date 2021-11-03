package domain

import (
	"time"
)

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

// WagerService represent the wager's use cases
type WagerService interface {
	ListWagers(pageInfo PageInfo) ([]Wager, error)
	PlaceWager(req PlaceWagerReq) (Wager, error)
	BuyWager(wagerID uint, req BuyWagerReq) (BuyWager, error)
}

// WagerRepo represent the wager's use cases
type WagerRepo interface {
	ListWagers(pageInfo PageInfo) ([]Wager, error)
	PlaceWager(req PlaceWagerReq) (Wager, error)
	BuyWager(wagerID uint, req BuyWagerReq) (BuyWager, error)
}

type PlaceWagerReq struct {
	TotalWagerValue   uint    `json:"total_wager_value"`
	Odds              uint    `json:"odds"`
	SellingPercentage uint    `json:"selling_percentage"`
	SellingPrice      float64 `json:"selling_price"`
}

type BuyWagerReq struct {
	BuyingPrice float64 `json:"buying_price"`
}

type PageInfo struct {
	Page    int
	Limit   int
	SortBy  string
	OrderBy string
}
