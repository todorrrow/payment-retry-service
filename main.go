package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Payment struct {
	ID        string
	Amount    float64
	Currency  string
	Status    string
	Retries   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type createPaymentRequest struct {
	Amount float64 `json:"amount"`
}

type createPaymentResponse struct {
	RequestID int     `json:"request_id"`
	Amount    float64 `json:"amount"`
	Approved  bool    `json:"approved"`
	Message   string  `json:"message"`
	LatencyMS int64   `json:"latency_ms"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// serve the static html file from disk
	http.ServeFile(w, r, "static/index.html")
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	bankReq := BankRequest{ID: int(time.Now().UnixNano() % 1000000), Amount: req.Amount}
	results := ProcessBankRequests(ctx, []BankRequest{bankReq}, 1)
	if len(results) == 0 {
		http.Error(w, "request timed out or cancelled", http.StatusGatewayTimeout)
		return
	}

	res := results[0]
	out := createPaymentResponse{
		RequestID: res.RequestID,
		Amount:    req.Amount,
		Approved:  res.Approved,
		Message:   res.Message,
		LatencyMS: res.Latency.Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	// serve static assets under /static/
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/api/payments", paymentHandler)

	addr := ":8080"
	log.Printf("server listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
