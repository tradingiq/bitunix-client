package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tradingiq/bitunix-client/model"
	"net/url"
	"strconv"
)

func (c *apiClient) GetOrderHistory(ctx context.Context, params model.OrderHistoryParams) (*model.OrderHistoryResponse, error) {
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
	if params.Type != "" {
		queryParams.Add("type", params.Type.String())
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

	responseBody, err := c.restClient.Get(ctx, "/api/v1/futures/trade/get_history_orders", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}

	response := &model.OrderHistoryResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}
