# TP/SL Order Channel Subscription

## Description

TP/SL Order

## Push Parameters

| Parameter      |                                                             | Type         | Description                                                                                                       |
|----------------|-------------------------------------------------------------|--------------|-------------------------------------------------------------------------------------------------------------------|
| ch             |                                                             | String       | Channel name: position                                                                                            |
| ts             |                                                             | Int64        | Time stamp                                                                                                        |
| data           | #broken not a List<Object> but Object                       | List<Object> | Subscription data                                                                                                 |
| > event        |                                                             | String       | Event: CREATE/UPDATE/CLOSE                                                                                        |
| > positionId   |                                                             | String       | position Id                                                                                                       |
| > orderId      |                                                             | String       | order Id                                                                                                          |
| > symbol       |                                                             | String       | Symbol                                                                                                            |
| > leverage     |                                                             | String       | Leverage                                                                                                          |
| > side         | #broken actual return value is "Sell" or "Buy"              | String       | Sell: BUY/SELL                                                                                                    |
| > positionMode |                                                             | String       | Position mode: ONE_WAY/HEDGE                                                                                      |
| > status       | #broken possible value *SYSTEM_CANCELED* missing            | String       | INIT: prepare status<br>NEW: pending<br>PART_FILLED: partially filled<br>CANCELED: canceled<br>FILLED: All filled |
| > ctime        |                                                             | String       | create timestamp                                                                                                  |
| > type         | #broken actual return value is "TPSL" or "POSITION_TPSL"    | String       | LIMIT/MARKET                                                                                                      |
| > tpQty        |                                                             | String       | Take-profit order quantity(base coin)<br>At least one of `tpQty` or `slQty` is required.                          |
| > slQty        | #broken documentation shows bool but a quantity is returned | Bool         | Stop-loss order quantity(base coin)<br>At least one of `tpQty` or `slQty` is required.                            |
| > tpStopType   | #broken actual return value is "Last"                       | String       | take profit trigger type: MARK_PRICE/LAST_PRICE                                                                   |
| > tpPrice      |                                                             | String       | take profit trigger price                                                                                         |
| > tpOrderType  | #broken actual return value is "Market" or "Limit"          | String       | take profit trigger place order type: LIMIT/MARKET                                                                |
| > tpOrderPrice |                                                             | String       | take profit trigger place order price: LIMIT/MARKET                                                               |
| > slStopType   | #broken actual return value is "Last"                       | String       | stop loss trigger type: MARK_PRICE/LAST_PRICE                                                                     |
| > slPrice      |                                                             | String       | stop loss trigger price                                                                                           |
| > slOrderType  | #broken actual return value is "Market" or "Limit"          | String       | stop loss trigger place order type: LIMIT/MARKET                                                                  |
| > slOrderPrice |                                                             | String       | stop loss trigger place order price: LIMIT/MARKET                                                                 |