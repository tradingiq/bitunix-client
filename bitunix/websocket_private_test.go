package bitunix

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/security"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tradingiq/bitunix-client/model"
)

type mockWebsocketServer struct {
	server      *httptest.Server
	handlers    map[string]func(message []byte, conn *websocket.Conn)
	connections []*websocket.Conn
	mu          sync.Mutex
}

func newMockWebsocketServer() *mockWebsocketServer {
	mock := &mockWebsocketServer{
		handlers:    make(map[string]func(message []byte, conn *websocket.Conn)),
		connections: make([]*websocket.Conn, 0),
	}

	mock.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}

		mock.mu.Lock()
		mock.connections = append(mock.connections, conn)
		mock.mu.Unlock()

		initialMsg := map[string]interface{}{
			"event": "connected",
			"msg":   "WebSocket connected",
		}
		conn.Write(context.Background(), websocket.MessageText, mustMarshal(initialMsg))

		for {
			msgType, data, err := conn.Read(context.Background())
			if err != nil {
				return
			}

			if msgType == websocket.MessageText {
				var msg map[string]interface{}
				if err := json.Unmarshal(data, &msg); err != nil {
					continue
				}

				if op, ok := msg["op"].(string); ok {
					if op == "login" {

						loginResponse := map[string]interface{}{
							"op": "login",
							"data": map[string]interface{}{
								"result": true,
							},
						}
						conn.Write(context.Background(), websocket.MessageText, mustMarshal(loginResponse))
						continue
					}

					if op == "ping" {

						pongResponse := map[string]interface{}{
							"op":   "pong",
							"pong": time.Now().Unix(),
						}
						conn.Write(context.Background(), websocket.MessageText, mustMarshal(pongResponse))
						continue
					}
				}

				mock.mu.Lock()
				for key, handler := range mock.handlers {
					if strings.Contains(string(data), key) {
						go handler(data, conn)
					}
				}
				mock.mu.Unlock()
			}
		}
	}))

	return mock
}

func (m *mockWebsocketServer) close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, conn := range m.connections {
		conn.Close(websocket.StatusNormalClosure, "server shutdown")
	}
	m.server.Close()
}

func (m *mockWebsocketServer) registerHandler(key string, handler func(message []byte, conn *websocket.Conn)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[key] = handler
}

func (m *mockWebsocketServer) broadcastToAll(msg interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var data []byte

	switch v := msg.(type) {
	case []byte:
		data = v
	default:
		data = mustMarshal(msg)
	}

	for _, conn := range m.connections {
		conn.Write(context.Background(), websocket.MessageText, data)
	}
}

func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

type testBalanceSubscriber struct {
	channel chan model.BalanceChannelMessage
}

func (s *testBalanceSubscriber) SubscribeBalance(msg *model.BalanceChannelMessage) {
	s.channel <- *msg
}

type testPositionSubscriber struct {
	channel chan model.PositionChannelMessage
}

func (s *testPositionSubscriber) SubscribePosition(msg *model.PositionChannelMessage) {
	s.channel <- *msg
}

type testOrderSubscriber struct {
	channel chan model.OrderChannelMessage
}

func (s *testOrderSubscriber) SubscribeOrder(msg *model.OrderChannelMessage) {
	s.channel <- *msg
}

type testTpSlOrderSubscriber struct {
	channel chan model.TpSlOrderChannelMessage
}

func (s *testTpSlOrderSubscriber) SubscribeTpSlOrder(msg *model.TpSlOrderChannelMessage) {
	s.channel <- *msg
}

