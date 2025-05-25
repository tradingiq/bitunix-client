package bitunix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/websocket"
	"sync"
	"sync/atomic"
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
		websocketClient: &websocketClient{
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

func TestPublicWebsocketClient_SubscribeKLine_NilSubscriber(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client: mockWs,
			uri:    "wss://test.com",
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}

	err := client.SubscribeKLine(nil)
	assert.Error(t, err, "Should return error when subscriber is nil")
}

func TestPublicWebsocketClient_SubscribeKLine_WriteError(t *testing.T) {
	mockWs := &mockWsClient{}
	mockErr := errors.New("write error")

	mockWs.writeFn = func(bytes []byte) error {
		return mockErr
	}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client: mockWs,
			uri:    "wss://test.com",
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}

	sub := &subTest{ch: make(chan struct{}, 1)}
	err := client.SubscribeKLine(sub)
	assert.Error(t, err, "Should return error when write fails")

	_, exists := client.klineHandlers[sub]
	assert.True(t, exists, "Handler should be added before write is attempted")
}

type subTest struct {
	called      bool
	msg         model.KLineChannelMessage
	ch          chan struct{}
	subscribeMu sync.Mutex
	symbol      *model.Symbol
}

func (s *subTest) SubscribeKLine(message *model.KLineChannelMessage) {
	s.subscribeMu.Lock()
	s.called = true
	s.msg = *message
	s.subscribeMu.Unlock()

	if s.ch != nil {
		s.ch <- struct{}{}
	}
}

func (s *subTest) SubscribeInterval() model.Interval {
	return model.Interval1Min
}

func (s *subTest) SubscribeSymbol() model.Symbol {
	if s.symbol == nil {
		return model.Symbol("BTCUSDT")
	}

	return *s.symbol
}

func (s *subTest) SubscribePriceType() model.PriceType {
	return model.PriceTypeMarket
}

func TestPublicWebsocketClient_UnsubscribeKLine(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
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
	client.klineHandlers[sub] = struct{}{}

	err := client.UnsubscribeKLine(sub)
	assert.NoError(t, err)

	var req SubscribeRequest
	err = json.Unmarshal(sentBytes, &req)
	assert.NoError(t, err)
	assert.Equal(t, "unsubscribe", req.Op)

	_, exists := client.klineHandlers[sub]
	assert.False(t, exists, "Handler should be removed after unsubscribe")
}

func TestPublicWebsocketClient_UnsubscribeKLine_NilSubscriber(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client: mockWs,
			uri:    "wss://test.com",
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}

	err := client.UnsubscribeKLine(nil)
	assert.Error(t, err, "Should return error when subscriber is nil")
}

func TestPublicWebsocketClient_UnsubscribeKLine_WriteError(t *testing.T) {
	mockWs := &mockWsClient{}
	mockErr := errors.New("write error")

	mockWs.writeFn = func(bytes []byte) error {
		return mockErr
	}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client: mockWs,
			uri:    "wss://test.com",
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}

	sub := &subTest{ch: make(chan struct{}, 1)}
	client.klineHandlers[sub] = struct{}{}

	err := client.UnsubscribeKLine(sub)
	assert.Error(t, err, "Should return error when write fails")

	_, exists := client.klineHandlers[sub]
	assert.False(t, exists, "Handler should be removed even if unsubscribe message fails")
}

func TestWebsocketClient_Connect_Error(t *testing.T) {
	mockWs := &mockWsClient{}
	mockErr := errors.New("connection error")

	mockWs.connectFn = func() error {
		return mockErr
	}

	client := &websocketClient{
		client: mockWs,
		uri:    "wss://test.com",
	}

	err := client.Connect()
	assert.Error(t, err, "Should return error when connect fails")
}

func TestWebsocketClient_Stream_Error(t *testing.T) {
	mockWs := &mockWsClient{}
	mockErr := errors.New("listen error")

	mockWs.listenFn = func(callback websocket.HandlerFunc) error {
		return mockErr
	}

	client := &websocketClient{
		client:       mockWs,
		uri:          "wss://test.com",
		messageQueue: make(chan []byte, 10),
	}

	err := client.Stream()
	assert.Error(t, err, "Should return error when listen fails")
}

