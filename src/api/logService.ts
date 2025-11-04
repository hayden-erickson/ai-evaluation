import api from './api';
import { Log } from '../types/types';

export const getHabitLogs = async (habitId: number): Promise<Log[]> => {
  try {
    const response = await api.get(`/habits/${habitId}/logs`);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch logs for habit ${habitId}`, error);
    throw error;
  }
};

export const createLog = async (habitId: number, logData: { notes: string }): Promise<Log> => {
  try {
    const response = await api.post(`/habits/${habitId}/logs`, logData);
    return response.data;
  } catch (error) {
    console.error('Failed to create log', error);
    throw error;
  }
};

export const updateLog = async (id: number, logData: { notes?: string }): Promise<Log> => {
  try {
    const response = await api.put(`/logs/${id}`, logData);
    return response.data;
  } catch (error) {
    console.error(`Failed to update log ${id}`, error);
    throw error;
  }
};

export const deleteLog = async (id: number): Promise<void> => {
  try {
    await api.delete(`/logs/${id}`);
  } catch (error) {
    console.error(`Failed to delete log ${id}`, error);
    throw error;
  }
};