func TestPrivateWebsocketClient(t *testing.T) {

	mockServer := newMockWebsocketServer()
	defer mockServer.close()

	serverURL := strings.TrimPrefix(mockServer.server.URL, "http://")
	wsURL := "ws://" + serverURL + "/private/"

	withCustomURI := func(ws *websocketClient) {
		ws.uri = wsURL
	}

	mockBalanceJSON := `{
		"ch": "balance",
		"ts": 1651234567890,
		"data": {
			"coin": "USDT",
			"available": "1000.0",
			"frozen": "100.0",
			"isolationFrozen": "50.0",
			"crossFrozen": "50.0",
			"margin": "200.0",
			"isolationMargin": "100.0",
			"crossMargin": "100.0",
			"expMoney": "50.0"
		}
	}`

	mockPositionJSON := `{
		"ch": "position",
		"ts": 1651234567890,
		"data": {
			"event": "UPDATE",
			"positionId": "pos123456",
			"marginMode": "ISOLATION",
			"positionMode": "ONE_WAY",
			"side": "LONG",
			"leverage": "10",
			"margin": "100.0",
			"ctime": "2023-01-01T00:00:00.000Z",
			"qty": "1.5",
			"entryValue": "50000.0",
			"symbol": "BTCUSDT",
			"realizedPNL": "100.0",
			"unrealizedPNL": "50.0",
			"funding": "10.0",
			"fee": "5.0"
		}
	}`

	mockOrderJSON := `{
		"ch": "order",
		"ts": 1651234567890,
		"data": {
			"event": "CREATE",
			"orderId": "ord123456",
			"symbol": "BTCUSDT",
			"positionType": "ISOLATION",
			"positionMode": "ONE_WAY",
			"side": "BUY",
			"type": "LIMIT",
			"qty": "1.0",
			"reductionOnly": false,
			"price": "50000.0",
			"ctime": "2023-01-01T00:00:00.000Z",
			"mtime": "2023-01-01T00:00:00.000Z",
			"leverage": "10",
			"orderStatus": "NEW",
			"fee": "0.1"
		}
	}`

	mockTpSlJSON := `{
		"ch": "tpsl",
		"ts": 1651234567890,
		"data": {
			"event": "CREATE",
			"positionId": "pos123456",
			"orderId": "tpsl123456",
			"symbol": "BTCUSDT",
			"leverage": "10",
			"side": "BUY",
			"positionMode": "ONE_WAY",
			"status": "NEW",
			"ctime": "2023-01-01T00:00:00.000Z",
			"type": "POSITION_TPSL",
			"tpQty": "1.0",
			"slQty": "1.0",
			"tpStopType": "LAST_PRICE",
			"tpPrice": "52000.0",
			"tpOrderType": "LIMIT",
			"tpOrderPrice": "52000.0",
			"slStopType": "LAST_PRICE",
			"slPrice": "48000.0",
			"slOrderType": "LIMIT",
			"slOrderPrice": "48000.0"
		}
	}`

	mockServer.registerHandler(model.ChannelBalance, func(message []byte, conn *websocket.Conn) {

	})

	ctx := context.Background()
	c, _ := NewPrivateWebsocket(ctx, "test_api_key", "test_api_secret", withCustomURI)
	require.NotNil(t, c)

	client := c.(*privateWebsocketClient)

	err := client.Connect()
	require.NoError(t, err)
	defer client.Disconnect()

	streamDone := make(chan struct{})
	go func() {
		defer close(streamDone)
		err := client.Stream()
		if !errors.Is(err, bitunix_errors.ErrConnectionClosed) {
			t.Logf("Streaming error: %v", err.Error())
		}
	}()

	time.Sleep(50 * time.Millisecond)

	t.Run("BalanceUpdates", func(t *testing.T) {

		balanceCh := make(chan model.BalanceChannelMessage, 1)

		balanceSubscriber := &testBalanceSubscriber{
			channel: balanceCh,
		}

		err := client.SubscribeBalance(balanceSubscriber)
		require.NoError(t, err)

		defer client.UnsubscribeBalance(balanceSubscriber)

		mockServer.broadcastToAll([]byte(mockBalanceJSON))

		select {
		case balance := <-balanceCh:
			assert.Equal(t, model.ChannelBalance, balance.Ch)
			assert.Equal(t, "USDT", balance.Data.Coin)
			assert.Equal(t, 1000.0, balance.Data.Available)
			assert.Equal(t, 100.0, balance.Data.Frozen)
		case <-time.After(2 * time.Second):
			t.Fatal("Timed out waiting for balance update in subscriber")
		}
	})

	t.Run("PositionUpdates", func(t *testing.T) {

		positionCh := make(chan model.PositionChannelMessage, 1)

		positionSubscriber := &testPositionSubscriber{
			channel: positionCh,
		}

		err := client.SubscribePositions(positionSubscriber)
		require.NoError(t, err)

		defer client.UnsubscribePosition(positionSubscriber)

		mockServer.broadcastToAll([]byte(mockPositionJSON))

		select {
		case position := <-positionCh:
			assert.Equal(t, model.ChannelPosition, position.Channel)
			assert.Equal(t, "pos123456", position.Data.PositionID)
			assert.Equal(t, model.PositionEventUpdate, position.Data.Event)
			assert.Equal(t, model.PositionSideLong, position.Data.Side)
			assert.Equal(t, "BTCUSDT", position.Data.Symbol.String())
		case <-time.After(2 * time.Second):
			t.Fatal("Timed out waiting for position update in subscriber")
		}
	})

	t.Run("OrderUpdates", func(t *testing.T) {

		orderCh := make(chan model.OrderChannelMessage, 1)

		orderSubscriber := &testOrderSubscriber{
			channel: orderCh,
		}

		err := client.SubscribeOrders(orderSubscriber)
		require.NoError(t, err)

		defer client.UnsubscribeOrders(orderSubscriber)

		mockServer.broadcastToAll([]byte(mockOrderJSON))

		select {
		case order := <-orderCh:
			assert.Equal(t, model.ChannelOrder, order.Channel)
			assert.Equal(t, "ord123456", order.Data.OrderID)
			assert.Equal(t, model.OrderEventCreate, order.Data.Event)
			assert.Equal(t, model.TradeSideBuy, order.Data.Side)
			assert.Equal(t, "BTCUSDT", order.Data.Symbol.String())
		case <-time.After(2 * time.Second):
			t.Fatal("Timed out waiting for order update in subscriber")
		}
	})

	t.Run("TpSlUpdates", func(t *testing.T) {

		tpslCh := make(chan model.TpSlOrderChannelMessage, 1)

		tpslSubscriber := &testTpSlOrderSubscriber{
			channel: tpslCh,
		}

		err := client.SubscribeTpSlOrders(tpslSubscriber)
		require.NoError(t, err)

		defer client.UnsubscribeTpSlOrders(tpslSubscriber)

		mockServer.broadcastToAll([]byte(mockTpSlJSON))

		select {
		case tpsl := <-tpslCh:
			assert.Equal(t, model.ChannelTpSl, tpsl.Channel)
			assert.Equal(t, "tpsl123456", tpsl.Data.OrderID)
			assert.Equal(t, "pos123456", tpsl.Data.PositionID)
			assert.Equal(t, model.TPSLEventCreate, tpsl.Data.Event)
			assert.Equal(t, model.TPSLTypeFull, tpsl.Data.Type)
			assert.Equal(t, "BTCUSDT", tpsl.Data.Symbol.String())
		case <-time.After(2 * time.Second):
			t.Fatal("Timed out waiting for tp/sl update in subscriber")
		}
	})

	t.Run("MultipleSubscriptions", func(t *testing.T) {

		balanceCh1 := make(chan model.BalanceChannelMessage, 1)
		balanceCh2 := make(chan model.BalanceChannelMessage, 1)

		subscriber1 := &testBalanceSubscriber{channel: balanceCh1}
		subscriber2 := &testBalanceSubscriber{channel: balanceCh2}

		err := client.SubscribeBalance(subscriber1)
		require.NoError(t, err)
		err = client.SubscribeBalance(subscriber2)
		require.NoError(t, err)

		defer client.UnsubscribeBalance(subscriber1)
		defer client.UnsubscribeBalance(subscriber2)

		mockServer.broadcastToAll([]byte(mockBalanceJSON))

		for i, ch := range []chan model.BalanceChannelMessage{balanceCh1, balanceCh2} {
			select {
			case balance := <-ch:
				assert.Equal(t, model.ChannelBalance, balance.Ch)
				assert.Equal(t, "USDT", balance.Data.Coin)
			case <-time.After(2 * time.Second):
				t.Fatalf("Timed out waiting for balance update on subscriber %d", i+1)
			}
		}
	})

	t.Run("Unsubscribe", func(t *testing.T) {

		balanceCh := make(chan model.BalanceChannelMessage, 1)
		subscriber := &testBalanceSubscriber{channel: balanceCh}

		err := client.SubscribeBalance(subscriber)
		require.NoError(t, err)

		client.balanceSubscriberMtx.Lock()
		_, exists := client.balanceSubscribers[subscriber]
		client.balanceSubscriberMtx.Unlock()
		assert.True(t, exists, "Subscriber should be in subscribers map")

		err = client.UnsubscribeBalance(subscriber)
		require.NoError(t, err)

		client.balanceSubscriberMtx.Lock()
		_, stillExists := client.balanceSubscribers[subscriber]
		client.balanceSubscriberMtx.Unlock()
		assert.False(t, stillExists, "Subscriber should not be in subscribers map after unsubscribe")
	})
}

