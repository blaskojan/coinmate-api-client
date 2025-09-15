# Coinmate API Go Client

HTTP Go client to communicate with [Coinmate.io API](https://coinmate.docs.apiary.io/#)

## Current Implementation Status

### ✅ Implemented Endpoints

**Public Endpoints:**
- `/ticker` - Get ticker data for a currency pair
- `/orderBook` - Get order book with optional grouping
- `/tradingPairs` - Get available trading pairs
- `/transactions` - Get recent transactions

**Secure Endpoints:**
- `/balances` - Get account balances
- `/orderHistory` - Get order history
- `/openOrders` - Get open orders
- `/cancelOrder` - Cancel order
- `/cancelOrderWithInfo` - Cancel order with detailed info
- `/buyLimit` - Place buy limit order
- `/sellLimit` - Place sell limit order
- `/buyInstant` - Place buy instant order
- `/sellInstant` - Place sell instant order

### ❌ Missing Endpoints

**Public Endpoints:**
- `/currencies` - Get available currencies
- `/currency-pairs` - Get currency pairs
- `/ticker-all` - Get all tickers
- `/system/get-server-time` - Get server time

**Secure Endpoints:**
- `/trader-fees` - Get trading fees
- `/trade-history` - Get trade history
- `/transaction-history` - Get transaction history
- `/transfers` - Transfer management
- `/order/get-order-by-orderid` - Get order by ID
- `/order/get-order-by-clientorderid` - Get order by client order ID
- `/order/replace-existing-order-by-*` - Replace existing orders
- `/order/cancel-all-open-orders` - Cancel all open orders

**Withdrawal/Deposit Endpoints (Completely Missing):**
- Bitcoin withdrawal/deposit operations
- Ethereum withdrawal/deposit operations
- Litecoin withdrawal/deposit operations
- Ripple withdrawal/deposit operations
- Cardano withdrawal/deposit operations
- Solana withdrawal/deposit operations
- USDT withdrawal/deposit operations
- Virtual currency withdrawal/deposit operations
- Fiat withdrawal operations

## Known Issues

1. **Typo**: `TraidingPairs` should be `TradingPairs`
2. **Incorrect method receiver**: `GetTradingPairs()` uses wrong receiver
3. **Missing error handling**: Some endpoints lack proper error handling
4. **Inconsistent response structures**: Some responses don't match API documentation
5. **Missing request validation**: No validation for required parameters
6. **Hardcoded timeout**: 2-second timeout might be insufficient

## Usage

```go
// Create client
client := coinmate.GetCoinmateClient(clientId, apiKey, privateKey)

// Public endpoints
ticker := &public.Ticker{Client: client}
tickerData, err := ticker.GetTicker("BTC_EUR")

// Secure endpoints
balances := &secure.Balances{Client: client}
balanceData, err := balances.GetBalances()
```

## Running tests

You can run tests locally (requires Go 1.25+) or inside Docker.

### Local

- **Quick run:**

```bash
make test
```

- **With coverage report:**

```bash
make test-coverage
# Opens HTML reports in coverage/*.html
```

- **Direct Go command:**

```bash
go test -v ./...
```

- **Helper script (detailed, per-package coverage):**

```bash
./run_tests_simple.sh
```

### Docker

- **One-shot tests in container:**

```bash
make docker-test
```

- **Quick tests (no coverage HTML):**

```bash
make quick-test
```

- **Alternative Docker test with coverage:**

```bash
make test-docker
```

- **Helper script (Docker, per-package coverage):**

```bash
./run_tests_docker_simple.sh
```

Coverage HTML reports (when enabled) are written to `coverage/client.html`, `coverage/public.html`, and `coverage/secure.html`.

## Contributing

This client is incomplete compared to the official API documentation. Contributions to add missing endpoints are welcome.