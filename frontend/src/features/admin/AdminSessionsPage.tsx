import React, { useEffect, useState } from 'react'
import { api } from '../../lib/api'

interface Session {
  id: string
  merchant_id: string
  status: string
  payment_status: string
  email: string
  created_at: string
  business_name: string
}

export default function AdminSessionsPage() {
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get('/admin/sessions').then((res) => {
      setSessions(res.data.data || [])
      setLoading(false)
    })
  }, [])

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">All Sessions</h1>

      {loading ? (
        <p className="text-gray-400">Loading...</p>
      ) : (
        <div className="space-y-3">
          {sessions.map((s) => (
            <div key={s.id} className="bg-gray-800 rounded-2xl p-4 flex items-center justify-between">
              <div>
                <p className="font-medium">{s.id.slice(0, 8)}... <span className="text-gray-400">by {s.business_name}</span></p>
                <p className="text-sm text-gray-400">{new Date(s.created_at).toLocaleString()}</p>
              </div>
              <div className="flex gap-2">
                <span className={`px-3 py-1 rounded-full text-xs font-medium ${s.status === 'COMPLETED' ? 'bg-green-900 text-green-300' : 'bg-yellow-900 text-yellow-300'}`}>
                  {s.status}
                </span>
                <span className={`px-3 py-1 rounded-full text-xs font-medium ${s.payment_status === 'PAID' ? 'bg-green-900 text-green-300' : 'bg-gray-700 text-gray-300'}`}>
                  {s.payment_status}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
