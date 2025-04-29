package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockWSServer struct {
	server      *httptest.Server
	connections []*websocket.Conn
	mu          sync.Mutex

	messageHandlers []func([]byte, *websocket.Conn)
}

func newMockWSServer() *mockWSServer {
	m := &mockWSServer{
		connections:     make([]*websocket.Conn, 0),
		messageHandlers: make([]func([]byte, *websocket.Conn), 0),
	}

	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}

		m.mu.Lock()
		m.connections = append(m.connections, conn)
		m.mu.Unlock()

		initialMsg := map[string]interface{}{
			"event": "connected",
			"msg":   "WebSocket connected",
		}
		data, _ := json.Marshal(initialMsg)
		conn.Write(context.Background(), websocket.MessageText, data)

		for {
			msgType, data, err := conn.Read(context.Background())
			if err != nil {
				return
			}

			if msgType == websocket.MessageText {
				m.mu.Lock()
				handlers := make([]func([]byte, *websocket.Conn), len(m.messageHandlers))
				copy(handlers, m.messageHandlers)
				m.mu.Unlock()

				for _, handler := range handlers {
					go handler(data, conn)
				}
			}
		}
	}))

	return m
}

func (m *mockWSServer) close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, conn := range m.connections {
		conn.Close(websocket.StatusNormalClosure, "server shutdown")
	}
	m.server.Close()
}

func (m *mockWSServer) addMessageHandler(handler func([]byte, *websocket.Conn)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messageHandlers = append(m.messageHandlers, handler)
}

func (m *mockWSServer) broadcastToAll(msg interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, _ := json.Marshal(msg)
	for _, conn := range m.connections {
		conn.Write(context.Background(), websocket.MessageText, data)
	}
}

func TestClientWithMockServer(t *testing.T) {

	mock := newMockWSServer()
	defer mock.close()

	mock.addMessageHandler(func(data []byte, conn *websocket.Conn) {
		var msg map[string]interface{}
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return
		}

		if op, ok := msg["op"].(string); ok && op == "login" {
			loginResponse := map[string]interface{}{
				"event": "login",
				"msg":   "Login successful",
			}
			responseData, _ := json.Marshal(loginResponse)
			conn.Write(context.Background(), websocket.MessageText, responseData)
		}

		if op, ok := msg["op"].(string); ok && op == "ping" {
			pongResponse := map[string]interface{}{
				"op":   "pong",
				"pong": time.Now().Unix(),
			}
			responseData, _ := json.Marshal(pongResponse)
			conn.Write(context.Background(), websocket.MessageText, responseData)
		}
	})

	serverURL := strings.TrimPrefix(mock.server.URL, "http://")
	wsURL := "ws://" + serverURL

	t.Run("ClientConnect", func(t *testing.T) {

		ctx := context.Background()
		client := New(ctx, wsURL)

		err := client.Connect()
		require.NoError(t, err)

		client.Close()
	})

	t.Run("ClientAuthentication", func(t *testing.T) {

		testLoginGenerator := func() ([]byte, error) {
			loginMsg := map[string]interface{}{
				"op": "login",
				"args": []map[string]interface{}{
					{
						"apiKey":    "test_key",
						"timestamp": time.Now().Unix(),
						"nonce":     "test_nonce",
						"sign":      "test_signature",
					},
				},
			}
			return json.Marshal(loginMsg)
		}

		ctx := context.Background()
		client := New(ctx, wsURL, WithAuthentication(testLoginGenerator))

		err := client.Connect()
		require.NoError(t, err)

		client.Close()
	})

	t.Run("ClientHeartbeat", func(t *testing.T) {

		pingCount := 0
		var pingMutex sync.Mutex

		pingHandler := func(data []byte, conn *websocket.Conn) {
			var msg map[string]interface{}
			err := json.Unmarshal(data, &msg)
			if err != nil {
				return
			}

			if op, ok := msg["op"].(string); ok && op == "ping" {
				pingMutex.Lock()
				pingCount++
				pingMutex.Unlock()
			}
		}
		mock.addMessageHandler(pingHandler)

		testHeartbeatGenerator := func() ([]byte, error) {
			pingMsg := map[string]interface{}{
				"op":   "ping",
				"ping": time.Now().Unix(),
			}
			return json.Marshal(pingMsg)
		}

		ctx := context.Background()
		client := New(ctx, wsURL, WithKeepAliveMonitor(100*time.Millisecond, testHeartbeatGenerator))

		err := client.Connect()
		require.NoError(t, err)

		time.Sleep(350 * time.Millisecond)

		pingMutex.Lock()
		assert.GreaterOrEqual(t, pingCount, 3, "Should have sent at least 3 heartbeat messages")
		pingMutex.Unlock()

		client.Close()
	})

	t.Run("ClientMessaging", func(t *testing.T) {

		mock.addMessageHandler(func(data []byte, conn *websocket.Conn) {
			var msg map[string]interface{}
			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			if _, ok := msg["test"]; ok {
				response := map[string]interface{}{
					"response": true,
					"echo":     msg,
				}
				responseData, _ := json.Marshal(response)
				conn.Write(context.Background(), websocket.MessageText, responseData)
			}
		})

		ctx := context.Background()
		client := New(ctx, wsURL)

		err := client.Connect()
		require.NoError(t, err)

		testMessage := map[string]interface{}{
			"test": true,
			"data": "hello",
		}
		testData, _ := json.Marshal(testMessage)

		receivedMessages := 0
		messageReceived := make(chan bool, 1)

		go func() {
			err := client.Listen(func(data []byte) error {
				var msg map[string]interface{}
				if err := json.Unmarshal(data, &msg); err != nil {
					return err
				}

				if _, ok := msg["response"]; ok {
					receivedMessages++
					messageReceived <- true
				}

				return nil
			})
			if err != nil {
				t.Logf("Listen error: %v", err)
			}
		}()

		time.Sleep(100 * time.Millisecond)

		err = client.Write(testData)
		require.NoError(t, err)

		select {
		case <-messageReceived:
			assert.Equal(t, 1, receivedMessages, "Should have received exactly one message")
		case <-time.After(1 * time.Second):
			t.Fatal("Timed out waiting for response")
		}

		client.Close()
	})

	t.Run("ConcurrentMessageHandling", func(t *testing.T) {

		ctx := context.Background()
		client := New(ctx, wsURL)

		err := client.Connect()
		require.NoError(t, err)
		defer client.Close()

		receivedCount := 0
		var receivedMutex sync.Mutex
		messageReceived := make(chan bool, 5)

		go func() {
			err := client.Listen(func(data []byte) error {
				var msg map[string]interface{}
				if err := json.Unmarshal(data, &msg); err != nil {
					return err
				}

				if _, ok := msg["broadcast"]; ok {
					receivedMutex.Lock()
					receivedCount++
					receivedMutex.Unlock()
					messageReceived <- true
				}

				return nil
			})
			if err != nil {
				t.Logf("Listen error: %v", err)
			}
		}()

		time.Sleep(100 * time.Millisecond)

		for i := 0; i < 5; i++ {
			mock.broadcastToAll(map[string]interface{}{
				"broadcast": true,
				"index":     i,
			})
			time.Sleep(10 * time.Millisecond)
		}

		for i := 0; i < 5; i++ {
			select {
			case <-messageReceived:

			case <-time.After(1 * time.Second):
				t.Fatalf("Timed out waiting for message %d", i)
				return
			}
		}

		receivedMutex.Lock()
		assert.Equal(t, 5, receivedCount, "Should have received all 5 broadcast messages")
		receivedMutex.Unlock()
	})
}

