package binance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"copytrader/internal/model"
)

type IdolPositionResponse struct {
	Positions       []model.Position `json:"otherPositionRetList"`
	UpdateTimestamp int64            `json:"updateTimeStamp"`
}

func GetIdolPositions(ctx context.Context, uid string) (*IdolPositionResponse, error) {
	req := map[string]string{
		"encryptedUid": uid,
		"tradeType":    "PERPETUAL",
	}
	resp, err := resty.New().R().SetBody(req).Post("https://www.binance.com/gateway-api/v1/public/future/leaderboard/getOtherPosition")
	if err != nil {
		return nil, err
	}
	var data struct {
		Data IdolPositionResponse `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}
	return &data.Data, nil
}

func printPositions(positions []model.Position) {
	fmt.Println(positions)
}
