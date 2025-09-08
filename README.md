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
- **Connection Resilience**: Automatic WebSocket reconnection with configurable retry logic and subscription restoration
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
ws, err := bitunix.NewPrivateWebsocket(ctx, "YOUR_API_KEY", "YOUR_SECRET_KEY")
if err != nil {
    log.Fatalf("Failed to create WebSocket client: %v", err)
}
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

### Working with Reconnecting WebSockets

The client provides reconnecting WebSocket wrappers that automatically handle connection failures and reestablish
subscriptions for both public and private WebSockets:

#### Reconnecting Private WebSocket

```go
ctx := context.Background()

// Create logger for reconnection events
logger, _ := zap.NewDevelopment()

// Create base private websocket client
baseClient, err := bitunix.NewPrivateWebsocket(ctx, "YOUR_API_KEY", "YOUR_SECRET_KEY",
	bitunix.WithWebsocketLogger(logger),
)
if err != nil {
	log.Fatalf("Failed to create base private client: %v", err)
}

// Create reconnecting private websocket client with custom options
reconnectingClient := bitunix.NewReconnectingPrivateWebsocket(ctx, baseClient,
	bitunix.WithPrivateMaxReconnectAttempts(0), // 0 = infinite attempts, or set a specific limit
	bitunix.WithPrivateReconnectDelay(10*time.Second), // delay between reconnection attempts
	bitunix.WithPrivateReconnectLogger(logger), // optional logger for reconnection events
)

// Connect
if err := reconnectingClient.Connect(); err != nil {
	log.Fatalf("Failed to connect: %v", err)
}

// Subscribe to private channels (subscriptions are automatically restored on reconnection)
balanceSubscriber := &ExampleBalanceSubscriber{}
positionSubscriber := &ExamplePositionSubscriber{}
orderSubscriber := &ExampleOrderSubscriber{}
tpSlOrderSubscriber := &ExampleTpSlOrderSubscriber{}

if err := reconnectingClient.SubscribeBalance(balanceSubscriber); err != nil {
	log.Fatalf("Failed to subscribe to balance: %v", err)
}

if err := reconnectingClient.SubscribePositions(positionSubscriber); err != nil {
	log.Fatalf("Failed to subscribe to positions: %v", err)
}

if err := reconnectingClient.SubscribeOrders(orderSubscriber); err != nil {
	log.Fatalf("Failed to subscribe to orders: %v", err)
}

if err := reconnectingClient.SubscribeTpSlOrders(tpSlOrderSubscriber); err != nil {
	log.Fatalf("Failed to subscribe to TP/SL orders: %v", err)
}

// Start streaming - this will automatically reconnect on failures
go func() {
	if err := reconnectingClient.Stream(); err != nil {
		log.Printf("Stream ended with error: %v", err)
	}
}()

// Keep running
select {}
```

#### Reconnecting Public WebSocket

```go
ctx := context.Background()

// Create logger for reconnection events
logger, _ := zap.NewDevelopment()

// Create base public websocket client
baseClient, err := bitunix.NewPublicWebsocket(ctx, 
	bitunix.WithLogger(logger),
)
if err != nil {
	log.Fatalf("Failed to create base public client: %v", err)
}

// Create reconnecting public websocket client with custom options
reconnectingClient := bitunix.NewReconnectingPublicWebsocket(ctx, baseClient,
	bitunix.WithMaxReconnectAttempts(0), // 0 = infinite attempts, or set a specific limit
	bitunix.WithReconnectDelay(10*time.Second), // delay between reconnection attempts
	bitunix.WithReconnectLogger(logger), // optional logger for reconnection events
)

// Connect
if err := reconnectingClient.Connect(); err != nil {
	log.Fatalf("Failed to connect: %v", err)
}

// Subscribe to public channels (subscriptions are automatically restored on reconnection)
subscriber := &ExampleKLineSubscriber{
	symbol:    model.ParseSymbol("BTCUSDT"),
	interval:  model.Interval1m,
	priceType: model.PriceTypeIndex,
}

if err := reconnectingClient.SubscribeKLine(subscriber); err != nil {
	log.Fatalf("Failed to subscribe: %v", err)
}

// Start streaming - this will automatically reconnect on failures
go func() {
	if err := reconnectingClient.Stream(); err != nil {
		log.Printf("Stream ended with error: %v", err)
	}
}()

// Keep running
select {}
```

