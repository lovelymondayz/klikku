import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'

// Auth
import LoginPage from '../features/auth/LoginPage'

// Dashboard
import DashboardLayout from '../features/dashboard/DashboardLayout'
import OverviewPage from '../features/dashboard/OverviewPage'
import CampaignsPage from '../features/campaigns/CampaignsPage'
import TemplatesPage from '../features/templates/TemplatesPage'
import SessionsPage from '../features/sessions/SessionsPage'
import SessionDetailPage from '../features/sessions/SessionDetailPage'
import BrandingPage from '../features/branding/BrandingPage'
import DevicesPage from '../features/devices/DevicesPage'

// Photobooth
import PhotoboothPage from '../features/photobooth/PhotoboothPage'

// Admin
import AdminLayout from '../features/admin/AdminLayout'
import AdminMerchantsPage from '../features/admin/AdminMerchantsPage'
import AdminSessionsPage from '../features/admin/AdminSessionsPage'
import AdminAnalyticsPage from '../features/admin/AdminAnalyticsPage'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const token = useAuthStore((s: any) => s.token)
  if (!token) return <Navigate to="/login" replace />
  return <>{children}</>
}

function AdminRoute({ children }: { children: React.ReactNode }) {
  const role = useAuthStore((s: any) => s.role)
  if (role !== 'SUPER_ADMIN') return <Navigate to="/dashboard" replace />
  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      {/* Public */}
      <Route path="/login" element={<LoginPage />} />

      {/* Photobooth (device-facing) */}
      <Route path="/photobooth/:deviceToken" element={<PhotoboothPage />} />

      {/* Merchant Dashboard */}
      <Route path="/dashboard" element={<ProtectedRoute><DashboardLayout /></ProtectedRoute>}>
        <Route index element={<OverviewPage />} />
        <Route path="campaigns" element={<CampaignsPage />} />
        <Route path="templates" element={<TemplatesPage />} />
        <Route path="sessions" element={<SessionsPage />} />
        <Route path="sessions/:id" element={<SessionDetailPage />} />
        <Route path="branding" element={<BrandingPage />} />
        <Route path="devices" element={<DevicesPage />} />
      </Route>

      {/* Super Admin */}
      <Route path="/admin" element={<ProtectedRoute><AdminRoute><AdminLayout /></AdminRoute></ProtectedRoute>}>
        <Route index element={<AdminMerchantsPage />} />
        <Route path="sessions" element={<AdminSessionsPage />} />
        <Route path="analytics" element={<AdminAnalyticsPage />} />
      </Route>

      {/* Default */}
      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  )
}
