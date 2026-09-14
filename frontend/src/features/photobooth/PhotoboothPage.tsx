import React from 'react'
import { AnimatePresence } from 'framer-motion'
import { usePhotoboothStore, STEPS } from '../../stores/photoboothStore'
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
    <div className="w-full h-full bg-surface overflow-hidden">
      <AnimatePresence mode="wait">
        {currentStep === STEPS.IDLE && <AttractScreen key={STEPS.IDLE} />}
        {currentStep === STEPS.TEMPLATE_SELECT && <TemplateSelection key={STEPS.TEMPLATE_SELECT} />}
        {currentStep === STEPS.PAYMENT && <PaymentScreen key={STEPS.PAYMENT} />}
        {currentStep === STEPS.CAPTURE && <CameraCapture key={STEPS.CAPTURE} />}
        {currentStep === STEPS.PHOTO_REVIEW && <PhotoReview key={STEPS.PHOTO_REVIEW} />}
        {currentStep === STEPS.EMAIL_QR && <EmailQRScreen key={STEPS.EMAIL_QR} />}
        {currentStep === STEPS.PROMOTION && <PromotionScreen key={STEPS.PROMOTION} />}
      </AnimatePresence>
    </div>
  )
}