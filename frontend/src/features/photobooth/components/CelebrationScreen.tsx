import React, { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { usePhotoboothStore, STEPS } from '../../../stores/photoboothStore'
import { Star, Heart, Sparkles, PartyPopper } from 'lucide-react'

const PARTICLE_COUNT = 30

export default function CelebrationScreen() {
  const { setStep } = usePhotoboothStore()
  const [particles, setParticles] = useState<Array<{ id: number; x: number; y: number; size: number }>>([])

  useEffect(() => {
    const newParticles = Array.from({ length: PARTICLE_COUNT }, (_, i) => ({
      id: i,
      x: Math.random() * 100,
      y: Math.random() * 100,
      size: Math.random() * 20 + 10,
    }))
    setParticles(newParticles)

    const timer = setTimeout(() => {
      setStep(STEPS.IDLE)
    }, 10000)
    return () => clearTimeout(timer)
  }, [setStep])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-subtle via-surface-alt to-info-subtle overflow-hidden relative">
      {/* Celebration particles */}
      {particles.map((p) => (
        <motion.div
          key={p.id}
          className="absolute rounded-full bg-primary"
          style={{
            left: `${p.x}%`,
            top: `${p.y}%`,
            width: p.size,
            height: p.size,
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
          <PartyPopper size={80} className="mx-auto text-primary mb-6" />
        </motion.div>
        <h1 className="text-5xl font-bold mb-4">Congratulations! 🎉</h1>
        <p className="text-xl text-text-muted mb-8">Your photos are being processed</p>
        <div className="flex justify-center gap-4">
          {[...Array(5)].map((_, i) => (
            <motion.div
              key={i}
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.8 + i * 0.1 }}
            >
              <Star size={30} className="text-warning fill-warning" />
            </motion.div>
          ))}
        </div>
      </motion.div>
    </div>
  )
}
