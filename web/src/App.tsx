import { useEffect, useMemo, useState } from 'react'
import Login from './components/Login'
import AddNewHabitButton from './components/AddNewHabitButton'
import HabitList from './components/HabitList'
import HabitDetailsModal from './components/HabitDetailsModal'
import LogDetailsModal from './components/LogDetailsModal'
import { useAuth } from './hooks/useAuth'
import { api, Habit, Log } from './api/client'

export default function App() {
  const { token, user, logout } = useAuth()
  const [habits, setHabits] = useState<Habit[]>([])
  const [editingHabit, setEditingHabit] = useState<Habit | null>(null)
  const [logsByHabit, setLogsByHabit] = useState<Record<number, Log[]>>({})
  const [editingLog, setEditingLog] = useState<Log | { habit_id: number } | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!token) return
    loadHabits()
  }, [token])

  async function loadHabits() {
    if (!token) return
    try {
      const list = await api.getHabits(token)
      setHabits(list)
      for (const h of list) {
        const logs = await api.getLogs(token, h.id)
        setLogsByHabit((prev: Record<number, Log[]>) => ({ ...prev, [h.id]: logs }))
      }
    } catch (e: any) {
      setError(e.message)
    }
  }

  if (!token || !user) return <Login />

  return (
    <div className="container">
      <div className="header">
        <div>
          <h2>Habit Streaks</h2>
          <div className="muted">Signed in as {user.name}</div>
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          <AddNewHabitButton onClick={() => setEditingHabit({ id: 0, user_id: user.id, name: '', description: '', duration_seconds: null, created_at: new Date().toISOString() })} />
          <button className="button secondary" onClick={logout}>Logout</button>
        </div>
      </div>

      {error && <div className="error" style={{ marginBottom: 12 }}>{error}</div>}

      <HabitList
        habits={habits}
        logsByHabit={logsByHabit}
        onEditHabit={(h: Habit) => setEditingHabit(h)}
        onDeleteHabit={async (h: Habit) => {
          try {
            await api.deleteHabit(token, h.id)
            setHabits((prev: Habit[]) => prev.filter((x: Habit) => x.id !== h.id))
          } catch (e: any) {
            setError(e.message)
          }
        }}
        onNewLog={(h: Habit) => setEditingLog({ habit_id: h.id })}
        onEditLog={(log: Log) => setEditingLog(log)}
      />

      {editingHabit && (
        <HabitDetailsModal
          initial={editingHabit.id ? editingHabit : null}
          onClose={() => setEditingHabit(null)}
          onSave={async (data: Partial<{ name: string; description: string; duration_seconds: number | null }>) => {
            try {
              if (editingHabit.id) {
                const updated = await api.updateHabit(token, editingHabit.id, data)
                setHabits((prev: Habit[]) => prev.map((h: Habit) => (h.id === updated.id ? updated : h)))
              } else {
                const created = await api.createHabit(token, data as { name: string; description?: string; duration_seconds?: number | null })
                setHabits((prev: Habit[]) => [created, ...prev])
                const logs = await api.getLogs(token, created.id)
                setLogsByHabit((p: Record<number, Log[]>) => ({ ...p, [created.id]: logs }))
              }
              setEditingHabit(null)
            } catch (e: any) {
              setError(e.message)
            }
          }}
        />
      )}

      {editingLog && (
        <LogDetailsModal
          habit={habits.find((h: Habit) => h.id === (editingLog as any).habit_id)}
          initial={('id' in editingLog) ? (editingLog as Log) : null}
          onClose={() => setEditingLog(null)}
          onSave={async (body: Partial<{ notes: string; duration_seconds: number | null }>) => {
            try {
              if ('id' in (editingLog as any)) {
                const updated = await api.updateLog(token, (editingLog as Log).id, body)
                const hid = updated.habit_id
                setLogsByHabit((prev: Record<number, Log[]>) => ({ ...prev, [hid]: (prev[hid] || []).map((l: Log) => (l.id === updated.id ? updated : l)) }))
              } else {
                const hid = (editingLog as { habit_id: number }).habit_id
                const created = await api.createLog(token, hid, body)
                setLogsByHabit((prev: Record<number, Log[]>) => ({ ...prev, [hid]: [created, ...(prev[hid] || [])] }))
              }
              setEditingLog(null)
            } catch (e: any) {
              setError(e.message)
            }
          }}
        />
      )}
    </div>
  )
}
