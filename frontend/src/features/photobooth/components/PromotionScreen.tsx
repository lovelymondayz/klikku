import React, { useEffect } from 'react'
import { motion } from 'framer-motion'
import { usePhotoboothStore, STEPS } from '../../../stores/photoboothStore'
import { Gift, Heart, Star, Instagram, Sparkles } from 'lucide-react'

export default function PromotionScreen() {
  const { setStep } = usePhotoboothStore()

  useEffect(() => {
    const timer = setTimeout(() => {
      setStep(STEPS.IDLE)
    }, 30000)
    return () => clearTimeout(timer)
  }, [])

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-pink-100 via-purple-50 to-blue-100 p-8">
      <motion.div
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        className="text-center max-w-lg"
      >
        <Sparkles className="w-16 h-16 text-yellow-400 mx-auto mb-4" />
        <h2 className="text-4xl font-bold mb-4">Thanks for visiting!</h2>
        <p className="text-gray-600 mb-6">Follow us for more fun moments</p>
        <div className="flex gap-4 justify-center">
          <Instagram className="w-10 h-10 text-pink-500" />
          <Heart className="w-10 h-10 text-red-500" />
          <Star className="w-10 h-10 text-yellow-500" />
        </div>
        <button
          onClick={() => setStep(STEPS.IDLE)}
          className="mt-8 px-6 py-3 bg-white rounded-2xl shadow-lg font-semibold"
        >
          Start Over
        </button>
      </motion.div>
    </div>
  )
}
