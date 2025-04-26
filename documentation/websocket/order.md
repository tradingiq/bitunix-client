# Order Channel Subscription

## Description

Subscribe the order channel.  
Data will be pushed when the following events occured:

1. Open/Close orders are created
2. Open/Close orders are filled
3. Orders canceled

## Push Parameters

| Parameter       |                                                 | Type         | Description                                                                                                                                                                    |
|-----------------|-------------------------------------------------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ch              |                                                 | String       | Channel name: position                                                                                                                                                         |
| ts              |                                                 | Int64        | Time stamp                                                                                                                                                                     |
| data            | #broken: returns Object instead of List<Object> | List<Object> | Subscription data                                                                                                                                                              |
| > event         |                                                 | String       | Event: CREATE/UPDATE/CLOSE                                                                                                                                                     |
| > orderId       |                                                 | String       | order Id                                                                                                                                                                       |
| > symbol        |                                                 | String       | Symbol                                                                                                                                                                         |
| > positionType  |                                                 | String       | Margin mode: ISOLATION/CROSS                                                                                                                                                   |
| > positionMode  |                                                 | String       | Position mode: ONE_WAY/HEDGE                                                                                                                                                   |
| > side          | #broken actual value is "Sell" or "Buy"         | String       | Sell: BUY/SELL                                                                                                                                                                 |
| > effect        |                                                 | String       | Order expiration date. Required if the orderType is limit.<br>IOC: Immediate or cancel<br>FOK: Fill or kill<br>GTC: Good till canceled (default value)<br>POST_ONLY: POST only |
| > type          | #broken actual value is "Market" or "Limit"     | String       | LIMIT/MARKET                                                                                                                                                                   |
| > qty           |                                                 | String       | Amount (base coin)                                                                                                                                                             |
| > reductionOnly |                                                 | Bool         | Reduction Only                                                                                                                                                                 |
| > price         |                                                 | String       | Price of the order. Required if the order type is **LIMIT**                                                                                                                    |
| > ctime         |                                                 | String       | create timestamp                                                                                                                                                               |
| > mtime         |                                                 | String       | create timestamp                                                                                                                                                               |
| > leverage      |                                                 | String       | Leverage                                                                                                                                                                       |
| > orderStatus   | #broken missing value SYSTEM_CANCELED           | String       | INIT: prepare status<br>NEW: pending<br>PART_FILLED: partially filled<br>CANCELED: canceled<br>FILLED: All filled                                                              |
| > fee           |                                                 | String       | Deducted transaction fees: transaction fees deducted during the position                                                                                                       |
| > tpStopType    |                                                 | String       | take profit trigger type: MARK_PRICE/LAST_PRICE                                                                                                                                |
| > tpPrice       |                                                 | String       | take profit trigger price                                                                                                                                                      |
| > tpOrderType   |                                                 | String       | take profit trigger place order type: LIMIT/MARKET                                                                                                                             |
| > tpOrderPrice  |                                                 | String       | take profit trigger place order price: LIMIT/MARKET                                                                                                                            |
| > slStopType    |                                                 | String       | stop loss trigger type: MARK_PRICE/LAST_PRICE                                                                                                                                  |
| > slPrice       |                                                 | String       | stop loss trigger price                                                                                                                                                        |
| > slOrderType   |                                                 | String       | stop loss trigger place order type: LIMIT/MARKET                                                                                                                               |
| > slOrderPrice  |                                                 | String       | stop loss trigger place order price: LIMIT/MARKET                                                                                                                              |