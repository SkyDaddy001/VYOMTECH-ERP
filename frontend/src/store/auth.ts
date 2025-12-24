import { create } from 'zustand';

interface User {
  user_id: string;
  email: string;
  tenant_id: string;
  role: string;
}

interface AuthStore {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  setUser: (user: User) => void;
  setToken: (token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  token: typeof window !== 'undefined' ? localStorage.getItem('access_token') : null,
  isAuthenticated: typeof window !== 'undefined' ? !!localStorage.getItem('access_token') : false,

  setUser: (user) => set({ user, isAuthenticated: true }),

  setToken: (token) => {
    localStorage.setItem('access_token', token);
    set({ token, isAuthenticated: true });
  },

  logout: () => {
    localStorage.removeItem('access_token');
    set({ user: null, token: null, isAuthenticated: false });
  },
}));
