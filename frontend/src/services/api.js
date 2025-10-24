// API base URL - defaults to localhost for development
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

/**
 * Get authentication token from localStorage
 */
const getAuthToken = () => {
  return localStorage.getItem('authToken');
};

/**
 * Set authentication token in localStorage
 */
const setAuthToken = (token) => {
  localStorage.setItem('authToken', token);
};

/**
 * Remove authentication token from localStorage
 */
const removeAuthToken = () => {
  localStorage.removeItem('authToken');
};

/**
 * Base fetch wrapper with error handling and authentication
 */
const apiFetch = async (endpoint, options = {}) => {
  const token = getAuthToken();
  
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
    
    // Handle 401 Unauthorized
    if (response.status === 401) {
      removeAuthToken();
      throw new Error('Session expired. Please login again.');
    }

    // Handle non-2xx responses
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || `HTTP error! status: ${response.status}`);
    }

    // Handle 204 No Content
    if (response.status === 204) {
      return null;
    }

    return await response.json();
  } catch (error) {
    // Network error or parsing error
    if (error.name === 'TypeError' && error.message === 'Failed to fetch') {
      throw new Error('Cannot connect to server. Please ensure the backend is running.');
    }
    throw error;
  }
};

// Authentication API
export const authAPI = {
  /**
   * Register a new user
   */
  register: async (userData) => {
    const response = await apiFetch('/users/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
    if (response.token) {
      setAuthToken(response.token);
    }
    return response;
  },

  /**
   * Login user
   */
  login: async (credentials) => {
    const response = await apiFetch('/users/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
    if (response.token) {
      setAuthToken(response.token);
    }
    return response;
  },

  /**
   * Logout user
   */
  logout: () => {
    removeAuthToken();
  },

  /**
   * Check if user is authenticated
   */
  isAuthenticated: () => {
    return !!getAuthToken();
  },
};

// Habits API
export const habitsAPI = {
  /**
   * Get all habits for the authenticated user
   */
  getAll: async () => {
    return await apiFetch('/habits', {
      method: 'GET',
    });
  },

  /**
   * Get a specific habit by ID
   */
  getById: async (habitId) => {
    return await apiFetch(`/habits/${habitId}`, {
      method: 'GET',
    });
  },

  /**
   * Create a new habit
   */
  create: async (habitData) => {
    return await apiFetch('/habits', {
      method: 'POST',
      body: JSON.stringify(habitData),
    });
  },

  /**
   * Update an existing habit
   */
  update: async (habitId, habitData) => {
    return await apiFetch(`/habits/${habitId}`, {
      method: 'PUT',
      body: JSON.stringify(habitData),
    });
  },

  /**
   * Delete a habit
   */
  delete: async (habitId) => {
    return await apiFetch(`/habits/${habitId}`, {
      method: 'DELETE',
    });
  },
};

// Logs API
export const logsAPI = {
  /**
   * Get all logs for a habit
   */
  getByHabit: async (habitId) => {
    return await apiFetch(`/habits/${habitId}/logs`, {
      method: 'GET',
    });
  },

  /**
   * Get a specific log by ID
   */
  getById: async (logId) => {
    return await apiFetch(`/logs/${logId}`, {
      method: 'GET',
    });
  },

  /**
   * Create a new log for a habit
   */
  create: async (habitId, logData) => {
    return await apiFetch(`/habits/${habitId}/logs`, {
      method: 'POST',
      body: JSON.stringify(logData),
    });
  },

  /**
   * Update an existing log
   */
  update: async (logId, logData) => {
    return await apiFetch(`/logs/${logId}`, {
      method: 'PUT',
      body: JSON.stringify(logData),
    });
  },

  /**
   * Delete a log
   */
  delete: async (logId) => {
    return await apiFetch(`/logs/${logId}`, {
      method: 'DELETE',
    });
  },
};
