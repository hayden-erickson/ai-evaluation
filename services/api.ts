const API_URL = 'http://10.0.2.2:8080';

let authToken: string | null = null;

interface Habit {
  id: string;
  name: string;
  description: string;
}

interface Log {
  id: string;
  habit_id: string;
  date: string;
  notes: string;
}

export const setAuthToken = (token: string | null) => {
  authToken = token;
};

const getAuthToken = (): string | null => {
  return authToken;
};

const apiRequest = async (endpoint: string, method: string, body?: any, requiresAuth = true): Promise<any> => {
  const headers: { [key: string]: string } = {
    'Content-Type': 'application/json',
  };

  if (requiresAuth) {
    const token = getAuthToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
  }

  const config: RequestInit = {
    method,
    headers,
    body: body ? JSON.stringify(body) : null,
  };

  try {
    const response = await fetch(`${API_URL}${endpoint}`, config);
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Something went wrong');
    }
    if (method === 'DELETE' || response.status === 204) {
      return;
    }
    return response.json();
  } catch (error) {
    console.error('API request error:', error);
    throw error;
  }
};

// Auth endpoints
export const login = (credentials: any): Promise<{ token: string }> => apiRequest('/users/login', 'POST', credentials, false);
export const register = (userInfo: any): Promise<{ token: string }> => apiRequest('/users/register', 'POST', userInfo, false);

// Habit endpoints
export const getHabits = (): Promise<Habit[]> => apiRequest('/habits', 'GET');
export const createHabit = (habit: Omit<Habit, 'id'>): Promise<Habit> => apiRequest('/habits', 'POST', habit);
export const updateHabit = (id: string, habit: Partial<Habit>): Promise<Habit> => apiRequest(`/habits/${id}`, 'PUT', habit);
export const deleteHabit = (id: string): Promise<void> => apiRequest(`/habits/${id}`, 'DELETE');

// Log endpoints
export const getLogs = (habitId: string): Promise<Log[]> => apiRequest(`/habits/${habitId}/logs`, 'GET');
export const createLog = (habitId: string, log: Omit<Log, 'id' | 'habit_id'>): Promise<Log> => apiRequest(`/habits/${habitId}/logs`, 'POST', log);
export const updateLog = (id: string, log: Partial<Log>): Promise<Log> => apiRequest(`/logs/${id}`, 'PUT', log);
export const deleteLog = (id: string): Promise<void> => apiRequest(`/logs/${id}`, 'DELETE');
