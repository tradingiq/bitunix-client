package bitunix

import (
	"context"
	"github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"net/url"
)

func (c *apiClient) GetAccountBalance(ctx context.Context, params model.AccountBalanceParams) (*model.AccountBalanceResponse, error) {
	if params.MarginCoin == "" {
		return nil, errors.NewValidationError("marginCoin", "is required", nil)
	}

	queryParams := url.Values{}
	queryParams.Add("marginCoin", params.MarginCoin.String())

	endpoint := "/api/v1/futures/account"
	responseBody, err := c.restClient.Get(ctx, endpoint, queryParams)
	if err != nil {
		return nil, err
	}

	response := &model.AccountBalanceResponse{}
	if err := handleAPIResponse(responseBody, endpoint, response); err != nil {
		return nil, err
	}

	return response, nil
}
