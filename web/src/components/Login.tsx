import { useState } from 'react'
import { useAuth } from '../hooks/useAuth'

export default function Login() {
  const { login, register, error, loading } = useAuth()
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [phone, setPhone] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [tz, setTz] = useState(Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC')

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (mode === 'login') await login(phone, password)
    else await register({ name, phone_number: phone, password, time_zone: tz })
  }

  return (
    <div className="container">
      <div className="card" style={{ maxWidth: 480, margin: '72px auto' }}>
        <div className="header"><h2>Habit Streaks</h2></div>
        <form onSubmit={onSubmit} className="grid">
          {mode === 'register' && (
            <input className="input" placeholder="Name" value={name} onChange={(e) => setName(e.target.value)} required />
          )}
          <input className="input" placeholder="Phone number" value={phone} onChange={(e) => setPhone(e.target.value)} required />
          <input className="input" type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} required />
          {mode === 'register' && (
            <input className="input" placeholder="Time zone" value={tz} onChange={(e) => setTz(e.target.value)} required />
          )}
          {error && <div className="error">{error}</div>}
          <button className="button" disabled={loading} type="submit">{mode === 'login' ? 'Login' : 'Create account'}</button>
        </form>
        <div style={{ marginTop: 12 }}>
          <button className="button secondary" onClick={() => setMode(mode === 'login' ? 'register' : 'login')}>
            {mode === 'login' ? 'Need an account? Register' : 'Have an account? Login'}
          </button>
        </div>
      </div>
    </div>
  )
}
