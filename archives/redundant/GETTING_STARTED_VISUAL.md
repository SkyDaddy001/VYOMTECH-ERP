# 🎬 Getting Started - Visual Guide

## Your Application Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    User's Browser                            │
│                  (http://localhost:3000)                     │
├─────────────────────────────────────────────────────────────┤
│                   React/Next.js Frontend                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Pages:                                             │   │
│  │  - /auth/login (public)                             │   │
│  │  - /auth/register (public)                          │   │
│  │  - /dashboard (protected)                           │   │
│  │  - /dashboard/agents                                │   │
│  │  - /dashboard/calls, leads, campaigns, reports      │   │
│  │                                                      │   │
│  │  Components: Forms, Layout, Cards                   │   │
│  │  Services: API Client with JWT                      │   │
│  │  State: Auth Context + React Query                  │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTP Requests
                            │ (with JWT Token)
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              Go REST API Server                              │
│            (http://localhost:8080)                           │
├─────────────────────────────────────────────────────────────┤
│  Routes:                    Services:                        │
│  ├─ /api/v1/auth/*         ├─ AuthService                   │
│  ├─ /api/v1/agents/*       ├─ AgentService                  │
│  ├─ /api/v1/leads/*        ├─ EmailService                  │
│  ├─ /api/v1/calls/*        ├─ PasswordResetService          │
│  ├─ /api/v1/campaigns/*    └─ AIOrchestrator                │
│  └─ /health                                                  │
│                                                              │
│  Middleware:               Database:                         │
│  ├─ JWT Auth               ├─ Tenants                       │
│  ├─ CORS                   ├─ Users                         │
│  ├─ Logging                ├─ Agents                        │
│  └─ Error Recovery         ├─ Leads                         │
│                            ├─ Calls                         │
│                            ├─ Campaigns                     │
│                            └─ Password Resets               │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ SQL Queries
                            ▼
                    ┌──────────────────┐
                    │  MySQL Database  │
                    │   (Port 3306)    │
                    │  Database: call  │
                    │  center (Podman) │
                    └──────────────────┘
```

---

## 🚀 Installation Steps (5 minutes)

### Step 1️⃣: Navigate to Frontend
```bash
cd frontend
```
**Time**: 5 seconds

### Step 2️⃣: Install Dependencies
```bash
npm install
```
**Time**: 2-3 minutes  
**What it does**: Downloads ~1000+ packages from npm registry

### Step 3️⃣: Create .env.local
```bash
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local
```
**Time**: 2 seconds

### Step 4️⃣: Start Development Server
```bash
npm run dev
```
**Time**: 10 seconds to start

### Step 5️⃣: Open in Browser
```
http://localhost:3000
```
**Time**: 2 seconds

---

## 📊 Before/After: What Changed

### BEFORE (Only Backend)
```
Project Structure:
├── cmd/main.go              ← Backend entry point
├── internal/                ← Backend services
├── pkg/                     ← Backend packages
├── migrations/              ← Database schema
├── frontend/                ← EMPTY (no UI)
└── ...

Workflow:
1. Start backend: ./bin/main
2. Test with curl: curl http://localhost:8080/health
3. No visual interface - API only
```

### AFTER (Full Stack)
```
Project Structure:
├── cmd/main.go              ← Backend entry point
├── internal/                ← Backend services
├── pkg/                     ← Backend packages
├── migrations/              ← Database schema
├── frontend/                ← ✨ NEW: React/Next.js
│   ├── app/                 ← Pages & routing
│   ├── components/          ← Reusable components
│   ├── services/            ← API client
│   ├── hooks/               ← Custom hooks
│   ├── types/               ← TypeScript types
│   ├── package.json         ← Dependencies
│   └── ...
└── ...

Workflow:
1. Terminal 1: ./bin/main (backend)
2. Terminal 2: cd frontend && npm run dev (frontend)
3. Open http://localhost:3000 (beautiful UI!)
4. Login and see dashboard
```

---

## 🎯 User Journey

### New User Flow
```
1. Visit http://localhost:3000
   ↓
2. See Login Page
   ↓
3. Click "Register here"
   ↓
4. Fill in:
   - Name: John Doe
   - Email: john@example.com
   - Password: MySecurePass123!
   ↓
5. Click "Sign Up"
   ↓
6. Redirected to Login
   ↓
7. Login with credentials
   ↓
8. See Dashboard
   ↓
9. Navigate sidebar:
   - Dashboard (home)
   - Agents
   - Calls
   - Leads
   - Campaigns
   - Reports
   ↓
10. Click Logout
   ↓
11. Back to Login
```

### Existing User Flow
```
1. Visit http://localhost:3000
   ↓
2. See Login Page (session expired)
   ↓
3. Login with email/password
   ↓
4. See Dashboard (token restored)
   ↓
5. Dashboard shows:
   - Online Agents: 12
   - Calls Today: 342
   - Avg Handle Time: 5m 32s
   - Customer Satisfaction: 94%
   - Revenue Today: $12,450
   - Queue Length: 8
   ↓
6. Quick actions:
   - ➕ New Agent
   - 📞 Start Call
   - 📋 New Campaign
   - 📊 View Reports
```

---

## 📁 File Organization Quick Reference

### Frontend Important Files
```
frontend/
├── app/
│   ├── layout.tsx          ← Root layout wrapper
│   ├── page.tsx            ← Redirect to dashboard
│   ├── globals.css         ← Global styles
│   └── auth/
│       ├── login/
│       │   └── page.tsx    ← LOGIN PAGE
│       └── register/
│           └── page.tsx    ← REGISTER PAGE
│
├── components/
│   ├── auth/
│   │   ├── LoginForm.tsx   ← Login form component
│   │   └── RegisterForm.tsx ← Register form component
│   ├── layouts/
│   │   └── DashboardLayout.tsx ← SIDEBAR + HEADER
│   └── providers/
│       ├── AuthProvider.tsx ← Auth context
│       └── ToasterProvider.tsx ← Toast notifications
│
├── services/
│   └── api.ts              ← ALL API CALLS HERE
│
├── hooks/
│   └── useAuth.ts          ← useAuth() hook
│
└── types/
    └── index.ts            ← All TypeScript types
```

---

## 🔑 Key Code Locations

### To Fix Login: `frontend/components/auth/LoginForm.tsx`
### To Add Features: `frontend/components/dashboard/`
### To Change API URL: `frontend/.env.local`
### To Modify Styles: `frontend/tailwind.config.js`
### To Call Backend API: `frontend/services/api.ts`

---

## 🧪 Quick Testing

### Test 1: Login Works
```
1. npm run dev
2. Open http://localhost:3000
3. Click login button
4. Should show form ✓
```

### Test 2: Registration Works
```
1. Click "Register here" link
2. Fill form
3. Should navigate to login ✓
```

### Test 3: API Connection Works
```
1. Open DevTools (F12)
2. Go to Network tab
3. Try to login
4. Should see POST /api/v1/auth/login ✓
```

### Test 4: Dashboard Shows
```
1. Login successfully
2. Should see stats cards ✓
3. Should see agents list ✓
4. Should see sidebar ✓
```

---

## ⚡ Performance Stats

| Item | Time |
|------|------|
| npm install | 2-3 min |
| npm run dev startup | 5-10 sec |
| Page load | <1 sec |
| Login request | 100-200 ms |
| Dashboard render | <500 ms |
| Build for production | 30-45 sec |

---

## 🛠️ Troubleshooting Quick Guide

### Problem: "Cannot find module"
```bash
rm -rf node_modules
npm install
```

### Problem: Port 3000 in use
```bash
npm run dev -- -p 3001
```

### Problem: Backend unreachable
```bash
# In terminal: ./bin/main
# Check: curl http://localhost:8080/health
```

### Problem: Login fails silently
```bash
# Open DevTools (F12)
# Go to Network tab
# Try login
# Check response in Network tab
```

---

## 📚 Documentation Quick Links

| Document | Purpose |
|----------|---------|
| **FRONTEND_SETUP.md** | Installation guide |
| **frontend/README.md** | Detailed frontend info |
| **FULLSTACK_COMPLETE.md** | Full architecture |
| **APPLICATION_RUNNING.md** | Backend status |

---

## ✅ Checklist: Are You Ready?

- [ ] Node.js 18+ installed
- [ ] npm 9+ installed
- [ ] Backend running (./bin/main)
- [ ] MySQL running (podman ps shows mysql-callcenter)
- [ ] In frontend folder
- [ ] Ran npm install
- [ ] Created .env.local
- [ ] Ran npm run dev
- [ ] Browser opened to http://localhost:3000

**All checked? You're ready to go! 🚀**

---

## 🎨 What You'll See

### Login Page
```
┌─────────────────────────────────┐
│    Call Center Login            │
├─────────────────────────────────┤
│                                 │
│  📧 Email: [____________]       │
│  🔑 Password: [____________]    │
│                                 │
│  [Sign In Button]               │
│                                 │
│  Don't have account? Register   │
└─────────────────────────────────┘
```

### Dashboard Page
```
┌────────┬──────────────────────────────────┐
│ MENU   │  📊 Dashboard                    │
│ 📊Dash │                                  │
│ 👥Agt  ├──────────────────────────────────┤
│ 📞Call │ Welcome to Call Center!          │
│ 📋Lead │                                  │
│ 🎯Cam  │ [Online] [Calls] [Time] [Satis] │
│ 📈Rep  │  Agents  Today   Avg    Score   │
│ 🚪Out  │   12     342    5m 32s  94%     │
│        │                                  │
│        │ Recent Calls | Active Agents     │
└────────┴──────────────────────────────────┘
```

---

## 🚀 Ready to Start?

```bash
# Copy this command:
cd frontend && npm install && npm run dev

# Then open:
http://localhost:3000

# Login with:
Email: testuser@example.com
Password: TestPassword123!
```

**That's it! Enjoy your new web dashboard! 🎉**
