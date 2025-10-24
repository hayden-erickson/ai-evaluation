/**
 * API client for communicating with the backend
 * Handles authentication, habits, and logs
 */

const API_BASE_URL = '';

/**
 * Get the authentication token from localStorage
 */
const getAuthToken = () => {
  return localStorage.getItem('authToken');
};

/**
 * Set the authentication token in localStorage
 */
const setAuthToken = (token) => {
  localStorage.setItem('authToken', token);
};

/**
 * Remove the authentication token from localStorage
 */
const removeAuthToken = () => {
  localStorage.removeItem('authToken');
};

/**
 * Make an authenticated API request
 */
const apiRequest = async (endpoint, options = {}) => {
  const token = getAuthToken();
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  // Add authorization header if token exists
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  // Handle different response types
  if (response.status === 204) {
    return null;
  }

  const data = await response.json().catch(() => null);

  if (!response.ok) {
    throw new Error(data?.message || response.statusText || 'An error occurred');
  }

  return data;
};

/**
 * Authentication API
 */
export const authAPI = {
  /**
   * Register a new user
   */
  register: async (userData) => {
    const response = await apiRequest('/users/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
    if (response.token) {
      setAuthToken(response.token);
    }
    return response;
  },

  /**
   * Login a user
   */
  login: async (credentials) => {
    const response = await apiRequest('/users/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
    if (response.token) {
      setAuthToken(response.token);
    }
    return response;
  },

  /**
   * Logout the current user
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

/**
 * Habits API
 */
export const habitsAPI = {
  /**
   * Get all habits for the authenticated user
   */
  getAll: async () => {
    return apiRequest('/habits', {
      method: 'GET',
    });
  },

  /**
   * Get a specific habit by ID
   */
  getById: async (habitId) => {
    return apiRequest(`/habits/${habitId}`, {
      method: 'GET',
    });
  },

  /**
   * Create a new habit
   */
  create: async (habitData) => {
    return apiRequest('/habits', {
      method: 'POST',
      body: JSON.stringify(habitData),
    });
  },

  /**
   * Update a habit
   */
  update: async (habitId, habitData) => {
    return apiRequest(`/habits/${habitId}`, {
      method: 'PUT',
      body: JSON.stringify(habitData),
    });
  },

  /**
   * Delete a habit
   */
  delete: async (habitId) => {
    return apiRequest(`/habits/${habitId}`, {
      method: 'DELETE',
    });
  },
};

/**
 * Logs API
 */
export const logsAPI = {
  /**
   * Get all logs for a specific habit
   */
  getByHabit: async (habitId) => {
    return apiRequest(`/habits/${habitId}/logs`, {
      method: 'GET',
    });
  },

  /**
   * Get a specific log by ID
   */
  getById: async (logId) => {
    return apiRequest(`/logs/${logId}`, {
      method: 'GET',
    });
  },

  /**
   * Create a new log for a habit
   */
  create: async (habitId, logData) => {
    return apiRequest(`/habits/${habitId}/logs`, {
      method: 'POST',
      body: JSON.stringify(logData),
    });
  },

  /**
   * Update a log
   */
  update: async (logId, logData) => {
    return apiRequest(`/logs/${logId}`, {
      method: 'PUT',
      body: JSON.stringify(logData),
    });
  },

  /**
   * Delete a log
   */
  delete: async (logId) => {
    return apiRequest(`/logs/${logId}`, {
      method: 'DELETE',
    });
  },
};
