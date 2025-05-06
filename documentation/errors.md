# Bitunix Client Error Reference

This document provides a comprehensive list of all possible errors that can be returned by each interface in the Bitunix
client library.

## Error Types

The Bitunix client uses the following error types:

| Error Type            | Description                               |
|-----------------------|-------------------------------------------|
| `ValidationError`     | Input validation errors                   |
| `NetworkError`        | Network-related errors                    |
| `APIError`            | Errors returned by the Bitunix API        |
| `AuthenticationError` | Authentication-related errors             |
| `WebsocketError`      | Websocket connection and operation errors |
| `InternalError`       | Internal client errors                    |
| `TimeoutError`        | Request timeout errors                    |

## API Error Codes

For the complete list of Bitunix API error codes, see [error_codes.md](error_codes.md)

## Errors by Interface

### PlaceOrder

```go
PlaceOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error)
```

| Error                 | Conditions                                                                                                                                                                                                                                                                                                                        |
|-----------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `ValidationError`     | - Missing symbol<br/>- Missing trade action (side)<br/>- Missing trade side<br/>- Missing or invalid quantity<br/>- Missing position ID when closing<br/>- Missing order type<br/>- Missing price for limit orders<br/>- Invalid take profit/stop loss configuration                                                              |
| `NetworkError`        | - API connection issues<br/>- HTTP request failures                                                                                                                                                                                                                                                                               |
| `APIError`            | - Market does not exist (10020)<br/>- Position limit exceeded (10022)<br/>- Insufficient balance (10023, 10028, 10210)<br/>- Invalid leverage (10025)<br/>- Order price issues (30001-30003)<br/>- Order quantity issues (30013, 30016, 30017, 30039)<br/>- Reduce-only settings (30018, 30019)<br/>- Client ID duplicate (30042) |
| `AuthenticationError` | - Missing API key (10003)<br/>- IP not in whitelist (10004)<br/>- Signature error (10007)                                                                                                                                                                                                                                         |
| `TimeoutError`        | - Request exceeds timeout duration                                                                                                                                                                                                                                                                                                |
| `InternalError`       | - JSON marshalling failures                                                                                                                                                                                                                                                                                                       |

### CancelOrders

```go
CancelOrders(ctx context.Context, request *model.CancelOrderRequest) (*model.CancelOrderResponse, error)
```

| Error                 | Conditions                                                                                |
|-----------------------|-------------------------------------------------------------------------------------------|
| `ValidationError`     | - Missing symbol<br/>- Empty order list<br/>- Missing both order ID and client ID         |
| `NetworkError`        | - API connection issues<br/>- HTTP request failures                                       |
| `APIError`            | - Order not found (20007)<br/>- Market does not exist (20001)                             |
| `AuthenticationError` | - Missing API key (10003)<br/>- IP not in whitelist (10004)<br/>- Signature error (10007) |
| `TimeoutError`        | - Request exceeds timeout duration                                                        |
| `InternalError`       | - JSON marshalling failures                                                               |

### GetTradeHistory

```go
GetTradeHistory(ctx context.Context, params model.TradeHistoryParams) (*model.TradeHistoryResponse, error)
```

| Error                 | Conditions                                                                                |
|-----------------------|-------------------------------------------------------------------------------------------|
| `NetworkError`        | - API connection issues<br/>- HTTP request failures                                       |
| `APIError`            | - Invalid parameters (10002)<br/>- Market does not exist (20001)                          |
| `AuthenticationError` | - Missing API key (10003)<br/>- IP not in whitelist (10004)<br/>- Signature error (10007) |
| `TimeoutError`        | - Request exceeds timeout duration                                                        |
| `InternalError`       | - Response parsing errors                                                                 |

### GetOrderHistory

```go
GetOrderHistory(ctx context.Context, params model.OrderHistoryParams) (*model.OrderHistoryResponse, error)
```

