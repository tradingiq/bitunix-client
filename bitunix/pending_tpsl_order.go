package bitunix

import (
	"context"
	"net/url"
	"strconv"

	"github.com/tradingiq/bitunix-client/model"
)

func (c *apiClient) GetPendingTPSLOrder(ctx context.Context, params model.PendingTPSLOrderParams) (*model.PendingTPSLOrderResponse, error) {
	queryParams := url.Values{}

	if params.Symbol != "" {
		queryParams.Add("symbol", params.Symbol.String())
	}
	if params.PositionID != "" {
		queryParams.Add("positionId", params.PositionID)
	}
	if params.Side != "" {
		queryParams.Add("side", string(params.Side))
	}
	if params.PositionMode != "" {
		queryParams.Add("positionMode", string(params.PositionMode))
	}
	if params.Skip > 0 {
		queryParams.Add("skip", strconv.FormatInt(params.Skip, 10))
	}
	if params.Limit > 0 {
		queryParams.Add("limit", strconv.FormatInt(params.Limit, 10))
	}

	endpoint := "/api/v1/futures/tpsl/get_pending_orders"
	responseBody, err := c.restClient.Get(ctx, endpoint, queryParams)
	if err != nil {
		return nil, err
	}

	response := &model.PendingTPSLOrderResponse{}
	if err := handleAPIResponse(responseBody, endpoint, response); err != nil {
		return nil, err
	}

	return response, nil
}
