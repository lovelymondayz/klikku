import { create } from 'zustand'
import { api } from '../lib/api'

export const STEPS = {
  IDLE: 'IDLE',
  TEMPLATE_SELECT: 'TEMPLATE_SELECT',
  PAYMENT: 'PAYMENT',
  CAPTURE: 'CAPTURE',
  REVIEW: 'REVIEW',
  PROCESSING: 'PROCESSING',
  FINAL: 'FINAL',
  EMAIL: 'EMAIL',
  PROMOTION: 'PROMOTION',
  DONE: 'DONE',
}

export const MOCK_TEMPLATES = [
  { id: '1', name: 'Classic Strip', photo_count: 4, price: 25000, preview_url: null, layout_config: { output_width: 1200, output_height: 1800 } },
  { id: '2', name: 'Polaroid', photo_count: 1, price: 15000, preview_url: null, layout_config: { output_width: 1200, output_height: 1800 } },
  { id: '3', name: 'Best Friends', photo_count: 3, price: 20000, preview_url: null, layout_config: { output_width: 1200, output_height: 1800 } },
  { id: '4', name: 'Summer Vibes', photo_count: 4, price: 30000, preview_url: null, layout_config: { output_width: 1200, output_height: 1800 } },
]

interface PhotoboothStore {
  currentStep: string
  sessionID: string | null
  merchantData: any
  selectedTemplate: any
  capturedPhotos: string[]
  finalImageURL: string | null
  email: string
  loading: boolean
  error: string | null
  countdown: number
  isFlashing: boolean
  customerEmail: string
  downloadUrl: string
  sessionData: any
  isLoading: boolean
  
  // Actions
  setStep: (step: string) => void
  setMerchantData: (data: any) => void
  selectTemplate: (template: any) => void
  addPhoto: (photo: string) => void
  setEmail: (email: string) => void
  reset: () => void
  setCountdown: (n: number) => void
  triggerFlash: () => void
  clearCapturedPhotos: () => void
  addCapturedPhoto: (photo: string) => void
  setCustomerEmail: (email: string) => void
  setLoading: (loading: boolean) => void
  
  // API calls
  fetchAttractData: (deviceToken: string) => Promise<void>
  createSession: (deviceToken?: string) => Promise<void>
  capturePhotos: () => Promise<void>
  finalizeSession: () => Promise<void>
  startIdleTimer: () => void
  getDownloadUrl: () => string
}

export const usePhotoboothStore = create<PhotoboothStore>((set, get) => ({
  currentStep: STEPS.IDLE,
  sessionID: null,
  merchantData: null,
  selectedTemplate: null,
  capturedPhotos: [],
  finalImageURL: null,
  email: '',
  loading: false,
  error: null,
  countdown: 0,
  isFlashing: false,
  customerEmail: '',
  downloadUrl: '',
  sessionData: null,
  isLoading: false,

  setStep: (step) => set({ currentStep: step }),
  setMerchantData: (data) => set({ merchantData: data }),
  selectTemplate: (template) => set({ selectedTemplate: template }),
  addPhoto: (photo: string) => set((s) => ({ capturedPhotos: [...s.capturedPhotos, photo] })),
  setEmail: (email) => set({ email }),
  setCountdown: (n) => set({ countdown: n }),
  triggerFlash: () => set({ isFlashing: true }),
  clearCapturedPhotos: () => set({ capturedPhotos: [] }),
  addCapturedPhoto: (photo: string) => set((s) => ({ capturedPhotos: [...s.capturedPhotos, photo] })),
  setCustomerEmail: (email) => set({ customerEmail: email }),
  setLoading: (loading) => set({ loading, isLoading: loading }),

  reset: () => set({
    currentStep: STEPS.IDLE,
    sessionID: null,
    selectedTemplate: null,
    capturedPhotos: [],
    finalImageURL: null,
    email: '',
    error: null,
    customerEmail: '',
    downloadUrl: '',
    sessionData: null,
    countdown: 0,
    isFlashing: false,
  }),

  fetchAttractData: async (deviceToken) => {
    set({ loading: true, isLoading: true, error: null })
    try {
      const { data } = await api.get(`/devices/${deviceToken}/attract`)
      set({ merchantData: data, loading: false, isLoading: false })
    } catch (err: any) {
      set({
        merchantData: {
          business_name: 'Klikku Photobooth',
          logo_url: null,
          primary_color: '#ff6b9d',
          welcome_message: 'Take fun photos and bring home your memories!',
        },
        loading: false,
        isLoading: false,
      })
    }
  },

  createSession: async (deviceToken) => {
    set({ loading: true, isLoading: true, error: null })
    try {
      const { data } = await api.post(`/devices/${deviceToken}/session`)
      set({ sessionID: data.session_id, loading: false, isLoading: false })
    } catch (err: any) {
      set({ error: err.response?.data?.error || 'Failed to create session', loading: false, isLoading: false })
    }
  },

  capturePhotos: async () => {
    const { sessionID, capturedPhotos } = get()
    if (!sessionID) return

    set({ loading: true, isLoading: true, error: null })
    try {
      const formData = new FormData()
      capturedPhotos.forEach((photo, i) => {
        // Convert data URL to Blob if needed
        if (photo.startsWith('data:')) {
          const arr = photo.split(',')
          const mimeMatch = arr[0].match(/:(.*?);/)
          const mime = mimeMatch ? mimeMatch[1] : 'image/jpeg'
          const bstr = atob(arr[1])
          const u8arr = new Uint8Array(bstr.length)
          for (let j = 0; j < bstr.length; j++) {
            u8arr[j] = bstr.charCodeAt(j)
          }
          const blob = new Blob([u8arr], { type: mime })
          formData.append('photos', blob, `photo_${i}.jpg`)
        } else {
          formData.append('photos', photo as any, `photo_${i}.jpg`)
        }
      })

      await api.post(`/sessions/${sessionID}/capture`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      set({ loading: false, isLoading: false })
    } catch (err: any) {
      set({ error: err.response?.data?.error || 'Failed to capture photos', loading: false, isLoading: false })
    }
  },

  finalizeSession: async () => {
    const { sessionID } = get()
    if (!sessionID) return

    set({ loading: true, isLoading: true, error: null, currentStep: STEPS.PROCESSING })
    try {
      const { data } = await api.post(`/sessions/${sessionID}/finalize`)
      set({
        finalImageURL: data.final_image_url,
        downloadUrl: data.download_url,
        loading: false,
        isLoading: false,
        currentStep: STEPS.FINAL,
      })
    } catch (err: any) {
      set({ error: err.response?.data?.error || 'Failed to finalize', loading: false, isLoading: false })
    }
  },

  startIdleTimer: () => {
    const timer = setTimeout(() => {
      get().reset()
    }, 30000)
    return () => clearTimeout(timer)
  },

  getDownloadUrl: () => {
    const { sessionID } = get()
    return sessionID ? `/api/download/${sessionID}` : ''
  },
}))