func TestClientErrorWithMock(t *testing.T) {

	failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err == nil {
			conn.Close(websocket.StatusInternalError, "Test error")
		}
	}))
	defer failServer.Close()

	serverURL := strings.TrimPrefix(failServer.URL, "http://")
	wsURL := "ws://" + serverURL

	ctx := context.Background()
	client := New(ctx, wsURL)

	err := client.Listen(func(data []byte) error {
		return nil
	})
	assert.Error(t, err, "Listen should fail with not connected error")

	err = client.Write([]byte("test"))
	assert.Error(t, err, "Write should fail with not connected error")

	badClient := New(ctx, "ws://invalid-url-that-doesnt-exist:12345")
	err = badClient.Connect()
	assert.Error(t, err, "Connect should fail with invalid URL")
}

func TestClientReconnectWithMock(t *testing.T) {

	mock1 := newMockWSServer()
	defer mock1.close()

	serverURL1 := strings.TrimPrefix(mock1.server.URL, "http://")
	wsURL1 := "ws://" + serverURL1

	ctx := context.Background()
	client := New(ctx, wsURL1)

	err := client.Connect()
	require.NoError(t, err)

	client.Close()

	mock2 := newMockWSServer()
	defer mock2.close()

	connected := make(chan bool, 1)
	mock2.addMessageHandler(func(data []byte, conn *websocket.Conn) {
		var msg map[string]interface{}
		if err := json.Unmarshal(data, &msg); err != nil {
			return
		}

		if _, ok := msg["reconnect_test"]; ok {
			connected <- true
		}
	})

	serverURL2 := strings.TrimPrefix(mock2.server.URL, "http://")
	wsURL2 := "ws://" + serverURL2

	client2 := New(ctx, wsURL2)
	err = client2.Connect()
	require.NoError(t, err)
	defer client2.Close()

	testMessage := map[string]interface{}{
		"reconnect_test": true,
	}
	testData, _ := json.Marshal(testMessage)

	err = client2.Write(testData)
	require.NoError(t, err)

	select {
	case <-connected:

	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for connection verification")
	}
}
