'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link'; // <-- Menggunakan Link Next.js
import { useAuthStore } from '../../store/useAuthStore';

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const { user, setUser, logout } = useAuthStore();
  const [isClient, setIsClient] = useState(false);

  useEffect(() => {
    setIsClient(true);
    const token = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');
    
    if (!token) {
      router.push('/login');
    } else if (!user && savedUser) {
      // Jika halaman di-refresh, pulihkan data user dari localStorage ke Zustand
      setUser(JSON.parse(savedUser));
    }
  }, [router, user, setUser]);

  if (!isClient || !user) {
    return <div className="p-10 text-white bg-slate-900 min-h-screen">Memuat sistem...</div>;
  }

  const handleLogout = () => {
    logout();
    localStorage.removeItem('user'); // Bersihkan data user saat logout
    router.push('/login');
  };

  return (
    <div className="flex min-h-screen bg-slate-100 text-slate-800">
      {/* Sidebar Kiri */}
      <aside className="w-64 bg-slate-900 text-white p-6 shadow-xl flex flex-col justify-between">
        <div>
          <h2 className="text-2xl font-bold mb-8 text-blue-400">E-Transact</h2>
          
          <nav className="space-y-4">
            {/* GANTI <a> MENJADI <Link> */}
            <Link href="/dashboard" className="block p-2 rounded hover:bg-slate-800 transition">Beranda</Link>
            <Link href="/dashboard/pos" className="block p-2 rounded hover:bg-slate-800 transition">Mesin Kasir (POS)</Link>
            <Link href="/dashboard/equipments" className="block p-2 rounded hover:bg-slate-800 transition">Stok Alat Berat</Link>
            
            {user.role === 'SUPERADMIN' && (
              <Link href="/dashboard/reports" className="block p-2 rounded text-yellow-400 hover:bg-slate-800 transition">
                Laporan Keuangan
              </Link>
            )}
          </nav>
        </div>
        
        <div>
          <button onClick={handleLogout} className="text-red-400 hover:text-red-300 flex items-center font-semibold">
            Keluar Sistem
          </button>
        </div>
      </aside>

      {/* Konten Utama Kanan */}
      <main className="flex-1 p-8 overflow-y-auto">
        <header className="flex justify-between items-center mb-8 bg-white p-4 rounded-xl shadow-sm border border-slate-200">
          <h1 className="text-xl font-semibold">Selamat Datang, {user.name}</h1>
          <span className="bg-blue-100 text-blue-800 px-4 py-1.5 rounded-full text-sm font-bold tracking-wide">
            {user.role}
          </span>
        </header>
        
        {children}
      </main>
    </div>
  );
}