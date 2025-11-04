/**
 * API Service for communicating with the backend
 * Handles all HTTP requests and responses
 */

import {
  User,
  Habit,
  Log,
  LoginRequest,
  LoginResponse,
  CreateUserRequest,
  CreateHabitRequest,
  UpdateHabitRequest,
  CreateLogRequest,
  UpdateLogRequest,
} from '../types';

// Default API base URL - can be configured via environment
const API_BASE_URL = 'http://localhost:8080';

class ApiService {
  private baseUrl: string;
  private token: string | null = null;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  /**
   * Set the authentication token for subsequent requests
   */
  setToken(token: string | null) {
    this.token = token;
  }

  /**
   * Get the current authentication token
   */
  getToken(): string | null {
    return this.token;
  }

  /**
   * Helper method to make authenticated requests
   */
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    // Add authorization header if token is available
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers,
    });

    // Handle error responses
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || `HTTP ${response.status}: ${response.statusText}`);
    }

    // Handle empty responses
    const text = await response.text();
    if (!text) {
      return {} as T;
    }

    return JSON.parse(text);
  }

  // ==================== Authentication ====================

  /**
   * Register a new user
   */
  async register(data: CreateUserRequest): Promise<User> {
    return this.request<User>('/users/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  /**
   * Login with phone number and password
   */
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await this.request<LoginResponse>('/users/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    // Store the token for subsequent requests
    this.setToken(response.token);
    return response;
  }

  /**
   * Logout (clear token)
   */
  logout() {
    this.setToken(null);
  }

  // ==================== User Operations ====================

  /**
   * Get user details by ID
   */
  async getUser(userId: number): Promise<User> {
    return this.request<User>(`/users/${userId}`);
  }

  /**
   * Update user details
   */
  async updateUser(userId: number, data: Partial<User>): Promise<User> {
    return this.request<User>(`/users/${userId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  /**
   * Delete user account
   */
  async deleteUser(userId: number): Promise<void> {
    return this.request<void>(`/users/${userId}`, {
      method: 'DELETE',
    });
  }

  // ==================== Habit Operations ====================

  /**
   * Get all habits for the authenticated user
   */
  async getHabits(): Promise<Habit[]> {
    return this.request<Habit[]>('/habits');
  }

  /**
   * Get a specific habit by ID
   */
  async getHabit(habitId: number): Promise<Habit> {
    return this.request<Habit>(`/habits/${habitId}`);
  }

  /**
   * Create a new habit
   */
  async createHabit(data: CreateHabitRequest): Promise<Habit> {
    return this.request<Habit>('/habits', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  /**
   * Update an existing habit
   */
  async updateHabit(habitId: number, data: UpdateHabitRequest): Promise<Habit> {
    return this.request<Habit>(`/habits/${habitId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  /**
   * Delete a habit
   */
  async deleteHabit(habitId: number): Promise<void> {
    return this.request<void>(`/habits/${habitId}`, {
      method: 'DELETE',
    });
  }

  // ==================== Log Operations ====================

  /**
   * Get all logs for a specific habit
   */
  async getHabitLogs(habitId: number): Promise<Log[]> {
    return this.request<Log[]>(`/habits/${habitId}/logs`);
  }

  /**
   * Get a specific log by ID
   */
  async getLog(logId: number): Promise<Log> {
    return this.request<Log>(`/logs/${logId}`);
  }

  /**
   * Create a new log for a habit
   */
  async createLog(habitId: number, data: CreateLogRequest): Promise<Log> {
    return this.request<Log>(`/habits/${habitId}/logs`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  /**
   * Update an existing log
   */
  async updateLog(logId: number, data: UpdateLogRequest): Promise<Log> {
    return this.request<Log>(`/logs/${logId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  /**
   * Delete a log
   */
  async deleteLog(logId: number): Promise<void> {
    return this.request<void>(`/logs/${logId}`, {
      method: 'DELETE',
    });
  }

  // ==================== Health Check ====================

  /**
   * Check API health
   */
  async healthCheck(): Promise<string> {
    const response = await fetch(`${this.baseUrl}/health`);
    return response.text();
  }
}

// Export a singleton instance
export const apiService = new ApiService();
export default ApiService;
