# Phase 3D Technical Implementation Specifications

**Version:** 1.0  
**Date:** November 24, 2025  
**Purpose:** Detailed technical specifications for Phase 3D feature implementation

---

## 🏗️ Architecture Overview

### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Next.js)                        │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  React Components │ Zustand Store │ API Client │ WebSocket   │
│  │  (Real-time UI)   │  (State Mgmt)  │  (HTTP)    │ (Real-time) │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              ↕ (REST + WebSocket)
┌─────────────────────────────────────────────────────────────────┐
│                    API Gateway / Load Balancer                   │
│                      (NGINX / Cloud Load)                        │
└─────────────────────────────────────────────────────────────────┘
                              ↕
┌─────────────────────────────────────────────────────────────────┐
│                    Backend (Go + Gorilla)                        │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  HTTP Handlers │ WebSocket Handler │ Authentication       │   │
│  │  (REST API)    │ (Real-time)       │ (JWT + Tenant)      │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Services Layer                                          │   │
│  │  ├─ Payment Service (Stripe, PayPal)                   │   │
│  │  ├─ Billing Service (Invoicing, Subscriptions)        │   │
│  │  ├─ Analytics Service (Events, Metrics)               │   │
│  │  ├─ Module Service                                     │   │
│  │  └─ Notification Service (WebSocket Broadcast)        │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Caching Layer (Redis)                                   │   │
│  │  ├─ Session Cache (TTL: 24h)                           │   │
│  │  ├─ User Preferences (TTL: 1h)                         │   │
│  │  ├─ Module Metadata (TTL: 6h)                          │   │
│  │  └─ Analytics Aggregates (TTL: 1h)                     │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  Message Queue (Redis Streams / RabbitMQ)               │   │
│  │  ├─ Payment Processing                                  │   │
│  │  ├─ Invoice Generation                                  │   │
│  │  ├─ Email Notifications                                 │   │
│  │  └─ Analytics Events                                    │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              ↕
┌─────────────────────────────────────────────────────────────────┐
│                     Data Layer                                   │
│  ┌────────────────────────────────────────────────────────┐    │
│  │  MySQL (Primary Database) - Main OLTP Database         │    │
│  │  ├─ Transactions                                      │    │
│  │  ├─ Multi-tenant Data                                │    │
│  │  ├─ User & Company Info                              │    │
│  │  └─ Billing & Invoices                               │    │
│  └────────────────────────────────────────────────────────┘    │
│  ┌────────────────────────────────────────────────────────┐    │
│  │  Redis (Caching & Pub/Sub)                            │    │
│  │  ├─ Session Cache                                    │    │
│  │  ├─ Real-time Pub/Sub                                │    │
│  │  └─ Message Queues                                   │    │
│  └────────────────────────────────────────────────────────┘    │
│  ┌────────────────────────────────────────────────────────┐    │
│  │  Analytics Database (Optional: Elasticsearch/Big Query) │    │
│  │  └─ Event data for analytics                          │    │
│  └────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              ↕
┌─────────────────────────────────────────────────────────────────┐
│                   External Services                              │
│  ├─ Stripe API (Payment Processing)                           │
│  ├─ PayPal API (Payment Processing)                           │
│  ├─ OpenExchangeRates API (Currency Exchange)                 │
│  ├─ SendGrid (Email)                                          │
│  ├─ DataDog (Monitoring)                                      │
│  └─ Sentry (Error Tracking)                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 💳 Payment Integration Architecture

### Stripe Integration Flow

