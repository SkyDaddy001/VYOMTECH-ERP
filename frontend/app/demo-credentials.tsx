'use client';

import React, { useState } from 'react';

interface DemoCredentialsProps {
  onSelectCredential?: (email: string, password: string) => void;
}

export const DemoCredentials = ({ onSelectCredential }: DemoCredentialsProps) => {
  const [selectedIdx, setSelectedIdx] = useState<number | null>(null);

  const credentials = [
    {
      category: 'System Admin',
      role: 'Master Admin',
      email: 'master.admin@vyomtech.com',
      password: 'demo123',
      description: 'Full system access for demo tenant'
    },
    {
      category: 'Call Center',
      role: 'Agent - Rajesh Kumar',
      email: 'rajesh@demo.vyomtech.com',
      password: 'demo123',
      description: 'Call center agent'
    },
    {
      category: 'Call Center',
      role: 'Agent - Priya Singh',
      email: 'priya@demo.vyomtech.com',
      password: 'demo123',
      description: 'Call center agent'
    },
    {
      category: 'Call Center',
      role: 'Agent - Arun Patel',
      email: 'arun@demo.vyomtech.com',
      password: 'demo123',
      description: 'Call center agent'
    },
    {
      category: 'Call Center',
      role: 'Agent - Neha Sharma',
      email: 'neha@demo.vyomtech.com',
      password: 'demo123',
      description: 'Call center agent'
    },
    {
      category: 'Partners',
      role: 'Portal Admin',
      email: 'demo@vyomtech.com',
      password: 'demo123',
      description: 'White-label portal administrator'
    },
    {
      category: 'Partners',
      role: 'Channel Partner',
      email: 'channel@demo.vyomtech.com',
      password: 'demo123',
      description: 'Channel partner manager'
    },
    {
      category: 'Partners',
      role: 'Vendor Admin',
      email: 'vendor@demo.vyomtech.com',
      password: 'demo123',
      description: 'Vendor/supplier administrator'
    },
    {
      category: 'Partners',
      role: 'Customer Manager',
      email: 'customer@demo.vyomtech.com',
      password: 'demo123',
      description: 'Direct customer account manager'
    }
  ];

  const groupedByCategory = credentials.reduce((acc, cred) => {
    if (!acc[cred.category]) acc[cred.category] = [];
    acc[cred.category].push(cred);
    return acc;
  }, {} as Record<string, typeof credentials>);

  const handleSelectCredential = (index: number, email: string, password: string) => {
    setSelectedIdx(index);
    if (onSelectCredential) {
      onSelectCredential(email, password);
    }
  };

  let credentialIndex = 0;

  return (
    <div className="bg-blue-50 border-l-4 border-blue-500 p-4 mt-6">
      <h3 className="text-lg font-semibold text-blue-900 mb-1">
        🎯 Demo Test Credentials
      </h3>
      <p className="text-xs text-blue-700 mb-4">
        Click any credential below to auto-fill and login instantly:
      </p>
      
      <div className="space-y-6">
        {Object.entries(groupedByCategory).map(([category, creds]) => (
          <div key={category}>
            <h4 className="text-sm font-bold text-blue-800 mb-2 uppercase tracking-wide">
              {category}
            </h4>
            <div className="space-y-2 ml-2">
              {creds.map((cred) => {
                const currentIndex = credentialIndex++;
                return (
                  <button
                    key={currentIndex}
                    onClick={() => handleSelectCredential(currentIndex, cred.email, cred.password)}
                    className={`w-full text-left p-3 rounded border-2 transition-all duration-200 cursor-pointer ${
                      selectedIdx === currentIndex
                        ? 'bg-blue-100 border-blue-500 shadow-md'
                        : 'bg-white border-blue-200 hover:border-blue-400 hover:bg-blue-50 hover:shadow-sm'
                    }`}
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <p className="text-sm font-semibold text-gray-800">{cred.role}</p>
                        <p className="text-xs text-gray-500 mb-2">{cred.description}</p>
                        <div className="space-y-1">
                          <p className="text-xs text-gray-600">
                            Email: <code className="bg-gray-100 px-1.5 py-0.5 rounded font-mono text-xs">{cred.email}</code>
                          </p>
                          <p className="text-xs text-gray-600">
                            Password: <code className="bg-gray-100 px-1.5 py-0.5 rounded font-mono text-xs">{cred.password}</code>
                          </p>
                        </div>
                      </div>
                      {selectedIdx === currentIndex && (
                        <div className="ml-2 mt-1 text-blue-500 text-lg">✓</div>
                      )}
                    </div>
                  </button>
                );
              })}
            </div>
          </div>
        ))}
      </div>

      <div className="mt-4 p-3 bg-amber-50 border border-amber-200 rounded">
        <p className="text-xs text-amber-800">
          <strong>🔐 Security Notice:</strong> These are demo credentials for testing purposes only. 
          Change all passwords in production. Demo data resets every 30 days automatically.
        </p>
      </div>
    </div>
  );
};
