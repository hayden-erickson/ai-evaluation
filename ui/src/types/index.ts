export interface User {
  id: number;
  username: string;
  email: string;
}

export interface Habit {
  id: number;
  user_id: number;
  name: string;
  description: string;
  created_at: string;
}

export interface Log {
  id: number;
  habit_id: number;
  date: string;
  notes: string;
  duration: number; // in minutes
}
