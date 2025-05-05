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

	err := NewAPIError(400, "Bad Request", "/api/endpoint", nil)

	expected := "API error on /api/endpoint: code=400, message=Bad Request"
	if err.Error() != expected {
		t.Errorf("Wrong error message. Expected '%s', got '%s'", expected, err.Error())
	}

	if !errors.Is(err, ErrAPI) {
		t.Error("errors.Is(err, ErrAPI) should be true")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Error("Type assertion to *APIError should succeed")
	}

	if apiErr.Code != 400 {
		t.Errorf("Expected Code=400, got %d", apiErr.Code)
	}
	if apiErr.Message != "Bad Request" {
		t.Errorf("Expected Message='Bad Request', got '%s'", apiErr.Message)
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

	if !errors.Is(apiErr, ErrAPI) {
		t.Error("errors.Is(apiErr, ErrAPI) should be true")
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
