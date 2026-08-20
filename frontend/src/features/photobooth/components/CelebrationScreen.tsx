import React, { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { usePhotoboothStore } from '../../../stores/photoboothStore'
import { Star, Heart, Sparkles, PartyPopper } from 'lucide-react'

export default function CelebrationScreen() {
  const { setStep } = usePhotoboothStore()
  const [particles, setParticles] = useState<Array<{ id: number; x: number; y: number; color: string; size: number }>>([])

  useEffect(() => {
    // Generate random celebration particles
    const newParticles = Array.from({ length: 30 }, (_, i) => ({
      id: i,
      x: Math.random() * 100,
      y: Math.random() * 100,
      color: ['#ff6b9d', '#c44dff', '#ffd93d', '#ff6b6b', '#4ecdc4'][Math.floor(Math.random() * 5)],
      size: Math.random() * 20 + 10,
    }))
    setParticles(newParticles)

    // Auto-return to idle after 10 seconds
    const timer = setTimeout(() => {
      setStep('IDLE')
    }, 10000)
    return () => clearTimeout(timer)
  }, [])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-pink-200 via-purple-100 to-blue-200 overflow-hidden relative">
      {/* Celebration particles */}
      {particles.map((p) => (
        <motion.div
          key={p.id}
          className="absolute rounded-full"
          style={{
            left: `${p.x}%`,
            top: `${p.y}%`,
            width: p.size,
            height: p.size,
            backgroundColor: p.color,
          }}
          initial={{ scale: 0, opacity: 0 }}
          animate={{
            scale: [0, 1.5, 1],
            opacity: [0, 1, 0.8],
            y: [0, -50, 0],
          }}
          transition={{
            duration: 2,
            repeat: Infinity,
            delay: Math.random() * 2,
          }}
        />
      ))}

      {/* Main content */}
      <motion.div
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ type: 'spring', stiffness: 200, delay: 0.3 }}
        className="text-center z-10"
      >
        <motion.div
          animate={{ rotate: [0, 10, -10, 0] }}
          transition={{ duration: 0.5, delay: 0.5 }}
        >
          <PartyPopper size={80} className="mx-auto text-pink-500 mb-6" />
        </motion.div>
        <h1 className="text-5xl font-bold mb-4">Congratulations! 🎉</h1>
        <p className="text-xl text-gray-600 mb-8">Your photos are being processed</p>
        <div className="flex justify-center gap-4">
          {[...Array(5)].map((_, i) => (
            <motion.div
              key={i}
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.8 + i * 0.1 }}
            >
              <Star size={30} className="text-yellow-400 fill-yellow-400" />
            </motion.div>
          ))}
        </div>
      </motion.div>
    </div>
  )
}
