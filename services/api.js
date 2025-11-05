// services/api.js
const API_URL = 'http://10.0.2.2:8080';

const getAuthToken = () => {
  // In a real app, you'd get this from AsyncStorage or a similar secure storage
  return 'your-jwt-token';
};

const apiRequest = async (endpoint, method, body) => {
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${getAuthToken()}`,
  };

  const config = {
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
    return response.json();
  } catch (error) {
    console.error('API request error:', error);
    throw error;
  }
};

// Habit endpoints
export const getHabits = () => apiRequest('/habits', 'GET');
export const createHabit = (habit) => apiRequest('/habits', 'POST', habit);
export const updateHabit = (id, habit) => apiRequest(`/habits/${id}`, 'PUT', habit);
export const deleteHabit = (id) => apiRequest(`/habits/${id}`, 'DELETE');

// Log endpoints
export const getLogs = (habitId) => apiRequest(`/habits/${habitId}/logs`, 'GET');
export const createLog = (habitId, log) => apiRequest(`/habits/${habitId}/logs`, 'POST', log);
export const updateLog = (id, log) => apiRequest(`/logs/${id}`, 'PUT', log);
export const deleteLog = (id) => apiRequest(`/logs/${id}`, 'DELETE');
