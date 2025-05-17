package bitunix

import (
	"context"
	"net/url"

	"github.com/tradingiq/bitunix-client/model"
)

func (c *apiClient) GetPendingPositions(ctx context.Context, params model.PendingPositionParams) (*model.PendingPositionResponse, error) {
	queryParams := url.Values{}

	if params.Symbol != "" {
		queryParams.Add("symbol", params.Symbol.String())
	}

	if params.PositionID != "" {
		queryParams.Add("positionId", params.PositionID)
	}

	endpoint := "/api/v1/futures/position/get_pending_positions"
	responseBody, err := c.restClient.Get(ctx, endpoint, queryParams)
	if err != nil {
		return nil, err
	}

	response := &model.PendingPositionResponse{}
	if err := handleAPIResponse(responseBody, endpoint, response); err != nil {
		return nil, err
	}

	return response, nil
}
