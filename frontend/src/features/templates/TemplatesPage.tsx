import React, { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import { Plus, Trash2, Paintbrush } from 'lucide-react'
import TemplateEditor from './TemplateEditor'

interface Template {
  id: string
  name: string
  photo_count: number
  price: number
  active: boolean
  output_width: number
  output_height: number
}

export default function TemplatesPage() {
  const [templates, setTemplates] = useState<Template[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingTemplateId, setEditingTemplateId] = useState<string | null>(null)
  const [form, setForm] = useState({ name: '', photo_count: 4, price: 0, output_width: 1200, output_height: 1800 })

  if (editingTemplateId !== null) {
    return (
      <div className="h-[calc(100vh-200px)]">
        <TemplateEditor
          templateId={editingTemplateId}
          onSave={() => { setEditingTemplateId(null); fetchTemplates() }}
        />
      </div>
    )
  }

  const fetchTemplates = () => {
    api.get('/templates').then((res) => {
      setTemplates(res.data.data || [])
      setLoading(false)
    })
  }

  useEffect(() => { fetchTemplates() }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await api.post('/templates', form)
    setShowForm(false)
    fetchTemplates()
  }

  const handleDelete = async (id: string) => {
    if (confirm('Delete this template?')) {
      await api.delete(`/templates/${id}`)
      fetchTemplates()
    }
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Templates</h1>
        <button onClick={() => setShowForm(true)} className="btn-primary flex items-center gap-2">
          <Plus size={20} /> New Template
        </button>
      </div>

      {showForm && (
        <div className="card mb-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <input placeholder="Template Name" className="input" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
            <div className="grid grid-cols-3 gap-4">
              <div>
                <label className="text-sm text-gray-600">Photo Count</label>
                <input type="number" min="1" max="10" className="input" value={form.photo_count} onChange={(e) => setForm({ ...form, photo_count: +e.target.value })} />
              </div>
              <div>
                <label className="text-sm text-gray-600">Price (Rp)</label>
                <input type="number" className="input" value={form.price} onChange={(e) => setForm({ ...form, price: +e.target.value })} />
              </div>
              <div>
                <label className="text-sm text-gray-600">Dimensions</label>
                <div className="flex gap-2">
                  <input type="number" className="input" value={form.output_width} onChange={(e) => setForm({ ...form, output_width: +e.target.value })} placeholder="W" />
                  <input type="number" className="input" value={form.output_height} onChange={(e) => setForm({ ...form, output_height: +e.target.value })} placeholder="H" />
                </div>
              </div>
            </div>
            <div className="flex gap-3">
              <button type="submit" className="btn-primary">Create</button>
              <button type="button" onClick={() => setShowForm(false)} className="btn-secondary">Cancel</button>
            </div>
          </form>
        </div>
      )}

      {loading ? (
        <p className="text-gray-500">Loading...</p>
      ) : templates.length === 0 ? (
        <p className="text-gray-500">No templates yet.</p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {templates.map((t) => (
            <div key={t.id} className="card">
              <div className="w-full h-40 bg-gradient-to-br from-pink-200 to-purple-200 rounded-2xl mb-4 flex items-center justify-center">
                <span className="text-4xl">📸</span>
              </div>
              <h3 className="font-semibold text-lg">{t.name}</h3>
              <p className="text-sm text-gray-500">{t.photo_count} photos • Rp {t.price.toLocaleString()}</p>
              <p className="text-xs text-gray-400">{t.output_width}×{t.output_height}px</p>
              <div className="flex justify-between mt-3">
                <button onClick={() => setEditingTemplateId(t.id)} className="btn-secondary flex items-center gap-1 text-sm">
                  <Paintbrush size={14} /> Edit
                </button>
                <button onClick={() => handleDelete(t.id)} className="p-2 text-red-500 hover:bg-red-50 rounded-lg">
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
