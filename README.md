# Bitunix Client

[![codecov](https://codecov.io/github/tradingiq/bitunix-client/graph/badge.svg?token=21PSK71QHJ)](https://codecov.io/github/tradingiq/bitunix-client)

A Go client library for interacting with the Bitunix cryptocurrency futures exchange API.

## Project Overview

This Go library provides a comprehensive client for interacting with the Bitunix cryptocurrency exchange API, focusing
on futures trading. It implements both REST API endpoints and WebSocket connections for real-time data streaming.

The client is designed to be reliable, efficient, and easy to use, with built-in support for request signing, WebSocket
connection management, and detailed error handling.

## Key Features

- **REST API Integration**: Full implementation of Bitunix's futures trading REST API endpoints
- **WebSocket Support**: Real-time data streaming for account balances, positions, orders, and take-profit/stop-loss
  orders
- **Flexible WebSocket Client**: Interface-based WebSocket client implementation for easy mocking and testing
- **Order Management**: Place, cancel, and query orders with support for various order types
- **Account Information**: Retrieve account balances and trading history
- **Position Management**: Track and manage trading positions
- **Comprehensive Type System**: Strong typing for all API objects and parameters
- **Builder Pattern**: Order builders to simplify creating complex order requests
- **Authentication**: Automatic request signing and authentication
- **Connection Resilience**: Improved WebSocket reconnection and error handling
- **Thread Safety**: Race-condition free implementation for concurrent usage
- **Work Group Pattern**: Limits resource extension for websocket connections

## Installation

```go
go get github.com/tradingiq/bitunix-client
```

## Usage

### Setting up the client

```go
// Create Bitunix client with authentication
client, _ := bitunix.NewApiClient("YOUR_API_KEY", "YOUR_SECRET_KEY")
```

### Getting account balance

```go
params := model.AccountBalanceParams{
MarginCoin: "USDT",
}

ctx := context.Background()
response, err := client.GetAccountBalance(ctx, params)
if err != nil {
log.Fatal(err)
}

fmt.Printf("Account Balance: %.6f %s\n", response.Data.Available, response.Data.MarginCoin)
```

### Placing an order

```go
// Create a limit order using the builder pattern
limitOrder := bitunix.NewOrderBuilder(
model.ParseSymbol("BTCUSDT"), // Symbol
model.TradeSideSell,          // Side (BUY/SELL)
model.SideOpen,               // Trade side (OPEN/CLOSE)
0.002,                        // Quantity
).WithOrderType(model.OrderTypeLimit).
WithPrice(100000.0).
WithTimeInForce(model.TimeInForcePostOnly).
Build()

// Submit the order
response, err := client.PlaceOrder(ctx, &limitOrder)
if err != nil {
log.Fatalf("Failed to place order: %v", err)
}

fmt.Printf("Order placed successfully: %+v\n", response)
```

### Working with WebSockets (Private)

```go
// Create a WebSocket client
ctx := context.Background()
ws := bitunix.NewPrivateWebsocket(ctx, "YOUR_API_KEY", "YOUR_SECRET_KEY")
defer ws.Disconnect()

// Connect to the WebSocket server
if err := ws.Connect(); err != nil {
log.Fatalf("Failed to connect to WebSocket: %v", err)
}

// Subscribe to balance updates
balanceHandler := &BalanceHandler{}
if err := ws.SubscribeBalance(balanceHandler); err != nil {
log.WithError(err).Fatal("Failed to subscribe to balance")
}

// Subscribe to position updates
positionHandler := &PositionHandler{}
if err := ws.SubscribePositions(positionHandler); err != nil {
log.WithError(err).Fatal("Failed to subscribe to positions")
}

// Subscribe to order updates
orderHandler := &OrderHandler{}
if err := ws.SubscribeOrders(orderHandler); err != nil {
log.WithError(err).Fatal("Failed to subscribe to orders")
}

// Subscribe to TP/SL order updates
tpslHandler := &TpSlOrderHandler{}
if err := ws.SubscribeTpSlOrders(tpslHandler); err != nil {
log.WithError(err).Fatal("Failed to subscribe to TP/SL orders")
}

// Start the WebSocket stream
if err := ws.Stream(); err != nil {
log.WithError(err).Fatal("Failed to stream")
}
```

### Working with WebSockets (Public)

```go
ctx := context.Background()
ws := bitunix.NewPublicWebsocket(ctx)
defer ws.Disconnect()

if err := ws.Connect(); err != nil {
log.Fatalf("Failed to connect to WebSocket: %v", err)
}

// Set up KLine subscription
interval, _ := model.ParseInterval("1m")
symbol := model.ParseSymbol("BTCUSDT")
priceType, _ := model.ParsePriceType("mark")

sub := &KLineSubscriber{
interval:  interval,
symbol:    symbol,
priceType: priceType,
}

err = ws.SubscribeKLine(sub)
if err != nil {
log.Fatalf("Failed to subscribe: %v", err)
}

if err := ws.Stream(); err != nil {
log.WithError(err).Fatal("Failed to stream")
}
```

## Documentation

The project includes detailed documentation in the `/documentation` directory:

### REST API Endpoints

- Account Balance: `/documentation/get_account_balance.md`
- Place Order: `/documentation/place_order.md`
- Cancel Order: `/documentation/cancel_order.md`
- Historical Orders: `/documentation/get_history_orders.md`
- Historical Positions: `/documentation/get_history_positions.md`
- Historical Trades: `/documentation/get_history_trades.md`
- Take-Profit/Stop-Loss Orders: `/documentation/place_tpsl_order.md`

### WebSocket Channels

- Balance: `/documentation/websocket/balance.md`
- KLine: `/documentation/websocket/kline.md`
- Order: `/documentation/websocket/order.md`
- Position: `/documentation/websocket/position.md`
- Take-Profit/Stop-Loss: `/documentation/websocket/tpsl.md`

## Sample Applications

The `/samples` directory contains example applications demonstrating the client's functionality:

- Account Balance: `/samples/get_account_balance/main.go`
- Place Order: `/samples/place_order/main.go`
- Cancel Order: `/samples/cancel_order/main.go`
- Order History: `/samples/get_order_history/main.go`
- Position History: `/samples/get_position_history/main.go`
- Trade History: `/samples/get_trade_history/main.go`
- Take-Profit Order: `/samples/place_tp/main.go`
- WebSocket Private Client: `/samples/websocket_private/main.go`
- WebSocket Public Client: `/samples/websocket_public/main.go`

## Authentication

The client handles authentication automatically by generating the required API signatures for both REST and WebSocket
connections. Configure the client with your API and Secret keys from Bitunix:

```go
// For REST API
client, _ := bitunix.NewApiClient("YOUR_API_KEY", "YOUR_SECRET_KEY")

// For WebSocket (private)
ws := bitunix.NewPrivateWebsocket(ctx, "YOUR_API_KEY", "YOUR_SECRET_KEY")
```

## Error Types

The package defines several error types for different categories of errors:

- `ValidationError`: Errors related to input validation
- `NetworkError`: Errors related to network operations
- `APIError`: Errors returned by the BitUnix API
- `AuthenticationError`: Errors related to authentication
- `WebsocketError`: Errors related to websocket operations
- `InternalError`: Internal client errors
- `TimeoutError`: Errors related to timeouts

## Sentinel Errors

The package defines sentinel errors that can be used with `errors.Is()` to check for specific error types:

- `ErrValidation`
- `ErrNetwork`
- `ErrAPI`
- `ErrAuthentication`
- `ErrWebsocket`
- `ErrInternal`
- `ErrTimeout`

### Checking Error Types

```go
// Check if an error is a validation error
if errors.Is(err, errors.ErrValidation) {
// Handle validation error
}

// Check if an error is a network error
if errors.Is(err, errors.ErrNetwork) {
// Handle network error
}

// Type assertion to get more details
if apiErr, ok := err.(*errors.APIError); ok {
fmt.Printf("API error code: %d, message: %s\n", apiErr.Code, apiErr.Message)
}

err := client.PlaceOrder(...)
if errors.Is(err, errors.ErrInsufficientBalance) {
// Handle insufficient balance case
}

```

## Configuration

The samples use a configuration mechanism to load API credentials from environment variables:

```go
// Configuration is loaded from environment variables
// API_KEY and SECRET_KEY must be set
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to
discuss what you would like to change.

## License

This project is licensed under the MIT License with Attribution - see the LICENSE file for details.

## Contributors

See [CONTRIBUTORS.md](CONTRIBUTORS.md) for a list of project contributors.