package errors

import (
	"errors"
	"testing"
)

func TestValidationError(t *testing.T) {

	err := NewValidationError("field", "message", nil)

	expected := "validation error: field field: message"
	if err.Error() != expected {
		t.Errorf("Wrong error message. Expected '%s', got '%s'", expected, err.Error())
	}

	if !errors.Is(err, ErrValidation) {
		t.Error("errors.Is(err, ErrValidation) should be true")
	}

	originalErr := errors.New("original error")
	wrappedErr := NewValidationError("field", "message", originalErr)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("errors.Is(wrappedErr, originalErr) should be true")
	}

	if _, ok := wrappedErr.(*ValidationError); !ok {
		t.Error("Type assertion to *ValidationError should succeed")
	}
}

func TestAPIError(t *testing.T) {

	err := NewAPIError(400, "Bad request", "/api/endpoint", nil)

	expected := "API error on /api/endpoint: code=400, message=Bad request"
	if err.Error() != expected {
		t.Errorf("Wrong error message. Expected '%s', got '%s'", expected, err.Error())
	}

	if !errors.Is(err, UnknownAPIError) {
		t.Error("errors.Is(err, UnknownAPIError) should be true")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Error("Type assertion to *APIError should succeed")
	}

	if apiErr.Code != 400 {
		t.Errorf("Expected Code=400, got %d", apiErr.Code)
	}
	if apiErr.Message != "Bad request" {
		t.Errorf("Expected Message='Bad request', got '%s'", apiErr.Message)
	}
	if apiErr.Endpoint != "/api/endpoint" {
		t.Errorf("Expected Endpoint='/api/endpoint', got '%s'", apiErr.Endpoint)
	}
}

func TestNetworkError(t *testing.T) {

	err := NewNetworkError("HTTP request", "connection refused", nil)

	expected := "network error during HTTP request: connection refused"
	if err.Error() != expected {
		t.Errorf("Wrong error message. Expected '%s', got '%s'", expected, err.Error())
	}

	if !errors.Is(err, ErrNetwork) {
		t.Error("errors.Is(err, ErrNetwork) should be true")
	}

	if _, ok := err.(*NetworkError); !ok {
		t.Error("Type assertion to *NetworkError should succeed")
	}
}

func TestConnectionClosedError(t *testing.T) {
	err := NewConnectionClosedError("WebSocket", "remote server closed connection", nil)

	expected := "connection closed error during WebSocket: remote server closed connection"
	if err.Error() != expected {
		t.Errorf("Wrong error message. Expected '%s', got '%s'", expected, err.Error())
	}

	if !errors.Is(err, ErrConnectionClosed) {
		t.Error("errors.Is(err, ErrConnectionClosed) should be true")
	}

	originalErr := errors.New("original error")
	wrappedErr := NewConnectionClosedError("WebSocket", "connection error", originalErr)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("errors.Is(wrappedErr, originalErr) should be true")
	}

	if _, ok := err.(*ConnectionClosedError); !ok {
		t.Error("Type assertion to *ConnectionClosedError should succeed")
	}
}

func TestWorkgroupExhaustedError(t *testing.T) {
	err := NewWorkgroupExhaustedError("Queue Processing", "queue is full", nil)

	expected := "work group exhausted error during Queue Processing: queue is full"
	if err.Error() != expected {
		t.Errorf("Wrong error message. Expected '%s', got '%s'", expected, err.Error())
	}

	if !errors.Is(err, ErrWorkgroupExhausted) {
		t.Error("errors.Is(err, ErrWorkgroupExhausted) should be true")
	}

	originalErr := errors.New("original error")
	wrappedErr := NewWorkgroupExhaustedError("Queue Processing", "queue is full", originalErr)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("errors.Is(wrappedErr, originalErr) should be true")
	}

	if _, ok := err.(*WorkgroupExhaustedError); !ok {
		t.Error("Type assertion to *WorkgroupExhaustedError should succeed")
	}
}

