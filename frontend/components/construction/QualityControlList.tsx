'use client';

import { useState, useEffect } from 'react';
import { QualityControl } from '@/types/construction';

export default function QualityControlList() {
  const [inspections, setInspections] = useState<QualityControl[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchInspections();
  }, []);

  const fetchInspections = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/v1/construction/quality', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      });

      if (!response.ok) throw new Error('Failed to fetch inspections');
      const data = await response.json();
      setInspections(data.data || []);
    } catch (error) {
      console.error('Failed to fetch inspections:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'passed':
        return 'bg-green-100 text-green-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      case 'partial':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-blue-100 text-blue-800';
    }
  };

  if (loading) return <div>Loading inspections...</div>;

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Quality Control Inspections</h2>
      {inspections.length > 0 ? (
        <div className="space-y-2">
          {inspections.slice(0, 5).map((inspection) => (
            <div key={inspection.id} className="border rounded p-3 hover:bg-gray-50">
              <div className="flex justify-between items-start">
                <div>
                  <p className="font-semibold">Inspector: {inspection.inspector_name}</p>
                  <p className="text-sm text-gray-600">{inspection.observations}</p>
                </div>
                <span className={`px-2 py-1 rounded text-xs font-semibold ${getStatusColor(inspection.quality_status)}`}>
                  {inspection.quality_status}
                </span>
              </div>
              <p className="text-xs text-gray-500 mt-1">{inspection.inspection_date}</p>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-8 text-gray-500">No inspections found</div>
      )}
    </div>
  );
}
