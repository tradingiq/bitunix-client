# Get History Trades

**Rate Limit:** 10 req/sec/uid

## Description

Get history trades, sort by create time desc

## HTTP Request

* GET `/api/v1/futures/trade/get_history_trades`

## Request Parameters

| Parameter  | Notes                         | Type   | Required | Description                                                               |
|------------|-------------------------------|--------|----------|---------------------------------------------------------------------------|
| symbol     |                               | string | false    | Trading pair                                                              |
| orderId    |                               | string | false    | order id                                                                  |
| positionId | #broken: doesnt filter trades | string | false    | position id                                                               |
| startTime  |                               | int64  | false    | Start timestamp Unix timestamp in milliseconds format, e.g. 1597026383085 |
| endTime    |                               | int64  | false    | Start timestamp Unix timestamp in milliseconds format, e.g. 1597026683085 |
| skip       |                               | int64  | false    | skip order count default: 0                                               |
| limit      |                               | int64  | false    | Number of queries: Maximum: 100, default: 10                              |

## Request Example

```bash
curl -X 'GET' --location 'https://fapi.bitunix.com/api/v1/futures/trade/get_history_trades?symbol=BTCUSDT' \
   -H "api-key:*******" \
   -H "sign:*" \
   -H "nonce:your-nonce" \
   -H "time:1659076670000" \
   -H "language:en-US" \
   -H "Content-Type: application/json"
```

## Response Parameters

| Parameter      | Type    | Description                                                                                                                                                                                  |
|----------------|---------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| tradeList      | list    | trade list                                                                                                                                                                                   |
| > tradeId      | string  | trade id                                                                                                                                                                                     |
| > orderId      | string  | order id                                                                                                                                                                                     |
| > symbol       | string  | Trading pair                                                                                                                                                                                 |
| > qty          | string  | Amount (base coin)                                                                                                                                                                           |
| > positionMode | string  | ONE_WAY or HEDGE                                                                                                                                                                             |
| > marginMode   | string  | ISOLATION or CROSS                                                                                                                                                                           |
| > leverage     | int     | leverage                                                                                                                                                                                     |
| > price        | string  | Price of the order. Required if the order type is **LIMIT**                                                                                                                                  |
| > side         | string  | Order direction buy: **BUY** sell: **SELL**                                                                                                                                                  |
| > orderType    | string  | Order type<br>**LIMIT**: limit orders<br>**MARKET**: market orders                                                                                                                           |
| > effect       | string  | Order expiration date. Required if the orderType is limit<br>**IOC**: Immediate or cancel<br>**FOK**: Fill or kill<br>**GTC**: Good till canceled(default value)<br>**POST_ONLY**: POST only |
| > clientId     | string  | Customize order ID                                                                                                                                                                           |
| > reduceOnly   | boolean | Whether or not to just reduce the position                                                                                                                                                   |
| > fee          | string  | fee                                                                                                                                                                                          |
| > realizedPNL  | string  | realized pnl                                                                                                                                                                                 |
| > ctime        | int64   | create timestamp                                                                                                                                                                             |
| > roleType     |         | string                                                                                                                                                                                       | Trader tag<br>**TAKER**: taker<br>**MAKER**: maker                                                                                                                                           |
| total          |         | string                                                                                                                                                                                       | total count                                                                                                                                                                                  |

## Response Example

```json
{
  "code": 0,
  "data": {
    "tradeList": [
      {
        "tradeId": "12345",
        "orderId": "123",
        "marginCoin": null,
        "symbol": "BTCUSDT",
        "qty": "0.0005",
        "positionMode": "HEDGE",
        "marginMode": "ISOLATION",
        "leverage": 20,
        "price": "91077.2",
        "side": "BUY",
        "orderType": "MARKET",
        "effect": null,
        "clientId": null,
        "reduceOnly": true,
        "status": null,
        "fee": "0.019126212",
        "realizedPNL": "-0.15985",
        "ctime": "1745335541000",
        "roleType": "MAKER"
      }
    ],
    "total": 1
  },
  "msg": "Success"
}
```

Note: The response example appears to be truncated in the original documentation. I've completed it with a basic
structure for clarity.