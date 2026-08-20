import React, { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import { Camera, Image, Printer, Mail, DollarSign, TrendingUp, Users, Monitor } from 'lucide-react'

interface AnalyticsData {
  sessions: number
  photos: number
  prints: number
  emails: number
  revenue: number
  merchants?: number
  active_devices?: number
}

export default function OverviewPage() {
  const [data, setData] = useState<AnalyticsData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get('/analytics/overview').then((res) => {
      setData(res.data.data)
      setLoading(false)
    })
  }, [])

  if (loading) return <p className="text-gray-500">Loading...</p>

  const stats = [
    { label: 'Total Sessions', value: data?.sessions || 0, icon: Camera, color: 'bg-pink-100 text-pink-600' },
    { label: 'Photos Generated', value: data?.photos || 0, icon: Image, color: 'bg-purple-100 text-purple-600' },
    { label: 'Prints Completed', value: data?.prints || 0, icon: Printer, color: 'bg-blue-100 text-blue-600' },
    { label: 'Emails Sent', value: data?.emails || 0, icon: Mail, color: 'bg-green-100 text-green-600' },
    { label: 'Revenue', value: `Rp ${(data?.revenue || 0).toLocaleString()}`, icon: DollarSign, color: 'bg-yellow-100 text-yellow-600' },
  ]

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Overview</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {stats.map((stat) => (
          <div key={stat.label} className="card">
            <div className="flex items-center gap-4">
              <div className={`w-12 h-12 rounded-2xl flex items-center justify-center ${stat.color}`}>
                <stat.icon size={24} />
              </div>
              <div>
                <p className="text-sm text-gray-500">{stat.label}</p>
                <p className="text-2xl font-bold">{stat.value}</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
