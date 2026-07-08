package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPaymentHandler_Succeeds(t *testing.T) {
	rand.Seed(1) // deterministic for test

	mux := http.NewServeMux()
	mux.HandleFunc("/api/payments", paymentHandler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	body := `{"amount":12.34}`
	resp, err := http.Post(ts.URL+"/api/payments", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var out createPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if out.RequestID == 0 {
		t.Fatalf("expected non-zero request id")
	}
	if out.Amount != 12.34 {
		t.Fatalf("expected amount 12.34, got %v", out.Amount)
	}
	if out.Message == "" {
		t.Fatalf("expected non-empty message")
	}
}

func TestPaymentHandler_InvalidJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/payments", paymentHandler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/api/payments", "application/json", strings.NewReader("notjson"))
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid json, got %d", resp.StatusCode)
	}
}
