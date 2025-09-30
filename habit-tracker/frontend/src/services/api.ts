import axios from 'axios';
import type { AuthResponse, Habit, HabitWithStreak, Log, Tag, LogFormData } from '../types';

const API_URL = 'http://localhost:8080';

const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add token to requests if available
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Auth API
export const loginWithGoogle = async (accessToken: string): Promise<AuthResponse> => {
    const response = await api.post('/auth/google/callback', { accessToken });
    return response.data;
};

// Habits API
export const getHabits = async (): Promise<Habit[]> => {
    const response = await api.get('/api/habits');
    return response.data;
};

export const getHabitsByTag = async (tag: string): Promise<Habit[]> => {
    const response = await api.get(`/api/habits/by-tag?value=${encodeURIComponent(tag)}`);
    return response.data;
};

export const getHabit = async (id: number): Promise<Habit> => {
    const response = await api.get(`/api/habits/${id}`);
    return response.data;
};

export const getHabitWithStreak = async (id: number): Promise<HabitWithStreak> => {
    const response = await api.get(`/api/habits/${id}/streak`);
    return response.data;
};

export const createHabit = async (habit: Partial<Habit>): Promise<Habit> => {
    const response = await api.post('/api/habits', habit);
    return response.data;
};

export const updateHabit = async (id: number, habit: Partial<Habit>): Promise<Habit> => {
    const response = await api.put(`/api/habits/${id}`, habit);
    return response.data;
};

export const deleteHabit = async (id: number): Promise<void> => {
    await api.delete(`/api/habits/${id}`);
};

// Logs API
export const getLogs = async (habitId: number): Promise<Log[]> => {
    const response = await api.get(`/api/habits/${habitId}/logs`);
    return response.data;
};

export const createLog = async (habitId: number, data: LogFormData): Promise<Log> => {
    const response = await api.post(`/api/habits/${habitId}/logs`, data);
    return response.data;
};

export const updateLog = async (habitId: number, logId: number, data: LogFormData): Promise<Log> => {
    const response = await api.put(`/api/habits/${habitId}/logs/${logId}`, data);
    return response.data;
};

export const deleteLog = async (habitId: number, logId: number): Promise<void> => {
    await api.delete(`/api/habits/${habitId}/logs/${logId}`);
};

// Tags API
export const getTags = async (habitId: number): Promise<Tag[]> => {
    const response = await api.get(`/api/habits/${habitId}/tags`);
    return response.data;
};

export const createTag = async (habitId: number, value: string): Promise<Tag> => {
    const response = await api.post(`/api/habits/${habitId}/tags`, { value });
    return response.data;
};

export const deleteTag = async (habitId: number, tagId: number): Promise<void> => {
    await api.delete(`/api/habits/${habitId}/tags/${tagId}`);
};