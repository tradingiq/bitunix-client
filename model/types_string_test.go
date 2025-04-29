package model

import (
	"testing"
)

func TestTypeString(t *testing.T) {

	if StopTypeLastPrice.String() != "LAST_PRICE" {
		t.Errorf("Expected StopTypeLastPrice.String() to return 'LAST_PRICE', got '%s'", StopTypeLastPrice.String())
	}
	if StopTypeMarkPrice.String() != "MARK_PRICE" {
		t.Errorf("Expected StopTypeMarkPrice.String() to return 'MARK_PRICE', got '%s'", StopTypeMarkPrice.String())
	}

	if OrderTypeLimit.String() != "LIMIT" {
		t.Errorf("Expected OrderTypeLimit.String() to return 'LIMIT', got '%s'", OrderTypeLimit.String())
	}
	if OrderTypeMarket.String() != "MARKET" {
		t.Errorf("Expected OrderTypeMarket.String() to return 'MARKET', got '%s'", OrderTypeMarket.String())
	}

	if TimeInForceIOC.String() != "IOC" {
		t.Errorf("Expected TimeInForceIOC.String() to return 'IOC', got '%s'", TimeInForceIOC.String())
	}
	if TimeInForceFOK.String() != "FOK" {
		t.Errorf("Expected TimeInForceFOK.String() to return 'FOK', got '%s'", TimeInForceFOK.String())
	}
	if TimeInForceGTC.String() != "GTC" {
		t.Errorf("Expected TimeInForceGTC.String() to return 'GTC', got '%s'", TimeInForceGTC.String())
	}
	if TimeInForcePostOnly.String() != "POST_ONLY" {
		t.Errorf("Expected TimeInForcePostOnly.String() to return 'POST_ONLY', got '%s'", TimeInForcePostOnly.String())
	}

	if SideOpen.String() != "OPEN" {
		t.Errorf("Expected SideOpen.String() to return 'OPEN', got '%s'", SideOpen.String())
	}
	if SideClose.String() != "CLOSE" {
		t.Errorf("Expected SideClose.String() to return 'CLOSE', got '%s'", SideClose.String())
	}

	if TradeSideBuy.String() != "BUY" {
		t.Errorf("Expected TradeSideBuy.String() to return 'BUY', got '%s'", TradeSideBuy.String())
	}
	if TradeSideSell.String() != "SELL" {
		t.Errorf("Expected TradeSideSell.String() to return 'SELL', got '%s'", TradeSideSell.String())
	}

	if MarginModeIsolation.String() != "ISOLATION" {
		t.Errorf("Expected MarginModeIsolation.String() to return 'ISOLATION', got '%s'", MarginModeIsolation.String())
	}
	if MarginModeCross.String() != "CROSS" {
		t.Errorf("Expected MarginModeCross.String() to return 'CROSS', got '%s'", MarginModeCross.String())
	}

	if TradeRoleTypeTaker.String() != "TAKER" {
		t.Errorf("Expected TradeRoleTypeTaker.String() to return 'TAKER', got '%s'", TradeRoleTypeTaker.String())
	}
	if TradeRoleTypeMaker.String() != "MAKER" {
		t.Errorf("Expected TradeRoleTypeMaker.String() to return 'MAKER', got '%s'", TradeRoleTypeMaker.String())
	}

	if PositionModeOneWay.String() != "ONE_WAY" {
		t.Errorf("Expected PositionModeOneWay.String() to return 'ONE_WAY', got '%s'", PositionModeOneWay.String())
	}
	if PositionModeHedge.String() != "HEDGE" {
		t.Errorf("Expected PositionModeHedge.String() to return 'HEDGE', got '%s'", PositionModeHedge.String())
	}

	if OrderStatusInit.String() != "INIT" {
		t.Errorf("Expected OrderStatusInit.String() to return 'INIT', got '%s'", OrderStatusInit.String())
	}
	if OrderStatusNew.String() != "NEW" {
		t.Errorf("Expected OrderStatusNew.String() to return 'NEW', got '%s'", OrderStatusNew.String())
	}
	if OrderStatusPartFilled.String() != "PART_FILLED" {
		t.Errorf("Expected OrderStatusPartFilled.String() to return 'PART_FILLED', got '%s'", OrderStatusPartFilled.String())
	}
	if OrderStatusCanceled.String() != "CANCELED" {
		t.Errorf("Expected OrderStatusCanceled.String() to return 'CANCELED', got '%s'", OrderStatusCanceled.String())
	}
	if OrderStatusSystemCanceled.String() != "SYSTEM_CANCELED" {
		t.Errorf("Expected OrderStatusSystemCanceled.String() to return 'SYSTEM_CANCELED', got '%s'", OrderStatusSystemCanceled.String())
	}
	if OrderStatusExpired.String() != "EXPIRED" {
		t.Errorf("Expected OrderStatusExpired.String() to return 'EXPIRED', got '%s'", OrderStatusExpired.String())
	}
	if OrderStatusFilled.String() != "FILLED" {
		t.Errorf("Expected OrderStatusFilled.String() to return 'FILLED', got '%s'", OrderStatusFilled.String())
	}

	if PositionSideShort.String() != "SHORT" {
		t.Errorf("Expected PositionSideShort.String() to return 'SHORT', got '%s'", PositionSideShort.String())
	}
	if PositionSideLong.String() != "LONG" {
		t.Errorf("Expected PositionSideLong.String() to return 'LONG', got '%s'", PositionSideLong.String())
	}

	if PositionEventOpen.String() != "OPEN" {
		t.Errorf("Expected PositionEventOpen.String() to return 'OPEN', got '%s'", PositionEventOpen.String())
	}
	if PositionEventUpdate.String() != "UPDATE" {
		t.Errorf("Expected PositionEventUpdate.String() to return 'UPDATE', got '%s'", PositionEventUpdate.String())
	}
	if PositionEventClose.String() != "CLOSE" {
		t.Errorf("Expected PositionEventClose.String() to return 'CLOSE', got '%s'", PositionEventClose.String())
	}

	if OrderEventCreate.String() != "CREATE" {
		t.Errorf("Expected OrderEventCreate.String() to return 'CREATE', got '%s'", OrderEventCreate.String())
	}
	if OrderEventUpdate.String() != "UPDATE" {
		t.Errorf("Expected OrderEventUpdate.String() to return 'UPDATE', got '%s'", OrderEventUpdate.String())
	}
	if OrderEventClose.String() != "CLOSE" {
		t.Errorf("Expected OrderEventClose.String() to return 'CLOSE', got '%s'", OrderEventClose.String())
	}

	if TPSLEventCreate.String() != "CREATE" {
		t.Errorf("Expected TPSLEventCreate.String() to return 'CREATE', got '%s'", TPSLEventCreate.String())
	}
	if TPSLEventUpdate.String() != "UPDATE" {
		t.Errorf("Expected TPSLEventUpdate.String() to return 'UPDATE', got '%s'", TPSLEventUpdate.String())
	}
	if TPSLEventClose.String() != "CLOSE" {
		t.Errorf("Expected TPSLEventClose.String() to return 'CLOSE', got '%s'", TPSLEventClose.String())
	}

	if TPSLTypeFull.String() != "POSITION_TPSL" {
		t.Errorf("Expected TPSLTypeFull.String() to return 'POSITION_TPSL', got '%s'", TPSLTypeFull.String())
	}
	if TPSLTypePartial.String() != "TPSL" {
		t.Errorf("Expected TPSLTypePartial.String() to return 'TPSL', got '%s'", TPSLTypePartial.String())
	}
}
