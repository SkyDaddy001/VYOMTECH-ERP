'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { ConstructionProject } from '@/types/construction';

export default function ProjectList() {
  const [projects, setProjects] = useState<ConstructionProject[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchProjects();
  }, []);

  const fetchProjects = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/v1/construction/projects', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      });

      if (!response.ok) throw new Error('Failed to fetch projects');
      const data = await response.json();
      setProjects(data.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div className="text-center py-8">Loading projects...</div>;
  if (error) return <div className="text-red-600 py-8">{error}</div>;

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Construction Projects</h2>
        <Link
          href="/dashboard/construction/projects/new"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          New Project
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {projects.length > 0 ? (
          projects.map((project) => (
            <div key={project.id} className="border rounded-lg p-4 hover:shadow-lg transition">
              <h3 className="font-bold text-lg mb-2">{project.project_name}</h3>
              <p className="text-gray-600 text-sm mb-3">{project.location}</p>
              <div className="space-y-1 text-sm mb-3">
                <p>Client: {project.client}</p>
                <p>Manager: {project.project_manager}</p>
                <p>Progress: <span className="font-semibold">{project.current_progress_percentage}%</span></p>
                <p>Status: <span className="font-semibold">{project.status}</span></p>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2 mb-3">
                <div
                  className="bg-blue-600 h-2 rounded-full"
                  style={{ width: `${project.current_progress_percentage}%` }}
                ></div>
              </div>
              <Link
                href={`/dashboard/construction/projects/${project.id}`}
                className="text-blue-600 hover:underline text-sm"
              >
                View Details →
              </Link>
            </div>
          ))
        ) : (
          <div className="col-span-full text-center py-8 text-gray-500">
            No projects found. Create one to get started.
          </div>
        )}
      </div>
    </div>
  );
}
