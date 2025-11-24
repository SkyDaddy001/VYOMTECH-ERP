// Edit Workflow Page - Phase 3B

import { Metadata } from 'next';
import { WorkflowProvider } from '@/contexts/WorkflowContext';
import { WorkflowEditor } from '@/components/workflows/WorkflowEditor';

export const metadata: Metadata = {
  title: 'Edit Workflow | AI Call Center',
  description: 'Edit workflow automation',
};

interface PageProps {
  params: {
    id: string;
  };
}

export default function EditWorkflowPage({ params }: PageProps) {
  return (
    <WorkflowProvider>
      <WorkflowEditor workflowId={params.id} />
    </WorkflowProvider>
  );
}
