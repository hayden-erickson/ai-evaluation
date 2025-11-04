/**
 * API Service for communicating with the backend
 * Handles all HTTP requests with proper error handling
 */

import {
  AuthResponse,
  CreateHabitRequest,
  CreateLogRequest,
  Habit,
  Log,
  LoginRequest,
  RegisterRequest,
  UpdateHabitRequest,
  UpdateLogRequest,
  User,
} from '../types';

// Configure your API base URL here
const API_BASE_URL = 'http://localhost:8080';

/**
 * Custom error class for API errors
 */
export class ApiError extends Error {
  constructor(
    message: string,
    public statusCode?: number,
    public details?: any,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Makes an HTTP request with proper error handling
 */
async function request<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    // Handle non-JSON responses
    const contentType = response.headers.get('content-type');
    const isJson = contentType?.includes('application/json');

    if (!response.ok) {
      const errorBody = isJson ? await response.json() : await response.text();
      const errorMessage =
        typeof errorBody === 'string'
          ? errorBody
          : errorBody.error || errorBody.message || 'An error occurred';

      throw new ApiError(errorMessage, response.status, errorBody);
    }

    // Handle empty responses (e.g., DELETE requests)
    if (response.status === 204 || !isJson) {
      return {} as T;
    }

    return await response.json();
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }

    // Network or other errors
    throw new ApiError(
      error instanceof Error ? error.message : 'Network error occurred',
    );
  }
}

/**
 * Authentication API methods
 */
export const authApi = {
  /**
   * Register a new user
   */
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    return request<AuthResponse>('/users/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  /**
   * Login an existing user
   */
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    return request<AuthResponse>('/users/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },
};

/**
 * User API methods
 */
export const userApi = {
  /**
   * Get user details by ID
   */
  getUser: async (userId: number, token: string): Promise<User> => {
    return request<User>(`/users/${userId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },

  /**
   * Delete user account
   */
  deleteUser: async (userId: number, token: string): Promise<void> => {
    return request<void>(`/users/${userId}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },
};

/**
 * Habit API methods
 */
export const habitApi = {
  /**
   * Get all habits for the authenticated user
   */
  getHabits: async (token: string): Promise<Habit[]> => {
    return request<Habit[]>('/habits', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },

  /**
   * Get a specific habit by ID
   */
  getHabit: async (habitId: number, token: string): Promise<Habit> => {
    return request<Habit>(`/habits/${habitId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },

  /**
   * Create a new habit
   */
  createHabit: async (
    data: CreateHabitRequest,
    token: string,
  ): Promise<Habit> => {
    return request<Habit>('/habits', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
  },

  /**
   * Update an existing habit
   */
  updateHabit: async (
    habitId: number,
    data: UpdateHabitRequest,
    token: string,
  ): Promise<Habit> => {
    return request<Habit>(`/habits/${habitId}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
  },

  /**
   * Delete a habit
   */
  deleteHabit: async (habitId: number, token: string): Promise<void> => {
    return request<void>(`/habits/${habitId}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },
};

/**
 * Log API methods
 */
export const logApi = {
  /**
   * Get all logs for a specific habit
   */
  getHabitLogs: async (habitId: number, token: string): Promise<Log[]> => {
    return request<Log[]>(`/habits/${habitId}/logs`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },

  /**
   * Get a specific log by ID
   */
  getLog: async (logId: number, token: string): Promise<Log> => {
    return request<Log>(`/logs/${logId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },

  /**
   * Create a new log for a habit
   */
  createLog: async (
    habitId: number,
    data: CreateLogRequest,
    token: string,
  ): Promise<Log> => {
    return request<Log>(`/habits/${habitId}/logs`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
  },

  /**
   * Update an existing log
   */
  updateLog: async (
    logId: number,
    data: UpdateLogRequest,
    token: string,
  ): Promise<Log> => {
    return request<Log>(`/logs/${logId}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
  },

  /**
   * Delete a log
   */
  deleteLog: async (logId: number, token: string): Promise<void> => {
    return request<void>(`/logs/${logId}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },
};
