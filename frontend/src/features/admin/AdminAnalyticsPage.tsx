import React, { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import { Users, Camera, Image, Printer, Mail, DollarSign } from 'lucide-react'

interface Analytics {
  merchants: number
  active_devices: number
  total_sessions: number
  photos: number
  prints: number
  emails: number
  revenue: number
}

export default function AdminAnalyticsPage() {
  const [data, setData] = useState<Analytics | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get('/admin/analytics').then((res) => {
      setData(res.data.data)
      setLoading(false)
    })
  }, [])

  if (loading) return <p className="text-text-subtle">Loading...</p>

  const stats = [
    { label: 'Merchants', value: data?.merchants || 0, icon: Users, color: 'bg-pink-500' },
    { label: 'Active Devices', value: data?.active_devices || 0, icon: Camera, color: 'bg-purple-500' },
    { label: 'Sessions', value: data?.total_sessions || 0, icon: Camera, color: 'bg-blue-500' },
    { label: 'Photos', value: data?.photos || 0, icon: Image, color: 'bg-green-500' },
    { label: 'Prints', value: data?.prints || 0, icon: Printer, color: 'bg-yellow-500' },
    { label: 'Emails', value: data?.emails || 0, icon: Mail, color: 'bg-indigo-500' },
    { label: 'Revenue', value: `Rp ${(data?.revenue || 0).toLocaleString()}`, icon: DollarSign, color: 'bg-emerald-500' },
  ]

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Platform Analytics</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat) => {
          const Icon = stat.icon
          return (
            <div key={stat.label} className="bg-surface rounded-lg p-6">
              <div className={`w-12 h-12 ${stat.color} rounded-lg flex items-center justify-center mb-4`}>
                <Icon className="text-text" size={24} />
              </div>
              <p className="text-2xl font-bold">{stat.value}</p>
              <p className="text-sm text-text-subtle">{stat.label}</p>
            </div>
          )
        })}
      </div>
    </div>
  )
}
