# Payment Retry Service

Lightweight Go service that simulates bank requests with retry-friendly processing and a minimal web UI.

## Requirements

- Go 1.20+ (module mode)

## Quick start

Run locally:

```bash
cd ../payment-retry-service
go run .
```

Open the UI in your browser:

- http://localhost:8080/

## API

POST /api/payments

Request (JSON):

```json
{ "amount": 49.99 }
```

Response (JSON):

```json
{
  "request_id": 12345,
  "amount": 49.99,
  "approved": true,
  "message": "approved",
  "latency_ms": 312
}
```

Example curl:

```bash
curl -s -X POST http://localhost:8080/api/payments \
  -H "Content-Type: application/json" \
  -d '{"amount":49.99}' | jq
```

## Project layout

- [main.go](main.go) — HTTP server and handlers
- [processor.go](processor.go) — request processing and concurrency logic
- [static/index.html](static/index.html) — minimal front-end
- [api_test.go](api_test.go), [processor_test.go](processor_test.go) — tests

## Tests

Run unit tests:

```bash
go test ./...
```

## Notes

- This project is intentionally minimal. If you want CI, persistent storage, or a Dockerfile, I can add them.
