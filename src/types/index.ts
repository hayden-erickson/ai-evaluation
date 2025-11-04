/**
 * Type definitions for the Habit Tracker application
 */

// User types
export interface User {
  id: number;
  name: string;
  phone_number: string;
  time_zone: string;
  profile_image_url?: string;
  created_at: string;
}

export interface LoginRequest {
  phone_number: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  phone_number: string;
  password: string;
  time_zone: string;
  profile_image_url?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

// Habit types
export interface Habit {
  id: number;
  user_id: number;
  name: string;
  description?: string;
  duration_seconds?: number;
  created_at: string;
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

// Log types
export interface Log {
  id: number;
  habit_id: number;
  notes?: string;
  duration_seconds?: number;
  created_at: string;
}

export interface CreateLogRequest {
  notes?: string;
  duration_seconds?: number;
}

export interface UpdateLogRequest {
  notes?: string;
  duration_seconds?: number;
}

// Streak calculation helper type
export interface StreakInfo {
  currentStreak: number;
  lastLogDate: Date | null;
  canContinue: boolean; // true if within the 1-day skip allowance
}
