/**
 * Type definitions for the Habit Tracker app
 */

export interface User {
  id: number;
  name: string;
  phone_number: string;
  time_zone: string;
  profile_image_url?: string;
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

export interface LoginResponse {
  token: string;
  user: User;
}

export interface CreateUserRequest {
  name: string;
  phone_number: string;
  password: string;
  time_zone: string;
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
