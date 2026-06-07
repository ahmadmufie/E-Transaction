'use client';
import { useState, useEffect } from 'react';
import api from '../../lib/api';
import { useAuthStore } from '../../store/useAuthStore';

export default function DashboardPage() {
  const { user } = useAuthStore();
  const [history, setHistory] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  // State untuk Ringkasan (Statistik)
  const [stats, setStats] = useState({
    totalRevenue: 0,
    totalTransactions: 0,
    pendingCount: 0,
    paidCount: 0
  });

  useEffect(() => {
    fetchHistory();
  }, []);

  const fetchHistory = async () => {
    try {
      const res = await api.get('/pos/history');
      const data = res.data.data || [];
      setHistory(data);

      // Kalkulasi Statistik Cepat
      let revenue = 0;
      let pending = 0;
      let paid = 0;

      data.forEach((trx: any) => {
        revenue += Number(trx.total);
        if (trx.status === 'PENDING') pending++;
        if (trx.status === 'PAID' || trx.status === 'LUNAS' || trx.status === 'SUCCESS') paid++;
      });

      setStats({
        totalRevenue: revenue,
        totalTransactions: data.length,
        pendingCount: pending,
        paidCount: paid
      });

    } catch (error) {
      console.error("Gagal mengambil data dashboard", error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div className="p-10 font-bold text-slate-600 animate-pulse">Menghidupkan Dashboard...</div>;

  return (
    <div className="space-y-6">
      {/* HEADER PESAN SAMBUTAN */}
      <div className="bg-gradient-to-r from-blue-900 to-blue-700 rounded-2xl p-8 text-white shadow-lg">
        <h1 className="text-3xl font-bold mb-2">Selamat Datang, {user?.name || 'Admin'}! 👋</h1>
        <p className="text-blue-200">Ini adalah ringkasan operasional PT Kontraktor Konstruksi hari ini.</p>
      </div>

      {/* KARTU STATISTIK (WIDGETS) */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200 border-l-4 border-l-blue-500">
          <p className="text-sm font-bold text-slate-500 mb-1">Total Pendapatan</p>
          <h2 className="text-2xl font-bold text-slate-800">Rp {stats.totalRevenue.toLocaleString('id-ID')}</h2>
        </div>
        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200 border-l-4 border-l-indigo-500">
          <p className="text-sm font-bold text-slate-500 mb-1">Total Transaksi</p>
          <h2 className="text-2xl font-bold text-slate-800">{stats.totalTransactions} Order</h2>
        </div>
        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200 border-l-4 border-l-amber-500">
          <p className="text-sm font-bold text-slate-500 mb-1">Menunggu Pembayaran (Pending)</p>
          <h2 className="text-2xl font-bold text-amber-600">{stats.pendingCount} Order</h2>
        </div>
        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200 border-l-4 border-l-green-500">
          <p className="text-sm font-bold text-slate-500 mb-1">Transaksi Lunas</p>
          <h2 className="text-2xl font-bold text-green-600">{stats.paidCount} Order</h2>
        </div>
      </div>

      {/* TABEL RIWAYAT TRANSAKSI */}
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
        <div className="p-6 border-b border-slate-100 flex justify-between items-center">
          <h2 className="text-xl font-bold text-slate-800">Riwayat Transaksi Terakhir</h2>
          <button onClick={fetchHistory} className="text-sm text-blue-600 font-bold hover:underline">🔄 Segarkan Data</button>
        </div>
        
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead className="bg-slate-50 border-b border-slate-200">
              <tr>
                <th className="p-4 font-bold text-slate-600 text-sm">ID Transaksi</th>
                <th className="p-4 font-bold text-slate-600 text-sm">Tanggal</th>
                <th className="p-4 font-bold text-slate-600 text-sm">Metode</th>
                <th className="p-4 font-bold text-slate-600 text-sm">Total (Rp)</th>
                <th className="p-4 font-bold text-slate-600 text-sm">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {history.map((trx, idx) => (
                <tr key={idx} className="hover:bg-slate-50 transition">
                  <td className="p-4 text-sm font-medium text-slate-800">
                    {trx.transaction_id.split('-')[0].toUpperCase()}...
                  </td>
                  <td className="p-4 text-sm text-slate-500">
                    {new Date(trx.date).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })}
                  </td>
                  <td className="p-4 text-sm font-bold text-slate-600">
                    {trx.method}
                  </td>
                  <td className="p-4 text-sm font-bold text-blue-600">
                    Rp {Number(trx.total).toLocaleString('id-ID')}
                  </td>
                  <td className="p-4 text-sm">
                    <span className={`px-3 py-1 rounded-full text-xs font-bold ${
                      trx.status === 'PENDING' ? 'bg-amber-100 text-amber-700' : 
                      trx.status === 'PAID' || trx.status === 'LUNAS' ? 'bg-green-100 text-green-700' : 
                      'bg-slate-100 text-slate-700'
                    }`}>
                      {trx.status}
                    </span>
                  </td>
                </tr>
              ))}
              {history.length === 0 && (
                <tr>
                  <td colSpan={5} className="p-8 text-center text-slate-400 italic">
                    Belum ada data transaksi. Silakan buat transaksi di menu Kasir (POS).
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}