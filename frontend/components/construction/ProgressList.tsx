'use client';

import { useState, useEffect } from 'react';
import { ProgressTracking } from '@/types/construction';

export default function ProgressList() {
  const [progress, setProgress] = useState<ProgressTracking[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchProgress();
  }, []);

  const fetchProgress = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/v1/construction/progress', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      });

      if (!response.ok) throw new Error('Failed to fetch progress');
      const data = await response.json();
      setProgress(data.data || []);
    } catch (error) {
      console.error('Failed to fetch progress:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>Loading progress...</div>;

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Recent Progress</h2>
      {progress.length > 0 ? (
        <div className="space-y-3">
          {progress.slice(0, 5).map((item) => (
            <div key={item.id} className="border rounded p-4 hover:bg-gray-50">
              <div className="flex justify-between items-start mb-2">
                <p className="font-semibold">{item.activity_description}</p>
                <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded text-xs font-semibold">
                  {item.percentage_complete}%
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2 mb-2">
                <div
                  className="bg-green-500 h-2 rounded-full"
                  style={{ width: `${item.percentage_complete}%` }}
                ></div>
              </div>
              <p className="text-xs text-gray-500">{item.date}</p>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-8 text-gray-500">No progress records found</div>
      )}
    </div>
  );
}
