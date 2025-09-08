# Get History TP/SL Order

**Rate Limit:** 10 req/sec/uid

## Description

Get History TP/SL Order

## HTTP Request

**GET** `/api/v1/futures/tpsl/get_history_orders`

## Request Parameters

| Parameter    | Type   | Required | Description                                                               |
|--------------|--------|----------|---------------------------------------------------------------------------|
| symbol       | string | false    | Trading pair                                                              |
| side         | int32  | false    | Order side                                                                |
| positionMode | int32  | false    | Order position mode                                                       |
| startTime    | int64  | false    | Start timestamp Unix timestamp in milliseconds format, e.g. 1597026383085 |
| endTime      | int64  | false    | End timestamp Unix timestamp in milliseconds format, e.g. 1597026683085   |
| skip         | int64  | false    | Skip order count default: 0                                               |
| limit        | int64  | false    | Number of queries: Maximum: 100, default: 10                              |

## Request Example

```bash
curl -X 'GET'  --location 'https://fapi.bitunix.com/api/v1/futures/tpsl/get_history_orders?symbol=BTCUSDT' \
   -H "api-key:*******" \
   -H "sign:*" \
   -H "nonce:your-nonce" \
   -H "timestamp:1659076670000" \
   -H "language:en-US" \
   -H "Content-Type: application/json"
```

## Response Parameters

| Parameter      | Type   | Description                                                                           |
|----------------|--------|---------------------------------------------------------------------------------------|
| orderList      | list   | TP/SL order List                                                                      |
| > id           | string | Order id                                                                              |
| > positionId   | string | Position id                                                                           |
| > symbol       | string | Coin pair                                                                             |
| > base         | string | Base                                                                                  |
| > quote        | string | Quote                                                                                 |
| > tpPrice      | string | Take-profit trigger price                                                             |
| > tpStopType   | string | Take-profit trigger type LAST_PRICE MARK_PRICE                                        |
| > slPrice      | string | Stop-loss trigger price                                                               |
| > slStopType   | string | Stop-loss trigger type LAST_PRICE MARK_PRICE                                          |
| > tpOrderType  | string | Take-profit order type LIMIT MARKET Default is market.                                |
| > tpOrderPrice | string | Take-profit order price                                                               |
| > slOrderType  | string | Stop-loss order type LIMIT MARKET Default is market.                                  |
| > slOrderPrice | string | Stop-loss order price                                                                 |
| > tpQty        | string | Take-profit order quantity(base coin) At least one of `tpQty` or `slQty` is required. |
| > slQty        | string | Stop-loss order quantity(base coin) At least one of `tpQty` or `slQty` is required.   |
| > status       | string | TP/SL order status                                                                    |
| > ctime        | int64  | Create timestamp                                                                      |
| > triggerTime  | int64  | Trigger time timestamp                                                                |
| total          | int64  | Total                                                                                 |

## Response Example

```json
{
  "code": 0,
  "data": [
    {
      "positionId": "12345678",
      "symbol": "BTCUSDT",
      "qty": "0.5",
      "entryValue": "30000"
    }
  ]
}
```