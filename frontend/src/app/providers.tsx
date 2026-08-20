import React, { createContext, useContext, useEffect, useState } from 'react'
import { useAuthStore } from '../stores/authStore'
import { api } from '../lib/api'

interface AuthContextType {
  loading: boolean
}

const AuthContext = createContext<AuthContextType>({ loading: true })

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { token, setAuth, clearAuth } = useAuthStore()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Try to restore session
    const storedToken = localStorage.getItem('klikku_token')
    const storedRefresh = localStorage.getItem('klikku_refresh')
    const storedRole = localStorage.getItem('klikku_role')
    const storedMerchant = localStorage.getItem('klikku_merchant')

    if (storedToken && storedRole) {
      api.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`
      setAuth({
        token: storedToken,
        refreshToken: storedRefresh || '',
        role: storedRole,
        merchantId: storedMerchant || '',
        email: '',
        name: '',
      })
    }
    setLoading(false)
  }, [])

  return (
    <AuthContext.Provider value={{ loading }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => useContext(AuthContext)
