package errors

import (
	"errors"
	"fmt"
)

var (
	ErrValidation          = errors.New("validation error")
	ErrNetwork             = errors.New("network error")
	ErrAPI                 = errors.New("API error")
	ErrAuthentication      = errors.New("authentication error")
	ErrWebsocket           = errors.New("websocket error")
	ErrInternal            = errors.New("internal error")
	ErrTimeout             = errors.New("timeout error")
	ErrParameterError      = errors.New("parameter error")
	ErrRateLimitExceeded   = errors.New("rate limit exceeded")
	ErrSignatureError      = errors.New("signature error")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrOrderNotFound       = errors.New("order not found")
	ErrPositionNotExist    = errors.New("position not exist")
	ErrMarketNotExists     = errors.New("market not exists")
	ErrAccountNotAllowed   = errors.New("account not allowed to trade")
	ErrInvalidLeverage     = errors.New("invalid leverage")
	ErrTPSLOrderError      = errors.New("take profit/stop loss order error")
	ErrDuplicateClientID   = errors.New("client ID duplicate")
)

type ValidationError struct {
	Field   string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error: field %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

func (e *ValidationError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return ErrValidation
}

func (e *ValidationError) Is(target error) bool {
	return target == ErrValidation
}

type NetworkError struct {
	Operation string
	Message   string
	Err       error
}

func (e *NetworkError) Error() string {
	if e.Operation != "" {
		return fmt.Sprintf("network error during %s: %s", e.Operation, e.Message)
	}
	return fmt.Sprintf("network error: %s", e.Message)
}

func (e *NetworkError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return ErrNetwork
}

func (e *NetworkError) Is(target error) bool {
	return target == ErrNetwork
}

type APIError struct {
	Code     int
	Message  string
	Endpoint string
	Err      error
}

func (e *APIError) Error() string {
	if e.Endpoint != "" {
		return fmt.Sprintf("API error on %s: code=%d, message=%s", e.Endpoint, e.Code, e.Message)
	}
	return fmt.Sprintf("API error: code=%d, message=%s", e.Code, e.Message)
}

func (e *APIError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return ErrAPI
}

func (e *APIError) Is(target error) bool {
	if target == ErrAPI {
		return true
	}
	if e.Err != nil {
		return errors.Is(e.Err, target)
	}
	return false
}

type AuthenticationError struct {
	Message string
	Err     error
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("authentication error: %s", e.Message)
}

func (e *AuthenticationError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return ErrAuthentication
}

func (e *AuthenticationError) Is(target error) bool {
	return target == ErrAuthentication
}

type WebsocketError struct {
	Operation string
	Message   string
	Err       error
}

func (e *WebsocketError) Error() string {
	if e.Operation != "" {
		return fmt.Sprintf("websocket error during %s: %s", e.Operation, e.Message)
	}
	return fmt.Sprintf("websocket error: %s", e.Message)
}

func (e *WebsocketError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return ErrWebsocket
}

func (e *WebsocketError) Is(target error) bool {
	return target == ErrWebsocket
}

type InternalError struct {
	Message string
	Err     error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("internal error: %s", e.Message)
}

func (e *InternalError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return ErrInternal
}

func (e *InternalError) Is(target error) bool {
	return target == ErrInternal
}

type TimeoutError struct {
	Operation string
	Timeout   string
	Err       error
}

func (e *TimeoutError) Error() string {
	if e.Operation != "" && e.Timeout != "" {
		return fmt.Sprintf("timeout error during %s: exceeded %s", e.Operation, e.Timeout)
	} else if e.Operation != "" {
		return fmt.Sprintf("timeout error during %s", e.Operation)
	}
	return "timeout error"
}

func (e *TimeoutError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return ErrTimeout
}

func (e *TimeoutError) Is(target error) bool {
	return target == ErrTimeout
}

func NewValidationError(field, message string, err error) error {
	return &ValidationError{
		Field:   field,
		Message: message,
		Err:     err,
	}
}

func NewNetworkError(operation, message string, err error) error {
	return &NetworkError{
		Operation: operation,
		Message:   message,
		Err:       err,
	}
}

func NewAPIError(code int, message, endpoint string, err error) error {
	return &APIError{
		Code:     code,
		Message:  message,
		Endpoint: endpoint,
		Err:      err,
	}
}

func NewAuthenticationError(message string, err error) error {
	return &AuthenticationError{
		Message: message,
		Err:     err,
	}
}

func NewWebsocketError(operation, message string, err error) error {
	return &WebsocketError{
		Operation: operation,
		Message:   message,
		Err:       err,
	}
}

func NewInternalError(message string, err error) error {
	return &InternalError{
		Message: message,
		Err:     err,
	}
}

func NewTimeoutError(operation, timeout string, err error) error {
	return &TimeoutError{
		Operation: operation,
		Timeout:   timeout,
		Err:       err,
	}
}