func TestPrivateWebsocketExceptions(t *testing.T) {

	failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conn, err := websocket.Accept(w, r, nil)
		if err == nil {
			conn.Close(websocket.StatusInternalError, "Test error")
		}
	}))
	defer failServer.Close()

	failURL := strings.TrimPrefix(failServer.URL, "http://")
	failWsURL := "ws://" + failURL + "/private/"

	withFailURI := func(ws *websocketClient) {
		ws.uri = failWsURL
	}

	ctx := context.Background()
	failClient, _ := NewPrivateWebsocket(ctx, "test_api_key", "test_api_secret", withFailURI)
	require.NotNil(t, failClient)

	err := failClient.Connect()
	assert.Error(t, err, "Expected connection to fail")
}

func TestCustomWebsocketURI(t *testing.T) {

	customURI := "wss://custom.example.com/ws/"

	withCustomURI := func(ws *websocketClient) {
		ws.uri = customURI
	}

	ctx := context.Background()
	c, _ := NewPrivateWebsocket(ctx, "test_api_key", "test_api_secret", withCustomURI)
	client := c.(*privateWebsocketClient)

	assert.Equal(t, customURI, client.websocketClient.uri)
}

func TestGenerateWebsocketSignature(t *testing.T) {
	apiKey := "test_api_key"
	apiSecret := "test_api_secret"
	timestamp := int64(1234567890)
	nonceBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	expectedNonceHex := "0102030405060708"
	expectedPreSign := expectedNonceHex + "1234567890" + apiKey
	expectedPreSignHash := security.Sha256Hex(expectedPreSign)
	expectedSign := security.Sha256Hex(expectedPreSignHash + apiSecret)

	actualSign, actualTimestamp := generateWebsocketSignature(apiKey, apiSecret, timestamp, nonceBytes)

	assert.Equal(t, expectedSign, actualSign, "The signature should match the expected value")
	assert.Equal(t, timestamp, actualTimestamp, "The timestamp should be returned unchanged")
}

