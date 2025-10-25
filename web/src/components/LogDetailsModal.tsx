import type React from 'react'
import { useState } from 'react'
import { Habit, Log } from '../api/client'

type LogBody = { notes?: string; duration_seconds?: number | null }

type Props = {
  habit: Habit | undefined
  initial: Log | null
  onClose: () => void
  onSave: (body: LogBody) => Promise<void> | void
}

export default function LogDetailsModal({ habit, initial, onClose, onSave }: Props) {
  const [notes, setNotes] = useState<string>(initial?.notes || '')
  const [duration, setDuration] = useState<string>(
    initial?.duration_seconds != null ? String(initial.duration_seconds) : '',
  )
  const isEditing = Boolean(initial && initial.id)
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  async function submit() {
    setError(null)
    const body: LogBody = {}
    if (notes.trim()) body.notes = notes.trim()
    if (habit?.duration_seconds != null || duration !== '') {
      const n = Number(duration)
      if (Number.isNaN(n) || n < 0) {
        setError('Duration must be a non-negative number')
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
        <div className="header"><h3>{isEditing ? 'Edit Log' : 'New Log'}</h3></div>
        <div className="grid">
          <textarea className="textarea" placeholder="Notes" value={notes} onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setNotes(e.target.value)} />
          <input className="input" placeholder="Duration seconds" value={duration} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDuration(e.target.value)} />
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
