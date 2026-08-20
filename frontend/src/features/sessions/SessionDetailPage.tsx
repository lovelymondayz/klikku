import React, { useEffect, useState } from 'react'
import { api } from '../../lib/api'
import { ArrowLeft, Download, RefreshCw, Mail, Printer, Trash2, CheckCircle, Clock, AlertCircle } from 'lucide-react'

interface Session {
  id: string
  status: string
  payment_status: string
  email: string
  final_image_url: string
  created_at: string
  completed_at: string
}

interface Photo {
  id: string
  original_url: string
  position: number
}

interface PrintJob {
  id: string
  status: string
  print_type: string
  copies: number
  printer_name: string
  error_message: string
  created_at: string
  printed_at: string
}

export default function SessionDetailPage() {
  const [session, setSession] = useState<Session | null>(null)
  const [photos, setPhotos] = useState<Photo[]>([])
  const [printJobs, setPrintJobs] = useState<PrintJob[]>([])
  const [loading, setLoading] = useState(true)
  const [downloading, setDownloading] = useState(false)
  const [sendingEmail, setSendingEmail] = useState(false)

  const id = window.location.pathname.split('/').pop()

  const fetchSession = () => {
    api.get(`/sessions/${id}`).then((res) => {
      const data = res.data.data
      setSession(data.session)
      setPhotos(data.photos || [])
      setLoading(false)
    })
  }

  const fetchPrintJobs = () => {
    api.get(`/print-jobs?session_id=${id}`).then((res) => {
      setPrintJobs(res.data.data || [])
    }).catch(() => {})
  }

  useEffect(() => {
    fetchSession()
    fetchPrintJobs()
  }, [])

  const handleGenerateLink = async () => {
    setDownloading(true)
    try {
      const res = await api.post(`/sessions/${id}/generate-link`)
      const downloadUrl = res.data.data.download_url
      alert(`Secure download link (24h):\n${downloadUrl}`)
    } catch (err) {
      alert('Failed to generate link')
    }
    setDownloading(false)
  }

  const handleSendEmail = async () => {
    setSendingEmail(true)
    try {
      await api.post(`/sessions/${id}/send-email`)
      alert('Email sent!')
    } catch (err) {
      alert('Failed to send email')
    }
    setSendingEmail(false)
  }

  const handleReprint = async () => {
    try {
      await api.post(`/sessions/${id}/reprint`)
      alert('Reprint job queued!')
      fetchPrintJobs()
    } catch (err) {
      alert('Failed to reprint')
    }
  }

  const handleAutoPrint = async () => {
    try {
      await api.post(`/sessions/${id}/auto-print`)
      alert('Auto-print started!')
      fetchPrintJobs()
    } catch (err) {
      alert('Failed to start auto-print')
    }
  }

  const handleDelete = async () => {
    if (confirm('Delete this session?')) {
      await api.delete(`/sessions/${id}`)
      window.location.href = '/sessions'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'QUEUED': return <Clock size={14} className="text-gray-500" />
      case 'PREPARING': return <Clock size={14} className="text-yellow-500" />
      case 'SENDING': return <Clock size={14} className="text-blue-500" />
      case 'PRINTING': return <Printer size={14} className="text-blue-500 animate-pulse" />
      case 'PRINT_COMPLETE': return <CheckCircle size={14} className="text-green-500" />
      case 'PRINTER_ERROR': return <AlertCircle size={14} className="text-red-500" />
      default: return <Clock size={14} className="text-gray-400" />
    }
  }

  if (loading) return <p className="text-gray-500">Loading...</p>
  if (!session) return <p className="text-gray-500">Session not found.</p>

  return (
    <div>
      <button onClick={() => window.history.back()} className="flex items-center gap-2 text-gray-600 mb-6">
        <ArrowLeft size={20} /> Back
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Image */}
        <div className="lg:col-span-2">
          <div className="card">
            <div className="aspect-[3/4] bg-gradient-to-br from-pink-100 to-purple-100 rounded-2xl flex items-center justify-center mb-4">
              {session.final_image_url ? (
                <img src={`/api/download/${session.id}`} alt="Final" className="max-w-full max-h-full object-contain rounded-2xl" />
              ) : (
                <span className="text-6xl">📸</span>
              )}
            </div>

            {/* Individual Photos */}
            {photos.length > 0 && (
              <div className="flex gap-3 overflow-x-auto pb-2">
                {photos.map((p) => (
                  <div key={p.id} className="w-24 h-24 bg-gray-100 rounded-xl flex-shrink-0 flex items-center justify-center">
                    <span className="text-2xl">🖼️</span>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* Actions */}
          <div className="card">
            <h3 className="font-semibold mb-3">Actions</h3>
            <div className="space-y-2">
              <button onClick={handleGenerateLink} disabled={downloading} className="btn-primary w-full flex items-center justify-center gap-2">
                <Download size={18} /> {downloading ? 'Generating...' : 'Get Download Link'}
              </button>
              <button onClick={handleSendEmail} disabled={sendingEmail} className="btn-secondary w-full flex items-center justify-center gap-2">
                <Mail size={18} /> {sendingEmail ? 'Sending...' : 'Send Email'}
              </button>
              <button onClick={handleAutoPrint} className="btn-secondary w-full flex items-center justify-center gap-2">
                <Printer size={18} /> Auto Print
              </button>
              <button onClick={handleReprint} className="btn-secondary w-full flex items-center justify-center gap-2">
                <RefreshCw size={18} /> Reprint
              </button>
              <button onClick={handleDelete} className="w-full py-2 text-red-500 hover:bg-red-50 rounded-xl flex items-center justify-center gap-2">
                <Trash2 size={18} /> Delete
              </button>
            </div>
          </div>

          {/* Session Info */}
          <div className="card">
            <h3 className="font-semibold mb-3">Session Info</h3>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-500">Status</span>
                <span className="font-medium">{session.status}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500">Payment</span>
                <span className={`font-medium ${session.payment_status === 'PAID' ? 'text-green-600' : 'text-yellow-600'}`}>
                  {session.payment_status}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500">Email</span>
                <span className="font-medium">{session.email || 'Not provided'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500">Created</span>
                <span className="font-medium">{new Date(session.created_at).toLocaleString()}</span>
              </div>
            </div>
          </div>

          {/* Print Jobs */}
          {printJobs.length > 0 && (
            <div className="card">
              <h3 className="font-semibold mb-3">Print Jobs</h3>
              <div className="space-y-2">
                {printJobs.map((job) => (
                  <div key={job.id} className="flex items-center justify-between p-2 bg-gray-50 rounded-lg">
                    <div className="flex items-center gap-2">
                      {getStatusIcon(job.status)}
                      <div>
                        <p className="text-sm font-medium">{job.print_type} × {job.copies}</p>
                        <p className="text-xs text-gray-500">{job.status}</p>
                      </div>
                    </div>
                    {job.printer_name && <span className="text-xs text-gray-400">{job.printer_name}</span>}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
