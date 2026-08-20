import React, { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import { Palette } from 'lucide-react'

interface BrandingData {
  business_name: string
  logo_url: string
  primary_color: string
  secondary_color: string
  font: string
  welcome_message: string
  idle_background_url: string
  email_design: Record<string, any>
  social_links: Record<string, any>
}

export default function BrandingPage() {
  const [branding, setBranding] = useState<BrandingData | null>(null)
  const [form, setForm] = useState<Partial<BrandingData>>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    api.get('/branding').then((res) => {
      const data = res.data.data as BrandingData
      setBranding(data)
      setForm(data)
      setLoading(false)
    })
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    await api.put('/branding', form)
    setSaving(false)
    setBranding(form as BrandingData)
  }

  if (loading) return <p className="text-gray-500">Loading...</p>

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Branding</h1>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="card">
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <Palette size={24} /> Visual Identity
          </h2>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Business Name</label>
              <input className="input" value={form.business_name || ''} onChange={(e) => setForm({ ...form, business_name: e.target.value })} />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Welcome Message</label>
              <input className="input" value={form.welcome_message || ''} onChange={(e) => setForm({ ...form, welcome_message: e.target.value })} placeholder="Capture Your Moment" />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Primary Color</label>
                <div className="flex gap-2">
                  <input type="color" className="w-12 h-10 rounded border" value={form.primary_color || '#ff6b9d'} onChange={(e) => setForm({ ...form, primary_color: e.target.value })} />
                  <input className="input" value={form.primary_color || ''} onChange={(e) => setForm({ ...form, primary_color: e.target.value })} />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Secondary Color</label>
                <div className="flex gap-2">
                  <input type="color" className="w-12 h-10 rounded border" value={form.secondary_color || '#c44dff'} onChange={(e) => setForm({ ...form, secondary_color: e.target.value })} />
                  <input className="input" value={form.secondary_color || ''} onChange={(e) => setForm({ ...form, secondary_color: e.target.value })} />
                </div>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Font</label>
              <select className="input" value={form.font || ''} onChange={(e) => setForm({ ...form, font: e.target.value })}>
                <option value="">Default</option>
                <option value="Poppins">Poppins</option>
                <option value="Inter">Inter</option>
                <option value="Playfair Display">Playfair Display</option>
                <option value="Montserrat">Montserrat</option>
              </select>
            </div>
          </div>
        </div>

        <div className="card">
          <h2 className="text-xl font-semibold mb-4">Logo</h2>
          <div className="flex items-center gap-4">
            {form.logo_url && (
              <img src={form.logo_url} alt="Logo" className="w-20 h-20 rounded-xl object-cover" />
            )}
            <div>
              <p className="text-sm text-gray-500 mb-2">Upload your logo to MinIO storage</p>
              <p className="text-xs text-gray-400">Use the Assets section to upload files</p>
            </div>
          </div>
        </div>

        <button type="submit" disabled={saving} className="btn-primary">
          {saving ? 'Saving...' : 'Save Changes'}
        </button>
      </form>
    </div>
  )
}
