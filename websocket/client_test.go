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

		typ, msg, err := c.Read(ctx)
		require.NoError(t, err)
		assert.Equal(t, websocket.MessageText, typ)
		receivedLogin = msg
		loginCalled = true

		err = c.Write(ctx, websocket.MessageText, []byte(`{"event":"login","success":true}`))
		require.NoError(t, err)

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

	time.Sleep(200 * time.Millisecond)

	assert.True(t, loginCalled)
	assert.Equal(t, `{"op":"login","args":[{"key":"test"}]}`, string(receivedLogin))
}

func TestHeartbeat(t *testing.T) {
	heartbeatCh := make(chan []byte, 1)

	srv := setupWebsocketServer(t, func(c *websocket.Conn, ctx context.Context) {

		typ, msg, err := c.Read(ctx)
		require.NoError(t, err)
		assert.Equal(t, websocket.MessageText, typ)
		heartbeatCh <- msg
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

	time.Sleep(300 * time.Millisecond)

	var receivedHeartbeat []byte
	select {
	case receivedHeartbeat = <-heartbeatCh:

	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for heartbeat")
	}

	assert.Equal(t, `{"op":"ping","args":[]}`, string(receivedHeartbeat))
}

func TestWrite(t *testing.T) {
	heartbeatCh := make(chan []byte, 1)

	srv := setupWebsocketServer(t, func(c *websocket.Conn, ctx context.Context) {

		typ, msg, err := c.Read(ctx)
		require.NoError(t, err)
		assert.Equal(t, websocket.MessageText, typ)
		heartbeatCh <- msg
	})
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := New(ctx, wsURL)

	err := client.Connect()
	require.NoError(t, err)
	defer client.Close()

	time.Sleep(100 * time.Millisecond)

	testMessage := []byte(`{"op":"subscribe","args":[{"channel":"test"}]}`)
	err = client.Write(testMessage)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	var receivedHeartbeat []byte
	select {
	case receivedHeartbeat = <-heartbeatCh:

	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for heartbeat")
	}

	assert.Equal(t, string(testMessage), string(receivedHeartbeat))
}

func TestListen(t *testing.T) {
	listenTestFinished := make(chan struct{})

	srv := setupWebsocketServer(t, func(c *websocket.Conn, ctx context.Context) {

		_, _, err := c.Read(ctx)
		require.NoError(t, err)

		err = c.Write(ctx, websocket.MessageText, []byte(`{"event":"response","data":"test"}`))
		require.NoError(t, err)

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

	time.Sleep(100 * time.Millisecond)

	err = client.Write([]byte(`{"ping":true}`))
	require.NoError(t, err)

	var (
		handlerCalled  bool
		handlerMessage []byte
		listenComplete = make(chan struct{})
	)

	go func() {
		client.Listen(func(msg []byte) error {
			handlerCalled = true
			handlerMessage = msg
			close(listenComplete)
			return nil
		})
	}()

	select {
	case <-listenComplete:

	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for handler to be called")
	}

	close(listenTestFinished)

	cancel()

	assert.True(t, handlerCalled)

	var jsonMsg map[string]interface{}
	err = json.Unmarshal(handlerMessage, &jsonMsg)
	require.NoError(t, err)
	assert.Equal(t, "response", jsonMsg["event"])
	assert.Equal(t, "test", jsonMsg["data"])
}

func TestClose(t *testing.T) {
	srv := setupWebsocketServer(t, func(c *websocket.Conn, ctx context.Context) {

		for {
			_, _, err := c.Read(ctx)
			if err != nil {

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

	time.Sleep(100 * time.Millisecond)

	client.Close()

	client.conn = nil

	err = client.Write([]byte(`{"test":true}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}
