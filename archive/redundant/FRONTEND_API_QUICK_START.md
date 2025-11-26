# Frontend API Integration - Quick Start Guide

**Last Updated:** November 22, 2025  
**Status:** ✅ Ready for Development  
**Build Status:** ✅ Backend Clean (0 errors)

---

## Quick Links

- 📚 **Complete API Reference:** `COMPLETE_API_REFERENCE.md`
- 🔗 **Integration Guide:** `FRONTEND_BACKEND_INTEGRATION.md`
- 📋 **Audit Summary:** `API_AUDIT_AND_FRONTEND_UPDATE_SUMMARY.md`

---

## Get Started in 5 Minutes

### 1. Start Backend
```bash
cd /c/Users/Skydaddy/Desktop/Developement
go run cmd/main.go
```
✅ Runs on `http://localhost:8080`

### 2. Start Frontend
```bash
cd frontend
npm install  # Only if needed
npm run dev
```
✅ Runs on `http://localhost:3000`

### 3. API is Ready
```
✅ Backend: http://localhost:8080/api/v1
✅ Health: http://localhost:8080/health
✅ Frontend: http://localhost:3000
```

---

## Import Services in Components

```typescript
// Option 1: Import specific service
import { leadService } from '@/services/api'

// Option 2: Import hook
import { useLeads } from '@/hooks/useLeads'

// Usage
const leads = await leadService.listLeads()
// OR
const { leads, fetchLeads } = useLeads()
```

---

## 14 Services Available

```typescript
import {
  authService,           // Login, register, token validation
  agentService,          // Agent management
  tenantService,         // Tenant operations
  leadService,           // Lead CRUD
  callService,           // Call management
  campaignService,       // Campaign CRUD
  webSocketService,      // Real-time WebSocket
  analyticsService,      // Reports and analytics
  automationService,     // Lead scoring & routing
  communicationService,  // Multi-channel messaging
  gamificationService,   // Points, badges, challenges
  complianceService,     // RBAC, audit, GDPR
  aiService,             // AI queries
  apiClient,             // Raw Axios client
} from '@/services/api'
```

---

## 7 React Hooks Ready to Use

### 1. useLeads()
```typescript
import { useLeads } from '@/hooks/useLeads'

function LeadsList() {
  const { leads, loading, fetchLeads } = useLeads()
  
  useEffect(() => {
    fetchLeads(1, 20)
  }, [])
  
  return (
    <div>
      {loading && <p>Loading...</p>}
      {leads.map(lead => <div key={lead.id}>{lead.name}</div>)}
    </div>
  )
}
```

### 2. useCalls()
```typescript
import { useCalls } from '@/hooks/useCalls'

function CallList() {
  const { calls, createCall, endCall } = useCalls()
  // Use the hook methods
}
```

### 3. useGamification()
```typescript
import { useGamification } from '@/hooks/useGamification'

function GamificationProfile() {
  const { profile, badges, leaderboard, fetchProfile } = useGamification()
  // Profile data and operations
}
```

### 4. useAnalytics()
```typescript
import { useAnalytics } from '@/hooks/useAnalytics'

function AnalyticsDashboard() {
  const { generateReport, exportReport } = useAnalytics()
  // Analytics operations
}
```

### 5. useAutomation()
```typescript
import { useAutomation } from '@/hooks/useAutomation'

function AutomationPanel() {
  const { rankLeads, routeLeadToAgent } = useAutomation()
  // Automation operations
}
```

### 6. useCompliance()
```typescript
import { useCompliance } from '@/hooks/useCompliance'

function ComplianceDashboard() {
  const { auditLogs, exportUserData } = useCompliance()
  // Compliance operations
}
```

---

## Common API Calls

### Fetch Leads
```typescript
// Using service
const leads = await leadService.listLeads(1, 20)

// Using hook
const { fetchLeads } = useLeads()
useEffect(() => {
  fetchLeads(1, 20)
}, [])
```

### Create a Lead
```typescript
// Using service
const newLead = await leadService.createLead({
  name: 'John Doe',
  email: 'john@example.com',
  source: 'campaign'
})

// Using hook
const { createLead } = useLeads()
const newLead = await createLead({ name, email, source })
```

### Award Points
```typescript
// Using service
await gamificationService.awardPoints(userId, 50, 'Call completed')

// Using hook
const { awardPoints } = useGamification()
await awardPoints(userId, 50, 'Call completed')
```

### Generate Report
```typescript
// Using service
const report = await analyticsService.generateReport(
  'agent_performance',
  '2025-01-01',
  '2025-01-31'
)

// Using hook
const { generateReport } = useAnalytics()
const report = await generateReport('agent_performance', startDate, endDate)
```

