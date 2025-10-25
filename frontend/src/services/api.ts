import {
  User,
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  Habit,
  CreateHabitRequest,
  UpdateHabitRequest,
  Log,
  CreateLogRequest,
  UpdateLogRequest,
} from '../types';

// Get API base URL from environment or default to empty string for proxied requests
const API_BASE_URL = process.env.REACT_APP_API_URL || '';

/**
 * Helper function to get auth headers
 */
const getAuthHeaders = (): HeadersInit => {
  const token = localStorage.getItem('token');
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
};

/**
 * Helper function to handle API responses
 */
const handleResponse = async <T>(response: Response): Promise<T> => {
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || `HTTP error! status: ${response.status}`);
  }
  
  // Handle empty responses (e.g., DELETE operations)
  if (response.status === 204) {
    return {} as T;
  }
  
  return response.json();
};

// Auth API
export const authAPI = {
  /**
   * Register a new user
   */
  async register(data: RegisterRequest): Promise<LoginResponse> {
    const response = await fetch(`${API_BASE_URL}/users/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return handleResponse<LoginResponse>(response);
  },

  /**
   * Login with existing credentials
   */
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await fetch(`${API_BASE_URL}/users/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return handleResponse<LoginResponse>(response);
  },
};

// User API
export const userAPI = {
  /**
   * Get user details by ID
   */
  async getUser(id: number): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/users/${id}`, {
      headers: getAuthHeaders(),
    });
    return handleResponse<User>(response);
  },

  /**
   * Delete user account
   */
  async deleteUser(id: number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/users/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    return handleResponse<void>(response);
  },
};

// Habit API
export const habitAPI = {
  /**
   * Get all habits for the authenticated user
   */
  async getUserHabits(): Promise<Habit[]> {
    const response = await fetch(`${API_BASE_URL}/habits`, {
      headers: getAuthHeaders(),
    });
    return handleResponse<Habit[]>(response);
  },

  /**
   * Get a specific habit by ID
   */
  async getHabit(id: number): Promise<Habit> {
    const response = await fetch(`${API_BASE_URL}/habits/${id}`, {
      headers: getAuthHeaders(),
    });
    return handleResponse<Habit>(response);
  },

  /**
   * Create a new habit
   */
  async createHabit(data: CreateHabitRequest): Promise<Habit> {
    const response = await fetch(`${API_BASE_URL}/habits`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    return handleResponse<Habit>(response);
  },

  /**
   * Update an existing habit
   */
  async updateHabit(id: number, data: UpdateHabitRequest): Promise<Habit> {
    const response = await fetch(`${API_BASE_URL}/habits/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    return handleResponse<Habit>(response);
  },

  /**
   * Delete a habit
   */
  async deleteHabit(id: number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/habits/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    return handleResponse<void>(response);
  },
};

// Log API
export const logAPI = {
  /**
   * Get all logs for a specific habit
   */
  async getHabitLogs(habitId: number): Promise<Log[]> {
    const response = await fetch(`${API_BASE_URL}/habits/${habitId}/logs`, {
      headers: getAuthHeaders(),
    });
    return handleResponse<Log[]>(response);
  },

  /**
   * Get a specific log by ID
   */
  async getLog(id: number): Promise<Log> {
    const response = await fetch(`${API_BASE_URL}/logs/${id}`, {
      headers: getAuthHeaders(),
    });
    return handleResponse<Log>(response);
  },

  /**
   * Create a new log for a habit
   */
  async createLog(habitId: number, data: CreateLogRequest): Promise<Log> {
    const response = await fetch(`${API_BASE_URL}/habits/${habitId}/logs`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    return handleResponse<Log>(response);
  },

  /**
   * Update an existing log
   */
  async updateLog(id: number, data: UpdateLogRequest): Promise<Log> {
    const response = await fetch(`${API_BASE_URL}/logs/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    return handleResponse<Log>(response);
  },

  /**
   * Delete a log
   */
  async deleteLog(id: number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/logs/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    return handleResponse<void>(response);
  },
};