```go
// File: internal/services/payment/stripe.go

package payment

import (
    "github.com/stripe/stripe-go/v76"
    "github.com/stripe/stripe-go/v76/webhook"
)

// StripeService handles all Stripe interactions
type StripeService struct {
    apiKey    string
    webhookKey string
    db        *db.Connection
    logger    *log.Logger
}

// CreatePaymentIntent creates a Stripe payment intent
func (s *StripeService) CreatePaymentIntent(ctx context.Context, req PaymentRequest) (*stripe.PaymentIntent, error) {
    // 1. Validate request
    // 2. Convert amount to cents
    // 3. Call Stripe API
    // 4. Store transaction record
    // 5. Return payment intent
}

// HandleWebhook processes Stripe webhook events
func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
    // 1. Verify webhook signature
    // 2. Parse event
    // 3. Route to appropriate handler
    // 4. Store event log
    // 5. Trigger follow-up actions
    
    // Event handlers:
    // - payment_intent.succeeded → Confirm invoice
    // - payment_intent.payment_failed → Trigger retry
    // - charge.refunded → Process refund
    // - customer.subscription.created → Activate subscription
}

// RetryFailedPayment implements exponential backoff retry
func (s *StripeService) RetryFailedPayment(ctx context.Context, transactionID string) error {
    // 1. Get transaction details
    // 2. Validate retry conditions
    // 3. Increment retry count
    // 4. Attempt retry
    // 5. Schedule next retry if failed (24h, 48h, 72h)
}

// ReconcilePayments reconciles Stripe with local database
func (s *StripeService) ReconcilePayments(ctx context.Context, startDate time.Time) error {
    // 1. Get all Stripe charges since startDate
    // 2. Compare with local transactions
    // 3. Identify discrepancies
    // 4. Log reconciliation report
    // 5. Alert on major discrepancies
}
```

### Database Schema for Payments

```sql
-- Payment-related tables

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    invoice_id UUID NOT NULL,
    company_id UUID NOT NULL,
    user_id UUID,
    amount DECIMAL(15, 2) NOT NULL,
    currency CHAR(3) DEFAULT 'INR',
    status ENUM('pending', 'processing', 'succeeded', 'failed', 'refunded') DEFAULT 'pending',
    payment_method ENUM('stripe', 'paypal', 'bank_transfer') NOT NULL,
    payment_processor_id VARCHAR(100),
    provider_transaction_id VARCHAR(100) UNIQUE,
    retry_count INT DEFAULT 0,
    last_retry_at TIMESTAMP,
    next_retry_at TIMESTAMP,
    error_message TEXT,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_invoice_id (invoice_id),
    INDEX idx_company_id (company_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);

CREATE TABLE payment_methods (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL,
    user_id UUID NOT NULL,
    type ENUM('card', 'bank_account', 'paypal') NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_id VARCHAR(100) UNIQUE NOT NULL,
    last_four VARCHAR(4),
    brand VARCHAR(50),
    exp_month INT,
    exp_year INT,
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_company_id (company_id),
    INDEX idx_user_id (user_id),
    UNIQUE KEY unique_default (company_id, is_default) WHERE is_default = TRUE
);

CREATE TABLE refunds (
    id UUID PRIMARY KEY,
    transaction_id UUID NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    reason ENUM('requested', 'chargeback', 'refund') NOT NULL,
    status ENUM('pending', 'processing', 'succeeded', 'failed') DEFAULT 'pending',
    provider_refund_id VARCHAR(100) UNIQUE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_status (status),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

CREATE TABLE payment_webhooks (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_event_id VARCHAR(100) UNIQUE NOT NULL,
    payload JSON NOT NULL,
    status ENUM('received', 'processing', 'completed', 'failed') DEFAULT 'received',
    error_message TEXT,
    retry_count INT DEFAULT 0,
    next_retry_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_event_type (event_type),
    INDEX idx_provider (provider),
    INDEX idx_created_at (created_at)
);
```

### API Endpoints for Payments

```go
// Payment API endpoints

// 1. Create payment intent
POST /api/v1/payments/create-intent
Request: {
    "invoice_id": "uuid",
    "amount": 10000,
    "currency": "INR",
    "payment_method_id": "uuid"
}
Response: {
    "client_secret": "pi_xxx",
    "status": "requires_action"
}

// 2. Get payment status
GET /api/v1/payments/{transaction_id}
Response: {
    "id": "uuid",
    "status": "succeeded",
    "amount": 10000,
    "created_at": "2025-11-24T10:00:00Z"
}

// 3. Refund payment
POST /api/v1/payments/{transaction_id}/refund
Request: {
    "amount": 10000,
    "reason": "customer_request",
    "notes": "Customer requested refund"
}
Response: {
    "refund_id": "uuid",
    "status": "processing"
}

// 4. List payment methods
GET /api/v1/payment-methods
Response: [
    {
        "id": "uuid",
        "type": "card",
        "last_four": "4242",
        "exp_month": 12,
        "exp_year": 2025,
        "is_default": true
    }
]

// 5. Add payment method
POST /api/v1/payment-methods
Request: {
    "type": "card",
    "token": "tok_xxx" // from Stripe.js
}
Response: {
    "id": "uuid",
    "type": "card",
    "last_four": "4242"
}

// 6. Webhook handler
POST /api/v1/webhooks/stripe
Headers: {
    "Stripe-Signature": "xxx"
}
Body: Stripe event JSON
Response: {"received": true}
```

