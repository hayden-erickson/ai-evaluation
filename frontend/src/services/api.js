/**
 * API Service Layer
 * Handles all HTTP requests to the backend API with authentication
 */

const API_BASE_URL = '';

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
 * Make an authenticated API request
 */
const apiRequest = async (url, options = {}) => {
  const token = getAuthToken();
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  // Add authorization header if token exists
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  try {
    const response = await fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers,
    });

    // Handle unauthorized responses
    if (response.status === 401) {
      removeAuthToken();
      window.location.href = '/login';
      throw new Error('Unauthorized');
    }

    // Parse response
    const contentType = response.headers.get('content-type');
    let data;
    if (contentType && contentType.includes('application/json')) {
      data = await response.json();
    } else {
      data = await response.text();
    }

    // Handle error responses
    if (!response.ok) {
      throw new Error(data.error || data || 'Request failed');
    }

    return data;
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
};

// Auth API
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
   * Login user
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
   * Get all habits for the current user
   */
  getAll: async () => {
    return await apiRequest('/habits', {
      method: 'GET',
    });
  },

  /**
   * Get a specific habit by ID
   */
  getById: async (habitId) => {
    return await apiRequest(`/habits/${habitId}`, {
      method: 'GET',
    });
  },

  /**
   * Create a new habit
   */
  create: async (habitData) => {
    return await apiRequest('/habits', {
      method: 'POST',
      body: JSON.stringify(habitData),
    });
  },

  /**
   * Update a habit
   */
  update: async (habitId, habitData) => {
    return await apiRequest(`/habits/${habitId}`, {
      method: 'PUT',
      body: JSON.stringify(habitData),
    });
  },

  /**
   * Delete a habit
   */
  delete: async (habitId) => {
    return await apiRequest(`/habits/${habitId}`, {
      method: 'DELETE',
    });
  },

  /**
   * Get logs for a specific habit
   */
  getLogs: async (habitId) => {
    return await apiRequest(`/habits/${habitId}/logs`, {
      method: 'GET',
    });
  },
};

// Logs API
export const logsAPI = {
  /**
   * Create a new log for a habit
   */
  create: async (habitId, logData) => {
    return await apiRequest(`/habits/${habitId}/logs`, {
      method: 'POST',
      body: JSON.stringify(logData),
    });
  },

  /**
   * Get a specific log by ID
   */
  getById: async (logId) => {
    return await apiRequest(`/logs/${logId}`, {
      method: 'GET',
    });
  },

  /**
   * Update a log
   */
  update: async (logId, logData) => {
    return await apiRequest(`/logs/${logId}`, {
      method: 'PUT',
      body: JSON.stringify(logData),
    });
  },

  /**
   * Delete a log
   */
  delete: async (logId) => {
    return await apiRequest(`/logs/${logId}`, {
      method: 'DELETE',
    });
  },
};
