'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import userService, { type User } from '@/services/user-service';
import Link from 'next/link';

export default function UserFormPage() {
  const router = useRouter();
  const params = useParams();
  const userId = params.id as string | undefined;
  const isEditing = !!userId && userId !== 'create';

  const [user, setUser] = useState<Partial<User>>({
    fullName: '',
    email: '',
    roleId: '',
    status: 'active',
  });
  const [isLoading, setIsLoading] = useState(isEditing);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [roles, setRoles] = useState<{ id: string; name: string }[]>([
    { id: 'admin', name: 'Admin' },
    { id: 'manager', name: 'Manager' },
    { id: 'user', name: 'User' },
    { id: 'viewer', name: 'Viewer' },
  ]);

  useEffect(() => {
    if (isEditing && userId) {
      loadUser();
    }
  }, [isEditing, userId]);

  const loadUser = async () => {
    try {
      if (userId) {
        const userData = await userService.getUser(userId);
        setUser(userData);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load user');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsSaving(true);

    try {
      if (isEditing && userId) {
        await userService.updateUser(userId, {
          fullName: user.fullName || '',
          roleId: user.roleId || '',
          status: user.status as any,
        });
        router.push('/users');
      } else {
        await userService.createUser({
          email: user.email || '',
          fullName: user.fullName || '',
          roleId: user.roleId || '',
          sendInvitation: true,
        });
        router.push('/users');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save user');
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <Link href="/users" className="text-blue-600 hover:text-blue-800 font-medium mb-4 inline-block">
            ← Back to Users
          </Link>
          <h1 className="text-3xl font-bold text-gray-900">
            {isEditing ? 'Edit User' : 'Create New User'}
          </h1>
        </div>

        {/* Error Alert */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
            {error}
          </div>
        )}

        {/* Form */}
        <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-8 space-y-6">
          {/* Full Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Full Name</label>
            <input
              type="text"
              value={user.fullName || ''}
              onChange={(e) => setUser({ ...user, fullName: e.target.value })}
              placeholder="John Doe"
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          {/* Email */}
          {!isEditing && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Email</label>
              <input
                type="email"
                value={user.email || ''}
                onChange={(e) => setUser({ ...user, email: e.target.value })}
                placeholder="john@example.com"
                required
                disabled={isEditing}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              />
              <p className="text-sm text-gray-600 mt-1">An invitation will be sent to this email</p>
            </div>
          )}

          {/* Role */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Role</label>
            <select
              value={user.roleId || ''}
              onChange={(e) => setUser({ ...user, roleId: e.target.value })}
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Select a role</option>
              {roles.map((role) => (
                <option key={role.id} value={role.id}>
                  {role.name}
                </option>
              ))}
            </select>
          </div>

          {/* Status */}
          {isEditing && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
              <select
                value={user.status || 'active'}
                onChange={(e) => setUser({ ...user, status: e.target.value as any })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
                <option value="suspended">Suspended</option>
              </select>
            </div>
          )}

          {/* Submit */}
          <div className="flex gap-3 pt-4">
            <button
              type="submit"
              disabled={isSaving}
              className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg transition"
            >
              {isSaving ? 'Saving...' : isEditing ? 'Update User' : 'Create User'}
            </button>
            <Link
              href="/users"
              className="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-900 font-semibold py-2 px-4 rounded-lg transition text-center"
            >
              Cancel
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
