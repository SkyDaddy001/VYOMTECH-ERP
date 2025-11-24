// Workflows Dashboard Page - Phase 3B

import { Metadata } from 'next';
import { WorkflowProvider } from '@/contexts/WorkflowContext';
import { WorkflowList } from '@/components/workflows/WorkflowList';

export const metadata: Metadata = {
  title: 'Workflows | AI Call Center',
  description: 'Manage workflow automations',
};

export default function WorkflowsPage() {
  return (
    <WorkflowProvider>
      <div className="space-y-8">
        <WorkflowList />
      </div>
    </WorkflowProvider>
  );
}
