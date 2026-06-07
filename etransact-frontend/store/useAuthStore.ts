import { create } from 'zustand';
import { persist } from 'zustand/middleware';

// 1. Tambahkan 'id' di cetakan (Interface)
interface User {
  id: string;  // <-- INI KUNCI UTAMANYA!
  name: string;
  role: string;
}

interface AuthState {
  user: User | null;
  setUser: (user: User | null) => void;
  logout: () => void;
}

// 2. Gunakan 'persist' agar kebal terhadap Refresh halaman
export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      setUser: (user) => set({ user }),
      logout: () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user'); // Bersihkan sisa manual jika ada
        set({ user: null });
      },
    }),
    {
      name: 'auth-storage', // Nama brankas penyimpanan rahasia di browser
    }
  )
);