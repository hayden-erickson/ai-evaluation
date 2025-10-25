import { Habit } from '../types';

const API_URL = '/habits';

const getAuthHeaders = () => {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('No token found');
  }
  return {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
  };
};

export const HabitService = {
  getHabits: async (): Promise<Habit[]> => {
    const response = await fetch(API_URL, {
      headers: getAuthHeaders(),
    });
    if (!response.ok) {
      throw new Error('Failed to fetch habits');
    }
    return response.json();
  },

  createHabit: async (habit: Omit<Habit, 'id' | 'user_id' | 'created_at'>): Promise<Habit> => {
    const response = await fetch(API_URL, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(habit),
    });
    if (!response.ok) {
      throw new Error('Failed to create habit');
    }
    return response.json();
  },

  updateHabit: async (habit: Omit<Habit, 'user_id' | 'created_at'>): Promise<Habit> => {
    const response = await fetch(`${API_URL}/${habit.id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(habit),
    });
    if (!response.ok) {
      throw new Error('Failed to update habit');
    }
    return response.json();
  },

  deleteHabit: async (id: number): Promise<void> => {
    const response = await fetch(`${API_URL}/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    if (!response.ok) {
      throw new Error('Failed to delete habit');
    }
  },
};
