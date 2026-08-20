import React, { useRef, useEffect, useState, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { usePhotoboothStore } from '../../../stores/photoboothStore'
import { Camera, RefreshCw, X } from 'lucide-react'

type CaptureState = 'preview' | 'countdown' | 'review'

export default function CameraCapture() {
  const videoRef = useRef<HTMLVideoElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const streamRef = useRef<MediaStream | null>(null)
  const {
    setStep,
    selectedTemplate,
    capturedPhotos,
    addCapturedPhoto,
    clearCapturedPhotos,
    countdown,
    setCountdown,
    triggerFlash,
    isFlashing,
    capturePhotos,
    startIdleTimer,
  } = usePhotoboothStore()

  const [captureState, setCaptureState] = useState<CaptureState>('preview')
  const [currentShot, setCurrentShot] = useState(0)
  const [cameraError, setCameraError] = useState<string | null>(null)

  const totalPhotos = selectedTemplate?.photo_count || 4

  // Start camera
  useEffect(() => {
    async function startCamera() {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({
          video: {
            facingMode: 'user',
            width: { ideal: 1280 },
            height: { ideal: 720 },
          },
          audio: false,
        })
        streamRef.current = stream
        if (videoRef.current) {
          videoRef.current.srcObject = stream
        }
      } catch {
        setCameraError('Camera not available. Using demo mode.')
      }
    }
    startCamera()

    return () => {
      if (streamRef.current) {
        streamRef.current.getTracks().forEach((t) => t.stop())
      }
    }
  }, [])

  const captureSinglePhoto = useCallback(() => {
    if (!videoRef.current || !canvasRef.current) return null

    const video = videoRef.current
    const canvas = canvasRef.current
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight

    const ctx = canvas.getContext('2d')
    if (!ctx) return null

    // Mirror the image (selfie mode)
    ctx.translate(canvas.width, 0)
    ctx.scale(-1, 1)
    ctx.drawImage(video, 0, 0, canvas.width, canvas.height)

    return canvas.toDataURL('image/jpeg', 0.9)
  }, [])

  const runCountdown = useCallback(() => {
    setCaptureState('countdown')
    setCountdown(3)

    let count = 3
    const interval = setInterval(() => {
      count--
      setCountdown(count)

      if (count === 0) {
        clearInterval(interval)
        // Capture
        triggerFlash()
        const photoData = captureSinglePhoto()
        if (photoData) {
          addCapturedPhoto(photoData)
        }

        setTimeout(() => {
          if (currentShot + 1 >= totalPhotos) {
            setCaptureState('review')
          } else {
            setCurrentShot((prev) => prev + 1)
            setCaptureState('preview')
          }
        }, 500)
      }
    }, 1000)
  }, [captureSinglePhoto, triggerFlash, addCapturedPhoto, currentShot, totalPhotos])

  const handleRetake = () => {
    clearCapturedPhotos()
    setCurrentShot(0)
    setCaptureState('preview')
  }

  const handleContinue = async () => {
    await capturePhotos()
    setStep('photo-review')
  }

  const handleBack = () => {
    setStep('payment')
  }

  if (!selectedTemplate) {
    setStep('template-selection')
    return null
  }

  return (
    <motion.div
      className="relative w-full h-full overflow-hidden bg-black"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.4 }}
    >
      {/* Camera preview */}
      <video
        ref={videoRef}
        autoPlay
        playsInline
        muted
        className="absolute inset-0 w-full h-full object-cover -scale-x-100"
      />

      {/* Hidden canvas for capture */}
      <canvas ref={canvasRef} className="hidden" />

      {/* Flash effect */}
      <AnimatePresence>
        {isFlashing && (
          <motion.div
            className="absolute inset-0 bg-white z-50"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.15 }}
          />
        )}
      </AnimatePresence>

      {/* Camera error overlay */}
      {cameraError && (
        <div className="absolute inset-0 bg-gradient-to-br from-gray-900 to-gray-800 flex items-center justify-center">
          <div className="text-center">
            <Camera size={64} className="mx-auto text-gray-500 mb-4" />
            <p className="text-gray-400 text-lg">{cameraError}</p>
          </div>
        </div>
      )}

      {/* Countdown overlay */}
      <AnimatePresence>
        {captureState === 'countdown' && (
          <motion.div
            className="absolute inset-0 flex items-center justify-center z-40"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            <motion.div
              key={countdown}
              className="text-9xl font-bold text-white drop-shadow-2xl"
              initial={{ scale: 2, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.5, opacity: 0 }}
              transition={{ duration: 0.3 }}
            >
              {countdown || '📸'}
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Top bar - progress */}
      <div className="absolute top-0 left-0 right-0 p-4 z-30">
        <div className="flex items-center justify-between">
          <motion.button
            onClick={handleBack}
            className="w-12 h-12 rounded-full bg-black/40 backdrop-blur-sm flex items-center justify-center"
            whileTap={{ scale: 0.9 }}
          >
            <X size={24} className="text-white" />
          </motion.button>

          <div className="flex items-center gap-2">
            {Array.from({ length: totalPhotos }).map((_, i) => (
              <motion.div
                key={i}
                className={`w-3 h-3 rounded-full ${
                  i < capturedPhotos.length
                    ? 'bg-primary'
                    : i === currentShot
                    ? 'bg-white'
                    : 'bg-white/30'
                }`}
                animate={i === currentShot ? { scale: [1, 1.3, 1] } : {}}
                transition={{ duration: 1, repeat: Infinity }}
              />
            ))}
          </div>

          <div className="w-12" />
        </div>
      </div>

      {/* Preview of captured photos */}
      {capturedPhotos.length > 0 && (
        <div className="absolute bottom-32 left-4 right-4 z-30">
          <div className="flex gap-2 overflow-x-auto pb-2">
            {capturedPhotos.map((photo, i) => (
              <motion.div
                key={i}
                className="flex-shrink-0 w-16 h-16 rounded-xl overflow-hidden border-2 border-white/50 shadow-lg"
                initial={{ scale: 0, y: 20 }}
                animate={{ scale: 1, y: 0 }}
                transition={{ type: 'spring', stiffness: 300, delay: i * 0.05 }}
              >
                <img src={photo} alt={`Shot ${i + 1}`} className="w-full h-full object-cover" />
              </motion.div>
            ))}
          </div>
        </div>
      )}

      {/* Bottom controls */}
      <div className="absolute bottom-0 left-0 right-0 p-6 z-30">
        {captureState === 'preview' && (
          <motion.div
            className="flex flex-col items-center gap-4"
            initial={{ y: 30, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
          >
            <p className="text-white text-lg font-medium drop-shadow-lg">
              Photo {currentShot + 1} of {totalPhotos} — Get ready!
            </p>
            <motion.button
              onClick={runCountdown}
              className="w-20 h-20 rounded-full bg-white shadow-2xl flex items-center justify-center"
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.9 }}
            >
              <div className="w-16 h-16 rounded-full border-4 border-primary flex items-center justify-center">
                <Camera size={28} className="text-primary" />
              </div>
            </motion.button>
          </motion.div>
        )}

        {captureState === 'review' && (
          <motion.div
            className="flex flex-col items-center gap-4"
            initial={{ y: 30, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.3 }}
          >
            <p className="text-white text-lg font-medium drop-shadow-lg">
              All done! {totalPhotos} photos captured 🎉
            </p>
            <div className="flex gap-4">
              <motion.button
                onClick={handleRetake}
                className="px-6 py-4 rounded-2xl bg-white/20 backdrop-blur-sm text-white font-semibold flex items-center gap-2"
                whileTap={{ scale: 0.95 }}
              >
                <RefreshCw size={20} />
                Retake All
              </motion.button>
              <motion.button
                onClick={handleContinue}
                className="px-8 py-4 rounded-2xl bg-gradient-to-r from-primary to-secondary text-white font-bold shadow-lg flex items-center gap-2"
                whileTap={{ scale: 0.95 }}
              >
                Continue ✨
              </motion.button>
            </div>
          </motion.div>
        )}
      </div>
    </motion.div>
  )
}