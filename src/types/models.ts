/**
 * Type definitions for API models
 */

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

export interface CreateUserRequest {
  profile_image_url?: string;
  name: string;
  time_zone: string;
  phone_number: string;
  password: string;
}

export interface LoginRequest {
  phone_number: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
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

/**
 * Extended types for UI state
 */
export interface HabitWithLogs extends Habit {
  logs: Log[];
  currentStreak: number;
  lastLogDate: Date | null;
}

