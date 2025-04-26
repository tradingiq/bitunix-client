package model

import (
	"encoding/json"
	"testing"
)

func TestBalanceDetailUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"coin": "USDT",
		"available": "1000.50",
		"frozen": "100.25",
		"isolationFrozen": "20.75",
		"crossFrozen": "30.25",
		"margin": "200.50",
		"isolationMargin": "50.75",
		"crossMargin": "150.25",
		"expMoney": "1200.75"
	}`

	var balance BalanceDetail
	err := json.Unmarshal([]byte(jsonStr), &balance)
	if err != nil {
		t.Fatalf("Failed to unmarshal balance detail: %v", err)
	}

	if balance.Coin != "USDT" {
		t.Errorf("Expected coin 'USDT', got %s", balance.Coin)
	}

	if balance.Available != 1000.50 {
		t.Errorf("Expected available 1000.50, got %f", balance.Available)
	}

	if balance.Frozen != 100.25 {
		t.Errorf("Expected frozen 100.25, got %f", balance.Frozen)
	}

	if balance.IsolationFrozen != 20.75 {
		t.Errorf("Expected isolationFrozen 20.75, got %f", balance.IsolationFrozen)
	}

	if balance.CrossFrozen != 30.25 {
		t.Errorf("Expected crossFrozen 30.25, got %f", balance.CrossFrozen)
	}

	if balance.Margin != 200.50 {
		t.Errorf("Expected margin 200.50, got %f", balance.Margin)
	}

	if balance.IsolationMargin != 50.75 {
		t.Errorf("Expected isolationMargin 50.75, got %f", balance.IsolationMargin)
	}

	if balance.CrossMargin != 150.25 {
		t.Errorf("Expected crossMargin 150.25, got %f", balance.CrossMargin)
	}

	if balance.ExpMoney != 1200.75 {
		t.Errorf("Expected expMoney 1200.75, got %f", balance.ExpMoney)
	}

	// Test invalid JSON
	invalidJSON := `{
		"coin": "USDT",
		"available": "not-a-number",
		"frozen": "100.25",
		"isolationFrozen": "20.75",
		"crossFrozen": "30.25",
		"margin": "200.50",
		"isolationMargin": "50.75",
		"crossMargin": "150.25",
		"expMoney": "1200.75"
	}`

	err = json.Unmarshal([]byte(invalidJSON), &balance)
	if err == nil {
		t.Fatal("Expected error for invalid number format, got none")
	}
}

func TestBalanceResponseUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"ch": "balance",
		"ts": 1642521600000,
		"data": {
			"coin": "USDT",
			"available": "1000.50",
			"frozen": "100.25",
			"isolationFrozen": "20.75",
			"crossFrozen": "30.25",
			"margin": "200.50",
			"isolationMargin": "50.75",
			"crossMargin": "150.25",
			"expMoney": "1200.75"
		}
	}`

	var response BalanceResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal balance response: %v", err)
	}

	if response.Ch != "balance" {
		t.Errorf("Expected ch 'balance', got %s", response.Ch)
	}

	if response.Ts != 1642521600000 {
		t.Errorf("Expected ts 1642521600000, got %d", response.Ts)
	}

	if response.Data.Coin != "USDT" {
		t.Errorf("Expected data.coin 'USDT', got %s", response.Data.Coin)
	}

	if response.Data.Available != 1000.50 {
		t.Errorf("Expected data.available 1000.50, got %f", response.Data.Available)
	}

	if response.Data.Frozen != 100.25 {
		t.Errorf("Expected data.frozen 100.25, got %f", response.Data.Frozen)
	}
}
