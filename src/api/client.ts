import AsyncStorage from '@react-native-async-storage/async-storage';
import { Platform } from 'react-native';

// Base API URL. On Android emulator, localhost is 10.0.2.2
const DEFAULT_BASE_URL = Platform.OS === 'android' ? 'http://10.0.2.2:8080' : 'http://localhost:8080';
let BASE_URL = DEFAULT_BASE_URL;

/**
 * Sets the base URL for API requests at runtime.
 */
export function setBaseUrl(url: string) {
  BASE_URL = url || DEFAULT_BASE_URL;
}

/**
 * Reads the stored JWT token.
 */
export async function getToken(): Promise<string | null> {
  return AsyncStorage.getItem('auth_token');
}

/**
 * Stores the JWT token securely.
 */
export async function setToken(token: string | null): Promise<void> {
  if (token) {
    await AsyncStorage.setItem('auth_token', token);
  } else {
    await AsyncStorage.removeItem('auth_token');
  }
}

/**
 * Wrapper around fetch that injects Authorization header and parses JSON.
 * Ensures consistent error handling and timeouts.
 */
export async function apiFetch<T = any>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 15000);

  try {
    const token = await getToken();
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    };
    if (token) headers['Authorization'] = `Bearer ${token}`;

    const res = await fetch(`${BASE_URL}${path}`, {
      ...options,
      headers,
      signal: controller.signal,
    });

    const text = await res.text();
    const data = text ? JSON.parse(text) : null;

    if (!res.ok) {
      const message = (data && (data.error || data.message)) || res.statusText;
      throw new Error(message);
    }

    return data as T;
  } catch (e: any) {
    if (e.name === 'AbortError') {
      throw new Error('Request timed out');
    }
    throw e;
  } finally {
    clearTimeout(timeout);
  }
}
