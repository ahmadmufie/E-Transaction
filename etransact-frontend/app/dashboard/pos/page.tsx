'use client';
import { useState, useEffect } from 'react';
import Script from 'next/script'; // Import wajib untuk Midtrans
import api from '../../../lib/api';
import { useAuthStore } from '../../../store/useAuthStore';
import Swal from 'sweetalert2';

export default function POSPage() {
  const { user } = useAuthStore();
  const [equipments, setEquipments] = useState<any[]>([]);
  const [cart, setCart] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [isProcessing, setIsProcessing] = useState(false);
  const [payMethod, setPayMethod] = useState('CASH');

  useEffect(() => { fetchEquipments(); }, []);

  const fetchEquipments = async () => {
    try {
      const res = await api.get('/equipments');
      const dataArray = res.data.data || res.data || [];
      setEquipments(Array.isArray(dataArray) ? dataArray : []);
    } catch (error) { console.error("Gagal memuat armada", error); } 
    finally { setLoading(false); }
  };

  const getSafePrice = (item: any) => Number(item.retail_price || item.base_price || 0);

  const addToCart = (item: any) => {
    const existing = cart.find(c => c.id === item.id);
    if (existing) setCart(cart.map(c => c.id === item.id ? { ...c, quantity: c.quantity + 1 } : c));
    else setCart([...cart, { ...item, quantity: 1 }]);
  };

  const removeFromCart = (id: string) => setCart(cart.filter(c => c.id !== id));
  const totalHarga = cart.reduce((total, item) => total + (getSafePrice(item) * item.quantity), 0);
const handleCheckout = async () => {
    if (cart.length === 0) return Swal.fire('Oops!', 'Keranjang masih kosong!', 'warning');
    
    const authUser = user as any;
    const validCashierId = authUser?.id || authUser?.ID;
    if (!validCashierId) return Swal.fire('Sesi Habis', 'Silakan logout dan login kembali.', 'error');

    setIsProcessing(true);
    try {
      const payload = {
        cashier_id: validCashierId,
        pay_method: payMethod,
        items: cart.map(c => ({ equipment_id: c.id, quantity: c.quantity }))
      };
      
      const res = await api.post('/pos/checkout', payload);

      // --- LOGIKA MIDTRANS SNAP ---
      if (res.data.snap_token && (payMethod === 'TRANSFER' || payMethod === 'EWALLET')) {
        if (!(window as any).snap) {
          Swal.fire('Mohon Tunggu', 'Sistem pembayaran sedang memuat...', 'info');
          setIsProcessing(false);
          return;
        }

        (window as any).snap.pay(res.data.snap_token, {
          onSuccess: function() {
            Swal.fire('Sukses!', 'Pembayaran Digital Berhasil!', 'success');
            setCart([]); fetchEquipments();
          },
          onPending: function() {
            Swal.fire('Tertunda', 'Menunggu pembayaran Anda diselesaikan...', 'info');
          },
          onError: function() {
            Swal.fire('Gagal!', 'Pembayaran Dibatalkan atau Gagal.', 'error');
          },
          onClose: function() {
            Swal.fire('Perhatian', 'Anda menutup jendela pembayaran sebelum selesai.', 'warning');
          }
        });
      } else {
      // Logika Transaksi Tunai (CASH)
      Swal.fire({
        title: 'Transaksi Sukses!',
        html: `Pembayaran <b>Tunai (CASH)</b> berhasil dicatat.<br/><br/>ID: <span class="text-sm text-slate-500">${res.data.transaction_id}</span>`,
        icon: 'success',
        showCancelButton: true,
        confirmButtonText: 'Tutup Saja',
        cancelButtonText: '🖨️ Cetak Kontrak (PDF)',
        cancelButtonColor: '#10b981', // Warna Hijau Emerald
        confirmButtonColor: '#64748b' // Warna Abu-abu
      }).then((result) => {
        // Jika kasir menekan tombol Cetak Kontrak (Cancel Button)
        if (result.dismiss === Swal.DismissReason.cancel) {
          downloadContractPDF({
            transaction_id: res.data.transaction_id,
            pay_method: 'TUNAI (CASH)',
            total_amount: totalHarga
          });
        }
      });
      setCart([]); fetchEquipments();
        }
    } catch (error: any) {
      Swal.fire('Error!', error.response?.data?.error || "Gagal memproses transaksi", 'error');
    } finally {
      setIsProcessing(false);
    }
  };
  
  const downloadContractPDF = async (transactionData: any) => {
      Swal.fire({ title: 'Merakit Kontrak...', html: 'Mohon tunggu, Puppeteer sedang bekerja', allowOutsideClick: false, didOpen: () => { Swal.showLoading(); } });
      
      try {
        const response = await fetch('/api/generate-contract', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(transactionData),
        });

        if (!response.ok) throw new Error("Gagal generate PDF");

        // Mengubah response menjadi format file Blob (Binary)
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `Kontrak-${transactionData.transaction_id}.pdf`;
        document.body.appendChild(a);
        a.click();
        a.remove();
        
        Swal.fire('Berhasil!', 'Kontrak PDF berhasil diunduh.', 'success');
      } catch (error) {
        Swal.fire('Error', 'Gagal mencetak dokumen.', 'error');
      }
    };
  
  if (loading) return <div className="p-10 font-bold text-slate-600">Memuat data alat berat...</div>;

  return (
    <>
      {/* --- INJECT SCRIPT MIDTRANS DI SINI --- */}
      {/* WAJIB GANTI data-client-key DENGAN CLIENT KEY ANDA */}
      <Script 
        src="https://app.sandbox.midtrans.com/snap/snap.js" 
        data-client-key="SB-Mid-client-XXXXX" 
        strategy="lazyOnload"
      />

      <div className="flex flex-col lg:flex-row gap-6 h-[calc(100vh-140px)]">
        {/* KIRI: Katalog Alat Berat */}
        <div className="flex-1 bg-white p-6 rounded-xl shadow-sm border border-slate-200 overflow-y-auto">
          <h2 className="text-xl font-bold mb-4 text-slate-800">Katalog Armada & Layanan</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
            {equipments.map((eq, idx) => {
              const safeId = eq.id || eq.ID || idx;
              const isOutOfStock = eq.stock <= 0;
              return (
                <div key={safeId} className={`border p-4 rounded-lg transition ${isOutOfStock ? 'bg-red-50 border-red-200 opacity-60' : 'bg-slate-50 border-slate-200 hover:border-blue-400'}`}>
                  <h3 className="font-bold text-slate-800">{eq.name || eq.Name || 'Tanpa Nama'}</h3>
                  <p className={`text-sm mb-2 font-medium ${isOutOfStock ? 'text-red-500' : 'text-slate-500'}`}>Stok: {eq.stock || eq.Stock || 0} Unit</p>
                  <p className="font-semibold text-blue-600 mb-4">Rp {getSafePrice(eq).toLocaleString('id-ID')}</p>
                  <button onClick={() => addToCart({ ...eq, id: safeId })} disabled={isOutOfStock} className="w-full bg-slate-800 hover:bg-slate-700 disabled:bg-slate-400 text-white py-2 rounded text-sm font-medium transition disabled:cursor-not-allowed">
                    {isOutOfStock ? 'Stok Habis' : '+ Tambah'}
                  </button>
                </div>
              );
            })}
          </div>
        </div>

        {/* KANAN: Keranjang Kasir */}
        <div className="w-full lg:w-[400px] bg-white p-6 rounded-xl shadow-sm border border-slate-200 flex flex-col">
          <h2 className="text-xl font-bold mb-4 text-slate-800">Detail Pembayaran</h2>
          <div className="mb-4">
            <label className="text-sm text-slate-600 font-medium block mb-2">Metode Pembayaran:</label>
            <select className="w-full p-2 border border-slate-300 rounded focus:ring-1 focus:ring-blue-500 bg-slate-50 text-sm" value={payMethod} onChange={(e) => setPayMethod(e.target.value)}>
              <option value="CASH">Tunai (CASH)</option>
              <option value="TRANSFER">Transfer Bank (VA)</option>
              <option value="EWALLET">Dompet Digital (QRIS/Gopay)</option>
            </select>
          </div>
          <div className="flex-1 overflow-y-auto mb-4 border-t border-b border-slate-100 py-4 space-y-3">
            {cart.map((c, idx) => (
              <div key={idx} className="flex justify-between items-center text-sm border-b border-slate-50 pb-2">
                <div><p className="font-semibold text-slate-800">{c.name}</p><p className="text-slate-500">{c.quantity}x @ Rp {getSafePrice(c).toLocaleString('id-ID')}</p></div>
                <button onClick={() => removeFromCart(c.id)} className="text-red-500 font-bold p-2 hover:bg-red-50 rounded">✕</button>
              </div>
            ))}
          </div>
          <div className="pt-2">
            <div className="flex justify-between items-end mb-4"><span className="text-slate-500 font-medium">Total Akhir</span><span className="text-2xl font-bold text-blue-600">Rp {totalHarga.toLocaleString('id-ID')}</span></div>
            <button onClick={handleCheckout} disabled={isProcessing || cart.length === 0} className="w-full bg-green-600 hover:bg-green-700 disabled:bg-slate-300 text-white font-bold py-3 rounded-lg transition disabled:cursor-not-allowed">
              {isProcessing ? 'Memproses Gateway...' : 'Bayar Sekarang'}
            </button>
          </div>
        </div>
      </div>
    </>
  );
}