/**
 * Style Guide & Component Library
 * Reference for all UI patterns used in VYOMTECH ERP
 */

'use client';

import React from 'react';
import { Plus, Trash2, Download, Settings, Search } from 'lucide-react';

const StyleGuide = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 p-8">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-4xl font-bold text-gray-900 mb-2">VYOMTECH UI Style Guide</h1>
        <p className="text-lg text-gray-600 mb-12">
          Spreadsheet-inspired design system for intuitive data management
        </p>

        {/* Color Palette */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Color Palette</h2>
          <div className="grid grid-cols-2 gap-6 lg:grid-cols-4">
            {[
              { name: 'Blue Primary', class: 'bg-blue-500', hex: '#3B82F6' },
              { name: 'Green Success', class: 'bg-green-500', hex: '#22C55E' },
              { name: 'Red Danger', class: 'bg-red-500', hex: '#EF4444' },
              { name: 'Gray Neutral', class: 'bg-gray-500', hex: '#6B7280' },
              { name: 'Blue Light', class: 'bg-blue-50', hex: '#EFF6FF' },
              { name: 'Green Light', class: 'bg-green-50', hex: '#F0FDF4' },
              { name: 'Red Light', class: 'bg-red-50', hex: '#FEF2F2' },
              { name: 'Gray Light', class: 'bg-gray-50', hex: '#F9FAFB' },
            ].map((color) => (
              <div key={color.name} className="text-center">
                <div className={`${color.class} h-32 rounded-lg shadow-md mb-3`} />
                <p className="font-semibold text-gray-800">{color.name}</p>
                <p className="text-sm text-gray-600">{color.hex}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Buttons */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Buttons</h2>
          <div className="space-y-6">
            {/* Primary Buttons */}
            <div>
              <h3 className="text-lg font-semibold text-gray-800 mb-3">Primary Actions</h3>
              <div className="flex flex-wrap gap-4">
                <button className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition">
                  Primary Button
                </button>
                <button className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition">
                  <Plus size={16} className="inline mr-2" /> Add Row
                </button>
                <button className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 transition">
                  Delete
                </button>
              </div>
            </div>

            {/* Secondary Buttons */}
            <div>
              <h3 className="text-lg font-semibold text-gray-800 mb-3">Secondary Actions</h3>
              <div className="flex flex-wrap gap-4">
                <button className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition">
                  Secondary
                </button>
                <button className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition flex items-center gap-2">
                  <Download size={16} /> Export
                </button>
                <button className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition flex items-center gap-2">
                  <Settings size={16} /> Settings
                </button>
              </div>
            </div>

            {/* Icon Buttons */}
            <div>
              <h3 className="text-lg font-semibold text-gray-800 mb-3">Icon Buttons</h3>
              <div className="flex flex-wrap gap-4">
                <button className="p-2 text-blue-500 hover:bg-blue-100 rounded transition">
                  <Search size={20} />
                </button>
                <button className="p-2 text-green-500 hover:bg-green-100 rounded transition">
                  <Plus size={20} />
                </button>
                <button className="p-2 text-red-500 hover:bg-red-100 rounded transition">
                  <Trash2 size={20} />
                </button>
              </div>
            </div>
          </div>
        </section>

        {/* Input Fields */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Input Fields</h2>
          <div className="space-y-4 max-w-md">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">Text Input</label>
              <input
                type="text"
                placeholder="Enter text..."
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">Number Input</label>
              <input
                type="number"
                placeholder="0.00"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">Search Input</label>
              <div className="relative">
                <Search size={18} className="absolute left-3 top-2.5 text-gray-400" />
                <input
                  type="text"
                  placeholder="Search..."
                  className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">Select</label>
              <select className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                <option>Option 1</option>
                <option>Option 2</option>
                <option>Option 3</option>
              </select>
            </div>
            <div className="flex items-center gap-2">
              <input type="checkbox" id="checkbox" className="w-4 h-4" />
              <label htmlFor="checkbox" className="text-sm text-gray-700">
                Checkbox option
              </label>
            </div>
          </div>
        </section>

        {/* Cards */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Card Components</h2>
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
            {/* Stats Card */}
            <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
              <div className="text-sm text-gray-600">Total Projects</div>
              <div className="text-2xl font-bold text-gray-900 mt-2">28</div>
              <div className="text-xs text-gray-500 mt-2">+5 this month</div>
            </div>

            {/* Status Card */}
            <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
              <div className="flex items-center justify-between mb-3">
                <h3 className="font-semibold text-gray-800">Active</h3>
                <span className="text-2xl font-bold text-green-600">15</span>
              </div>
              <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                <div className="h-full w-3/4 bg-green-500 rounded-full" />
              </div>
            </div>

            {/* Feature Card */}
            <div className="bg-gradient-to-br from-blue-50 to-green-50 p-4 rounded-lg shadow-sm border border-blue-200">
              <h3 className="font-semibold text-gray-800 mb-2">Key Feature</h3>
              <p className="text-sm text-gray-600">
                Spreadsheet-style interface familiar to Excel users
              </p>
            </div>
          </div>
        </section>

        {/* Table Styling */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Table / Grid Example</h2>
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
            <table className="w-full">
              <thead>
                <tr className="bg-gray-100 border-b border-gray-300">
                  <th className="px-4 py-3 text-left font-semibold text-gray-700">#</th>
                  <th className="px-4 py-3 text-left font-semibold text-gray-700">Name</th>
                  <th className="px-4 py-3 text-left font-semibold text-gray-700">Status</th>
                  <th className="px-4 py-3 text-left font-semibold text-gray-700">Progress</th>
                </tr>
              </thead>
              <tbody>
                {[
                  { name: 'Project A', status: 'Active', progress: 65 },
                  { name: 'Project B', status: 'Planning', progress: 15 },
                  { name: 'Project C', status: 'Completed', progress: 100 },
                ].map((row, i) => (
                  <tr
                    key={i}
                    className={`border-b border-gray-200 hover:bg-blue-50 ${i % 2 === 0 ? 'bg-white' : 'bg-gray-50'}`}
                  >
                    <td className="px-4 py-3 text-sm text-gray-700">{i + 1}</td>
                    <td className="px-4 py-3 text-sm font-medium text-gray-900">{row.name}</td>
                    <td className="px-4 py-3 text-sm">
                      <span
                        className={`px-2 py-1 rounded text-xs font-semibold ${
                          row.status === 'Active'
                            ? 'bg-green-100 text-green-800'
                            : row.status === 'Planning'
                              ? 'bg-blue-100 text-blue-800'
                              : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {row.status}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm">
                      <div className="h-2 bg-gray-200 rounded-full w-24 overflow-hidden">
                        <div
                          className="h-full bg-blue-500 rounded-full"
                          style={{ width: `${row.progress}%` }}
                        />
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* Typography */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Typography</h2>
          <div className="space-y-4">
            <div>
              <p className="text-4xl font-bold text-gray-900">Heading 1 - 36px Bold</p>
            </div>
            <div>
              <p className="text-3xl font-bold text-gray-900">Heading 2 - 30px Bold</p>
            </div>
            <div>
              <p className="text-2xl font-bold text-gray-900">Heading 3 - 24px Bold</p>
            </div>
            <div>
              <p className="text-lg font-semibold text-gray-800">Subheading - 18px Semibold</p>
            </div>
            <div>
              <p className="text-base text-gray-700">Body Text - 16px Regular</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Small Text - 14px Regular</p>
            </div>
            <div>
              <p className="text-xs text-gray-500">Extra Small - 12px Regular</p>
            </div>
          </div>
        </section>

        {/* Alerts */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Alert Messages</h2>
          <div className="space-y-4">
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <p className="text-sm font-semibold text-blue-900">ℹ️ Info Message</p>
              <p className="text-sm text-blue-800">This is an informational message for users.</p>
            </div>
            <div className="bg-green-50 border border-green-200 rounded-lg p-4">
              <p className="text-sm font-semibold text-green-900">✓ Success Message</p>
              <p className="text-sm text-green-800">Operation completed successfully!</p>
            </div>
            <div className="bg-amber-50 border border-amber-200 rounded-lg p-4">
              <p className="text-sm font-semibold text-amber-900">⚠️ Warning Message</p>
              <p className="text-sm text-amber-800">Please review this information carefully.</p>
            </div>
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-sm font-semibold text-red-900">✕ Error Message</p>
              <p className="text-sm text-red-800">Something went wrong. Please try again.</p>
            </div>
          </div>
        </section>

        {/* Spacing */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Spacing Scale</h2>
          <div className="space-y-3">
            <div className="flex items-center gap-4">
              <span className="w-20 text-sm font-semibold text-gray-600">4px</span>
              <div className="bg-blue-500" style={{ width: '16px', height: '16px' }} />
            </div>
            <div className="flex items-center gap-4">
              <span className="w-20 text-sm font-semibold text-gray-600">8px</span>
              <div className="bg-blue-500" style={{ width: '32px', height: '16px' }} />
            </div>
            <div className="flex items-center gap-4">
              <span className="w-20 text-sm font-semibold text-gray-600">12px</span>
              <div className="bg-blue-500" style={{ width: '48px', height: '16px' }} />
            </div>
            <div className="flex items-center gap-4">
              <span className="w-20 text-sm font-semibold text-gray-600">16px</span>
              <div className="bg-blue-500" style={{ width: '64px', height: '16px' }} />
            </div>
            <div className="flex items-center gap-4">
              <span className="w-20 text-sm font-semibold text-gray-600">24px</span>
              <div className="bg-blue-500" style={{ width: '96px', height: '16px' }} />
            </div>
          </div>
        </section>

        {/* Footer */}
        <footer className="mt-16 pt-8 border-t border-gray-300 text-center text-sm text-gray-600">
          <p>VYOMTECH ERP - Spreadsheet-Style UI Component Library</p>
          <p className="mt-2">Last Updated: December 3, 2025</p>
        </footer>
      </div>
    </div>
  );
};

export default StyleGuide;
