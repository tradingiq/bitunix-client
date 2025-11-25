package bitunix

import (
	"context"
	"net/url"
	"strconv"

	"github.com/tradingiq/bitunix-client/model"
)

func (c *apiClient) GetPendingOrder(ctx context.Context, params model.PendingOrderParams) (*model.PendingOrderResponse, error) {
	queryParams := url.Values{}

	if params.Symbol != "" {
		queryParams.Add("symbol", params.Symbol.String())
	}
	if params.OrderID != "" {
		queryParams.Add("orderId", params.OrderID)
	}
	if params.ClientID != "" {
		queryParams.Add("clientId", params.ClientID)
	}
	if params.Status != "" {
		queryParams.Add("status", params.Status.String())
	}
	if params.StartTime != nil {
		queryParams.Add("startTime", strconv.FormatInt(params.StartTime.UnixMilli(), 10))
	}
	if params.EndTime != nil {
		queryParams.Add("endTime", strconv.FormatInt(params.EndTime.UnixMilli(), 10))
	}
	if params.Skip > 0 {
		queryParams.Add("skip", strconv.FormatInt(params.Skip, 10))
	}
	if params.Limit > 0 {
		queryParams.Add("limit", strconv.FormatInt(params.Limit, 10))
	}

	endpoint := "/api/v1/futures/trade/get_pending_orders"
	responseBody, err := c.restClient.Get(ctx, endpoint, queryParams)
	if err != nil {
		return nil, err
	}

	response := &model.PendingOrderResponse{}
	if err := handleAPIResponse(responseBody, endpoint, response); err != nil {
		return nil, err
	}

	return response, nil
}