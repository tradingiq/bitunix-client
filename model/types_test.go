package model

import (
	"testing"
)

func TestTypeConstants(t *testing.T) {
	if StopTypeLastPrice != "LAST_PRICE" {
		t.Errorf("Expected StopTypeLastPrice to be 'LAST_PRICE', got '%s'", StopTypeLastPrice)
	}

	if StopTypeMarkPrice != "MARK_PRICE" {
		t.Errorf("Expected StopTypeMarkPrice to be 'MARK_PRICE', got '%s'", StopTypeMarkPrice)
	}

	if OrderTypeLimit != "LIMIT" {
		t.Errorf("Expected OrderTypeLimit to be 'LIMIT', got '%s'", OrderTypeLimit)
	}

	if OrderTypeMarket != "MARKET" {
		t.Errorf("Expected OrderTypeMarket to be 'MARKET', got '%s'", OrderTypeMarket)
	}

	if TimeInForceIOC != "IOC" {
		t.Errorf("Expected TimeInForceIOC to be 'IOC', got '%s'", TimeInForceIOC)
	}

	if TimeInForceFOK != "FOK" {
		t.Errorf("Expected TimeInForceFOK to be 'FOK', got '%s'", TimeInForceFOK)
	}

	if TimeInForceGTC != "GTC" {
		t.Errorf("Expected TimeInForceGTC to be 'GTC', got '%s'", TimeInForceGTC)
	}

	if TimeInForcePostOnly != "POST_ONLY" {
		t.Errorf("Expected TimeInForcePostOnly to be 'POST_ONLY', got '%s'", TimeInForcePostOnly)
	}

	if SideOpen != "OPEN" {
		t.Errorf("Expected SideOpen to be 'OPEN', got '%s'", SideOpen)
	}

	if SideClose != "CLOSE" {
		t.Errorf("Expected SideClose to be 'CLOSE', got '%s'", SideClose)
	}

	if TradeSideBuy != "BUY" {
		t.Errorf("Expected TradeSideBuy to be 'BUY', got '%s'", TradeSideBuy)
	}

	if TradeSideSell != "SELL" {
		t.Errorf("Expected TradeSideSell to be 'SELL', got '%s'", TradeSideSell)
	}

	if MarginModeIsolation != "ISOLATION" {
		t.Errorf("Expected MarginModeIsolation to be 'ISOLATION', got '%s'", MarginModeIsolation)
	}

	if MarginModeCross != "CROSS" {
		t.Errorf("Expected MarginModeCross to be 'CROSS', got '%s'", MarginModeCross)
	}

	if PositionModeOneWay != "ONE_WAY" {
		t.Errorf("Expected PositionModeOneWay to be 'ONE_WAY', got '%s'", PositionModeOneWay)
	}

	if PositionModeHedge != "HEDGE" {
		t.Errorf("Expected PositionModeHedge to be 'HEDGE', got '%s'", PositionModeHedge)
	}

	if TradeRoleTypeTaker != "TAKER" {
		t.Errorf("Expected TradeRoleTypeTaker to be 'TAKER', got '%s'", TradeRoleTypeTaker)
	}

	if TradeRoleTypeMaker != "MAKER" {
		t.Errorf("Expected TradeRoleTypeMaker to be 'MAKER', got '%s'", TradeRoleTypeMaker)
	}
}

