import React, { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { usePhotoboothStore } from '../../../stores/photoboothStore'
import { CreditCard, QrCode, Smartphone, CheckCircle, Loader2 } from 'lucide-react'

export default function PaymentScreen() {
  const { setStep, selectedTemplate, createSession, isLoading } = usePhotoboothStore()
  const [paymentMethod, setPaymentMethod] = useState<'qr' | 'ewallet' | 'card' | null>(null)
  const [processing, setProcessing] = useState(false)
  const [paid, setPaid] = useState(false)

  const handlePayment = async (method: string) => {
    setPaymentMethod(method as any)
    setProcessing(true)

    // Simulate payment processing
    await new Promise(resolve => setTimeout(resolve, 2000))

    // In production, this would call the payment provider API
    // For now, we simulate success
    setProcessing(false)
    setPaid(true)

    // Create session after payment
    await createSession('demo-token')

    setTimeout(() => {
      setStep('CAPTURE')
    }, 1500)
  }

  const handleBack = () => {
    setStep('TEMPLATE_SELECT')
  }

  if (paid) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-green-50 to-emerald-100">
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          className="text-center"
        >
          <CheckCircle size={80} className="mx-auto text-green-500 mb-4" />
          <h2 className="text-3xl font-bold text-green-800 mb-2">Payment Successful!</h2>
          <p className="text-green-600">Preparing your photobooth...</p>
        </motion.div>
      </div>
    )
  }

  if (processing) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="text-center"
        >
          <Loader2 size={60} className="mx-auto text-blue-500 animate-spin mb-4" />
          <h2 className="text-2xl font-bold text-blue-800 mb-2">Processing Payment...</h2>
          <p className="text-blue-600">Please wait</p>
        </motion.div>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-pink-50 via-purple-50 to-blue-50 p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="max-w-md w-full"
      >
        <button onClick={handleBack} className="text-gray-600 mb-6 flex items-center gap-2">
          ← Back
        </button>

        <div className="card text-center mb-6">
          <h2 className="text-2xl font-bold mb-2">Complete Payment</h2>
          {selectedTemplate && (
            <div className="bg-gray-50 rounded-2xl p-4 mt-4">
              <p className="font-semibold">{selectedTemplate.name}</p>
              <p className="text-2xl font-bold text-primary mt-2">Rp {selectedTemplate.price?.toLocaleString()}</p>
            </div>
          )}
        </div>

        <div className="space-y-3">
          <button
            onClick={() => handlePayment('qr')}
            className="w-full card flex items-center gap-4 hover:border-primary transition"
          >
            <QrCode size={32} className="text-purple-500" />
            <div className="text-left">
              <p className="font-semibold">QRIS / QR Payment</p>
              <p className="text-sm text-gray-500">Scan with any e-wallet</p>
            </div>
          </button>

          <button
            onClick={() => handlePayment('ewallet')}
            className="w-full card flex items-center gap-4 hover:border-primary transition"
          >
            <Smartphone size={32} className="text-green-500" />
            <div className="text-left">
              <p className="font-semibold">E-Wallet</p>
              <p className="text-sm text-gray-500">GoPay, OVO, Dana, LinkAja</p>
            </div>
          </button>

          <button
            onClick={() => handlePayment('card')}
            className="w-full card flex items-center gap-4 hover:border-primary transition"
          >
            <CreditCard size={32} className="text-blue-500" />
            <div className="text-left">
              <p className="font-semibold">Credit / Debit Card</p>
              <p className="text-sm text-gray-500">Visa, Mastercard, JCB</p>
            </div>
          </button>
        </div>

        <p className="text-center text-xs text-gray-400 mt-6">
          🔒 Secure payment powered by Klikku
        </p>
      </motion.div>
    </div>
  )
}
