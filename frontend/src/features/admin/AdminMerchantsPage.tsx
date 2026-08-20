import React, { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import { Plus, Trash2, Edit2 } from 'lucide-react'

interface Merchant {
  id: string
  business_name: string
  slug: string
  subscription_status: string
  created_at: string
}

export default function AdminMerchantsPage() {
  const [merchants, setMerchants] = useState<Merchant[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ business_name: '', admin_email: '', admin_password: '', admin_name: '' })

  const fetchMerchants = () => {
    api.get('/admin/merchants').then((res) => {
      setMerchants(res.data.data || [])
      setLoading(false)
    })
  }

  useEffect(() => { fetchMerchants() }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await api.post('/admin/merchants', form)
    setShowForm(false)
    setForm({ business_name: '', admin_email: '', admin_password: '', admin_name: '' })
    fetchMerchants()
  }

  const handleDelete = async (id: string) => {
    if (confirm('Delete this merchant? This cannot be undone.')) {
      await api.delete(`/admin/merchants/${id}`)
      fetchMerchants()
    }
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Merchants</h1>
        <button onClick={() => setShowForm(true)} className="bg-primary text-white px-4 py-2 rounded-xl flex items-center gap-2">
          <Plus size={20} /> Add Merchant
        </button>
      </div>

      {showForm && (
        <div className="bg-gray-800 rounded-2xl p-6 mb-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <input placeholder="Business Name" className="w-full px-4 py-3 rounded-xl bg-gray-700 text-white border border-gray-600" value={form.business_name} onChange={(e) => setForm({ ...form, business_name: e.target.value })} required />
            <input placeholder="Admin Name" className="w-full px-4 py-3 rounded-xl bg-gray-700 text-white border border-gray-600" value={form.admin_name} onChange={(e) => setForm({ ...form, admin_name: e.target.value })} />
            <input type="email" placeholder="Admin Email" className="w-full px-4 py-3 rounded-xl bg-gray-700 text-white border border-gray-600" value={form.admin_email} onChange={(e) => setForm({ ...form, admin_email: e.target.value })} required />
            <input type="password" placeholder="Admin Password" className="w-full px-4 py-3 rounded-xl bg-gray-700 text-white border border-gray-600" value={form.admin_password} onChange={(e) => setForm({ ...form, admin_password: e.target.value })} required />
            <div className="flex gap-3">
              <button type="submit" className="bg-primary text-white px-6 py-2 rounded-xl">Create</button>
              <button type="button" onClick={() => setShowForm(false)} className="bg-gray-700 text-white px-6 py-2 rounded-xl">Cancel</button>
            </div>
          </form>
        </div>
      )}

      {loading ? (
        <p className="text-gray-400">Loading...</p>
      ) : (
        <div className="space-y-3">
          {merchants.map((m) => (
            <div key={m.id} className="bg-gray-800 rounded-2xl p-4 flex items-center justify-between">
              <div>
                <h3 className="font-semibold text-lg">{m.business_name}</h3>
                <p className="text-sm text-gray-400">{m.slug}</p>
              </div>
              <div className="flex items-center gap-2">
                <span className={`px-3 py-1 rounded-full text-xs font-medium ${m.subscription_status === 'active' ? 'bg-green-900 text-green-300' : 'bg-red-900 text-red-300'}`}>
                  {m.subscription_status}
                </span>
                <button onClick={() => handleDelete(m.id)} className="p-2 text-red-400 hover:bg-red-900/30 rounded-lg">
                  <Trash2 size={18} />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
