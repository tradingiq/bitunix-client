package api

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"net/url"
)

// MockClient is a mock implementation of the API client for testing
type MockClient struct {
	mock.Mock
	Client
}

// Request mocks the Request method
func (m *MockClient) Request(ctx context.Context, method, path string, query url.Values, bodyBytes []byte) ([]byte, error) {
	args := m.Called(ctx, method, path, query, bodyBytes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// Get mocks the Get method
func (m *MockClient) Get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	args := m.Called(ctx, path, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// Post mocks the Post method
func (m *MockClient) Post(ctx context.Context, path string, query url.Values, body []byte) ([]byte, error) {
	args := m.Called(ctx, path, query, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// RequestError represents an API request error
type RequestError struct {
	Message    string
	StatusCode int
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

// NewRequestError creates a new API request error
func NewRequestError(message string, statusCode int) *RequestError {
	return &RequestError{
		Message:    message,
		StatusCode: statusCode,
	}
}
