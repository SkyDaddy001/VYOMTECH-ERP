'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useLeadStore } from '@/store/leads';
import { apiClient } from '@/lib/api-client';
import CreateLeadModal from '@/components/CreateLeadModal';
import LeadTable from '@/components/LeadTable';

export default function LeadsPage() {
  const router = useRouter();
  const { leads, setLeads, isLoading, setIsLoading, error, setError } = useLeadStore();
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [totalPages, setTotalPages] = useState(1);
  const [statusFilter, setStatusFilter] = useState('');

  useEffect(() => {
    fetchLeads();
  }, [page, pageSize, statusFilter]);

  const fetchLeads = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const filters: any = {};
      if (statusFilter) filters.status = statusFilter;

      const response = await apiClient.listLeads(page, pageSize, filters);

      if (response.data) {
        setLeads(response.data);
        if (response.total_pages) {
          setTotalPages(response.total_pages);
        }
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to fetch leads');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateLead = async (leadData: any) => {
    try {
      await apiClient.createLead(leadData);
      setShowCreateModal(false);
      await fetchLeads();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create lead');
    }
  };

  const handleDeleteLead = async (id: string) => {
    if (confirm('Are you sure you want to delete this lead?')) {
      try {
        await apiClient.deleteLead(id);
        await fetchLeads();
      } catch (err: any) {
        setError(err.response?.data?.message || 'Failed to delete lead');
      }
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Sales Leads</h1>
          <button
            onClick={() => setShowCreateModal(true)}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            + New Lead
          </button>
        </div>

        {error && (
          <div className="mb-4 rounded-md bg-red-50 p-4">
            <p className="text-sm font-medium text-red-800">{error}</p>
          </div>
        )}

        <div className="mb-6 flex items-center gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Filter by Status
            </label>
            <select
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value);
                setPage(1);
              }}
              className="px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Statuses</option>
              <option value="new">New</option>
              <option value="qualified">Qualified</option>
              <option value="contacted">Contacted</option>
              <option value="proposal_sent">Proposal Sent</option>
              <option value="negotiation">Negotiation</option>
              <option value="converted">Converted</option>
              <option value="lost">Lost</option>
            </select>
          </div>
        </div>

        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-gray-500">Loading leads...</p>
          </div>
        ) : leads.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500">No leads found. Create your first lead.</p>
          </div>
        ) : (
          <LeadTable leads={leads} onDelete={handleDeleteLead} />
        )}

        {totalPages > 1 && (
          <div className="mt-8 flex items-center justify-between">
            <button
              onClick={() => setPage(Math.max(1, page - 1))}
              disabled={page === 1}
              className="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
            >
              Previous
            </button>
            <span className="text-sm text-gray-600">
              Page {page} of {totalPages}
            </span>
            <button
              onClick={() => setPage(Math.min(totalPages, page + 1))}
              disabled={page === totalPages}
              className="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
            >
              Next
            </button>
          </div>
        )}
      </div>

      {showCreateModal && (
        <CreateLeadModal
          onClose={() => setShowCreateModal(false)}
          onCreate={handleCreateLead}
        />
      )}
    </div>
  );
}
