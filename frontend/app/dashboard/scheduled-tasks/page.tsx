// Scheduled Tasks Page - Phase 3B

import { Metadata } from 'next';
import { ScheduledTasks } from '@/components/workflows/ScheduledTasks';

export const metadata: Metadata = {
  title: 'Scheduled Tasks | AI Call Center',
  description: 'Manage scheduled tasks and cron jobs',
};

export default function ScheduledTasksPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Scheduled Tasks</h1>
        <p className="text-gray-600 mt-1">Manage background jobs and automated tasks</p>
      </div>
      <ScheduledTasks />
    </div>
  );
}