func TestWebsocketClient_Stream_WorkgroupExhausted(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &websocketClient{
		client:       mockWs,
		uri:          "wss://test.com",
		messageQueue: make(chan []byte, 1),
	}

	client.messageQueue <- []byte("test message")

	var listenCallback websocket.HandlerFunc
	mockWs.listenFn = func(callback websocket.HandlerFunc) error {
		listenCallback = callback
		return nil
	}

	err := client.Stream()
	assert.NoError(t, err)

	err = listenCallback([]byte("another message"))
	assert.Error(t, err, "Should return error when queue is full")
}

func TestPublicWebsocketClient_Stream_KLine(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client:           mockWs,
			uri:              "wss://test.com",
			workerPoolSize:   10,
			workerBufferSize: 100,
			messageQueue:     make(chan []byte, 100),
			quit:             make(chan struct{}),
			processFunc:      nil,
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}
	client.websocketClient.processFunc = client.processMessage

	var listenCallback websocket.HandlerFunc
	mockWs.listenFn = func(callback websocket.HandlerFunc) error {
		listenCallback = callback
		return nil
	}

	client.websocketClient.startWorkerPool(context.Background())

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

	select {
	case <-ch:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message to be processed")
	}

	klineMsg := sub.msg

	assert.True(t, sub.called)
	assert.Equal(t, "market_kline_1min", klineMsg.Channel)
	assert.Equal(t, model.ParseSymbol("BTCUSDT"), klineMsg.Symbol)
	assert.Equal(t, int64(1732178884994), klineMsg.Ts)
	assert.Equal(t, 0.0010, klineMsg.Data.OpenPrice)
}

func TestPublicWebsocketClient_ProcessInvalidJSON(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client:           mockWs,
			uri:              "wss://test.com",
			workerPoolSize:   1,
			workerBufferSize: 10,
			messageQueue:     make(chan []byte, 10),
			quit:             make(chan struct{}),
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}
	client.websocketClient.processFunc = client.processMessage

	invalidJSON := []byte(`{ invalid json }`)
	client.processMessage(invalidJSON)

}

func TestPublicWebsocketClient_ProcessErrorMessage(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client:           mockWs,
			uri:              "wss://test.com",
			workerPoolSize:   1,
			workerBufferSize: 10,
			messageQueue:     make(chan []byte, 10),
			quit:             make(chan struct{}),
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}
	client.websocketClient.processFunc = client.processMessage

	errorMsg := []byte(`{ "error": "subscription failed" }`)
	client.processMessage(errorMsg)

}

func TestPublicWebsocketClient_ProcessInvalidChannel(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client:           mockWs,
			uri:              "wss://test.com",
			workerPoolSize:   1,
			workerBufferSize: 10,
			messageQueue:     make(chan []byte, 10),
			quit:             make(chan struct{}),
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}
	client.websocketClient.processFunc = client.processMessage

	invalidChannelMsg := []byte(`{ "ch": "invalid_channel_format", "symbol": "BTCUSDT" }`)
	client.processMessage(invalidChannelMsg)

}

func TestStartWorkerPool_NilProcessFunc(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &websocketClient{
		client:       mockWs,
		uri:          "wss://test.com",
		messageQueue: make(chan []byte, 10),
		processFunc:  nil,
	}

	err := client.startWorkerPool(context.Background())
	assert.Error(t, err, "Should return error when processFunc is nil")
}

func TestNewPublicWebsocket_WithOptions(t *testing.T) {
	customURI := "wss://custom.example.com/"
	customWorkerPoolSize := 5
	customBufferSize := 50

	uriOption := func(ws *websocketClient) {
		ws.uri = customURI
	}

	client, err := NewPublicWebsocket(context.Background(), uriOption, WithWorkerBufferSize(customBufferSize), WithWorkerPoolSize(customWorkerPoolSize))
	assert.NoError(t, err)

	pubClient, ok := client.(*publicWebsocketClient)
	assert.True(t, ok)

	assert.Equal(t, customURI, pubClient.websocketClient.uri)
	assert.Equal(t, customWorkerPoolSize, pubClient.websocketClient.workerPoolSize)
	assert.Equal(t, customBufferSize, cap(pubClient.websocketClient.messageQueue))

	pubClient.Disconnect()
}

