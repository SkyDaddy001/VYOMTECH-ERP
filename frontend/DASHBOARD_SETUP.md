# VYOMTECH ERP - Frontend Dashboard

A modern, responsive ERP dashboard built with Next.js 16, React 19, and Tailwind CSS for the VYOMTECH multi-tenant AI Call Center system.

## Features

### 📊 Dashboard
- **Stats Overview**: Real-time KPIs including leads, calls, conversion rates, and revenue
- **Recent Leads**: Quick view of the latest leads with scoring
- **Active Campaigns**: Campaign performance metrics and ROI tracking
- **Agent Performance**: Top agents by rating and call volume

### 📋 Leads Management
- Complete lead database view
- Lead scoring system
- Multi-status tracking (New, Qualified, Contacted, Converted, Lost)
- Bulk actions and filtering

### 🎯 Campaigns
- Campaign creation and management
- Budget tracking
- ROI calculation
- Lead attribution
- Campaign status monitoring

### 👥 Agents
- Agent performance dashboard
- Call statistics
- Quality ratings
- Status monitoring

### 🔧 Technical Stack
- **Framework**: Next.js 16.0.7 with Turbopack
- **UI Library**: React 19.2.0
- **Styling**: Tailwind CSS v4
- **State Management**: Zustand (optional)
- **HTTP Client**: Axios
- **Icons**: React Icons
- **Utilities**: date-fns, clsx

## Project Structure

```
frontend/
├── app/                           # Next.js App Router
│   ├── dashboard/                # Dashboard pages
│   │   ├── page.tsx             # Main dashboard
│   │   ├── leads/               # Leads management
│   │   ├── campaigns/           # Campaign management
│   │   ├── calls/               # Call management
│   │   ├── agents/              # Agent management
│   │   ├── reports/             # Reports
│   │   └── settings/            # Settings
│   ├── layout.tsx               # Root layout
│   ├── page.tsx                 # Home page (redirects to dashboard)
│   └── globals.css              # Global styles
├── components/
│   ├── dashboard/               # Dashboard components
│   │   ├── stats-overview.tsx   # KPI cards
│   │   ├── recent-leads.tsx     # Recent leads table
│   │   ├── campaigns-overview.tsx
│   │   └── agents-performance.tsx
│   └── layout/                  # Layout components
│       ├── sidebar.tsx          # Navigation sidebar
│       └── header.tsx           # Top header
├── hooks/                       # React hooks
│   └── use-dashboard.ts         # Dashboard data fetching hooks
├── lib/
│   ├── api-client.ts           # Axios API client wrapper
│   └── utils.ts                # Utility functions
├── services/
│   └── dashboard.ts            # Dashboard API service
└── types/                       # TypeScript type definitions
```

## Setup Instructions

### 1. Install Dependencies

```bash
cd frontend
npm install
```

This will install all required packages:
- `axios` - HTTP client for API calls
- `react-icons` - Icon library
- `date-fns` - Date utilities
- `zustand` - Optional state management
- `chart.js` & `react-chartjs-2` - Charts (for future analytics)

### 2. Environment Configuration

Create a `.env.local` file in the frontend directory:

```env
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080

# Optional: For production
# NEXT_PUBLIC_API_URL=https://api.yourdomain.com
```

### 3. Running the Development Server

```bash
npm run dev
```

The application will be available at `http://localhost:3000`

### 4. Building for Production

```bash
npm run build
npm start
```

## API Integration

The dashboard integrates with the VYOMTECH backend API endpoints:

### Dashboard Stats
- `GET /api/v1/dashboard/stats` - Overall statistics

### Leads
- `GET /api/v1/leads` - List all leads
- `GET /api/v1/leads/:id` - Get specific lead
- `POST /api/v1/leads` - Create new lead
- `PUT /api/v1/leads/:id` - Update lead

### Campaigns
- `GET /api/v1/campaigns` - List campaigns
- `GET /api/v1/campaigns/:id` - Get campaign details
- `POST /api/v1/campaigns` - Create campaign

