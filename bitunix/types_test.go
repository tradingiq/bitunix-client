package bitunix

import (
	"testing"
)

func TestTypes(t *testing.T) {
	// Test StopType constants
	if StopTypeLastPrice != "LAST_PRICE" {
		t.Errorf("Expected StopTypeLastPrice to be 'LAST_PRICE', got '%s'", StopTypeLastPrice)
	}

	if StopTypeMarkPrice != "MARK_PRICE" {
		t.Errorf("Expected StopTypeMarkPrice to be 'MARK_PRICE', got '%s'", StopTypeMarkPrice)
	}

	// Test OrderType constants
	if OrderTypeLimit != "LIMIT" {
		t.Errorf("Expected OrderTypeLimit to be 'LIMIT', got '%s'", OrderTypeLimit)
	}

	if OrderTypeMarket != "MARKET" {
		t.Errorf("Expected OrderTypeMarket to be 'MARKET', got '%s'", OrderTypeMarket)
	}

	// Test TimeInForce constants
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

	// Test TradeSide constants
	if TradeSideOpen != "OPEN" {
		t.Errorf("Expected TradeSideOpen to be 'OPEN', got '%s'", TradeSideOpen)
	}

	if TradeSideClose != "CLOSE" {
		t.Errorf("Expected TradeSideClose to be 'CLOSE', got '%s'", TradeSideClose)
	}

	// Test TradeAction constants
	if TradeActionBuy != "BUY" {
		t.Errorf("Expected TradeActionBuy to be 'BUY', got '%s'", TradeActionBuy)
	}

	if TradeActionSell != "SELL" {
		t.Errorf("Expected TradeActionSell to be 'SELL', got '%s'", TradeActionSell)
	}

	// Test MarginMode constants
	if MarginModeIsolation != "ISOLATION" {
		t.Errorf("Expected MarginModeIsolation to be 'ISOLATION ', got '%s'", MarginModeIsolation)
	}

	if MarginModeCross != "CROSS" {
		t.Errorf("Expected MarginModeCross to be 'CROSS', got '%s'", MarginModeCross)
	}

	// Test TradePositionMode constants
	if TradePositionModeOneWay != "ONE_WAY" {
		t.Errorf("Expected TradePositionModeOneWay to be 'ONE_WAY', got '%s'", TradePositionModeOneWay)
	}

	if TradePositionModeHedge != "HEDGE" {
		t.Errorf("Expected TradePositionModeHedge to be 'HEDGE', got '%s'", TradePositionModeHedge)
	}

	// Test TradeRoleType constants
	if TradeRoleTypeTaker != "TAKER" {
		t.Errorf("Expected TradeRoleTypeTaker to be 'TAKER', got '%s'", TradeRoleTypeTaker)
	}

	if TradeRoleTypeMaker != "MAKER" {
		t.Errorf("Expected TradeRoleTypeMaker to be 'MAKER', got '%s'", TradeRoleTypeMaker)
	}
}
