import {Platform} from 'react-native';

/**
 * Configuration file for the app
 * Handles environment-specific settings
 */

/**
 * Get the appropriate API base URL based on the platform
 * - iOS Simulator: localhost works
 * - Android Emulator: Use 10.0.2.2 to access host machine
 * - Physical devices: Use your computer's local IP address
 */
export const getApiBaseUrl = (): string => {
  // For production, use your production API URL
  if (__DEV__) {
    // Development mode
    if (Platform.OS === 'android') {
      // Android emulator uses 10.0.2.2 to access host machine's localhost
      return 'http://10.0.2.2:8080';
    } else {
      // iOS simulator can use localhost
      return 'http://localhost:8080';
    }
  } else {
    // Production mode - replace with your production API URL
    return 'https://your-production-api.com';
  }
};

/**
 * Other configuration constants
 */
export const config = {
  apiBaseUrl: getApiBaseUrl(),
  tokenKey: '@habit_tracker_token',
  apiTimeout: 30000, // 30 seconds
};
