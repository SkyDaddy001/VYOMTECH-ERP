'use client'

import React, { useState, useEffect } from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'
import { inventoryDashboardService } from '@/services/api'
import { Package, Home, TrendingUp, AlertTriangle } from 'lucide-react'

export default function InventoryPresentationDashboard() {
  // State for inventory data
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [inventorySummary, setInventorySummary] = useState<any>(null)
  const [warehouseData, setWarehouseData] = useState<any>(null)
  const [realEstateData, setRealEstateData] = useState<any>(null)
  const [inventoryByWarehouse, setInventoryByWarehouse] = useState<any>(null)

  // Fetch inventory data on mount
  useEffect(() => {
    const fetchInventoryData = async () => {
      try {
        setLoading(true)

        // Fetch inventory summary
        const summaryRes = await inventoryDashboardService.getInventorySummary()
        setInventorySummary(summaryRes.data)

        // Fetch warehouse distribution
        const warehouseRes = await inventoryDashboardService.getWarehouseDistribution()
        setWarehouseData(warehouseRes.data)

        // Fetch real estate summary
        const realEstateRes = await inventoryDashboardService.getRealEstateSummary()
        setRealEstateData(realEstateRes.data)

        // Fetch inventory by warehouse
        const invByWhRes = await inventoryDashboardService.getInventoryByWarehouse()
        setInventoryByWarehouse(invByWhRes.data)

        setError(null)
      } catch (err: any) {
        console.error('Failed to fetch inventory data:', err)
        setError(err.message || 'Failed to load inventory data')
      } finally {
        setLoading(false)
      }
    }

    fetchInventoryData()
  }, [])

  // Use real data or fallback values
  const inventoryValue = inventorySummary?.total_value || 850000000
  const stockUnits = inventorySummary?.total_units || 12450
  const warehouses = inventorySummary?.warehouse_count || 24
  const spaceUtilization = inventorySummary?.space_utilization || 92

  const stockHealth = inventorySummary?.stock_health || [
    { status: 'Optimal Stock', count: 4280, percent: 34 },
    { status: 'Low Stock', count: 3120, percent: 25 },
    { status: 'Critical Stock', count: 1890, percent: 15 },
    { status: 'Excess Stock', count: 3160, percent: 26 }
  ]

  const criticalItems = inventorySummary?.critical_items || [
    { sku: 'SKU-2845', name: 'Bearing Assembly A', stock: 45, minLevel: 150, daysToStock: 3 },
    { sku: 'SKU-1203', name: 'Steel Plate 10mm', stock: 120, minLevel: 200, daysToStock: 2 },
    { sku: 'SKU-5678', name: 'Motor Coupling', stock: 28, minLevel: 80, daysToStock: 1 },
    { sku: 'SKU-4523', name: 'Control Panel Box', stock: 15, minLevel: 50, daysToStock: 1 }
  ]

  const whDistribution = warehouseData?.warehouses || [
    { name: 'Main Warehouse - Mumbai', stock: 3500, capacity: 4000, utilization: 87, type: 'Primary' },
    { name: 'Secondary Hub - Delhi', stock: 2200, capacity: 2500, utilization: 88, type: 'Secondary' },
    { name: 'Regional Center - Bangalore', stock: 1800, capacity: 2200, utilization: 82, type: 'Secondary' },
    { name: 'Distribution Center - Pune', stock: 1950, capacity: 2200, utilization: 89, type: 'Distribution' },
    { name: 'Transit Hub - Chennai', stock: 950, capacity: 1500, utilization: 63, type: 'Distribution' },
    { name: 'Field Depot - Hyderabad', stock: 1050, capacity: 1200, utilization: 87, type: 'Field' }
  ]

  const properties = realEstateData?.properties || [
    { name: 'Head Office - Mumbai', area: 45000, type: 'Owned', value: 4500000000, status: 'Active' },
    { name: 'Manufacturing Plant - Pune', area: 120000, type: 'Owned', value: 3500000000, status: 'Active' },
    { name: 'R&D Center - Bangalore', area: 25000, type: 'Leased', value: 800000000, status: 'Active' },
    { name: 'Regional Office - Delhi', area: 12000, type: 'Leased', value: 200000000, status: 'Active' }
  ]

  const logistics = inventorySummary?.logistics || {
    inbound: 2840,
    outbound: 3210,
    turnover: 8.3,
    on_time_rate: 96,
    damage_rate: 0.8,
    accuracy: 99.2
  }

  const formatCurrency = (value: number) => {
    if (value >= 10000000) {
      return `₹${(value / 10000000).toFixed(1)} Cr`
    } else if (value >= 100000) {
      return `₹${(value / 100000).toFixed(0)} L`
    }
    return `₹${value.toLocaleString('en-IN')}`
  }

  const slides: Slide[] = [
    {
      id: 'cover',
      title: 'Inventory & Real Estate',
      subtitle: 'Stock Management & Property Portfolio Overview',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <Package className="w-20 h-20 text-blue-600" />
          <div className="grid grid-cols-2 gap-6 w-full max-w-2xl">
            <div className="bg-blue-50 p-6 rounded-lg border border-blue-200">
              <div className="text-3xl font-bold text-blue-700">{formatCurrency(inventoryValue)}</div>
              <div className="text-sm text-gray-600 mt-1">Inventory Value</div>
            </div>
            <div className="bg-green-50 p-6 rounded-lg border border-green-200">
              <div className="text-3xl font-bold text-green-700">{stockUnits.toLocaleString('en-IN')}</div>
              <div className="text-sm text-gray-600 mt-1">Stock Units</div>
            </div>
            <div className="bg-purple-50 p-6 rounded-lg border border-purple-200">
              <div className="text-3xl font-bold text-purple-700">{warehouses}</div>
              <div className="text-sm text-gray-600 mt-1">Warehouses</div>
            </div>
            <div className="bg-orange-50 p-6 rounded-lg border border-orange-200">
              <div className="text-3xl font-bold text-orange-700">{spaceUtilization}%</div>
              <div className="text-sm text-gray-600 mt-1">Space Utilization</div>
            </div>
          </div>
          {error && <p className="text-red-600 text-sm">{error}</p>}
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'stock-health',
      title: 'Stock Health & Levels',
      subtitle: 'Current inventory status and critical items',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Inventory Health Summary</h3>
            {stockHealth.map((item: any, i: number) => {
              const colors = ['green', 'yellow', 'red', 'blue']
              return (
                <div key={i} className={`bg-${colors[i]}-50 border-l-4 border-${colors[i]}-500 p-3 rounded`}>
                  <div className="flex justify-between items-center">
                    <span className="font-semibold text-gray-800 text-sm">{item.status}</span>
                    <div className="text-right">
                      <div className="text-lg font-bold text-gray-800">{item.count}</div>
                      <div className="text-xs text-gray-600">{item.percent}% of total</div>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Critical Stock Items (Action Required)</h3>
            {criticalItems.map((item: any, i: number) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-red-300">
                <div className="flex justify-between items-start mb-1">
                  <div>
                    <div className="font-bold text-gray-800 text-sm">{item.name}</div>
                    <div className="text-xs text-gray-600">{item.sku}</div>
                  </div>
                  <span className="text-xs bg-red-100 text-red-800 px-2 py-0.5 rounded font-bold">{item.daysToStock}d</span>
                </div>
                <div className="flex justify-between text-xs text-gray-600">
                  <span>Current: {item.stock}</span>
                  <span>Min: {item.minLevel}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-yellow-50'
    },
    {
      id: 'warehouse',
      title: 'Warehouse Distribution',
      subtitle: 'Stock across warehouse locations',
      content: (
        <div className="space-y-3 h-full overflow-y-auto">
          <div className="grid grid-cols-3 gap-3">
            {whDistribution.map((wh: any, i: number) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200">
                <div className="flex items-start justify-between mb-2">
                  <div>
                    <h4 className="font-bold text-gray-800 text-sm">{wh.name}</h4>
                    <div className="text-xs text-gray-500 mt-0.5">{wh.type}</div>
                  </div>
                  <div className="text-right">
                    <div className="text-lg font-bold text-blue-600">{wh.stock}</div>
                    <div className="text-xs text-gray-600">units</div>
                  </div>
                </div>
                <div className="w-full bg-gray-200 h-2 rounded-full overflow-hidden mb-1">
                  <div className="bg-blue-500 h-full" style={{ width: `${wh.utilization}%` }}></div>
                </div>
                <div className="text-xs text-gray-600">{wh.utilization}% of {wh.capacity} capacity</div>
              </div>
            ))}
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'realestate',
      title: 'Real Estate Portfolio',
      subtitle: 'Property holdings and lease analysis',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Property Assets</h3>
            {properties.map((prop: any, i: number) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200">
                <div className="flex justify-between items-start">
                  <div>
                    <h4 className="font-bold text-gray-800 text-sm">{prop.name}</h4>
                    <div className="text-xs text-gray-600 mt-1">{prop.area.toLocaleString()} sq ft</div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-bold text-blue-600">{formatCurrency(prop.value)}</div>
                    <div className="text-xs px-1.5 py-0.5 rounded mt-1 bg-blue-100 text-blue-800">{prop.type}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Portfolio Analysis</h3>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="text-sm mb-3">
                <span className="font-semibold text-gray-800">Total Property Value</span>
                <span className="float-right font-bold text-blue-600">{formatCurrency(9000000000)}</span>
              </div>
              <div className="space-y-2 text-sm mb-3">
                <div className="flex justify-between">
                  <span className="text-gray-600">Owned Properties: {formatCurrency(8000000000)} (89%)</span>
                  <div className="w-16 bg-gray-200 h-2 rounded-full overflow-hidden">
                    <div className="bg-green-500 h-full" style={{ width: '100%' }}></div>
                  </div>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Leased Properties: {formatCurrency(1000000000)} (11%)</span>
                  <div className="w-16 bg-gray-200 h-2 rounded-full overflow-hidden">
                    <div className="bg-blue-500 h-full" style={{ width: '100%' }}></div>
                  </div>
                </div>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <h4 className="font-bold text-gray-800 text-sm mb-2">Lease Renewals Due</h4>
              <div className="space-y-2">
                <div className="flex justify-between text-xs">
                  <span>R&D Center Lease</span>
                  <span className="text-orange-600 font-bold">3 months</span>
                </div>
                <div className="flex justify-between text-xs">
                  <span>Regional Office</span>
                  <span className="text-red-600 font-bold">2 months</span>
                </div>
                <div className="flex justify-between text-xs">
                  <span>Sales Office - Pune</span>
                  <span className="text-red-600 font-bold">1 month</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-green-50'
    },
    {
      id: 'logistics',
      title: 'Logistics & Movement',
      subtitle: 'Inbound/outbound shipments and turns',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">This Month Activity</h3>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="text-sm text-gray-600">Inbound Shipments</div>
              <div className="text-3xl font-bold text-blue-700 mt-2">{logistics.inbound.toLocaleString()}</div>
              <div className="text-xs text-gray-600 mt-1">units received</div>
            </div>
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="text-sm text-gray-600">Outbound Shipments</div>
              <div className="text-3xl font-bold text-green-700 mt-2">{logistics.outbound.toLocaleString()}</div>
              <div className="text-xs text-gray-600 mt-1">units dispatched</div>
            </div>
            <div className="bg-purple-50 border-l-4 border-purple-500 p-4 rounded">
              <div className="text-sm text-gray-600">Avg Inventory Turns</div>
              <div className="text-3xl font-bold text-purple-700 mt-2">{logistics.turnover}x</div>
              <div className="text-xs text-gray-600 mt-1">per year (industry: 6x)</div>
            </div>
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Logistics Performance</h3>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="space-y-2">
                <div>
                  <div className="flex justify-between mb-1">
                    <span className="text-sm font-semibold">Delivery On-Time Rate</span>
                    <span className="text-sm font-bold">{logistics.on_time_rate}%</span>
                  </div>
                  <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                    <div className="bg-green-500 h-full" style={{ width: `${logistics.on_time_rate}%` }}></div>
                  </div>
                </div>
                <div>
                  <div className="flex justify-between mb-1">
                    <span className="text-sm font-semibold">Damage Rate</span>
                    <span className="text-sm font-bold">{logistics.damage_rate}%</span>
                  </div>
                  <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                    <div className="bg-orange-500 h-full" style={{ width: `${logistics.damage_rate * 10}%` }}></div>
                  </div>
                </div>
                <div>
                  <div className="flex justify-between mb-1">
                    <span className="text-sm font-semibold">Warehouse Accuracy</span>
                    <span className="text-sm font-bold">{logistics.accuracy}%</span>
                  </div>
                  <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                    <div className="bg-blue-500 h-full" style={{ width: `${logistics.accuracy}%` }}></div>
                  </div>
                </div>
              </div>
            </div>
            <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
              <div className="font-bold text-yellow-900">⚡ Improvement Opportunity</div>
              <div className="text-sm text-gray-700 mt-1">Implement RFID tags for real-time inventory. Potential 2% efficiency gain = {formatCurrency(17000000)} annual savings.</div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-purple-50'
    },
    {
      id: 'summary',
      title: 'Summary & Strategic Actions',
      subtitle: 'Key achievements and optimization focus',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="font-bold text-green-900">✓ Optimal Utilization</div>
              <div className="text-sm text-gray-700 mt-2">{spaceUtilization}% warehouse utilization indicates efficient space management. {logistics.turnover}x inventory turns exceeds industry benchmark (6x).</div>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="font-bold text-blue-900">📊 Strong Logistics</div>
              <div className="text-sm text-gray-700 mt-2">{logistics.on_time_rate}% on-time delivery, {logistics.damage_rate}% damage rate, {logistics.accuracy}% inventory accuracy. Top-quartile performance across KPIs.</div>
            </div>
            <div className="bg-purple-50 border-l-4 border-purple-500 p-4 rounded">
              <div className="font-bold text-purple-900">🏢 Real Estate Value</div>
              <div className="text-sm text-gray-700 mt-2">{formatCurrency(9000000000)} portfolio. 89% owned properties = strategic asset. Strong balance sheet strengthener.</div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded">
              <div className="font-bold text-red-900">⚠️ Immediate Actions</div>
              <div className="text-sm text-gray-700 mt-2">
                • Replenish 4 critical SKUs (by Dec 10)<br/>
                • Review excess stock items<br/>
                • Renew Sales Office lease (urgent)<br/>
                • Execute RFID pilot program
              </div>
            </div>
            <div className="bg-orange-50 border-l-4 border-orange-500 p-4 rounded">
              <div className="font-bold text-orange-900">🎯 Q1 2025 Initiatives</div>
              <div className="text-sm text-gray-700 mt-2">
                • RFID implementation (3-month ROI)<br/>
                • Warehouse automation study<br/>
                • Lease optimization review<br/>
                • Inventory forecasting system upgrade
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    }
  ]

  return <PresentationDashboard slides={slides} title="Inventory & Real Estate Dashboard" showSlideNumbers={true} />
}