---

## 🔌 WebSocket Real-Time Architecture

### WebSocket Server Implementation

```go
// File: internal/websocket/hub.go

package websocket

import (
    "sync"
    "github.com/gorilla/websocket"
)

// Hub maintains active client connections
type Hub struct {
    clients      map[*Client]bool
    broadcast    chan *Message
    register     chan *Client
    unregister   chan *Client
    mu           sync.RWMutex
    presenceMu   sync.RWMutex
    presence     map[string]bool // company_id -> online
    redis        *redis.Client
}

// Client represents a WebSocket connection
type Client struct {
    hub       *Hub
    conn      *websocket.Conn
    send      chan *Message
    companyID string
    userID    string
    sessionID string
}

// Message represents a WebSocket message
type Message struct {
    Type      string      `json:"type"` // "invoice", "payment", "alert"
    CompanyID string      `json:"company_id"`
    UserID    string      `json:"user_id,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
    Data      interface{} `json:"data"`
}

// RunHub runs the hub event loop
func (h *Hub) Run(ctx context.Context) {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            h.presenceMu.Lock()
            h.presence[client.companyID] = true
            h.presenceMu.Unlock()
            
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()
            
        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                if client.companyID == message.CompanyID {
                    select {
                    case client.send <- message:
                    default:
                        // Channel full, skip message
                    }
                }
            }
            h.mu.RUnlock()
            
            // Also publish to Redis for horizontal scaling
            h.redis.Publish(ctx, "channel:"+message.CompanyID, message)
        }
    }
}

// UpgradeConnection upgrades HTTP to WebSocket
func (h *Hub) UpgradeConnection(w http.ResponseWriter, r *http.Request) error {
    // 1. Extract JWT token
    token := r.Header.Get("Authorization")
    claims, err := auth.ValidateToken(token)
    if err != nil {
        return err
    }
    
    // 2. Upgrade connection
    upgrader := websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true // Validate origin in production
        },
    }
    
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return err
    }
    
    // 3. Create client
    client := &Client{
        hub:       h,
        conn:      conn,
        send:      make(chan *Message, 256),
        companyID: claims.CompanyID,
        userID:    claims.UserID,
        sessionID: uuid.New().String(),
    }
    
    // 4. Register client
    h.register <- client
    
    // 5. Start read/write loops
    go client.readPump()
    go client.writePump()
    
    return nil
}
```

### Real-Time Event Types

```go
// Real-time event definitions

const (
    // Invoice events
    EventInvoiceCreated = "invoice.created"
    EventInvoiceUpdated = "invoice.updated"
    EventInvoicePaid    = "invoice.paid"
    
    // Payment events
    EventPaymentStarted   = "payment.started"
    EventPaymentSucceeded = "payment.succeeded"
    EventPaymentFailed    = "payment.failed"
    
    // Module events
    EventModuleSubscribed   = "module.subscribed"
    EventModuleUnsubscribed = "module.unsubscribed"
    EventModuleUsageAlert   = "module.usage_alert"
    
    // System events
    EventSystemAlert = "system.alert"
    EventUserOnline  = "user.online"
    EventUserOffline = "user.offline"
)

// Event payload examples
type InvoiceCreatedEvent struct {
    InvoiceID string
    CompanyID string
    Amount    decimal.Decimal
    Currency  string
    CreatedAt time.Time
}

type PaymentSucceededEvent struct {
    TransactionID string
    InvoiceID     string
    CompanyID     string
    Amount        decimal.Decimal
    Timestamp     time.Time
}

