# 🚀 Quick Start Guide - Testing with Demo Credentials

**Last Updated**: December 3, 2025

---

## 🎯 Quick Start (2 Minutes)

### Step 1: Start the Application
```bash
# Terminal 1: Start Backend (if not running)
cd /d/VYOMTECH-ERP
go run cmd/main.go

# Terminal 2: Start Frontend (if not running)
cd /d/VYOMTECH-ERP/frontend
npm run dev
```

### Step 2: Open Login Page
```
Open browser: http://localhost:3000/auth/login
```

### Step 3: Login with Demo Credentials
You'll see a green **"Demo Test Credentials"** card with all available accounts.

**Click any credential** to auto-fill and login instantly!

Or manually enter:
- Email: `demo@vyomtech.com`
- Password: `DemoPass@123`
- Click: **Sign In**

### Step 4: Explore Dashboard
After login, you'll have access to:
- Dashboard with 8 sample leads
- 4 call records with AI metrics
- 3 campaigns with performance data
- 2 agent profiles
- Complete multi-tenant demo

---

## 👥 Available Test Accounts

### 1. Admin User (Full Access)
```
Email:    demo@vyomtech.com
Password: DemoPass@123
Role:     Admin
Access:   Everything - system administration, all modules
```

### 2. Agent User (Call Management)
```
Email:    agent@vyomtech.com
Password: AgentPass@123
Role:     Agent
Access:   Call handling, leads, customer interactions
Skills:   Sales, Support, Billing
Status:   Online (available for calls)
Calls:    125 total calls, 4.6★ satisfaction
```

### 3. Manager User (Team Management)
```
Email:    manager@vyomtech.com
Password: ManagerPass@123
Role:     Supervisor/Manager
Access:   Team management, reporting, analytics
Skills:   Management, Training, Quality
Status:   Offline
Calls:    256 total calls, 4.8★ satisfaction
```

### 4. Sales User (Pipeline Management)
```
Email:    sales@vyomtech.com
Password: SalesPass@123
Role:     Sales User
Access:   Sales pipeline, leads, opportunities
Usage:    Sales module demonstration
```

### 5. HR User (Employee Management)
```
Email:    hr@vyomtech.com
Password: HRPass@123
Role:     HR Staff
Access:   Employee management, payroll, compliance
Usage:    HR module demonstration
```

---

## 📊 Sample Data Available

### Leads (8 Total)
| Name | Email | Status | Priority | Source | Notes |
|------|-------|--------|----------|--------|-------|
| John Smith | john.smith@company.com | New | High | Direct | 45 days old |
| Sarah Johnson | sarah.johnson@company.com | Contacted | Medium | Referral | 30 days old |
| Michael Brown | michael.brown@company.com | Qualified | High | Website | 15 days old |
| Emily Davis | emily.davis@company.com | Converted | High | Email | Recent |
| Robert Wilson | robert.wilson@company.com | New | Low | Phone | 2 days old |
| Jessica Martinez | jessica.martinez@company.com | Contacted | Medium | LinkedIn | Week old |
| James Taylor | james.taylor@company.com | Qualified | High | Partner | 20 days old |
| Amanda Anderson | amanda.anderson@company.com | Lost | Low | Trade Show | 60 days old |

### Calls (4 Sample Calls)
| Direction | Status | Duration | AI Used | Sentiment | Agent |
|-----------|--------|----------|---------|-----------|-------|
| Inbound | Completed | 7m 30s | Yes (OpenAI) | 4.2/5 | John Agent |
| Outbound | Completed | 10m 20s | Yes (GPT-4) | 3.8/5 | John Agent |
| Inbound | Completed | 6m 20s | No | 4.5/5 | John Agent |
| Outbound | Failed | — | Yes (OpenAI) | — | John Agent |

### Campaigns (3 Active Campaigns)
| Name | Type | Status | Recipients | Sent | Opened | Converted | Budget |
|------|------|--------|------------|------|--------|-----------|--------|
| Summer Sale Campaign | Email | Running | 100 | 85 | 45 | 3 | $500 |
| Product Launch - Email | Email | Scheduled | 100 | 0 | 0 | 0 | $500 |
| Q4 Outreach Campaign | Call | Completed | 100 | 50 | 25 | 2 | $500 |

### AI Requests (5 Samples)
| Query | Provider | Tokens | Time | Cost |
|-------|----------|--------|------|------|
| Sentiment Analysis | OpenAI | 450 | 280ms | $0.0018 |
| Email Response | GPT-4 | 320 | 195ms | $0.0032 |
| Lead Classification | OpenAI | 200 | 150ms | $0.0008 |
| Recommendations | Claude | 600 | 425ms | $0.0045 |
| Call Transcription | Whisper | 800 | 650ms | $0.0065 |

---

## 🎮 What to Try

### 1. Multi-Role Testing
```
1. Login as demo@vyomtech.com (Admin)
2. View all system features
3. Logout and login as agent@vyomtech.com
4. Notice the different dashboard for agents
5. Repeat with other roles
```

### 2. Lead Management
```
1. Go to Leads section
2. View the 8 sample leads
3. Try filtering by status (New, Contacted, Qualified, Converted)
4. Try sorting by priority
5. Open individual lead details
```

### 3. Call Tracking
```
1. Go to Calls section
2. See 4 sample call records
3. View call details including:
   - Duration and timestamps
   - AI provider and sentiment
   - Call transcription/notes
   - Agent information
```

