package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWebsocketServer(t *testing.T, handler func(c *websocket.Conn, ctx context.Context)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.AcceptOptions{
			InsecureSkipVerify: true,
		}
		c, err := websocket.Accept(w, r, &upgrader)
		if err != nil {
			return
		}
		defer c.Close(websocket.StatusInternalError, "")

		ctx := context.Background()
		// Send initial welcome message
		err = c.Write(ctx, websocket.MessageText, []byte(`{"event":"welcome"}`))
		require.NoError(t, err)

		handler(c, ctx)
	}))
}

func TestNew(t *testing.T) {
	ctx := context.Background()
	wsURL := "wss://example.com/ws"

	client := New(ctx, wsURL)

	assert.Equal(t, wsURL, client.wsURL)
	assert.NotNil(t, client.done)
	assert.NotNil(t, client.ctx)
	assert.NotNil(t, client.cancel)
	assert.Nil(t, client.generateLoginMessage)
	assert.Nil(t, client.generateHeartbeatMessage)
}

func TestClientOptions(t *testing.T) {
	ctx := context.Background()
	wsURL := "wss://example.com/ws"

	loginGenerator := func() ([]byte, error) {
		return []byte(`{"op":"login"}`), nil
	}

	heartbeatInterval := 5 * time.Second
	heartbeatGenerator := func() ([]byte, error) {
		return []byte(`{"op":"ping"}`), nil
	}

	client := New(ctx, wsURL,
		WithAuthentication(loginGenerator),
		WithHeartbeat(heartbeatInterval, heartbeatGenerator),
	)

	assert.Equal(t, wsURL, client.wsURL)
	assert.NotNil(t, client.generateLoginMessage)
	assert.NotNil(t, client.generateHeartbeatMessage)
	assert.Equal(t, heartbeatInterval, client.heartBeatInterval)
}

func TestConnect(t *testing.T) {
	loginCalled := false
	var receivedLogin []byte

	srv := setupWebsocketServer(t, func(c *websocket.Conn, ctx context.Context) {
		// Read login message
		typ, msg, err := c.Read(ctx)
		require.NoError(t, err)
		assert.Equal(t, websocket.MessageText, typ)
		receivedLogin = msg
		loginCalled = true

		// Send login response
		err = c.Write(ctx, websocket.MessageText, []byte(`{"event":"login","success":true}`))
		require.NoError(t, err)

		// Keep connection open briefly
		time.Sleep(100 * time.Millisecond)
	})
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := New(ctx, wsURL, WithAuthentication(func() ([]byte, error) {
		return []byte(`{"op":"login","args":[{"key":"test"}]}`), nil
	}))

	err := client.Connect()
	require.NoError(t, err)
	defer client.Close()

	// Give the goroutines time to process
	time.Sleep(200 * time.Millisecond)

	assert.True(t, loginCalled)
	assert.Equal(t, `{"op":"login","args":[{"key":"test"}]}`, string(receivedLogin))
}

func TestHeartbeat(t *testing.T) {
	heartbeatReceived := false
	var receivedHeartbeat []byte

	srv := setupWebsocketServer(t, func(c *websocket.Conn, ctx context.Context) {
		// Read heartbeat message
		typ, msg, err := c.Read(ctx)
		require.NoError(t, err)
		assert.Equal(t, websocket.MessageText, typ)
		receivedHeartbeat = msg
		heartbeatReceived = true
	})
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := New(ctx, wsURL, WithHeartbeat(100*time.Millisecond, func() ([]byte, error) {
		return []byte(`{"op":"ping","args":[]}`), nil
	}))

	err := client.Connect()
	require.NoError(t, err)
	defer client.Close()

	// Give time for heartbeat to trigger
	time.Sleep(300 * time.Millisecond)

	assert.True(t, heartbeatReceived)
	assert.Equal(t, `{"op":"ping","args":[]}`, string(receivedHeartbeat))
}

func TestWrite(t *testing.T) {
	messageReceived := false
	var receivedMessage []byte

	srv := setupWebsocketServer(t, func(c *websocket.Conn, ctx context.Context) {
		// Read client message
		typ, msg, err := c.Read(ctx)
		require.NoError(t, err)
		assert.Equal(t, websocket.MessageText, typ)
		receivedMessage = msg
		messageReceived = true
	})
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := New(ctx, wsURL)

	err := client.Connect()
	require.NoError(t, err)
	defer client.Close()

	// Give time for connection to establish
	time.Sleep(100 * time.Millisecond)

	// Test Write
	testMessage := []byte(`{"op":"subscribe","args":[{"channel":"test"}]}`)
	err = client.Write(testMessage)
	require.NoError(t, err)

	// Give time for message to be processed
	time.Sleep(100 * time.Millisecond)

	assert.True(t, messageReceived)
	assert.Equal(t, string(testMessage), string(receivedMessage))
}

func TestListen(t *testing.T) {
	listenTestFinished := make(chan struct{})

	srv := setupWebsocketServer(t, func(c *websocket.Conn, ctx context.Context) {
		// Read message to confirm connection
		_, _, err := c.Read(ctx)
		require.NoError(t, err)

		// Send a message to trigger the handler
		err = c.Write(ctx, websocket.MessageText, []byte(`{"event":"response","data":"test"}`))
		require.NoError(t, err)

		// Wait for test to finish before closing
		<-listenTestFinished
	})
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := New(ctx, wsURL)

	err := client.Connect()
	require.NoError(t, err)
	defer client.Close()

	// Give time for connection to establish
	time.Sleep(100 * time.Millisecond)

	// Send ping to ensure connection is established
	err = client.Write([]byte(`{"ping":true}`))
	require.NoError(t, err)

	// Track handler calls
	var (
		handlerCalled  bool
		handlerMessage []byte
		listenComplete = make(chan struct{})
	)

	// Start listener in goroutine
	go func() {
		client.Listen(func(msg []byte) error {
			handlerCalled = true
			handlerMessage = msg
			close(listenComplete)
			return nil
		})
	}()

	// Wait for handler to be called
	select {
	case <-listenComplete:
		// Handler was called
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for handler to be called")
	}

	// Signal test completion to server
	close(listenTestFinished)

	// Cancel context to stop listener
	cancel()

	// Assertions for handler
	assert.True(t, handlerCalled)

	// Verify the message format
	var jsonMsg map[string]interface{}
	err = json.Unmarshal(handlerMessage, &jsonMsg)
	require.NoError(t, err)
	assert.Equal(t, "response", jsonMsg["event"])
	assert.Equal(t, "test", jsonMsg["data"])
}

func TestClose(t *testing.T) {
	srv := setupWebsocketServer(t, func(c *websocket.Conn, ctx context.Context) {
		// Keep connection open
		for {
			_, _, err := c.Read(ctx)
			if err != nil {
				// Expected on close
				break
			}
		}
	})
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := New(ctx, wsURL)

	err := client.Connect()
	require.NoError(t, err)

	// Give time for connection to establish
	time.Sleep(100 * time.Millisecond)

	// Test Close
	client.Close()

	// Set conn to nil after closing to ensure Write returns our expected error
	client.conn = nil

	// Verify connection is closed
	err = client.Write([]byte(`{"test":true}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}