func TestAuthenticationError(t *testing.T) {
	err := NewAuthenticationError("invalid API key", nil)

	expected := "authentication error: invalid API key"
	if err.Error() != expected {
		t.Errorf("Wrong error message. Expected '%s', got '%s'", expected, err.Error())
	}

	if !errors.Is(err, ErrAuthentication) {
		t.Error("errors.Is(err, ErrAuthentication) should be true")
	}

	originalErr := errors.New("original error")
	wrappedErr := NewAuthenticationError("invalid API key", originalErr)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("errors.Is(wrappedErr, originalErr) should be true")
	}

	if _, ok := err.(*AuthenticationError); !ok {
		t.Error("Type assertion to *AuthenticationError should succeed")
	}
}

func TestWebsocketError(t *testing.T) {
	err := NewWebsocketError("subscription", "failed to subscribe to topic", nil)

	expected := "websocket error during subscription: failed to subscribe to topic"
	if err.Error() != expected {
		t.Errorf("Wrong error message. Expected '%s', got '%s'", expected, err.Error())
	}

	if !errors.Is(err, ErrWebsocket) {
		t.Error("errors.Is(err, ErrWebsocket) should be true")
	}

	originalErr := errors.New("original error")
	wrappedErr := NewWebsocketError("subscription", "failed to subscribe to topic", originalErr)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("errors.Is(wrappedErr, originalErr) should be true")
	}

	if _, ok := err.(*WebsocketError); !ok {
		t.Error("Type assertion to *WebsocketError should succeed")
	}
}

func TestInternalError(t *testing.T) {
	err := NewInternalError("unexpected state", nil)

	expected := "internal error: unexpected state"
	if err.Error() != expected {
		t.Errorf("Wrong error message. Expected '%s', got '%s'", expected, err.Error())
	}

	if !errors.Is(err, ErrInternal) {
		t.Error("errors.Is(err, ErrInternal) should be true")
	}

	originalErr := errors.New("original error")
	wrappedErr := NewInternalError("unexpected state", originalErr)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("errors.Is(wrappedErr, originalErr) should be true")
	}

	if _, ok := err.(*InternalError); !ok {
		t.Error("Type assertion to *InternalError should succeed")
	}
}

func TestTimeoutError(t *testing.T) {

	err := NewTimeoutError("operation", "5s", nil)

	expected := "timeout error during operation: exceeded 5s"
	if err.Error() != expected {
		t.Errorf("Wrong error message. Expected '%s', got '%s'", expected, err.Error())
	}

	if !errors.Is(err, ErrTimeout) {
		t.Error("errors.Is(err, ErrTimeout) should be true")
	}

	if _, ok := err.(*TimeoutError); !ok {
		t.Error("Type assertion to *TimeoutError should succeed")
	}
}

func TestErrorHierarchy(t *testing.T) {

	originalErr := errors.New("original error")
	networkErr := NewNetworkError("HTTP", "connection error", originalErr)
	apiErr := NewAPIError(500, "Internal Server Error", "/api", networkErr)

	if !errors.Is(apiErr, UnknownAPIError) {
		t.Error("errors.Is(apiErr, UnknownAPIError) should be true")
	}
	if !errors.Is(apiErr, ErrNetwork) {
		t.Error("errors.Is(apiErr, ErrNetwork) should be true")
	}
	if !errors.Is(apiErr, originalErr) {
		t.Error("errors.Is(apiErr, originalErr) should be true")
	}

	unwrapped1 := errors.Unwrap(apiErr)
	if !errors.Is(unwrapped1, ErrNetwork) {
		t.Error("First unwrap should result in a NetworkError")
	}

	unwrapped2 := errors.Unwrap(unwrapped1)
	if unwrapped2 != originalErr {
		t.Error("Second unwrap should result in the original error")
	}
}
