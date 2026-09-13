import React from 'react'
import { motion } from 'framer-motion'
import { usePhotoboothStore } from '../../../stores/photoboothStore'
import { ArrowLeft, RotateCcw, Download, Sparkles, Image } from 'lucide-react'

export default function PhotoReview() {
  const {
    setStep,
    selectedTemplate,
    capturedPhotos,
    clearCapturedPhotos,
  } = usePhotoboothStore()

  const handleRetake = () => {
    clearCapturedPhotos()
    setStep('camera-capture')
  }

  const handleContinue = () => {
    setStep('email-qr')
  }

  if (!selectedTemplate) {
    setStep('template-selection')
    return null
  }

  return (
    <motion.div
      className="relative w-full h-full overflow-hidden bg-surface to-surface"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.5 }}
    >
      {/* Animated background */}
      <div className="absolute inset-0">
        <motion.div
          className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-primary-subtle to-primary-subtle"
          animate={{ opacity: [0.5, 0.8, 0.5] }}
          transition={{ duration: 4, repeat: Infinity }}
        />
        {[...Array(20)].map((_, i) => (
          <motion.div
            key={i}
            className="absolute w-1 h-1 bg-surface/30 rounded-full"
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
            }}
            animate={{
              y: [0, -30, 0],
              opacity: [0, 1, 0],
            }}
            transition={{
              duration: 3 + Math.random() * 2,
              repeat: Infinity,
              delay: Math.random() * 2,
            }}
          />
        ))}
      </div>

      <div className="relative z-10 flex flex-col h-full p-6 md:p-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <motion.button
            onClick={handleRetake}
            className="w-14 h-14 rounded-full bg-surface-alt  flex items-center justify-center"
            whileTap={{ scale: 0.9 }}
          >
            <ArrowLeft size={24} className="text-text" />
          </motion.button>

          <h1 className="text-2xl md:text-3xl font-bold text-text">
            Your Photos
          </h1>

          <div className="w-14" />
        </div>

        {/* Photo Preview Area */}
        <motion.div
          className="flex-1 flex items-center justify-center"
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ delay: 0.2, duration: 0.5 }}
        >
          <div className="relative w-full max-w-2xl aspect-[4/3] bg-surface-alt rounded-lg overflow-hidden border-4 border-white/20 shadow-md">
            {/* Placeholder composed image */}
            {capturedPhotos.length > 0 ? (
              <img
                src={capturedPhotos[0]}
                alt="Composed photo"
                className="w-full h-full object-cover"
              />
            ) : (
              <div className="w-full h-full flex flex-col items-center justify-center bg-gradient-to-br from-surface-alt to-border">
                <Image size={80} className="text-text-muted mb-4" />
                <p className="text-text-subtle text-lg font-medium">Your composed photo</p>
                <p className="text-text-muted text-sm mt-1">
                  {selectedTemplate.name} • {selectedTemplate.photo_count} shots
                </p>
              </div>
            )}

            {/* Decorative corners */}
            <div className="absolute top-4 left-4 w-8 h-8 border-t-2 border-l-2 border-accent/60 rounded-tl-lg" />
            <div className="absolute top-4 right-4 w-8 h-8 border-t-2 border-r-2 border-accent/60 rounded-tr-lg" />
            <div className="absolute bottom-4 left-4 w-8 h-8 border-b-2 border-l-2 border-accent/60 rounded-bl-lg" />
            <div className="absolute bottom-4 right-4 w-8 h-8 border-b-2 border-r-2 border-accent/60 rounded-br-lg" />

            {/* Watermark */}
            <div className="absolute bottom-4 left-1/2 -translate-x-1/2 px-4 py-1 bg-surface-alt  rounded-full">
              <span className="text-text/60 text-xs font-medium">KLICKU</span>
            </div>
          </div>
        </motion.div>

        {/* Bottom Controls */}
        <motion.div
          className="mt-6 flex flex-col md:flex-row gap-4"
          initial={{ y: 30, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.4 }}
        >
          <motion.button
            onClick={handleRetake}
            className="flex-1 py-5 rounded-lg bg-surface-alt  text-text font-semibold text-lg flex items-center justify-center gap-3 active:scale-95 transition-transform"
            whileTap={{ scale: 0.95 }}
          >
            <RotateCcw size={22} />
            Retake Photos
          </motion.button>
          <motion.button
            onClick={handleContinue}
            className="flex-1 py-5 rounded-lg bg-primary text-text font-bold text-lg shadow-md flex items-center justify-center gap-3 active:scale-95 transition-transform"
            whileTap={{ scale: 0.95 }}
          >
            <Sparkles size={22} />
            Looks Great!
          </motion.button>
        </motion.div>

        {/* Photo strip preview */}
        <div className="mt-4 flex justify-center gap-2">
          {capturedPhotos.map((photo, i) => (
            <motion.div
              key={i}
              className="w-12 h-12 rounded-lg overflow-hidden border-2 border-white/30 shadow-md"
              initial={{ scale: 0, rotate: -10 }}
              animate={{ scale: 1, rotate: 0 }}
              transition={{ type: 'spring', delay: 0.5 + i * 0.1 }}
            >
              <img src={photo} alt={`Shot ${i + 1}`} className="w-full h-full object-cover" />
            </motion.div>
          ))}
        </div>
      </div>
    </motion.div>
  )
}