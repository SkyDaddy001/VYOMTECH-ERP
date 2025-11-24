# Multi-Tenant AI-Based Call Center & Lead Management System Development Document

## Executive Summary

This document outlines the design and implementation of a production-ready, enterprise-grade multi-tenant AI-based call center and lead management system. The platform leverages Asterisk for telephony, Golang 1.24.x microservices, MySQL 8 for data persistence, and integrates multiple AI providers (Claude, Gemini, OpenAI, Ollama) for intelligent call handling.

The architecture supports inbound/outbound calls, AI agent interactions, human escalation, knowledge base retrieval, and marketing attribution. Deployment uses Kubernetes with Helm for orchestration, ensuring high availability and elastic scaling. Services are fully containerized with Docker, promoting modularity, independent deployment, and scalability. Each microservice is self-contained, with clear APIs, enabling plug-and-play updates and horizontal scaling. The enterprise edition is optimized for high concurrency, supporting up to 5000 concurrent callers per tenant, with advanced load balancing, auto-scaling, and distributed caching. This document serves as a blueprint for engineering teams to estimate, build, and operate the system, with actionable examples in Go, SQL, OpenAPI, Docker, Helm, and infrastructure configurations.

## Goals & Non-Goals

### Goals
- Deliver a scalable, secure multi-tenant platform handling up to 5000 concurrent callers per tenant, scaling to 100 tenants with advanced load balancing and auto-scaling.
- Integrate AI providers with intelligent routing based on cost, latency, and quality.
- Provide real-time call handling with WebSocket/gRPC for agent UIs and AI agents.
- Track marketing spends (e.g., Google Ads, FB Ads) and reconcile against credits with automated alerts.
- Enable multi-channel customer communication via Email, SMS, WhatsApp, and Slack with open-source tools (Postfix, Kannel, open-source WhatsApp libraries, Slack Bolt) and pluggable external provider integrations (e.g., Twilio, SendGrid, Slack API) for flexibility and failover.
- Support campaign management and tracking for Email, SMS, WhatsApp, and Slack, including scheduling, targeting, analytics, and conversion tracking integrated with lead management.
- Integrate Google Ads, Facebook Ads, and Slack APIs for analytics tracking, audience targeting, and automated ad campaign management (creation, optimization, budgeting) with real-time performance monitoring and A/B testing.
- Configure per-tenant DNS records (subdomains, MX, SPF/DKIM/DMARC, CNAME) for email routing, webhooks, and custom domains.
- Ensure observability with Prometheus/Grafana for metrics, Loki for logs, and Jaeger for traces.
- Achieve 99.9% uptime with sub-5s latency for call initiation, AI responses, and message delivery.
- Support enterprise compliance (PCI/PII) with audit logging and data encryption.
- Enable fully containerized, modular deployment for easy scaling, maintenance, and CI/CD.

### Non-Goals
- Direct integration with non-AI telephony providers beyond Asterisk.
- Real-time video calling or multimedia beyond voice/transcripts.
- Custom AI model training; focus on orchestration and fine-tuning via providers.
- Mobile app development; prioritize web-based admin/tenant UIs.
- Legacy system migrations; assume greenfield deployment.
- Proprietary-only integrations; prioritize open-source with external fallbacks.

## Scope & Major Capabilities

### High-Level Feature List
- **Telephony Integration**: Asterisk-based trunks, IVR, call recording, routing for inbound/outbound calls.
- **AI Orchestration**: Dynamic routing to AI providers with failover, caching, rate limiting, and cost tracking.
- **Call Handling**: Real-time WebSocket/gRPC for agent UI, conversation state machines, sentiment analysis, queue priority management for inbound/outbound calls.
- **Queue Priority Management**: Ability to set calling queue priorities for inbound/outbound calls.
- **Lead Management**: CRUD operations, scoring, assignment workflows, lifecycle tracking; Open API integrations with various portals for lead generation and tracking.
- **Marketing Spend Tracking**: Ingestion from APIs (Google Ads, FB), manual uploads; reconciliation against credits.
- **Knowledge Base**: Versioned KB with RAG pipeline, embedding stores (e.g., Pinecone or local vector DB).
- **Billing & Credits**: Per-tenant accounting, alerts for overspend, exportable reports.
- **Admin UI**: Dashboards for campaigns, credits, analytics; tenant-specific views.
- **Reporting**: CSV/PDF exports, real-time dashboards for calls, leads, spends; multi-project campaign reports with individual ROI tracking per project.
- **Observability**: Metrics for AI performance (latency, hallucination rate), SRE runbooks.
- **Containerization & Modularity**: Each service as a Docker container; Helm charts for K8s deployment; modular APIs for extensibility.

## Non-Functional Requirements

- **Scalability**: Support 10 tenants initially, scaling to 100 with horizontal pod scaling in K8s; Asterisk capacity for 5000 concurrent calls per tenant via advanced load balancing (e.g., NGINX or Envoy), auto-scaling policies (HPA with custom metrics), and distributed caching (Redis Cluster for session state, in-memory data grids for high-throughput operations).
- **Latency**: <5s for call initiation; <2s for AI response generation; <1s for lead CRUD.
- **Availability**: 99.9% uptime; multi-zone deployment for resilience.
- **Compliance**: PCI DSS for payment data; GDPR/CCPA for PII; SOX for financial reconciliation.
- **Security**: TLS 1.3 everywhere; RBAC with SSO (OAuth2); encrypted storage/transit for recordings/transcripts.
- **Cost**: Per-tenant attribution; AI token/call limits; infra optimization (e.g., spot instances for non-critical workloads).
- **Modularity**: Services decoupled via APIs; containerized for independent updates; event-driven where possible.

## Software Engineering Best Practices

### SOLID Principles
The system adheres to SOLID principles to ensure maintainable, extensible, and testable code:

- **Single Responsibility Principle (SRP)**: Each microservice (e.g., Call Handler, AI Orchestrator) has a single, well-defined responsibility. For example, the Lead Management Service only handles lead CRUD and workflows, not call processing.
- **Open-Closed Principle (OCP)**: Services are open for extension but closed for modification. AI providers are added via adapter pattern without changing core orchestrator logic.
- **Liskov Substitution Principle (LSP)**: All AI provider adapters implement the same interface, allowing seamless substitution (e.g., switching from OpenAI to Claude).
- **Interface Segregation Principle (ISP)**: Interfaces are client-specific and minimal. For example, the AIProvider interface only exposes Call() and CostPerToken() methods.
- **Dependency Inversion Principle (DIP)**: High-level modules (e.g., Call Handler) depend on abstractions (interfaces) rather than concrete implementations, enabling easy testing and swapping of components like databases or caches.

### ACID Properties
Database transactions follow ACID properties to ensure data integrity, especially critical for multi-tenant billing and lead management:

- **Atomicity**: All operations in a transaction (e.g., credit deduction and spend recording) succeed or fail together. Implemented via MySQL transactions in Go services.
- **Consistency**: Transactions move the database from one valid state to another. Schema per tenant ensures tenant isolation; foreign key constraints maintain referential integrity.
- **Isolation**: Concurrent transactions don't interfere. MySQL's InnoDB engine provides snapshot isolation to prevent dirty reads, non-repeatable reads, and phantom reads.
- **Durability**: Committed transactions survive failures. MySQL's redo logs and binary logs ensure durability; backups and replication provide additional resilience.

