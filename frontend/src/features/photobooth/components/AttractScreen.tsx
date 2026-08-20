import React, { useEffect, useRef } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { usePhotoboothStore } from '../../../stores/photoboothStore'
import { Camera, Sparkles, QrCode } from 'lucide-react'

export default function AttractScreen() {
  const { setStep, merchantData, fetchAttractData } = usePhotoboothStore()
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const token = window.location.pathname.split('/').pop() || 'demo-token'
    fetchAttractData(token)
  }, [])

  const handleStart = () => {
    setStep('template-selection')
  }

  return (
    <motion.div
      ref={containerRef}
      className="relative w-full h-full overflow-hidden cursor-pointer select-none"
      onClick={handleStart}
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.5 }}
    >
      <div className="absolute inset-0 bg-gradient-to-br from-primary via-secondary to-purple-600" />

      <motion.div
        className="absolute -top-20 -left-20 w-80 h-80 rounded-full bg-accent/20 blur-3xl"
        animate={{ x: [0, 60, 0], y: [0, 40, 0], scale: [1, 1.2, 1] }}
        transition={{ duration: 8, repeat: Infinity, ease: 'easeInOut' }}
      />
      <motion.div
        className="absolute -bottom-20 -right-20 w-96 h-96 rounded-full bg-pink-400/20 blur-3xl"
        animate={{ x: [0, -50, 0], y: [0, -30, 0], scale: [1, 1.3, 1] }}
        transition={{ duration: 10, repeat: Infinity, ease: 'easeInOut' }}
      />

      <div className="relative z-10 flex flex-col items-center justify-center h-full px-8 text-white">
        <AnimatePresence>
          <motion.div
            initial={{ scale: 0, rotate: -180 }}
            animate={{ scale: 1, rotate: 0 }}
            transition={{ type: 'spring', stiffness: 200, damping: 15, delay: 0.2 }}
            className="mb-6"
          >
            <div className="w-32 h-32 rounded-3xl bg-white/20 backdrop-blur-sm flex items-center justify-center shadow-2xl border border-white/30">
              <Camera size={56} className="text-white" />
            </div>
          </motion.div>
        </AnimatePresence>

        <motion.h1
          className="text-5xl md:text-7xl font-bold mb-4 text-center drop-shadow-lg"
          initial={{ y: 30, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.4, duration: 0.6 }}
        >
          {merchantData?.business_name || 'Klikku'}
        </motion.h1>

        <motion.p
          className="text-xl md:text-2xl mb-10 text-white/90 text-center font-medium"
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.6, duration: 0.6 }}
        >
          {merchantData?.welcome_message || 'Strike a pose, snap a memory!'}
        </motion.p>

        <motion.button
          className="relative px-16 py-6 bg-white text-primary font-bold text-2xl md:text-3xl rounded-full shadow-2xl active:scale-95 transition-transform"
          initial={{ y: 40, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.8, duration: 0.5 }}
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
        >
          <span className="relative flex items-center gap-3">
            <Sparkles size={28} />
            Touch to Start
            <Sparkles size={28} />
          </span>
        </motion.button>

        <motion.div
          className="absolute bottom-8 right-8 flex flex-col items-center gap-2"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1.2 }}
        >
          <div className="w-16 h-16 bg-white rounded-2xl p-2 shadow-lg">
            <QrCode size={40} className="text-gray-800 w-full h-full" />
          </div>
          <span className="text-xs text-white/70 font-medium">Follow Us</span>
        </motion.div>
      </div>
    </motion.div>
  )
}
