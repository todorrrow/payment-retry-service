package main

import (
	"context"
	"testing"
	"time"
)

func TestSimulateBankRequest_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res := simulateBankRequest(ctx, BankRequest{ID: 42, Amount: 100.0})

	if res.RequestID != 42 {
		t.Fatalf("expected request id 42, got %d", res.RequestID)
	}
	if res.Approved {
		t.Fatalf("expected cancelled request to be not approved")
	}
	if res.Message != "cancelled" {
		t.Fatalf("expected message cancelled, got %q", res.Message)
	}
	if res.Latency < 0 {
		t.Fatalf("expected non-negative latency, got %s", res.Latency)
	}
}

func TestProcessBankRequests_ProcessesAllRequests(t *testing.T) {
	requests := []BankRequest{
		{ID: 1, Amount: 12.3},
		{ID: 2, Amount: 45.6},
		{ID: 3, Amount: 78.9},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	responses := ProcessBankRequests(ctx, requests, 2)
	if len(responses) != len(requests) {
		t.Fatalf("expected %d responses, got %d", len(requests), len(responses))
	}

	seen := make(map[int]bool, len(requests))
	for _, r := range responses {
		seen[r.RequestID] = true
		if r.Message == "" {
			t.Fatalf("expected non-empty message for request %d", r.RequestID)
		}
	}

	for _, req := range requests {
		if !seen[req.ID] {
			t.Fatalf("missing response for request %d", req.ID)
		}
	}
}

func TestProcessBankRequests_AlreadyCancelledContext(t *testing.T) {
	requests := []BankRequest{
		{ID: 1, Amount: 10},
		{ID: 2, Amount: 20},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	responses := ProcessBankRequests(ctx, requests, 0)
	if len(responses) != 0 {
		t.Fatalf("expected 0 responses for already-cancelled context, got %d", len(responses))
	}
}
