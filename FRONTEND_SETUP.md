# 🚀 Frontend Installation & Quick Start

## Prerequisites

- Node.js 18+ and npm 9+
- Go backend running on port 8080
- MySQL database running

## Quick Setup

### 1. Navigate to Frontend Directory
```bash
cd frontend
```

### 2. Install Dependencies
```bash
npm install
```

### 3. Create Environment File
Create `.env.local`:
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 4. Start Development Server
```bash
npm run dev
```

Visit: **http://localhost:3000**

---

## Default Test Credentials

Use these to test the application:

```
Email: testuser@example.com
Password: TestPassword123!
Role: user
Tenant: default-tenant
```

Or register a new account at the registration page.

---

## Features Implemented

✅ **Authentication**
- User login & registration
- JWT token management
- Protected dashboard routes

✅ **Dashboard**
- Real-time statistics cards
- Active agents list
- Recent calls history
- Quick action buttons

✅ **Agent Management**
- List all agents
- View agent details (stub)
- Update agent status (stub)

✅ **Responsive Design**
- Mobile-friendly UI
- Collapsible sidebar
- Adaptive grid layouts

✅ **Error Handling**
- Toast notifications
- Form validation
- API error messages

---

## Project Structure

```
frontend/
├── app/                    # Next.js pages
│   ├── auth/
│   │   ├── login/
│   │   └── register/
│   ├── dashboard/
│   │   ├── agents/
│   │   ├── calls/
│   │   ├── leads/
│   │   ├── campaigns/
│   │   └── reports/
│   └── layout.tsx
├── components/
│   ├── auth/              # LoginForm, RegisterForm
│   ├── dashboard/         # DashboardContent
│   ├── layouts/           # DashboardLayout
│   └── providers/         # AuthProvider, ToasterProvider
├── hooks/
│   └── useAuth.ts        # useAuth hook
├── services/
│   └── api.ts            # API client
├── types/
│   └── index.ts          # TypeScript interfaces
└── package.json
```

---

## Available Scripts

```bash
# Development
npm run dev                # Start dev server on port 3000

# Production
npm run build             # Build for production
npm start                 # Start production server

# Code Quality
npm run lint              # Run ESLint
npm test                  # Run tests

# Testing
npm run test:watch        # Run tests in watch mode
```

---

## API Integration

### Authentication Service
```typescript
import { authService } from '@/services/api'

// Login
const response = await authService.login(email, password)
const token = response.token
const user = response.user

// Register
await authService.register(email, password, role, tenant_id)

// Logout
authService.logout()
```

### Agent Service
```typescript
import { agentService } from '@/services/api'

// Get all agents
const agents = await agentService.listAgents()

// Get agent details
const agent = await agentService.getAgent(agentId)

// Update availability
await agentService.updateAvailability(agentId, 'online')

// Get agent stats
const stats = await agentService.getAgentStats(agentId)
```

---

## Component Architecture

### Layout Hierarchy
```
RootLayout
├── AuthProvider
├── ToasterProvider
└── Content
    ├── DashboardLayout (for protected pages)
    │   ├── Sidebar Navigation
    │   ├── Top Header
    │   └── Main Content Area
    └── Auth Pages
        ├── LoginForm
        └── RegisterForm
```

### State Management
- **Auth State**: React Context (AuthProvider)
- **Notifications**: react-hot-toast
- **Future**: Zustand (installed, ready to use)

---

## Styling

**Framework**: Tailwind CSS

### Key Classes
```css
/* Layout */
flex, grid, w-full, h-screen, p-6, m-4

/* Colors */
bg-blue-600, text-gray-800, border-gray-300

/* Utilities */
hover:, focus:, disabled:, transition

/* Responsive */
md:, lg:, xl: breakpoints
```

### Custom Utilities
```css
.truncate-2    /* Truncate to 2 lines */
.animate-pulse /* Pulsing animation */
```

---

## Authentication Flow

```
1. User visits http://localhost:3000
2. Redirected to /auth/login (no token)
3. User logs in or registers
4. Token stored in localStorage
5. Redirected to /dashboard
6. Protected routes check token
7. If expired, redirected to login
```

---

## Troubleshooting

### Backend Connection Error
```
Error: connect ECONNREFUSED 127.0.0.1:8080
```
**Solution**: Ensure Go backend is running
```bash
# In another terminal, from project root
./bin/main
```

### Port 3000 Already in Use
```bash
# Use different port
npm run dev -- -p 3001
```

### Node Modules Issues
```bash
rm -rf node_modules package-lock.json
npm install
```

### Clear Cache
```bash
rm -rf .next
npm run dev
```

### Token Expired
- Frontend will redirect to login
- Clear localStorage: `localStorage.clear()`
- Login again

---

## Testing the Full Stack

### 1. Start Backend (Terminal 1)
```bash
cd /c/Users/Skydaddy/Desktop/Developement
./startup.sh start
# Or just: ./bin/main
```

### 2. Start Frontend (Terminal 2)
```bash
cd /c/Users/Skydaddy/Desktop/Developement/frontend
npm install  # First time only
npm run dev
```

### 3. Test in Browser
1. Open http://localhost:3000
2. Click "Register here"
3. Create new account (or use testuser@example.com)
4. Login
5. View dashboard with real stats

---

## Next Steps

### Quick Wins (Easy)
- [ ] Add dark mode toggle
- [ ] Add profile page
- [ ] Add settings page
- [ ] Add password change
- [ ] Add logout confirmation

### Medium Difficulty
- [ ] Complete Agent CRUD
- [ ] Add Call management pages
- [ ] Add Lead management pages
- [ ] Add Campaign pages
- [ ] Add Charts & Analytics

### Advanced
- [ ] Real-time updates (Socket.io)
- [ ] Call recording playback
- [ ] Advanced filtering & search
- [ ] Data export (CSV/PDF)
- [ ] User preferences storage

---

## Useful Resources

- **[Next.js Docs](https://nextjs.org/docs)**
- **[React Docs](https://react.dev)**
- **[TypeScript Docs](https://www.typescriptlang.org/docs)**
- **[Tailwind CSS](https://tailwindcss.com/docs)**
- **[Axios Docs](https://axios-http.com/docs)**

---

## Development Tips

### VS Code Extensions (Recommended)
- ES7+ React/Redux/React-Native snippets
- Tailwind CSS IntelliSense
- TypeScript Vue Plugin
- Prettier - Code formatter

### Debugging
```typescript
// Use console.log
console.log('Debug:', data)

// Use React DevTools in browser
// Use Next.js Dev Tools at top-left of page
```

### Performance
- Images: Next.js Image component
- Code splitting: Automatic by Next.js
- Lazy loading: dynamic imports
- Caching: React Query with cache time

---

**Status**: ✅ Frontend ready to develop  
**Version**: 1.0.0-alpha  
**Last Updated**: 2025-11-21
