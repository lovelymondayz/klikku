import React from 'react'
import { AnimatePresence } from 'framer-motion'
import { usePhotoboothStore } from '../../stores/photoboothStore'
import AttractScreen from './components/AttractScreen'
import TemplateSelection from './components/TemplateSelection'
import PaymentScreen from './components/PaymentScreen'
import CameraCapture from './components/CameraCapture'
import PhotoReview from './components/PhotoReview'
import EmailQRScreen from './components/EmailQRScreen'
import PromotionScreen from './components/PromotionScreen'

export default function PhotoboothPage() {
  const { currentStep } = usePhotoboothStore()

  return (
    <div className="w-full h-full bg-gray-900 overflow-hidden">
      <AnimatePresence mode="wait">
        {currentStep === 'attract' && <AttractScreen key="attract" />}
        {currentStep === 'template-selection' && <TemplateSelection key="template-selection" />}
        {currentStep === 'payment' && <PaymentScreen key="payment" />}
        {currentStep === 'camera-capture' && <CameraCapture key="camera-capture" />}
        {currentStep === 'photo-review' && <PhotoReview key="photo-review" />}
        {currentStep === 'email-qr' && <EmailQRScreen key="email-qr" />}
        {currentStep === 'promotion' && <PromotionScreen key="promotion" />}
      </AnimatePresence>
    </div>
  )
}