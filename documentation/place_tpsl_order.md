# Place TP/SL Order

**Rate Limit:** 10 req/sec/UID

## Description
Place TP/SL Order

## HTTP Request
* POST `/api/v1/futures/tpsl/place_order`

## Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | true | Trading pair |
| positionId | string | true | Position ID associated with take-profit and stop-loss |
| tpPrice | string | false | Take-profit trigger price<br>At least one of `tpPrice` or `slPrice` is required. |
| tpStopType | string | false | Take-profit trigger type<br>LAST_PRICE<br>MARK_PRICE<br>Default is market price. |
| slPrice | string | false | Stop-loss trigger price<br>At least one of `tpPrice` or `slPrice` is required. |
| slStopType | string | false | Stop-loss trigger type<br>LAST_PRICE<br>MARK_PRICE<br>Default is market price. |
| tpOrderType | string | false | Take-profit order type<br>LIMIT<br>MARKET<br>Default is market. |
| tpOrderPrice | string | false | Take-profit order price |
| slOrderType | string | false | Stop-loss order type<br>LIMIT<br>MARKET<br>Default is market. |
| slOrderPrice | string | false | Stop-loss order price |
| tpQty | string | false | Take-profit order quantity(base coin)<br>At least one of `tpQty` or `slQty` is required. |
| slQty | string | false | Stop-loss order quantity(base coin)<br>At least one of `tpQty` or `slQty` is required. |

## Request Example

```bash
curl -X 'POST'  --location 'https://fapi.bitunix.com/api/v1/futures/tpsl/place_order' \
   -H "api-key:*******" \
   -H "sign:*" \
   -H "nonce:your-nonce" \
   -H "time:1659076670000" \
   -H "language:en-US" \
   -H "Content-Type: application/json" \
 --data '{"symbol":"BTCUSDT","positionId":"111","tpPrice":"12","tpStopType":"LAST_PRICE","slPrice":"9","slStopType":"LAST_PRICE","tpOrderType":"LIMIT","tpOrderPrice":"11","slOrderType":"LIMIT","slOrderPrice":"8","tpQty":"1","slQty":"1"}'
```

## Response Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| orderId | string | TP/SL Order ID |

## Response Example

```json
{"code":0,"data":{"orderId":"11111"},"msg":"Success"}
```