### Calls
- `GET /api/v1/calls` - List calls
- `GET /api/v1/calls/:id` - Get call details

### Agents
- `GET /api/v1/agents` - List all agents
- `GET /api/v1/agents/:id` - Get agent details

### Gamification
- `GET /api/v1/gamification/:userId` - Get user gamification stats

## Authentication

Authentication tokens are stored in `localStorage` with the key `auth_token`. The API client automatically includes the token in request headers.

### Login Flow
1. User provides credentials
2. Backend returns JWT token
3. Token stored in localStorage
4. Token automatically added to all API requests
5. On 401 response, user is redirected to login

## Component Usage Examples

### Using Dashboard Hooks

```typescript
import { useDashboardStats, useLeads } from '@/hooks/use-dashboard'

export const MyComponent = () => {
  const { stats, loading, error } = useDashboardStats()
  const { leads } = useLeads({ limit: 10 })

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error: {error}</div>

  return (
    <div>
      <p>Total Leads: {stats?.totalLeads}</p>
      {leads.map(lead => <div key={lead.id}>{lead.name}</div>)}
    </div>
  )
}
```

### Using Dashboard Service Directly

```typescript
import dashboardService from '@/services/dashboard'

const data = await dashboardService.getDashboardStats()
const leads = await dashboardService.getLeads({ status: 'qualified' })
```

## Styling Guide

### Color Palette
- **Primary**: `#3b82f6` (Blue)
- **Secondary**: `#10b981` (Green)
- **Danger**: `#ef4444` (Red)
- **Warning**: `#f59e0b` (Amber)

### Utility Classes
- `.transition-smooth` - Smooth transitions
- `.truncate-2` / `.truncate-3` - Multi-line text truncation
- `.shadow-sm-custom` / `.shadow-md-custom` / `.shadow-lg-custom` - Custom shadows

## Responsive Design

The dashboard is fully responsive with breakpoints:
- Mobile: < 640px
- Tablet: 768px - 1024px
- Desktop: > 1024px

The sidebar collapses on mobile with a hamburger menu.

## Performance Optimization

- **Code Splitting**: Automatic with Next.js and dynamic imports
- **Image Optimization**: Using Next.js Image component
- **API Caching**: Implement with SWR or React Query for production
- **Bundle Size**: Monitored with Build Analysis

## Development Commands

```bash
# Development server
npm run dev

# Build production bundle
npm run build

# Start production server
npm start

# Lint code
npm run lint

# Format code
npm run format  # (if prettier is added)
```

## Troubleshooting

### API Connection Issues
1. Verify `NEXT_PUBLIC_API_URL` in `.env.local`
2. Check backend is running on configured port
3. Verify CORS headers are set correctly
4. Check browser Network tab for blocked requests

### Build Errors
1. Clear `.next` directory: `rm -rf .next`
2. Reinstall dependencies: `rm -rf node_modules && npm install`
3. Check TypeScript errors: `npx tsc --noEmit`

### Styling Issues
1. Verify Tailwind CSS is properly imported in `globals.css`
2. Clear Tailwind cache: `rm -rf .next`
3. Check for CSS conflicts with custom styles

## Future Enhancements

- [ ] Analytics dashboards with charts
- [ ] Real-time notifications with WebSocket
- [ ] Advanced filtering and search
- [ ] Export reports to PDF/Excel
- [ ] Mobile app optimization
- [ ] Dark mode support
- [ ] Internationalization (i18n)
- [ ] Accessibility improvements (a11y)

## Contributing

When adding new pages or components:

1. Follow the existing folder structure
2. Use TypeScript for type safety
3. Implement loading and error states
4. Use the API client wrapper for consistency
5. Add comments for complex logic
6. Test responsive design

## Support

For issues or questions:
1. Check existing issue tracker
2. Review API documentation
3. Contact development team

## License

© 2025 VYOMTECH. All rights reserved.
