/**
 * API client for communicating with the backend server
 */

import AsyncStorage from '@react-native-async-storage/async-storage';
import {
  User,
  Habit,
  Log,
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  CreateHabitRequest,
  UpdateHabitRequest,
  CreateLogRequest,
  UpdateLogRequest,
} from '../types';

// API base URL - change this to your backend server URL
const API_URL = 'http://localhost:8080';

// Storage key for auth token
const TOKEN_KEY = '@auth_token';

/**
 * Custom error class for API errors
 */
export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Get the stored authentication token
 */
export const getToken = async (): Promise<string | null> => {
  try {
    return await AsyncStorage.getItem(TOKEN_KEY);
  } catch (error) {
    console.error('Error getting token:', error);
    return null;
  }
};

/**
 * Store the authentication token
 */
export const setToken = async (token: string): Promise<void> => {
  try {
    await AsyncStorage.setItem(TOKEN_KEY, token);
  } catch (error) {
    console.error('Error setting token:', error);
  }
};

/**
 * Remove the authentication token
 */
export const removeToken = async (): Promise<void> => {
  try {
    await AsyncStorage.removeItem(TOKEN_KEY);
  } catch (error) {
    console.error('Error removing token:', error);
  }
};

/**
 * Make an authenticated API request
 */
const apiRequest = async <T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> => {
  const token = await getToken();
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({
      error: 'An error occurred',
    }));
    throw new ApiError(
      response.status,
      errorData.error || errorData.message || 'An error occurred',
    );
  }

  // Handle empty responses (like 204 No Content)
  const text = await response.text();
  return text ? JSON.parse(text) : null;
};

// Authentication API calls

/**
 * Register a new user
 */
export const register = async (
  data: RegisterRequest,
): Promise<LoginResponse> => {
  return apiRequest<LoginResponse>('/users/register', {
    method: 'POST',
    body: JSON.stringify(data),
  });
};

/**
 * Login with phone number and password
 */
export const login = async (data: LoginRequest): Promise<LoginResponse> => {
  return apiRequest<LoginResponse>('/users/login', {
    method: 'POST',
    body: JSON.stringify(data),
  });
};

// User API calls

/**
 * Get user details by ID
 */
export const getUser = async (userId: number): Promise<User> => {
  return apiRequest<User>(`/users/${userId}`);
};

// Habit API calls

/**
 * Get all habits for the current user
 */
export const getHabits = async (): Promise<Habit[]> => {
  return apiRequest<Habit[]>('/habits');
};

/**
 * Get a specific habit by ID
 */
export const getHabit = async (habitId: number): Promise<Habit> => {
  return apiRequest<Habit>(`/habits/${habitId}`);
};

/**
 * Create a new habit
 */
export const createHabit = async (data: CreateHabitRequest): Promise<Habit> => {
  return apiRequest<Habit>('/habits', {
    method: 'POST',
    body: JSON.stringify(data),
  });
};

/**
 * Update a habit
 */
export const updateHabit = async (
  habitId: number,
  data: UpdateHabitRequest,
): Promise<Habit> => {
  return apiRequest<Habit>(`/habits/${habitId}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
};

/**
 * Delete a habit
 */
export const deleteHabit = async (habitId: number): Promise<void> => {
  return apiRequest<void>(`/habits/${habitId}`, {
    method: 'DELETE',
  });
};

// Log API calls

/**
 * Get all logs for a specific habit
 */
export const getHabitLogs = async (habitId: number): Promise<Log[]> => {
  return apiRequest<Log[]>(`/habits/${habitId}/logs`);
};

/**
 * Get a specific log by ID
 */
export const getLog = async (logId: number): Promise<Log> => {
  return apiRequest<Log>(`/logs/${logId}`);
};

/**
 * Create a new log for a habit
 */
export const createLog = async (
  habitId: number,
  data: CreateLogRequest,
): Promise<Log> => {
  return apiRequest<Log>(`/habits/${habitId}/logs`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
};

/**
 * Update a log
 */
export const updateLog = async (
  logId: number,
  data: UpdateLogRequest,
): Promise<Log> => {
  return apiRequest<Log>(`/logs/${logId}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
};

/**
 * Delete a log
 */
export const deleteLog = async (logId: number): Promise<void> => {
  return apiRequest<void>(`/logs/${logId}`, {
    method: 'DELETE',
  });
};