| Error                 | Conditions                                                                                |
|-----------------------|-------------------------------------------------------------------------------------------|
| `NetworkError`        | - API connection issues<br/>- HTTP request failures                                       |
| `APIError`            | - Invalid parameters (10002)<br/>- Market does not exist (20001)                          |
| `AuthenticationError` | - Missing API key (10003)<br/>- IP not in whitelist (10004)<br/>- Signature error (10007) |
| `TimeoutError`        | - Request exceeds timeout duration                                                        |
| `InternalError`       | - Response parsing errors                                                                 |

### GetPositionHistory

```go
GetPositionHistory(ctx context.Context, params model.PositionHistoryParams) (*model.PositionHistoryResponse, error)
```

| Error                 | Conditions                                                                                |
|-----------------------|-------------------------------------------------------------------------------------------|
| `NetworkError`        | - API connection issues<br/>- HTTP request failures                                       |
| `APIError`            | - Invalid parameters (10002)<br/>- Market does not exist (20001)                          |
| `AuthenticationError` | - Missing API key (10003)<br/>- IP not in whitelist (10004)<br/>- Signature error (10007) |
| `TimeoutError`        | - Request exceeds timeout duration                                                        |
| `InternalError`       | - Response parsing errors                                                                 |

### PlaceTpSlOrder

```go
PlaceTpSlOrder(ctx context.Context, request *model.TPSLOrderRequest) (*model.TpSlOrderResponse, error)
```

| Error                 | Conditions                                                                                                                                                               |
|-----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `ValidationError`     | - Missing symbol<br/>- Missing position ID<br/>- Neither TP nor SL set<br/>- Invalid TP/SL price (≤ 0)<br/>- Missing TP/SL quantity<br/>- Missing price for limit orders |
| `NetworkError`        | - API connection issues<br/>- HTTP request failures                                                                                                                      |
| `APIError`            | - Position not exist (30004)<br/>- TP/SL errors (30005-30038)<br/>- Invalid trigger price (30041)<br/>- Market does not exist (20001)                                    |
| `AuthenticationError` | - Missing API key (10003)<br/>- IP not in whitelist (10004)<br/>- Signature error (10007)                                                                                |
| `TimeoutError`        | - Request exceeds timeout duration                                                                                                                                       |
| `InternalError`       | - JSON marshalling failures                                                                                                                                              |

### GetAccountBalance

```go
GetAccountBalance(ctx context.Context, params model.AccountBalanceParams) (*model.AccountBalanceResponse, error)
```

| Error                 | Conditions                                                                                |
|-----------------------|-------------------------------------------------------------------------------------------|
| `ValidationError`     | - Missing margin coin                                                                     |
| `NetworkError`        | - API connection issues<br/>- HTTP request failures                                       |
| `APIError`            | - Account not allowed to trade (20011)<br/>- Account inactive or deleted (20013, 20014)   |
| `AuthenticationError` | - Missing API key (10003)<br/>- IP not in whitelist (10004)<br/>- Signature error (10007) |
| `TimeoutError`        | - Request exceeds timeout duration                                                        |
| `InternalError`       | - Response parsing errors                                                                 |

### Public Websocket Client

```go
type PublicWebsocketClient interface {
Stream() error
Connect() error
Disconnect()
SubscribeKLine(subscriber KLineSubscriber) error
UnsubscribeKLine(subscriber KLineSubscriber) error
}
```

| Method               | Error             | Conditions                                        |
|----------------------|-------------------|---------------------------------------------------|
| `Connect()`          | `WebsocketError`  | - Connection failure<br/>- Network issues         |
| `Stream()`           | `WebsocketError`  | - Listener initialization failure                 |
| `SubscribeKLine()`   | `ValidationError` | - Nil subscriber<br/>- Invalid channel parameters |
|                      | `InternalError`   | - JSON marshalling failures                       |
|                      | `WebsocketError`  | - Subscription request sending failure            |
| `UnsubscribeKLine()` | `ValidationError` | - Nil subscriber<br/>- Invalid channel parameters |
|                      | `InternalError`   | - JSON marshalling failures                       |
|                      | `WebsocketError`  | - Unsubscription request sending failure          |

