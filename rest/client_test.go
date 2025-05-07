package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestClientOptions(t *testing.T) {
	client := &client{
		httpClient: &http.Client{},
	}

	expectedTimeout := 30 * time.Second
	WithDefaultTimeout(expectedTimeout)(client)
	if client.httpClient.Timeout != expectedTimeout {
		t.Errorf("WithDefaultTimeout option did not set timeout correctly, expected %v, got %v",
			expectedTimeout, client.httpClient.Timeout)
	}

	called := false
	testSigner := func(req *http.Request, body []byte) error {
		called = true
		return nil
	}
	WithRequestSigner(testSigner)(client)

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	client.signRequest(req, nil)
	if !called {
		t.Error("WithRequestSigner option did not set request signer correctly")
	}
}

func TestNewClient(t *testing.T) {
	_, err := New("https://api.example.com")
	if err != nil {
		t.Fatalf("New client with valid URL returned error: %v", err)
	}

	_, err = New("://invalid")
	if err == nil {
		t.Fatal("New client with invalid URL should return error")
	}
}

func TestClientRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test/path" {
			t.Errorf("Expected path /test/path, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("param") != "value" {
			t.Errorf("Expected query param=value, got %s", r.URL.Query().Get("param"))
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header application/json, got %s",
				r.Header.Get("Content-Type"))
		}
		if r.Header.Get("test-header") != "test-value" {
			t.Errorf("Expected test-header header test-value, got %s",
				r.Header.Get("test-header"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	client, err := New(server.URL, WithRequestSigner(func(req *http.Request, body []byte) error {
		req.Header.Set("test-header", "test-value")
		return nil
	}))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	query := url.Values{}
	query.Set("param", "value")

	resp, err := client.Get(context.Background(), "/test/path", query)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	if string(resp) != `{"success":true}` {
		t.Errorf("Expected response {\"success\":true}, got %s", string(resp))
	}

	body := []byte(`{"data":"test"}`)
	resp, err = client.Post(context.Background(), "/test/path", query, body)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	if string(resp) != `{"success":true}` {
		t.Errorf("Expected response {\"success\":true}, got %s", string(resp))
	}
}

func TestClientRequestErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Test error"}`))
	}))
	defer server.Close()

	client, err := New(server.URL, WithRequestSigner(func(req *http.Request, body []byte) error {
		return nil
	}))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Get(context.Background(), "/test/path", nil)
	if err == nil {
		t.Fatal("Expected error for 400 response, got nil")
	}
}
