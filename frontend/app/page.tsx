export default function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-blue-100 flex items-center justify-center p-4">
      <div className="text-center">
        <h1 className="text-5xl font-bold text-gray-900 mb-4">🎉 Welcome to VYOM ERP</h1>
        <p className="text-xl text-gray-600 mb-8">User Count & Seat Management System</p>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-6xl">
          <a
            href="/user-count"
            className="bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition"
          >
            <div className="text-4xl mb-3">👥</div>
            <h2 className="text-lg font-semibold text-gray-900">User Count</h2>
            <p className="text-sm text-gray-600 mt-2">Real-time user metrics</p>
          </a>

          <a
            href="/user-activity"
            className="bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition"
          >
            <div className="text-4xl mb-3">📊</div>
            <h2 className="text-lg font-semibold text-gray-900">Activity Log</h2>
            <p className="text-sm text-gray-600 mt-2">Login/logout tracking</p>
          </a>

          <a
            href="/seat-management"
            className="bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition"
          >
            <div className="text-4xl mb-3">💺</div>
            <h2 className="text-lg font-semibold text-gray-900">Seat Mgmt</h2>
            <p className="text-sm text-gray-600 mt-2">Manage user limits</p>
          </a>

          <a
            href="/billing"
            className="bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition"
          >
            <div className="text-4xl mb-3">💳</div>
            <h2 className="text-lg font-semibold text-gray-900">Billing</h2>
            <p className="text-sm text-gray-600 mt-2">Overage charges</p>
          </a>
        </div>
      </div>
    </div>
  );
}
