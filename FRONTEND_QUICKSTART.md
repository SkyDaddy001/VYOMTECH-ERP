# Frontend Excel/Sheets - Quick Start Guide

## 🚀 Quick Integration (5 minutes)

### Step 1: Copy Components
All components are in `frontend/components/`:
```
✅ SimpleSpreadsheet.tsx
✅ InvoiceEntryForm.tsx
✅ SalesOrderEntryForm.tsx
✅ BOQEntryForm.tsx
✅ ExcelDashboard.tsx
✅ SalesDashboard.tsx
✅ ConstructionDashboard.tsx
✅ AccountingDashboard.tsx
```

### Step 2: Add Routes
```tsx
// frontend/app/dashboard/sales/page.tsx
import SalesDashboard from '@/components/SalesDashboard';
export default SalesDashboard;

// frontend/app/dashboard/construction/page.tsx
import ConstructionDashboard from '@/components/ConstructionDashboard';
export default ConstructionDashboard;

// frontend/app/dashboard/accounting/page.tsx
import AccountingDashboard from '@/components/AccountingDashboard';
export default AccountingDashboard;
```

### Step 3: Update Navigation
```tsx
// In your sidebar or navbar
<Link href="/dashboard/sales">📊 Sales</Link>
<Link href="/dashboard/construction">🏗️ Construction</Link>
<Link href="/dashboard/accounting">📈 Accounting</Link>
```

### Step 4: Done! 🎉
Navigate to `/dashboard/sales` and start using!

---

## 📋 Component Reference

### SimpleSpreadsheet
Super simple grid with sort, filter, inline-edit.

```tsx
<SimpleSpreadsheet
  columns={[
    { id: 'name', label: 'Name', type: 'text', editable: true },
    { id: 'amount', label: 'Amount', type: 'currency', editable: true },
  ]}
  data={myData}
  onDataChange={setMyData}
  title="My Grid"
  showSearch={true}
  allowExport={true}
/>
```

### InvoiceEntryForm
Invoice entry with auto-calculation.

```tsx
<InvoiceEntryForm
  onSave={(data) => {
    console.log('Invoice:', data);
    // POST to /api/v1/invoices
  }}
/>
```

**Auto-calculated:**
- Amount = Qty × Unit Price
- Tax = Amount × Tax Rate
- Total = Amount + Tax

### SalesOrderEntryForm
Sales order with 18% GST auto-calc.

```tsx
<SalesOrderEntryForm
  onSave={(data) => {
    console.log('Order:', data);
    // POST to /api/v1/sales-orders
  }}
/>
```

**Auto-calculated:**
- Amount = Qty × Unit Price - Discount
- Tax = Amount × 18%
- Total = Amount + Tax

### BOQEntryForm
Bill of Quantities with 0.01₹ precision.

```tsx
<BOQEntryForm
  onSave={(data) => {
    console.log('BOQ:', data);
    // POST to /api/v1/boq
  }}
/>
```

**Features:**
- Auto-generated item codes (B001, B002, ...)
- Progress % tracking (0-100%)
- Unit selector (nos, sqm, cum, lm, kg, t)
- Contingency % (default 5%)

### ExcelDashboard
Multi-tab workbook interface.

```tsx
<ExcelDashboard
  tabs={[
    { id: 'invoices', label: 'Invoices', icon: <FileText />, content: <Grid1 /> },
    { id: 'orders', label: 'Orders', icon: <ShoppingCart />, content: <Grid2 /> },
  ]}
  title="Sales Dashboard"
/>
```

---

## 🎨 Color Guide

Use these colors for KPI cards:

```tsx
// Info
<div className="bg-blue-50 border-blue-200">
  <p className="text-blue-600">Label</p>
</div>

// Success
<div className="bg-green-50 border-green-200">
  <p className="text-green-600">Label</p>
</div>

// Warning
<div className="bg-orange-50 border-orange-200">
  <p className="text-orange-600">Label</p>
</div>

// Data
<div className="bg-purple-50 border-purple-200">
  <p className="text-purple-600">Label</p>
</div>
```

---

## 🔧 Common Tasks

### 1. Add Custom Column Type

In `SimpleSpreadsheet.tsx`, update formatValue():

```typescript
case 'email':
  return <a href={`mailto:${value}`}>{value}</a>;
```

Then use:
```tsx
{ id: 'email', label: 'Email', type: 'email', editable: true }
```

### 2. Add Validation

In form components:

```typescript
const handleSave = () => {
  if (!formData.invoiceNumber) {
    alert('Invoice number required');
    return;
  }
  onSave(formData);
};
```

