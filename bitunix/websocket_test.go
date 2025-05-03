package bitunix

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/websocket"
	"testing"
	"time"
)

func TestHeartbeatMessage(t *testing.T) {
	timestamp := time.Now().Unix()
	message := heartbeatMessage{
		Op:   "ping",
		Ping: timestamp,
	}

	bytes, err := json.Marshal(message)
	require.NoError(t, err)

	var decoded heartbeatMessage
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "ping", decoded.Op)
	assert.Equal(t, timestamp, decoded.Ping)
}

func TestKLineUnmarshal(t *testing.T) {
	jsonData := `{
		"ch": "mark_kline_1min",
		"symbol": "BTCUSDT",
		"ts": 1732178884994,
		"data": {
			"o": "0.0010",
			"c": "0.0020",
			"h": "0.0025",
			"l": "0.0015",
			"b": "1.01",
			"q": "1.09"
		}
	}`

	var msg model.KLineChannelMessage
	err := json.Unmarshal([]byte(jsonData), &msg)
	assert.NoError(t, err)
	assert.Equal(t, "mark_kline_1min", msg.Channel)
	assert.Equal(t, model.ParseSymbol("BTCUSDT"), msg.Symbol)
	assert.Equal(t, int64(1732178884994), msg.Ts)
	assert.Equal(t, 0.0010, msg.Data.OpenPrice)
	assert.Equal(t, 0.0020, msg.Data.ClosePrice)
	assert.Equal(t, 0.0025, msg.Data.HighPrice)
	assert.Equal(t, 0.0015, msg.Data.LowPrice)
	assert.Equal(t, 1.01, msg.Data.BaseVolume)
	assert.Equal(t, 1.09, msg.Data.QuoteVolume)
}

type mockWsClient struct {
	connectFn func() error
	listenFn  func(callback websocket.HandlerFunc) error
	closeFn   func()
	writeFn   func([]byte) error
}

func (m *mockWsClient) Connect() error {
	if m.connectFn != nil {
		return m.connectFn()
	}
	return nil
}

func (m *mockWsClient) Listen(callback websocket.HandlerFunc) error {
	if m.listenFn != nil {
		return m.listenFn(callback)
	}
	return nil
}

func (m *mockWsClient) Close() {
	if m.closeFn != nil {
		m.closeFn()
	}
}

func (m *mockWsClient) Write(bytes []byte) error {
	if m.writeFn != nil {
		return m.writeFn(bytes)
	}
	return nil
}

func TestPublicWebsocketClient_SubscribeKLine(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &publicWebsocketClient{
		websocketClient: websocketClient{
			client: mockWs,
			uri:    "wss://test.com",
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}

	var sentBytes []byte
	mockWs.writeFn = func(bytes []byte) error {
		sentBytes = bytes
		return nil
	}

	sub := &subTest{ch: make(chan struct{}, 1)}
	err := client.SubscribeKLine(sub)
	assert.NoError(t, err)

	var req SubscribeRequest
	err = json.Unmarshal(sentBytes, &req)
	assert.NoError(t, err)
	assert.Equal(t, "subscribe", req.Op)

	_, exists := client.klineHandlers[sub]
	assert.True(t, exists)
}

type subTest struct {
	called bool
	msg    model.KLineChannelMessage
	ch     chan struct{}
}

func (s *subTest) Handle(message *model.KLineChannelMessage) {
	s.called = true
	s.msg = *message
	// Signal that the message has been processed
	if s.ch != nil {
		s.ch <- struct{}{}
	}
}

func (s *subTest) Interval() model.Interval {
	return model.Interval1Min
}

func (s *subTest) Symbol() model.Symbol {
	return model.Symbol("BTCUSDT")
}

func (s *subTest) PriceType() model.PriceType {
	return model.PriceTypeMarket
}

func TestPublicWebsocketClient_Stream_KLine(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &publicWebsocketClient{
		websocketClient: websocketClient{
			client: mockWs,
			uri:    "wss://test.com",
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}

	var listenCallback websocket.HandlerFunc
	mockWs.listenFn = func(callback websocket.HandlerFunc) error {
		listenCallback = callback
		return nil
	}

	// Create notification channel
	ch := make(chan struct{}, 1)
	sub := &subTest{ch: ch}
	err := client.SubscribeKLine(sub)
	assert.NoError(t, err)

	err = client.Stream()
	assert.NoError(t, err)

	klineData := `{
		"ch": "market_kline_1min",
		"symbol": "BTCUSDT",
		"ts": 1732178884994,
		"data": {
			"o": "0.0010",
			"c": "0.0020",
			"h": "0.0025",
			"l": "0.0015",
			"b": "1.01",
			"q": "1.09"
		}
	}`

	err = listenCallback([]byte(klineData))
	assert.NoError(t, err)

	// Wait for processing to complete with timeout
	select {
	case <-ch:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message to be processed")
	}

	// Access the data after synchronization
	klineMsg := sub.msg

	assert.True(t, sub.called)
	assert.Equal(t, "market_kline_1min", klineMsg.Channel)
	assert.Equal(t, model.ParseSymbol("BTCUSDT"), klineMsg.Symbol)
	assert.Equal(t, int64(1732178884994), klineMsg.Ts)
	assert.Equal(t, 0.0010, klineMsg.Data.OpenPrice)
}
