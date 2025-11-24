// Create Workflow Page - Phase 3B

import { Metadata } from 'next';
import { WorkflowProvider } from '@/contexts/WorkflowContext';
import { WorkflowEditor } from '@/components/workflows/WorkflowEditor';

export const metadata: Metadata = {
  title: 'Create Workflow | AI Call Center',
  description: 'Create a new workflow automation',
};

export default function CreateWorkflowPage() {
  return (
    <WorkflowProvider>
      <WorkflowEditor />
    </WorkflowProvider>
  );
}