## WebSocket Connection Resilience

The library provides robust WebSocket connection management with automatic reconnection capabilities:

### Reconnection Features

- **Automatic Reconnection**: When a WebSocket connection fails, the reconnecting client automatically attempts to
  reestablish the connection
- **Configurable Retry Logic**: Set maximum retry attempts (0 for infinite) and delay between attempts
- **Subscription Persistence**: All active subscriptions are automatically restored after successful reconnection
- **Graceful Error Handling**: Connection errors are logged and handled gracefully without terminating the application
- **Thread-Safe Operations**: The reconnection logic is thread-safe and can be used in concurrent applications

### Reconnection Behavior

1. **Connection Monitoring**: The client continuously monitors the WebSocket connection status
2. **Failure Detection**: When a connection failure is detected (network error, server disconnect, etc.), the client
   immediately marks itself as disconnected
3. **Reconnection Loop**: The client enters a reconnection loop that:
    - Waits for the configured delay period (`WithReconnectDelay`)
    - Attempts to reconnect to the WebSocket server
    - If successful, resubscribes to all previously active channels
    - If failed, increments the attempt counter and retries (unless max attempts reached)
4. **Subscription Restoration**: After successful reconnection, all subscribers are automatically resubscribed to their
   respective channels
5. **Logging**: All reconnection events, failures, and successes are logged for monitoring and debugging

### Configuration Options

**For Public WebSocket:**

```go
// Configure public websocket reconnection behavior
reconnectingClient := bitunix.NewReconnectingPublicWebsocket(ctx, baseClient,
    // Maximum number of reconnection attempts (0 = infinite)
    bitunix.WithMaxReconnectAttempts(5),
    
    // Delay between reconnection attempts
    bitunix.WithReconnectDelay(30*time.Second),
    
    // Logger for reconnection events
    bitunix.WithReconnectLogger(logger),
)
```

**For Private WebSocket:**

```go
// Configure private websocket reconnection behavior
reconnectingClient := bitunix.NewReconnectingPrivateWebsocket(ctx, baseClient,
    // Maximum number of reconnection attempts (0 = infinite)
    bitunix.WithPrivateMaxReconnectAttempts(5),
    
    // Delay between reconnection attempts
    bitunix.WithPrivateReconnectDelay(30*time.Second),
    
    // Logger for reconnection events
    bitunix.WithPrivateReconnectLogger(logger),
)
```

### Error Handling

The reconnecting client handles several types of connection errors:

- Network connectivity issues
- Server-side disconnections
- Context cancellation
- Connection timeout errors

When maximum reconnection attempts are reached, the client returns an error and stops attempting to reconnect.
Applications should monitor this error to decide whether to restart the connection or handle the failure gracefully.

### Best Practices

1. **Use Appropriate Delays**: Set reasonable reconnection delays (5-30 seconds) to avoid overwhelming the server
2. **Monitor Logs**: Enable logging to track connection health and reconnection events
3. **Handle Final Failures**: Implement error handling for cases where reconnection ultimately fails
4. **Resource Management**: Properly disconnect clients when shutting down to release resources
5. **Context Management**: Use proper context cancellation to gracefully stop reconnection attempts

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

The `/samples` and `/examples` directories contain example applications demonstrating the client's functionality:

### REST API Examples

- Account Balance: `/samples/get_account_balance/main.go`
- Place Order: `/samples/place_order/main.go`
- Cancel Order: `/samples/cancel_order/main.go`
- Order History: `/samples/get_order_history/main.go`
- Position History: `/samples/get_position_history/main.go`
- Trade History: `/samples/get_trade_history/main.go`
- Take-Profit Order: `/samples/place_tp/main.go`

### WebSocket Examples

- WebSocket Private Client: `/samples/websocket_private/main.go`
- WebSocket Public Client: `/samples/websocket_public/main.go`
- **Reconnecting Public WebSocket Client**: `/examples/reconnecting_websocket.go` - Demonstrates automatic reconnection
  with subscription restoration for public WebSocket
- **Reconnecting Private WebSocket Client**: `/examples/reconnecting_private_websocket.go` - Demonstrates automatic
  reconnection with subscription restoration for private WebSocket

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