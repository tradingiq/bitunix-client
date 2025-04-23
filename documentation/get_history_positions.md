# Get History Positions

**Rate Limit:** 10 req/sec/uid

## Description

Get History Positions

## HTTP Request

* GET `/api/v1/futures/position/get_history_positions`

## Request Parameters

| Parameter  | Type   | Required | Description                                                                                     |
|------------|--------|----------|-------------------------------------------------------------------------------------------------|
| symbol     | string | false    | Trading pair                                                                                    |
| positionId | string | false    | position id                                                                                     |
| startTime  | int64  | false    | Start timestamp(position create time) Unix timestamp in milliseconds format, e.g. 1597026383085 |
| endTime    | int64  | false    | Start timestamp(position create time) Unix timestamp in milliseconds format, e.g. 1597026683085 |
| skip       | int64  | false    | skip order count default: 0                                                                     |
| limit      | int64  | false    | Number of queries: Maximum: 100, default: 10                                                    |

## Request Example

```bash
curl -X 'GET' --location 'https://fapi.bitunix.com/api/v1/futures/position/get_history_positions?symbol=BTCUSDT' \
   -H "api-key:*******" \
   -H "sign:*" \
   -H "nonce:your-nonce" \
   -H "time:1659076670000" \
   -H "language:en-US" \
   -H "Content-Type: application/json"
```

## Response Parameters

| Parameter      | Type   | Description                                                                                                                           |
|----------------|--------|---------------------------------------------------------------------------------------------------------------------------------------|
| positionList   | list   | position list                                                                                                                         |
| > positionId   | string | position id                                                                                                                           |
| > symbol       | string | Trading pair                                                                                                                          |
| > maxQty       | string | max position amount                                                                                                                   |
| > entryPrice   | string | average entry price                                                                                                                   |
| > closePrice   | string | average close price                                                                                                                   |
| > liqQty       | string | liquidate quantity                                                                                                                    |
| > side         | string | **LONG** **SHORT**                                                                                                                    |
| > marginMode   | string | **ISOLATION** **CROSS**                                                                                                               |
| > positionMode | string | **ONE_WAY** **HEDGE**                                                                                                                 |
| > leverage     | int32  | leverage                                                                                                                              |
| > fee          | string | Deducted transaction fees: transaction fees deducted during the position                                                              |
| > funding      | string | total funding fee during the position                                                                                                 |
| > realizedPNL  | string | Realized PnL(exclude funding fee and transaction fee)                                                                                 |
| > liqPrice     | string | Estimated liquidation price<br>If the value <= 0, it means the position is at low risk and there is no liquidation price at this time |
| > ctime        | int64  | create timestamp                                                                                                                      |
| > mtime        | int64  | latest modify timestamp                                                                                                               |
| total          | int64  | total count                                                                                                                           |

## Response Example

```json
{
  "code": 0,
  "data": {
    "positionList": [
      {
        "positionId": "12345678",
        "symbol": "BTCUSDT",
        "maxQty": "0.5",
        "entryPrice": "30000",
        "closePrice": "35000",
        "liqQty": "0",
        "side": "LONG",
        "marginMode": "CROSS",
        "positionMode": "ONE_WAY",
        "leverage": 10,
        "fee": "2.5",
        "funding": "0.75",
        "realizedPNL": "2500",
        "liqPrice": "27500",
        "ctime": 1659076670000,
        "mtime": 1659086670000
      }
    ],
    "total": 1
  },
  "msg": "Success"
}
```