type ModuleUsageAlertEvent struct {
    CompanyID string
    ModuleID  string
    Usage     int
    Threshold int
    AlertType string // "warning", "critical"
}
```

---

## 📊 Analytics Architecture

### Analytics Data Collection

```go
// File: internal/analytics/collector.go

package analytics

import "context"

// EventCollector collects usage events
type EventCollector struct {
    db    *db.Connection
    redis *redis.Client
    queue chan *Event
}

// Event represents an analytics event
type Event struct {
    ID        string
    Type      string // "subscription", "usage", "payment", "error"
    CompanyID string
    UserID    string
    ModuleID  string
    Metadata  map[string]interface{}
    Timestamp time.Time
}

// TrackEvent records an event
func (ec *EventCollector) TrackEvent(ctx context.Context, event *Event) error {
    event.Timestamp = time.Now()
    event.ID = uuid.New().String()
    
    // 1. Add to async queue (non-blocking)
    select {
    case ec.queue <- event:
    default:
        // Queue full, log warning but don't block
        log.Warn("analytics queue full, dropping event")
    }
    
    // 2. Return immediately
    return nil
}

// ProcessQueue processes queued events in batches
func (ec *EventCollector) ProcessQueue(ctx context.Context) {
    batch := make([]*Event, 0, 100)
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case event := <-ec.queue:
            batch = append(batch, event)
            if len(batch) >= 100 {
                ec.flushBatch(ctx, batch)
                batch = make([]*Event, 0, 100)
            }
            
        case <-ticker.C:
            if len(batch) > 0 {
                ec.flushBatch(ctx, batch)
                batch = make([]*Event, 0, 100)
            }
        }
    }
}

// flushBatch writes batch to database
func (ec *EventCollector) flushBatch(ctx context.Context, batch []*Event) error {
    // Insert all events in one transaction
    // Also update aggregated metrics in Redis
}
```

### Analytics Database Schema

```sql
CREATE TABLE analytics_events (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    company_id UUID NOT NULL,
    user_id UUID,
    module_id UUID,
    metadata JSON,
    timestamp DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_company_timestamp (company_id, timestamp),
    INDEX idx_event_type (event_type),
    INDEX idx_module_id (module_id)
);

CREATE TABLE analytics_daily_metrics (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(15, 4) NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_metric (company_id, metric_name, date),
    INDEX idx_company_date (company_id, date)
);

CREATE TABLE analytics_reports (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL,
    report_name VARCHAR(255) NOT NULL,
    report_type VARCHAR(50) NOT NULL, -- "revenue", "usage", "custom"
    query_params JSON,
    results JSON,
    generated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_company_created (company_id, created_at)
);
```

---

## 🔴 Caching Strategy (Redis)

### Redis Caching Implementation

```go
// File: pkg/cache/cache.go

package cache

import "github.com/redis/go-redis/v9"

// CacheManager manages all caching
type CacheManager struct {
    redis *redis.Client
}

// CachingPatterns defines TTLs and invalidation rules
const (
    // TTLs
    SessionTTL        = 24 * time.Hour
    UserPrefsTTL      = 1 * time.Hour
    ModuleMetaTTL     = 6 * time.Hour
    AnalyticsTTL      = 1 * time.Hour
    PermissionCacheTTL = 30 * time.Minute
    
    // Cache key prefixes
    SessionPrefix        = "session:"
    UserPrefsPrefix      = "user_prefs:"
    ModuleMetaPrefix     = "module_meta:"
    CompanyMetricsPrefix = "company_metrics:"
)

// GetOrFetch gets value from cache or fetches from source
func (cm *CacheManager) GetOrFetch(ctx context.Context, key string, ttl time.Duration, fetchFunc func() (interface{}, error)) (interface{}, error) {
    // 1. Try to get from cache
    val, err := cm.redis.Get(ctx, key).Result()
    if err == nil {
        return val, nil
    }
    
    // 2. Cache miss, fetch from source
    value, err := fetchFunc()
    if err != nil {
        return nil, err
    }
    
    // 3. Store in cache
    cm.redis.Set(ctx, key, value, ttl)
    
    return value, nil
}

