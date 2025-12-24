'use client';

import { useEffect } from 'react';
import { useAuthStore } from '@/store/auth';
import { apiClient } from '@/lib/api-client';
import { useRouter } from 'next/navigation';

export default function Home() {
  const router = useRouter();
  const { isAuthenticated, user } = useAuthStore();

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
    }
  }, [isAuthenticated, router]);

  if (!isAuthenticated) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">
          Welcome to VYOM LMS
        </h1>
        <p className="text-lg text-gray-600 mb-8">
          Lead Management System for Sales Teams
        </p>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-2">
              Lead Management
            </h2>
            <p className="text-gray-600 mb-4">
              Create, track, and manage sales leads
            </p>
            <a
              href="/leads"
              className="inline-block px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Go to Leads
            </a>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-2">
              Marketing
            </h2>
            <p className="text-gray-600 mb-4">
              Manage marketing campaigns
            </p>
            <a
              href="/marketing"
              className="inline-block px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Go to Marketing
            </a>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-2">
              Site Visits
            </h2>
            <p className="text-gray-600 mb-4">
              Schedule property site visits
            </p>
            <a
              href="/site-visits"
              className="inline-block px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Go to Site Visits
            </a>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-2">
              Sales
            </h2>
            <p className="text-gray-600 mb-4">
              Track property sales and commissions
            </p>
            <a
              href="/sales"
              className="inline-block px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Go to Sales
            </a>
          </div>
        </div>

        {/* Real Estate Workflow */}
        <div className="mt-12 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-lg p-8">
          <h2 className="text-2xl font-bold mb-4">Real Estate Developer Workflow</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="text-center">
              <div className="text-4xl mb-2">📢</div>
              <h3 className="font-semibold mb-2">Marketing</h3>
              <p className="text-sm text-blue-100">Create campaigns to generate leads</p>
            </div>
            <div className="text-center">
              <div className="text-4xl mb-2">→</div>
              <p className="text-blue-100"></p>
            </div>
            <div className="text-center">
              <div className="text-4xl mb-2">👥</div>
              <h3 className="font-semibold mb-2">Leads</h3>
              <p className="text-sm text-blue-100">Track interested buyers</p>
            </div>
            <div className="text-center">
              <div className="text-4xl mb-2">→</div>
              <p className="text-blue-100"></p>
            </div>
            <div className="text-center">
              <div className="text-4xl mb-2">🏘️</div>
              <h3 className="font-semibold mb-2">Site Visits</h3>
              <p className="text-sm text-blue-100">Schedule property tours</p>
            </div>
            <div className="text-center">
              <div className="text-4xl mb-2">→</div>
              <p className="text-blue-100"></p>
            </div>
            <div className="text-center">
              <div className="text-4xl mb-2">💰</div>
              <h3 className="font-semibold mb-2">Sales</h3>
              <p className="text-sm text-blue-100">Close and manage sales</p>
            </div>
          </div>
        </div>

        <div className="mt-12 bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-semibold text-gray-900 mb-4">
            Sprint 1 Status
          </h2>
          <ul className="space-y-2 text-gray-700">
            <li>✅ Backend infrastructure setup</li>
            <li>✅ Authentication endpoints</li>
            <li>✅ Lead CRUD endpoints (mock)</li>
            <li>⏳ Prisma ORM integration</li>
            <li>⏳ Database connectivity</li>
            <li>⏳ Frontend components</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
