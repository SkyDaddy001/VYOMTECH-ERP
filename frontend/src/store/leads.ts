import { create } from 'zustand';

interface Lead {
  id: string;
  lead_code: string;
  lead_name: string;
  email?: string;
  phone?: string;
  company?: string;
  lead_source?: string;
  lead_value?: number;
  lead_status: string;
  description?: string;
  assigned_to?: string;
  created_at: string;
  updated_at: string;
}

interface LeadStore {
  leads: Lead[];
  selectedLead: Lead | null;
  isLoading: boolean;
  error: string | null;
  setLeads: (leads: Lead[]) => void;
  setSelectedLead: (lead: Lead | null) => void;
  setIsLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  addLead: (lead: Lead) => void;
  updateLead: (lead: Lead) => void;
  deleteLead: (id: string) => void;
}

export const useLeadStore = create<LeadStore>((set) => ({
  leads: [],
  selectedLead: null,
  isLoading: false,
  error: null,

  setLeads: (leads) => set({ leads }),
  setSelectedLead: (lead) => set({ selectedLead: lead }),
  setIsLoading: (isLoading) => set({ isLoading }),
  setError: (error) => set({ error }),

  addLead: (lead) => set((state) => ({ leads: [lead, ...state.leads] })),

  updateLead: (lead) =>
    set((state) => ({
      leads: state.leads.map((l) => (l.id === lead.id ? lead : l)),
      selectedLead: state.selectedLead?.id === lead.id ? lead : state.selectedLead,
    })),

  deleteLead: (id) =>
    set((state) => ({
      leads: state.leads.filter((l) => l.id !== id),
      selectedLead: state.selectedLead?.id === id ? null : state.selectedLead,
    })),
}));
