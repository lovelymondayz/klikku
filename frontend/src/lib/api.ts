import axios from 'axios'

export const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Clear auth and redirect to login
      localStorage.removeItem('klikku_token')
      localStorage.removeItem('klikku_refresh')
      localStorage.removeItem('klikku_role')
      localStorage.removeItem('klikku_merchant')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)
