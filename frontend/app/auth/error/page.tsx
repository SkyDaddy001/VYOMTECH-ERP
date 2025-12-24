/**
 * OAuth Error Page
 */

'use client';

import Link from 'next/link';
import { useSearchParams } from 'next/navigation';

export default function ErrorPage() {
  const searchParams = useSearchParams();
  const error = searchParams.get('error');
  const provider = searchParams.get('provider');

  const errorMessages: Record<string, string> = {
    missing_params: 'Invalid request parameters',
    callback_failed: 'OAuth callback failed',
    invalid_code: 'Invalid authorization code',
    access_denied: 'Access denied by user',
    server_error: 'Server error during authentication',
  };

  const message = errorMessages[error || ''] || 'Unknown error occurred';

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-red-50 to-orange-100">
      <div className="w-full max-w-md p-8 bg-white rounded-lg shadow-lg">
        <div className="text-center mb-6">
          <div className="text-6xl mb-4">❌</div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">
            Authentication Error
          </h1>
          <p className="text-gray-600 mb-4">{message}</p>
          {provider && (
            <p className="text-sm text-gray-500">
              Provider: {provider.charAt(0).toUpperCase() + provider.slice(1)}
            </p>
          )}
        </div>

        <div className="space-y-4">
          <Link href="/auth/login">
            <button className="w-full py-3 px-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition duration-200">
              Try Again
            </button>
          </Link>
          <Link href="/">
            <button className="w-full py-3 px-4 bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold rounded-lg transition duration-200">
              Go Home
            </button>
          </Link>
        </div>
      </div>
    </div>
  );
}
