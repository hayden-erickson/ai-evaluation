import { useEffect, useMemo, useState } from 'react'
import { api, LoginResponse, User } from '../api/client'

export function useAuth() {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'))
  const [user, setUser] = useState<User | null>(() => {
    const raw = localStorage.getItem('user')
    return raw ? (JSON.parse(raw) as User) : null
  })
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  function saveAuth(resp: LoginResponse) {
    localStorage.setItem('token', resp.token)
    localStorage.setItem('user', JSON.stringify(resp.user))
    setToken(resp.token)
    setUser(resp.user)
  }
  function logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setToken(null)
    setUser(null)
  }

  async function login(phone_number: string, password: string) {
    setError(null)
    setLoading(true)
    try {
      const resp = await api.login({ phone_number, password })
      saveAuth(resp)
    } catch (e: any) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  async function register(data: { name: string; phone_number: string; password: string; time_zone: string; profile_image_url?: string }) {
    setError(null)
    setLoading(true)
    try {
      await api.register(data)
      const resp = await api.login({ phone_number: data.phone_number, password: data.password })
      saveAuth(resp)
    } catch (e: any) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  return useMemo(() => ({ token, user, error, loading, login, register, logout }), [token, user, error, loading])
}
