import { Log } from '../api/client'

export function computeStreak(logs: Log[], now = new Date()): number {
  const days = new Set(
    logs.map((l) => new Date(l.created_at)).map((d) => new Date(Date.UTC(d.getUTCFullYear(), d.getUTCMonth(), d.getUTCDate())).getTime()),
  )
  let streak = 0
  let skips = 0
  let cursor = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate())).getTime()
  while (true) {
    const has = days.has(cursor)
    if (has) {
      streak += 1
    } else if (skips === 0) {
      skips = 1
    } else {
      break
    }
    cursor -= 24 * 60 * 60 * 1000
  }
  return streak
}

export function buildRecentDays(logs: Log[], daysBack = 30, now = new Date()): { ts: number; status: 'on' | 'skip' | 'off' }[] {
  const set = new Set(
    logs.map((l) => new Date(l.created_at)).map((d) => new Date(Date.UTC(d.getUTCFullYear(), d.getUTCMonth(), d.getUTCDate())).getTime()),
  )
  const list: { ts: number; status: 'on' | 'skip' | 'off' }[] = []
  let skips = 0
  let cursor = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate())).getTime()
  for (let i = 0; i < daysBack; i++) {
    if (set.has(cursor)) list.push({ ts: cursor, status: 'on' })
    else if (skips === 0) {
      list.push({ ts: cursor, status: 'skip' })
      skips = 1
    } else list.push({ ts: cursor, status: 'off' })
    cursor -= 24 * 60 * 60 * 1000
  }
  return list.reverse()
}