func TestPublicWebsocketClient_MultipleSubscribers(t *testing.T) {
	mockWs := &mockWsClient{}
	writeCount := 0
	var writtenData [][]byte

	mockWs.writeFn = func(data []byte) error {
		writeCount++
		writtenData = append(writtenData, data)
		return nil
	}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client:           mockWs,
			uri:              "wss://test.com",
			workerPoolSize:   10,
			workerBufferSize: 100,
			messageQueue:     make(chan []byte, 100),
			quit:             make(chan struct{}),
		},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}

	// Create multiple subscribers for the same symbol
	sub1 := &subTest{ch: make(chan struct{}, 1)}
	sub2 := &subTest{ch: make(chan struct{}, 1)}
	sub3 := &subTest{ch: make(chan struct{}, 1)}

	// Subscribe first subscriber
	err := client.SubscribeKLine(sub1)
	assert.NoError(t, err)
	assert.Equal(t, 1, writeCount, "Should send subscribe request for first subscriber")

	// Subscribe second subscriber for same symbol
	err = client.SubscribeKLine(sub2)
	assert.NoError(t, err)
	assert.Equal(t, 1, writeCount, "Should NOT send another subscribe request for second subscriber")

	// Subscribe third subscriber for same symbol
	err = client.SubscribeKLine(sub3)
	assert.NoError(t, err)
	assert.Equal(t, 1, writeCount, "Should NOT send another subscribe request for third subscriber")

	// Reset write count
	writeCount = 0
	writtenData = nil

	// Unsubscribe first subscriber
	err = client.UnsubscribeKLine(sub1)
	assert.NoError(t, err)
	assert.Equal(t, 0, writeCount, "Should NOT send unsubscribe request when other subscribers exist")

	// Unsubscribe second subscriber
	err = client.UnsubscribeKLine(sub2)
	assert.NoError(t, err)
	assert.Equal(t, 0, writeCount, "Should NOT send unsubscribe request when other subscribers exist")

	// Unsubscribe third (last) subscriber
	err = client.UnsubscribeKLine(sub3)
	assert.NoError(t, err)
	assert.Equal(t, 1, writeCount, "SHOULD send unsubscribe request when last subscriber is removed")

	// Verify the unsubscribe message
	var unsubReq SubscribeRequest
	err = json.Unmarshal(writtenData[0], &unsubReq)
	assert.NoError(t, err)
	assert.Equal(t, "unsubscribe", unsubReq.Op)
}

func TestParseChannel(t *testing.T) {
	testCases := []struct {
		name             string
		channelStr       string
		expectError      bool
		expectedChannel  model.Channel
		expectedType     model.PriceType
		expectedInterval model.Interval
	}{
		{
			name:             "Valid Market KLine Channel",
			channelStr:       "market_kline_1min",
			expectError:      false,
			expectedChannel:  model.ChannelKline,
			expectedType:     model.PriceTypeMarket,
			expectedInterval: model.Interval1Min,
		},
		{
			name:             "Valid Mark KLine Channel",
			channelStr:       "mark_kline_5min",
			expectError:      false,
			expectedChannel:  model.ChannelKline,
			expectedType:     model.PriceTypeMark,
			expectedInterval: model.Interval5Min,
		},
		{
			name:        "Invalid Channel Format - Too Few Parts",
			channelStr:  "market_kline",
			expectError: true,
		},
		{
			name:        "Invalid Channel Format - Invalid Price Type",
			channelStr:  "invalid_kline_1min",
			expectError: true,
		},
		{
			name:        "Invalid Channel Format - Invalid Channel",
			channelStr:  "market_invalid_1min",
			expectError: true,
		},
		{
			name:        "Invalid Channel Format - Invalid Interval",
			channelStr:  "market_kline_invalid",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interval, channel, priceType, err := parseChannel(tc.channelStr)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedInterval, interval)
				assert.Equal(t, tc.expectedChannel, channel)
				assert.Equal(t, tc.expectedType, priceType)
			}
		})
	}
}

