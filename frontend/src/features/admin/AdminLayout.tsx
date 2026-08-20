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
    <div className="min-h-screen bg-gray-900 text-white flex">
      <aside className="w-64 bg-gray-800 flex flex-col">
        <div className="p-6 border-b border-gray-700">
          <h1 className="text-2xl font-bold">Klikku Admin</h1>
          <p className="text-sm text-gray-400">Super Admin</p>
        </div>

        <nav className="flex-1 p-4 space-y-1">
          {navItems.map((item) => {
            const Icon = item.icon
            const active = location.pathname === item.path
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${
                  active ? 'bg-primary/20 text-primary' : 'text-gray-400 hover:bg-gray-700 hover:text-white'
                }`}
              >
                <Icon size={20} />
                {item.label}
              </Link>
            )
          })}
        </nav>

        <div className="p-4 border-t border-gray-700">
          <Link to="/dashboard" className="flex items-center gap-2 text-gray-400 hover:text-white mb-3">
            <ArrowLeft size={16} /> Back to Dashboard
          </Link>
          <button onClick={handleLogout} className="flex items-center gap-2 text-gray-400 hover:text-white">
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
