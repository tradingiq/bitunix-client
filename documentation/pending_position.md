# Get Pending Positions

> Rate Limit: 10 req/sec/uid

## Description
Get Pending Positions from the futures API.

## HTTP Request
```
GET /api/v1/futures/position/get_pending_positions
```

## Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | false | Trading pair |
| positionId | string | false | Position ID |

## Request Example

```bash
curl -X 'GET'  --location 'https://fapi.bitunix.com/api/v1/futures/position/get_pending_positions?symbol=BTCUSDT' \
   -H "api-key:*******" \
   -H "sign:*" \
   -H "nonce:your-nonce" \
   -H "timestamp:1659076670000" \
   -H "language:en-US" \
   -H "Content-Type: application/json"
```

## Response Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| positionId | string | Position ID |
| symbol | string | Trading pair |
| qty | string | Position amount |
| entryValue | string | Available amount for positions |
| side | string | **LONG** or **SHORT** |
| marginMode | string | **ISOLATION** or **CROSS** |
| positionMode | string | **ONE_WAY** or **HEDGE** |
| leverage | int32 | Leverage |
| fees | string | Deducted transaction fees: transaction fees deducted during the position |
| funding | string | Total funding fee during the position |
| realizedPNL | string | Realized PnL (exclude funding fee and transaction fee) |
| margin | string | Locked asset of the position |
| unrealizedPNL | string | Unrealized PnL |
| liqPrice | string | Estimated liquidation price. If the value <= 0, it means the position is at low risk and there is no liquidation price at this time |
| marginRate | string | Margin ratio |
| avgOpenPrice | string | Average open price |
| ctime | int64 | Create timestamp |
| mtime | int64 | Latest modify timestamp |

## Response Example

```json
{
  "code": 0,
  "data": [
    {
      "positionId": "12345678",
      "symbol": "BTCUSDT",
      "qty": "0.5",
      "entryValue": "30000",
      "side": "LONG",
      "positionMode": "HEDGE",
      "marginMode": "ISOLATION",
      "leverage": 100,
      "fees": "1.5",
      "funding": "0.25",
      "realizedPNL": "150.75",
      "margin": "300",
      "unrealizedPNL": "75.5",
      "liqPrice": "29100.5",
      "marginRate": "0.1",
      "avgOpenPrice": "60000",
      "ctime": 1659076670000,
      "mtime": 1659086670000
    }
  ]
}
```