/**
 * Configuration file for environment-specific settings
 */

/**
 * Get the API base URL based on the environment
 * 
 * Configuration:
 * - For iOS Simulator: Use localhost:8080
 * - For Android Emulator: Use 10.0.2.2:8080
 * - For Physical Devices: Use your computer's IP address
 * 
 * To change the API URL, modify the API_BASE_URL below
 */

// Default configuration for different platforms
// Change this to match your environment
export const API_BASE_URL = __DEV__ 
  ? 'http://localhost:8080'  // Development mode (iOS Simulator)
  : 'https://your-production-api.com';  // Production mode

/**
 * Alternative configurations (uncomment the one you need):
 */
// For Android Emulator:
// export const API_BASE_URL = 'http://10.0.2.2:8080';

// For Physical Device (replace with your computer's IP):
// export const API_BASE_URL = 'http://192.168.1.100:8080';

/**
 * JWT token expiration time (in seconds)
 * Should match backend configuration
 */
export const JWT_EXPIRATION = 24 * 60 * 60; // 24 hours

/**
 * Maximum number of days to display in streak list
 */
export const STREAK_DISPLAY_DAYS = 7;

/**
 * Maximum streak skip days allowed before reset
 */
export const MAX_SKIP_DAYS = 1;
