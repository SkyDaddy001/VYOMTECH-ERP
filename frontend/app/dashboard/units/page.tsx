'use client'

import { useState, useEffect } from 'react'
import UnitList from '@/components/modules/Units/UnitList'
import UnitForm from '@/components/modules/Units/UnitForm'
import CostSheetForm from '@/components/modules/Units/CostSheetForm'
import { PropertyUnit, UnitCostSheet } from '@/types/unit'
import { unitService } from '@/services/unit.service'
import toast from 'react-hot-toast'

export default function UnitsPage() {
  const [units, setUnits] = useState<PropertyUnit[]>([])
  const [costSheets, setCostSheets] = useState<Map<string, UnitCostSheet>>(new Map())
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showUnitForm, setShowUnitForm] = useState(false)
  const [editingUnit, setEditingUnit] = useState<PropertyUnit | undefined>()
  const [showCostSheet, setShowCostSheet] = useState(false)
  const [selectedUnit, setSelectedUnit] = useState<PropertyUnit | undefined>()
  const [currentCostSheet, setCurrentCostSheet] = useState<UnitCostSheet | undefined>()

  useEffect(() => {
    loadUnits()
  }, [])

  const loadUnits = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await unitService.getUnits()
      setUnits(data)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to load units'
      setError(errorMsg)
      toast.error(errorMsg)
      console.error('Error loading units:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateUnit = async (data: Partial<PropertyUnit>) => {
    try {
      // Note: projectId should be passed from context or params
      const projectId = 'default-project'
      await unitService.createUnit(projectId, data)
      toast.success('Unit created successfully!')
      await loadUnits()
      setShowUnitForm(false)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to create unit'
      toast.error(errorMsg)
      throw err
    }
  }

  const handleEditUnit = async (data: Partial<PropertyUnit>) => {
    try {
      if (!editingUnit) return
      await unitService.updateUnit(editingUnit.id, data)
      toast.success('Unit updated successfully!')
      await loadUnits()
      setShowUnitForm(false)
      setEditingUnit(undefined)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to update unit'
      toast.error(errorMsg)
      throw err
    }
  }

  const handleDeleteUnit = async (unit: PropertyUnit) => {
    try {
      if (confirm(`Are you sure you want to delete unit ${unit.unit_number}?`)) {
        await unitService.deleteUnit(unit.id)
        toast.success('Unit deleted successfully!')
        await loadUnits()
      }
    } catch (err) {
      toast.error('Failed to delete unit')
    }
  }

  const handleViewCostSheet = async (unit: PropertyUnit) => {
    try {
      setSelectedUnit(unit)
      const costSheet = await unitService.getCostSheet(unit.id)
      setCurrentCostSheet(costSheet)
      setCostSheets(new Map(costSheets).set(unit.id, costSheet))
      setShowCostSheet(true)
    } catch (err) {
      toast.error('Failed to load cost sheet')
    }
  }

  const handleUpdateCostSheet = async (data: Partial<UnitCostSheet>) => {
    try {
      if (!selectedUnit) return
      await unitService.updateCostSheet(selectedUnit.id, data)
      toast.success('Cost sheet updated successfully!')
      setShowCostSheet(false)
      setSelectedUnit(undefined)
      setCurrentCostSheet(undefined)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to update cost sheet'
      toast.error(errorMsg)
      throw err
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Units</h1>
          <p className="mt-2 text-gray-600">Manage property units and their pricing</p>
        </div>
        <button
          onClick={() => {
            setEditingUnit(undefined)
            setShowUnitForm(true)
          }}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"
        >
          + New Unit
        </button>
      </div>

      {error && (
        <div className="rounded-md bg-red-50 p-4">
          <p className="text-sm font-medium text-red-800">{error}</p>
        </div>
      )}

      <UnitList
        units={units}
        loading={loading}
        onEdit={(unit) => {
          setEditingUnit(unit)
          setShowUnitForm(true)
        }}
        onDelete={handleDeleteUnit}
        onViewCostSheet={handleViewCostSheet}
      />

      {showUnitForm && (
        <UnitForm
          unit={editingUnit}
          onSubmit={editingUnit ? handleEditUnit : handleCreateUnit}
          onCancel={() => {
            setShowUnitForm(false)
            setEditingUnit(undefined)
          }}
        />
      )}

      {showCostSheet && (
        <CostSheetForm
          costSheet={currentCostSheet}
          onSubmit={handleUpdateCostSheet}
          onCancel={() => {
            setShowCostSheet(false)
            setSelectedUnit(undefined)
            setCurrentCostSheet(undefined)
          }}
        />
      )}
    </div>
  )
}