func TestGenerateWebsocketSignatureFormat(t *testing.T) {
	apiKey := "test_api_key"
	apiSecret := "test_api_secret"
	timestamp := int64(1234567890)
	nonceBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	sign, _ := generateWebsocketSignature(apiKey, apiSecret, timestamp, nonceBytes)

	assert.Equal(t, 64, len(sign), "Signature should be 64 characters long (32 bytes in hex)")

	_, err := hex.DecodeString(sign)
	assert.NoError(t, err, "Signature should be a valid hex string")
}

func TestGenerateWebsocketSignatureDifferentInputs(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     string
		apiSecret  string
		timestamp  int64
		nonceBytes []byte
	}{
		{
			name:       "Empty inputs",
			apiKey:     "",
			apiSecret:  "",
			timestamp:  0,
			nonceBytes: []byte{},
		},
		{
			name:       "Different apiClient key",
			apiKey:     "different_api_key",
			apiSecret:  "test_api_secret",
			timestamp:  1234567890,
			nonceBytes: []byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:       "Different apiClient secret",
			apiKey:     "test_api_key",
			apiSecret:  "different_api_secret",
			timestamp:  1234567890,
			nonceBytes: []byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:       "Different timestamp",
			apiKey:     "test_api_key",
			apiSecret:  "test_api_secret",
			timestamp:  9876543210,
			nonceBytes: []byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:       "Different nonce",
			apiKey:     "test_api_key",
			apiSecret:  "test_api_secret",
			timestamp:  1234567890,
			nonceBytes: []byte{0x08, 0x07, 0x06, 0x05},
		},
	}

	signatures := make(map[string]bool)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sign, returnedTimestamp := generateWebsocketSignature(tc.apiKey, tc.apiSecret, tc.timestamp, tc.nonceBytes)

			assert.Equal(t, tc.timestamp, returnedTimestamp)

			assert.Equal(t, 64, len(sign))
			_, err := hex.DecodeString(sign)
			assert.NoError(t, err)

			signatures[sign] = true
		})
	}

	assert.GreaterOrEqual(t, len(signatures), len(tests)-1,
		"Different inputs should produce different signatures in most cases")
}

