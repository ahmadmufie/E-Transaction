'use client';
import { useState, useEffect } from 'react';
import api from '../../../lib/api';
import { useAuthStore } from '../../../store/useAuthStore';

export default function EquipmentsPage() {
  const { user } = useAuthStore();
  const [equipments, setEquipments] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingItem, setEditingItem] = useState<any>(null);

  // Form State
  const [formData, setFormData] = useState({
    name: '',
    type: 'HEAVY_MACHINERY',
    base_price: 0,
    retail_price: 0,
    stock: 0
  });

  useEffect(() => { fetchEquipments(); }, []);

  const fetchEquipments = async () => {
    try {
      const res = await api.get('/equipments');
      const data = res.data.data || res.data || [];
      setEquipments(Array.isArray(data) ? data : []);
    } catch (error) { console.error("Gagal load data", error); }
    finally { setLoading(false); }
  };

  const handleOpenModal = (item: any = null) => {
    setEditingItem(item);
    if (item) {
      setFormData({
        name: item.name,
        type: item.type,
        base_price: item.base_price,
        retail_price: item.retail_price,
        stock: item.stock
      });
    } else {
      setFormData({ name: '', type: 'HEAVY_MACHINERY', base_price: 0, retail_price: 0, stock: 0 });
    }
    setShowModal(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingItem) {
        await api.put(`/equipments/${editingItem.id}`, formData);
      } else {
        await api.post('/equipments', formData);
      }
      setShowModal(false);
      fetchEquipments();
    } catch (error) { alert("Gagal menyimpan data"); }
  };

  const handleDelete = async (id: string) => {
    if (confirm("Hapus alat berat ini dari database?")) {
      try {
        await api.delete(`/equipments/${id}`);
        fetchEquipments();
      } catch (error) { alert("Gagal menghapus"); }
    }
  };

  if (loading) return <div className="p-10 font-bold">Memuat armada...</div>;

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-slate-800">Manajemen Armada & Stok</h2>
        <button 
          onClick={() => handleOpenModal()}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-bold transition shadow-lg"
        >
          + Tambah
        </button>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
        <table className="w-full text-left">
          <thead className="bg-slate-50 border-b border-slate-200">
            <tr>
              <th className="p-4 font-bold text-slate-600">Nama Armada/Alat/Service</th>
              <th className="p-4 font-bold text-slate-600">Kategori</th>
              <th className="p-4 font-bold text-slate-600">Harga Modal</th>
              <th className="p-4 font-bold text-slate-600">Harga Retail</th>
              <th className="p-4 font-bold text-slate-600">Stok</th>
              <th className="p-4 font-bold text-slate-600 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {equipments.map((eq) => (
              <tr key={eq.id} className="hover:bg-slate-50 transition">
                <td className="p-4 font-semibold text-slate-800">{eq.name}</td>
                <td className="p-4"><span className="bg-slate-100 text-slate-600 px-3 py-1 rounded-full text-xs font-bold">{eq.type}</span></td>
                <td className="p-4 font-bold text-blue-600">Rp {Number(eq.base_price).toLocaleString()}</td>
                <td className="p-4 font-bold text-blue-600">Rp {Number(eq.retail_price).toLocaleString()}</td>
                <td className={`p-4 font-bold ${eq.stock < 2 ? 'text-red-500' : 'text-green-600'}`}>{eq.stock} Unit</td>
                <td className="p-4 text-right space-x-2">
                  <button onClick={() => handleOpenModal(eq)} className="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-lg font-bold transition shadow-lg">Edit</button>
                  <button onClick={() => handleDelete(eq.id)} className="bg-red-600 hover:bg-red-700 text-white px-6 py-2 rounded-lg font-bold transition shadow-lg">Hapus</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* MODAL CRUD */}
      {showModal && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md p-8">
            <h3 className="text-xl font-bold mb-6 text-slate-800">{editingItem ? 'Edit Armada' : 'Tambah Armada Baru'}</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-bold text-slate-600 mb-1">Nama Armada</label>
                <input type="text" required className="w-full p-2 border rounded-lg" value={formData.name} onChange={(e)=>setFormData({...formData, name: e.target.value})} />
              </div>
              <div>
                <label className="block text-sm font-bold text-slate-600 mb-1">Kategori</label>
                <select 
                  required 
                  className="w-full p-2 border rounded-lg bg-white"
                  value={formData.type} 
                  onChange={(e)=>setFormData({...formData, type: e.target.value})}
                >
                  <option value="ARMADA">Alat Berat (Heavy Machinery)</option>
                  <option value="SPAREPART">Suku Cadang (Spare Part)</option>
                  <option value="SERVICE">Layanan (Service)</option>
                </select>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-bold text-slate-600 mb-1">Harga Modal</label>
                  <input type="number" required className="w-full p-2 border rounded-lg" value={formData.base_price} onChange={(e)=>setFormData({...formData, base_price: Number(e.target.value)})} />
                </div>
                <div>
                  <label className="block text-sm font-bold text-slate-600 mb-1">Harga Retail</label>
                  <input type="number" required className="w-full p-2 border rounded-lg" value={formData.retail_price} onChange={(e)=>setFormData({...formData, retail_price: Number(e.target.value)})} />
                </div>
              </div>
              <div>
                <label className="block text-sm font-bold text-slate-600 mb-1">Jumlah Stok</label>
                <input type="number" required className="w-full p-2 border rounded-lg" value={formData.stock} onChange={(e)=>setFormData({...formData, stock: Number(e.target.value)})} />
              </div>
              <div className="flex gap-4 pt-4">
                <button type="button" onClick={() => setShowModal(false)} className="flex-1 py-3 font-bold text-slate-500 hover:bg-slate-50 rounded-xl">Batal</button>
                <button type="submit" className="flex-1 py-3 bg-blue-600 text-white font-bold rounded-xl shadow-lg hover:bg-blue-700 transition">Simpan</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}