### Calculate Lead Score
```typescript
// Using service
const score = await automationService.calculateLeadScore(leadId)

// Using hook
const { calculateLeadScore } = useAutomation()
const score = await calculateLeadScore(leadId)
```

### Export User Data (GDPR)
```typescript
// Using service
const data = await complianceService.exportUserData()

// Using hook
const { exportUserData } = useCompliance()
const data = await exportUserData()
```

---

## Error Handling

### Using Services
```typescript
try {
  const leads = await leadService.listLeads()
} catch (error) {
  console.error('Failed to fetch leads:', error)
}
```

### Using Hooks
```typescript
const { leads, error, loading } = useLeads()

return (
  <div>
    {error && <p>Error: {error}</p>}
    {loading && <p>Loading...</p>}
    {/* Display leads */}
  </div>
)
```

---

## Feature Examples

### Real-time WebSocket
```typescript
import { webSocketService } from '@/services/api'
import { useAuth } from '@/hooks/useAuth'

function RealtimeUpdates() {
  const { token } = useAuth()
  
  useEffect(() => {
    const ws = webSocketService.connect(token)
    ws.onmessage = (event) => {
      console.log('Real-time update:', event.data)
    }
  }, [token])
}
```

### Multi-Channel Communication
```typescript
import { communicationService } from '@/services/api'

async function sendNotification() {
  // Send SMS
  await communicationService.sendMessage(
    '+1234567890',
    'sms',
    'template_123'
  )
  
  // Send Email
  await communicationService.sendMessage(
    'user@example.com',
    'email',
    'template_456'
  )
}
```

### Lead Scoring & Routing
```typescript
import { automationService } from '@/services/api'

async function handleNewLead(leadId) {
  // Calculate score
  const score = await automationService.calculateLeadScore(leadId)
  console.log(`Lead score: ${score.score}`)
  
  // Route to agent
  const routing = await automationService.routeLeadToAgent(leadId)
  console.log(`Assigned to agent: ${routing.agent_id}`)
}
```

### Compliance & Audit
```typescript
import { complianceService } from '@/services/api'

async function getAuditTrail(userId) {
  const logs = await complianceService.getAuditLogs(userId)
  console.log('User actions:', logs)
  
  const report = await complianceService.getComplianceReport()
  console.log('Compliance report:', report)
}
```

---

## API Endpoints by Feature

| Feature | Base URL | Endpoints |
|---------|----------|-----------|
| WebSocket | `/api/v1/ws` | 2 |
| Analytics | `/api/v1/analytics` | 4 |
| Automation | `/api/v1/automation` | 6 |
| Communication | `/api/v1/communication` | 5 |
| Gamification | `/api/v1/gamification` | 17 |
| Compliance | `/api/v1/compliance` | 14 |
| Core (Leads, Calls, Campaigns, etc.) | `/api/v1/` | 42 |
| **Total** | | **100+** |

---

## Common Issues & Solutions

### Issue: 401 Unauthorized
**Solution:** Token expired. Login again or check localStorage.
```typescript
import { authService } from '@/services/api'
await authService.login(email, password)
```

### Issue: CORS Error
**Solution:** Make sure backend is running on port 8080.
```bash
go run cmd/main.go  # Check it's listening on :8080
```

### Issue: WebSocket Connection Failed
**Solution:** Check token is valid and WebSocket path is correct.
```typescript
const ws = webSocketService.connect(validToken)
```

### Issue: Database Connection Error
**Solution:** Ensure MySQL is running and connection details are correct.
```bash
# Check MySQL
mysql -u root -p  # Enter password when prompted
```

---

## TypeScript Support

All services have proper TypeScript types:

```typescript
import { 
  LeadType, 
  CallType, 
  CampaignType 
} from '@/types'

// Type-safe API calls
const lead: LeadType = await leadService.createLead({
  name: 'John',
  email: 'john@example.com'
})

// Hook returns typed data
const { leads }: { leads: LeadType[] } = useLeads()
```

---

## Environment Setup

### Create `.env.local` in Frontend
```bash
cd frontend
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local
```

### Verify It Works
```bash
npm run dev  # Should connect to backend on 8080
```

---

## Documentation

📚 **For detailed information, see:**
- `COMPLETE_API_REFERENCE.md` - All 100+ endpoints
- `FRONTEND_BACKEND_INTEGRATION.md` - Integration patterns
- `API_AUDIT_AND_FRONTEND_UPDATE_SUMMARY.md` - What was done

---

## Support

For issues or questions:
1. Check the API reference for endpoint details
2. Review hook examples in documentation
3. Check browser console for error messages
4. Verify backend is running: `http://localhost:8080/health`

---

**Everything is ready to go! Start building! 🚀**
