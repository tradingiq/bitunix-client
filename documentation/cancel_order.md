# Cancel Orders

**Rate Limit:** 5 req/sec/uid

## Description

Cancel orders. Successful interface response is not necessarily equal to the success of the operation, please use the
websocket push message as an accurate judgment of the success of the operation.

## HTTP Request

* POST `/api/v1/futures/trade/cancel_orders`

## Request Parameters

| Parameter | Type   | Required | Description                                                                                          |
|-----------|--------|----------|------------------------------------------------------------------------------------------------------|
| symbol    | string | true     | Trading pair                                                                                         |
| orderList | list   | true     | order parameter list                                                                                 |
| orderId   | string | false    | Order ID<br>Either orderId or clientId is required. If both are entered, orderId prevails.           |
| clientId  | string | false    | Customize order ID<br>Either orderId or clientId is required. If both are entered, orderId prevails. |

## Request Example

```bash
curl -X 'POST'  --location 'https://fapi.bitunix.com/api/v1/futures/trade/cancel_orders' \
   -H "api-key:*******" \
   -H "sign:*" \
   -H "nonce:your-nonce" \
   -H "time:1659076670000" \
   -H "language:en-US" \
   -H "Content-Type: application/json" \
 --data '{"symbol":"BTCUSDT","orderList":[{"orderId":"11111"},{"clientId":"22223"}]}'
```

## Response Parameters

| Parameter   | Type   | Description           |
|-------------|--------|-----------------------|
| successList | list   | Successful order list |
| > orderId   | string | order id              |
| > clientId  | string | client id             |
| failureList | list   | Failed order list     |
| > orderId   | string | order id              |
| > clientId  | string | client id             |
| > errorMsg  | string | error msg             |
| > errorCode | string | error code            |

## Response Example

```json
{
  "code": 0,
  "data": {
    "successList": [
      {
        "orderId": "11111",
        "clientId": "22222"
      }
    ],
    "failureList": [
      {
        "orderId": "33333",
        "clientId": "44444",
        "errorMsg": "Order not found",
        "errorCode": "10001"
      }
    ]
  },
  "msg": "Success"
}
```

Note: The response example appears to be truncated in the original documentation. I've completed it with assumed values
for clarity.