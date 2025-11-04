import api from './api';
import { Habit } from '../types/types';

export const getHabits = async (): Promise<Habit[]> => {
  try {
    const response = await api.get('/habits');
    return response.data;
  } catch (error) {
    console.error('Failed to fetch habits', error);
    throw error;
  }
};

export const createHabit = async (habitData: { name: string; description: string }): Promise<Habit> => {
  try {
    const response = await api.post('/habits', habitData);
    return response.data;
  } catch (error) {
    console.error('Failed to create habit', error);
    throw error;
  }
};

export const updateHabit = async (id: number, habitData: { name?: string; description?: string }): Promise<Habit> => {
  try {
    const response = await api.put(`/habits/${id}`, habitData);
    return response.data;
  } catch (error) {
    console.error(`Failed to update habit ${id}`, error);
    throw error;
  }
};

export const deleteHabit = async (id: number): Promise<void> => {
  try {
    await api.delete(`/habits/${id}`);
  } catch (error) {
    console.error(`Failed to delete habit ${id}`, error);
    throw error;
  }
};
