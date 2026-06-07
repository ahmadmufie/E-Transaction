'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import api from '../../lib/api';
import { useAuthStore } from '../../store/useAuthStore';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');
  
  const router = useRouter();
  const setUser = useAuthStore((state) => state.setUser);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setErrorMsg('');

    try {
      const response = await api.post('/login', { email, password });
      
      const { token, user } = response.data;
      
      // Simpan Token untuk otorisasi API
      localStorage.setItem('token', token);
      
      // Pastikan backend mengirimkan ID. Jika backend Golang memakai 'ID' (huruf besar), kita format menjadi 'id' huruf kecil.
      const formattedUser = {
        id: user.id || user.ID, 
        name: user.name || user.Name,
        role: user.role || user.Role
      };
      
      // Zustand akan otomatis menyimpannya ke memori dan LocalStorage!
      setUser(formattedUser);
      
      alert(`Selamat datang, ${formattedUser.name}!`);
      router.push('/dashboard');
      
    } catch (err: any) {
      setErrorMsg(err.response?.data?.error || 'Gagal terhubung ke server');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-900 text-white">
      <div className="bg-slate-800 p-8 rounded-xl shadow-2xl w-96 border border-slate-700">
        <h1 className="text-2xl font-bold mb-6 text-center text-blue-400">E-Transact Login</h1>
        
        {errorMsg && (
          <div className="bg-red-500/20 border border-red-500 text-red-300 p-3 rounded mb-4 text-sm text-center">
            {errorMsg}
          </div>
        )}

        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Email Perusahaan</label>
            <input 
              type="email" 
              className="w-full p-2.5 rounded bg-slate-700 border border-slate-600 text-white focus:ring-2 focus:ring-blue-500 outline-none"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Password</label>
            <input 
              type="password" 
              className="w-full p-2.5 rounded bg-slate-700 border border-slate-600 text-white focus:ring-2 focus:ring-blue-500 outline-none"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button 
            type="submit" 
            disabled={loading}
            className="w-full bg-blue-600 hover:bg-blue-700 font-bold py-2.5 rounded transition disabled:opacity-50"
          >
            {loading ? 'Memeriksa...' : 'Masuk ke Sistem'}
          </button>
        </form>
      </div>
    </div>
  );
}