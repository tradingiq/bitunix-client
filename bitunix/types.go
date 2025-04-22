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
	TimeInForceIOC      TimeInForce = "IOC"
	TimeInForceFOK      TimeInForce = "FOK"
	TimeInForceGTC      TimeInForce = "GTC"
	TimeInForcePostOnly TimeInForce = "POST_ONLY"
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

type MarginMode string

const (
	MarginModeIsolation MarginMode = "ISOLATION"
	MarginModeCross     MarginMode = "CROSS"
)

type TradePositionMode string

const (
	TradePositionModeOneWay TradePositionMode = "ONE_WAY"
	TradePositionModeHedge  TradePositionMode = "HEDGE"
)

type TradeRoleType string

const (
	TradeRoleTypeTaker TradeRoleType = "TAKER"
	TradeRoleTypeMaker TradeRoleType = "MAKER"
)
