package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tradingiq/bitunix-client/model"
	"net/url"
)

func (c *apiClient) GetAccountBalance(ctx context.Context, params model.AccountBalanceParams) (*model.AccountBalanceResponse, error) {
	if params.MarginCoin == "" {
		return nil, fmt.Errorf("marginCoin is required")
	}

	queryParams := url.Values{}
	queryParams.Add("marginCoin", params.MarginCoin.String())

	responseBody, err := c.restClient.Get(ctx, "/api/v1/futures/account", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balance: %w", err)
	}

	response := &model.AccountBalanceResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account balance response: %w", err)
	}

	return response, nil
}
