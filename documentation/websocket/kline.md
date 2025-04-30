## Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| op | String | Yes | Operation, subscribe unsubscribe |
| args | List\<Object\> | Yes | List of channels to request subscription |
| > ch | String | Yes | Channel name, The subscription channel is: Price Type\*klineTime Interval；The price types include market price and marked price；market\*\_\*kline\*\_\*1min, mark\*\_\*kline\*\_\*1min, market\*\_\*kline\*\_\*3min, mark\*\_\*kline\*\_\*3min, market\*\_\*kline\*\_\*5min, mark\*\_\*kline\*\_\*5min, market\*\_\*kline\*\_\*15min, mark\*\_\*kline\*\_\*15min, market\*\_\*kline\*\_\*30min, mark\*\_\*kline\*\_\*30min, market\*\_\*kline\*\_\*60min, mark\*\_\*kline\*\_\*60min, market\*\_\*kline\*\_\*2h, mark\*\_\*kline\*\_\*2h, market\*\_\*kline\*\_\*4h, mark\*\_\*kline\*\_\*4h, market\*\_\*kline\*\_\*6h, mark\*\_\*kline\*\_\*6h, market\*\_\*kline\*\_\*8h, mark\*\_\*kline\*\_\*8h, market\*\_\*kline\*\_\*12h, mark\*\_\*kline\*\_\*12h, market\*\_\*kline\*\_\*1day, mark\*\_\*kline\*\_\*1day, market\*\_\*kline\*\_\*3day, mark\*\_\*kline\*\_\*3day, market\*\_\*kline\*\_\*1week, mark\*\_\*kline\*\_\*1week, market\*\_\*kline\*\_\*1month, mark\*\_\*kline\*\_\*1month* |
| > symbol | String | Yes | Product ID E.g. ETHUSDT |

### Request Example:

```json
{
    "op":"subscribe",
    "args":[
        {
            "symbol":"BTCUSDT",
            "ch":"market_kline_1min" 
        }
    ]
}
```

## Push Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| ch | String | Channel name |
| symbol | String | Product ID E.g. ETHUSDT |
| ts | int64 | Time stamp |
| data | Object | Subscription data |
| > o | String | Opening price |
| > h | String | Highest price |
| > l | String | Lowest price |
| > c | String | Closing price |
| > b | String | Trading volume of the coin |
| > q | String | Trading volume of quote currency |

### Push Data Example:

```json
{ 
  "ch": "mark_kline_1min",
  "symbol": "BNBUSDT",
  "ts": 1732178884994,                   
  "data":{
      "o": "0.0010",                     
      "c": "0.0020",                     
      "h": "0.0025",                     
      "l": "0.0015",                    
      "b": "1.01",                     
      "q": "1.09"                         
  }
}
```