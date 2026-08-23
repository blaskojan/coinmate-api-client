# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

A Go HTTP client for the [Coinmate.io API](https://coinmate.docs.apiary.io/#). Requires Go 1.25+. The `main.go` at the root is a demo harness that exercises endpoints; the reusable client lives under `coinmate/`.

## Module layout

The Go module is named `coinmate` (see `go.mod`), and the library packages live under the `coinmate/` directory, so internal imports are `coinmate/coinmate`, `coinmate/coinmate/public`, and `coinmate/coinmate/secure` (module path + directory). The package name is `coinmate` while the import path is `coinmate/coinmate`.

## Commands

```bash
make test            # go test -v ./...
make test-coverage   # tests + HTML coverage -> coverage.html
make build           # go build -o main .
make run             # go run . (demo; reads COINMATE_* env vars)
make fmt             # go fmt ./...
make lint            # golangci-lint run
go test -v ./coinmate/secure/                       # single package
go test -v -run TestGetTickerSuccess ./coinmate/public/   # single test
```

Docker equivalents (`make docker-test`, `quick-test`, `test-docker`) run against `golang:1.21-alpine`. The demo reads `COINMATE_CLIENT_ID`, `COINMATE_API_KEY`, `COINMATE_PRIVATE_KEY` from the environment; secure endpoints are skipped when they are unset.

## Architecture

Three packages, split by authentication requirement:

- **`coinmate/`** (`client.go`) — the core `CoinmateClient` and the `ClientInterface` it implements. Owns HTTP transport, base URL, nonce generation, and HMAC signing.
- **`coinmate/public/`** — endpoints needing no credentials (ticker, order book, trading pairs, transactions, currencies, server time).
- **`coinmate/secure/`** — authenticated endpoints (balances, orders). Requests are POST with a signed form body.

**`ClientInterface` is the key abstraction.** Every endpoint is a struct that holds `Client coinmate.ClientInterface` (e.g. `Ticker{Client: client}`, `Order{Client: client}`), and calls the client only through this interface. This is what makes endpoints unit-testable without network access — tests inject a `MockClient` implementing the same interface.

**Authentication flow (secure requests):** `GetRequestBody` builds a form body containing `clientId`, `publicKey`, `nonce`, and `signature`. The signature is `HMAC-SHA256(privateKey, nonce+clientId+apiKey)` uppercased (`GetSignature`). Nonces are monotonic (`GetNonce` bumps past `lastNonce`) so concurrent-ish calls don't collide. Secure requests send `Content-Type: application/x-www-form-urlencoded`.

## Endpoint implementation pattern

When adding an endpoint, mirror the existing files (`coinmate/public/ticker.go`, `coinmate/secure/order.go`):

1. Struct `{Client coinmate.ClientInterface}` with an exported method per API call.
2. Response structs with `json:"..."` tags matching the API; the standard shape is `{ Error bool, ErrorMessage string, Data ... }`.
3. Validate inputs early and return a wrapped error (e.g. empty `currencyPair`).
4. Build a `coinmate.Request`, call `MakePublicRequest` / `MakeSecureRequest`.
5. Check `err`, then check `response.StatusCode != http.StatusOK`, then `json.Unmarshal`. **Wrap every error with `fmt.Errorf(..., %w)`** — this is an established convention across the codebase.
6. Endpoint paths and param names are declared as package-level `const` (see the block at the top of `order.go`).

## Testing pattern

Each package defines its own `MockClient` embedding `coinmate.ClientInterface` and overriding the request methods to return a canned `coinmate.Response` or error (see `coinmate/public/ticker_test.go`). Tests cover success, API-error response (`error: true`), non-200 HTTP status, network error, and invalid JSON. Follow this same set of cases for new endpoints.

## Status

Many secure endpoints (withdrawals/deposits, trade history, transfers, replace/cancel-all orders, get-order-by-id) are not yet implemented — see README.md for the current list. The `GetRequestBody` in `client.go` still prints the encoded body to stdout (`fmt.Println`), a debug leftover to be aware of when working near request signing.
