import React, { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import { Plus, Trash2, Wifi, WifiOff, Settings, Camera } from 'lucide-react'

interface Device {
  id: string
  name: string
  device_token: string
  status: string
  current_campaign_id: string
  last_seen_at: string
  printer_config: any
}

export default function DevicesPage() {
  const [devices, setDevices] = useState<Device[]>([])
  const [campaigns, setCampaigns] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ name: '', device_token: '', current_campaign_id: '' })

  const fetchDevices = () => {
    api.get('/devices').then((res) => {
      setDevices(res.data.data || [])
      setLoading(false)
    })
  }

  const fetchCampaigns = () => {
    api.get('/campaigns').then((res) => {
      setCampaigns(res.data.data || [])
    })
  }

  useEffect(() => { fetchDevices(); fetchCampaigns() }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await api.post('/devices', { ...form, printer_config: {} })
    setShowForm(false)
    setForm({ name: '', device_token: '', current_campaign_id: '' })
    fetchDevices()
  }

  const handleDelete = async (id: string) => {
    if (confirm('Remove this device?')) {
      await api.delete(`/devices/${id}`)
      fetchDevices()
    }
  }

  const handleAssignCampaign = async (deviceId: string, campaignId: string) => {
    await api.put(`/devices/${devices.find(d => d.id === deviceId)?.id}`, { current_campaign_id: campaignId })
    fetchDevices()
  }

  const getCampaignName = (id: string) => {
    const campaign = campaigns.find(c => c.id === id)
    return campaign?.name || 'None assigned'
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold">Devices</h1>
          <p className="text-gray-500 mt-1">Manage photobooth devices and their campaigns</p>
        </div>
        <button onClick={() => setShowForm(true)} className="btn-primary flex items-center gap-2">
          <Plus size={20} /> Add Device
        </button>
      </div>

      {showForm && (
        <div className="card mb-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <input placeholder="Device Name (e.g., Main Store iPad)" className="input" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
            <input placeholder="Device Token (leave empty to auto-generate)" className="input" value={form.device_token} onChange={(e) => setForm({ ...form, device_token: e.target.value })} />
            <select className="input" value={form.current_campaign_id} onChange={(e) => setForm({ ...form, current_campaign_id: e.target.value })}>
              <option value="">Select Campaign (optional)</option>
              {campaigns.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
            </select>
            <div className="flex gap-3">
              <button type="submit" className="btn-primary">Add</button>
              <button type="button" onClick={() => setShowForm(false)} className="btn-secondary">Cancel</button>
            </div>
          </form>
        </div>
      )}

      {loading ? (
        <p className="text-gray-500">Loading...</p>
      ) : devices.length === 0 ? (
        <div className="card text-center py-12">
          <Camera size={48} className="mx-auto text-gray-300 mb-4" />
          <p className="text-gray-500">No devices registered. Add your first device to get started.</p>
        </div>
      ) : (
        <div className="grid gap-4">
          {devices.map((d) => (
            <div key={d.id} className="card">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className={`w-12 h-12 rounded-2xl flex items-center justify-center ${d.status === 'online' ? 'bg-green-100' : 'bg-gray-100'}`}>
                    {d.status === 'online' ? <Wifi className="text-green-600" size={24} /> : <WifiOff className="text-gray-400" size={24} />}
                  </div>
                  <div>
                    <h3 className="font-semibold">{d.name}</h3>
                    <p className="text-sm text-gray-500 font-mono">{d.device_token.slice(0, 16)}...</p>
                    <p className="text-xs text-gray-400">Last seen: {d.last_seen_at ? new Date(d.last_seen_at).toLocaleString() : 'Never'}</p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <span className={`px-3 py-1 rounded-full text-xs font-medium ${d.status === 'online' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'}`}>
                    {d.status}
                  </span>
                  <button onClick={() => handleDelete(d.id)} className="p-2 text-red-500 hover:bg-red-50 rounded-lg">
                    <Trash2 size={18} />
                  </button>
                </div>
              </div>
              <div className="mt-4 pt-4 border-t border-gray-100 flex items-center gap-3">
                <Settings size={16} className="text-gray-400" />
                <span className="text-sm text-gray-500">Campaign:</span>
                <select
                  className="text-sm border border-gray-200 rounded-lg px-3 py-1"
                  value={d.current_campaign_id || ''}
                  onChange={(e) => handleAssignCampaign(d.id, e.target.value)}
                >
                  <option value="">None assigned</option>
                  {campaigns.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                </select>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
