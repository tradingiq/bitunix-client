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

	if PositionSideShort != "SHORT" {
		t.Errorf("Expected PositionSideShort to be 'SHORT', got '%s'", PositionSideShort)
	}

	if PositionSideLong != "LONG" {
		t.Errorf("Expected PositionSideLong to be 'LONG', got '%s'", PositionSideLong)
	}

	if PositionEventOpen != "OPEN" {
		t.Errorf("Expected PositionEventOpen to be 'OPEN', got '%s'", PositionEventOpen)
	}

	if PositionEventUpdate != "UPDATE" {
		t.Errorf("Expected PositionEventUpdate to be 'UPDATE', got '%s'", PositionEventUpdate)
	}

	if PositionEventClose != "CLOSE" {
		t.Errorf("Expected PositionEventClose to be 'CLOSE', got '%s'", PositionEventClose)
	}

	if OrderEventCreate != "CREATE" {
		t.Errorf("Expected OrderEventCreate to be 'CREATE', got '%s'", OrderEventCreate)
	}

	if OrderEventUpdate != "UPDATE" {
		t.Errorf("Expected OrderEventUpdate to be 'UPDATE', got '%s'", OrderEventUpdate)
	}

	if OrderEventClose != "CLOSE" {
		t.Errorf("Expected OrderEventClose to be 'CLOSE', got '%s'", OrderEventClose)
	}

	if TPSLEventCreate != "CREATE" {
		t.Errorf("Expected TPSLEventCreate to be 'CREATE', got '%s'", TPSLEventCreate)
	}

	if TPSLEventUpdate != "UPDATE" {
		t.Errorf("Expected TPSLEventUpdate to be 'UPDATE', got '%s'", TPSLEventUpdate)
	}

	if TPSLEventClose != "CLOSE" {
		t.Errorf("Expected TPSLEventClose to be 'CLOSE', got '%s'", TPSLEventClose)
	}

	if TPSLTypeFull != "POSITION_TPSL" {
		t.Errorf("Expected TPSLTypeFull to be 'POSITION_TPSL', got '%s'", TPSLTypeFull)
	}

	if TPSLTypePartial != "TPSL" {
		t.Errorf("Expected TPSLTypePartial to be 'TPSL', got '%s'", TPSLTypePartial)
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

func TestPositionSideNormalize(t *testing.T) {
	tests := []struct {
		input    PositionSide
		expected PositionSide
	}{
		{PositionSide("SHORT"), PositionSideShort},
		{PositionSide("short"), PositionSideShort},
		{PositionSide("Short"), PositionSideShort},
		{PositionSide("LONG"), PositionSideLong},
		{PositionSide("long"), PositionSideLong},
		{PositionSide("Long"), PositionSideLong},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("PositionSide.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestPositionSideParse(t *testing.T) {
	tests := []struct {
		input       string
		expected    PositionSide
		expectError bool
	}{
		{"SHORT", PositionSideShort, false},
		{"short", PositionSideShort, false},
		{"Short", PositionSideShort, false},
		{"LONG", PositionSideLong, false},
		{"long", PositionSideLong, false},
		{"Long", PositionSideLong, false},
		{"INVALID", PositionSide("INVALID"), true},
	}

	for _, test := range tests {
		result, err := ParsePositionSide(test.input)
		if test.expectError && err == nil {
			t.Errorf("ParsePositionSide(%s): expected error, got nil", test.input)
		} else if !test.expectError && err != nil {
			t.Errorf("ParsePositionSide(%s): unexpected error: %v", test.input, err)
		} else if !test.expectError && result != test.expected {
			t.Errorf("ParsePositionSide(%s): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

func TestPositionEventTypeNormalize(t *testing.T) {
	tests := []struct {
		input    PositionEventType
		expected PositionEventType
	}{
		{PositionEventType("OPEN"), PositionEventOpen},
		{PositionEventType("open"), PositionEventOpen},
		{PositionEventType("Open"), PositionEventOpen},
		{PositionEventType("UPDATE"), PositionEventUpdate},
		{PositionEventType("update"), PositionEventUpdate},
		{PositionEventType("Update"), PositionEventUpdate},
		{PositionEventType("CLOSE"), PositionEventClose},
		{PositionEventType("close"), PositionEventClose},
		{PositionEventType("Close"), PositionEventClose},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("PositionEventType.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestPositionEventTypeParse(t *testing.T) {
	tests := []struct {
		input       string
		expected    PositionEventType
		expectError bool
	}{
		{"OPEN", PositionEventOpen, false},
		{"open", PositionEventOpen, false},
		{"Open", PositionEventOpen, false},
		{"UPDATE", PositionEventUpdate, false},
		{"update", PositionEventUpdate, false},
		{"Update", PositionEventUpdate, false},
		{"CLOSE", PositionEventClose, false},
		{"close", PositionEventClose, false},
		{"Close", PositionEventClose, false},
		{"INVALID", PositionEventType("INVALID"), true},
	}

	for _, test := range tests {
		result, err := ParsePositionEvent(test.input)
		if test.expectError && err == nil {
			t.Errorf("ParsePositionEvent(%s): expected error, got nil", test.input)
		} else if !test.expectError && err != nil {
			t.Errorf("ParsePositionEvent(%s): unexpected error: %v", test.input, err)
		} else if !test.expectError && result != test.expected {
			t.Errorf("ParsePositionEvent(%s): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

func TestOrderEventTypeNormalize(t *testing.T) {
	tests := []struct {
		input    OrderEventType
		expected OrderEventType
	}{
		{OrderEventType("CREATE"), OrderEventCreate},
		{OrderEventType("create"), OrderEventCreate},
		{OrderEventType("Create"), OrderEventCreate},
		{OrderEventType("UPDATE"), OrderEventUpdate},
		{OrderEventType("update"), OrderEventUpdate},
		{OrderEventType("Update"), OrderEventUpdate},
		{OrderEventType("CLOSE"), OrderEventClose},
		{OrderEventType("close"), OrderEventClose},
		{OrderEventType("Close"), OrderEventClose},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("OrderEventType.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestOrderEventTypeParse(t *testing.T) {
	tests := []struct {
		input       string
		expected    OrderEventType
		expectError bool
	}{
		{"CREATE", OrderEventCreate, false},
		{"create", OrderEventCreate, false},
		{"Create", OrderEventCreate, false},
		{"UPDATE", OrderEventUpdate, false},
		{"update", OrderEventUpdate, false},
		{"Update", OrderEventUpdate, false},
		{"CLOSE", OrderEventClose, false},
		{"close", OrderEventClose, false},
		{"Close", OrderEventClose, false},
		{"INVALID", OrderEventType("INVALID"), true},
	}

	for _, test := range tests {
		result, err := ParseOrderEvent(test.input)
		if test.expectError && err == nil {
			t.Errorf("ParseOrderEvent(%s): expected error, got nil", test.input)
		} else if !test.expectError && err != nil {
			t.Errorf("ParseOrderEvent(%s): unexpected error: %v", test.input, err)
		} else if !test.expectError && result != test.expected {
			t.Errorf("ParseOrderEvent(%s): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

func TestTpSlEventTypeNormalize(t *testing.T) {
	tests := []struct {
		input    TpSlEventType
		expected TpSlEventType
	}{
		{TpSlEventType("CREATE"), TPSLEventCreate},
		{TpSlEventType("create"), TPSLEventCreate},
		{TpSlEventType("Create"), TPSLEventCreate},
		{TpSlEventType("UPDATE"), TPSLEventUpdate},
		{TpSlEventType("update"), TPSLEventUpdate},
		{TpSlEventType("Update"), TPSLEventUpdate},
		{TpSlEventType("CLOSE"), TPSLEventClose},
		{TpSlEventType("close"), TPSLEventClose},
		{TpSlEventType("Close"), TPSLEventClose},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("TpSlEventType.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestTpSlEventTypeParse(t *testing.T) {
	tests := []struct {
		input       string
		expected    TpSlEventType
		expectError bool
	}{
		{"CREATE", TPSLEventCreate, false},
		{"create", TPSLEventCreate, false},
		{"Create", TPSLEventCreate, false},
		{"UPDATE", TPSLEventUpdate, false},
		{"update", TPSLEventUpdate, false},
		{"Update", TPSLEventUpdate, false},
		{"CLOSE", TPSLEventClose, false},
		{"close", TPSLEventClose, false},
		{"Close", TPSLEventClose, false},
		{"INVALID", TpSlEventType("INVALID"), true},
	}

	for _, test := range tests {
		result, err := ParseTPSLEvent(test.input)
		if test.expectError && err == nil {
			t.Errorf("ParseTPSLEvent(%s): expected error, got nil", test.input)
		} else if !test.expectError && err != nil {
			t.Errorf("ParseTPSLEvent(%s): unexpected error: %v", test.input, err)
		} else if !test.expectError && result != test.expected {
			t.Errorf("ParseTPSLEvent(%s): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

func TestTpSlTypeNormalize(t *testing.T) {
	tests := []struct {
		input    TpSlType
		expected TpSlType
	}{
		{TpSlType("POSITION_TPSL"), TPSLTypeFull},
		{TpSlType("position_tpsl"), TPSLTypeFull},
		{TpSlType("Position_Tpsl"), TPSLTypeFull},
		{TpSlType("TPSL"), TPSLTypePartial},
		{TpSlType("tpsl"), TPSLTypePartial},
		{TpSlType("Tpsl"), TPSLTypePartial},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("TpSlType.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestTpSlTypeParse(t *testing.T) {
	tests := []struct {
		input       string
		expected    TpSlType
		expectError bool
	}{
		{"POSITION_TPSL", TPSLTypeFull, false},
		{"position_tpsl", TPSLTypeFull, false},
		{"Position_Tpsl", TPSLTypeFull, false},
		{"TPSL", TPSLTypePartial, false},
		{"tpsl", TPSLTypePartial, false},
		{"Tpsl", TPSLTypePartial, false},
		{"INVALID", TpSlType("INVALID"), true},
	}

	for _, test := range tests {
		result, err := ParseTPSLType(test.input)
		if test.expectError && err == nil {
			t.Errorf("ParseTPSLType(%s): expected error, got nil", test.input)
		} else if !test.expectError && err != nil {
			t.Errorf("ParseTPSLType(%s): unexpected error: %v", test.input, err)
		} else if !test.expectError && result != test.expected {
			t.Errorf("ParseTPSLType(%s): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

func TestSymbolNormalize(t *testing.T) {
	tests := []struct {
		input    Symbol
		expected Symbol
	}{
		{Symbol("BTCUSDT"), Symbol("BTCUSDT")},
		{Symbol("btcusdt"), Symbol("BTCUSDT")},
		{Symbol("BtCuSdT"), Symbol("BTCUSDT")},
		{Symbol(" BTCUSDT "), Symbol("BTCUSDT")},
		{Symbol("  btcusdt  "), Symbol("BTCUSDT")},
	}

	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("Symbol.Normalize() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestParseSymbol(t *testing.T) {
	tests := []struct {
		input    string
		expected Symbol
	}{
		{"BTCUSDT", Symbol("BTCUSDT")},
		{"btcusdt", Symbol("BTCUSDT")},
		{"BtCuSdT", Symbol("BTCUSDT")},
		{" BTCUSDT ", Symbol("BTCUSDT")},
		{"  btcusdt  ", Symbol("BTCUSDT")},
	}

	for _, test := range tests {
		result := ParseSymbol(test.input)
		if result != test.expected {
			t.Errorf("ParseSymbol(%s): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

func TestSymbolString(t *testing.T) {
	tests := []struct {
		input    Symbol
		expected string
	}{
		{Symbol("BTCUSDT"), "BTCUSDT"},
	}

	for _, test := range tests {
		result := test.input.String()
		if result != test.expected {
			t.Errorf("Symbol.String() with input '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestIntervalNormalize(t *testing.T) {
	testCases := []struct {
		input    string
		expected Interval
	}{
		// Standard format
		{"1min", Interval1Min},
		{"3min", Interval3Min},
		{"5min", Interval5Min},
		{"15min", Interval15Min},
		{"30min", Interval30Min},
		{"60min", Interval60Min},
		{"2h", Interval2H},
		{"4h", Interval4H},
		{"6h", Interval6H},
		{"8h", Interval8H},
		{"12h", Interval12H},
		{"1day", Interval1Day},
		{"3day", Interval3Day},
		{"1week", Interval1Week},
		{"1month", Interval1Month},

		// Alternative formats
		{"1m", Interval1Min},
		{"3m", Interval3Min},
		{"5m", Interval5Min},
		{"15m", Interval15Min},
		{"30m", Interval30Min},
		{"60m", Interval60Min},
		{"1h", Interval60Min},
		{"1d", Interval1Day},
		{"3d", Interval3Day},
		{"1w", Interval1Week},
		{"1mo", Interval1Month},

		// With whitespace
		{" 1min ", Interval1Min},
		{" 1m ", Interval1Min},
		{" 1day ", Interval1Day},
		{"  2h  ", Interval2H},

		// Mixed case
		{"1Min", Interval1Min},
		{"1DAY", Interval1Day},
		{"1WeEk", Interval1Week},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := Interval(tc.input).Normalize()
			if result != tc.expected {
				t.Errorf("Interval.Normalize() with input '%s': expected '%s', got '%s'",
					tc.input, tc.expected, result)
			}

			// Also test the Parse function
			parsedResult, err := ParseInterval(tc.input)
			if err != nil {
				t.Errorf("ParseInterval(%s) unexpected error: %v", tc.input, err)
			}
			if parsedResult != tc.expected {
				t.Errorf("ParseInterval(%s): expected %s, got %s", tc.input, tc.expected, parsedResult)
			}
		})
	}
}

func TestPriceTypeNormalize(t *testing.T) {
	testCases := []struct {
		input    string
		expected PriceType
	}{
		// Standard format
		{"mark", PriceTypeMark},
		{"market", PriceTypeMarket},

		// Alternative formats
		{"markprice", PriceTypeMark},
		{"mark_price", PriceTypeMark},
		{"marketprice", PriceTypeMarket},
		{"market_price", PriceTypeMarket},

		// With whitespace
		{" mark ", PriceTypeMark},
		{" market ", PriceTypeMarket},
		{"  markprice  ", PriceTypeMark},

		// Mixed case
		{"Mark", PriceTypeMark},
		{"MARKET", PriceTypeMarket},
		{"MarkPrice", PriceTypeMark},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := PriceType(tc.input).Normalize()
			if result != tc.expected {
				t.Errorf("PriceType.Normalize() with input '%s': expected '%s', got '%s'",
					tc.input, tc.expected, result)
			}

			// Also test the Parse function
			parsedResult, err := ParsePriceType(tc.input)
			if err != nil {
				t.Errorf("ParsePriceType(%s) unexpected error: %v", tc.input, err)
			}
			if parsedResult != tc.expected {
				t.Errorf("ParsePriceType(%s): expected %s, got %s", tc.input, tc.expected, parsedResult)
			}
		})
	}
}

func TestChannelNormalize(t *testing.T) {
	testCases := []struct {
		input    string
		expected Channel
	}{
		// Standard format
		{"kline", ChannelKline},

		// Alternative formats
		{"k", ChannelKline},
		{"candle", ChannelKline},
		{"candlestick", ChannelKline},

		// With whitespace
		{" kline ", ChannelKline},
		{" k ", ChannelKline},
		{"  candle  ", ChannelKline},

		// Mixed case
		{"Kline", ChannelKline},
		{"CANDLE", ChannelKline},
		{"CandleStick", ChannelKline},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := Channel(tc.input).Normalize()
			if result != tc.expected {
				t.Errorf("Channel.Normalize() with input '%s': expected '%s', got '%s'",
					tc.input, tc.expected, result)
			}

			// Also test the Parse function
			parsedResult, err := ParseChannel(tc.input)
			if err != nil {
				t.Errorf("ParseChannel(%s) unexpected error: %v", tc.input, err)
			}
			if parsedResult != tc.expected {
				t.Errorf("ParseChannel(%s): expected %s, got %s", tc.input, tc.expected, parsedResult)
			}
		})
	}
}

// Test for whitespace handling across all types
func TestWhitespaceNormalization(t *testing.T) {
	t.Run("StopType", func(t *testing.T) {
		input := " LAST_PRICE "
		expected := StopTypeLastPrice
		result := StopType(input).Normalize()
		if result != expected {
			t.Errorf("StopType.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("OrderType", func(t *testing.T) {
		input := " LIMIT "
		expected := OrderTypeLimit
		result := OrderType(input).Normalize()
		if result != expected {
			t.Errorf("OrderType.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("TimeInForce", func(t *testing.T) {
		input := " GTC "
		expected := TimeInForceGTC
		result := TimeInForce(input).Normalize()
		if result != expected {
			t.Errorf("TimeInForce.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("Side", func(t *testing.T) {
		input := " OPEN "
		expected := SideOpen
		result := Side(input).Normalize()
		if result != expected {
			t.Errorf("Side.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("TradeSide", func(t *testing.T) {
		input := " BUY "
		expected := TradeSideBuy
		result := TradeSide(input).Normalize()
		if result != expected {
			t.Errorf("TradeSide.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("MarginMode", func(t *testing.T) {
		input := " ISOLATION "
		expected := MarginModeIsolation
		result := MarginMode(input).Normalize()
		if result != expected {
			t.Errorf("MarginMode.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("TradeRoleType", func(t *testing.T) {
		input := " TAKER "
		expected := TradeRoleTypeTaker
		result := TradeRoleType(input).Normalize()
		if result != expected {
			t.Errorf("TradeRoleType.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("PositionMode", func(t *testing.T) {
		input := " ONE_WAY "
		expected := PositionModeOneWay
		result := PositionMode(input).Normalize()
		if result != expected {
			t.Errorf("PositionMode.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("OrderStatus", func(t *testing.T) {
		input := " NEW "
		expected := OrderStatusNew
		result := OrderStatus(input).Normalize()
		if result != expected {
			t.Errorf("OrderStatus.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("PositionSide", func(t *testing.T) {
		input := " LONG "
		expected := PositionSideLong
		result := PositionSide(input).Normalize()
		if result != expected {
			t.Errorf("PositionSide.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("PositionEventType", func(t *testing.T) {
		input := " OPEN "
		expected := PositionEventOpen
		result := PositionEventType(input).Normalize()
		if result != expected {
			t.Errorf("PositionEventType.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("OrderEventType", func(t *testing.T) {
		input := " CREATE "
		expected := OrderEventCreate
		result := OrderEventType(input).Normalize()
		if result != expected {
			t.Errorf("OrderEventType.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("TpSlEventType", func(t *testing.T) {
		input := " CREATE "
		expected := TPSLEventCreate
		result := TpSlEventType(input).Normalize()
		if result != expected {
			t.Errorf("TpSlEventType.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})

	t.Run("TpSlType", func(t *testing.T) {
		input := " POSITION_TPSL "
		expected := TPSLTypeFull
		result := TpSlType(input).Normalize()
		if result != expected {
			t.Errorf("TpSlType.Normalize() with input '%s': expected '%s', got '%s'",
				input, expected, result)
		}
	})
}
