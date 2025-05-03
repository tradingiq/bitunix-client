package model

import (
	"fmt"
	"strings"
)

type StopType string

const (
	StopTypeLastPrice StopType = "LAST_PRICE"
	StopTypeMarkPrice StopType = "MARK_PRICE"
)

func (s StopType) IsValid() bool {
	switch s {
	case StopTypeLastPrice, StopTypeMarkPrice:
		return true
	}
	return false
}

func (s StopType) String() string {
	return string(s)
}

func ParseStopType(s string) (StopType, error) {
	stopType := StopType(s)
	stopType = stopType.Normalize()

	if !stopType.IsValid() {
		return stopType, fmt.Errorf("%s is not a valid StopType", s)
	}

	return stopType, nil
}

func (s StopType) Normalize() StopType {
	upper := strings.ToUpper(string(s))

	switch upper {
	case "LAST":
		upper = string(StopTypeLastPrice)
	case "MARK":
		upper = string(StopTypeMarkPrice)
	}

	return StopType(upper)
}

type OrderType string

const (
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"
)

func (o OrderType) IsValid() bool {
	switch o {
	case OrderTypeLimit, OrderTypeMarket:
		return true
	}
	return false
}

func (o OrderType) String() string {
	return string(o)
}

func ParseOrderType(s string) (OrderType, error) {
	orderType := OrderType(s)
	orderType = orderType.Normalize()

	if !orderType.IsValid() {
		return orderType, fmt.Errorf("%s is not a valid OrderType", s)
	}

	return orderType, nil
}

func (o OrderType) Normalize() OrderType {
	return OrderType(strings.ToUpper(string(o)))
}

type TimeInForce string

const (
	TimeInForceIOC      TimeInForce = "IOC"
	TimeInForceFOK      TimeInForce = "FOK"
	TimeInForceGTC      TimeInForce = "GTC"
	TimeInForcePostOnly TimeInForce = "POST_ONLY"
)

func (t TimeInForce) IsValid() bool {
	switch t {
	case TimeInForceIOC, TimeInForceFOK, TimeInForceGTC, TimeInForcePostOnly:
		return true
	}
	return false
}

func (t TimeInForce) String() string {
	return string(t)
}

func ParseTimeInForce(s string) (TimeInForce, error) {
	timeInForce := TimeInForce(s)
	timeInForce = timeInForce.Normalize()

	if !timeInForce.IsValid() {
		return timeInForce, fmt.Errorf("%s is not a valid TimeInForce", s)
	}

	return timeInForce, nil
}

func (t TimeInForce) Normalize() TimeInForce {
	return TimeInForce(strings.ToUpper(string(t)))
}

type Side string

const (
	SideOpen  Side = "OPEN"
	SideClose Side = "CLOSE"
)

func (s Side) IsValid() bool {
	switch s {
	case SideOpen, SideClose:
		return true
	}
	return false
}

func (s Side) String() string {
	return string(s)
}

func ParseSide(str string) (Side, error) {
	side := Side(str)
	side = side.Normalize()

	if !side.IsValid() {
		return side, fmt.Errorf("%s is not a valid Side", str)
	}

	return side, nil
}

func (s Side) Normalize() Side {
	return Side(strings.ToUpper(string(s)))
}

type TradeSide string

const (
	TradeSideBuy  TradeSide = "BUY"
	TradeSideSell TradeSide = "SELL"
)

func (t TradeSide) IsValid() bool {
	switch t {
	case TradeSideBuy, TradeSideSell:
		return true
	}
	return false
}

func (t TradeSide) String() string {
	return string(t)
}

func ParseTradeSide(s string) (TradeSide, error) {
	tradeSide := TradeSide(s)
	tradeSide = tradeSide.Normalize()

	if !tradeSide.IsValid() {
		return tradeSide, fmt.Errorf("%s is not a valid TradeSide", s)
	}

	return tradeSide, nil
}

func (t TradeSide) Normalize() TradeSide {
	return TradeSide(strings.ToUpper(string(t)))
}

type MarginMode string

const (
	MarginModeIsolation MarginMode = "ISOLATION"
	MarginModeCross     MarginMode = "CROSS"
)

func (m MarginMode) IsValid() bool {
	switch m {
	case MarginModeIsolation, MarginModeCross:
		return true
	}
	return false
}

func (m MarginMode) String() string {
	return string(m)
}

func ParseMarginMode(s string) (MarginMode, error) {
	marginMode := MarginMode(s)
	marginMode = marginMode.Normalize()

	if !marginMode.IsValid() {
		return marginMode, fmt.Errorf("%s is not a valid MarginMode", s)
	}

	return marginMode, nil
}

func (m MarginMode) Normalize() MarginMode {
	return MarginMode(strings.ToUpper(string(m)))
}

type TradeRoleType string

const (
	TradeRoleTypeTaker TradeRoleType = "TAKER"
	TradeRoleTypeMaker TradeRoleType = "MAKER"
)

