'use client';

/**
 * Accounting Dashboard
 * Excel-like GL transaction and account tracking
 */

import React, { useState } from 'react';
import SimpleSpreadsheet, { SimpleColumn } from '@/components/SimpleSpreadsheet';
import ExcelDashboard, { DashboardTab } from '@/components/ExcelDashboard';
import { CreditCard, BarChart3 } from 'lucide-react';

const AccountingDashboard = () => {
  // Sample GL transaction data
  const [transactions, setTransactions] = useState([
    {
      id: '1',
      date: '2024-01-15',
      refNumber: 'JE-001',
      description: 'Sales invoice recorded',
      debitAccount: '1101 - Cash',
      debitAmount: 53100,
      creditAccount: '4001 - Sales Revenue',
      creditAmount: 45000,
      taxAmount: 8100,
      posted: true,
    },
    {
      id: '2',
      date: '2024-01-16',
      refNumber: 'JE-002',
      description: 'Payment received',
      debitAccount: '1101 - Cash',
      debitAmount: 37760,
      creditAccount: '1201 - Accounts Receivable',
      creditAmount: 37760,
      taxAmount: 0,
      posted: true,
    },
    {
      id: '3',
      date: '2024-01-17',
      refNumber: 'JE-003',
      description: 'Expense recorded',
      debitAccount: '5101 - Rent Expense',
      debitAmount: 50000,
      creditAccount: '1101 - Cash',
      creditAmount: 50000,
      taxAmount: 0,
      posted: false,
    },
  ]);

  // Sample chart of accounts
  const [accounts, setAccounts] = useState([
    {
      id: '1',
      accountCode: '1101',
      accountName: 'Cash',
      accountType: 'ASSET',
      balance: 485860,
      status: 'ACTIVE',
    },
    {
      id: '2',
      accountCode: '1201',
      accountName: 'Accounts Receivable',
      accountType: 'ASSET',
      balance: 88500,
      status: 'ACTIVE',
    },
    {
      id: '3',
      accountCode: '2101',
      accountName: 'Accounts Payable',
      accountType: 'LIABILITY',
      balance: -125000,
      status: 'ACTIVE',
    },
    {
      id: '4',
      accountCode: '3101',
      accountName: 'Capital',
      accountType: 'EQUITY',
      balance: 1000000,
      status: 'ACTIVE',
    },
    {
      id: '5',
      accountCode: '4001',
      accountName: 'Sales Revenue',
      accountType: 'REVENUE',
      balance: 45000,
      status: 'ACTIVE',
    },
    {
      id: '6',
      accountCode: '5101',
      accountName: 'Rent Expense',
      accountType: 'EXPENSE',
      balance: -50000,
      status: 'ACTIVE',
    },
  ]);

  // GL transaction columns
  const transactionColumns: SimpleColumn[] = [
    { id: 'date', label: 'Date', type: 'date', width: 110, editable: true },
    { id: 'refNumber', label: 'Ref #', type: 'text', width: 100, editable: false },
    { id: 'description', label: 'Description', type: 'text', width: 180, editable: true },
    { id: 'debitAccount', label: 'Debit Account', type: 'text', width: 160, editable: true },
    { id: 'debitAmount', label: 'Debit Amount', type: 'currency', width: 140, editable: true },
    { id: 'creditAccount', label: 'Credit Account', type: 'text', width: 160, editable: true },
    { id: 'creditAmount', label: 'Credit Amount', type: 'currency', width: 140, editable: true },
    { id: 'posted', label: 'Posted', type: 'select', width: 100, editable: true },
  ];

  // Chart of accounts columns
  const accountColumns: SimpleColumn[] = [
    { id: 'accountCode', label: 'Code', type: 'text', width: 100, editable: false },
    { id: 'accountName', label: 'Account Name', type: 'text', width: 160, editable: true },
    { id: 'accountType', label: 'Type', type: 'select', width: 120, editable: true },
    { id: 'balance', label: 'Balance', type: 'currency', width: 140, editable: false },
    { id: 'status', label: 'Status', type: 'select', width: 110, editable: true },
  ];

  // Calculate trial balance
  const debits = transactions.filter(t => t.posted).reduce((sum, t) => sum + t.debitAmount, 0);
  const credits = transactions.filter(t => t.posted).reduce((sum, t) => sum + t.creditAmount, 0);
  const isBalanced = Math.abs(debits - credits) < 0.01;

  // Get account totals
  const totalAssets = accounts.filter(a => a.accountType === 'ASSET').reduce((sum, a) => sum + a.balance, 0);
  const totalLiabilities = accounts.filter(a => a.accountType === 'LIABILITY').reduce((sum, a) => sum + a.balance, 0);
  const totalEquity = accounts.filter(a => a.accountType === 'EQUITY').reduce((sum, a) => sum + a.balance, 0);

  // Dashboard tabs
  const tabs: DashboardTab[] = [
    {
      id: 'journal',
      label: 'Journal Entries',
      icon: <CreditCard size={20} />,
      content: (
        <div>
          <div className="mb-4">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">General Ledger Transactions</h3>
            <div className="grid grid-cols-4 gap-4 mb-4">
              <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                <p className="text-sm text-blue-600 font-semibold">Total Entries</p>
                <p className="text-2xl font-bold text-blue-900">{transactions.length}</p>
              </div>
              <div className={`p-4 rounded-lg border ${isBalanced ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'}`}>
                <p className={`text-sm font-semibold ${isBalanced ? 'text-green-600' : 'text-red-600'}`}>
                  {isBalanced ? 'Trial Balance' : 'Out of Balance'}
                </p>
                <p className={`text-2xl font-bold ${isBalanced ? 'text-green-900' : 'text-red-900'}`}>
                  {isBalanced ? '✓ Balanced' : '✗ Imbalanced'}
                </p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
                <p className="text-sm text-purple-600 font-semibold">Total Debits</p>
                <p className="text-2xl font-bold text-purple-900">₹ {debits.toLocaleString('en-IN', { maximumFractionDigits: 0 })}</p>
              </div>
              <div className="bg-orange-50 p-4 rounded-lg border border-orange-200">
                <p className="text-sm text-orange-600 font-semibold">Total Credits</p>
                <p className="text-2xl font-bold text-orange-900">₹ {credits.toLocaleString('en-IN', { maximumFractionDigits: 0 })}</p>
              </div>
            </div>
          </div>
          <SimpleSpreadsheet
            columns={transactionColumns}
            data={transactions}
            onDataChange={setTransactions}
            onDeleteRow={(idx) => setTransactions(transactions.filter((_, i) => i !== idx))}
            onAddRow={() => setTransactions([...transactions, {
              id: `je-${Date.now()}`,
              date: new Date().toISOString().split('T')[0],
              refNumber: `JE-${String(transactions.length + 1).padStart(3, '0')}`,
              description: '',
              debitAccount: '',
              debitAmount: 0,
              creditAccount: '',
              creditAmount: 0,
              taxAmount: 0,
              posted: false,
            }])}
            title="Journal Entry Register"
            showSearch={true}
            allowExport={true}
          />
        </div>
      ),
    },
    {
      id: 'accounts',
      label: 'Chart of Accounts',
      icon: <BarChart3 size={20} />,
      content: (
        <div>
          <div className="mb-4">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">Chart of Accounts</h3>
            <div className="grid grid-cols-4 gap-4 mb-4">
              <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                <p className="text-sm text-blue-600 font-semibold">Total Assets</p>
                <p className="text-2xl font-bold text-blue-900">₹ {totalAssets.toLocaleString('en-IN', { maximumFractionDigits: 0 })}</p>
              </div>
              <div className="bg-red-50 p-4 rounded-lg border border-red-200">
                <p className="text-sm text-red-600 font-semibold">Total Liabilities</p>
                <p className="text-2xl font-bold text-red-900">₹ {Math.abs(totalLiabilities).toLocaleString('en-IN', { maximumFractionDigits: 0 })}</p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
                <p className="text-sm text-purple-600 font-semibold">Total Equity</p>
                <p className="text-2xl font-bold text-purple-900">₹ {totalEquity.toLocaleString('en-IN', { maximumFractionDigits: 0 })}</p>
              </div>
              <div className={`p-4 rounded-lg border ${Math.abs(totalAssets - (Math.abs(totalLiabilities) + totalEquity)) < 0.01 ? 'bg-green-50 border-green-200' : 'bg-yellow-50 border-yellow-200'}`}>
                <p className={`text-sm font-semibold ${Math.abs(totalAssets - (Math.abs(totalLiabilities) + totalEquity)) < 0.01 ? 'text-green-600' : 'text-yellow-600'}`}>
                  Accounting Equation
                </p>
                <p className={`text-2xl font-bold ${Math.abs(totalAssets - (Math.abs(totalLiabilities) + totalEquity)) < 0.01 ? 'text-green-900' : 'text-yellow-900'}`}>
                  {Math.abs(totalAssets - (Math.abs(totalLiabilities) + totalEquity)) < 0.01 ? '✓ Balanced' : '⚠ Check'}
                </p>
              </div>
            </div>
          </div>
          <SimpleSpreadsheet
            columns={accountColumns}
            data={accounts}
            onDataChange={setAccounts}
            onDeleteRow={(idx) => setAccounts(accounts.filter((_, a) => a !== accounts[idx]))}
            onAddRow={() => setAccounts([...accounts, {
              id: `acc-${Date.now()}`,
              accountCode: '',
              accountName: '',
              accountType: 'ASSET',
              balance: 0,
              status: 'ACTIVE',
            }])}
            title="Accounts Register"
            showSearch={true}
            allowExport={true}
          />
        </div>
      ),
    },
  ];

  return (
    <div className="h-screen flex flex-col">
      <ExcelDashboard
        tabs={tabs}
        title="Accounting Dashboard"
        subtitle="GL transactions and chart of accounts • Simple Excel-like interface"
      />
    </div>
  );
};

export default AccountingDashboard;