// InvalidatePattern invalidates all keys matching pattern
func (cm *CacheManager) InvalidatePattern(ctx context.Context, pattern string) error {
    iter := cm.redis.Scan(ctx, 0, pattern, 1000).Iterator()
    
    for iter.Next(ctx) {
        cm.redis.Del(ctx, iter.Val())
    }
    
    return iter.Err()
}

// Cache invalidation triggers on write operations
func (cm *CacheManager) OnInvoiceCreated(ctx context.Context, companyID string) {
    cm.InvalidatePattern(ctx, CompanyMetricsPrefix+companyID+":*")
}

func (cm *CacheManager) OnPaymentProcessed(ctx context.Context, companyID string) {
    cm.InvalidatePattern(ctx, CompanyMetricsPrefix+companyID+":*")
}
```

---

## 🎯 Performance Optimization Details

### Query Optimization Example

```sql
-- BEFORE: N+1 query problem
SELECT * FROM companies WHERE tenant_id = ?;
-- Then for each company: SELECT SUM(amount) FROM invoices WHERE company_id = ?;

-- AFTER: Single optimized query with JOIN and GROUP BY
SELECT 
    c.id,
    c.name,
    COUNT(i.id) as invoice_count,
    COALESCE(SUM(i.amount), 0) as total_billed,
    COALESCE(SUM(CASE WHEN i.status = 'paid' THEN i.amount ELSE 0 END), 0) as total_paid
FROM companies c
LEFT JOIN invoices i ON c.id = i.company_id
WHERE c.tenant_id = ?
GROUP BY c.id
ORDER BY total_billed DESC;

-- Added indexes:
CREATE INDEX idx_invoices_company_status ON invoices(company_id, status);
CREATE INDEX idx_companies_tenant ON companies(tenant_id);
```

### Frontend Performance Optimization

```typescript
// File: frontend/hooks/useQuery.ts

export function useQuery<T>(
  queryKey: string[],
  queryFn: () => Promise<T>,
  options?: QueryOptions<T>
) {
  // 1. Check Zustand store cache
  const cachedData = useStore((state) => state.queryCache?.[queryKey.join('/')])
  
  // 2. Return cached if valid
  if (cachedData && !isStale(cachedData)) {
    return { data: cachedData, isLoading: false }
  }
  
  // 3. Use stale-while-revalidate pattern
  const [data, setData] = useState(cachedData?.data)
  const [isLoading, setIsLoading] = useState(!cachedData)
  
  useEffect(() => {
    let cancelled = false
    
    queryFn()
      .then((result) => {
        if (!cancelled) {
          setData(result)
          // Update cache with new data
          useStore.setState({
            queryCache: {
              ...useStore.getState().queryCache,
              [queryKey.join('/')]: {
                data: result,
                timestamp: Date.now(),
              },
            },
          })
        }
      })
      .finally(() => !cancelled && setIsLoading(false))
    
    return () => { cancelled = true }
  }, [queryKey.join('/')])
  
  return { data, isLoading }
}
```

---

## 📚 Integration Checklists

### Stripe Integration Checklist

```
Development:
  ☐ Set up Stripe test account
  ☐ Generate API keys (publish, secret)
  ☐ Install stripe-go SDK
  ☐ Implement PaymentIntent creation
  ☐ Implement webhook handlers
  ☐ Implement retry logic
  ☐ Write unit tests (80%+ coverage)
  ☐ Write integration tests with Stripe test API
  ☐ Test refund flow
  ☐ Test dispute handling

Testing:
  ☐ Test with Stripe test cards
  ☐ Test webhook failures and retries
  ☐ Test webhook signature verification
  ☐ Load test with 1000 concurrent payments
  ☐ Test edge cases (expired cards, declined, etc)

Staging:
  ☐ Deploy to staging
  ☐ Enable Stripe test webhooks
  ☐ Do end-to-end test
  ☐ Monitor logs for errors
  ☐ Verify reconciliation

Production:
  ☐ Switch to live Stripe account
  ☐ Update API keys in production
  ☐ Enable live webhooks
  ☐ Verify Stripe dashboard configuration
  ☐ Monitor first 100 transactions
  ☐ Verify reconciliation working
  ☐ Set up alerting for payment failures
