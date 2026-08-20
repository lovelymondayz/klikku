import { create } from 'zustand'
import { api } from '../lib/api'

interface User {
  token: string
  refreshToken: string
  role: string
  merchantId: string
  email: string
  name: string
}

interface AuthState {
  token: string | null
  refreshToken: string | null
  role: string | null
  merchantId: string | null
  email: string | null
  name: string | null
  setAuth: (user: User) => void
  clearAuth: () => void
  login: (email: string, password: string) => Promise<void>
  register: (data: { name: string; email: string; password: string; business_name: string }) => Promise<void>
  logout: () => void
}

export const useAuthStore = create<AuthState>((set, get) => ({
  token: null,
  refreshToken: null,
  role: null,
  merchantId: null,
  email: null,
  name: null,

  setAuth: (user) => {
    localStorage.setItem('klikku_token', user.token)
    localStorage.setItem('klikku_refresh', user.refreshToken)
    localStorage.setItem('klikku_role', user.role)
    localStorage.setItem('klikku_merchant', user.merchantId)
    api.defaults.headers.common['Authorization'] = `Bearer ${user.token}`
    set({
      token: user.token,
      refreshToken: user.refreshToken,
      role: user.role,
      merchantId: user.merchantId,
      email: user.email,
      name: user.name,
    })
  },

  clearAuth: () => {
    localStorage.removeItem('klikku_token')
    localStorage.removeItem('klikku_refresh')
    localStorage.removeItem('klikku_role')
    localStorage.removeItem('klikku_merchant')
    delete api.defaults.headers.common['Authorization']
    set({
      token: null,
      refreshToken: null,
      role: null,
      merchantId: null,
      email: null,
      name: null,
    })
  },

  login: async (email, password) => {
    const res = await api.post('/auth/login', { email, password })
    // Backend wraps: { success: true, data: { access_token, ... } }
    const data = res.data.data || res.data
    get().setAuth({
      token: data.access_token,
      refreshToken: data.refresh_token,
      role: data.role,
      merchantId: data.merchant_id,
      email: data.email,
      name: data.name,
    })
  },

  register: async (data) => {
    const res = await api.post('/auth/register', data)
    const result = res.data.data || res.data
    get().setAuth({
      token: result.access_token,
      refreshToken: result.refresh_token,
      role: result.role,
      merchantId: result.merchant_id,
      email: result.email,
      name: data.name,
    })
  },

  logout: () => {
    get().clearAuth()
  },
}))
