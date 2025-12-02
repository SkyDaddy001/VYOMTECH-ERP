'use client';

import { useState, useEffect } from 'react';
import { SafetyIncident } from '@/types/civil';

export default function IncidentsList() {
  const [incidents, setIncidents] = useState<SafetyIncident[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchIncidents();
  }, []);

  const fetchIncidents = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/v1/civil/incidents', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      });

      if (!response.ok) throw new Error('Failed to fetch incidents');
      const data = await response.json();
      setIncidents(data.data || []);
    } catch (error) {
      console.error('Failed to fetch incidents:', error);
    } finally {
      setLoading(false);
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'bg-red-100 text-red-800';
      case 'high':
        return 'bg-orange-100 text-orange-800';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-blue-100 text-blue-800';
    }
  };

  if (loading) return <div>Loading incidents...</div>;

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Recent Incidents</h2>
      {incidents.length > 0 ? (
        <div className="space-y-2">
          {incidents.slice(0, 5).map((incident) => (
            <div key={incident.id} className="border rounded p-3 hover:bg-gray-50">
              <div className="flex justify-between items-start">
                <div>
                  <p className="font-semibold">{incident.incident_number}</p>
                  <p className="text-sm text-gray-600">{incident.description}</p>
                </div>
                <span className={`px-2 py-1 rounded text-xs font-semibold ${getSeverityColor(incident.severity)}`}>
                  {incident.severity}
                </span>
              </div>
              <p className="text-xs text-gray-500 mt-1">{incident.incident_date}</p>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-8 text-gray-500">No incidents found</div>
      )}
    </div>
  );
}
