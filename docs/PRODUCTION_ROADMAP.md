# VYOMTECH Lead Management Module Specification

## 1. Overview

The VYOMTECH Lead Management Module is the revenue engine of the platform. It captures, qualifies, nurtures, and converts prospects into business opportunities while integrating with automation, campaigns, and customer lifecycle processes.

This module is directly informed by the patterns proven in the propertyscouts repository, especially:
- lead capture from forms
- n8n-like workflow automation
- WhatsApp and email follow-up
- tiered qualification and funnel progression
- source tracking and campaign attribution
- conversion reporting

## 2. Goals

- capture leads from all sources
- classify and qualify leads rapidly
- assign leads to team members or queues
- automate nurture sequences
- track communication history
- manage opportunities and next steps
- measure conversion rates and sales efficiency
- connect with CRM, marketing, and ERP workflows

## 3. Users and actors

- marketing user
- sales rep
- sales manager
- account owner
- support agent
- tenant admin
- system automation

## 4. Core domain model

### Entities
- Lead
- LeadSource
- LeadStatus
- LeadNote
- Campaign
- CampaignMember
- Contact
- Account
- Opportunity
- Task
- Activity
- Product
- Quote
- Contract

## 5. Exact tables

### 5.1 leads
- id
- tenant_id
- company_id
- source_id
- status_id
- owner_user_id
- account_id (nullable)
- contact_id (nullable)
- campaign_id (nullable)
- first_name
- last_name
- email
- phone
- alternate_phone
- company_name
- job_title
- city
- state
- country
- lead_score
- priority
- assigned_at
- created_at
- updated_at
- converted_at
- source_url
- utm_source
- utm_medium
- utm_campaign
- notes

### 5.2 lead_sources
- id
- tenant_id
- name
- type
- description
- is_active
- created_at

### 5.3 lead_statuses
- id
- tenant_id
- name
- workflow_stage
- sequence_order
- is_default
- created_at

### 5.4 lead_notes
- id
- tenant_id
- lead_id
- user_id
- note_type
- content
- created_at

### 5.5 campaigns
- id
- tenant_id
- name
- type
- start_date
- end_date
- owner_user_id
- budget
- status
- created_at

### 5.6 campaign_members
- id
- campaign_id
- lead_id
- contact_id
- status
- created_at

### 5.7 contacts
- id
- tenant_id
- first_name
- last_name
- email
- phone
- company_name
- title
- address
- city
- state
- country
- created_at

### 5.8 accounts
- id
- tenant_id
- name
- industry
- website
- billing_email
- phone
- assigned_owner_id
- account_type
- created_at

### 5.9 opportunities
- id
- tenant_id
- lead_id
- account_id
- owner_user_id
- name
- amount
- stage
- probability
- expected_close_date
- created_at
- updated_at

### 5.10 tasks
- id
- tenant_id
- lead_id
- opportunity_id
- assigned_to_user_id
- task_type
- title
- description
- due_date
- status
- created_at

### 5.11 activities
- id
- tenant_id
- lead_id
- opportunity_id
- user_id
- activity_type
- subject
- details
- occurred_at

## 6. Lead lifecycle states

New
- incoming lead
- no qualification yet

Qualified
- valid contact info and clear need

Nurturing
- automated sequence in progress

Contacted
- first interaction happened

Follow-up
- further contact scheduled

Proposal
- quote or proposal in process

Negotiation
- deal discussion active

Won
- deal closed

Lost
- deal rejected or dead

## 7. Workflow definitions

### 7.1 Inbound lead capture workflow
1. website form or API submits lead payload
2. lead is validated for required fields
3. lead is deduplicated against email/phone
4. lead is enriched with source metadata
5. automation rule creates task to assign owner
6. lead score is calculated
7. email or WhatsApp nurture sequence starts
8. lead status changes to Qualified or Nurturing

### 7.2 Qualification workflow
1. sales rep reviews lead
2. lead is classified by fit, urgency, and buying intent
3. score is updated
4. owner is assigned or changed
5. tasks are created for follow-up
6. opportunity can be created if suitable

### 7.3 Nurture workflow
1. initial acknowledgment is sent
2. email/WhatsApp sequence is triggered based on source and persona
3. next steps are scheduled automatically
4. if no response occurs, re-engagement rules fire
5. all activities are logged in audit trail

### 7.4 Opportunity conversion workflow
1. lead qualifies as sales-ready
2. opportunity record is created
3. quote is generated
4. contract is issued
5. deal progress tracked through stages
6. closed-won or closed-lost events update lead state

## 8. Automation rules

### Rule examples
- if lead comes from website and score > 70 -> assign sales rep within 5 minutes
- if lead has no contact after 24 hours -> send follow-up email
- if lead has high value and positive intent -> create opportunity
- if no response after 3 follow-ups -> move to dormant
- if deal is won -> create customer account and invoice flow

## 9. APIs

### 9.1 Leads
- POST /api/v1/leads
- GET /api/v1/leads
- GET /api/v1/leads/:id
- PATCH /api/v1/leads/:id
- DELETE /api/v1/leads/:id
- POST /api/v1/leads/:id/assign
- POST /api/v1/leads/:id/convert

### 9.2 Activities
- POST /api/v1/leads/:id/activities
- GET /api/v1/leads/:id/activities

### 9.3 Forms and webhooks
- POST /api/v1/webhooks/lead-capture
- POST /api/v1/lead-capture/form

### 9.4 Opportunities
- POST /api/v1/opportunities
- GET /api/v1/opportunities
- PATCH /api/v1/opportunities/:id

## 10. Lead scoring model

A simple scoring model should include:
- source quality
- buyer intent
- engagement level
- urgency
- fit to ICP
- whether contact was established
- budget signal
- company size or deal size

Example scoring:
- source website = +15
- high-intent form = +25
- replied to first follow-up = +10
- requested pricing = +20
- matched ICP = +20
- has budget info = +15

Thresholds:
- 0–29: cold / nurture
- 30–59: warm / qualifying
- 60–79: sales-ready
- 80+: hot lead

## 11. Reporting requirements

- leads by source
- conversion by salesperson
- conversion by campaign
- response rate by channel
- lead aging report
- pipeline value by stage
- lost reason analysis
- revenue from closed-won opportunities

## 12. Integration requirements

- WhatsApp automation integration
- email integration
- calendar integration
- webhooks from external landing pages
- CRM sync to accounting when won
- dashboard analytics feed to exec views

## 13. Recommended implementation order

1. lead schema and validation
2. lead source/status configuration
3. lead capture endpoint and deduplication
4. assignment and ownership rules
5. lead notes and activities
6. nurture automation
7. opportunity creation and pipeline stage tracking
8. analytics and reporting
9. outbound campaign management

## 14. Final outcome

The VYOMTECH Lead Management Module will become the single revenue operations center for all customer acquisition and sales processes. It will connect to ERP, LMS, and support workflows and serve as the foundation for higher-level business intelligence.
