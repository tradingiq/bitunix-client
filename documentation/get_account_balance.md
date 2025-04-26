# Get Single Account

**Rate Limit:** 10 req/sec/uid

## Description

Get account details with the given 'marginCoin'

## HTTP Request

* GET `/api/v1/futures/account`

## Request Parameters

| Parameter  | Type   | Required | Description |
|------------|--------|----------|-------------|
| marginCoin | string | true     | Margin coin |

## Request Example

```bash
curl -X 'GET'  --location 'https://fapi.bitunix.com/api/v1/futures/account?marginCoin=USDT' \
   -H "api-key:*******" \
   -H "sign:*" \
   -H "nonce:your-nonce" \
   -H "time:1659076670000" \
   -H "language:en-US" \
   -H "Content-Type: application/json"
```

## Response Parameters

| Parameter              | Type   | Description                                                                                     |
|------------------------|--------|-------------------------------------------------------------------------------------------------|
| marginCoin             | string | Margin Coin                                                                                     |
| available              | string | Available quantity in the account. This field + crossUnrealizedPNL = Actual maximum open amount |
| frozen                 | string | locked quantity of orders                                                                       |
| margin                 | string | locked quantity of positions                                                                    |
| transfer               | string | Maximum transferable amount                                                                     |
| positionMode           | string | Position mode<br>**ONE_WAY**<br>**HEDGE**                                                       |
| crossUnrealizedPNL     | string | unrealizedPNL for cross positions                                                               |
| isolationUnrealizedPNL | string | unrealizedPNL for isolation positions                                                           |
| bonus                  | string | Futures Bonus                                                                                   |

## Response Example

#### #broken example shows data.[]{} while it actually is data.{}

```json
{
  "code": 0,
  "data": [
    {
      "marginCoin": "USDT",
      "available": "1000",
      "frozen": "0",
      "margin": "10",
      "transfer": "1000",
      "positionMode": "HEDGE",
      "crossUnrealizedPNL": "2",
      "isolationUnrealizedPNL": "0",
      "bonus": "0"
    }
  ],
  "msg": "Success"
}
```