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

	switch s {
	case "Last":
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
