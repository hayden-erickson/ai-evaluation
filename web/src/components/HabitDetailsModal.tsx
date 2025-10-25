import type React from 'react'
import { useEffect, useState } from 'react'
import { Habit } from '../api/client'

type HabitBody = { name?: string; description?: string; duration_seconds?: number | null }

type Props = {
  initial: Habit | null
  onClose: () => void
  onSave: (body: HabitBody) => Promise<void> | void
}

export default function HabitDetailsModal({ initial, onClose, onSave }: Props) {
  const [name, setName] = useState(initial?.name || '')
  const [description, setDescription] = useState(initial?.description || '')
  const [duration, setDuration] = useState<string>(
    initial?.duration_seconds != null ? String(initial.duration_seconds) : '',
  )
  const isEditing = Boolean(initial && initial.id)
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  async function submit() {
    setError(null)
    if (!name.trim()) {
      setError('Name is required')
      return
    }
    const body: HabitBody = { name: name.trim() }
    if (description.trim()) body.description = description.trim()
    if (duration !== '') {
      const n = Number(duration)
      if (Number.isNaN(n) || n < 0) {
        setError('Duration must be a non-negative number or empty')
        return
      }
      body.duration_seconds = n
    } else {
      body.duration_seconds = null
    }
    setSaving(true)
    try {
      await onSave(body)
    } catch (e: any) {
      setError(e.message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="card modal" onClick={(e: React.MouseEvent<HTMLDivElement>) => e.stopPropagation()}>
        <div className="header"><h3>{isEditing ? 'Edit Habit' : 'New Habit'}</h3></div>
        <div className="grid">
          <input className="input" placeholder="Name" value={name} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setName(e.target.value)} />
          <textarea className="textarea" placeholder="Description" value={description} onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setDescription(e.target.value)} />
          <input className="input" placeholder="Duration seconds (optional)" value={duration} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDuration(e.target.value)} />
          {error && <div className="error">{error}</div>}
          <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end' }}>
            <button className="button secondary" onClick={onClose} disabled={saving}>Cancel</button>
            <button className="button" onClick={submit} disabled={saving}>{saving ? 'Saving...' : 'Save'}</button>
          </div>
        </div>
      </div>
    </div>
  )
}