func TestKeepAliveMonitor(t *testing.T) {
	heartbeatGenerator := KeepAliveMonitor()
	require.NotNil(t, heartbeatGenerator)

	bytes, err := heartbeatGenerator()
	require.NoError(t, err)
	require.NotNil(t, bytes)

	var message heartbeatMessage
	err = json.Unmarshal(bytes, &message)
	require.NoError(t, err)

	assert.Equal(t, "ping", message.Op)
	assert.NotZero(t, message.Ping)
	assert.LessOrEqual(t, message.Ping, time.Now().Unix())
	assert.GreaterOrEqual(t, message.Ping, time.Now().Unix()-5)
}

func TestWebsocketSigner(t *testing.T) {
	apiKey := "test_api_key"
	apiSecret := "test_api_secret"

	signer := WebsocketSigner(apiKey, apiSecret)
	require.NotNil(t, signer)

	bytes, err := signer()
	require.NoError(t, err)
	require.NotNil(t, bytes)

	var message loginMessage
	err = json.Unmarshal(bytes, &message)
	require.NoError(t, err)

	assert.Equal(t, "login", message.Op)
	assert.Len(t, message.Args, 1)

	loginParams := message.Args[0]
	assert.Equal(t, apiKey, loginParams.ApiKey, "The apiClient key in the message should match the provided key")
	assert.NotEmpty(t, loginParams.Nonce)
	assert.NotEmpty(t, loginParams.Sign)
	assert.NotZero(t, loginParams.Timestamp)

	nonceBytes, err := hex.DecodeString(loginParams.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(nonceBytes))

	_, err = hex.DecodeString(loginParams.Sign)
	assert.NoError(t, err)
	assert.Equal(t, 64, len(loginParams.Sign))

	differentApiKey := "different_api_key"
	differentApiSecret := "different_api_secret"

	differentSigner := WebsocketSigner(differentApiKey, differentApiSecret)
	differentBytes, err := differentSigner()
	require.NoError(t, err)

	var differentMessage loginMessage
	err = json.Unmarshal(differentBytes, &differentMessage)
	require.NoError(t, err)

	assert.Equal(t, differentApiKey, differentMessage.Args[0].ApiKey)

	knownTimestamp := int64(1234567890)
	knownNonce := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	sign1, _ := generateWebsocketSignature(apiKey, apiSecret, knownTimestamp, knownNonce)
	sign2, _ := generateWebsocketSignature(differentApiKey, differentApiSecret, knownTimestamp, knownNonce)

	assert.NotEqual(t, sign1, sign2,
		"Different apiClient credentials should produce different signatures")
}

func TestLoginMessage(t *testing.T) {
	message := loginMessage{
		Op: "login",
		Args: []loginParams{
			{
				ApiKey:    "test_key",
				Timestamp: 1234567890,
				Nonce:     "0102030405060708",
				Sign:      "test_signature",
			},
		},
	}

	bytes, err := json.Marshal(message)
	require.NoError(t, err)

	var decoded loginMessage
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "login", decoded.Op)
	assert.Len(t, decoded.Args, 1)
	assert.Equal(t, "test_key", decoded.Args[0].ApiKey)
	assert.Equal(t, int64(1234567890), decoded.Args[0].Timestamp)
	assert.Equal(t, "0102030405060708", decoded.Args[0].Nonce)
	assert.Equal(t, "test_signature", decoded.Args[0].Sign)
}