func (t TradeRoleType) IsValid() bool {
	switch t {
	case TradeRoleTypeTaker, TradeRoleTypeMaker:
		return true
	}
	return false
}

func (t TradeRoleType) String() string {
	return string(t)
}

func ParseTradeRoleType(s string) (TradeRoleType, error) {
	roleType := TradeRoleType(s)
	roleType = roleType.Normalize()

	if !roleType.IsValid() {
		return roleType, fmt.Errorf("%s is not a valid TradeRoleType", s)
	}

	return roleType, nil
}

func (t TradeRoleType) Normalize() TradeRoleType {
	return TradeRoleType(strings.ToUpper(string(t)))
}

type PositionMode string

const (
	PositionModeOneWay PositionMode = "ONE_WAY"
	PositionModeHedge  PositionMode = "HEDGE"
)

func (p PositionMode) IsValid() bool {
	switch p {
	case PositionModeOneWay, PositionModeHedge:
		return true
	}
	return false
}

func (p PositionMode) String() string {
	return string(p)
}

func ParsePositionMode(s string) (PositionMode, error) {
	mode := PositionMode(s)
	mode = mode.Normalize()

	if !mode.IsValid() {
		return mode, fmt.Errorf("%s is not a valid PositionMode", s)
	}

	return mode, nil
}

func (p PositionMode) Normalize() PositionMode {
	return PositionMode(strings.ToUpper(string(p)))
}

const ChannelBalance = "balance"
const ChannelPosition = "position"
const ChannelOrder = "order"
const ChannelTpSl = "tpsl"

type OrderStatus string

const (
	OrderStatusInit           OrderStatus = "INIT"
	OrderStatusNew            OrderStatus = "NEW"
	OrderStatusPartFilled     OrderStatus = "PART_FILLED"
	OrderStatusCanceled       OrderStatus = "CANCELED"
	OrderStatusSystemCanceled OrderStatus = "SYSTEM_CANCELED"
	OrderStatusExpired        OrderStatus = "EXPIRED"
	OrderStatusFilled         OrderStatus = "FILLED"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusInit, OrderStatusNew, OrderStatusPartFilled, OrderStatusCanceled, OrderStatusSystemCanceled, OrderStatusExpired, OrderStatusFilled:
		return true
	}
	return false
}

func (s OrderStatus) String() string {
	return string(s)
}

func ParseOrderStatus(s string) (OrderStatus, error) {
	status := OrderStatus(s)
	status = status.Normalize()

	if !status.IsValid() {
		return status, fmt.Errorf("%s is not a valid OrderStatus", s)
	}

	return status, nil
}

func (s OrderStatus) Normalize() OrderStatus {
	return OrderStatus(strings.ToUpper(string(s)))
}

type PositionSide string

const (
	PositionSideShort PositionSide = "SHORT"
	PositionSideLong  PositionSide = "LONG"
)

func (s PositionSide) IsValid() bool {
	switch s {
	case PositionSideShort, PositionSideLong:
		return true
	}
	return false
}

func (s PositionSide) String() string {
	return string(s)
}

func ParsePositionSide(s string) (PositionSide, error) {
	status := PositionSide(s)
	status = status.Normalize()

	if !status.IsValid() {
		return status, fmt.Errorf("%s is not a valid PositionSide", s)
	}

	return status, nil
}

func (s PositionSide) Normalize() PositionSide {
	return PositionSide(strings.ToUpper(string(s)))
}

type PositionEventType string

const (
	PositionEventOpen   PositionEventType = "OPEN"
	PositionEventUpdate PositionEventType = "UPDATE"
	PositionEventClose  PositionEventType = "CLOSE"
)

func (s PositionEventType) IsValid() bool {
	switch s {
	case PositionEventOpen, PositionEventUpdate, PositionEventClose:
		return true
	}
	return false
}

func (s PositionEventType) String() string {
	return string(s)
}

func ParsePositionEvent(s string) (PositionEventType, error) {
	status := PositionEventType(s)
	status = status.Normalize()

	if !status.IsValid() {
		return status, fmt.Errorf("%s is not a valid PositionEventType", s)
	}

	return status, nil
}

func (s PositionEventType) Normalize() PositionEventType {
	return PositionEventType(strings.ToUpper(string(s)))
}

type OrderEventType string

const (
	OrderEventCreate OrderEventType = "CREATE"
	OrderEventUpdate OrderEventType = "UPDATE"
	OrderEventClose  OrderEventType = "CLOSE"
)

func (s OrderEventType) IsValid() bool {
	switch s {
	case OrderEventCreate, OrderEventUpdate, OrderEventClose:
		return true
	}
	return false
}

func (s OrderEventType) String() string {
	return string(s)
}

func ParseOrderEvent(s string) (OrderEventType, error) {
	status := OrderEventType(s)
	status = status.Normalize()

	if !status.IsValid() {
		return status, fmt.Errorf("%s is not a valid PositionEventType", s)
	}

	return status, nil
}

func (s OrderEventType) Normalize() OrderEventType {
	return OrderEventType(strings.ToUpper(string(s)))
}

