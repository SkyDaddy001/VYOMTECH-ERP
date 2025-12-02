'use client'

import { useState, useEffect } from 'react'
import UserList from '@/components/modules/Users/UserList'
import InviteUserForm from '@/components/modules/Users/InviteUserForm'
import RoleList from '@/components/modules/Users/RoleList'
import RoleForm from '@/components/modules/Users/RoleForm'
import { TenantMember, Role, Permission } from '@/types/user'
import { userService, rbacService } from '@/services/user.service'
import toast from 'react-hot-toast'

type Tab = 'users' | 'roles'

export default function UsersPage() {
  const [activeTab, setActiveTab] = useState<Tab>('users')
  const [users, setUsers] = useState<TenantMember[]>([])
  const [roles, setRoles] = useState<Role[]>([])
  const [permissions, setPermissions] = useState<Permission[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showInviteForm, setShowInviteForm] = useState(false)
  const [showRoleForm, setShowRoleForm] = useState(false)
  const [editingRole, setEditingRole] = useState<Role | undefined>()

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      setError(null)
      await Promise.all([loadUsers(), loadRoles(), loadPermissions()])
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to load data'
      setError(errorMsg)
      toast.error(errorMsg)
      console.error('Error loading data:', err)
    } finally {
      setLoading(false)
    }
  }

  const loadUsers = async () => {
    const data = await userService.getTenantMembers()
    setUsers(data)
  }

  const loadRoles = async () => {
    const data = await rbacService.getRoles()
    setRoles(data)
  }

  const loadPermissions = async () => {
    const data = await rbacService.getPermissions()
    setPermissions(data)
  }

  const handleInviteUser = async (data: { email: string; role: string }) => {
    try {
      await userService.inviteUser(data.email, data.role)
      toast.success('User invited successfully!')
      await loadUsers()
      setShowInviteForm(false)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to invite user'
      toast.error(errorMsg)
      throw err
    }
  }

  const handleRemoveUser = async (user: TenantMember) => {
    try {
      if (confirm(`Are you sure you want to remove ${user.email}?`)) {
        await userService.removeMember(user.user_id)
        toast.success('User removed successfully!')
        await loadUsers()
      }
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to remove user'
      toast.error(errorMsg)
    }
  }

  const handleChangeRole = async (user: TenantMember, newRole: string) => {
    try {
      await userService.updateMemberRole(user.user_id, newRole)
      toast.success('Role updated successfully!')
      await loadUsers()
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to update role'
      toast.error(errorMsg)
    }
  }

  const handleCreateRole = async (data: Partial<Role>) => {
    try {
      if (editingRole) {
        await rbacService.updateRole(editingRole.id, data)
        toast.success('Role updated successfully!')
      } else {
        await rbacService.createRole(data)
        toast.success('Role created successfully!')
      }
      await loadRoles()
      setShowRoleForm(false)
      setEditingRole(undefined)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to save role'
      toast.error(errorMsg)
      throw err
    }
  }

  const handleDeleteRole = async (role: Role) => {
    try {
      if (confirm(`Are you sure you want to delete the role "${role.name}"?`)) {
        await rbacService.deleteRole(role.id)
        toast.success('Role deleted successfully!')
        await loadRoles()
      }
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to delete role'
      toast.error(errorMsg)
    }
  }

  const handleTogglePermission = async (roleId: string, permissionId: string, enable: boolean) => {
    try {
      if (enable) {
        await rbacService.assignPermission(roleId, permissionId)
        toast.success('Permission assigned!')
      } else {
        await rbacService.removePermission(roleId, permissionId)
        toast.success('Permission removed!')
      }
      await loadRoles()
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to update permission'
      toast.error(errorMsg)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Users & Roles</h1>
          <p className="mt-2 text-gray-600">Manage users and their permissions</p>
        </div>
        {activeTab === 'users' && (
          <button
            onClick={() => setShowInviteForm(true)}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"
          >
            + Invite User
          </button>
        )}
        {activeTab === 'roles' && (
          <button
            onClick={() => {
              setEditingRole(undefined)
              setShowRoleForm(true)
            }}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"
          >
            + New Role
          </button>
        )}
      </div>

      {error && (
        <div className="rounded-md bg-red-50 p-4">
          <p className="text-sm font-medium text-red-800">{error}</p>
        </div>
      )}

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <div className="flex space-x-8">
          <button
            onClick={() => setActiveTab('users')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'users'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            Users ({users.length})
          </button>
          <button
            onClick={() => setActiveTab('roles')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'roles'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            Roles ({roles.length})
          </button>
        </div>
      </div>

      {/* Tab Content */}
      {activeTab === 'users' && (
        <UserList
          users={users}
          loading={loading}
          onEdit={(user) => console.log('Edit user:', user)}
          onDelete={handleRemoveUser}
          onChangeRole={handleChangeRole}
        />
      )}

      {activeTab === 'roles' && (
        <RoleList
          roles={roles}
          permissions={permissions}
          loading={loading}
          onEdit={(role) => {
            setEditingRole(role)
            setShowRoleForm(true)
          }}
          onDelete={handleDeleteRole}
          onTogglePermission={handleTogglePermission}
        />
      )}

      {/* Modals */}
      {showInviteForm && (
        <InviteUserForm
          onSubmit={handleInviteUser}
          onCancel={() => setShowInviteForm(false)}
        />
      )}

      {showRoleForm && (
        <RoleForm
          role={editingRole}
          permissions={permissions}
          onSubmit={handleCreateRole}
          onCancel={() => {
            setShowRoleForm(false)
            setEditingRole(undefined)
          }}
        />
      )}
    </div>
  )
}
