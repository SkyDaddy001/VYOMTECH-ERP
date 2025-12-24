'use client';

import { formatDistanceToNow } from 'date-fns';
import Link from 'next/link';

interface Lead {
  id: string;
  lead_code: string;
  lead_name: string;
  email?: string;
  company?: string;
  lead_status: string;
  lead_value?: number;
  created_at: string;
}

interface LeadTableProps {
  leads: Lead[];
  onDelete: (id: string) => void;
}

const statusColors = {
  new: 'bg-blue-100 text-blue-800',
  qualified: 'bg-green-100 text-green-800',
  contacted: 'bg-yellow-100 text-yellow-800',
  proposal_sent: 'bg-purple-100 text-purple-800',
  negotiation: 'bg-orange-100 text-orange-800',
  converted: 'bg-teal-100 text-teal-800',
  lost: 'bg-red-100 text-red-800',
};

export default function LeadTable({ leads, onDelete }: LeadTableProps) {
  return (
    <div className="bg-white rounded-lg shadow overflow-hidden">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Lead Code
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Name
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Company
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Value
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Created
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {leads.map((lead) => (
            <tr key={lead.id} className="hover:bg-gray-50">
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                {lead.lead_code}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {lead.lead_name}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                {lead.company || '-'}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span
                  className={`px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${
                    statusColors[lead.lead_status as keyof typeof statusColors] ||
                    'bg-gray-100 text-gray-800'
                  }`}
                >
                  {lead.lead_status}
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                {lead.lead_value ? `$${lead.lead_value.toLocaleString()}` : '-'}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                {formatDistanceToNow(new Date(lead.created_at), { addSuffix: true })}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                <Link
                  href={`/leads/${lead.id}`}
                  className="text-blue-600 hover:text-blue-900"
                >
                  View
                </Link>
                <button
                  onClick={() => onDelete(lead.id)}
                  className="text-red-600 hover:text-red-900"
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