type TpSlEventType string

const (
	TPSLEventCreate TpSlEventType = "CREATE"
	TPSLEventUpdate TpSlEventType = "UPDATE"
	TPSLEventClose  TpSlEventType = "CLOSE"
)

func (s TpSlEventType) IsValid() bool {
	switch s {
	case TPSLEventCreate, TPSLEventUpdate, TPSLEventClose:
		return true
	}
	return false
}

func (s TpSlEventType) String() string {
	return string(s)
}

func ParseTPSLEvent(s string) (TpSlEventType, error) {
	status := TpSlEventType(s)
	status = status.Normalize()

	if !status.IsValid() {
		return status, fmt.Errorf("%s is not a valid TpSlEventType", s)
	}

	return status, nil
}

func (s TpSlEventType) Normalize() TpSlEventType {
	return TpSlEventType(strings.ToUpper(string(s)))
}

type TpSlType string

const (
	TPSLTypeFull    TpSlType = "POSITION_TPSL"
	TPSLTypePartial TpSlType = "TPSL"
)

func (s TpSlType) IsValid() bool {
	switch s {
	case TPSLTypeFull, TPSLTypePartial:
		return true
	}
	return false
}

func (s TpSlType) String() string {
	return string(s)
}

func ParseTPSLType(s string) (TpSlType, error) {
	status := TpSlType(s)
	status = status.Normalize()

	if !status.IsValid() {
		return status, fmt.Errorf("%s is not a valid TpSlType", s)
	}

	return status, nil
}

func (s TpSlType) Normalize() TpSlType {
	return TpSlType(strings.ToUpper(string(s)))
}

type Symbol string

func (s Symbol) String() string {
	return string(s)
}

func ParseSymbol(s string) Symbol {
	return Symbol(s).Normalize()
}

func (s Symbol) Normalize() Symbol {
	return Symbol(strings.ToUpper(strings.TrimSpace(string(s))))
}

type MarginCoin string

func (s MarginCoin) String() string {
	return string(s)
}

func ParseMarginCoin(s string) MarginCoin {
	return MarginCoin(s).Normalize()
}

func (s MarginCoin) Normalize() MarginCoin {
	return MarginCoin(strings.ToUpper(strings.TrimSpace(string(s))))
}

const (
	Interval1Min   Interval = "1min"
	Interval3Min   Interval = "3min"
	Interval5Min   Interval = "5min"
	Interval15Min  Interval = "15min"
	Interval30Min  Interval = "30min"
	Interval60Min  Interval = "60min"
	Interval2H     Interval = "2h"
	Interval4H     Interval = "4h"
	Interval6H     Interval = "6h"
	Interval8H     Interval = "8h"
	Interval12H    Interval = "12h"
	Interval1Day   Interval = "1day"
	Interval3Day   Interval = "3day"
	Interval1Week  Interval = "1week"
	Interval1Month Interval = "1month"
)

type Interval string

func (s Interval) String() string {
	return string(s)
}

func ParseInterval(s string) (Interval, error) {
	interval := Interval(s).Normalize()

	if !interval.IsValid() {
		return interval, fmt.Errorf("%s is not a valid interval", s)
	}

	return interval, nil
}

func (s Interval) IsValid() bool {
	switch s {
	case Interval1Min, Interval3Min, Interval5Min, Interval15Min, Interval30Min, Interval60Min,
		Interval2H, Interval4H, Interval6H, Interval8H, Interval12H,
		Interval1Day, Interval3Day,
		Interval1Week,
		Interval1Month:
		return true
	}
	return false
}

func (s Interval) Normalize() Interval {
	return Interval(strings.ToLower(strings.TrimSpace(string(s))))
}

const PriceTypeMark PriceType = "mark"
const PriceTypeMarket PriceType = "market"

type PriceType string

func (s PriceType) String() string {
	return string(s)
}

func ParsePriceType(s string) (PriceType, error) {
	priceType := PriceType(s).Normalize()

	if !priceType.IsValid() {
		return priceType, fmt.Errorf("%s is not a valid pricetype", s)
	}

	return priceType, nil
}

func (s PriceType) IsValid() bool {
	switch s {
	case PriceTypeMark, PriceTypeMarket:
		return true
	}

	return false
}

func (s PriceType) Normalize() PriceType {
	return PriceType(strings.ToLower(strings.TrimSpace(string(s))))
}

const ChannelKline Channel = "kline"

type Channel string

func (s Channel) String() string {
	return string(s)
}

func ParseChannel(s string) (Channel, error) {
	Channel := Channel(s).Normalize()

	if !Channel.IsValid() {
		return Channel, fmt.Errorf("%s is not a valid Channel", s)
	}

	return Channel, nil
}

func (s Channel) IsValid() bool {
	switch s {
	case ChannelKline:
		return true
	}

	return false
}

func (s Channel) Normalize() Channel {
	return Channel(strings.ToLower(strings.TrimSpace(string(s))))
}
