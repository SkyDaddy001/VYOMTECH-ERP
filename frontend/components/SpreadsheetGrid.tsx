/**
 * SpreadsheetGrid Component
 * A spreadsheet-like data grid with inline editing, sorting, and filtering
 * Designed for users familiar with Excel
 */

import React, { useState, useCallback, useMemo } from 'react';
import { ChevronUp, ChevronDown, Plus, Trash2, Download, Upload } from 'lucide-react';

export interface Column {
  id: string;
  header: string;
  accessor: string;
  type?: 'text' | 'number' | 'date' | 'select' | 'checkbox';
  width?: number;
  editable?: boolean;
  sortable?: boolean;
  filterOptions?: { label: string; value: string }[];
}

interface SpreadsheetGridProps {
  columns: Column[];
  data: any[];
  onDataChange: (updatedData: any[]) => void;
  onAddRow?: () => void;
  onDeleteRow?: (rowIndex: number) => void;
  title?: string;
  showRowNumbers?: boolean;
  densePacking?: boolean;
}

const SpreadsheetGrid: React.FC<SpreadsheetGridProps> = ({
  columns,
  data,
  onDataChange,
  onAddRow,
  onDeleteRow,
  title,
  showRowNumbers = true,
  densePacking = true,
}) => {
  const [editingCell, setEditingCell] = useState<{ rowIndex: number; colId: string } | null>(null);
  const [sortConfig, setSortConfig] = useState<{ key: string; direction: 'asc' | 'desc' } | null>(null);
  const [filterValues, setFilterValues] = useState<Record<string, string>>({});

  // Filter data
  const filteredData = useMemo(() => {
    return data.filter(row =>
      Object.entries(filterValues).every(([colId, filterValue]) => {
        if (!filterValue) return true;
        const column = columns.find(c => c.id === colId);
        if (!column) return true;
        const cellValue = String(row[column.accessor] || '').toLowerCase();
        return cellValue.includes(filterValue.toLowerCase());
      })
    );
  }, [data, filterValues, columns]);

  // Sort data
  const sortedData = useMemo(() => {
    if (!sortConfig) return filteredData;
    
    const sorted = [...filteredData].sort((a, b) => {
      const column = columns.find(c => c.id === sortConfig.key);
      if (!column) return 0;
      
      const aVal = a[column.accessor];
      const bVal = b[column.accessor];
      
      if (aVal < bVal) return sortConfig.direction === 'asc' ? -1 : 1;
      if (aVal > bVal) return sortConfig.direction === 'asc' ? 1 : -1;
      return 0;
    });
    
    return sorted;
  }, [filteredData, sortConfig, columns]);

  const handleCellChange = (rowIndex: number, colId: string, value: any) => {
    const column = columns.find(c => c.id === colId);
    if (!column) return;

    const newData = [...data];
    const realRowIndex = data.findIndex(row => row === sortedData[rowIndex]);
    newData[realRowIndex] = {
      ...newData[realRowIndex],
      [column.accessor]: value,
    };

    onDataChange(newData);
    setEditingCell(null);
  };

  const handleSort = (colId: string) => {
    setSortConfig(prev =>
      prev?.key === colId && prev.direction === 'asc'
        ? { key: colId, direction: 'desc' }
        : { key: colId, direction: 'asc' }
    );
  };

  const handleFilterChange = (colId: string, value: string) => {
    setFilterValues(prev => ({
      ...prev,
      [colId]: value,
    }));
  };

  const renderCell = (row: any, column: Column, rowIndex: number) => {
    const isEditing = editingCell?.rowIndex === rowIndex && editingCell?.colId === column.id;
    const value = row[column.accessor] || '';

    if (isEditing && column.editable) {
      return (
        <input
          autoFocus
          type={column.type === 'number' ? 'number' : 'text'}
          value={value}
          onChange={e => handleCellChange(rowIndex, column.id, e.target.value)}
          onBlur={() => setEditingCell(null)}
          onKeyDown={e => {
            if (e.key === 'Enter') handleCellChange(rowIndex, column.id, e.currentTarget.value);
            if (e.key === 'Escape') setEditingCell(null);
          }}
          className="w-full px-2 py-1 border border-blue-400 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      );
    }

    return (
      <div
        onClick={() => column.editable && setEditingCell({ rowIndex, colId: column.id })}
        className={`px-3 py-2 cursor-cell ${column.editable ? 'hover:bg-blue-50' : ''}`}
      >
        {column.type === 'checkbox' ? (
          <input
            type="checkbox"
            checked={Boolean(value)}
            onChange={e => handleCellChange(rowIndex, column.id, e.target.checked)}
            className="w-4 h-4"
          />
        ) : (
          String(value)
        )}
      </div>
    );
  };

  return (
    <div className="w-full h-full flex flex-col bg-white rounded-lg shadow-sm border border-gray-200">
      {/* Header */}
      {title && (
        <div className="px-4 py-3 border-b border-gray-200 flex justify-between items-center bg-gray-50">
          <h2 className="text-lg font-semibold text-gray-800">{title}</h2>
          <div className="flex gap-2">
            {onAddRow && (
              <button
                onClick={onAddRow}
                className="px-3 py-1 bg-green-500 text-white rounded hover:bg-green-600 flex items-center gap-1 text-sm"
              >
                <Plus size={16} /> Add Row
              </button>
            )}
            <button className="px-3 py-1 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 text-sm">
              <Download size={16} />
            </button>
          </div>
        </div>
      )}

      {/* Toolbar */}
      <div className="px-4 py-2 bg-gray-50 border-b border-gray-200 flex gap-2 items-center text-xs text-gray-600">
        <span>{sortedData.length} rows</span>
        {Object.keys(filterValues).length > 0 && (
          <span className="ml-auto text-blue-600">
            {Object.keys(filterValues).length} filter(s) active
          </span>
        )}
      </div>

      {/* Table */}
      <div className="flex-1 overflow-auto">
        <table className="w-full border-collapse text-sm">
          <thead>
            <tr className="bg-gray-50 border-b border-gray-300 sticky top-0 z-10">
              {showRowNumbers && (
                <th className="w-12 px-3 py-2 text-center bg-gray-100 border-r border-gray-300 font-medium text-gray-700">
                  #
                </th>
              )}
              {columns.map(column => (
                <th
                  key={column.id}
                  style={{ width: column.width ? `${column.width}px` : 'auto', minWidth: '100px' }}
                  className="px-3 py-2 text-left bg-gray-100 border-r border-gray-300 font-semibold text-gray-700 group"
                >
                  <div className="flex items-center justify-between gap-2">
                    <span>{column.header}</span>
                    {column.sortable && (
                      <button
                        onClick={() => handleSort(column.id)}
                        className="p-0.5 opacity-0 group-hover:opacity-100 hover:bg-gray-200 rounded"
                      >
                        {sortConfig?.key === column.id ? (
                          sortConfig.direction === 'asc' ? (
                            <ChevronUp size={14} />
                          ) : (
                            <ChevronDown size={14} />
                          )
                        ) : (
                          <ChevronUp size={14} className="opacity-50" />
                        )}
                      </button>
                    )}
                  </div>
                </th>
              ))}
              {onDeleteRow && <th className="w-10 bg-gray-100 border-l border-gray-300" />}
            </tr>

            {/* Filter Row */}
            <tr className="bg-white border-b border-gray-200">
              {showRowNumbers && <td className="w-12 bg-gray-50" />}
              {columns.map(column => (
                <td key={column.id} className="px-2 py-1 border-r border-gray-300">
                  <input
                    type="text"
                    placeholder="Filter..."
                    value={filterValues[column.id] || ''}
                    onChange={e => handleFilterChange(column.id, e.target.value)}
                    className="w-full px-2 py-1 text-xs border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-blue-400"
                  />
                </td>
              ))}
              {onDeleteRow && <td className="w-10" />}
            </tr>
          </thead>

          <tbody>
            {sortedData.length === 0 ? (
              <tr>
                <td colSpan={columns.length + (showRowNumbers ? 1 : 0) + (onDeleteRow ? 1 : 0)} className="px-4 py-8 text-center text-gray-500">
                  No data
                </td>
              </tr>
            ) : (
              sortedData.map((row, displayIndex) => {
                const realIndex = data.findIndex(r => r === row);
                return (
                  <tr key={displayIndex} className={`border-b border-gray-200 hover:bg-blue-50 ${displayIndex % 2 === 0 ? 'bg-white' : 'bg-gray-50'}`}>
                    {showRowNumbers && (
                      <td className="w-12 px-3 py-1 text-center text-gray-500 bg-gray-50 border-r border-gray-300 font-mono text-xs select-none">
                        {displayIndex + 1}
                      </td>
                    )}
                    {columns.map(column => (
                      <td key={column.id} style={{ width: column.width ? `${column.width}px` : 'auto', minWidth: '100px' }} className="px-0 py-0 border-r border-gray-300">
                        {renderCell(row, column, displayIndex)}
                      </td>
                    ))}
                    {onDeleteRow && (
                      <td className="w-10 px-2 py-1 text-center border-l border-gray-300">
                        <button
                          onClick={() => onDeleteRow(realIndex)}
                          className="p-1 text-red-500 hover:bg-red-100 rounded"
                        >
                          <Trash2 size={14} />
                        </button>
                      </td>
                    )}
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>

      {/* Footer */}
      <div className="px-4 py-2 bg-gray-50 border-t border-gray-200 text-xs text-gray-600">
        Showing {sortedData.length} of {data.length} rows
      </div>
    </div>
  );
};

export default SpreadsheetGrid;
