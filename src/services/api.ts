/**
 * API service for communicating with the backend
 */

import AsyncStorage from '@react-native-async-storage/async-storage';
import {
  User,
  Habit,
  Log,
  CreateUserRequest,
  LoginRequest,
  LoginResponse,
  CreateHabitRequest,
  UpdateHabitRequest,
  CreateLogRequest,
  UpdateLogRequest,
} from '../types/models';

const API_BASE_URL = 'http://localhost:8080'; // Change this for production
const TOKEN_KEY = '@habit_tracker_token';

/**
 * HTTP client wrapper with authentication
 */
class APIClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  /**
   * Get stored authentication token
   */
  async getToken(): Promise<string | null> {
    try {
      return await AsyncStorage.getItem(TOKEN_KEY);
    } catch (error) {
      console.error('Error getting token:', error);
      return null;
    }
  }

  /**
   * Store authentication token
   */
  async setToken(token: string): Promise<void> {
    try {
      await AsyncStorage.setItem(TOKEN_KEY, token);
    } catch (error) {
      console.error('Error storing token:', error);
      throw new Error('Failed to store authentication token');
    }
  }

  /**
   * Clear authentication token
   */
  async clearToken(): Promise<void> {
    try {
      await AsyncStorage.removeItem(TOKEN_KEY);
    } catch (error) {
      console.error('Error clearing token:', error);
    }
  }

  /**
   * Make an authenticated HTTP request
   */
  async request<T>(
    endpoint: string,
    options: RequestInit = {},
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const token = await this.getToken();

    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    // Add authentication header if token exists
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      // Handle empty responses (e.g., 204 No Content)
      if (response.status === 204) {
        return {} as T;
      }

      // Parse response body
      const data = await response.json().catch(() => null);

      // Handle error responses
      if (!response.ok) {
        const errorMessage = data?.error || data || response.statusText || 'Request failed';
        throw new Error(errorMessage);
      }

      return data as T;
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error('Network request failed');
    }
  }

  /**
   * Make a GET request
   */
  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  /**
   * Make a POST request
   */
  async post<T>(endpoint: string, data: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  /**
   * Make a PUT request
   */
  async put<T>(endpoint: string, data: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  /**
   * Make a DELETE request
   */
  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

// Create singleton instance
const apiClient = new APIClient(API_BASE_URL);

/**
 * User API methods
 */
export const userAPI = {
  /**
   * Register a new user
   */
  register: async (data: CreateUserRequest): Promise<User> => {
    return apiClient.post<User>('/users/register', data);
  },

  /**
   * Login user and get authentication token
   */
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/users/login', data);
    // Store the token
    await apiClient.setToken(response.token);
    return response;
  },

  /**
   * Logout user (clear token)
   */
  logout: async (): Promise<void> => {
    await apiClient.clearToken();
  },

  /**
   * Get user by ID
   */
  getUser: async (userId: number): Promise<User> => {
    return apiClient.get<User>(`/users/${userId}`);
  },
};

/**
 * Habit API methods
 */
export const habitAPI = {
  /**
   * Get all habits for authenticated user
   */
  getUserHabits: async (): Promise<Habit[]> => {
    return apiClient.get<Habit[]>('/habits');
  },

  /**
   * Get a single habit by ID
   */
  getHabit: async (habitId: number): Promise<Habit> => {
    return apiClient.get<Habit>(`/habits/${habitId}`);
  },

  /**
   * Create a new habit
   */
  createHabit: async (data: CreateHabitRequest): Promise<Habit> => {
    return apiClient.post<Habit>('/habits', data);
  },

  /**
   * Update an existing habit
   */
  updateHabit: async (habitId: number, data: UpdateHabitRequest): Promise<Habit> => {
    return apiClient.put<Habit>(`/habits/${habitId}`, data);
  },

  /**
   * Delete a habit
   */
  deleteHabit: async (habitId: number): Promise<void> => {
    return apiClient.delete<void>(`/habits/${habitId}`);
  },
};

/**
 * Log API methods
 */
export const logAPI = {
  /**
   * Get all logs for a habit
   */
  getHabitLogs: async (habitId: number): Promise<Log[]> => {
    return apiClient.get<Log[]>(`/habits/${habitId}/logs`);
  },

  /**
   * Get a single log by ID
   */
  getLog: async (logId: number): Promise<Log> => {
    return apiClient.get<Log>(`/logs/${logId}`);
  },

  /**
   * Create a new log for a habit
   */
  createLog: async (habitId: number, data: CreateLogRequest): Promise<Log> => {
    return apiClient.post<Log>(`/habits/${habitId}/logs`, data);
  },

  /**
   * Update an existing log
   */
  updateLog: async (logId: number, data: UpdateLogRequest): Promise<Log> => {
    return apiClient.put<Log>(`/logs/${logId}`, data);
  },

  /**
   * Delete a log
   */
  deleteLog: async (logId: number): Promise<void> => {
    return apiClient.delete<void>(`/logs/${logId}`);
  },
};

export { apiClient };

