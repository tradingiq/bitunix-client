package bitunix

type StopType string

const (
	StopTypeLastPrice StopType = "LAST_PRICE"
	StopTypeMarkPrice StopType = "MARK_PRICE"
)

type OrderType string

const (
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"
)

type TimeInForce string

const (
	TimeInForceIOC      TimeInForce = "IOC"       // Immediate or cancel
	TimeInForceFOK      TimeInForce = "FOK"       // Fill or kill
	TimeInForceGTC      TimeInForce = "GTC"       // Good till canceled (default value)
	TimeInForcePostOnly TimeInForce = "POST_ONLY" // POST only
)

type TradeSide string

const (
	TradeSideOpen  TradeSide = "OPEN"
	TradeSideClose TradeSide = "CLOSE"
)

type TradeAction string

const (
	TradeActionBuy  TradeAction = "BUY"
	TradeActionSell TradeAction = "SELL"
)