```

### WebSocket Integration Checklist

```
Backend:
  ☐ Implement WebSocket upgrade handler
  ☐ Implement message routing
  ☐ Implement Redis Pub/Sub integration
  ☐ Implement connection pooling
  ☐ Implement heartbeat/ping-pong
  ☐ Write integration tests
  ☐ Test with 100+ concurrent connections
  ☐ Test message broadcast
  ☐ Test connection drops and reconnection
  ☐ Load test with 1000+ connections

Frontend:
  ☐ Create WebSocket connection manager
  ☐ Implement auto-reconnect logic
  ☐ Implement message handlers for all event types
  ☐ Create UI notifications for events
  ☐ Implement presence tracking UI
  ☐ Write component tests
  ☐ Test browser reconnection
  ☐ Test browser tab synchronization

Integration:
  ☐ End-to-end test: emit event, receive in UI
  ☐ Test with slow network (DevTools throttling)
  ☐ Test with connection drops (DevTools offline)
  ☐ Test across multiple browser tabs
  ☐ Performance test with 100+ events/second
  ☐ Verify no memory leaks (long-running connections)
```

---

## 🚨 Error Handling & Retry Strategies

### Payment Retry Logic

```go
// RetryStrategy defines how to retry failed operations
type RetryStrategy struct {
    MaxAttempts int
    BackoffFunc func(attempt int) time.Duration
}

// ExponentialBackoff returns exponential backoff durations
func ExponentialBackoff(attempt int) time.Duration {
    // 1 min, 2 min, 4 min, 8 min
    return time.Duration(math.Pow(2, float64(attempt))) * time.Minute
}

// ProcessPaymentWithRetry processes payment with automatic retries
func (ps *PaymentService) ProcessPaymentWithRetry(ctx context.Context, txn *Transaction) error {
    strategy := &RetryStrategy{
        MaxAttempts: 4,
        BackoffFunc: ExponentialBackoff,
    }
    
    for attempt := 0; attempt < strategy.MaxAttempts; attempt++ {
        err := ps.ProcessPayment(ctx, txn)
        if err == nil {
            return nil
        }
        
        // Check if retryable error
        if !isRetryable(err) {
            return err
        }
        
        // Sleep before retry
        if attempt < strategy.MaxAttempts-1 {
            backoff := strategy.BackoffFunc(attempt)
            select {
            case <-time.After(backoff):
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    }
    
    return fmt.Errorf("payment failed after %d attempts", strategy.MaxAttempts)
}

// Retryable errors
func isRetryable(err error) bool {
    // Retry on network errors, timeouts, provider temporary errors
    // Don't retry on validation errors, insufficient funds, etc
}
```

---

## 📝 Configuration & Environment

### Backend Configuration (Phase 3D)

```go
// Config for Phase 3D features
type Config struct {
    // Stripe
    Stripe struct {
        APIKey     string
        WebhookKey string
    }
    
    // PayPal
    PayPal struct {
        ClientID     string
        ClientSecret string
        WebhookID    string
    }
    
    // Redis
    Redis struct {
        Host     string
        Port     int
        Password string
        DB       int
    }
    
    // Analytics
    Analytics struct {
        EnableTracking bool
        BatchSize      int
        FlushInterval  time.Duration
    }
    
    // WebSocket
    WebSocket struct {
        MaxConnections int
        ReadBufferSize int
        WriteBufferSize int
    }
}
```

### Environment Variables (.env.production)

```bash
# Payment Processing
STRIPE_API_KEY=sk_live_...
STRIPE_WEBHOOK_KEY=whsec_...
PAYPAL_CLIENT_ID=...
PAYPAL_CLIENT_SECRET=...

# Caching & Real-time
REDIS_HOST=redis-prod.example.com
REDIS_PORT=6379
REDIS_PASSWORD=...

# Analytics
ANALYTICS_ENABLED=true
ANALYTICS_BATCH_SIZE=100

# WebSocket
WS_MAX_CONNECTIONS=5000
```

---

## ✅ Conclusion

This technical specification provides the detailed implementation roadmap for all Phase 3D features. Each section includes:
- Architecture diagrams
- Code examples
- Database schemas
- API specifications
- Integration checklists
- Performance considerations

Use this document alongside `PHASE3D_DEVELOPMENT_IMPLEMENTATION.md` for complete Phase 3D execution.

