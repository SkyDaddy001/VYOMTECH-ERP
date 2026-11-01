# VYOMTECH ERP + LMS

VYOMTECH-ERP is now the master platform for the entire VYOMTECH business operating system.

This repository is the canonical home for the complete VYOMTECH ERP + LMS architecture, combining:
- the ERP foundation already developed in this repo
- the lead management and automation patterns from the propertyscouts project
- the company-level operating system vision from company-os
- the growth and customer acquisition patterns from qucogroup_blog

The decision is clear: we do not build fragmented SaaS products in separate repos. We consolidate everything into VYOMTECH-ERP and make it the unified operating system for business, learning, sales, and operations.

## Mission

Build one complete, modular, multi-tenant business platform that supports:
- sales and CRM
- finance and accounting
- HR and payroll
- procurement and inventory
- projects and assets
- LMS and training
- support and service operations
- analytics and AI-powered insights
- workflow automation and integrations

## Core strategy

VYOMTECH-ERP will be structured as a single SaaS platform with modular domains, not as disconnected apps.

Principles:
- single tenant-aware identity layer
- shared audit, permissions, and workflow engine
- domain-driven modules with clean APIs
- API-first integration architecture
- multi-tenant and white-label ready
- modular LMS and ERP features under one platform
- automation-first operations

## Reused learnings from other repos

### propertyscouts
Provides the strongest lead management and sales automation patterns:
- lead capture
- lead scoring and assignment
- nurturing workflows
- WhatsApp/email automation
- customer persona-based funnels
- campaign tracking
- pipeline conversion analytics

### company-os
Provides the operating system vision:
- unified dashboard
- company command center
- cross-functional visibility
- executive overview

### qucogroup_blog
Provides the growth and customer acquisition model:
- content-driven funnels
- marketing and onboarding flows
- audience segmentation
- branded properties and thought leadership

## What this repo becomes

VYOMTECH-ERP becomes:
- ERP core for operations and finance
- CRM + sales pipeline engine
- HRMS + payroll
- procurement + inventory
- LMS + academy platform
- service desk + support
- analytics + executive dashboards
- workflow orchestration layer
- integration gateway for all business systems

## High-level architecture

- Applications: web app, admin portal, customer portal, LMS frontend
- API Layer: REST + GraphQL-ready service layer
- Domain Modules: CRM, HR, finance, procurement, inventory, LMS, support, analytics
- Workflow Engine: n8n-style automation and custom rules
- Data Layer: PostgreSQL + Prisma + Redis + object storage
- Identity Layer: multi-tenant auth, RBAC, SSO
- Integrations: email, WhatsApp, payment, analytics, ERP connectors

## Documentation set

- docs/ERP_LMS_ARCHITECTURE_BLUEPRINT.md — complete system architecture blueprint
- docs/LEAD_MANAGEMENT_MODULE_SPEC.md — exact lead workflow and schema specification
- docs/PRODUCTION_ROADMAP.md — roadmap to finish the repo and launch production

## Immediate objective

The near-term focus is:
1. complete the ERP core foundation
2. finalize the lead management module
3. merge LMS into the main platform
4. connect workflows and automation
5. harden security, audit, and multi-tenant isolation
6. build analytics and executive dashboards

## Recommended platform direction

Build VYOMTECH as a modular enterprise operating system with these major layers:
- Core platform
- Customer growth and CRM
- Finance and accounting
- People operations
- Learning and enablement
- Operations and logistics
- Intelligence and automation

This is the correct strategic direction for VYOMTECH.