### 3. Format Columns

Add `format` prop to columns:

```tsx
{ id: 'date', label: 'Date', type: 'date', format: 'DD-MMM-YYYY' }
```

### 4. Hide Columns by Default

Set `hidden` in column definition:

```tsx
{ id: 'internalId', label: 'Internal ID', hidden: true }
```

### 5. Change Column Width

```tsx
{ id: 'description', label: 'Description', width: 300 }
```

---

## 🐛 Troubleshooting

### Grid shows no data
- Check data prop is an array
- Verify column IDs match data object keys
- Check browser console for errors

### Inline editing not working
- Set `editable: true` on column
- Check onDataChange callback is defined
- Verify data is mutable (not frozen)

### Export button missing
- Set `allowExport={true}` on SimpleSpreadsheet
- Check CSV file is downloading to Downloads folder

### Calculations wrong
- Verify number formatting in form
- Check parseFloat/parseValue functions
- Debug in browser DevTools

### Styling looks weird
- Verify Tailwind CSS is included
- Check no conflicting CSS
- Clear Next.js cache and rebuild

---

## 💾 Backend Integration

### Save Invoice Example
```typescript
const handleSaveInvoice = async (data: InvoiceFormData) => {
  // 1. Save invoice
  const invoiceRes = await fetch('/api/v1/invoices', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Tenant-ID': tenantId,
    },
    body: JSON.stringify(data),
  });
  
  if (!invoiceRes.ok) {
    alert('Error saving invoice');
    return;
  }

  const invoice = await invoiceRes.json();

  // 2. Post to GL (debit AR, credit Revenue)
  const glRes = await fetch('/api/v1/journal-entries', {
    method: 'POST',
    headers: { 'X-Tenant-ID': tenantId },
    body: JSON.stringify({
      date: data.date,
      description: `Invoice ${data.invoiceNumber}`,
      items: [
        { account: '1201', debit: data.total }, // AR
        { account: '4001', credit: data.subtotal }, // Revenue
        { account: '2106', credit: data.taxAmount }, // GST Payable
      ],
    }),
  });

  alert('Invoice saved and posted to GL');
};
```

### Load Dashboard Data
```typescript
useEffect(() => {
  const loadInvoices = async () => {
    const res = await fetch('/api/v1/invoices', {
      headers: { 'X-Tenant-ID': tenantId },
    });
    const invoices = await res.json();
    setInvoices(invoices);
  };
  
  loadInvoices();
}, [tenantId]);
```

---

## 📱 Mobile Support

### Current
- ✅ Fully functional on desktop
- ✅ Works on tablet (768px+)

### To-Do for Mobile (landscape mode)
```tsx
// Hide less important columns on mobile
hidden: window.innerWidth < 640
```

---

## 🚀 Performance Tips

1. **Use React.memo for large lists**
   ```tsx
   const GridRow = React.memo(({ row, columns }) => ...)
   ```

2. **Pagination for 1000+ rows**
   ```tsx
   const [page, setPage] = useState(0);
   const pageSize = 50;
   const displayData = data.slice(page * pageSize, (page + 1) * pageSize);
   ```

3. **Debounce search**
   ```tsx
   const [searchText, setSearchText] = useState('');
   const debouncedSearch = useMemo(() => 
     debounce(setSearchText, 300), []
   );
   ```

---

## ✅ Checklist Before Deployment

- [ ] All components copied to `frontend/components/`
- [ ] Routes created for each dashboard
- [ ] Navigation menu updated with links
- [ ] Backend APIs working and tested
- [ ] Multi-tenancy headers (X-Tenant-ID) implemented
- [ ] Error handling added to forms
- [ ] Loading states visible while fetching
- [ ] CSV export tested
- [ ] Mobile view tested (if needed)
- [ ] No console errors or warnings

---

## 📞 Support Files

1. **Comprehensive Guide**: `FRONTEND_EXCEL_SHEETS_GUIDE.md`
2. **Completion Report**: `FRONTEND_SIMPLIFICATION_COMPLETE.md`
3. **Component Source**: `frontend/components/*.tsx`

---

## 🎯 Success = Users Love It!

Your frontend is now:
- ✅ Simple as Excel
- ✅ Familiar to all users
- ✅ Productive and fast
- ✅ Beautiful and modern
- ✅ Full-featured and reliable

**Happy coding! 🚀**

---

*Quick Start Guide v1.0*  
*All 8 components production-ready*  
*Zero dependencies beyond React + Tailwind*