func TestConcurrentKLineSubscriptions(t *testing.T) {

	numSubscribers := 3

	mockWs := &mockWsClient{}
	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client:           mockWs,
			uri:              "wss://test.com",
			workerPoolSize:   3,
			workerBufferSize: 100,
			messageQueue:     make(chan []byte, 100),
			quit:             make(chan struct{}),
		},
		subscriberMtx: sync.Mutex{},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}
	client.websocketClient.processFunc = client.processMessage

	var callCount int32
	mockWs.writeFn = func(bytes []byte) error {
		atomic.AddInt32(&callCount, 1)
		return nil
	}

	subscribers := make([]*subTest, numSubscribers)
	var wg sync.WaitGroup
	wg.Add(numSubscribers)

	msgNum := 5

	symbols := []model.Symbol{model.Symbol("BTCUSDT"), model.Symbol("ETHUSDT"), model.Symbol("SUIUSDT")}

	for i := 0; i < numSubscribers; i++ {
		subscribers[i] = &subTest{
			symbol:      &symbols[i],
			called:      false,
			msg:         model.KLineChannelMessage{},
			ch:          make(chan struct{}, msgNum),
			subscribeMu: sync.Mutex{},
		}

		go func(idx int) {
			defer wg.Done()
			err := client.SubscribeKLine(subscribers[idx])
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(numSubscribers), callCount, "Expected all subscription requests to be processed")

	client.subscriberMtx.Lock()
	assert.Equal(t, numSubscribers, len(client.klineHandlers), "All subscribers should be in the handlers map")
	client.subscriberMtx.Unlock()

	klineData := `{
		"ch": "market_kline_1min",
		"symbol": "%s",
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

	client.websocketClient.startWorkerPool(context.Background())

	for _, symbol := range symbols {
		for i := 0; i < msgNum; i++ {
			client.websocketClient.messageQueue <- []byte(fmt.Sprintf(klineData, symbol.String()))
			time.Sleep(10 * time.Millisecond)
		}
	}

	time.Sleep(200 * time.Millisecond)

	for i, subscriber := range subscribers {
		count := 0
		for {
			select {
			case <-subscriber.ch:
				count++
			default:
				// No more messages
				assert.Equal(t, msgNum, count, "sub %d should have received %d messages", i, msgNum)
				goto nextSubscriber
			}
		}
	nextSubscriber:
	}

	client.Disconnect()
}

func TestConcurrentUnsubscribe(t *testing.T) {
	numSubscribers := 3

	mockWs := &mockWsClient{}
	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client:           mockWs,
			uri:              "wss://test.com",
			workerPoolSize:   3,
			workerBufferSize: 100,
			messageQueue:     make(chan []byte, 100),
			quit:             make(chan struct{}),
		},
		subscriberMtx: sync.Mutex{},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}

	var writeMutex sync.Mutex
	var subscribeCount int32
	var unsubscribeCount int32
	mockWs.writeFn = func(bytes []byte) error {
		writeMutex.Lock()
		defer writeMutex.Unlock()

		var req SubscribeRequest
		if err := json.Unmarshal(bytes, &req); err != nil {
			return err
		}

		if req.Op == "subscribe" {
			atomic.AddInt32(&subscribeCount, 1)
		} else if req.Op == "unsubscribe" {
			atomic.AddInt32(&unsubscribeCount, 1)
		}

		return nil
	}
	symbols := []model.Symbol{model.Symbol("BTCUSDT"), model.Symbol("ETHUSDT"), model.Symbol("SUIUSDT")}

	subscribers := make([]*subTest, numSubscribers)
	for i := 0; i < numSubscribers; i++ {
		subscribers[i] = &subTest{symbol: &symbols[i], ch: make(chan struct{}, 1)}
		err := client.SubscribeKLine(subscribers[i])
		assert.NoError(t, err)
	}

	assert.Equal(t, int32(numSubscribers), subscribeCount)

	client.subscriberMtx.Lock()
	assert.Equal(t, numSubscribers, len(client.klineHandlers), "All subscribers should be in the map")
	client.subscriberMtx.Unlock()

	// Then concurrently unsubscribe to test thread safety
	var wg sync.WaitGroup
	wg.Add(numSubscribers)

	for i := 0; i < numSubscribers; i++ {
		go func(idx int) {
			defer wg.Done()
			err := client.UnsubscribeKLine(subscribers[idx])
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(numSubscribers), unsubscribeCount, "All unsubscribe requests should be processed")

	client.subscriberMtx.Lock()
	assert.Equal(t, 0, len(client.klineHandlers), "All subscribers should be removed from handlers map")
	client.subscriberMtx.Unlock()
}

func TestConcurrentMessageProcessing(t *testing.T) {
	mockWs := &mockWsClient{}

	client := &publicWebsocketClient{
		websocketClient: &websocketClient{
			client:           mockWs,
			uri:              "wss://test.com",
			workerPoolSize:   5,
			workerBufferSize: 100,
			messageQueue:     make(chan []byte, 100),
			quit:             make(chan struct{}),
		},
		subscriberMtx: sync.Mutex{},
		klineHandlers: make(map[KLineSubscriber]struct{}),
	}
	client.websocketClient.processFunc = client.processMessage

	numSubscribers := 20
	subscribers := make([]*subTest, numSubscribers)
	for i := 0; i < numSubscribers; i++ {
		subscribers[i] = &subTest{ch: make(chan struct{}, 10)}
		client.klineHandlers[subscribers[i]] = struct{}{}
	}

	client.websocketClient.startWorkerPool(context.Background())

	numMessages := 50
	messages := make([][]byte, numMessages)

	for i := 0; i < numMessages; i++ {
		klineData := fmt.Sprintf(`{
			"ch": "market_kline_1min",
			"symbol": "BTCUSDT",
			"ts": %d,
			"data": {
				"o": "%0.4f",
				"c": "%0.4f",
				"h": "%0.4f",
				"l": "%0.4f",
				"b": "%0.2f",
				"q": "%0.2f"
			}
		}`, time.Now().UnixNano()/int64(time.Millisecond),
			0.0010+float64(i)*0.0001,
			0.0020+float64(i)*0.0001,
			0.0025+float64(i)*0.0001,
			0.0015+float64(i)*0.0001,
			1.01+float64(i)*0.01,
			1.09+float64(i)*0.01)

		messages[i] = []byte(klineData)
	}

	var wg sync.WaitGroup
	wg.Add(numMessages)

	for i := 0; i < numMessages; i++ {
		go func(msg []byte) {
			defer wg.Done()
			select {
			case client.websocketClient.messageQueue <- msg:
			case <-time.After(100 * time.Millisecond):
			}
		}(messages[i])
	}

	wg.Wait()

	time.Sleep(200 * time.Millisecond)

	receivedMessages := 0
	for _, sub := range subscribers {
		count := 0
		for {
			select {
			case <-sub.ch:
				count++
			default:
				// No more messages
				receivedMessages += count
				goto nextSubscriber
			}
		}
	nextSubscriber:
	}

	client.Disconnect()

	t.Logf("Processed %d messages across subscribers", receivedMessages)
}

func TestWaitForDisconnection(t *testing.T) {

	mockWs := &mockWsClient{}

	var closeMethodCalled bool
	mockWs.closeFn = func() {
		closeMethodCalled = true
	}

	client := &websocketClient{
		client:       mockWs,
		uri:          "wss://test.com",
		messageQueue: make(chan []byte, 10),
		quit:         make(chan struct{}),
	}

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-client.quit:
				close(done)
				return
			}
		}
	}()

	client.Disconnect()

	assert.True(t, closeMethodCalled, "The close method should have been called")

	select {
	case <-done:

	case <-time.After(100 * time.Millisecond):
		t.Fatal("Worker goroutine did not exit in time")
	}

	_, ok := <-client.messageQueue
	assert.False(t, ok, "Message queue should be closed")

	_, ok = <-client.quit
	assert.False(t, ok, "Quit channel should be closed")
}
