package wager

import (
	"errors"
	"math"
	"test-wage/domain"
)

type Service struct {
	repo domain.WagerRepo
}

func NewWagerService(r domain.WagerRepo) domain.WagerService {
	return &Service{
		repo: r,
	}
}

func (s *Service) ListWagers(pageInfo domain.PageInfo) ([]domain.Wager, error) {
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}
	if pageInfo.OrderBy == "" {
		pageInfo.OrderBy = "asc"
	}

	if pageInfo.SortBy == "" {
		pageInfo.SortBy = "id"
	}

	return s.repo.ListWagers(pageInfo)
}

func (s *Service) PlaceWager(req domain.PlaceWagerReq) (domain.Wager, error) {
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

	return s.repo.PlaceWager(req)
}

func (s *Service) BuyWager(wagerID uint, req domain.BuyWagerReq) (domain.BuyWager, error) {
	if req.BuyingPrice <= 0 {
		return domain.BuyWager{}, errors.New("invalid_buying_price__must_greater_zero")
	}

	return s.repo.BuyWager(wagerID, req)
}
