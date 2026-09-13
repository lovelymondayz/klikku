import React from 'react'
import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../../stores/authStore'
import { Users, Camera, BarChart3, LogOut, ArrowLeft } from 'lucide-react'

const navItems = [
  { path: '/admin', label: 'Merchants', icon: Users },
  { path: '/admin/sessions', label: 'All Sessions', icon: Camera },
  { path: '/admin/analytics', label: 'Analytics', icon: BarChart3 },
]

export default function AdminLayout() {
  const location = useLocation()
  const navigate = useNavigate()
  const { name, email, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="min-h-screen bg-surface text-text flex">
      <aside className="w-64 bg-surface flex flex-col">
        <div className="p-6 border-b border-border-strong">
          <h1 className="text-2xl font-bold">Klikku Admin</h1>
          <p className="text-sm text-text-subtle">Super Admin</p>
        </div>

        <nav className="flex-1 p-4 space-y-1">
          {navItems.map((item) => {
            const Icon = item.icon
            const active = location.pathname === item.path
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`flex items-center gap-3 px-4 py-3 rounded-md transition-colors ${
                  active ? 'bg-primary/20 text-primary' : 'text-text-subtle hover:bg-surface hover:text-text'
                }`}
              >
                <Icon size={20} />
                {item.label}
              </Link>
            )
          })}
        </nav>

        <div className="p-4 border-t border-border-strong">
          <Link to="/dashboard" className="flex items-center gap-2 text-text-subtle hover:text-text mb-3">
            <ArrowLeft size={16} /> Back to Dashboard
          </Link>
          <button onClick={handleLogout} className="flex items-center gap-2 text-text-subtle hover:text-text">
            <LogOut size={16} /> Sign Out
          </button>
        </div>
      </aside>

      <main className="flex-1 p-8 overflow-auto">
        <Outlet />
      </main>
    </div>
  )
}
