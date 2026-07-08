package main

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type BankRequest struct {
	ID     int
	Amount float64
}

type BankResponse struct {
	RequestID int
	Approved  bool
	Message   string
	Latency   time.Duration
}

func simulateBankRequest(ctx context.Context, req BankRequest) BankResponse {
	start := time.Now()

	// Simulate network and processing delay.
	delay := time.Duration(100+rand.Intn(500)) * time.Millisecond
	select {
	case <-ctx.Done():
		return BankResponse{
			RequestID: req.ID,
			Approved:  false,
			Message:   "cancelled",
			Latency:   time.Since(start),
		}
	case <-time.After(delay):
	}

	approved := rand.Intn(100) < 80
	msg := "declined"
	if approved {
		msg = "approved"
	}

	return BankResponse{
		RequestID: req.ID,
		Approved:  approved,
		Message:   msg,
		Latency:   time.Since(start),
	}
}

func worker(ctx context.Context, id int, jobs <-chan BankRequest, results chan<- BankResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-jobs:
			if !ok {
				return
			}
			_ = id // retained for potential worker-level logging/metrics.
			results <- simulateBankRequest(ctx, req)
		}
	}
}

func ProcessBankRequests(ctx context.Context, requests []BankRequest, workerCount int) []BankResponse {
	if workerCount <= 0 {
		workerCount = 1
	}

	jobs := make(chan BankRequest)
	results := make(chan BankResponse, len(requests))

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(ctx, i+1, jobs, results, &wg)
	}

	go func() {
		defer close(jobs)
		for _, req := range requests {
			select {
			case <-ctx.Done():
				return
			case jobs <- req:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	responses := make([]BankResponse, 0, len(requests))
	for res := range results {
		responses = append(responses, res)
	}

	return responses
}
