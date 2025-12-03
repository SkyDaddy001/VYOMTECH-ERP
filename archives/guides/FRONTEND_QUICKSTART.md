# Frontend Development - Quick Start Guide

## Overview
The VYOMTECH frontend has been completely redesigned with a **spreadsheet-first interface** that makes it intuitive for users familiar with Excel. No learning curve required!

---

## What Was Built

### 1. **Core Component: SpreadsheetGrid**
A fully-featured data grid component that mimics Excel with:
- ✅ Inline cell editing
- ✅ Column sorting (click headers)
- ✅ Row filtering (type in filter row)
- ✅ Add/Delete rows
- ✅ Row numbering
- ✅ Multiple data types (text, number, date, select, checkbox)
- ✅ Auto-calculated fields
- ✅ Responsive layout

### 2. **Sample Dashboards**
Ready-to-use pages for all major modules:
- **Projects** (`/dashboard/projects`) - Manage all projects
- **Sites** (`/dashboard/sites`) - Track construction sites
- **BOQ** (`/dashboard/boq`) - Bill of Quantities with auto-calc
- **Progress** (`/dashboard/progress`) - Daily progress tracking

### 3. **Reusable Components**
- `SpreadsheetGrid.tsx` - Main data grid
- `SpreadsheetToolbar.tsx` - Search/Export/Import toolbar
- `DashboardLayout.tsx` - Sidebar navigation

---

## File Structure

```
frontend/
├── components/
│   ├── SpreadsheetGrid.tsx       ⭐ Main component
│   ├── SpreadsheetToolbar.tsx    📋 Toolbar
│   └── DashboardLayout.tsx       📐 Layout
├── app/dashboard/
│   ├── projects.tsx              📊 Projects page
│   ├── sites.tsx                 🏗️ Sites page
│   ├── boq.tsx                   📋 BOQ page
│   └── progress.tsx              📈 Progress page
└── docs/
    └── FRONTEND_DESIGN.md         📖 Full documentation
```

---

## Getting Started

### Install Dependencies
```bash
cd frontend
npm install
```

### Run Development Server
```bash
npm run dev
```

Visit `http://localhost:3000` in your browser.

### Build for Production
```bash
npm run build
npm start
```

---

## Component Usage

### Basic Example: Create a Data Grid

```typescript
import SpreadsheetGrid, { Column } from '@/components/SpreadsheetGrid';

function MyPage() {
  const [data, setData] = useState([
    { id: '1', name: 'Item 1', value: 100 },
    { id: '2', name: 'Item 2', value: 200 },
  ]);

  const columns: Column[] = [
    {
      id: 'name',
      header: 'Name',
      accessor: 'name',
      type: 'text',
      width: 200,
      editable: true,
      sortable: true,
    },
    {
      id: 'value',
      header: 'Value',
      accessor: 'value',
      type: 'number',
      width: 120,
      editable: true,
      sortable: true,
    },
  ];

  return (
    <SpreadsheetGrid
      title="My Items"
      columns={columns}
      data={data}
      onDataChange={setData}
      onAddRow={() => setData([...data, { id: '', name: '', value: 0 }])}
      onDeleteRow={(idx) => setData(data.filter((_, i) => i !== idx))}
      showRowNumbers={true}
    />
  );
}
```

---

## Key Features Explained

### 1. **Inline Editing**
- Click any cell to edit
- Blue highlight indicates editable cells
- Press Enter to save, Escape to cancel
- Tab key moves to next cell

### 2. **Sorting**
- Click column header to sort ascending
- Click again to sort descending
- Sort indicator (chevron) shows active sort
- Hover over header to see sort button

### 3. **Filtering**
- Type in the filter row (below headers)
- Matches any part of the cell value
- Real-time filtering
- Multiple filters work together

### 4. **Data Types**
- **text**: Free-form text input
- **number**: Numeric input with validation
- **date**: Date picker (future enhancement)
- **select**: Dropdown with predefined options
- **checkbox**: Boolean toggle

---

## Adding a New Dashboard Page

### Step 1: Create Page File
Create `app/dashboard/mynewpage.tsx`:

```typescript
'use client';
import React, { useState } from 'react';
import SpreadsheetGrid, { Column } from '@/components/SpreadsheetGrid';

const MyNewPage = () => {
  const [data, setData] = useState([
    // Your initial data
  ]);

  const columns: Column[] = [
    // Your columns
  ];

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">My New Module</h1>
      <SpreadsheetGrid
        columns={columns}
        data={data}
        onDataChange={setData}
      />
    </div>
  );
};

export default MyNewPage;
```

### Step 2: Update Navigation
Add to `DashboardLayout.tsx`:

```typescript
const dashboardMenuItems = [
  // ... existing items
  { label: 'My Module', href: '/dashboard/mynewpage', icon: '🆕' },
];
```

That's it! Your new page is live.

---

## Styling Guide

### Colors
- **Primary (Blue)**: `text-blue-600`, `bg-blue-500`
- **Success (Green)**: `text-green-600`, `bg-green-500`
- **Danger (Red)**: `text-red-600`, `bg-red-500`
- **Gray**: `text-gray-600`, `bg-gray-100`

### Common Patterns
```typescript
// Card layout
<div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">

// Button
<button className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">

// Header
<h1 className="text-3xl font-bold text-gray-900">

// Stats card
<div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
  <div className="text-sm text-gray-600">Label</div>
  <div className="text-2xl font-bold text-gray-900">Value</div>
</div>
```

---

## API Integration

### Connect to Backend
Update `services/api.ts`:

```typescript
import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const projectsAPI = {
  getAll: () => axios.get(`${API_URL}/api/v1/projects`),
  create: (data) => axios.post(`${API_URL}/api/v1/projects`, data),
  update: (id, data) => axios.put(`${API_URL}/api/v1/projects/${id}`, data),
  delete: (id) => axios.delete(`${API_URL}/api/v1/projects/${id}`),
};
```

### Use in Component
```typescript
import { projectsAPI } from '@/services/api';
import { useQuery, useMutation } from '@tanstack/react-query';

const { data, isLoading } = useQuery({
  queryKey: ['projects'],
  queryFn: () => projectsAPI.getAll(),
});

const mutation = useMutation({
  mutationFn: (newData) => projectsAPI.update(id, newData),
  onSuccess: () => {
    // Refetch data
  },
});
```

---

## Testing

### Run Tests
```bash
npm test
npm run test:watch
```

### Example Test
```typescript
import { render, screen } from '@testing-library/react';
import SpreadsheetGrid from '@/components/SpreadsheetGrid';

describe('SpreadsheetGrid', () => {
  it('renders with data', () => {
    render(<SpreadsheetGrid columns={[]} data={[]} onDataChange={() => {}} />);
    expect(screen.getByRole('table')).toBeInTheDocument();
  });
});
```

---

## Deployment

### Build
```bash
npm run build
```

### Docker
The `Dockerfile` in the root will automatically build the frontend.

### Docker Compose
```bash
docker-compose up frontend
```

---

## Performance Tips

1. **Use useMemo** for expensive calculations
2. **Use useCallback** for handlers
3. **Lazy load heavy components** with React.lazy
4. **Optimize images** before using
5. **Virtual scrolling** for large datasets (future)

---

## Common Issues

### Issue: Grid not updating
**Solution**: Use `onDataChange` to update state properly
```typescript
const handleDataChange = useCallback((updatedData) => {
  setData(updatedData); // Always update state
}, []);
```

### Issue: Styles not applied
**Solution**: Ensure Tailwind CSS is imported
```typescript
// In app/globals.css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

### Issue: API errors
**Solution**: Check NEXT_PUBLIC_API_URL environment variable
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev
```

---

## Next Steps

1. **Connect to Backend API** - Update pages to use real data
2. **Add More Pages** - Safety, Compliance, Equipment, Permits
3. **Enhance Features** - Export/Import, Bulk operations, Charts
4. **Mobile Optimization** - Responsive design improvements
5. **Authentication** - Integrate with backend auth

---

## Useful Resources

- [Next.js Docs](https://nextjs.org/docs)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [React Hooks](https://react.dev/reference/react)
- [TypeScript](https://www.typescriptlang.org/docs)

---

## Support

For questions or issues:
1. Check the `FRONTEND_DESIGN.md` for detailed documentation
2. Review component props in the `.tsx` files
3. Check browser console for errors
4. Run `npm run lint` to find code issues

---

**Happy Coding! 🎉**

Last Updated: December 3, 2025
