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
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Log {
  id: number;
  habit_id: number;
  notes: string;
  created_at: string;
  updated_at: string;
}
