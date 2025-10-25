import { Habit, Log } from '../api/client'
import { buildRecentDays, computeStreak } from '../utils/streak'

type Props = {
  habits: Habit[]
  logsByHabit: Record<number, Log[]>
  onEditHabit: (habit: Habit) => void
  onDeleteHabit: (habit: Habit) => void
  onNewLog: (habit: Habit) => void
  onEditLog: (log: Log) => void
}

export default function HabitList({ habits, logsByHabit, onEditHabit, onDeleteHabit, onNewLog, onEditLog }: Props) {
  return (
    <div className="habits">
      {habits.map((h) => {
        const logs = logsByHabit[h.id] || []
        const streak = computeStreak(logs)
        const days = buildRecentDays(logs, 30)
        return (
          <div key={h.id} className="card">
            <div className="header" style={{ marginBottom: 8 }}>
              <div>
                <div style={{ fontWeight: 700 }}>{h.name}</div>
                {h.description && <div className="muted" style={{ marginTop: 4 }}>{h.description}</div>}
              </div>
              <div style={{ display: 'flex', gap: 8 }}>
                <button className="button secondary" onClick={() => onEditHabit(h)}>Edit</button>
                <button className="button danger" onClick={() => onDeleteHabit(h)}>Delete</button>
              </div>
            </div>
            <div className="streak">Current streak: {streak} day{streak === 1 ? '' : 's'} (one skip allowed)</div>
            <div className="streak-list">
              {days.map((d) => (
                <div key={d.ts} className={`day ${d.status === 'on' ? 'on' : d.status === 'skip' ? 'skip' : ''}`}></div>
              ))}
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 12 }}>
              <div style={{ fontSize: 12 }} className="muted">{h.duration_seconds != null ? 'Requires duration logs' : 'Free-form logs'}</div>
              <button className="button" onClick={() => onNewLog(h)}>New Log</button>
            </div>
            {logs.length > 0 && (
              <div style={{ marginTop: 12 }}>
                <div className="muted" style={{ marginBottom: 6, fontSize: 12 }}>Recent logs</div>
                <div className="grid">
                  {logs.slice(0, 5).map((l) => (
                    <div key={l.id} className="card" style={{ padding: 10, cursor: 'pointer' }} onClick={() => onEditLog(l)}>
                      <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                        <div style={{ fontWeight: 600 }}>{new Date(l.created_at).toLocaleString()}</div>
                        {l.duration_seconds != null && <div className="muted">{l.duration_seconds}s</div>}
                      </div>
                      {l.notes && <div className="muted" style={{ marginTop: 4 }}>{l.notes}</div>}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )
      })}
    </div>
  )
}
