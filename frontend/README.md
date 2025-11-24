# 🎨 Frontend Setup Guide

The frontend is built with **Next.js 15**, **React 19**, and **TypeScript**.

## Project Structure

```
frontend/
├── app/                          # Next.js App Router
│   ├── auth/
│   │   ├── login/
│   │   │   └── page.tsx         # Login page
│   │   └── register/
│   │       └── page.tsx         # Registration page
│   ├── dashboard/
│   │   ├── page.tsx             # Dashboard home
│   │   ├── agents/              # Agent management
│   │   ├── calls/               # Call management
│   │   ├── leads/               # Lead management
│   │   └── campaigns/           # Campaign management
│   ├── layout.tsx               # Root layout
│   └── globals.css              # Global styles
├── components/                  # Reusable React components
│   ├── auth/                    # Authentication forms
│   ├── dashboard/               # Dashboard components
│   ├── layouts/                 # Layout components
│   └── providers/               # Context providers
├── hooks/                       # Custom React hooks
│   └── useAuth.ts              # Authentication hook
├── services/                    # API services
│   └── api.ts                  # Axios client & API functions
├── types/                       # TypeScript type definitions
│   └── index.ts                # Type definitions
├── package.json
├── tsconfig.json
├── tailwind.config.js
└── next.config.js
```

## Installation

```bash
cd frontend

# Install dependencies
npm install

# Or with yarn
yarn install

# Or with pnpm
pnpm install
```

## Development

```bash
# Start dev server (http://localhost:3000)
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Run linting
npm run lint
```

## Environment Variables

Create `.env.local` in the `frontend` directory:

```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Key Features

✅ **Authentication**
- Login & Registration
- JWT token management
- Protected routes

✅ **Dashboard**
- Real-time statistics
- Agent management
- Call tracking
- Lead management
- Campaign management

✅ **Styling**
- Tailwind CSS for utility-first styling
- Responsive design
- Dark/Light mode ready

✅ **State Management**
- React Context for auth
- Zustand ready (dependency installed)
- React Query for data fetching

## API Integration

All API calls go through `services/api.ts`:

```typescript
import { authService, agentService } from '@/services/api'

// Login
await authService.login(email, password)

// Get agents
await agentService.listAgents()

// Update availability
await agentService.updateAvailability(agentId, 'online')
```

## Components

### Auth Components
- `LoginForm` - User login
- `RegisterForm` - User registration

### Layout Components
- `DashboardLayout` - Main dashboard wrapper with sidebar

### Dashboard Components
- `DashboardContent` - Dashboard home with stats

## Hooks

### useAuth
Access authentication state and functions:

```typescript
const { user, loading, login, logout } = useAuth()
```

## Types

All TypeScript types defined in `types/index.ts`:

```typescript
interface User {
  id: number
  email: string
  role: 'admin' | 'agent' | 'supervisor' | 'user'
  tenant_id: string
}

interface Agent extends User {
  status: 'active' | 'inactive'
  availability: 'online' | 'offline' | 'busy'
  skills: string[]
  // ...
}

// And many more...
```

## Next Steps

1. **Install dependencies**: `npm install`
2. **Ensure backend is running**: http://localhost:8080
3. **Start dev server**: `npm run dev`
4. **Visit**: http://localhost:3000

## Testing

```bash
# Run tests
npm test

# Run tests in watch mode
npm run test:watch
```

## Building for Production

```bash
# Build
npm run build

# Start
npm start
```

## Troubleshooting

### Backend connection fails
- Ensure Go backend is running on port 8080
- Check `NEXT_PUBLIC_API_URL` environment variable
- Verify CORS is enabled on backend

### Port 3000 already in use
```bash
npm run dev -- -p 3001
```

### Node modules issues
```bash
rm -rf node_modules package-lock.json
npm install
```

## More Information

- [Next.js Documentation](https://nextjs.org/docs)
- [React Documentation](https://react.dev)
- [TypeScript Documentation](https://www.typescriptlang.org/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
