/**
 * API service layer for communicating with the backend
 * Handles authentication, habits, and logs
 */

const API_BASE_URL = '/api';

/**
 * Get the stored authentication token from localStorage
 * @returns {string|null} The JWT token or null if not found
 */
const getAuthToken = () => {
  return localStorage.getItem('auth_token');
};

/**
 * Store authentication token in localStorage
 * @param {string} token - The JWT token to store
 */
const setAuthToken = (token) => {
  localStorage.setItem('auth_token', token);
};

/**
 * Remove authentication token from localStorage
 */
const clearAuthToken = () => {
  localStorage.removeItem('auth_token');
};

/**
 * Make an authenticated API request
 * @param {string} endpoint - The API endpoint
 * @param {object} options - Fetch options
 * @returns {Promise<object>} The response data
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

  const config = {
    ...options,
    headers,
  };

  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, config);

    // Handle different response types
    if (response.status === 204) {
      return null; // No content
    }

    const data = await response.json().catch(() => null);

    if (!response.ok) {
      // Handle authentication errors
      if (response.status === 401) {
        clearAuthToken();
        throw new Error('Session expired. Please login again.');
      }

      // Throw error with message from server or default message
      throw new Error(data?.message || data || `Request failed: ${response.statusText}`);
    }

    return data;
  } catch (error) {
    // Re-throw errors to be handled by the caller
    if (error.message) {
      throw error;
    }
    throw new Error('Network error. Please check your connection.');
  }
};

// ============= Authentication API =============

/**
 * Register a new user
 * @param {object} userData - User registration data
 * @returns {Promise<object>} The created user
 */
export const register = async (userData) => {
  return apiRequest('/users/register', {
    method: 'POST',
    body: JSON.stringify(userData),
  });
};

/**
 * Login user
 * @param {string} phoneNumber - User's phone number
 * @param {string} password - User's password
 * @returns {Promise<object>} Login response with token and user data
 */
export const login = async (phoneNumber, password) => {
  const response = await apiRequest('/users/login', {
    method: 'POST',
    body: JSON.stringify({ phone_number: phoneNumber, password }),
  });
  
  // Store the token
  if (response.token) {
    setAuthToken(response.token);
  }
  
  return response;
};

/**
 * Logout user
 */
export const logout = () => {
  clearAuthToken();
};

/**
 * Check if user is authenticated
 * @returns {boolean} True if user has a valid token
 */
export const isAuthenticated = () => {
  return !!getAuthToken();
};

// ============= Habits API =============

/**
 * Get all habits for the authenticated user
 * @returns {Promise<Array>} Array of habits
 */
export const getHabits = async () => {
  return apiRequest('/habits', {
    method: 'GET',
  });
};

/**
 * Get a single habit by ID
 * @param {number} habitId - The habit ID
 * @returns {Promise<object>} The habit data
 */
export const getHabit = async (habitId) => {
  return apiRequest(`/habits/${habitId}`, {
    method: 'GET',
  });
};

/**
 * Create a new habit
 * @param {object} habitData - Habit data (name, description, duration_seconds)
 * @returns {Promise<object>} The created habit
 */
export const createHabit = async (habitData) => {
  return apiRequest('/habits', {
    method: 'POST',
    body: JSON.stringify(habitData),
  });
};

/**
 * Update an existing habit
 * @param {number} habitId - The habit ID
 * @param {object} habitData - Updated habit data
 * @returns {Promise<object>} The updated habit
 */
export const updateHabit = async (habitId, habitData) => {
  return apiRequest(`/habits/${habitId}`, {
    method: 'PUT',
    body: JSON.stringify(habitData),
  });
};

/**
 * Delete a habit
 * @param {number} habitId - The habit ID
 * @returns {Promise<null>}
 */
export const deleteHabit = async (habitId) => {
  return apiRequest(`/habits/${habitId}`, {
    method: 'DELETE',
  });
};

// ============= Logs API =============

/**
 * Get all logs for a habit
 * @param {number} habitId - The habit ID
 * @returns {Promise<Array>} Array of logs
 */
export const getHabitLogs = async (habitId) => {
  return apiRequest(`/habits/${habitId}/logs`, {
    method: 'GET',
  });
};

/**
 * Get a single log by ID
 * @param {number} logId - The log ID
 * @returns {Promise<object>} The log data
 */
export const getLog = async (logId) => {
  return apiRequest(`/logs/${logId}`, {
    method: 'GET',
  });
};

/**
 * Create a new log for a habit
 * @param {number} habitId - The habit ID
 * @param {object} logData - Log data (notes, duration_seconds)
 * @returns {Promise<object>} The created log
 */
export const createLog = async (habitId, logData) => {
  return apiRequest(`/habits/${habitId}/logs`, {
    method: 'POST',
    body: JSON.stringify(logData),
  });
};

/**
 * Update an existing log
 * @param {number} logId - The log ID
 * @param {object} logData - Updated log data
 * @returns {Promise<object>} The updated log
 */
export const updateLog = async (logId, logData) => {
  return apiRequest(`/logs/${logId}`, {
    method: 'PUT',
    body: JSON.stringify(logData),
  });
};

/**
 * Delete a log
 * @param {number} logId - The log ID
 * @returns {Promise<null>}
 */
export const deleteLog = async (logId) => {
  return apiRequest(`/logs/${logId}`, {
    method: 'DELETE',
  });
};

