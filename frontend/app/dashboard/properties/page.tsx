'use client'

import { useState, useMemo } from 'react'
import { useProperties } from '@/hooks/use-api'
import { ProtectedRoute } from '@/hooks/use-auth'
import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { FiSearch, FiPlus, FiMapPin, FiDollarSign, FiHome, FiFilter } from 'react-icons/fi'
import { format } from 'date-fns'

interface Property {
  id: string | number
  name: string
  address?: string
  city?: string
  property_type?: string
  status?: string
  price?: number
  area?: number
  bedrooms?: number
  bathrooms?: number
  description?: string
  listed_date?: string
  created_at: string
}

export default function PropertiesPage() {
  const { data: properties, loading, error } = useProperties({ limit: 100 })
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [sortBy, setSortBy] = useState<'created' | 'price' | 'name' | 'area'>('created')

  const filteredProperties = useMemo(() => {
    let result = [...(properties || [])]

    if (filterStatus !== 'all') {
      result = result.filter(p => p.status === filterStatus)
    }

    if (searchTerm) {
      const query = searchTerm.toLowerCase()
      result = result.filter(p =>
        (p.name || '').toLowerCase().includes(query) ||
        (p.address || '').toLowerCase().includes(query) ||
        (p.city || '').toLowerCase().includes(query)
      )
    }

    result.sort((a, b) => {
      switch (sortBy) {
        case 'price':
          return (b.price || 0) - (a.price || 0)
        case 'area':
          return (b.area || 0) - (a.area || 0)
        case 'name':
          return (a.name || '').localeCompare(b.name || '')
        case 'created':
        default:
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      }
    })

    return result
  }, [properties, searchTerm, filterStatus, sortBy])

  const getStatusColor = (status?: string) => {
    switch (status?.toLowerCase()) {
      case 'available':
      case 'active':
        return 'bg-green-100 text-green-800'
      case 'sold':
        return 'bg-blue-100 text-blue-800'
      case 'rented':
        return 'bg-purple-100 text-purple-800'
      case 'inactive':
        return 'bg-gray-100 text-gray-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getTypeColor = (type?: string) => {
    switch (type?.toLowerCase()) {
      case 'apartment':
        return 'bg-blue-500'
      case 'house':
        return 'bg-green-500'
      case 'commercial':
        return 'bg-yellow-500'
      case 'land':
        return 'bg-orange-500'
      default:
        return 'bg-gray-500'
    }
  }

  return (
    <ProtectedRoute>
      <div className="flex h-screen bg-gray-50">
        <Sidebar />
        <div className="flex-1 flex flex-col lg:ml-64">
          <Header />
          <main className="flex-1 overflow-auto pt-20 pb-6">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
              {/* Header */}
              <div className="mb-8 flex items-center justify-between">
                <div>
                  <h1 className="text-3xl font-bold text-gray-900">Properties</h1>
                  <p className="text-gray-600 mt-2">Manage real estate properties and listings</p>
                </div>
                <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg flex items-center gap-2 transition">
                  <FiPlus className="text-lg" />
                  Add Property
                </button>
              </div>

              {/* Filters & Search */}
              <div className="bg-white rounded-lg shadow p-4 mb-6">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {/* Search */}
                  <div className="relative">
                    <FiSearch className="absolute left-3 top-3 text-gray-400" />
                    <input
                      type="text"
                      placeholder="Search properties..."
                      value={searchTerm}
                      onChange={(e) => setSearchTerm(e.target.value)}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>

                  {/* Status Filter */}
                  <select
                    value={filterStatus}
                    onChange={(e) => setFilterStatus(e.target.value)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="all">All Status</option>
                    <option value="available">Available</option>
                    <option value="sold">Sold</option>
                    <option value="rented">Rented</option>
                    <option value="inactive">Inactive</option>
                  </select>

                  {/* Sort */}
                  <select
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value as any)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="created">Newest First</option>
                    <option value="price">By Price (High to Low)</option>
                    <option value="area">By Area (Largest First)</option>
                    <option value="name">By Name</option>
                  </select>
                </div>

                <div className="mt-4 text-sm text-gray-600">
                  Showing {filteredProperties.length} of {(properties || []).length} properties
                </div>
              </div>

              {/* Properties Grid */}
              {filteredProperties.length === 0 ? (
                <div className="bg-white rounded-lg shadow p-12 text-center">
                  <div className="text-gray-500 mb-4 text-5xl">🏠</div>
                  <p className="text-gray-600 text-lg font-medium">No properties found</p>
                  <p className="text-gray-500 mt-1">Add a new property to get started</p>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {filteredProperties.map((property) => (
                    <div key={property.id} className="bg-white rounded-lg shadow hover:shadow-md transition overflow-hidden">
                      {/* Card Header with Image Placeholder */}
                      <div className="bg-gradient-to-r from-blue-400 to-blue-600 h-48 flex items-center justify-center">
                        <FiHome className="text-white text-6xl opacity-40" />
                      </div>

                      {/* Card Body */}
                      <div className="p-6">
                        {/* Title and Status */}
                        <div className="flex items-start justify-between mb-3">
                          <h3 className="text-lg font-semibold text-gray-900">{property.name}</h3>
                          <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(property.status)}`}>
                            {property.status || 'Unlisted'}
                          </span>
                        </div>

                        {/* Type Badge */}
                        {property.property_type && (
                          <div className="mb-3">
                            <span className={`px-2 py-1 rounded text-xs font-medium text-white ${getTypeColor(property.property_type)}`}>
                              {property.property_type}
                            </span>
                          </div>
                        )}

                        {/* Location */}
                        <div className="flex items-start gap-2 text-sm text-gray-600 mb-4">
                          <FiMapPin className="text-gray-400 mt-1" />
                          <div>
                            <div>{property.address}</div>
                            {property.city && <div className="text-xs text-gray-500">{property.city}</div>}
                          </div>
                        </div>

                        {/* Key Details */}
                        <div className="grid grid-cols-3 gap-2 mb-4 py-4 border-t border-b border-gray-200">
                          {property.bedrooms && (
                            <div className="text-center">
                              <div className="text-sm text-gray-600">Beds</div>
                              <div className="text-lg font-bold text-gray-900">{property.bedrooms}</div>
                            </div>
                          )}
                          {property.bathrooms && (
                            <div className="text-center">
                              <div className="text-sm text-gray-600">Baths</div>
                              <div className="text-lg font-bold text-gray-900">{property.bathrooms}</div>
                            </div>
                          )}
                          {property.area && (
                            <div className="text-center">
                              <div className="text-sm text-gray-600">Area (sqft)</div>
                              <div className="text-lg font-bold text-gray-900">{(property.area / 1000).toFixed(1)}K</div>
                            </div>
                          )}
                        </div>

                        {/* Price */}
                        <div className="mb-4">
                          <div className="flex items-center gap-2 text-gray-600 mb-1">
                            <FiDollarSign className="text-gray-400" />
                            <span className="text-sm">Listed Price</span>
                          </div>
                          <div className="text-2xl font-bold text-gray-900">
                            ${(property.price || 0).toLocaleString()}
                          </div>
                        </div>

                        {/* Description */}
                        {property.description && (
                          <p className="text-sm text-gray-600 line-clamp-2 mb-4">
                            {property.description}
                          </p>
                        )}

                        {/* Listed Date */}
                        {property.listed_date && (
                          <div className="text-xs text-gray-500 mb-4">
                            Listed: {format(new Date(property.listed_date), 'MMM d, yyyy')}
                          </div>
                        )}

                        {/* Action Buttons */}
                        <div className="flex gap-2 pt-4 border-t border-gray-200">
                          <button className="flex-1 px-3 py-2 bg-blue-50 hover:bg-blue-100 text-blue-600 text-sm font-medium rounded-lg transition">
                            View Details
                          </button>
                          <button className="flex-1 px-3 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 text-sm font-medium rounded-lg transition">
                            Edit
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </main>
        </div>
      </div>
    </ProtectedRoute>
  )
}
