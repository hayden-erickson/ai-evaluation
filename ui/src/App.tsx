import { useState } from "react";
import "./App.css";
import { AuthProvider, useAuth } from "./lib/auth";
import { HabitUI } from "./components/HabitUI";

function Login() {
  const { login } = useAuth();
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  return (
    <div className="login">
      <h1>Sign in</h1>
      {error && <p className="error">{error}</p>}
      <div className="form">
        <label>
          Phone number
          <input value={phone} onChange={(e) => setPhone(e.target.value)} placeholder="+1234567890" />
        </label>
        <label>
          Password
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" />
        </label>
      </div>
      <button
        className="primary"
        disabled={loading || !phone.trim() || !password.trim()}
        onClick={async () => {
          setLoading(true);
          setError(null);
          try {
            await login({ phone_number: phone.trim(), password: password });
          } catch (e: any) {
            setError(e.message || "Login failed");
          } finally {
            setLoading(false);
          }
        }}
      >
        Sign in
      </button>
    </div>
  );
}

function Shell() {
  const { token, logout } = useAuth();
  if (!token) return <Login />;
  return (
    <div>
      <nav className="topbar">
        <div className="brand">Habit Streaks</div>
        <button className="ghost" onClick={logout}>Sign out</button>
      </nav>
      <HabitUI />
    </div>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <Shell />
    </AuthProvider>
  );
}

