import React, { useState } from 'react'
import { motion } from 'framer-motion'
import { usePhotoboothStore } from '../../../stores/photoboothStore'
import { Mail, QrCode, SkipForward, Send, CheckCircle } from 'lucide-react'
import { QRCodeSVG } from 'qrcode.react'

export default function EmailQRScreen() {
  const { setStep, customerEmail, setCustomerEmail, getDownloadUrl, downloadUrl, sessionData, startIdleTimer } = usePhotoboothStore()
  const [isEmailSent, setIsEmailSent] = useState(false)
  const [localEmail, setLocalEmail] = useState(customerEmail)

  const handleEmailSubmit = async () => {
    if (localEmail && localEmail.includes('@')) {
      setCustomerEmail(localEmail)
      await getDownloadUrl()
      setIsEmailSent(true)
    }
  }

  const handleSkip = async () => {
    await getDownloadUrl()
    setStep('promotion')
  }

  const handleContinue = () => {
    setStep('promotion')
  }

  return (
    <motion.div
      className="relative w-full h-full overflow-hidden bg-gradient-to-br from-gray-50 via-white to-purple-50"
      initial={{ opacity: 0, x: 100 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -100 }}
      transition={{ duration: 0.5 }}
    >
      {/* Background decoration */}
      <div className="absolute top-0 right-0 w-72 h-72 bg-gradient-to-bl from-primary/10 to-transparent rounded-full -translate-y-1/2 translate-x-1/2" />
      <div className="absolute bottom-0 left-0 w-56 h-56 bg-gradient-to-tr from-secondary/10 to-transparent rounded-full translate-y-1/2 -translate-x-1/2" />

      <div className="relative z-10 flex flex-col h-full p-6 md:p-10">
        {/* Header */}
        <div className="text-center mb-8">
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: 'spring', stiffness: 300, delay: 0.2 }}
            className="w-20 h-20 mx-auto mb-4 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center shadow-lg"
          >
            <CheckCircle size={40} className="text-white" />
          </motion.div>
          <h1 className="text-3xl md:text-4xl font-bold text-gray-800 mb-2">
            Get Your Photos!
          </h1>
          <p className="text-gray-500 text-lg">
            Enter your email or scan the QR code
          </p>
        </div>

        {/* Main content */}
        <div className="flex-1 flex flex-col md:flex-row items-center justify-center gap-8">
          {/* Email input section */}
          <motion.div
            className="w-full max-w-md"
            initial={{ x: -30, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.3 }}
          >
            <div className="bg-white rounded-3xl shadow-xl p-6 border border-gray-100">
              <div className="flex items-center gap-3 mb-4">
                <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
                  <Mail size={20} className="text-primary" />
                </div>
                <h2 className="text-xl font-bold text-gray-800">Email Delivery</h2>
              </div>

              {!isEmailSent ? (
                <>
                  <div className="relative mb-4">
                    <input
                      type="email"
                      value={localEmail}
                      onChange={(e) => setLocalEmail(e.target.value)}
                      placeholder="your@email.com"
                      className="w-full px-5 py-4 rounded-2xl border-2 border-gray-200 focus:border-primary focus:outline-none text-lg transition-colors"
                    />
                    <Mail size={20} className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400" />
                  </div>
                  <motion.button
                    onClick={handleEmailSubmit}
                    disabled={!localEmail || !localEmail.includes('@')}
                    className="w-full py-4 rounded-2xl bg-gradient-to-r from-primary to-secondary text-white font-bold text-lg shadow-lg disabled:opacity-50 disabled:cursor-not-allowed active:scale-95 transition-transform"
                    whileTap={{ scale: 0.98 }}
                  >
                    <span className="flex items-center justify-center gap-2">
                      <Send size={20} />
                      Send to Email
                    </span>
                  </motion.button>
                </>
              ) : (
                <motion.div
                  initial={{ scale: 0.9, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  className="text-center py-4"
                >
                  <CheckCircle size={48} className="mx-auto text-green-500 mb-3" />
                  <p className="text-gray-800 font-semibold">Email sent!</p>
                  <p className="text-gray-500 text-sm mt-1">Check your inbox shortly</p>
                </motion.div>
              )}
            </div>
          </motion.div>

          {/* Divider */}
          <div className="hidden md:flex flex-col items-center gap-2">
            <div className="w-px h-16 bg-gray-200" />
            <span className="text-gray-400 font-medium">OR</span>
            <div className="w-px h-16 bg-gray-200" />
          </div>
          <div className="md:hidden flex items-center gap-2 w-full">
            <div className="flex-1 h-px bg-gray-200" />
            <span className="text-gray-400 font-medium">OR</span>
            <div className="flex-1 h-px bg-gray-200" />
          </div>

          {/* QR Code section */}
          <motion.div
            className="w-full max-w-md"
            initial={{ x: 30, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.4 }}
          >
            <div className="bg-white rounded-3xl shadow-xl p-6 border border-gray-100 text-center">
              <div className="flex items-center justify-center gap-3 mb-4">
                <div className="w-10 h-10 rounded-full bg-secondary/10 flex items-center justify-center">
                  <QrCode size={20} className="text-secondary" />
                </div>
                <h2 className="text-xl font-bold text-gray-800">Scan to Download</h2>
              </div>

              <div className="w-48 h-48 mx-auto mb-4 p-4 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200 flex items-center justify-center">
                <QRCodeSVG
                  value={downloadUrl || `https://klikku.app/download/${sessionData?.session_id || 'demo'}`}
                  size={160}
                  level="H"
                  includeMargin={false}
                  fgColor="#ff6b9d"
                />
              </div>

              <p className="text-gray-500 text-sm">
                Scan with your phone camera to download instantly
              </p>
            </div>
          </motion.div>
        </div>

        {/* Bottom buttons */}
        <motion.div
          className="mt-6 flex gap-4"
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.5 }}
        >
          <motion.button
            onClick={handleSkip}
            className="flex-1 py-4 rounded-2xl bg-gray-100 text-gray-700 font-semibold text-lg flex items-center justify-center gap-2 active:scale-95 transition-transform"
            whileTap={{ scale: 0.95 }}
          >
            <SkipForward size={20} />
            Skip for Now
          </motion.button>
          <motion.button
            onClick={handleContinue}
            className="flex-1 py-4 rounded-2xl bg-gradient-to-r from-primary to-secondary text-white font-bold text-lg shadow-lg flex items-center justify-center gap-2 active:scale-95 transition-transform"
            whileTap={{ scale: 0.95 }}
          >
            Continue ✨
          </motion.button>
        </motion.div>
      </div>
    </motion.div>
  )
}