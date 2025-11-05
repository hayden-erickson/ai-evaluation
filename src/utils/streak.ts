import type { LogEntry } from '../types';

/** Formats a Date as YYYY-MM-DD in local time for day comparisons. */
export function dayKey(d: Date): string {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${y}-${m}-${day}`;
}

/** Returns the number of whole days between two dates in local time. */
export function daysBetween(a: Date, b: Date): number {
  const da = new Date(a.getFullYear(), a.getMonth(), a.getDate());
  const db = new Date(b.getFullYear(), b.getMonth(), b.getDate());
  const diff = da.getTime() - db.getTime();
  return Math.round(diff / (1000 * 60 * 60 * 24));
}

/** Indexes logs by local day key for quick lookups. */
export function indexLogsByDay(logs: LogEntry[]): Record<string, LogEntry> {
  const map: Record<string, LogEntry> = {};
  for (const l of logs) {
    const d = new Date(l.created_at);
    map[dayKey(d)] = l;
  }
  return map;
}

/**
 * Computes the current streak length with a rule that allows skipping a single day
 * between consecutive logged days without breaking the streak. Skipped days do not
 * increment the streak count themselves; only logged days count.
 */
export function computeStreak(logs: LogEntry[], today = new Date()): number {
  if (!logs.length) return 0;
  const byDay = indexLogsByDay(logs);

  // Start from today or yesterday if today isn't logged but yesterday is
  let current = new Date(today.getFullYear(), today.getMonth(), today.getDate());
  let count = 0;

  // If neither today nor yesterday has a log, there's no ongoing streak.
  if (byDay[dayKey(current)]) {
    count++;
  } else {
    const y = new Date(current);
    y.setDate(y.getDate() - 1);
    if (byDay[dayKey(y)]) {
      current = y;
      count++;
    } else {
      return 0;
    }
  }

  // Continue walking back, allowing single-day gaps anywhere in sequence
  // Stop when the gap between the last considered day and the previous log > 2 days
  while (true) {
    const d1 = new Date(current);
    d1.setDate(d1.getDate() - 1);
    const d2 = new Date(current);
    d2.setDate(d2.getDate() - 2);

    if (byDay[dayKey(d1)]) {
      current = d1;
      count++;
      continue;
    }
    if (byDay[dayKey(d2)]) {
      current = d2; // consume one skip day
      count++;
      continue;
    }
    break;
  }

  return count;
}
