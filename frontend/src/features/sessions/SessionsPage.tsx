import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../../lib/api'
import { Eye, Trash2, Download } from 'lucide-react'

interface Session {
  id: string
  status: string
  payment_status: string
  email: string
  created_at: string
  final_image_url: string
}

export default function SessionsPage() {
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    api.get('/sessions').then((res) => {
      setSessions(res.data.data || [])
      setLoading(false)
    })
  }, [])

  const handleDelete = async (id: string) => {
    if (confirm('Delete this session?')) {
      await api.delete(`/sessions/${id}`)
      setSessions(sessions.filter((s) => s.id !== id))
    }
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Sessions</h1>

      {loading ? (
        <p className="text-gray-500">Loading...</p>
      ) : sessions.length === 0 ? (
        <p className="text-gray-500">No sessions yet.</p>
      ) : (
        <div className="space-y-3">
          {sessions.map((s) => (
            <div key={s.id} className="card flex items-center justify-between">
              <div className="flex items-center gap-4">
                <div className="w-16 h-16 bg-gray-100 rounded-xl flex items-center justify-center">
                  {s.final_image_url ? (
                    <img src={s.final_image_url} alt="" className="w-full h-full object-cover rounded-xl" />
                  ) : (
                    <span className="text-2xl">📷</span>
                  )}
                </div>
                <div>
                  <p className="font-medium">{s.id.slice(0, 8)}...</p>
                  <p className="text-sm text-gray-500">{new Date(s.created_at).toLocaleString()}</p>
                  <div className="flex gap-2 mt-1">
                    <span className={`text-xs px-2 py-0.5 rounded-full ${s.status === 'COMPLETED' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'}`}>
                      {s.status}
                    </span>
                    <span className={`text-xs px-2 py-0.5 rounded-full ${s.payment_status === 'PAID' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'}`}>
                      {s.payment_status}
                    </span>
                  </div>
                </div>
              </div>
              <div className="flex gap-2">
                <button onClick={() => navigate(`/dashboard/sessions/${s.id}`)} className="p-2 hover:bg-gray-100 rounded-lg">
                  <Eye size={18} />
                </button>
                <button onClick={() => handleDelete(s.id)} className="p-2 text-red-500 hover:bg-red-50 rounded-lg">
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
