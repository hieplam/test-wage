package wager

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"test-wage/domain"
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewWagerRepo(db *gorm.DB) domain.WagerRepo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) ListWagers(pageInfo domain.PageInfo) ([]domain.Wager, error) {
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}
	if pageInfo.OrderBy == "" {
		pageInfo.OrderBy = "asc"
	}

	if pageInfo.SortBy == "" {
		pageInfo.SortBy = "id"
	}

	var models []domain.Wager

	offset := pageInfo.Page*pageInfo.Limit - pageInfo.Limit
	order := fmt.Sprintf("%s %s", strings.ToLower(pageInfo.SortBy), strings.ToLower(pageInfo.OrderBy))

	query := r.db.Table("wagers")
	query = query.Limit(pageInfo.Limit).Offset(offset).Order(order)

	err := query.Find(&models).Error
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (r *Repo) PlaceWager(req domain.PlaceWagerReq) (domain.Wager, error) {
	if req.SellingPercentage < 1 || req.SellingPercentage > 100 || req.Odds <= 0 || req.SellingPrice <= 0 {
		return domain.Wager{}, errors.New("invalid_params")
	}

	if math.Abs(req.SellingPrice*100-math.Round(req.SellingPrice*100)) > 1e-9 {
		return domain.Wager{}, errors.New("invalid_selling_price_format")
	}

	sellPrice := float64(req.TotalWagerValue) * (float64(req.SellingPercentage) / 100)
	if req.SellingPrice <= sellPrice {
		return domain.Wager{}, errors.New("invalid_selling_price")
	}

	var mod domain.Wager
	mod, err := ConvertToWagerModel(req)
	if err != nil {
		return domain.Wager{}, err
	}
	mod.PlacedAt = time.Now().UTC()
	mod.CurrentSellingPrice = mod.SellingPrice

	createResult := r.db.Create(&mod)
	if createResult.Error != nil {
		return domain.Wager{}, nil
	}

	return mod, nil

}
func ConvertToWagerModel(req domain.PlaceWagerReq) (domain.Wager, error) {
	bs, err := json.Marshal(req)
	if err != nil {
		return domain.Wager{}, err
	}

	var mod domain.Wager
	err = json.Unmarshal(bs, &mod)
	if err != nil {
		return domain.Wager{}, err
	}

	return mod, nil
}
func (r *Repo) BuyWager(wagerID uint, req domain.BuyWagerReq) (domain.BuyWager, error) {
	if req.BuyingPrice <= 0 {
		return domain.BuyWager{}, errors.New("invalid_buying_price__must_greater_zero")
	}
	var wagerInDb domain.Wager
	if err := r.db.First(&wagerInDb, wagerID).Error; err != nil {
		return domain.BuyWager{}, err
	}
	if req.BuyingPrice > wagerInDb.CurrentSellingPrice {
		return domain.BuyWager{}, errors.New("invalid_buying_price__must_less_than_selling_price")
	}

	//TODO put these action to transaction
	buyWagerEntity := domain.BuyWager{
		WagerID:     wagerID,
		BuyingPrice: req.BuyingPrice,
		BoughtAt:    time.Now().UTC(),
	}
	if err := r.db.Create(&buyWagerEntity).Error; err != nil {
		return domain.BuyWager{}, err
	}

	// update corresponding wager
	var updatedWager domain.Wager
	if err := r.db.First(&updatedWager, wagerID).Error; err != nil {
		return domain.BuyWager{}, err
	}
	updatedWager.CurrentSellingPrice -= req.BuyingPrice
	updatedWager.AmountSold += req.BuyingPrice
	updatedWager.PercentageSold = uint(math.Round(updatedWager.AmountSold / updatedWager.SellingPrice * 100))
	if err := r.db.Save(updatedWager).Error; err != nil {
		return domain.BuyWager{}, err
	}

	return buyWagerEntity, nil
}
