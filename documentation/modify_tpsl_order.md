# Modify TP/SL Order

**Rate Limit:** 10 req/sec/UID

## Description
Modify TP/SL Order

## HTTP Request
* POST `/api/v1/futures/tpsl/modify_order`

## Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| orderId | string | true | TP/SL Order ID |
| tpPrice | string | false | Take-profit trigger price At least one of `tpPrice` or `slPrice` is required. |
| tpStopType | string | false | Take-profit trigger type LAST_PRICE MARK_PRICE Default is market price. |
| slPrice | string | false | Stop-loss trigger price At least one of `tpPrice` or `slPrice` is required. |
| slStopType | string | false | Stop-loss trigger type LAST_PRICE MARK_PRICE Default is market price. |
| tpOrderType | string | false | Take-profit order type LIMIT MARKET Default is market. |
| tpOrderPrice | string | false | Take-profit order price |
| slOrderType | string | false | Stop-loss order type LIMIT MARKET Default is market. |
| slOrderPrice | string | false | Stop-loss order price |
| tpQty | string | false | Take-profit order quantity(base coin) At least one of `tpQty` or `slQty` is required. |
| slQty | string | false | Stop-loss order quantity(base coin) At least one of `tpQty` or `slQty` is required. |

## Request Example

```bash
curl -X 'POST'  --location 'https://fapi.bitunix.com/api/v1/futures/tpsl/modify_order' \
   -H "api-key:*******" \
   -H "sign:*" \
   -H "nonce:your-nonce" \
   -H "time:1659076670000" \
   -H "language:en-US" \
   -H "Content-Type: application/json" \
 --data '{"orderId":"123","tpPrice":"12","tpStopType":"LAST_PRICE","slPrice":"9","slStopType":"LAST_PRICE","tpOrderType":"LIMIT","tpOrderPrice":"11","slOrderType":"LIMIT","slOrderPrice":"8","tpQty":"1","slQty":"1"}'
```

## Response Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| orderId | string | TP/SL Order ID |

## Response Example

```json
{"code":0,"data":{"orderId":"11111"},"msg":"Success"}
```