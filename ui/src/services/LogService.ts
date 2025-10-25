import { Log, Habit } from '../types';

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

export const LogService = {
    getLogs: async (habitId: number): Promise<Log[]> => {
        const response = await fetch(`${API_URL}/${habitId}/logs`, {
            headers: getAuthHeaders(),
        });
        if (!response.ok) {
            throw new Error('Failed to fetch logs');
        }
        return response.json();
    },

    createLog: async (habitId: number, log: Omit<Log, 'id' | 'habit_id'>): Promise<Log> => {
        const response = await fetch(`${API_URL}/${habitId}/logs`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify(log),
        });
        if (!response.ok) {
            throw new Error('Failed to create log');
        }
        return response.json();
    },

    updateLog: async (log: Omit<Log, 'habit_id'>): Promise<Log> => {
        const response = await fetch(`/logs/${log.id}`, {
            method: 'PUT',
            headers: getAuthHeaders(),
            body: JSON.stringify(log),
        });
        if (!response.ok) {
            throw new Error('Failed to update log');
        }
        return response.json();
    },

    deleteLog: async (id: number): Promise<void> => {
        const response = await fetch(`/logs/${id}`, {
            method: 'DELETE',
            headers: getAuthHeaders(),
        });
        if (!response.ok) {
            throw new Error('Failed to delete log');
        }
    },
};