### 4. Campaign Analytics
```
1. Go to Campaigns section
2. View 3 sample campaigns
3. Check campaign metrics:
   - Email open rates
   - Click-through rates
   - Conversions
   - Budget vs actual cost
4. View campaign recipient details
```

### 5. Agent Dashboard
```
1. Login as agent@vyomtech.com
2. View agent statistics:
   - 125 total calls handled
   - 8.5 min average handle time
   - 4.6★ satisfaction score
   - Current availability: Online
```

### 6. Multi-Tenant Isolation
```
1. Login as any user
2. All data shown belongs to "Demo Organization" tenant
3. Users from different tenants won't see each other's data
4. Verify data security and isolation
```

---

## 🔍 Testing Checklist

Use this checklist to verify all features are working:

### Authentication
- [ ] Login page loads with test credentials visible
- [ ] Can click any credential to auto-fill
- [ ] Login with admin account works
- [ ] Login with agent account works
- [ ] Login with manager account works
- [ ] Invalid credentials show error
- [ ] Logout works properly
- [ ] Session expires show redirect to login

### Dashboard
- [ ] Dashboard loads for each role
- [ ] Different views for different roles
- [ ] Widgets display sample data
- [ ] Charts render correctly
- [ ] Real-time metrics update

### Leads Module
- [ ] All 8 leads display
- [ ] Can filter by status
- [ ] Can filter by priority
- [ ] Can search by name/email
- [ ] Lead details page works
- [ ] Edit lead functionality works

### Calls Module
- [ ] All 4 sample calls display
- [ ] Call details show all information
- [ ] Can filter by direction (inbound/outbound)
- [ ] Can filter by status
- [ ] AI metrics display correctly
- [ ] Sentiment scores show

### Campaigns Module
- [ ] All 3 campaigns display
- [ ] Campaign details show metrics
- [ ] Recipients list displays
- [ ] Performance charts render
- [ ] Campaign status updates work

### Agent Profiles
- [ ] 2 agent profiles visible
- [ ] Agent stats display
- [ ] Skills list shows
- [ ] Availability status shows
- [ ] Call history accessible

### Tenant Management
- [ ] Tenant info displays
- [ ] Settings accessible
- [ ] Multi-tenant isolation verified
- [ ] No data leakage between tenants

---

## 🛠️ Troubleshooting

### Login Page Not Showing Credentials
**Solution**: 
- Hard refresh browser (Ctrl+Shift+R or Cmd+Shift+R)
- Clear browser cache
- Check if frontend is running on port 3000

### Can't Login Even with Correct Credentials
**Solution**:
- Verify backend is running on port 8080
- Check database connection
- Verify migrations have been applied
- Check browser console for errors

### Dashboard Shows No Data
**Solution**:
- Ensure migration 020_comprehensive_test_data.sql has been applied
- Check database for test data:
  ```sql
  SELECT COUNT(*) FROM `lead` WHERE tenant_id = 'demo-tenant';
  SELECT COUNT(*) FROM `call` WHERE tenant_id = 'demo-tenant';
  ```
- Refresh page

### Session Expired Immediately
**Solution**:
- Check JWT token expiration in config
- Verify backend is setting tokens correctly
- Check localStorage for auth_token
- Clear cookies and try again

---

## 📱 API Testing (Advanced)

### Test Login via cURL
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@vyomtech.com",
    "password": "DemoPass@123"
  }'
```

**Response**:
```json
{
  "token": "eyJ0eXAiOiJKV1QiLCJhbGc...",
  "user": {
    "id": 1,
    "email": "demo@vyomtech.com",
    "role": "admin",
    "tenant_id": "demo-tenant"
  }
}
```

### Test Protected Endpoint
```bash
# Using the token from login response
curl -X GET http://localhost:8080/api/v1/leads \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "X-Tenant-ID: demo-tenant"
```

### Get Agent Stats
```bash
curl -X GET http://localhost:8080/api/v1/agents/1/stats \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "X-Tenant-ID: demo-tenant"
```

---

## 📚 Documentation Links

- Backend API Docs: See `/docs/api.md`
- Database Schema: See `/docs/schema.md`
- Module Features: See individual module docs
- Deployment Guide: See `DEPLOYMENT_CHECKLIST.md`

---

## ✅ Success Indicators

You'll know everything is working when:
- ✅ Login page shows green credentials card
- ✅ Can login with any demo account
- ✅ Dashboard displays sample data
- ✅ All 8 leads are visible
- ✅ Can view call records and campaigns
- ✅ Agent stats display correctly
- ✅ Different roles see different features
- ✅ Multi-tenant isolation works

---

## 🎓 Learn More

### Understanding the Test Data
1. **Tenant**: Demo Organization (multi-tenant isolation)
2. **Users**: 5 demo accounts with different roles
3. **Leads**: Various stages of the sales pipeline
4. **Calls**: Real-world call scenarios
5. **Campaigns**: Email and call-based campaigns
6. **AI Integration**: Sample AI provider usage

### Next Steps After Testing
1. Customize for your use case
2. Add your own users and data
3. Configure AI providers
4. Set up integrations
5. Deploy to production

---

## 💬 Support

If you encounter any issues:
1. Check the logs in the terminal
2. Review browser console (F12)
3. Verify database connection
4. Run migrations again
5. Clear cache and restart

---

**Happy Testing! 🎉**

For more information, see:
- `BACKEND_DATABASE_COMPLETION.md` - Detailed completion summary
- `VERIFICATION_CHECKLIST.md` - Full verification details
- Individual module documentation in `/docs`

---

**Note**: Test credentials are for development/testing only. Remove before production deployment.