### Security Best Practices
Beyond basic TLS and RBAC, the system implements advanced security measures:

- **Defense in Depth**: Multi-layered security with network segmentation (VPCs), service mesh (Istio mTLS), and application-level controls (input validation, rate limiting).
- **Zero Trust Architecture**: Every request is authenticated and authorized, even within the cluster. JWT tokens include tenant context for multi-tenant isolation.
- **Secure Coding**: Input sanitization, parameterized queries to prevent SQL injection, and regular dependency scanning for vulnerabilities.
- **Secrets Management**: HashiCorp Vault for API keys and credentials; no secrets in code or configs.
- **Audit Logging**: All data access and modifications logged immutably for compliance and forensics.
- **Encryption**: AES-256 for data at rest; TLS 1.3 for transit; PII masking in transcripts.

### Performance Optimization
To achieve sub-5s latency and handle 5000 concurrent calls per tenant:

- **Caching Strategy**: Redis Cluster for session state, AI responses, and frequently accessed data. Cache-aside pattern with TTL-based expiration.
- **Database Indexing**: Composite indexes on tenant_id + status/date for efficient queries. Partitioning by tenant for large datasets.
- **Asynchronous Processing**: Event-driven architecture with Kafka for non-blocking operations like lead scoring and spend reconciliation.
- **Connection Pooling**: Database connection pools (e.g., via Go's sql.DB) to reuse connections and reduce overhead.
- **Horizontal Scaling**: Stateless services scale via K8s HPA; stateful services (MySQL) use read replicas for load distribution.
- **Resource Optimization**: Spot instances for non-critical workloads; GPU support for AI embeddings; efficient Docker images with multi-stage builds.

### Reliability Patterns
Ensuring 99.9% uptime and resilience:

- **Circuit Breaker**: Implemented in Go services to fail fast during AI provider outages, preventing cascade failures.
- **Retry with Backoff**: Exponential backoff for transient failures (e.g., network timeouts); idempotent operations for safe retries.
- **Graceful Degradation**: System continues operating with reduced functionality (e.g., fallback to cached responses during AI downtime).
- **Health Checks**: Liveness and readiness probes in K8s; dependency health monitoring (e.g., DB connectivity).
- **Chaos Engineering**: Regular injection of failures (pod kills, network delays) to test resilience.
- **Monitoring & Alerting**: Prometheus metrics with SLO-based alerts; Jaeger traces for debugging distributed issues.

## Multi-Tenancy Design Options & Recommended Approach

### Options
1. **Row-Level Multitenancy**: Single schema with tenant_id columns; pros: simple, shared resources; cons: noisy neighbor risks, complex queries.
2. **Schema Per Tenant**: Dedicated schema per tenant; pros: strong isolation; cons: higher overhead, migration complexity.
3. **Database Per Tenant**: Separate DB instances; pros: ultimate isolation; cons: high cost, management burden.

### Recommended Approach
Schema per tenant for strong isolation, with tenant-specific tables. Use dedicated schemas in MySQL for each tenant. Pros: strong isolation; cons: higher overhead, migration complexity. Decision needed: Confirm tenant count growth (e.g., >100 may require schema management automation). Implement tenant context in Go services via middleware for request scoping, ensuring modularity in tenant handling.

## System Architecture

### Component Diagram (Mermaid)
```mermaid
graph TD
    A[Asterisk Telephony] --> B[Call Handler Service (Go)]
    B --> C[AI Orchestrator (Go)]
    C --> D[Claude/Gemini/OpenAI/Ollama Adapters]
    B --> E[Conversation Engine (Go)]
    E --> F[Knowledge Base Service (Go)]
    B --> G[Lead Management Service (Go)]
    H[Marketing Spend Tracker (Go)] --> I[Reconciliation Engine]
    J[Email Service (Go)] --> K[Postfix/Open-Source SMTP]
    L[SMS Service (Go)] --> M[Kannel/Open-Source SMS Gateway]
    N[WhatsApp Service (Go)] --> O[Open-Source WhatsApp Libraries]
    P[Slack Service (Go)] --> Q[Slack Bolt/Open-Source Libraries]
    R[Campaign Manager Service (Go)] --> S[Multi-Channel Orchestrator]
    S --> J
    S --> L
    S --> N
    S --> P
    T[Ad Management Service (Go)] --> U[Google Ads API]
    T --> V[Facebook Ads API]
    T --> W[Analytics & Targeting Engine]
    X[Admin/Tenant UI (React)] --> Y[API Gateway (Go)]
    Y --> B
    Y --> G
    Y --> H
    Y --> J
    Y --> L
    Y --> N
    Y --> P
    Y --> R
    Y --> T
    Z[MySQL 8] --> All Services
    AA[Vector DB (e.g., Pinecone)] --> F
    BB[Prometheus/Grafana] --> Observability
    CC[Jaeger] --> Traces
    DD[Docker Containers] --> EE[K8s Pods]
    EE --> FF[Helm Charts]
    GG[Event Bus (Kafka)] --> Event-Driven Services
    HH[External Providers] --> J
    HH --> L
    HH --> N
    HH --> P
    II[DNS/Webhooks] --> HH
```

### Data Flows
- Inbound Call: Asterisk → Call Handler → AI Orchestrator → AI Provider → Conversation Engine → Lead Creation.
- Marketing Data: External APIs → Spend Tracker → Reconciliation → Billing Alerts.
- UI Interactions: Admin UI → API Gateway → Services → MySQL/Vector DB.
- Events: Services publish to Kafka for async processing (e.g., lead scoring).

### Deployment Topology
Kubernetes cluster with Helm charts: stateless services (e.g., Call Handler) in Deployments; stateful (MySQL, Vector DB) in StatefulSets. Ingress for external access; Service Mesh (Istio) for traffic management. Each service containerized with multi-stage Docker builds for efficiency. Modular Helm sub-charts for each service.

## Core Components and Responsibilities

- **Asterisk**: Handles SIP trunks, IVR, recording, routing. Containerized with persistent volumes for configs/recordings.
- **Call Handler Service**: Manages real-time calls via WebSocket/gRPC; interfaces with Asterisk. Modular with separate handlers for inbound/outbound.
- **AI Orchestrator**: Routes to AI providers; handles failover, caching, rate limiting. Pluggable adapters for extensibility.
- **Conversation Engine**: Maintains call state, transcripts, sentiment; integrates KB. State machine library (e.g., statemachine) for flows.
- **Lead Management**: CRUD for leads; scoring/assignment logic. Workflow engine (e.g., Temporal) for async processes.
- **Marketing Spend Tracker**: Ingests spends; reconciles with credits. Batch/streaming pipelines.
- **Knowledge Base Service**: Manages KB versions, RAG retrieval. Separate ingestion and query modules.
- **Billing Service**: Tracks per-tenant credits/spends; generates reports. Ledger-based accounting.
- **API Gateway**: Authenticates requests; routes to services. Kong or Traefik for API management.
- **Admin/Tenant UI**: Dashboards for management and analytics. React with micro-frontends for modularity.

## Telephony Layer (Asterisk)

- **Trunks**: SIP trunks for inbound/outbound; config for multiple providers.
- **Endpoints**: PJSIP for agents; registration and authentication.
- **IVR**: DTMF-based menus; integration with Go services for dynamic prompts.
- **Recording**: WAV/MP3 storage; encrypted at rest.
- **Call Routing**: Based on DID, time, or AI decision.

Example Asterisk pjsip.conf fragment:
```
[endpoint-template](!)
type=endpoint
context=default
disallow=all
allow=ulaw
auth=auth-template
aors=aor-template

[alice]
type=endpoint
context=default
disallow=all
allow=ulaw
auth=alice-auth
aors=alice-aor
```

Dockerfile for Asterisk:
```dockerfile
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y asterisk && rm -rf /var/lib/apt/lists/*
COPY pjsip.conf /etc/asterisk/
EXPOSE 5060 10000-20000/udp
VOLUME ["/var/lib/asterisk/sounds", "/var/spool/asterisk/monitor"]
CMD ["asterisk", "-f", "-U", "asterisk"]
```

Helm Chart Snippet for Asterisk:
```yaml
apiVersion: v2
name: asterisk
description: Asterisk Telephony Service
type: application
version: 0.1.0
appVersion: "18"

dependencies:
  - name: mysql
    version: "8.0"
    repository: "https://charts.bitnami.com/bitnami"

spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: asterisk
        image: asterisk:18
        ports:
        - containerPort: 5060
        volumeMounts:
        - name: config
          mountPath: /etc/asterisk
      volumes:
      - name: config
        configMap:
          name: asterisk-config
```

## Real-Time Call Handling Service in Go

Built with Gorilla WebSocket and gRPC. Interfaces: WebSocket for agent UI (real-time transcripts); gRPC for AI agents (streaming responses). Containerized with health checks.

Example Go struct:
```go
type CallSession struct {
    ID       string
    TenantID string
    Status   string // active, escalated
    Transcript []Message
}

func (s *CallHandler) HandleInboundCall(ctx context.Context, req *pb.InboundRequest) (*pb.CallResponse, error) {
    tenantID := middleware.GetTenantID(ctx)
    // Route to AI or agent
    return &pb.CallResponse{CallID: req.CallID}, nil
}
```

Dockerfile:
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o call-handler .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/call-handler .
EXPOSE 8080 9090
CMD ["./call-handler"]
```

Helm Chart Snippet:
```yaml
apiVersion: v2
name: call-handler
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: call-handler
        image: call-handler:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
```

## AI Orchestrator

Routes based on policy (cost first, then latency). Caches responses; throttles per provider. Redis for caching.

Example routing logic:
```go
func (o *Orchestrator) RouteRequest(req *AIRequest) (*AIResponse, error) {
    provider := o.selectProvider(req.TenantID, req.Priority)
    cached, err := o.cache.Get(req.Key)
    if err == nil {
        return cached, nil
    }
    resp, err := provider.Call(req)
    if err == nil {
        o.cache.Set(req.Key, resp)
    }
    return resp, err
}
```

Dockerfile similar to Call Handler, with Redis client.

Helm Chart: Includes Redis sub-chart.

## Conversation Engine & State Machine

Finite state machine for call flow: Greeting → Question → Response → Escalation. Tracks sentiment via AI analysis. Use go-statemachine library.

Dockerfile: Standard Go multi-stage.

## Lead Management Service

RESTful CRUD with workflows: New → Qualified → Assigned → Closed. Scoring via ML model (optional).

Example scoring: Based on call duration, sentiment, keywords.

Dockerfile: Standard Go multi-stage.

Helm Chart: With Temporal for workflows.

## Marketing Spend Tracker

Ingests via batch APIs or streaming; maps to campaigns/tenants. Kafka for streaming.

Dockerfile: Go with Kafka client.

## Knowledge Base Service

Versions KB; RAG with embeddings. Retrains weekly. Separate containers for ingestion and query.

Dockerfile: With Python for embeddings if needed.

Helm Chart: With Vector DB dependency.

## Billing & Credit Accounting

Per-tenant ledger; alerts on negative balance. Event-sourcing for immutability.

Dockerfile: Go.

## Admin & Tenant UI

React-based with micro-frontends for modularity; dashboards for calls, spends, leads, campaigns, and credits. Includes live display screens for real-time call status and agent status, updated via WebSocket connections to the Call Handler service. Tenant-specific views with RBAC filtering.

### Gamification and User Experience Design
To make the platform engaging and enjoyable to use, transforming work into an enjoyable experience that agents look forward to, the UI incorporates extensive gamification elements and friendly design principles:

- **Achievement System**: Agents earn badges, points, and titles for milestones like "100 Calls Handled," "Fastest Resolution Time," "High Customer Satisfaction Hero," or "Team Player." Leaderboards display top performers with animated celebrations, confetti effects, and virtual trophies. Seasonal challenges add excitement with limited-time rewards.
- **Progress Bars and Visual Feedback**: Interactive progress bars for daily/weekly goals, with color-coded indicators (green for on-track, yellow for close, red for behind). Smooth animations, particle effects, and micro-interactions provide immediate, satisfying feedback on every action. Progress rings fill up with each completed task, creating a sense of accomplishment.
- **Friendly Interface**: Warm color palettes (blues and greens), intuitive icons resembling friendly characters, and conversational language throughout. Tooltips and onboarding wizards feature animated mascots that guide new users with encouraging messages and helpful hints. Error messages are gentle and constructive, like "Oops! Let's try that again together."
- **Game-Like Elements**: Call handling feels like an epic mission - agents "level up" with experience points, unlock custom themes, avatars, and soundtracks. Team challenges encourage collaboration, with shared goals and group celebrations. Mini-games during breaks (like quick puzzles or relaxation exercises) prevent burnout and keep energy high.
- **Narrative and Storytelling**: Each agent has a "hero journey" with story arcs based on performance. Virtual pets or mascots react to achievements, providing companionship and motivation. Social features include shoutouts, kudos walls, and team stories that highlight collective successes.
- **Personalization**: Customizable dashboards with drag-and-drop widgets, favorite shortcuts, and personalized greetings. Unlockable themes, wallpapers, and UI skins based on achievements. Dark/light mode toggles, accessibility features, and ergonomic layouts ensure comfort and reduce fatigue.
- **Motivational Features**: Daily streaks for consistent performance with streak bonuses, virtual rewards (e.g., coffee icons for breaks, gift cards for milestones), and positive reinforcement messages like "Amazing work! You're a star!" Sound effects for successes, gentle nudges for breaks, and ambient background music that adapts to performance levels.
- **Community and Social Elements**: Agent profiles with photos, bios, and achievement showcases. Team chat with emojis and stickers, recognition posts, and virtual high-fives. Monthly "hall of fame" features top performers with spotlight stories.

### Live Display Screens
- **Call Status Dashboard**: Real-time grid showing active calls (inbound/outbound), duration, status (ringing, connected, on-hold), AI provider used, and escalation flags. Filters by tenant/campaign. Includes gamified elements like call "health bars" and quick-action buttons with satisfying animations.
- **Agent Status Dashboard**: Live view of agent availability (online, busy, offline), current call assignments, queue lengths, and performance metrics (calls handled, avg. handle time). Supports multi-tenant isolation. Features team spirit elements like group progress trackers and collaborative goal celebrations.

Example React Component for Call Status:
```jsx
import React, { useEffect, useState } from 'react';
import { w3cwebsocket as W3CWebSocket } from 'websocket';

const CallStatusDashboard = () => {
  const [calls, setCalls] = useState([]);

  useEffect(() => {
    const client = new W3CWebSocket('ws://call-handler:8080/ws/calls');
    client.onmessage = (message) => {
      const data = JSON.parse(message.data);
      setCalls(data.calls);
    };
    return () => client.close();
  }, []);

  return (
    <div>
      <h2>Live Call Status</h2>
      <table>
        <thead>
          <tr>
            <th>Call ID</th>
            <th>Status</th>
            <th>Duration</th>
            <th>Agent</th>
          </tr>
        </thead>
        <tbody>
          {calls.map(call => (
            <tr key={call.id}>
              <td>{call.id}</td>
              <td>{call.status}</td>
              <td>{call.duration}</td>
              <td>{call.agent}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default CallStatusDashboard;
```

Dockerfile:
```dockerfile
FROM node:18-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
```

Helm Chart: With Ingress and WebSocket support.
```yaml
apiVersion: v2
name: admin-ui
spec:
  template:
    spec:
      containers:
      - name: admin-ui
        image: admin-ui:latest
        ports:
        - containerPort: 80
        env:
        - name: WS_URL
          value: "ws://call-handler:8080"
```

## Reporting & Analytics

Grafana for real-time; export via API. Presto for ad-hoc queries.

## Data Model

### ER Diagram (Mermaid)
```mermaid
erDiagram
    TENANT ||--o{ USER : has
    TENANT ||--o{ CALL : owns
    CALL ||--o{ TRANSCRIPT : contains
    CALL ||--o{ AI_CALL : triggers
    TENANT ||--o{ CAMPAIGN : manages
    CAMPAIGN ||--o{ SPEND_EVENT : incurs
    TENANT ||--o{ CREDIT_EVENT : receives
    SPEND_EVENT ||--o{ RECONCILIATION : reconciles
```

### Key Tables with Example DDL

```sql
-- Schema per tenant, e.g., CREATE SCHEMA tenant_1;

CREATE TABLE tenant (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    email VARCHAR(255) UNIQUE,
    role ENUM('admin', 'agent'),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id)
);

CREATE TABLE call (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    direction ENUM('inbound', 'outbound'),
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    status ENUM('completed', 'escalated'),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id)
);

CREATE TABLE transcript (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    call_id BIGINT NOT NULL,
    speaker ENUM('agent', 'caller'),
    text TEXT,
    timestamp TIMESTAMP,
    FOREIGN KEY (call_id) REFERENCES call(id)
);

CREATE TABLE ai_call (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    call_id BIGINT NOT NULL,
    provider VARCHAR(50),
    tokens_used INT,
    cost DECIMAL(10,4),
    FOREIGN KEY (call_id) REFERENCES call(id)
);

CREATE TABLE campaign (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(255),
    platform VARCHAR(50),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id)
);

CREATE TABLE spend_event (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    campaign_id BIGINT NOT NULL,
    amount DECIMAL(10,2),
    date DATE,
    FOREIGN KEY (campaign_id) REFERENCES campaign(id)
);
```

## APIs (OpenAPI Snippets)

### Telephony Webhook Endpoints
```yaml
paths:
  /webhooks/call-start:
    post:
      summary: Notify call start
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                callId: {type: string}
                tenantId: {type: string}
```

### Real-Time WebSocket/gRPC Contract
gRPC proto:
```proto
service CallService {
  rpc StreamTranscript (StreamRequest) returns (stream TranscriptMessage);
}
```

### Lead CRUD and Workflows
```yaml
paths:
  /leads:
    post:
      summary: Create lead
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                name: {type: string}
                score: {type: number}
```

### Marketing Spend Ingestion API
Batch endpoint for spends.

### Admin Endpoints for Credits and Reconciliations
```yaml
paths:
  /admin/credits:
    post:
      summary: Add credits
```

## AI Integration Details

### Supported Providers & Adapter Pattern
Adapters for Claude, Gemini, OpenAI, Ollama. Interface:
```go
type AIProvider interface {
    Call(req *AIRequest) (*AIResponse, error)
    CostPerToken() float64
}
```

### Per-Provider Config
Stored in DB; includes API keys, regions.

### Routing Policy
Cost-first; fallback on failure.

### Retry & Throttling
Exponential backoff; per-provider limits.

Example adapter:
```go
type OpenAIAdapter struct {
    client *openai.Client
}

func (a *OpenAIAdapter) Call(req *AIRequest) (*AIResponse, error) {
    resp, err := a.client.CreateChatCompletion(ctx, req.ToOpenAI())
    return &AIResponse{Text: resp.Choices[0].Message.Content}, err
}
```

## Data Privacy (PII Handling, Retention)

PII masked in transcripts; retention: 7 years for compliance. Encryption: AES-256 at rest. Tokenization for sensitive data.

## Call Flow Examples

### Inbound Call to AI Agent + Escalation (Mermaid Sequence Diagram)
```mermaid
sequenceDiagram
    Caller->>Asterisk: Call
    Asterisk->>CallHandler: Webhook
    CallHandler->>AIOrchestrator: Route
    AIOrchestrator->>AIProvider: Query
    AIProvider-->>CallHandler: Response
    CallHandler->>AgentUI: Stream
    AgentUI->>CallHandler: Escalate
    CallHandler->>HumanAgent: Connect
```

### Outbound Campaign Dialing + AI Agent Conversation + Lead Creation
```mermaid
sequenceDiagram
    Campaign->>CallHandler: Dial
    CallHandler->>Asterisk: Initiate
    Asterisk->>Caller: Ring
    Caller->>CallHandler: Answer
    CallHandler->>AIOrchestrator: Start Conversation
    AIOrchestrator->>LeadService: Create Lead
```

## Knowledge Base & RAG

Ingestion: Batch from docs; versioning via Git-like tags. Embedding: OpenAI or local. Vector DB: Pinecone for cloud, Chroma for on-prem. Containerized with GPU support for embeddings.

Health: Freshness checks; retrain on KB updates.

## AI Performance Monitoring & KPIs

KPIs: Latency (p95 <2s), Throughput (calls/min), Hallucination Rate (manual review %), Accuracy (KB query match %), Escalation Rate (<10%), NPS proxy (sentiment score).

Metrics: `ai_response_latency`, `ai_hallucination_count`.

Grafana: Panels for latency trends, error rates.

## Marketing Spend Tracking & Reconciliation

Data Sources: Google Ads API, FB Graph, manual CSV.

Event Model: SpendEvent vs CreditEvent.

Reconciliation: SQL join on campaign/date; alert on discrepancy >5%.

Example SQL:
```sql
SELECT c.name, se.amount AS spend, ce.amount AS credit, (ce.amount - se.amount) AS discrepancy
FROM campaign c
JOIN spend_event se ON c.id = se.campaign_id
LEFT JOIN credit_event ce ON c.id = ce.campaign_id AND se.date = ce.date;
```

Alerts: Prometheus rule for negative balance.

Report Fields: tenant, campaign_id, platform, spend_date, spend_amount, credited_amount, discrepancy.

## Security & Compliance

Auth: OAuth2 SSO; RBAC per tenant. Vault for secrets.

Network: VPC segmentation; TLS everywhere. mTLS in service mesh.

PCI/PII: Recordings encrypted; retention 1 year. Data masking in logs.

Audit: All actions logged to ELK. Immutable audit trails.

## Observability & SRE Playbook

Metrics: Prometheus; Logs: Loki; Traces: Jaeger.

Runbooks: Call degradation (scale pods); AI outage (switch provider); Billing anomaly (manual review).

SLOs: 99.9% call success; Alert thresholds: latency >5s.

## Deployment & Infra

### Containerization Strategy
Each service has a Dockerfile with multi-stage builds. Use Docker Compose for local dev.

Example docker-compose.yml for local development:
```yaml
version: '3.8'
services:
  call-handler:
    build: ./call-handler
    ports:
      - "8080:8080"
    environment:
      - DB_URL=mysql://root:password@mysql:3306/db
    depends_on:
      - mysql
  mysql:
    image: mysql:8
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: db
    ports:
      - "3306:3306"
  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
```

Helm Charts: Modular sub-charts for each service.

Full Helm Chart Example for Call Handler:
```yaml
apiVersion: v2
name: call-handler
description: Call Handler Microservice
version: 0.1.0

dependencies:
  - name: redis
    version: "17.0"
    repository: "https://charts.bitnami.com/bitnami"

spec:
  replicas: {{ .Values.replicas }}
  template:
    spec:
      containers:
      - name: call-handler
        image: {{ .Values.image }}
        ports:
        - containerPort: 8080
        env:
        - name: DB_URL
          valueFrom:
            secretKeyRef:
              name: {{ .Values.secretName }}
              key: dbUrl
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: {{ .Values.secretName }}
              key: redisUrl
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: config
        configMap:
          name: call-handler-config
```

Stateful: MySQL StatefulSet with PVCs.

Scaling: Asterisk with SBC; capacity: 200 calls via load testing.

CI/CD: GitHub Actions for build/test/deploy. ArgoCD for GitOps.

## Testing Strategy

Unit: Go tests for adapters. Mock AI providers.

Integration: gRPC contract tests. Testcontainers for DB.

Load: 200 calls/min with AI mocks. k6 for scripting.

Chaos: Provider latency injection; DB partial outage. LitmusChaos.

## Cost Model & Estimation Guidance

Costs: AI (tokens), Telephony (minutes), Infra (K8s nodes).

Attribution: Per-tenant DB queries. Tags on resources.

Example: Tenant A: 1000 calls/month @ $0.10/call = $100.

Reconciliation: CSV export with spend/credit diffs.

## Roadmap & Milestones

MVP: Inbound calls, AI routing, lead CRUD, basic spend tracking.

Phased: Q1: Core telephony; Q2: AI full; Q3: Marketing recon.

Acceptance: 90% call success rate; <5% hallucination.

## Appendix

### Helpful Code Snippets
Go AI Adapter Interface: As above.

Sample SQL: As in DDL.

Asterisk Fragment: As above.

### Example OpenAPI YAML for Lead Endpoints
As in APIs section.

### Example Prometheus Alert Rule
```
alert: HighAILatency
expr: ai_response_latency > 5
for: 5m
```

### Grafana Panel JSON
Description: Time series for latency.

### Glossary
- Tenant: Isolated customer instance.
- RAG: Retrieval-Augmented Generation.

## 2-Week Sprint MVP Checklist
- [ ] Set up K8s cluster with Helm.
- [ ] Implement Call Handler service with WebSocket.
- [ ] Integrate Asterisk for inbound calls.
- [ ] Add AI Orchestrator with one provider (OpenAI).
- [ ] Create Lead Management CRUD.
- [ ] Basic spend ingestion API.
- [ ] Unit tests for core services.
- [ ] Deploy MVP; test 10 calls.

## Detailed Design Documentation

### Detailed Component Diagrams

#### Call Handler Service Detailed Diagram
```mermaid
graph TD
    A[WebSocket Handler] --> B[Call Session Manager]
    B --> C[Transcript Processor]
    C --> D[Sentiment Analyzer]
    D --> E[Escalation Logic]
    E --> F[Agent Connector]
    G[gRPC Handler] --> H[AI Streaming Interface]
    H --> I[Conversation State Machine]
    I --> J[Knowledge Base Querier]
    J --> K[Lead Creator]
```

#### AI Orchestrator Detailed Diagram
```mermaid
graph TD
    A[Request Router] --> B[Provider Selector]
    B --> C[Cost Calculator]
    C --> D[Latency Checker]
    D --> E[Quality Assessor]
    E --> F[Cache Layer]
    F --> G[Rate Limiter]
    G --> H[Adapter Executor]
    H --> I[Response Processor]
    I --> J[Metrics Collector]
```

### Sequence Diagrams for Key Flows

#### Inbound Call Flow Detailed Sequence
```mermaid
sequenceDiagram
    participant Caller
    participant Asterisk
    participant CallHandler
    participant AIOrchestrator
    participant AIProvider
    participant ConversationEngine
    participant LeadService
    participant DB

    Caller->>Asterisk: Dial Number
    Asterisk->>CallHandler: Webhook (call-start)
    CallHandler->>CallHandler: Create CallSession
    CallHandler->>AIOrchestrator: Route to AI
    AIOrchestrator->>AIProvider: Query with Context
    AIProvider-->>AIOrchestrator: AI Response
    AIOrchestrator-->>CallHandler: Response
    CallHandler->>ConversationEngine: Update State
    ConversationEngine->>LeadService: Create Lead
    LeadService->>DB: Insert Lead
    CallHandler->>Caller: Play Response
```

#### Lead Scoring and Assignment Flow
```mermaid
sequenceDiagram
    participant CallHandler
    participant LeadService
    participant ScoringEngine
    participant AssignmentEngine
    participant AgentUI
    participant DB

    CallHandler->>LeadService: New Lead Event
    LeadService->>ScoringEngine: Score Lead
    ScoringEngine->>ScoringEngine: Analyze Call Data
    ScoringEngine-->>LeadService: Score Value
    LeadService->>AssignmentEngine: Assign Agent
    AssignmentEngine->>DB: Query Available Agents
    AssignmentEngine-->>LeadService: Assigned Agent
    LeadService->>DB: Update Lead
    LeadService->>AgentUI: Notify Assignment
```

### API Specifications

#### Call Handler API
```yaml
openapi: 3.0.0
info:
  title: Call Handler API
  version: 1.0.0
paths:
  /calls:
    post:
      summary: Initiate outbound call
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                tenantId: {type: string}
                phoneNumber: {type: string}
                campaignId: {type: string}
    get:
      summary: Get active calls
      parameters:
        - name: tenantId
          in: query
          schema: {type: string}
  /calls/{callId}/escalate:
    post:
      summary: Escalate call to human agent
      parameters:
        - name: callId
          in: path
          schema: {type: string}
```

#### Lead Management API
```yaml
openapi: 3.0.0
info:
  title: Lead Management API
  version: 1.0.0
paths:
  /leads:
    get:
      summary: List leads
      parameters:
        - name: tenantId
          in: query
          schema: {type: string}
        - name: status
          in: query
          schema: {type: string, enum: [new, qualified, assigned, closed]}
    post:
      summary: Create lead
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                tenantId: {type: string}
                name: {type: string}
                phone: {type: string}
                score: {type: number}
                source: {type: string}
  /leads/{leadId}:
    put:
      summary: Update lead
      parameters:
        - name: leadId
          in: path
          schema: {type: string}
    delete:
      summary: Delete lead
      parameters:
        - name: leadId
          in: path
          schema: {type: string}
  /leads/{leadId}/assign:
    post:
      summary: Assign lead to agent
      parameters:
        - name: leadId
          in: path
          schema: {type: string}
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                agentId: {type: string}
```

### Database Schema Details

#### Tenant Schema Structure
```sql
-- Per-tenant schema creation
CREATE SCHEMA tenant_{tenant_id};

-- Switch to tenant schema
USE tenant_{tenant_id};

-- Core tables
CREATE TABLE tenant (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE user (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL DEFAULT 1, -- For shared schema fallback
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('admin', 'agent', 'supervisor') NOT NULL,
    status ENUM('active', 'inactive') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_email (tenant_id, email)
);

CREATE TABLE call (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL DEFAULT 1,
    direction ENUM('inbound', 'outbound') NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NULL,
    duration_seconds INT NULL,
    status ENUM('ringing', 'connected', 'completed', 'escalated', 'failed') NOT NULL,
    agent_id BIGINT NULL,
    campaign_id BIGINT NULL,
    ai_provider VARCHAR(50) NULL,
    cost DECIMAL(10,4) DEFAULT 0,
    recording_url VARCHAR(500) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_tenant_start_time (tenant_id, start_time),
    FOREIGN KEY (agent_id) REFERENCES user(id),
    FOREIGN KEY (campaign_id) REFERENCES campaign(id)
);

CREATE TABLE transcript (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    call_id BIGINT NOT NULL,
    sequence_number INT NOT NULL,
    speaker ENUM('agent', 'caller', 'ai') NOT NULL,
    text TEXT NOT NULL,
    confidence DECIMAL(3,2) NULL,
    timestamp TIMESTAMP NOT NULL,
    sentiment ENUM('positive', 'neutral', 'negative') NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_call_sequence (call_id, sequence_number),
    FOREIGN KEY (call_id) REFERENCES call(id)
);

CREATE TABLE ai_call (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    call_id BIGINT NOT NULL,
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,
    tokens_used INT NOT NULL,
    cost DECIMAL(10,4) NOT NULL,
    response_time_ms INT NOT NULL,
    success BOOLEAN NOT NULL,
    error_message TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_call_provider (call_id, provider),
    FOREIGN KEY (call_id) REFERENCES call(id)
);

CREATE TABLE lead (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL DEFAULT 1,
    name VARCHAR(255) NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255) NULL,
    score DECIMAL(5,2) DEFAULT 0,
    status ENUM('new', 'qualified', 'assigned', 'contacted', 'converted', 'closed') DEFAULT 'new',
    source ENUM('call', 'web', 'referral', 'campaign') DEFAULT 'call',
    campaign_id BIGINT NULL,
    assigned_agent_id BIGINT NULL,
    assigned_at TIMESTAMP NULL,
    last_contacted_at TIMESTAMP NULL,
    notes TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_tenant_score (tenant_id, score DESC),
    FOREIGN KEY (campaign_id) REFERENCES campaign(id),
    FOREIGN KEY (assigned_agent_id) REFERENCES user(id)
);

CREATE TABLE campaign (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL DEFAULT 1,
    name VARCHAR(255) NOT NULL,
    platform VARCHAR(50) NOT NULL, -- google_ads, facebook, manual
    budget DECIMAL(10,2) NULL,
    status ENUM('active', 'paused', 'completed') DEFAULT 'active',
    start_date DATE NULL,
    end_date DATE NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_status (tenant_id, status)
);

CREATE TABLE spend_event (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL DEFAULT 1,
    campaign_id BIGINT NOT NULL,
    platform VARCHAR(50) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    date DATE NOT NULL,
    external_id VARCHAR(255) NULL, -- Platform-specific ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_campaign_date (tenant_id, campaign_id, date),
    FOREIGN KEY (campaign_id) REFERENCES campaign(id)
);

CREATE TABLE credit_event (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL DEFAULT 1,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    date DATE NOT NULL,
    type ENUM('purchase', 'adjustment', 'refund') NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_date (tenant_id, date)
);

CREATE TABLE reconciliation (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL DEFAULT 1,
    spend_event_id BIGINT NOT NULL,
    credit_event_id BIGINT NULL,
    discrepancy DECIMAL(10,2) DEFAULT 0,
    status ENUM('matched', 'unmatched', 'adjusted') DEFAULT 'unmatched',
    adjusted_at TIMESTAMP NULL,
    notes TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_status (tenant_id, status),
    FOREIGN KEY (spend_event_id) REFERENCES spend_event(id),
    FOREIGN KEY (credit_event_id) REFERENCES credit_event(id)
);
```

### Security Design

#### Authentication & Authorization
- **JWT Tokens**: Stateless authentication with tenant context
- **RBAC**: Role-based access control with tenant isolation
- **API Keys**: For external integrations (Google Ads, FB)
- **OAuth2**: SSO integration with enterprise providers

#### Data Encryption
- **At Rest**: AES-256 encryption for sensitive data
- **In Transit**: TLS 1.3 for all communications
- **Secrets**: HashiCorp Vault for API keys and credentials

#### Network Security
- **VPC Isolation**: Separate VPCs per environment
- **Service Mesh**: Istio with mTLS between services
- **WAF**: Cloudflare or similar for external traffic
- **Rate Limiting**: Per-tenant and per-endpoint limits

#### Compliance Controls
- **PII Masking**: Automatic detection and masking in transcripts
- **Audit Logging**: All data access and modifications logged
- **Retention Policies**: Automated data lifecycle management
- **Access Reviews**: Regular review of user permissions

## Detailed Testing Document

### Unit Testing

#### Go Service Unit Tests
```go
// Example test for AI Orchestrator
func TestOrchestrator_RouteRequest(t *testing.T) {
    tests := []struct {
        name     string
        req      *AIRequest
        expected *AIResponse
        wantErr  bool
    }{
        {
            name: "successful routing",
            req: &AIRequest{
                TenantID: "tenant1",
                Priority: "cost",
                Text:     "Hello",
            },
            expected: &AIResponse{Text: "Hi there!"},
            wantErr:  false,
        },
        {
            name: "provider failure fallback",
            req: &AIRequest{
                TenantID: "tenant1",
                Priority: "latency",
                Text:     "Test",
            },
            expected: &AIResponse{Text: "Fallback response"},
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            o := &Orchestrator{
                providers: []AIProvider{&mockProvider{}},
                cache:     &mockCache{},
            }
            got, err := o.RouteRequest(tt.req)
            if (err != nil) != tt.wantErr {
                t.Errorf("Orchestrator.RouteRequest() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.expected) {
                t.Errorf("Orchestrator.RouteRequest() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

#### React Component Unit Tests
```javascript
// Example test for CallStatusDashboard
import { render, screen, waitFor } from '@testing-library/react';
import CallStatusDashboard from './CallStatusDashboard';

const mockWebSocket = {
  onmessage: jest.fn(),
  close: jest.fn(),
};

jest.mock('websocket', () => ({
  w3cwebsocket: jest.fn(() => mockWebSocket),
}));

test('renders call status dashboard', async () => {
  render(<CallStatusDashboard />);

  expect(screen.getByText('Live Call Status')).toBeInTheDocument();

  // Simulate WebSocket message
  mockWebSocket.onmessage({
    data: JSON.stringify({
      calls: [
        { id: '1', status: 'connected', duration: '00:05:30', agent: 'Alice' },
      ],
    }),
  });

  await waitFor(() => {
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('connected')).toBeInTheDocument();
  });
});
```

### Integration Testing

#### API Integration Test Cases
```go
// Test suite for Call Handler API
func TestCallHandlerAPI(t *testing.T) {
    suite.Run(t, new(CallHandlerAPITestSuite))
}

type CallHandlerAPITestSuite struct {
    suite.Suite
    server *httptest.Server
    client *http.Client
}

func (s *CallHandlerAPITestSuite) SetupTest() {
    // Setup test server and database
    s.server = httptest.NewServer(setupRouter())
    s.client = s.server.Client()
}

func (s *CallHandlerAPITestSuite) TearDownTest() {
    s.server.Close()
}

func (s *CallHandlerAPITestSuite) TestCreateOutboundCall() {
    reqBody := `{
        "tenantId": "tenant1",
        "phoneNumber": "+1234567890",
        "campaignId": "camp1"
    }`

    resp, err := s.client.Post(
        s.server.URL+"/calls",
        "application/json",
        strings.NewReader(reqBody),
    )
    s.NoError(err)
    s.Equal(http.StatusCreated, resp.StatusCode)

    var response map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&response)
    s.NoError(err)
    s.Contains(response, "callId")
}

func (s *CallHandlerAPITestSuite) TestGetActiveCalls() {
    resp, err := s.client.Get(s.server.URL + "/calls?tenantId=tenant1")
    s.NoError(err)
    s.Equal(http.StatusOK, resp.StatusCode)

    var calls []map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&calls)
    s.NoError(err)
    // Verify call structure
}
```

#### Database Integration Tests
```go
func TestLeadRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := &LeadRepository{db: db}

    lead := &Lead{
        TenantID: 1,
        Name:     "John Doe",
        Phone:    "+1234567890",
        Score:    85.5,
        Status:   "new",
    }

    // Test Create
    err := repo.Create(context.Background(), lead)
    assert.NoError(t, err)
    assert.NotZero(t, lead.ID)

    // Test GetByID
    retrieved, err := repo.GetByID(context.Background(), lead.ID)
    assert.NoError(t, err)
    assert.Equal(t, lead.Name, retrieved.Name)

    // Test Update
    lead.Score = 90.0
    err = repo.Update(context.Background(), lead)
    assert.NoError(t, err)

    updated, err := repo.GetByID(context.Background(), lead.ID)
    assert.NoError(t, err)
    assert.Equal(t, 90.0, updated.Score)

    // Test List
    leads, err := repo.List(context.Background(), &LeadFilter{
        TenantID: 1,
        Status:   "new",
    })
    assert.NoError(t, err)
    assert.Len(t, leads, 1)
}
```

### Load Testing

#### K6 Load Test Script
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 }, // Ramp up to 100 users
    { duration: '5m', target: 100 }, // Stay at 100 users
    { duration: '2m', target: 200 }, // Ramp up to 200 users
    { duration: '5m', target: 200 }, // Stay at 200 users
    { duration: '2m', target: 0 },   // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.1'],    // Error rate should be below 10%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Simulate creating a lead
  let payload = JSON.stringify({
    tenantId: 'tenant1',
    name: 'Test User',
    phone: '+1234567890',
    score: Math.random() * 100,
  });

  let params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + __ENV.API_TOKEN,
    },
  };

  let response = http.post(`${BASE_URL}/leads`, payload, params);

  check(response, {
    'status is 201': (r) => r.status === 201,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });

  sleep(1);
}
```

#### AI Provider Load Test
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 50,
  duration: '10m',
  thresholds: {
    http_req_duration: ['p(95)<2000'], // AI responses should be <2s
  },
};

export default function () {
  let payload = JSON.stringify({
    tenantId: 'tenant1',
    message: 'Hello, I need help with my account',
    context: {
      callId: 'call_' + Math.random().toString(36).substr(2, 9),
      history: [],
    },
  });

  let response = http.post('http://localhost:8080/ai/chat', payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(response, {
    'AI response status is 200': (r) => r.status === 200,
    'response contains text': (r) => JSON.parse(r.body).text !== undefined,
  });

  sleep(0.5);
}
```

### Chaos Engineering Tests

#### Pod Failure Simulation
```yaml
# chaos-pod-failure.yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: pod-failure-chaos
  namespace: default
spec:
  action: pod-failure
  mode: one
  selector:
    namespaces:
      - default
    labelSelectors:
      app: call-handler
  duration: 5m
  scheduler:
    cron: '@every 10m'
```

#### Network Latency Injection
```yaml
# chaos-network-delay.yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: network-delay-chaos
  namespace: default
spec:
  action: delay
  mode: all
  selector:
    namespaces:
      - default
    labelSelectors:
      app: ai-orchestrator
  delay:
    latency: 1000ms
    correlation: 100
    jitter: 200ms
  duration: 3m
```

#### AI Provider Failure Simulation
```go
// chaos-ai-provider.go - Simulate AI provider failures
func (m *MockAIProvider) Call(req *AIRequest) (*AIResponse, error) {
    // Inject random failures
    if rand.Float32() < 0.1 { // 10% failure rate
        return nil, errors.New("simulated provider failure")
    }

    // Inject random latency
    delay := time.Duration(rand.Intn(2000)) * time.Millisecond
    time.Sleep(delay)

    return &AIResponse{
        Text: "This is a mock AI response",
        Tokens: len(strings.Fields(req.Text)),
    }, nil
}
```

### End-to-End Testing

#### E2E Test Scenarios
1. **Inbound Call to Lead Creation**
   - Simulate incoming call via Asterisk
   - Verify call handler receives webhook
   - Check AI orchestrator routes to provider
   - Validate conversation flow
   - Confirm lead creation in database
   - Verify real-time updates in UI

2. **Outbound Campaign Call**
   - Create campaign in system
   - Initiate outbound call
   - Verify call connection
   - Test AI conversation
   - Check lead qualification
   - Validate spend tracking

3. **Escalation to Human Agent**
   - Start AI conversation
   - Trigger escalation condition
   - Verify agent notification
   - Test handoff process
   - Confirm transcript continuity

#### E2E Test Implementation
```go
func TestInboundCallE2E(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test")
    }

    // Setup test environment
    env := setupTestEnvironment(t)
    defer env.cleanup()

    // Simulate incoming call
    callID := simulateIncomingCall(t, env.asteriskClient, "+1234567890")

    // Wait for call processing
    time.Sleep(2 * time.Second)

    // Verify call record
    call := getCallByID(t, env.db, callID)
    assert.Equal(t, "completed", call.Status)
    assert.NotZero(t, call.DurationSeconds)

    // Verify lead creation
    leads := getLeadsByPhone(t, env.db, "+1234567890")
    assert.Len(t, leads, 1)
    assert.Equal(t, "qualified", leads[0].Status)

    // Verify AI call logging
    aiCalls := getAICallsByCallID(t, env.db, callID)
    assert.True(t, len(aiCalls) > 0)
    assert.True(t, aiCalls[0].Cost > 0)
}
```

### Test Data Management

#### Test Data Setup
```sql
-- test_data.sql
INSERT INTO tenant (id, name) VALUES (1, 'Test Tenant');

INSERT INTO user (id, tenant_id, email, role) VALUES
(1, 1, 'admin@test.com', 'admin'),
(2, 1, 'agent@test.com', 'agent');

INSERT INTO campaign (id, tenant_id, name, platform) VALUES
(1, 1, 'Test Campaign', 'manual');

-- Sample leads for testing
INSERT INTO lead (tenant_id, name, phone, score, status) VALUES
(1, 'John Doe', '+1234567890', 85.5, 'qualified'),
(1, 'Jane Smith', '+1234567891', 72.3, 'new');
```

#### Test Fixtures
```go
// test_fixtures.go
func LoadTestFixtures(db *sql.DB) error {
    fixtures := []string{
        "test_data/tenants.sql",
        "test_data/users.sql",
        "test_data/campaigns.sql",
        "test_data/leads.sql",
    }

    for _, fixture := range fixtures {
        sql, err := ioutil.ReadFile(fixture)
        if err != nil {
            return err
        }

        if _, err := db.Exec(string(sql)); err != nil {
            return err
        }
    }

    return nil
}
```

### Performance Testing KPIs

#### Response Time Targets
- API endpoints: p95 < 200ms
- AI responses: p95 < 2000ms
- Database queries: p95 < 100ms
- WebSocket messages: < 50ms

#### Throughput Targets
- API calls: 1000 req/sec per service
- AI requests: 100 req/sec per provider
- Database connections: 100 concurrent
- WebSocket connections: 5000 concurrent

#### Resource Utilization
- CPU: < 70% average
- Memory: < 80% average
- Disk I/O: < 1000 IOPS
- Network: < 100 Mbps

### Test Automation Pipeline

#### CI/CD Test Stages
```yaml
# .github/workflows/test.yml
name: Test Suite
on: [push, pull_request]

jobs:
  unit-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.24
      - name: Run unit tests
        run: go test ./... -v -coverprofile=coverage.out
      - name: Upload coverage
        uses: codecov/codecov-action@v2

  integration-test:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8
        env:
          MYSQL_ROOT_PASSWORD: password
    steps:
      - uses: actions/checkout@v2
      - name: Run integration tests
        run: go test ./... -tags=integration -v

  e2e-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set up Kubernetes
        uses: engineerd/setup-kind@v0.5.0
      - name: Deploy to test cluster
        run: |
          helm install test-release ./helm
      - name: Run E2E tests
        run: go test ./e2e -v

  load-test:
    runs-on: ubuntu-latest
    needs: [e2e-test]
    steps:
      - uses: actions/checkout@v2
      - name: Run K6 load tests
        uses: grafana/k6-action@v0.2.0
        with:
          filename: load-tests/main.js
```

### Test Reporting and Monitoring

#### Test Results Dashboard
- Unit test coverage trends
- Integration test pass/fail rates
- Performance benchmarks over time
- E2E test success rates

#### Alerting
- Test failures in CI/CD
- Performance regression detection
- Coverage drops below threshold
- Flaky test identification

## Implementation Guide

### Getting Started

#### Prerequisites
- Go 1.24.x
- Docker & Docker Compose
- Kubernetes cluster (for production)
- Helm 3.x
- MySQL 8.x
- Redis 7.x

#### Local Development Setup
```bash
# Clone the repository
git clone <repository-url>
cd multi-tenant-ai-call-center

# Start local services
docker-compose up -d mysql redis

# Install dependencies
go mod download

# Run database migrations
go run cmd/migrate/main.go

# Start services
go run cmd/call-handler/main.go &
go run cmd/ai-orchestrator/main.go &
go run cmd/lead-service/main.go &

# Start UI
cd ui && npm install && npm start
```

#### Repository Structure
```
├── cmd/                    # Application entrypoints
│   ├── call-handler/
│   ├── ai-orchestrator/
│   └── lead-service/
├── internal/               # Private application code
│   ├── call/
│   ├── ai/
│   ├── lead/
│   └── common/
├── pkg/                    # Public packages
│   ├── api/
│   ├── database/
│   └── middleware/
├── ui/                     # React frontend
├── helm/                   # Kubernetes Helm charts
├── docker/                 # Docker files
├── docs/                   # Documentation
├── scripts/                # Build and deployment scripts
└── test/                   # Test files and fixtures
```

### Infrastructure as Code

#### Terraform Configuration
```hcl
# main.tf
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

module "vpc" {
  source = "./modules/vpc"
  name   = "ai-call-center"
  cidr   = "10.0.0.0/16"
}

module "eks" {
  source          = "./modules/eks"
  cluster_name    = "ai-call-center-prod"
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnets
  node_groups = {
    general = {
      desired_capacity = 3
      max_capacity     = 10
      min_capacity     = 1
      instance_type    = "t3.medium"
    }
  }
}

module "rds" {
  source              = "./modules/rds"
  identifier          = "ai-call-center-prod"
  engine              = "mysql"
  engine_version      = "8.0"
  instance_class      = "db.t3.medium"
  allocated_storage   = 100
  vpc_id              = module.vpc.vpc_id
  subnet_ids          = module.vpc.database_subnets
  security_group_ids  = [aws_security_group.rds.id]
}
```

#### Docker Compose for Local Development
```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: aicallcenter
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  call-handler:
    build:
      context: .
      dockerfile: docker/Dockerfile.call-handler
    ports:
      - "8080:8080"
    environment:
      - DB_DSN=mysql://root:rootpassword@mysql:3306/aicallcenter
      - REDIS_URL=redis://redis:6379
    depends_on:
      - mysql
      - redis

  ai-orchestrator:
    build:
      context: .
      dockerfile: docker/Dockerfile.ai-orchestrator
    ports:
      - "8081:8081"
    environment:
      - REDIS_URL=redis://redis:6379
    depends_on:
      - redis

  lead-service:
    build:
      context: .
      dockerfile: docker/Dockerfile.lead-service
    ports:
      - "8082:8082"
    environment:
      - DB_DSN=mysql://root:rootpassword@mysql:3306/aicallcenter
    depends_on:
      - mysql

volumes:
  mysql_data:
