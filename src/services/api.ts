import AsyncStorage from '@react-native-async-storage/async-storage';
import {config} from '../config/config';

// Configuration
const API_BASE_URL = config.apiBaseUrl;
const TOKEN_KEY = config.tokenKey;

// Types
export interface User {
  id: number;
  profile_image_url?: string;
  name: string;
  time_zone: string;
  phone_number: string;
  created_at: string;
}

export interface Habit {
  id: number;
  user_id: number;
  name: string;
  description?: string;
  duration_seconds?: number;
  created_at: string;
}

export interface Log {
  id: number;
  habit_id: number;
  notes?: string;
  duration_seconds?: number;
  created_at: string;
}

export interface LoginRequest {
  phone_number: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  time_zone: string;
  phone_number: string;
  password: string;
  profile_image_url?: string;
}

export interface CreateHabitRequest {
  name: string;
  description?: string;
  duration_seconds?: number;
}

export interface UpdateHabitRequest {
  name?: string;
  description?: string;
  duration_seconds?: number;
}

export interface CreateLogRequest {
  notes?: string;
  duration_seconds?: number;
}

export interface UpdateLogRequest {
  notes?: string;
  duration_seconds?: number;
}

// API Error class
export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

// Token management
export const getToken = async (): Promise<string | null> => {
  try {
    return await AsyncStorage.getItem(TOKEN_KEY);
  } catch (error) {
    console.error('Error getting token:', error);
    return null;
  }
};

export const setToken = async (token: string): Promise<void> => {
  try {
    await AsyncStorage.setItem(TOKEN_KEY, token);
  } catch (error) {
    console.error('Error setting token:', error);
    throw error;
  }
};

export const removeToken = async (): Promise<void> => {
  try {
    await AsyncStorage.removeItem(TOKEN_KEY);
  } catch (error) {
    console.error('Error removing token:', error);
    throw error;
  }
};

// Generic fetch wrapper with authentication
const fetchWithAuth = async (
  endpoint: string,
  options: RequestInit = {},
): Promise<Response> => {
  const token = await getToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new ApiError(response.status, errorText || response.statusText);
  }

  return response;
};

// Authentication API
export const login = async (
  credentials: LoginRequest,
): Promise<{token: string; user: User}> => {
  const response = await fetch(`${API_BASE_URL}/users/login`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(credentials),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new ApiError(response.status, errorText || 'Login failed');
  }

  const data = await response.json();
  await setToken(data.token);
  return data;
};

export const register = async (
  userData: RegisterRequest,
): Promise<{token: string; user: User}> => {
  const response = await fetch(`${API_BASE_URL}/users/register`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new ApiError(response.status, errorText || 'Registration failed');
  }

  const data = await response.json();
  await setToken(data.token);
  return data;
};

export const logout = async (): Promise<void> => {
  await removeToken();
};

// Habits API
export const getHabits = async (): Promise<Habit[]> => {
  const response = await fetchWithAuth('/habits');
  return response.json();
};

export const getHabit = async (habitId: number): Promise<Habit> => {
  const response = await fetchWithAuth(`/habits/${habitId}`);
  return response.json();
};

export const createHabit = async (
  habitData: CreateHabitRequest,
): Promise<Habit> => {
  const response = await fetchWithAuth('/habits', {
    method: 'POST',
    body: JSON.stringify(habitData),
  });
  return response.json();
};

export const updateHabit = async (
  habitId: number,
  habitData: UpdateHabitRequest,
): Promise<Habit> => {
  const response = await fetchWithAuth(`/habits/${habitId}`, {
    method: 'PUT',
    body: JSON.stringify(habitData),
  });
  return response.json();
};

export const deleteHabit = async (habitId: number): Promise<void> => {
  await fetchWithAuth(`/habits/${habitId}`, {
    method: 'DELETE',
  });
};

// Logs API
export const getHabitLogs = async (habitId: number): Promise<Log[]> => {
  const response = await fetchWithAuth(`/habits/${habitId}/logs`);
  return response.json();
};

export const getLog = async (logId: number): Promise<Log> => {
  const response = await fetchWithAuth(`/logs/${logId}`);
  return response.json();
};

export const createLog = async (
  habitId: number,
  logData: CreateLogRequest,
): Promise<Log> => {
  const response = await fetchWithAuth(`/habits/${habitId}/logs`, {
    method: 'POST',
    body: JSON.stringify(logData),
  });
  return response.json();
};

export const updateLog = async (
  logId: number,
  logData: UpdateLogRequest,
): Promise<Log> => {
  const response = await fetchWithAuth(`/logs/${logId}`, {
    method: 'PUT',
    body: JSON.stringify(logData),
  });
  return response.json();
};

export const deleteLog = async (logId: number): Promise<void> => {
  await fetchWithAuth(`/logs/${logId}`, {
    method: 'DELETE',
  });
};
