'use client';

/**
 * SimpleSpreadsheet Component
 * Ultra-simple, Excel/Google Sheets-like interface
 * Minimal complexity, maximum usability
 */

import React, { useState, useCallback, useMemo } from 'react';
import {
  Plus, Trash2, Download, Upload, Copy, Undo, Redo, Settings,
  Eye, EyeOff, Filter, SortAsc, MoreVertical
} from 'lucide-react';

export interface SimpleColumn {
  id: string;
  label: string;
  type?: 'text' | 'number' | 'date' | 'currency' | 'percentage' | 'select';
  width?: number;
  editable?: boolean;
  hidden?: boolean;
  format?: string;
}

export interface SimpleSpreadsheetProps {
  columns: SimpleColumn[];
  data: any[];
  onDataChange: (data: any[]) => void;
  onAddRow?: () => void;
  onDeleteRow?: (index: number) => void;
  title?: string;
  showSearch?: boolean;
  densePacking?: boolean;
  allowExport?: boolean;
}

const SimpleSpreadsheet: React.FC<SimpleSpreadsheetProps> = ({
  columns,
  data,
  onDataChange,
  onAddRow,
  onDeleteRow,
  title,
  showSearch = true,
  densePacking = true,
  allowExport = true,
}) => {
  const [searchText, setSearchText] = useState('');
  const [editCell, setEditCell] = useState<{ row: number; col: string } | null>(null);
  const [history, setHistory] = useState<any[][]>([data]);
  const [historyIndex, setHistoryIndex] = useState(0);
  const [hiddenColumns, setHiddenColumns] = useState<Set<string>>(new Set());

  // Format cell value for display
  const formatValue = (value: any, type?: string, format?: string) => {
    if (value === null || value === undefined) return '';

    switch (type) {
      case 'currency':
        return new Intl.NumberFormat('en-IN', {
          style: 'currency',
          currency: 'INR',
          minimumFractionDigits: 2,
        }).format(Number(value));
      case 'percentage':
        return `${Number(value).toFixed(2)}%`;
      case 'date':
        return new Date(value).toLocaleDateString('en-IN');
      case 'number':
        return Number(value).toLocaleString('en-IN');
      default:
        return String(value);
    }
  };

  // Parse input value
  const parseValue = (value: string, type?: string) => {
    if (!value) return null;

    switch (type) {
      case 'number':
      case 'currency':
      case 'percentage':
        return parseFloat(value.replace(/[^\d.-]/g, ''));
      case 'date':
        return new Date(value).toISOString();
      default:
        return value;
    }
  };

  // Update cell
  const updateCell = (rowIndex: number, colId: string, value: any) => {
    const newData = data.map((row, idx) => 
      idx === rowIndex 
        ? { ...row, [colId]: value }
        : row
    );
    onDataChange(newData);
    addToHistory(newData);
  };

  // Update history
  const addToHistory = (newData: any[]) => {
    const newHistory = history.slice(0, historyIndex + 1);
    newHistory.push(JSON.parse(JSON.stringify(newData)));
    setHistory(newHistory);
    setHistoryIndex(newHistory.length - 1);
  };

  // Undo
  const undo = () => {
    if (historyIndex > 0) {
      const newIndex = historyIndex - 1;
      setHistoryIndex(newIndex);
      onDataChange(JSON.parse(JSON.stringify(history[newIndex])));
    }
  };

  // Redo
  const redo = () => {
    if (historyIndex < history.length - 1) {
      const newIndex = historyIndex + 1;
      setHistoryIndex(newIndex);
      onDataChange(JSON.parse(JSON.stringify(history[newIndex])));
    }
  };

  // Filter data
  const filteredData = useMemo(() => {
    if (!searchText) return data;
    const query = searchText.toLowerCase();
    return data.filter(row =>
      columns.some(col => 
        !hiddenColumns.has(col.id) &&
        String(row[col.id] || '').toLowerCase().includes(query)
      )
    );
  }, [data, searchText, hiddenColumns, columns]);

  // Visible columns
  const visibleColumns = useMemo(
    () => columns.filter(col => !hiddenColumns.has(col.id)),
    [columns, hiddenColumns]
  );

  // Export to CSV
  const handleExport = () => {
    const headers = visibleColumns.map(c => c.label).join(',');
    const rows = filteredData.map(row =>
      visibleColumns.map(col => {
        const val = row[col.id];
        return typeof val === 'string' && val.includes(',') ? `"${val}"` : val;
      }).join(',')
    );
    const csv = [headers, ...rows].join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `export-${Date.now()}.csv`;
    a.click();
  };

  // Toggle column visibility
  const toggleColumnVisibility = (colId: string) => {
    const newHidden = new Set(hiddenColumns);
    if (newHidden.has(colId)) {
      newHidden.delete(colId);
    } else {
      newHidden.add(colId);
    }
    setHiddenColumns(newHidden);
  };

  return (
    <div className="flex flex-col h-full bg-white rounded-lg shadow-sm border border-gray-200">
      {/* Header */}
      {(title || showSearch || allowExport) && (
        <div className="px-4 py-3 border-b border-gray-200 bg-gray-50">
          <div className="flex items-center justify-between mb-3">
            {title && <h3 className="text-lg font-semibold text-gray-900">{title}</h3>}
            <div className="flex gap-2">
              {/* Undo/Redo */}
              <button
                onClick={undo}
                disabled={historyIndex === 0}
                className="p-2 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed rounded"
                title="Undo"
              >
                <Undo size={16} />
              </button>
              <button
                onClick={redo}
                disabled={historyIndex === history.length - 1}
                className="p-2 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed rounded"
                title="Redo"
              >
                <Redo size={16} />
              </button>

              {/* Column visibility */}
              <div className="relative group">
                <button className="p-2 hover:bg-gray-200 rounded" title="Columns">
                  <Eye size={16} />
                </button>
                <div className="absolute right-0 top-full hidden group-hover:block bg-white border border-gray-200 rounded shadow-lg z-10 min-w-48 p-2">
                  {columns.map(col => (
                    <label key={col.id} className="flex items-center gap-2 px-3 py-2 hover:bg-gray-100 cursor-pointer rounded">
                      <input
                        type="checkbox"
                        checked={!hiddenColumns.has(col.id)}
                        onChange={() => toggleColumnVisibility(col.id)}
                        className="w-4 h-4"
                      />
                      <span className="text-sm">{col.label}</span>
                    </label>
                  ))}
                </div>
              </div>

              {/* Export */}
              {allowExport && (
                <button
                  onClick={handleExport}
                  className="p-2 hover:bg-gray-200 rounded"
                  title="Export CSV"
                >
                  <Download size={16} />
                </button>
              )}
            </div>
          </div>

          {/* Search */}
          {showSearch && (
            <input
              type="text"
              placeholder="🔍 Search all columns..."
              value={searchText}
              onChange={(e) => setSearchText(e.target.value)}
              className="w-full px-3 py-2 text-sm border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          )}
        </div>
      )}

      {/* Table */}
      <div className="overflow-auto flex-1">
        <table className="w-full border-collapse text-sm">
          {/* Header Row */}
          <thead>
            <tr className="bg-gray-100 border-b border-gray-300 sticky top-0 z-10">
              <th className="w-12 px-3 py-2 text-center text-gray-600 font-semibold text-xs bg-gray-200 border-r border-gray-300">
                #
              </th>
              {visibleColumns.map(col => (
                <th
                  key={col.id}
                  style={{ minWidth: col.width || 120 }}
                  className="px-3 py-2 text-left text-gray-700 font-semibold bg-gray-100 border-r border-gray-300 cursor-pointer hover:bg-gray-200 transition-colors"
                >
                  {col.label}
                </th>
              ))}
              {(onDeleteRow || onAddRow) && (
                <th className="w-12 px-3 py-2 text-center bg-gray-200 border-l border-gray-300"></th>
              )}
            </tr>
          </thead>

          {/* Data Rows */}
          <tbody>
            {filteredData.length === 0 ? (
              <tr>
                <td
                  colSpan={visibleColumns.length + 2}
                  className="px-4 py-8 text-center text-gray-400"
                >
                  No data • Click + to add rows
                </td>
              </tr>
            ) : (
              filteredData.map((row, idx) => (
                <tr
                  key={idx}
                  className={`border-b border-gray-200 hover:bg-blue-50 transition-colors ${
                    idx % 2 === 0 ? 'bg-white' : 'bg-gray-50'
                  }`}
                >
                  {/* Row number */}
                  <td className="w-12 px-3 py-2 text-center text-gray-500 bg-gray-50 border-r border-gray-300 text-xs font-mono select-none">
                    {idx + 1}
                  </td>

                  {/* Data cells */}
                  {visibleColumns.map(col => (
                    <td
                      key={col.id}
                      style={{ minWidth: col.width || 120 }}
                      className="px-3 py-2 border-r border-gray-200 cursor-cell"
                      onClick={() => col.editable && setEditCell({ row: idx, col: col.id })}
                    >
                      {editCell?.row === idx && editCell?.col === col.id ? (
                        <input
                          type={col.type === 'number' || col.type === 'currency' ? 'number' : 'text'}
                          autoFocus
                          defaultValue={row[col.id] || ''}
                          onBlur={(e) => {
                            updateCell(idx, col.id, parseValue(e.target.value, col.type));
                            setEditCell(null);
                          }}
                          onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                              updateCell(idx, col.id, parseValue((e.target as HTMLInputElement).value, col.type));
                              setEditCell(null);
                            } else if (e.key === 'Escape') {
                              setEditCell(null);
                            }
                          }}
                          className="w-full px-2 py-1 border border-blue-500 rounded focus:outline-none focus:ring-2 focus:ring-blue-600"
                        />
                      ) : (
                        <span className={col.editable ? 'hover:bg-blue-100 px-1 rounded' : ''}>
                          {formatValue(row[col.id], col.type, col.format)}
                        </span>
                      )}
                    </td>
                  ))}

                  {/* Actions */}
                  {(onDeleteRow || onAddRow) && (
                    <td className="w-12 px-2 py-2 text-center border-l border-gray-300">
                      {onDeleteRow && (
                        <button
                          onClick={() => onDeleteRow(data.indexOf(row))}
                          className="p-1 text-red-500 hover:bg-red-100 rounded transition-colors"
                        >
                          <Trash2 size={14} />
                        </button>
                      )}
                    </td>
                  )}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Footer */}
      <div className="px-4 py-2 bg-gray-50 border-t border-gray-200 text-xs text-gray-600 flex items-center justify-between">
        <div>
          Showing {filteredData.length} of {data.length} rows
          {searchText && ` (filtered from ${data.length})`}
        </div>
        {onAddRow && (
          <button
            onClick={onAddRow}
            className="flex items-center gap-2 px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
          >
            <Plus size={16} />
            Add Row
          </button>
        )}
      </div>
    </div>
  );
};

export default SimpleSpreadsheet;
