# Multi-Channel Communication System - Complete Guide

## Overview

A comprehensive multi-channel communication system integrating **Telegram**, **WhatsApp**, **SMS**, and **Email** with click-to-call capabilities. This system enables organizations to manage customer conversations across multiple channels from a single unified platform.

## Table of Contents

1. [Architecture](#architecture)
2. [Database Schema](#database-schema)
3. [API Endpoints](#api-endpoints)
4. [Provider Integration](#provider-integration)
5. [Features](#features)
6. [Implementation Guide](#implementation-guide)
7. [Configuration](#configuration)
8. [Webhook Handling](#webhook-handling)
9. [Examples](#examples)
10. [Best Practices](#best-practices)

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend / Mobile App                     │
├─────────────────────────────────────────────────────────────┤
│                   API Gateway / Router                       │
├─────────────────────────────────────────────────────────────┤
│          MultiChannelCommunicationHandler (HTTP)             │
├─────────────────────────────────────────────────────────────┤
│        MultiChannelCommunicationService (Business Logic)     │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┬────────────┬──────────────┬─────────┐         │
│  │ Telegram │ WhatsApp   │ SMS (Twilio) │ Email   │         │
│  │ Service  │ Business   │ & AWS SNS    │ Service │         │
│  │          │ API        │ & Vonage     │         │         │
│  └──────────┴────────────┴──────────────┴─────────┘         │
├─────────────────────────────────────────────────────────────┤
│              Database (MySQL/PostgreSQL)                     │
│  ┌─────────────────────────────────────────────────┐        │
│  │ • communication_channel                         │        │
│  │ • communication_session                         │        │
│  │ • communication_message                         │        │
│  │ • message_template                              │        │
│  │ • contact_communication_preference              │        │
│  │ • bulk_message_campaign                         │        │
│  │ • message_automation_rule                       │        │
│  │ • communication_analytics                       │        │
│  │ • scheduled_message                             │        │
│  └─────────────────────────────────────────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

## Database Schema

### 1. communication_channel
Stores configuration for each communication channel (provider credentials, API endpoints, webhook URLs).

**Key Fields:**
- `id`: Unique channel identifier
- `tenant_id`: Multi-tenant isolation
- `channel_type`: TELEGRAM, WHATSAPP, SMS, EMAIL
- `api_provider`: twilio, sendgrid, telegram, whatsapp_business, aws_sns, vonage, mailgun
- `api_key`, `api_secret`, `auth_token`: Authentication credentials
- `webhook_url`: URL for provider to send callbacks
- `sender_id`: From address (phone number, email, bot ID, etc.)
- `is_active`: Channel status
- `priority`: For routing multiple channels

### 2. communication_session
Represents a conversation with a contact across one channel.

**Key Fields:**
- `id`: Session identifier
- `channel_type`: Channel used for this session
- `sender_id`: Internal sender address
- `recipient_id`: Contact's address (phone/email/ID)
- `agent_id`: Assigned agent
- `lead_id`: Associated lead/opportunity
- `status`: INITIATED, SENT, DELIVERED, READ, FAILED, COMPLETED
- `direction`: INBOUND (from customer), OUTBOUND (from agent), INTERNAL
- `message_count`: Total messages in session
- `first_message_at`, `last_message_at`: Session timeline
- `is_archived`: For conversation archival

### 3. communication_message
Individual messages within sessions.

**Key Fields:**
- `session_id`: Parent session
- `message_type`: TEXT, IMAGE, VIDEO, AUDIO, FILE, LOCATION, TEMPLATE
- `from_address`, `to_address`: Message routing
- `message_body`: Message content
- `status`: QUEUED, SENT, DELIVERED, READ, FAILED, BOUNCED
- `external_message_id`: Provider's message ID
- `cost`: Usage cost tracking
- `sent_at`, `delivered_at`, `read_at`: Status timestamps
- `retry_count`, `max_retries`: Automatic retry logic

### 4. message_template
Reusable message templates for campaigns and automation.

**Key Fields:**
- `template_name`: Template identifier
- `channel_type`: Which channels use this template
- `template_body`: Message content with variables
- `template_variables`: JSON with variable definitions
- `subject`: For email templates
- `usage_count`: Analytics tracking
- `language`: For multi-language support

### 5. contact_communication_preference
User preferences for communication channels.

**Key Fields:**
- `contact_id`: Reference to contact
- `preferred_channel`: Primary communication channel
- `allow_telegram`, `allow_whatsapp`, `allow_sms`, `allow_email`: Channel permissions
- `opt_in_sms`, `opt_in_marketing`: Compliance flags
- `do_not_contact`: Global opt-out

### 6. bulk_message_campaign
Manages bulk messaging campaigns.

**Key Fields:**
- `campaign_name`: Campaign identifier
- `campaign_status`: DRAFT, SCHEDULED, RUNNING, COMPLETED, CANCELLED
- `total_recipients`: Campaign size
- `sent_count`, `delivered_count`, `failed_count`: Campaign metrics
- `estimated_cost`, `actual_cost`: Budget tracking
- `scheduled_at`: Campaign start time

### 7. bulk_message_recipient
Recipients within bulk campaigns.

**Key Fields:**
- `campaign_id`: Parent campaign
- `recipient_address`: Email/phone/ID
- `template_variables`: Personalization data
- `send_status`: PENDING, SENT, DELIVERED, FAILED, SKIPPED

### 8. message_automation_rule
Rules for automated message triggering.

**Key Fields:**
- `rule_type`: TRIGGER_ON_EVENT, SCHEDULED, WORKFLOW, DRIP_CAMPAIGN
- `trigger_event`: Event name to trigger automation
- `condition_json`: Conditions for rule execution
- `action_json`: Actions to perform

### 9. scheduled_message
Messages scheduled for future delivery.

**Key Fields:**
- `scheduled_for`: Delivery timestamp
- `recurrence_pattern`: once, daily, weekly, monthly
- `recurrence_end_date`: When to stop recurring

### 10. communication_webhook_log
Audit trail for provider webhook events.

**Key Fields:**
- `webhook_event_type`: Event type from provider
- `webhook_payload`: Full event data
- `webhook_signature`: For signature verification
- `processing_status`: PENDING, PROCESSED, FAILED

### 11. communication_analytics
Aggregated metrics for reporting.

**Key Fields:**
- `metric_date`: Date for metrics
- `total_messages`, `sent_messages`, `delivered_messages`: Message counts
- `total_cost`: Channel usage cost
- `engagement_rate`: Read/response metrics
- `unique_recipients`: Unique contact count

### 12. user_communication_permission
User-level permissions for channels.

**Key Fields:**
- `user_id`: User authorization
- `can_send_*`: Per-channel send permissions
- `can_manage_templates`, `can_manage_campaigns`: Administrative permissions

### 13-16. Additional Tables
- `communication_attachment`: File attachments in messages
- `call_compliance_rule`: GDPR/CCPA compliance tracking
- `call_rate_config`: Pricing per destination
- `call_usage_billing`: Usage-based billing

## API Endpoints

### Communication Channels

#### Create Channel
```http
POST /api/v1/communication/channels
Content-Type: application/json
X-Tenant-ID: <tenant_id>

{
  "channel_type": "TELEGRAM",
  "channel_name": "Support Bot",
  "api_provider": "telegram",
  "auth_token": "your_bot_token",
  "sender_id": "your_bot_id",
  "webhook_url": "https://your-domain.com/webhooks/telegram"
}
```

**Response:** `201 Created`
```json
{
  "id": "ch_123456",
  "tenant_id": "tenant_1",
  "channel_type": "TELEGRAM",
  "channel_name": "Support Bot",
  "api_provider": "telegram",
  "is_active": true,
  "created_at": "2025-12-03T10:00:00Z"
}
```

#### Get Channel
```http
GET /api/v1/communication/channels/{id}
X-Tenant-ID: <tenant_id>
```

#### List Channels
```http
GET /api/v1/communication/channels?limit=50&offset=0
X-Tenant-ID: <tenant_id>
```

### Messages

#### Send Message
```http
POST /api/v1/communication/messages
Content-Type: application/json
X-Tenant-ID: <tenant_id>

{
  "channel_type": "EMAIL",
  "to_address": "customer@example.com",
  "message_body": "Hello! How can we help?",
  "message_type": "TEXT",
  "subject": "Support Reply",
  "contact_id": "contact_123",
  "lead_id": "lead_456",
  "priority": 1
}
```

**Response:** `201 Created`
```json
{
  "id": "msg_789",
  "session_id": "sess_123",
  "channel_type": "EMAIL",
  "from_address": "support@company.com",
  "to_address": "customer@example.com",
  "message_body": "Hello! How can we help?",
  "status": "SENT",
  "external_message_id": "twilio_msg_id",
  "sent_at": "2025-12-03T10:05:00Z",
  "created_at": "2025-12-03T10:05:00Z"
}
```

### Bulk Messaging

#### Send Bulk Messages
```http
POST /api/v1/communication/bulk-send
Content-Type: application/json
X-Tenant-ID: <tenant_id>

{
  "campaign_name": "Holiday Sale Announcement",
  "channel_type": "SMS",
  "template_id": "template_holiday",
  "recipients": [
    {
      "address": "+15551234567",
      "name": "John Doe",
      "contact_id": "contact_1",
      "template_variables": {
        "discount": "20%",
        "code": "HOLIDAY20"
      }
    }
  ],
  "schedule_for": "2025-12-10T10:00:00Z"
}
```

**Response:** `201 Created`
```json
{
  "id": "camp_001",
  "campaign_name": "Holiday Sale Announcement",
  "channel_type": "SMS",
  "total_recipients": 1,
  "campaign_status": "SCHEDULED",
  "estimated_cost": 0.50,
  "scheduled_at": "2025-12-10T10:00:00Z"
}
```

### Message Templates

#### Create Template
```http
POST /api/v1/communication/templates
Content-Type: application/json
X-Tenant-ID: <tenant_id>

{
  "template_name": "Welcome Email",
  "channel_type": "EMAIL",
  "template_body": "Welcome {{first_name}}! Your account is ready.",
  "subject": "Welcome to {{company_name}}",
  "language": "en",
  "template_variables": {
    "first_name": "string",
    "company_name": "string"
  }
}
```

#### List Templates
```http
GET /api/v1/communication/templates?limit=20&offset=0
X-Tenant-ID: <tenant_id>
```

### Contact Preferences

#### Update Preferences
```http
POST /api/v1/communication/contacts/{contact_id}/preferences
Content-Type: application/json
X-Tenant-ID: <tenant_id>

{
  "email_address": "john@example.com",
  "phone_number": "+15551234567",
  "telegram_id": "telegram_user_123",
  "whatsapp_number": "+15551234567",
  "preferred_channel": "WHATSAPP",
  "allow_telegram": true,
  "allow_whatsapp": true,
  "allow_sms": true,
  "allow_email": true,
  "opt_in_marketing": true
}
```

#### Get Preferences
```http
GET /api/v1/communication/contacts/{contact_id}/preferences
X-Tenant-ID: <tenant_id>
```

### Conversations

#### Get Session
```http
GET /api/v1/communication/sessions/{session_id}
X-Tenant-ID: <tenant_id>
```

#### List Sessions
```http
GET /api/v1/communication/sessions?limit=50&offset=0&channel_type=EMAIL&status=ACTIVE
X-Tenant-ID: <tenant_id>
```

### Analytics

#### Get Communication Stats
```http
GET /api/v1/communication/stats?start_date=2025-11-01&end_date=2025-12-03
X-Tenant-ID: <tenant_id>
```

**Response:**
```json
{
  "analytics": [
    {
      "channel_type": "EMAIL",
      "metric_date": "2025-12-03",
      "total_messages": 150,
      "sent_messages": 148,
      "delivered_messages": 145,
      "failed_messages": 2,
      "read_messages": 89,
      "total_cost": 15.50,
      "engagement_rate": 61.4
    }
  ],
  "start_date": "2025-11-01",
  "end_date": "2025-12-03"
}
```

### Webhooks

#### Receive Webhook
```http
POST /api/v1/communication/webhooks/{provider-type}
Content-Type: application/json
X-Webhook-Signature: sha256=...
X-Tenant-ID: <tenant_id>

{
  "event": "delivered",
  "message_id": "twilio_msg_id",
  "timestamp": "2025-12-03T10:05:30Z",
  "status": "delivered"
}
```

## Provider Integration

### Email Providers

#### SendGrid
```json
{
  "channel_type": "EMAIL",
  "api_provider": "sendgrid",
  "api_key": "SG.xxxxx",
  "sender_id": "noreply@company.com"
}
```

#### Mailgun
```json
{
  "channel_type": "EMAIL",
  "api_provider": "mailgun",
  "api_key": "key-xxxxx",
  "account_id": "sandbox123.mailgun.org",
  "sender_id": "noreply@company.com"
}
```

### SMS Providers

#### Twilio
```json
{
  "channel_type": "SMS",
  "api_provider": "twilio",
  "account_id": "AC123456",
  "api_key": "AC123456",
  "api_secret": "auth_token",
  "sender_id": "+15551234567"
}
```

#### AWS SNS
```json
{
  "channel_type": "SMS",
  "api_provider": "aws_sns",
  "api_key": "AKIAIOSFODNN7EXAMPLE",
  "api_secret": "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY",
  "sender_id": "+15551234567"
}
```

#### Vonage (Nexmo)
```json
{
  "channel_type": "SMS",
  "api_provider": "vonage",
  "api_key": "xxxxx",
  "api_secret": "yyyyy",
  "sender_id": "Company"
}
```

### Telegram
```json
{
  "channel_type": "TELEGRAM",
  "api_provider": "telegram",
  "auth_token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
  "sender_id": "123456789"
}
```

### WhatsApp
```json
{
  "channel_type": "WHATSAPP",
  "api_provider": "whatsapp_business",
  "account_id": "100123456789",
  "api_key": "xxx",
  "api_secret": "yyy",
  "sender_id": "+15551234567"
}
```

## Features

### 1. Multi-Channel Support
- **Email**: SendGrid, Mailgun
- **SMS**: Twilio, AWS SNS, Vonage
- **Telegram**: Native Bot API
- **WhatsApp**: WhatsApp Business API
- Extensible for additional providers

### 2. Unified Inbox
- Single conversation view across all channels
- Channel routing based on contact preferences
- Fallback to secondary channels

### 3. Message Templates
- Reusable templates per channel
- Variable substitution ({{name}}, {{code}}, etc.)
- Multi-language support
- Template usage analytics

### 4. Bulk Messaging
- Campaign management
- Recipient segmentation
- Scheduled delivery
- Personalization support
- Cost tracking

### 5. Contact Preferences
- Per-contact communication preferences
- Channel permission management
- Opt-in/opt-out tracking
- Do Not Contact compliance

### 6. Automation
- Event-triggered messaging
- Scheduled recurring messages
- Workflow automation
- Drip campaigns

### 7. Compliance & Security
- GDPR/CCPA compliance tracking
- Webhook signature verification
- Encrypted credential storage
- Audit trail for all messages
- User-level permissions

### 8. Analytics & Reporting
- Daily message metrics
- Channel performance comparison
- Delivery rate tracking
- Engagement analytics
- Cost analysis
- Response time metrics

### 9. Webhooks
- Real-time event notifications
- Message status updates (sent, delivered, read)
- Signature verification
- Automatic retry on failure

### 10. Rich Media Support
- Image attachments
- Video messages
- Audio attachments
- File uploads
- Location sharing (provider-dependent)

## Implementation Guide

### Step 1: Create Communication Channels

```bash
curl -X POST http://localhost:8080/api/v1/communication/channels \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_1" \
  -d '{
    "channel_type": "EMAIL",
    "channel_name": "SendGrid Email",
    "api_provider": "sendgrid",
    "api_key": "SG.your_sendgrid_key",
    "sender_id": "noreply@yourcompany.com"
  }'
```

### Step 2: Create Message Templates

```bash
curl -X POST http://localhost:8080/api/v1/communication/templates \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_1" \
  -d '{
    "template_name": "Order Confirmation",
    "channel_type": "EMAIL",
    "template_body": "Thank you for your order #{{order_id}}",
    "subject": "Order Confirmation - #{{order_id}}",
    "template_variables": {"order_id": "string"}
  }'
```

### Step 3: Update Contact Preferences

```bash
curl -X POST http://localhost:8080/api/v1/communication/contacts/contact_123/preferences \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_1" \
  -d '{
    "email_address": "john@example.com",
    "phone_number": "+15551234567",
    "preferred_channel": "EMAIL"
  }'
```

### Step 4: Send Messages

```bash
curl -X POST http://localhost:8080/api/v1/communication/messages \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_1" \
  -d '{
    "channel_type": "EMAIL",
    "to_address": "john@example.com",
    "message_body": "Order shipped! Track at: https://track.example.com/123",
    "contact_id": "contact_123"
  }'
```

### Step 5: Set Up Webhooks

Configure provider webhooks to point to:
```
https://your-api.com/api/v1/communication/webhooks/{provider-type}
```

For example:
- SendGrid: `https://your-api.com/api/v1/communication/webhooks/email`
- Twilio: `https://your-api.com/api/v1/communication/webhooks/sms`

## Configuration

### Database Setup
```sql
-- Run migration 020
mysql -h localhost -u root -p your_database < migrations/020_multi_channel_communication.sql
```

### Environment Variables
```bash
# Email Providers
SENDGRID_API_KEY=SG.xxxxx
MAILGUN_API_KEY=key-xxxxx
MAILGUN_DOMAIN=sandbox123.mailgun.org

# SMS Providers
TWILIO_ACCOUNT_SID=AC123456
TWILIO_AUTH_TOKEN=xxxxx
VONAGE_API_KEY=xxxxx
VONAGE_API_SECRET=yyyyy
AWS_SNS_REGION=us-east-1

# Telegram
TELEGRAM_BOT_TOKEN=123456:ABC-DEF...

# WhatsApp
WHATSAPP_BUSINESS_ACCOUNT_ID=100123456789
WHATSAPP_API_KEY=xxxxx

# Webhook
WEBHOOK_SIGNATURE_SECRET=your_secret_key
```

### Application Configuration
```go
// In main.go
multiChannelService := services.NewMultiChannelCommunicationService(db, logger)
handler := handlers.NewMultiChannelCommunicationHandler(multiChannelService, logger)

// Register routes
router.HandleFunc("POST", "/api/v1/communication/channels", handler.CreateCommunicationChannel)
router.HandleFunc("GET", "/api/v1/communication/channels/{id}", handler.GetCommunicationChannel)
router.HandleFunc("GET", "/api/v1/communication/channels", handler.ListCommunicationChannels)
router.HandleFunc("POST", "/api/v1/communication/messages", handler.SendMessage)
router.HandleFunc("POST", "/api/v1/communication/bulk-send", handler.SendBulkMessages)
router.HandleFunc("POST", "/api/v1/communication/templates", handler.CreateMessageTemplate)
router.HandleFunc("GET", "/api/v1/communication/templates/{id}", handler.GetMessageTemplate)
router.HandleFunc("GET", "/api/v1/communication/templates", handler.ListMessageTemplates)
router.HandleFunc("POST", "/api/v1/communication/contacts/{id}/preferences", handler.UpdateContactPreference)
router.HandleFunc("GET", "/api/v1/communication/contacts/{id}/preferences", handler.GetContactPreference)
router.HandleFunc("POST", "/api/v1/communication/webhooks/{provider-type}", handler.HandleWebhookEvent)
router.HandleFunc("GET", "/api/v1/communication/sessions/{id}", handler.GetCommunicationSession)
router.HandleFunc("GET", "/api/v1/communication/sessions", handler.ListCommunicationSessions)
router.HandleFunc("GET", "/api/v1/communication/stats", handler.GetCommunicationStats)
```

## Webhook Handling

### Expected Webhook Format

**From SendGrid:**
```json
{
  "email": "example@example.com",
  "timestamp": 1513299569,
  "smtp-id": "<message id string>",
  "event": "delivered",
  "sg_event_id": "event_id",
  "sg_message_id": "message_id"
}
```

**From Twilio:**
```json
{
  "MessageSid": "SMxxxxxxxxxxxxxxxxxxxxx",
  "MessageStatus": "delivered",
  "ChannelToAddress": "+15551234567",
  "EventType": "MessageStatus"
}
```

**From Telegram:**
```json
{
  "update_id": 10000000,
  "message": {
    "message_id": 1,
    "from": {"id": 123, "is_bot": false},
    "chat": {"id": 123},
    "text": "Hello, world!"
  }
}
```

## Examples

### React Frontend Example

```typescript
// Send Message
async function sendMessage(channelType: string, toAddress: string, body: string) {
  const response = await fetch('/api/v1/communication/messages', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Tenant-ID': tenantId,
    },
    body: JSON.stringify({
      channel_type: channelType,
      to_address: toAddress,
      message_body: body,
      message_type: 'TEXT',
    }),
  });
  return response.json();
}

// Get Conversations
async function listConversations() {
  const response = await fetch('/api/v1/communication/sessions?limit=20&offset=0', {
    headers: {
      'X-Tenant-ID': tenantId,
    },
  });
  return response.json();
}

// Get Contact Preferences
async function getContactPref(contactId: string) {
  const response = await fetch(`/api/v1/communication/contacts/${contactId}/preferences`, {
    headers: {
      'X-Tenant-ID': tenantId,
    },
  });
  return response.json();
}
```

### Go Example

```go
service := services.NewMultiChannelCommunicationService(db, logger)

// Send Email
msg, err := service.SendMessage(ctx, tenantID, &models.SendMessageRequest{
    ChannelType: "EMAIL",
    ToAddress:   "user@example.com",
    MessageBody: "Hello from system",
    Subject:     "Notification",
})

// Send Bulk SMS
campaign, err := service.SendBulkMessages(ctx, tenantID, &models.BulkSendRequest{
    CampaignName: "Promo",
    ChannelType:  "SMS",
    Recipients: []models.BulkRecipient{
        {
            Address: "+15551234567",
            Name:    "John",
        },
    },
})

// Get Analytics
analytics, err := service.GetCommunicationAnalytics(
    ctx, tenantID, 
    time.Now().AddDate(0, 0, -30), 
    time.Now(),
)
```

## Best Practices

### 1. Channel Selection
- Use contact preferences to determine primary channel
- Implement fallback channels if primary fails
- Track user feedback on preferred channels

### 2. Template Management
- Create templates for common scenarios (welcome, confirmation, reminder)
- Use consistent variable naming across channels
- Test templates across different devices and email clients

### 3. Rate Limiting
- Respect provider rate limits
- Implement exponential backoff for retries
- Monitor cost usage to prevent unexpected bills

### 4. Compliance
- Always verify opt-in before sending marketing messages
- Provide easy opt-out mechanisms
- Log all communication for audit trails
- Implement DNC (Do Not Contact) list checking

### 5. Monitoring
- Track delivery rates by channel
- Monitor message failures and error codes
- Set up alerts for high failure rates
- Regular webhook signature verification

### 6. Security
- Store API keys securely (use vault/secrets manager)
- Verify webhook signatures before processing
- Use HTTPS for all API calls
- Implement rate limiting on API endpoints

### 7. Performance
- Use bulk sending for campaigns
- Implement message queue for async processing
- Cache templates and channel configs
- Batch webhook processing

### 8. Cost Optimization
- Choose cost-effective channels per use case
- Monitor per-channel costs
- Implement delivery windows to reduce premium rates
- Track unused channels and consolidate providers

## Troubleshooting

### Messages Not Sending
1. Verify channel credentials are correct
2. Check webhook signature validation
3. Review error logs in communication_webhook_log table
4. Validate recipient address format for channel type

### Delivery Failures
1. Check contact's opt-in status
2. Verify contact address is valid
3. Review provider-specific error messages
4. Check for rate limiting issues

### Missing Webhook Updates
1. Verify webhook URL is publicly accessible
2. Check firewall/security group rules
3. Validate webhook signature secret
4. Review webhook_log table for processing errors

## Next Steps

1. **Frontend UI**: Build inbox UI for conversations
2. **Agent Dashboard**: Create agent workspace
3. **Automation**: Implement workflow engine
4. **Analytics**: Build reporting dashboards
5. **Integration**: Connect with CRM/ERP systems
6. **Performance**: Set up message queuing (Redis/RabbitMQ)
7. **Monitoring**: Implement observability stack

## Support

For issues, questions, or feature requests, contact the development team or file an issue in the project repository.
