import React, { useState, useRef, useEffect } from 'react'
import { api } from '../../lib/api'
import { Save, Plus, Trash2, Move, Image, Type, RotateCcw } from 'lucide-react'

interface PhotoSlot {
  id: string
  x: number
  y: number
  width: number
  height: number
  rotation: number
}

interface TextOverlay {
  id: string
  text: string
  x: number
  y: number
  fontSize: number
  color: string
}

interface TemplateEditorProps {
  templateId?: string
  onSave: () => void
}

export default function TemplateEditor({ templateId, onSave }: TemplateEditorProps) {
  const canvasRef = useRef<HTMLDivElement>(null)
  const [template, setTemplate] = useState({
    name: '',
    photo_count: 4,
    output_width: 1200,
    output_height: 1800,
    background_color: '#ffffff',
    overlay_url: '',
  })
  const [photoSlots, setPhotoSlots] = useState<PhotoSlot[]>([])
  const [textOverlays, setTextOverlays] = useState<TextOverlay[]>([])
  const [selectedSlot, setSelectedSlot] = useState<string | null>(null)
  const [selectedText, setSelectedText] = useState<string | null>(null)
  const [dragging, setDragging] = useState<{ id: string; type: 'slot' | 'text'; offsetX: number; offsetY: number } | null>(null)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (templateId) {
      api.get(`/templates/${templateId}`).then((res) => {
        const data = res.data.data
        setTemplate(data)
        if (data.layout_config?.photoSlots) {
          setPhotoSlots(data.layout_config.photoSlots)
        }
        if (data.layout_config?.textOverlays) {
          setTextOverlays(data.layout_config.textOverlays)
        }
      })
    } else {
      // Initialize with default slots based on photo_count
      initDefaultSlots(template.photo_count)
    }
  }, [templateId])

  const initDefaultSlots = (count: number) => {
    const slots: PhotoSlot[] = []
    const cols = count <= 2 ? 1 : 2
    const rows = Math.ceil(count / cols)
    const slotWidth = 400 / cols
    const slotHeight = 400 / rows

    for (let i = 0; i < count; i++) {
      const col = i % cols
      const row = Math.floor(i / cols)
      slots.push({
        id: `slot_${i}`,
        x: 50 + col * slotWidth,
        y: 50 + row * slotHeight,
        width: slotWidth - 20,
        height: slotHeight - 20,
        rotation: 0,
      })
    }
    setPhotoSlots(slots)
  }

  const handleMouseDown = (e: React.MouseEvent, id: string, type: 'slot' | 'text') => {
    e.stopPropagation()
    const rect = canvasRef.current?.getBoundingClientRect()
    if (!rect) return

    if (type === 'slot') {
      const slot = photoSlots.find(s => s.id === id)
      if (slot) {
        setDragging({ id, type, offsetX: e.clientX - rect.left - slot.x, offsetY: e.clientY - rect.top - slot.y })
        setSelectedSlot(id)
        setSelectedText(null)
      }
    } else {
      const text = textOverlays.find(t => t.id === id)
      if (text) {
        setDragging({ id, type, offsetX: e.clientX - rect.left - text.x, offsetY: e.clientY - rect.top - text.y })
        setSelectedText(id)
        setSelectedSlot(null)
      }
    }
  }

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!dragging) return
    const rect = canvasRef.current?.getBoundingClientRect()
    if (!rect) return

    const x = e.clientX - rect.left - dragging.offsetX
    const y = e.clientY - rect.top - dragging.offsetY

    if (dragging.type === 'slot') {
      setPhotoSlots(prev => prev.map(s => s.id === dragging.id ? { ...s, x, y } : s))
    } else {
      setTextOverlays(prev => prev.map(t => t.id === dragging.id ? { ...t, x, y } : t))
    }
  }

  const handleMouseUp = () => {
    setDragging(null)
  }

  const addTextOverlay = () => {
    const newText: TextOverlay = {
      id: `text_${Date.now()}`,
      text: 'Your Text Here',
      x: 100,
      y: 100,
      fontSize: 24,
      color: '#ffffff',
    }
    setTextOverlays(prev => [...prev, newText])
  }

  const updateText = (id: string, updates: Partial<TextOverlay>) => {
    setTextOverlays(prev => prev.map(t => t.id === id ? { ...t, ...updates } : t))
  }

  const deleteText = (id: string) => {
    setTextOverlays(prev => prev.filter(t => t.id !== id))
    setSelectedText(null)
  }

  const deleteSlot = (id: string) => {
    setPhotoSlots(prev => prev.filter(s => s.id !== id))
    setSelectedSlot(null)
  }

  const handleSave = async () => {
    setSaving(true)
    const layoutConfig = {
      photoSlots,
      textOverlays,
      backgroundColor: template.background_color,
    }

    const payload = {
      ...template,
      layout_config: layoutConfig,
    }

    try {
      if (templateId) {
        await api.put(`/templates/${templateId}`, payload)
      } else {
        await api.post('/templates', payload)
      }
      onSave()
    } catch (err) {
      console.error('Save failed:', err)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="flex h-full">
      {/* Canvas Area */}
      <div className="flex-1 flex items-center justify-center bg-gray-100 p-8">
        <div
          ref={canvasRef}
          className="relative bg-white shadow-2xl"
          style={{
            width: 500,
            height: 750,
            backgroundColor: template.background_color,
          }}
          onMouseMove={handleMouseMove}
          onMouseUp={handleMouseUp}
          onMouseLeave={handleMouseUp}
        >
          {/* Photo Slots */}
          {photoSlots.map((slot) => (
            <div
              key={slot.id}
              className={`absolute border-2 border-dashed flex items-center justify-center cursor-move ${
                selectedSlot === slot.id ? 'border-primary bg-primary/10' : 'border-gray-300 bg-gray-50'
              }`}
              style={{
                left: slot.x,
                top: slot.y,
                width: slot.width,
                height: slot.height,
                transform: `rotate(${slot.rotation}deg)`,
              }}
              onMouseDown={(e) => handleMouseDown(e, slot.id, 'slot')}
            >
              <Move size={24} className="text-gray-400" />
              {selectedSlot === slot.id && (
                <button
                  onClick={() => deleteSlot(slot.id)}
                  className="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center"
                >
                  <Trash2 size={12} />
                </button>
              )}
            </div>
          ))}

          {/* Text Overlays */}
          {textOverlays.map((text) => (
            <div
              key={text.id}
              className={`absolute cursor-move ${
                selectedText === text.id ? 'ring-2 ring-primary' : ''
              }`}
              style={{
                left: text.x,
                top: text.y,
                fontSize: text.fontSize,
                color: text.color,
              }}
              onMouseDown={(e) => handleMouseDown(e, text.id, 'text')}
            >
              {text.text}
              {selectedText === text.id && (
                <button
                  onClick={() => deleteText(text.id)}
                  className="absolute -top-2 -right-2 w-6 h-6 bg-red-500 text-white rounded-full flex items-center justify-center"
                >
                  <Trash2 size={12} />
                </button>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Sidebar */}
      <div className="w-80 bg-white border-l border-gray-200 p-6 overflow-y-auto">
        <h2 className="text-xl font-bold mb-6">Template Editor</h2>

        {/* Template Name */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">Template Name</label>
          <input
            className="input"
            value={template.name}
            onChange={(e) => setTemplate({ ...template, name: e.target.value })}
            placeholder="e.g., Classic Strip"
          />
        </div>

        {/* Photo Count */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">Photo Count</label>
          <select
            className="input"
            value={template.photo_count}
            onChange={(e) => {
              const count = parseInt(e.target.value)
              setTemplate({ ...template, photo_count: count })
              initDefaultSlots(count)
            }}
          >
            {[1, 2, 3, 4, 5, 6].map(n => <option key={n} value={n}>{n} photos</option>)}
          </select>
        </div>

        {/* Background Color */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">Background Color</label>
          <input
            type="color"
            className="w-full h-10 rounded-lg border border-gray-200"
            value={template.background_color}
            onChange={(e) => setTemplate({ ...template, background_color: e.target.value })}
          />
        </div>

        {/* Add Text Button */}
        <div className="mb-6">
          <button onClick={addTextOverlay} className="btn-secondary w-full flex items-center justify-center gap-2">
            <Type size={18} /> Add Text Overlay
          </button>
        </div>

        {/* Selected Text Properties */}
        {selectedText && (
          <div className="mb-6 p-4 bg-gray-50 rounded-xl">
            <h3 className="font-semibold mb-3">Text Properties</h3>
            <input
              className="input mb-3"
              value={textOverlays.find(t => t.id === selectedText)?.text || ''}
              onChange={(e) => updateText(selectedText, { text: e.target.value })}
              placeholder="Enter text"
            />
            <div className="flex gap-3">
              <input
                type="number"
                className="input w-1/2"
                value={textOverlays.find(t => t.id === selectedText)?.fontSize || 24}
                onChange={(e) => updateText(selectedText, { fontSize: parseInt(e.target.value) })}
                placeholder="Size"
              />
              <input
                type="color"
                className="w-1/2 h-10 rounded-lg border border-gray-200"
                value={textOverlays.find(t => t.id === selectedText)?.color || '#ffffff'}
                onChange={(e) => updateText(selectedText, { color: e.target.value })}
              />
            </div>
          </div>
        )}

        {/* Actions */}
        <div className="space-y-3">
          <button
            onClick={handleSave}
            disabled={saving || !template.name}
            className="btn-primary w-full flex items-center justify-center gap-2 disabled:opacity-50"
          >
            <Save size={18} /> {saving ? 'Saving...' : 'Save Template'}
          </button>
          <button
            onClick={onSave}
            className="btn-secondary w-full flex items-center justify-center gap-2"
          >
            <RotateCcw size={18} /> Cancel
          </button>
        </div>
      </div>
    </div>
  )
}
