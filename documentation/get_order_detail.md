# Get Order Detail API

## Overview

**Rate Limit:** 10 req/sec/uid
**Description:** Retrieve detailed information about a specific order

## HTTP Request

```
GET /api/v1/futures/trade/get_order_detail
```

## Request Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| orderId   | string | false    | Order ID    |
| clientId  | string | false    | Client ID   |

**Note:** At least one of `orderId` or `clientId` is required.

## Request Example

```bash
curl -X 'GET' \
  --location 'https://fapi.bitunix.com/api/v1/futures/trade/get_order_detail?orderId=12345' \
  -H "api-key:*******" \
  -H "sign:*" \
  -H "nonce:your-nonce" \
  -H "timestamp:1659076670000" \
  -H "language:en-US" \
  -H "Content-Type: application/json"
```

## Response Parameters

| Parameter    | Type    | Description                                                                                                                                                                                   |
|--------------|---------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| orderId      | string  | Order ID                                                                                                                                                                                      |
| symbol       | string  | Trading pair                                                                                                                                                                                  |
| qty          | string  | Amount (base coin)                                                                                                                                                                            |
| tradeQty     | string  | Fill amount (base coin)                                                                                                                                                                       |
| positionMode | string  | ONE_WAY or HEDGE                                                                                                                                                                              |
| marginMode   | string  | ISOLATION or CROSS                                                                                                                                                                            |
| leverage     | int     | Leverage                                                                                                                                                                                      |
| price        | string  | Price of the order. Required if the order type is **LIMIT**                                                                                                                                   |
| side         | string  | Order direction buy: **BUY** sell: **SELL**                                                                                                                                                   |
| orderType    | string  | Order type **LIMIT**: limit orders **MARKET**: market orders                                                                                                                                  |
| effect       | string  | Order expiration date. Required if the orderType is limit<br>**IOC**: Immediate or cancel<br>**FOK**: Fill or kill<br>**GTC**: Good till canceled (default value)<br>**POST_ONLY**: POST only |
| clientId     | string  | Customize order ID                                                                                                                                                                            |
| reduceOnly   | boolean | Whether or not to just reduce the position                                                                                                                                                    |
| status       | string  | **INIT**: prepare status<br>**NEW**: pending<br>**PART_FILLED**: partially filled<br>**CANCELED**: canceled<br>**FILLED**: All filled                                                         |
| fee          | string  | Fee                                                                                                                                                                                           |
| realizedPNL  | string  | Realized PNL                                                                                                                                                                                  |
| tpPrice      | string  | Take profit trigger price                                                                                                                                                                     |
| tpStopType   | string  | Take profit trigger type **MARK_PRICE** **LAST_PRICE**                                                                                                                                        |
| tpOrderType  | string  | Take profit trigger place order type **LIMIT** **MARKET**                                                                                                                                     |
| tpOrderPrice | string  | Take profit trigger place order price. Required if tpOrderType is **LIMIT**                                                                                                                   |
| slPrice      | string  | Stop loss trigger price                                                                                                                                                                       |
| slStopType   | string  | Stop loss trigger type **MARK_PRICE** **LAST_PRICE**                                                                                                                                          |
| slOrderType  | string  | Stop loss trigger place order type **LIMIT** **MARKET**                                                                                                                                       |
| slOrderPrice | string  | Stop loss trigger place order price. Required if slOrderType is **LIMIT**                                                                                                                     |
| ctime        | int64   | Create timestamp                                                                                                                                                                              |
| mtime        | int64   | Latest modify timestamp                                                                                                                                                                       |

## Response Example

```json
{
  "code": 0,
  "data": {
    "orderId": "11111",
    "qty": "1",
    "tradeQty": "0.5",
    "price": "60000",
    "symbol": "BTCUSDT"
    // Additional fields as per response parameters
  }
}
```