### Private Websocket Client

```go
type PrivateWebsocketClient interface {
SubscribeBalance(subscriber BalanceSubscriber) error
UnsubscribeBalance(subscriber BalanceSubscriber) error
SubscribePositions(subscriber PositionSubscriber) error
UnsubscribePosition(subscriber PositionSubscriber) error
SubscribeOrders(subscriber OrderSubscriber) error
UnsubscribeOrders(subscriber OrderSubscriber) error
SubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error
UnsubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error
Stream() error
Connect() error
Disconnect()
}
```

| Method                    | Error                 | Conditions                                              |
|---------------------------|-----------------------|---------------------------------------------------------|
| `Connect()`               | `WebsocketError`      | - Connection failure<br/>- Network issues               |
|                           | `AuthenticationError` | - Authentication failures<br/>- Invalid API credentials |
| `Stream()`                | `WebsocketError`      | - Listener initialization failure                       |
| `SubscribeBalance()`      | `ValidationError`     | - Nil subscriber                                        |
| `UnsubscribeBalance()`    | `ValidationError`     | - Nil subscriber                                        |
| `SubscribePositions()`    | `ValidationError`     | - Nil subscriber                                        |
| `UnsubscribePosition()`   | `ValidationError`     | - Nil subscriber                                        |
| `SubscribeOrders()`       | `ValidationError`     | - Nil subscriber                                        |
| `UnsubscribeOrders()`     | `ValidationError`     | - Nil subscriber                                        |
| `SubscribeTpSlOrders()`   | `ValidationError`     | - Nil subscriber                                        |
| `UnsubscribeTpSlOrders()` | `ValidationError`     | - Nil subscriber                                        |

## Common Error Categories by Operation Type

### Order-Related Errors

- `ErrOrderNotFound`: When trying to cancel or modify a non-existent order
- `ErrDuplicateClientID`: When placing an order with a client ID that already exists
- `ErrOrderPriceIssue`: Price-related issues for orders
- `ErrOrderQuantityIssue`: Quantity-related issues for orders
- `ErrTPSLOrderError`: Take profit/stop loss order validation errors

### Account-Related Errors

- `ErrInsufficientBalance`: Insufficient funds to perform the operation
- `ErrAccountNotAllowed`: Account not permitted to trade
- `ErrAccountInactive`: Account is inactive or has been deleted
- `ErrPositionLimitExceeded`: Position amount exceeded maximum open limit
- `ErrInvalidLeverage`: Invalid leverage setting

### Position-Related Errors

- `ErrPositionNotExist`: Position does not exist
- `ErrPositionsModeChange`: Position mode cannot be updated
- `ErrOpenOrdersExist`: Open orders prevent another operation

### Market-Related Errors

- `ErrMarketNotExists`: Market specified does not exist
- `ErrFuturesNotSupported`: Futures not supported or allowed for this operation

### Authentication Errors

- `ErrAuthentication`: General authentication error
- `ErrSignatureError`: Signature validation failed
- `ErrIPNotAllowed`: IP not in API key whitelist

### Network Errors

- `ErrNetwork`: General network error
- `ErrTimeout`: Request timeout

### Validation Errors

- `ErrValidation`: General validation error
- `ErrParameterError`: Invalid parameter values
- `ErrInvalidValue`: Value does not comply with rule

### Rate Limiting Errors

- `ErrRateLimitExceeded`: Rate limit exceeded, retry after waiting

## Error Handling Recommendations

1. **Authentication Errors**: Check API key and secret configuration
2. **Validation Errors**: Verify input parameters before making API calls
3. **Network Errors**: Implement retry logic with exponential backoff
4. **API Errors**: Handle based on specific operation context
5. **Rate Limit Errors**: Implement rate limiting on the client side
6. **Timeout Errors**: Consider adjusting timeout settings