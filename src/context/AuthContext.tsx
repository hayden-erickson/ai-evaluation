import React, { createContext, useContext, useEffect, useMemo, useState } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { apiFetch, setToken as persistToken } from '../api/client';
import type { LoginResponse, User } from '../types';

interface AuthState {
  user: User | null;
  token: string | null;
  loading: boolean;
  error: string | null;
  login: (phone: string, password: string) => Promise<void>;
  register: (name: string, phone: string, password: string, timeZone: string, profileImageUrl?: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthState | undefined>(undefined);

/** Provides authentication state and actions to the app. */
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Hydrate token/user from storage on boot
  useEffect(() => {
    (async () => {
      try {
        const [t, u] = await Promise.all([
          AsyncStorage.getItem('auth_token'),
          AsyncStorage.getItem('auth_user'),
        ]);
        if (t) setToken(t);
        if (u) setUser(JSON.parse(u));
      } catch (e) {
        // If hydration fails, continue without crashing
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  /** Attempts to login and persists token/user on success. */
  const login = async (phone: string, password: string) => {
    setError(null);
    try {
      const resp = await apiFetch<LoginResponse>('/users/login', {
        method: 'POST',
        body: JSON.stringify({ phone_number: phone, password }),
      });
      setToken(resp.token);
      setUser(resp.user);
      await persistToken(resp.token);
      await AsyncStorage.setItem('auth_user', JSON.stringify(resp.user));
    } catch (e: any) {
      setError(e.message || 'Login failed');
      throw e;
    }
  };

  /** Registers then logs in the user for a smooth UX. */
  const register = async (
    name: string,
    phone: string,
    password: string,
    timeZone: string,
    profileImageUrl?: string,
  ) => {
    setError(null);
    try {
      await apiFetch('/users/register', {
        method: 'POST',
        body: JSON.stringify({
          name,
          phone_number: phone,
          password,
          time_zone: timeZone,
          profile_image_url: profileImageUrl,
        }),
      });
      await login(phone, password);
    } catch (e: any) {
      setError(e.message || 'Registration failed');
      throw e;
    }
  };

  /** Clears token/user and storage. */
  const logout = async () => {
    setUser(null);
    setToken(null);
    await persistToken(null);
    await AsyncStorage.removeItem('auth_user');
  };

  const value = useMemo(
    () => ({ user, token, loading, error, login, register, logout }),
    [user, token, loading, error],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

/** Hook to access auth state. */
export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
