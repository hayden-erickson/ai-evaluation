const API_BASE: string = (import.meta as any).env?.VITE_API_BASE || '/api'

export type User = { id: number; profile_image_url?: string; name: string; time_zone: string; phone_number: string; created_at: string }
export type LoginResponse = { token: string; user: User }
export type Habit = { id: number; user_id: number; name: string; description?: string; duration_seconds?: number | null; created_at: string }
export type Log = { id: number; habit_id: number; notes?: string; duration_seconds?: number | null; created_at: string }

export function authHeaders(token?: string): Record<string, string> {
  return token ? { Authorization: `Bearer ${token}` } : {}
}

type SimpleInit = Omit<RequestInit, 'headers' | 'body'> & { headers?: Record<string, string>; body?: any }

async function request<T>(path: string, opts: SimpleInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(opts.headers || {}),
  }
  const res = await fetch(`${API_BASE}${path}`, {
    ...opts,
    headers,
    body: typeof opts.body === 'string' ? opts.body : opts.body != null ? JSON.stringify(opts.body) : undefined,
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Request failed: ${res.status}`)
  }
  if (res.status === 204) return undefined as unknown as T
  return res.json()
}

export const api = {
  register: (body: { profile_image_url?: string; name: string; time_zone: string; phone_number: string; password: string }) =>
    request<User>('/users/register', { method: 'POST', body: JSON.stringify(body) }),
  login: (body: { phone_number: string; password: string }) =>
    request<LoginResponse>('/users/login', { method: 'POST', body: JSON.stringify(body) }),
  getHabits: (token: string) => request<Habit[]>('/habits', { headers: authHeaders(token) }),
  createHabit: (token: string, body: { name: string; description?: string; duration_seconds?: number | null }) =>
    request<Habit>('/habits', { method: 'POST', headers: authHeaders(token), body: JSON.stringify(body) }),
  updateHabit: (token: string, id: number, body: Partial<{ name: string; description: string; duration_seconds: number | null }>) =>
    request<Habit>(`/habits/${id}`, { method: 'PUT', headers: authHeaders(token), body: JSON.stringify(body) }),
  deleteHabit: (token: string, id: number) => request<void>(`/habits/${id}`, { method: 'DELETE', headers: authHeaders(token) }),
  getLogs: (token: string, habitId: number) => request<Log[]>(`/habits/${habitId}/logs`, { headers: authHeaders(token) }),
  createLog: (token: string, habitId: number, body: { notes?: string; duration_seconds?: number | null }) =>
    request<Log>(`/habits/${habitId}/logs`, { method: 'POST', headers: authHeaders(token), body: JSON.stringify(body) }),
  updateLog: (token: string, id: number, body: Partial<{ notes: string; duration_seconds: number | null }>) =>
    request<Log>(`/logs/${id}`, { method: 'PUT', headers: authHeaders(token), body: JSON.stringify(body) }),
  deleteLog: (token: string, id: number) => request<void>(`/logs/${id}`, { method: 'DELETE', headers: authHeaders(token) }),
}
