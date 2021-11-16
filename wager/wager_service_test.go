package wager_test

import (
	"fmt"
	"test-wage/domain"
	"test-wage/domain/mocks"
	"test-wage/wager"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//TestListWager_MockListWagerRepo_ShouldReturnSameData unit test, mock repo here
func TestListWager_MockListWagerRepo_ShouldReturnSameData(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mockRepo := mocks.NewMockWagerRepo(ctr)

	mockRepo.
		EXPECT().
		ListWagers(gomock.Any()).
		Return([]domain.Wager{
			{
				ID:                  1,
				TotalWagerValue:     1,
				Odds:                1,
				SellingPercentage:   2.0,
				SellingPrice:        2.0,
				CurrentSellingPrice: 2.0,
				PercentageSold:      2.0,
				AmountSold:          2.0,
				PlacedAt:            time.Now(),
			},
		}, nil)

	wagerService := wager.NewWagerService(mockRepo)

	pageInfo := domain.PageInfo{}
	wagers, err := wagerService.ListWagers(pageInfo)
	fmt.Println("-----------------", wagers)

	assert.Nil(t, err)
	assert.NotNil(t, wagers)
	assert.Equal(t, 1, len(wagers))

}
