// Type definitions shared across the app

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
  duration_seconds?: number | null;
  created_at: string;
}

export interface LogEntry {
  id: number;
  habit_id: number;
  notes?: string;
  duration_seconds?: number | null;
  created_at: string; // ISO
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface CreateHabitRequest {
  name: string;
  description?: string;
  duration_seconds?: number | null;
}

export interface UpdateHabitRequest {
  name?: string;
  description?: string;
  duration_seconds?: number | null;
}

export interface CreateLogRequest {
  notes?: string;
  duration_seconds?: number | null;
}

export interface UpdateLogRequest {
  notes?: string;
  duration_seconds?: number | null;
}