func TestStopTypeNormalize(t *testing.T) {
	tests := []struct {
		input    StopType
		expected StopType
	}{
		{StopType("LAST_PRICE"), StopTypeLastPrice},
		{StopType("last_price"), StopTypeLastPrice},
		{StopType("MARK_PRICE"), StopTypeMarkPrice},
		{StopType("mark_price"), StopTypeMarkPrice},

		{StopType("LAST"), StopTypeLastPrice},
		{StopType("last"), StopTypeLastPrice},
		{StopType("Last"), StopTypeLastPrice},
		{StopType("MARK"), StopTypeMarkPrice},
		{StopType("mark"), StopTypeMarkPrice},
		{StopType("Mark"), StopTypeMarkPrice},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("StopType.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestStopTypeParse(t *testing.T) {
	tests := []struct {
		input       string
		expected    StopType
		expectError bool
	}{
		{"LAST_PRICE", StopTypeLastPrice, false},
		{"last_price", StopTypeLastPrice, false},
		{"MARK_PRICE", StopTypeMarkPrice, false},
		{"mark_price", StopTypeMarkPrice, false},

		{"LAST", StopTypeLastPrice, false},
		{"last", StopTypeLastPrice, false},
		{"Last", StopTypeLastPrice, false},
		{"MARK", StopTypeMarkPrice, false},
		{"mark", StopTypeMarkPrice, false},
		{"Mark", StopTypeMarkPrice, false},

		{"INVALID", StopType("INVALID"), true},
	}

	for _, test := range tests {
		result, err := ParseStopType(test.input)
		if test.expectError && err == nil {
			t.Errorf("ParseStopType(%s): expected error, got nil", test.input)
		} else if !test.expectError && err != nil {
			t.Errorf("ParseStopType(%s): unexpected error: %v", test.input, err)
		} else if !test.expectError && result != test.expected {
			t.Errorf("ParseStopType(%s): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

func TestOrderTypeNormalize(t *testing.T) {
	tests := []struct {
		input    OrderType
		expected OrderType
	}{
		{OrderType("LIMIT"), OrderTypeLimit},
		{OrderType("limit"), OrderTypeLimit},
		{OrderType("Limit"), OrderTypeLimit},
		{OrderType("MARKET"), OrderTypeMarket},
		{OrderType("market"), OrderTypeMarket},
		{OrderType("Market"), OrderTypeMarket},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("OrderType.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestOrderTypeParse(t *testing.T) {
	tests := []struct {
		input       string
		expected    OrderType
		expectError bool
	}{
		{"LIMIT", OrderTypeLimit, false},
		{"limit", OrderTypeLimit, false},
		{"Limit", OrderTypeLimit, false},
		{"MARKET", OrderTypeMarket, false},
		{"market", OrderTypeMarket, false},
		{"Market", OrderTypeMarket, false},

		{"INVALID", OrderType("INVALID"), true},
	}

	for _, test := range tests {
		result, err := ParseOrderType(test.input)
		if test.expectError && err == nil {
			t.Errorf("ParseOrderType(%s): expected error, got nil", test.input)
		} else if !test.expectError && err != nil {
			t.Errorf("ParseOrderType(%s): unexpected error: %v", test.input, err)
		} else if !test.expectError && result != test.expected {
			t.Errorf("ParseOrderType(%s): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

func TestTimeInForceNormalize(t *testing.T) {
	tests := []struct {
		input    TimeInForce
		expected TimeInForce
	}{
		{TimeInForce("IOC"), TimeInForceIOC},
		{TimeInForce("ioc"), TimeInForceIOC},
		{TimeInForce("FOK"), TimeInForceFOK},
		{TimeInForce("fok"), TimeInForceFOK},
		{TimeInForce("GTC"), TimeInForceGTC},
		{TimeInForce("gtc"), TimeInForceGTC},
		{TimeInForce("POST_ONLY"), TimeInForcePostOnly},
		{TimeInForce("post_only"), TimeInForcePostOnly},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("TimeInForce.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestSideNormalize(t *testing.T) {
	tests := []struct {
		input    Side
		expected Side
	}{
		{Side("OPEN"), SideOpen},
		{Side("open"), SideOpen},
		{Side("Open"), SideOpen},
		{Side("CLOSE"), SideClose},
		{Side("close"), SideClose},
		{Side("Close"), SideClose},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("Side.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestTradeSideNormalize(t *testing.T) {
	tests := []struct {
		input    TradeSide
		expected TradeSide
	}{
		{TradeSide("BUY"), TradeSideBuy},
		{TradeSide("buy"), TradeSideBuy},
		{TradeSide("Buy"), TradeSideBuy},
		{TradeSide("SELL"), TradeSideSell},
		{TradeSide("sell"), TradeSideSell},
		{TradeSide("Sell"), TradeSideSell},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("TradeSide.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestMarginModeNormalize(t *testing.T) {
	tests := []struct {
		input    MarginMode
		expected MarginMode
	}{
		{MarginMode("ISOLATION"), MarginModeIsolation},
		{MarginMode("isolation"), MarginModeIsolation},
		{MarginMode("Isolation"), MarginModeIsolation},
		{MarginMode("CROSS"), MarginModeCross},
		{MarginMode("cross"), MarginModeCross},
		{MarginMode("Cross"), MarginModeCross},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("MarginMode.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestPositionModeNormalize(t *testing.T) {
	tests := []struct {
		input    PositionMode
		expected PositionMode
	}{
		{PositionMode("ONE_WAY"), PositionModeOneWay},
		{PositionMode("one_way"), PositionModeOneWay},
		{PositionMode("One_Way"), PositionModeOneWay},
		{PositionMode("HEDGE"), PositionModeHedge},
		{PositionMode("hedge"), PositionModeHedge},
		{PositionMode("Hedge"), PositionModeHedge},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("PositionMode.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestTradeRoleTypeNormalize(t *testing.T) {
	tests := []struct {
		input    TradeRoleType
		expected TradeRoleType
	}{
		{TradeRoleType("TAKER"), TradeRoleTypeTaker},
		{TradeRoleType("taker"), TradeRoleTypeTaker},
		{TradeRoleType("Taker"), TradeRoleTypeTaker},
		{TradeRoleType("MAKER"), TradeRoleTypeMaker},
		{TradeRoleType("maker"), TradeRoleTypeMaker},
		{TradeRoleType("Maker"), TradeRoleTypeMaker},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("TradeRoleType.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestOrderStatusNormalize(t *testing.T) {
	tests := []struct {
		input    OrderStatus
		expected OrderStatus
	}{
		{OrderStatus("INIT"), OrderStatusInit},
		{OrderStatus("init"), OrderStatusInit},
		{OrderStatus("NEW"), OrderStatusNew},
		{OrderStatus("new"), OrderStatusNew},
		{OrderStatus("PART_FILLED"), OrderStatusPartFilled},
		{OrderStatus("part_filled"), OrderStatusPartFilled},
		{OrderStatus("CANCELED"), OrderStatusCanceled},
		{OrderStatus("canceled"), OrderStatusCanceled},
		{OrderStatus("SYSTEM_CANCELED"), OrderStatusSystemCanceled},
		{OrderStatus("system_canceled"), OrderStatusSystemCanceled},
		{OrderStatus("EXPIRED"), OrderStatusExpired},
		{OrderStatus("expired"), OrderStatusExpired},
		{OrderStatus("FILLED"), OrderStatusFilled},
		{OrderStatus("filled"), OrderStatusFilled},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("OrderStatus.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}
