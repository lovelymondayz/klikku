import React, { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import { Plus, Edit2, Trash2 } from 'lucide-react'

interface Campaign {
  id: string
  name: string
  description: string
  start_date: string
  end_date: string
  status: string
}

export default function CampaignsPage() {
  const [campaigns, setCampaigns] = useState<Campaign[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ name: '', description: '', start_date: '', end_date: '', status: 'draft' })

  const fetchCampaigns = () => {
    api.get('/campaigns').then((res) => {
      setCampaigns(res.data.data || [])
      setLoading(false)
    })
  }

  useEffect(() => { fetchCampaigns() }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await api.post('/campaigns', form)
    setShowForm(false)
    setForm({ name: '', description: '', start_date: '', end_date: '', status: 'draft' })
    fetchCampaigns()
  }

  const handleDelete = async (id: string) => {
    if (confirm('Delete this campaign?')) {
      await api.delete(`/campaigns/${id}`)
      fetchCampaigns()
    }
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Campaigns</h1>
        <button onClick={() => setShowForm(true)} className="btn-primary flex items-center gap-2">
          <Plus size={20} /> New Campaign
        </button>
      </div>

      {showForm && (
        <div className="card mb-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <input placeholder="Campaign Name" className="input" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
            <input placeholder="Description" className="input" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
            <div className="grid grid-cols-2 gap-4">
              <input type="date" className="input" value={form.start_date} onChange={(e) => setForm({ ...form, start_date: e.target.value })} required />
              <input type="date" className="input" value={form.end_date} onChange={(e) => setForm({ ...form, end_date: e.target.value })} required />
            </div>
            <select className="input" value={form.status} onChange={(e) => setForm({ ...form, status: e.target.value })}>
              <option value="draft">Draft</option>
              <option value="scheduled">Scheduled</option>
              <option value="active">Active</option>
              <option value="ended">Ended</option>
            </select>
            <div className="flex gap-3">
              <button type="submit" className="btn-primary">Create</button>
              <button type="button" onClick={() => setShowForm(false)} className="btn-secondary">Cancel</button>
            </div>
          </form>
        </div>
      )}

      {loading ? (
        <p className="text-gray-500">Loading...</p>
      ) : campaigns.length === 0 ? (
        <p className="text-gray-500">No campaigns yet. Create your first campaign!</p>
      ) : (
        <div className="grid gap-4">
          {campaigns.map((c) => (
            <div key={c.id} className="card flex items-center justify-between">
              <div>
                <h3 className="font-semibold text-lg">{c.name}</h3>
                <p className="text-sm text-gray-500">{c.description}</p>
                <p className="text-xs text-gray-400 mt-1">{c.start_date} → {c.end_date}</p>
              </div>
              <div className="flex items-center gap-2">
                <span className={`px-3 py-1 rounded-full text-xs font-medium ${c.status === 'active' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'}`}>
                  {c.status}
                </span>
                <button onClick={() => handleDelete(c.id)} className="p-2 text-red-500 hover:bg-red-50 rounded-lg">
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
