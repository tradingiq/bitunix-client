package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type BalanceChannelMessage struct {
	Ch   string       `json:"ch"`
	Ts   int64        `json:"ts"`
	Data BalanceEvent `json:"data"`
}

type BalanceEvent struct {
	Coin            string  `json:"coin"`
	Available       float64 `json:"-"`
	Frozen          float64 `json:"-"`
	IsolationFrozen float64 `json:"-"`
	CrossFrozen     float64 `json:"-"`
	Margin          float64 `json:"-"`
	IsolationMargin float64 `json:"-"`
	CrossMargin     float64 `json:"-"`
	ExpMoney        float64 `json:"-"`
}

func (p *BalanceEvent) UnmarshalJSON(data []byte) error {
	type Alias BalanceEvent
	aux := &struct {
		Available       string `json:"available"`
		Frozen          string `json:"frozen"`
		IsolationFrozen string `json:"isolationFrozen"`
		CrossFrozen     string `json:"crossFrozen"`
		Margin          string `json:"margin"`
		IsolationMargin string `json:"isolationMargin"`
		CrossMargin     string `json:"crossMargin"`
		ExpMoney        string `json:"expMoney"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	available, err := strconv.ParseFloat(aux.Available, 64)
	if err == nil {
		p.Available = available
	} else {
		return fmt.Errorf("failed to parse available: %w", err)
	}

	frozen, err := strconv.ParseFloat(aux.Frozen, 64)
	if err == nil {
		p.Frozen = frozen
	} else {
		return fmt.Errorf("failed to parse frozen: %w", err)
	}

	isolationFrozen, err := strconv.ParseFloat(aux.IsolationFrozen, 64)
	if err == nil {
		p.IsolationFrozen = isolationFrozen
	} else {
		return fmt.Errorf("failed to parse IisolationFrozen: %w", err)
	}

	crossFrozen, err := strconv.ParseFloat(aux.CrossFrozen, 64)
	if err == nil {
		p.CrossFrozen = crossFrozen
	} else {
		return fmt.Errorf("failed to parse crossFrozen: %w", err)
	}

	Margin, err := strconv.ParseFloat(aux.Margin, 64)
	if err == nil {
		p.Margin = Margin
	} else {
		return fmt.Errorf("failed to parse Margin: %w", err)
	}

	isolationMargin, err := strconv.ParseFloat(aux.IsolationMargin, 64)
	if err == nil {
		p.IsolationMargin = isolationMargin
	} else {
		return fmt.Errorf("failed to parse isolationMargin: %w", err)
	}

	crossMargin, err := strconv.ParseFloat(aux.CrossMargin, 64)
	if err == nil {
		p.CrossMargin = crossMargin
	} else {
		return fmt.Errorf("failed to parse crossMargin: %w", err)
	}

	expMoney, err := strconv.ParseFloat(aux.ExpMoney, 64)
	if err == nil {
		p.ExpMoney = expMoney
	} else {
		return fmt.Errorf("failed to parse expMoney: %w", err)
	}

	return nil
}

type PositionEvent struct {
	Event         PositionEventType `json:"-"`
	PositionID    string            `json:"positionId"`
	MarginMode    MarginMode        `json:"-"`
	PositionMode  PositionMode      `json:"-"`
	Side          PositionSide      `json:"-"`
	Leverage      int               `json:"-"`
	Margin        float64           `json:"-"`
	CreateTime    time.Time         `json:"-"`
	Quantity      float64           `json:"-"`
	EntryValue    float64           `json:"-"`
	Symbol        Symbol            `json:"symbol"`
	RealizedPNL   float64           `json:"-"`
	UnrealizedPNL float64           `json:"-"`
	Funding       float64           `json:"-"`
	Fee           float64           `json:"-"`
}

func (p *PositionEvent) UnmarshalJSON(data []byte) error {
	type Alias PositionEvent
	aux := &struct {
		Event         string `json:"event"`
		MarginMode    string `json:"marginMode"`
		PositionMode  string `json:"positionMode"`
		Side          string `json:"side"`
		Leverage      string `json:"leverage"`
		Margin        string `json:"margin"`
		CreateTime    string `json:"ctime"`
		Quantity      string `json:"qty"`
		EntryValue    string `json:"entryValue"`
		RealizedPNL   string `json:"realizedPNL"`
		UnrealizedPNL string `json:"unrealizedPNL"`
		Funding       string `json:"funding"`
		Fee           string `json:"fee"`
		Symbol        string `json:"symbol"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	p.Symbol = ParseSymbol(aux.Symbol)

	event, err := ParsePositionEvent(aux.Event)
	if err != nil {
		return fmt.Errorf("invalid position event: %w", err)
	}
	p.Event = event

	side, err := ParsePositionSide(aux.Side)
	if err != nil {
		return fmt.Errorf("invalid side: %w", err)
	}
	p.Side = side

	posMode, err := ParsePositionMode(aux.PositionMode)
	if err != nil {
		return fmt.Errorf("invalid position mode: %w", err)
	}
	p.PositionMode = posMode

	marginMode, err := ParseMarginMode(aux.MarginMode)
	if err != nil {
		return fmt.Errorf("invalid margin mode: %w", err)
	}
	p.MarginMode = marginMode

	if aux.CreateTime != "" {
		t, err := time.Parse(time.RFC3339Nano, aux.CreateTime)
		if err == nil {
			p.CreateTime = t
		} else {
			return fmt.Errorf("invalid create time: %w", err)
		}
	}

	if aux.Leverage != "" {
		val, err := strconv.Atoi(aux.Leverage)
		if err != nil {
			return fmt.Errorf("failed to parse leverage: %w", err)
		}
		p.Leverage = val
	}

	if aux.Margin != "" {
		margin, err := strconv.ParseFloat(aux.Margin, 64)
		if err == nil {
			p.Margin = margin
		} else {
			return fmt.Errorf("failed to parse margin: %w", err)
		}
	}

	if aux.Quantity != "" {
		qty, err := strconv.ParseFloat(aux.Quantity, 64)
		if err == nil {
			p.Quantity = qty
		} else {
			return fmt.Errorf("failed to parse qty: %w", err)
		}
	}

	if aux.EntryValue != "" {
		val, err := strconv.ParseFloat(aux.EntryValue, 64)
		if err == nil {
			p.EntryValue = val
		} else {
			return fmt.Errorf("failed to parse EntryValue: %w", err)
		}
	}

	if aux.RealizedPNL != "" {
		val, err := strconv.ParseFloat(aux.RealizedPNL, 64)
		if err == nil {
			p.RealizedPNL = val
		} else {
			return fmt.Errorf("failed to parse RealizedPNL: %w", err)
		}
	}

	if aux.UnrealizedPNL != "" {
		val, err := strconv.ParseFloat(aux.UnrealizedPNL, 64)
		if err == nil {
			p.UnrealizedPNL = val
		} else {
			return fmt.Errorf("failed to parse UnrealizedPNL: %w", err)
		}
	}

	if aux.Funding != "" {
		val, err := strconv.ParseFloat(aux.Funding, 64)
		if err == nil {
			p.Funding = val
		} else {
			return fmt.Errorf("failed to parse Funding: %w", err)
		}
	}

	if aux.Fee != "" {
		val, err := strconv.ParseFloat(aux.Fee, 64)
		if err == nil {
			p.Fee = val
		} else {
			return fmt.Errorf("failed to parse Fee: %w", err)
		}
	}

	return nil
}

type PositionChannelMessage struct {
	Channel   string        `json:"ch"`
	TimeStamp int64         `json:"ts"`
	Data      PositionEvent `json:"data"`
}

type OrderEvent struct {
	Event        OrderEventType `json:"-"`
	OrderID      string         `json:"orderId"`
	Symbol       Symbol         `json:"symbol"`
	PositionType MarginMode     `json:"-"`
	PositionMode PositionMode   `json:"-"`
	Side         TradeSide      `json:"-"`
	/*
		Currently broken, returns "SHORT"/"LONG"
		Effect        TimeInForce    `json:"effect"`
	*/
	Type          OrderType   `json:"-"`
	Quantity      float64     `json:"-"`
	ReductionOnly bool        `json:"reductionOnly"`
	Price         float64     `json:"-"`
	CreateTime    time.Time   `json:"-"`
	ModifyTime    time.Time   `json:"-"`
	Leverage      int         `json:"-"`
	OrderStatus   OrderStatus `json:"-"`
	Fee           float64     `json:"-"`
	/*
		Currently not consistently propagated, use the TPSL channel instead
		TPStopType    StopType    `json:"-"`
		TPPrice       float64     `json:"-"`
		TPOrderType   OrderType   `json:"-"`
		TPOrderPrice  float64     `json:"-"`
		SLStopType    StopType    `json:"-"`
		SLPrice float64 `json:"-"`
		SLOrderType OrderType `json:"-"`
		SLOrderPrice float64 `json:"-"`
	*/
}

func (o *OrderEvent) UnmarshalJSON(data []byte) error {
	type Alias OrderEvent
	aux := &struct {
		Event        string `json:"event"`
		PositionType string `json:"positionType"`
		PositionMode string `json:"positionMode"`
		Side         string `json:"side"`
		Type         string `json:"type"`
		Quantity     string `json:"qty"`
		Price        string `json:"price"`
		CreateTime   string `json:"ctime"`
		ModifyTime   string `json:"mtime"`
		Leverage     string `json:"leverage"`
		OrderStatus  string `json:"orderStatus"`
		Fee          string `json:"fee"`
		Symbol       string `json:"symbol"`
		/*
			TPStopType   string `json:"tpStopType,omitempty"`
			TPPrice      string `json:"tpPrice,omitempty"`
			TPOrderType  string `json:"tpOrderType,omitempty"`
			TPOrderPrice string `json:"tpOrderPrice,omitempty"`
			SLStopType   string `json:"slStopType,omitempty"`
			SLPrice      string `json:"slPrice,omitempty"`
			SLOrderType  string `json:"slOrderType,omitempty"`
			SLOrderPrice string `json:"slOrderPrice,omitempty"`
		*/
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	o.Symbol = ParseSymbol(aux.Symbol)

	event, err := ParseOrderEvent(aux.Event)
	if err != nil {
		return fmt.Errorf("invalid order event: %w", err)
	}
	o.Event = event

	posType, err := ParseMarginMode(aux.PositionType)
	if err != nil {
		return fmt.Errorf("invalid position type: %w", err)
	}
	o.PositionType = posType

	posMode, err := ParsePositionMode(aux.PositionMode)
	if err != nil {
		return fmt.Errorf("invalid position mode: %w", err)
	}
	o.PositionMode = posMode

	side, err := ParseTradeSide(aux.Side)
	if err != nil {
		return fmt.Errorf("invalid side: %w", err)
	}
	o.Side = side

	orderType, err := ParseOrderType(aux.Type)
	if err != nil {
		return fmt.Errorf("invalid order type: %w", err)
	}
	o.Type = orderType

	status, err := ParseOrderStatus(aux.OrderStatus)
	if err != nil {
		return fmt.Errorf("invalid order status: %w", err)
	}
	o.OrderStatus = status

	if aux.CreateTime != "" {
		t, err := time.Parse(time.RFC3339Nano, aux.CreateTime)
		if err == nil {
			o.CreateTime = t
		} else {
			return fmt.Errorf("invalid create time: %w", err)
		}
	}

	if aux.ModifyTime != "" {
		t, err := time.Parse(time.RFC3339Nano, aux.ModifyTime)
		if err == nil {
			o.ModifyTime = t
		} else {
			return fmt.Errorf("invalid modify time: %w", err)
		}
	}

	if aux.Quantity != "" {
		qty, err := strconv.ParseFloat(aux.Quantity, 64)
		if err == nil {
			o.Quantity = qty
		} else {
			return fmt.Errorf("failed to parse quantity: %w", err)
		}
	}

	if aux.Price != "" {
		price, err := strconv.ParseFloat(aux.Price, 64)
		if err == nil {
			o.Price = price
		} else {
			return fmt.Errorf("failed to parse price: %w", err)
		}
	}

	if aux.Leverage != "" {
		lev, err := strconv.Atoi(aux.Leverage)
		if err == nil {
			o.Leverage = lev
		} else {
			return fmt.Errorf("failed to parse leverage: %w", err)
		}
	}

	if aux.Fee != "" {
		fee, err := strconv.ParseFloat(aux.Fee, 64)
		if err == nil {
			o.Fee = fee
		} else {
			return fmt.Errorf("failed to parse fee: %w", err)
		}
	}

	return nil
}

type OrderChannelMessage struct {
	Channel   string     `json:"ch"`
	TimeStamp int64      `json:"ts"`
	Data      OrderEvent `json:"data"`
}

type TpSlOrderChannelMessage struct {
	Channel   string         `json:"ch"`
	Timestamp int64          `json:"ts"`
	Data      TpSlOrderEvent `json:"data"`
}

type KLineChannelMessage struct {
	Channel string      `json:"ch"`
	Symbol  string      `json:"symbol"`
	Ts      int64       `json:"ts"`
	Data    KLineEvent  `json:"data"`
}

type KLineEvent struct {
	OpenPrice   float64 `json:"-"`
	HighPrice   float64 `json:"-"`
	LowPrice    float64 `json:"-"`
	ClosePrice  float64 `json:"-"`
	BaseVolume  float64 `json:"-"`
	QuoteVolume float64 `json:"-"`
}

func (k *KLineEvent) UnmarshalJSON(data []byte) error {
	type Alias KLineEvent
	aux := &struct {
		OpenPrice   string `json:"o"`
		HighPrice   string `json:"h"`
		LowPrice    string `json:"l"`
		ClosePrice  string `json:"c"`
		BaseVolume  string `json:"b"`
		QuoteVolume string `json:"q"`
		*Alias
	}{
		Alias: (*Alias)(k),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	openPrice, err := strconv.ParseFloat(aux.OpenPrice, 64)
	if err == nil {
		k.OpenPrice = openPrice
	} else {
		return fmt.Errorf("failed to parse open price: %w", err)
	}

	highPrice, err := strconv.ParseFloat(aux.HighPrice, 64)
	if err == nil {
		k.HighPrice = highPrice
	} else {
		return fmt.Errorf("failed to parse high price: %w", err)
	}

	lowPrice, err := strconv.ParseFloat(aux.LowPrice, 64)
	if err == nil {
		k.LowPrice = lowPrice
	} else {
		return fmt.Errorf("failed to parse low price: %w", err)
	}

	closePrice, err := strconv.ParseFloat(aux.ClosePrice, 64)
	if err == nil {
		k.ClosePrice = closePrice
	} else {
		return fmt.Errorf("failed to parse close price: %w", err)
	}

	baseVolume, err := strconv.ParseFloat(aux.BaseVolume, 64)
	if err == nil {
		k.BaseVolume = baseVolume
	} else {
		return fmt.Errorf("failed to parse base volume: %w", err)
	}

	quoteVolume, err := strconv.ParseFloat(aux.QuoteVolume, 64)
	if err == nil {
		k.QuoteVolume = quoteVolume
	} else {
		return fmt.Errorf("failed to parse quote volume: %w", err)
	}

	return nil
}

type TpSlOrderEvent struct {
	Event        TpSlEventType `json:"-"`
	PositionID   string        `json:"positionId"`
	OrderID      string        `json:"orderId"`
	Symbol       Symbol        `json:"symbol"`
	Leverage     int           `json:"-"`
	Side         TradeSide     `json:"-"`
	PositionMode PositionMode  `json:"-"`
	Status       OrderStatus   `json:"-"`
	CreateTime   time.Time     `json:"-"`
	Type         TpSlType      `json:"-"`
	TPQuantity   float64       `json:"-"`
	SLQuantity   float64       `json:"-"`
	TPStopType   StopType      `json:"-"`
	TPPrice      float64       `json:"-"`
	TPOrderType  OrderType     `json:"-"`
	TPOrderPrice float64       `json:"-"`
	SLStopType   StopType      `json:"-"`
	SLPrice      float64       `json:"-"`
	SLOrderType  OrderType     `json:"-"`
	SLOrderPrice float64       `json:"-"`
}

func (t *TpSlOrderEvent) UnmarshalJSON(data []byte) error {
	type Alias TpSlOrderEvent
	aux := &struct {
		Event        string `json:"event"`
		Leverage     string `json:"leverage"`
		Side         string `json:"side"`
		PositionMode string `json:"positionMode"`
		Status       string `json:"status"`
		CreateTime   string `json:"ctime"`
		Type         string `json:"type"`
		TPQuantity   string `json:"tpQty"`
		SLQuantity   string `json:"slQty"`
		TPStopType   string `json:"tpStopType"`
		TPPrice      string `json:"tpPrice"`
		TPOrderType  string `json:"tpOrderType"`
		TPOrderPrice string `json:"tpOrderPrice"`
		SLStopType   string `json:"slStopType"`
		SLPrice      string `json:"slPrice"`
		SLOrderType  string `json:"slOrderType"`
		SLOrderPrice string `json:"slOrderPrice"`
		Symbol       string `json:"symbol"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Symbol = ParseSymbol(aux.Symbol)

	event, err := ParseTPSLEvent(aux.Event)
	if err != nil {
		return fmt.Errorf("invalid TPSL event: %w", err)
	}
	t.Event = event

	side, err := ParseTradeSide(aux.Side)
	if err != nil {
		return fmt.Errorf("invalid side: %w", err)
	}
	t.Side = side

	posMode, err := ParsePositionMode(aux.PositionMode)
	if err != nil {
		return fmt.Errorf("invalid position mode: %w", err)
	}
	t.PositionMode = posMode

	status, err := ParseOrderStatus(aux.Status)
	if err != nil {
		return fmt.Errorf("invalid order status: %w", err)
	}
	t.Status = status

	tpslType, err := ParseTPSLType(aux.Type)
	if err != nil {
		return fmt.Errorf("invalid TPSL type: %w", err)
	}
	t.Type = tpslType

	if aux.CreateTime != "" {
		timestamp, err := time.Parse(time.RFC3339Nano, aux.CreateTime)
		if err == nil {
			t.CreateTime = timestamp
		} else {
			return fmt.Errorf("invalid create time: %w", err)
		}
	}

	if aux.Leverage != "" {
		lev, err := strconv.Atoi(aux.Leverage)
		if err == nil {
			t.Leverage = lev
		} else {
			return fmt.Errorf("failed to parse leverage: %w", err)
		}
	}

	if aux.TPQuantity != "" {
		qty, err := strconv.ParseFloat(aux.TPQuantity, 64)
		if err == nil {
			t.TPQuantity = qty
		} else {
			return fmt.Errorf("failed to parse TP quantity: %w", err)
		}
	}

	if aux.SLQuantity != "" {
		qty, err := strconv.ParseFloat(aux.SLQuantity, 64)
		if err == nil {
			t.SLQuantity = qty
		} else {
			return fmt.Errorf("failed to parse SL quantity: %w", err)
		}
	}

	if aux.TPStopType != "" {
		tpStopType, err := ParseStopType(aux.TPStopType)
		if err == nil {
			t.TPStopType = tpStopType
		} else {
			return fmt.Errorf("failed to parse TP stop type: %w", err)
		}
	}

	if aux.TPPrice != "" {
		tpPrice, err := strconv.ParseFloat(aux.TPPrice, 64)
		if err == nil {
			t.TPPrice = tpPrice
		} else {
			return fmt.Errorf("failed to parse TP price: %w", err)
		}
	}

	if aux.TPOrderType != "" {
		tpOrderType, err := ParseOrderType(aux.TPOrderType)
		if err == nil {
			t.TPOrderType = tpOrderType
		} else {
			return fmt.Errorf("failed to parse TP order type: %w", err)
		}
	}

	if aux.TPOrderPrice != "" {
		tpOrderPrice, err := strconv.ParseFloat(aux.TPOrderPrice, 64)
		if err == nil {
			t.TPOrderPrice = tpOrderPrice
		} else {
			return fmt.Errorf("failed to parse TP order price: %w", err)
		}
	}

	if aux.SLStopType != "" {
		slStopType, err := ParseStopType(aux.SLStopType)
		if err == nil {
			t.SLStopType = slStopType
		} else {
			return fmt.Errorf("failed to parse SL stop type: %w", err)
		}
	}

	if aux.SLPrice != "" {
		slPrice, err := strconv.ParseFloat(aux.SLPrice, 64)
		if err == nil {
			t.SLPrice = slPrice
		} else {
			return fmt.Errorf("failed to parse SL price: %w", err)
		}
	}

	if aux.SLOrderType != "" {
		slOrderType, err := ParseOrderType(aux.SLOrderType)
		if err == nil {
			t.SLOrderType = slOrderType
		} else {
			return fmt.Errorf("failed to parse SL order type: %w", err)
		}
	}

	if aux.SLOrderPrice != "" {
		slOrderPrice, err := strconv.ParseFloat(aux.SLOrderPrice, 64)
		if err == nil {
			t.SLOrderPrice = slOrderPrice
		} else {
			return fmt.Errorf("failed to parse SL order price: %w", err)
		}
	}

	return nil
}
