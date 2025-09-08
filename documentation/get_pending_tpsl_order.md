# Get Pending TP/SL Order

**Rate Limit:** 10 req/sec/uid

## Description

Get Pending TP/SL Order

## HTTP Request

`GET /api/v1/futures/tpsl/get_pending_orders`

## Request Parameters

| Parameter    | Type   | Required | Description                                  |
|--------------|--------|----------|----------------------------------------------|
| symbol       | string | false    | Trading pair                                 |
| positionId   | string | false    | Position ID                                  |
| side         | int32  | false    | Order side                                   |
| positionMode | int32  | false    | Order position mode                          |
| skip         | int64  | false    | Skip order count (default: 0)                |
| limit        | int64  | false    | Number of queries: Maximum: 100, default: 10 |

## Request Example

```bash
curl -X 'GET' --location 'https://fapi.bitunix.com/api/v1/futures/tpsl/get_pending_orders?symbol=BTCUSDT' \
   -H "api-key:*******" \
   -H "sign:*" \
   -H "nonce:your-nonce" \
   -H "timestamp:1659076670000" \
   -H "language:en-US" \
   -H "Content-Type: application/json"
```

## Response Parameters

| Parameter    | Type   | Description                                                 |
|--------------|--------|-------------------------------------------------------------|
| id           | string | Order ID                                                    |
| positionId   | string | Position ID                                                 |
| symbol       | string | Coin pair                                                   |
| base         | string | Base currency                                               |
| quote        | string | Quote currency                                              |
| tpPrice      | string | Take-profit trigger price                                   |
| tpStopType   | string | Take-profit trigger type: `LAST_PRICE`, `MARK_PRICE`        |
| slPrice      | string | Stop-loss trigger price                                     |
| slStopType   | string | Stop-loss trigger type: `LAST_PRICE`, `MARK_PRICE`          |
| tpOrderType  | string | Take-profit order type: `LIMIT`, `MARKET` (default: market) |
| tpOrderPrice | string | Take-profit order price                                     |
| slOrderType  | string | Stop-loss order type: `LIMIT`, `MARKET` (default: market)   |
| slOrderPrice | string | Stop-loss order price                                       |
| tpQty        | string | Take-profit order quantity (base coin)*                     |
| slQty        | string | Stop-loss order quantity (base coin)*                       |

*At least one of `tpQty` or `slQty` is required.

## Response Example

```json
{
  "code": 0,
  "data": [
    {
      "id": "123",
      "positionId": "12345678",
      "symbol": "BTCUSDT",
      "base": "BTC",
      "quote": "USDT",
      "tpPrice": "50000",
      "tpStopType": "LAST_PRICE",
      "slPrice": "70000",
      "slStopType": "LAST_PRICE",
      "tpOrderType": "LIMIT",
      "tpOrderPrice": "50000",
      "slOrderType": "LIMIT",
      "slOrderPrice": "70000",
      "tpQty": "0.01",
      "slQty": "0.01"
    }
  ],
  "msg": "Success"
}
```