# Get History Orders

**Rate Limit:** 10 req/sec/uid

## Description
Get history orders, sort by create time desc

## HTTP Request
* GET `/api/v1/futures/trade/get_history_orders`

## Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | false | Trading pair |
| orderId | string | false | order id |
| clientId | string | false | client id |
| status | string | false | Order status<br>**FILLED**<br>**CANCELED**<br>**PART_FILLED_CANCELED**<br>**EXPIRED** |
| type | string | false | Order type<br>**LIMIT**<br>**MARKET**<br>default all |
| startTime | int64 | false | Start timestamp Unix timestamp in milliseconds format, e.g. 1597026383085 |
| endTime | int64 | false | Start timestamp Unix timestamp in milliseconds format, e.g. 1597026683085 |
| skip | int64 | false | skip order count default: 0 |
| limit | int64 | false | Number of queries: Maximum: 100, default: 10 |

## Request Example

```bash
curl -X 'GET'  --location 'https://fapi.bitunix.com/api/v1/futures/trade/get_history_orders?symbol=BTCUSDT' \
   -H "api-key:*******" \
   -H "sign:*" \
   -H "nonce:your-nonce" \
   -H "time:1659076670000" \
   -H "language:en-US" \
   -H "Content-Type: application/json"
```

## Response Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| orderList | list | order list |
| > orderId | string | order id |
| > symbol | string | Trading pair |
| > qty | string | Amount (base coin) |
| > tradeQty | string | Fill amount (base coin) |
| > positionMode | string | ONE_WAY or HEDGE |
| > marginMode | string | ISOLATION or CROSS |
| > leverage | int | leverage |
| > price | string | Price of the order. Required if the order type is **LIMIT** |
| > side | string | Order direction buy: **BUY** sell: **SELL** |
| > orderType | string | Order type<br>**LIMIT**: limit orders<br>**MARKET**: market orders |
| > effect | string | Order expiration date. Required if the orderType is limit<br>**IOC**: Immediate or cancel<br>**FOK**: Fill or kill<br>**GTC**: Good till canceled(default value)<br>**POST_ONLY**: POST only |
| > clientId | string | Customize order ID |
| > reduceOnly | boolean | Whether or not to just reduce the position |
| > status | string | **INIT**: prepare status<br>**NEW**: pending<br>**PART_FILLED**: partially filled<br>**CANCELED**: canceled<br>**FILLED**: All filled |
| > fee | string | fee |
| > realizedPNL | string | realized pnl |
| > tpPrice | string | take profit trigger price |
| > tpStopType | string | take profit trigger type<br>**MARK_PRICE**<br>**LAST_PRICE** |
| > tpOrderType | string | take profit trigger place order type<br>**LIMIT**<br>**MARKET** |
| > tpOrderPrice | string | take profit trigger place order price<br>**LIMIT**<br>**MARKET**<br>required if tpOrderType is **LIMIT** |
| > slPrice | string | stop loss trigger price |
| > slStopType | string | stop loss trigger type<br>**MARK_PRICE**<br>**LAST_PRICE** |
| > slOrderType | string | stop loss trigger place order type<br>**LIMIT**<br>**MARKET** |
| > slOrderPrice | string | stop loss trigger place order price<br>**LIMIT**<br>**MARKET**<br>required if slOrderType is **LIMIT** |
| > ctime | int64 | create timestamp |
| > mtime | int64 | latest modify timestamp |
| total | int64 | total count |

## Response Example

```json
{"code":0,"data":{"orderList":[{"orderId":"11111","qty":"1","tradeQty":"0.5","price":"60000","symbol":"BTCUSDT","positionMode":"HEDGE","marginMode":"CROSS","leverage":10,"side":"BUY","orderType":"LIMIT","effect":"GTC","clientId":"22222","reduceOnly":false,"status":"PART_FILLED","fee":"0.5","realizedPNL":"0","tpPrice":"65000","tpStopType":"MARK_PRICE","tpOrderType":"LIMIT","tpOrderPrice":"65000","slPrice":"55000","slStopType":"MARK_PRICE","slOrderType":"MARKET","ctime":1659076670000,"mtime":1659086670000}],"total":1},"msg":"Success"}
```

Note: The response example appears to be truncated in the original documentation. I've completed it with assumed values for clarity.