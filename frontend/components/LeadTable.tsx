import React from 'react';
import { Badge } from './ui/badge';
import { Button } from './ui/button';

interface Lead {
  id: string;
  name: string;
  email: string;
  phone?: string;
  company?: string;
  status: string;
  createdAt?: string;
}

interface LeadTableProps {
  leads: Lead[];
  onEdit?: (lead: Lead) => void;
  onDelete?: (id: string) => void;
}

export const LeadTable: React.FC<LeadTableProps> = ({ leads, onEdit, onDelete }) => {
  const getStatusBadgeVariant = (status: string) => {
    const variantMap: Record<string, any> = {
      new: 'default',
      qualified: 'primary',
      converted: 'success',
      lost: 'danger',
      contacted: 'warning',
    };
    return variantMap[status] || 'default';
  };

  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="border-b border-gray-200 bg-gray-50">
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Name</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Email</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Company</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Actions</th>
          </tr>
        </thead>
        <tbody>
          {leads.length === 0 ? (
            <tr>
              <td colSpan={5} className="px-6 py-4 text-center text-gray-500">
                No leads found
              </td>
            </tr>
          ) : (
            leads.map((lead) => (
              <tr key={lead.id} className="border-b border-gray-200 hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-medium text-gray-900">{lead.name}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{lead.email}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{lead.company || '-'}</td>
                <td className="px-6 py-4 text-sm">
                  <Badge variant={getStatusBadgeVariant(lead.status)}>
                    {lead.status}
                  </Badge>
                </td>
                <td className="px-6 py-4 text-sm flex gap-2">
                  {onEdit && (
                    <Button size="sm" variant="primary" onClick={() => onEdit(lead)}>
                      Edit
                    </Button>
                  )}
                  {onDelete && (
                    <Button size="sm" variant="danger" onClick={() => onDelete(lead.id)}>
                      Delete
                    </Button>
                  )}
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
};

export default LeadTable;
