// Workflow Executions Page - Phase 3B

import { Metadata } from 'next';
import { WorkflowProvider } from '@/contexts/WorkflowContext';
import { WorkflowExecutions } from '@/components/workflows/WorkflowExecutions';

export const metadata: Metadata = {
  title: 'Workflow Executions | AI Call Center',
  description: 'View workflow execution history',
};

interface PageProps {
  params: {
    id: string;
  };
}

export default function WorkflowExecutionsPage({ params }: PageProps) {
  return (
    <WorkflowProvider>
      <WorkflowExecutions workflowId={params.id} />
    </WorkflowProvider>
  );
}
