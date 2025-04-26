package model

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

type Side string

const (
	SideOpen  Side = "OPEN"
	SideClose Side = "CLOSE"
)

type TradeSide string

const (
	TradeSideBuy  TradeSide = "BUY"
	TradeSideSell TradeSide = "SELL"
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

type PositionMode string

const (
	PositionModeOneWay PositionMode = "ONE_WAY"
	PositionModeHedge  PositionMode = "HEDGE"
)

const ChannelBalance = "balance"
