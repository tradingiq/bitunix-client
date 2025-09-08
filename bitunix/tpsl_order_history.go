package bitunix

import (
	"context"
	"net/url"
	"strconv"

	"github.com/tradingiq/bitunix-client/model"
)

func (c *apiClient) GetTPSLOrderHistory(ctx context.Context, params model.TPSLOrderHistoryParams) (*model.TPSLOrderHistoryResponse, error) {
	queryParams := url.Values{}

	if params.Symbol != "" {
		queryParams.Add("symbol", params.Symbol.String())
	}
	if params.Side != "" {
		queryParams.Add("side", string(params.Side))
	}
	if params.PositionMode != "" {
		queryParams.Add("positionMode", string(params.PositionMode))
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

	endpoint := "/api/v1/futures/tpsl/get_history_orders"
	responseBody, err := c.restClient.Get(ctx, endpoint, queryParams)
	if err != nil {
		return nil, err
	}

	response := &model.TPSLOrderHistoryResponse{}
	if err := handleAPIResponse(responseBody, endpoint, response); err != nil {
		return nil, err
	}

	return response, nil
}
