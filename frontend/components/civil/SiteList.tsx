'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { Site } from '@/types/civil';

export default function SiteList() {
  const [sites, setSites] = useState<Site[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchSites();
  }, []);

  const fetchSites = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/v1/civil/sites', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      });

      if (!response.ok) throw new Error('Failed to fetch sites');
      const data = await response.json();
      setSites(data.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div className="text-center py-8">Loading sites...</div>;
  if (error) return <div className="text-red-600 py-8">{error}</div>;

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Sites</h2>
        <Link
          href="/dashboard/civil/sites/new"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          New Site
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {sites.length > 0 ? (
          sites.map((site) => (
            <div key={site.id} className="border rounded-lg p-4 hover:shadow-lg transition">
              <h3 className="font-bold text-lg mb-2">{site.site_name}</h3>
              <p className="text-gray-600 text-sm mb-3">{site.location}</p>
              <div className="space-y-1 text-sm mb-3">
                <p>Manager: {site.site_manager}</p>
                <p>Status: <span className="font-semibold">{site.current_status}</span></p>
                <p>Workforce: {site.workforce_count}</p>
              </div>
              <Link
                href={`/dashboard/civil/sites/${site.id}`}
                className="text-blue-600 hover:underline text-sm"
              >
                View Details →
              </Link>
            </div>
          ))
        ) : (
          <div className="col-span-full text-center py-8 text-gray-500">
            No sites found. Create one to get started.
          </div>
        )}
      </div>
    </div>
  );
}
