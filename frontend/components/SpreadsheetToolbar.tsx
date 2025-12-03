/**
 * Toolbar Component
 * Common spreadsheet operations: search, filter, export, import
 */

import React from 'react';
import { Search, Download, Upload, Settings } from 'lucide-react';

interface ToolbarProps {
  onSearch?: (query: string) => void;
  onExport?: () => void;
  onImport?: () => void;
  onSettings?: () => void;
  title?: string;
}

const SpreadsheetToolbar: React.FC<ToolbarProps> = ({
  onSearch,
  onExport,
  onImport,
  onSettings,
  title,
}) => {
  return (
    <div className="bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
      {title && <h2 className="text-lg font-semibold text-gray-800">{title}</h2>}
      
      <div className="flex items-center gap-2">
        {onSearch && (
          <div className="relative">
            <Search size={18} className="absolute left-2 top-2.5 text-gray-400" />
            <input
              type="text"
              placeholder="Search..."
              onChange={e => onSearch(e.target.value)}
              className="pl-8 pr-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        )}

        {onExport && (
          <button
            onClick={onExport}
            className="px-3 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 flex items-center gap-1 text-sm"
          >
            <Download size={16} /> Export
          </button>
        )}

        {onImport && (
          <button
            onClick={onImport}
            className="px-3 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 flex items-center gap-1 text-sm"
          >
            <Upload size={16} /> Import
          </button>
        )}

        {onSettings && (
          <button
            onClick={onSettings}
            className="px-3 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 text-sm"
          >
            <Settings size={16} />
          </button>
        )}
      </div>
    </div>
  );
};

export default SpreadsheetToolbar;
