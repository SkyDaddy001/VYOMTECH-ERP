'use client'

import { Role, Permission } from '@/types/user'

interface RoleListProps {
  roles: Role[]
  permissions: Permission[]
  loading: boolean
  onEdit: (role: Role) => void
  onDelete: (role: Role) => void
  onTogglePermission: (roleId: string, permissionId: string, enable: boolean) => void
}

export default function RoleList({
  roles,
  permissions,
  loading,
  onEdit,
  onDelete,
  onTogglePermission,
}: RoleListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading roles...</p>
        </div>
      </div>
    )
  }

  if (roles.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No roles yet. Create your first role to get started.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 gap-6">
        {roles.map((role) => (
          <div key={role.id} className="bg-white rounded-lg shadow border border-gray-200 p-6">
            <div className="flex items-start justify-between mb-4">
              <div>
                <h3 className="text-lg font-bold text-gray-900">{role.name}</h3>
                <p className="text-sm text-gray-600 mt-1">{role.description}</p>
              </div>
              <div className="space-x-2">
                <button
                  onClick={() => onEdit(role)}
                  className="text-blue-600 hover:text-blue-900 font-medium text-sm"
                >
                  Edit
                </button>
                <button
                  onClick={() => onDelete(role)}
                  className="text-red-600 hover:text-red-900 font-medium text-sm"
                >
                  Delete
                </button>
              </div>
            </div>

            {/* Permissions Grid */}
            <div className="border-t pt-4">
              <p className="text-sm font-medium text-gray-700 mb-3">Permissions</p>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                {permissions.map((permission) => {
                  const hasPermission = role.permissions.some((p) => p.id === permission.id)
                  return (
                    <label key={permission.id} className="flex items-center space-x-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={hasPermission}
                        onChange={(e) => onTogglePermission(role.id, permission.id, e.target.checked)}
                        className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
                      />
                      <span className="text-sm text-gray-700">{permission.name}</span>
                    </label>
                  )
                })}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
