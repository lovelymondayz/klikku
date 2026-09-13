import React from 'react'
import { motion } from 'framer-motion'
import { usePhotoboothStore, MOCK_TEMPLATES } from '../../../stores/photoboothStore'
import { ArrowLeft, Camera, Image, Grid3X3, Columns, LayoutGrid } from 'lucide-react'

const LAYOUT_ICONS: Record<string, React.ReactNode> = {
  strip: <Columns size={32} />,
  square: <Grid3X3 size={32} />,
  grid: <LayoutGrid size={32} />,
  duo: <Image size={32} />,
}

export default function TemplateSelection() {
  const { setStep, selectTemplate, selectedTemplate } = usePhotoboothStore()

  const handleSelect = (template: typeof MOCK_TEMPLATES[0]) => {
    selectTemplate(template)
  }

  const handleConfirm = () => {
    if (selectedTemplate) {
      setStep('payment')
    }
  }

  const handleBack = () => {
    setStep('attract')
  }

  return (
    <motion.div
      className="relative w-full h-full overflow-hidden bg-gradient-to-br from-bg to-primary-subtle"
      initial={{ opacity: 0, x: 100 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -100 }}
      transition={{ duration: 0.5 }}
    >
      {/* Background decoration */}
      <div className="absolute top-0 right-0 w-64 h-64 bg-gradient-to-bl from-primary-subtle rounded-full -translate-y-1/2 translate-x-1/2" />
      <div className="absolute bottom-0 left-0 w-48 h-48 bg-gradient-to-tr from-primary-subtle rounded-full translate-y-1/2 -translate-x-1/2" />

      <div className="relative z-10 flex flex-col h-full p-6 md:p-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <motion.button
            onClick={handleBack}
            className="w-14 h-14 rounded-full bg-surface shadow-md flex items-center justify-center active:scale-90 transition-transform"
            whileTap={{ scale: 0.9 }}
          >
            <ArrowLeft size={24} className="text-text" />
          </motion.button>

          <h1 className="text-2xl md:text-3xl font-bold text-text">
            Choose Your Style
          </h1>

          <div className="w-14" /> {/* Spacer */}
        </div>

        {/* Template Grid */}
        <div className="flex-1 grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-6 content-center">
          {MOCK_TEMPLATES.map((template, index) => (
            <motion.button
              key={template.id}
              onClick={() => handleSelect(template)}
              className={`relative flex flex-col items-center justify-center p-6 rounded-lg border-4 transition-colors duration-base min-h-[200px] ${
                selectedTemplate?.id === template.id
                  ? 'border-primary bg-surface shadow-xl shadow-primary/20'
                  : 'border-transparent bg-surface/80 shadow-md hover:shadow-xl'
              }`}
              initial={{ opacity: 0, y: 30 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1, duration: 0.4 }}
              whileHover={{ scale: 1.03 }}
              whileTap={{ scale: 0.97 }}
            >
              {/* Selected indicator */}
              {selectedTemplate?.id === template.id && (
                <motion.div
                  className="absolute -top-2 -right-2 w-8 h-8 bg-primary rounded-full flex items-center justify-center"
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  transition={{ type: 'spring', stiffness: 500 }}
                >
                  <span className="text-text text-lg">✓</span>
                </motion.div>
              )}

              {/* Layout icon */}
              <div className={`mb-4 p-4 rounded-lg ${
                selectedTemplate?.id === template.id
                  ? 'bg-primary/10 text-primary'
                  : 'bg-surface-alt text-text-muted'
              }`}>
                {LAYOUT_ICONS[template.layout_config?.output_width > 1000 ? 'strip' : 'square'] || <Camera size={32} />}
              </div>

              {/* Template name */}
              <h3 className="text-lg font-bold text-text mb-1">
                {template.name}
              </h3>

              {/* Photo count */}
              <div className="flex items-center gap-1 text-sm font-semibold text-primary">
                <Camera size={14} />
                <span>{template.photo_count} photos</span>
              </div>

              {/* Price */}
              <div className="mt-2 px-4 py-1 bg-accent/20 rounded-full">
                <span className="text-sm font-bold text-text">
                  ${template.price.toFixed(2)}
                </span>
              </div>
            </motion.button>
          ))}
        </div>

        {/* Bottom CTA */}
        <motion.button
          onClick={handleConfirm}
          disabled={!selectedTemplate}
          className={`mt-6 w-full py-5 rounded-lg text-xl font-bold shadow-md transition-colors duration-base ${
            selectedTemplate
              ? 'bg-primary text-text active:scale-95'
              : 'bg-border text-text-subtle cursor-not-allowed'
          }`}
          whileHover={selectedTemplate ? { scale: 1.02 } : {}}
          whileTap={selectedTemplate ? { scale: 0.98 } : {}}
        >
          {selectedTemplate
            ? `Continue with ${selectedTemplate.name} — $${selectedTemplate.price.toFixed(2)}`
            : 'Select a template to continue'}
        </motion.button>
      </div>
    </motion.div>
  )
}