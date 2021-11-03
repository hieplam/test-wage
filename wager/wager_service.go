package wager

import "test-wage/domain"

type Service struct {
	repo domain.WagerRepo
}

func NewWagerService(r domain.WagerRepo) domain.WagerService {
	return &Service{
		repo: r,
	}
}

func (s *Service) ListWagers(pageInfo domain.PageInfo) ([]domain.Wager, error) {
	return s.repo.ListWagers(pageInfo)
}

func (s *Service) PlaceWager(req domain.PlaceWagerReq) (domain.Wager, error) {
	return s.repo.PlaceWager(req)
}

func (s *Service) BuyWager(wagerID uint, req domain.BuyWagerReq) (domain.BuyWager, error) {
	return s.repo.BuyWager(wagerID, req)
}
