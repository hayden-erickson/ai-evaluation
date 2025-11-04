/**
 * Authentication Context
 * Manages user authentication state and provides auth methods to the app
 */

import React, {
  createContext,
  useState,
  useContext,
  useEffect,
  ReactNode,
} from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';
import {User} from '../types';
import {authApi, ApiError} from '../services/api';

interface AuthContextType {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
  login: (phoneNumber: string, password: string) => Promise<void>;
  register: (
    name: string,
    phoneNumber: string,
    password: string,
    timeZone: string,
  ) => Promise<void>;
  logout: () => Promise<void>;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const TOKEN_KEY = '@habit_tracker_token';
const USER_KEY = '@habit_tracker_user';

/**
 * Authentication Provider Component
 */
export function AuthProvider({children}: {children: ReactNode}) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  /**
   * Load saved auth state on app startup
   */
  useEffect(() => {
    loadAuthState();
  }, []);

  /**
   * Loads authentication state from AsyncStorage
   */
  const loadAuthState = async () => {
    try {
      const savedToken = await AsyncStorage.getItem(TOKEN_KEY);
      const savedUser = await AsyncStorage.getItem(USER_KEY);

      if (savedToken && savedUser) {
        setToken(savedToken);
        setUser(JSON.parse(savedUser));
      }
    } catch (err) {
      console.error('Failed to load auth state:', err);
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Saves authentication state to AsyncStorage
   */
  const saveAuthState = async (authToken: string, authUser: User) => {
    try {
      await AsyncStorage.setItem(TOKEN_KEY, authToken);
      await AsyncStorage.setItem(USER_KEY, JSON.stringify(authUser));
    } catch (err) {
      console.error('Failed to save auth state:', err);
    }
  };

  /**
   * Clears authentication state from AsyncStorage
   */
  const clearAuthState = async () => {
    try {
      await AsyncStorage.removeItem(TOKEN_KEY);
      await AsyncStorage.removeItem(USER_KEY);
    } catch (err) {
      console.error('Failed to clear auth state:', err);
    }
  };

  /**
   * Login user with phone number and password
   */
  const login = async (phoneNumber: string, password: string) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await authApi.login({phone_number: phoneNumber, password});
      setToken(response.token);
      setUser(response.user);
      await saveAuthState(response.token, response.user);
    } catch (err) {
      const errorMessage =
        err instanceof ApiError
          ? err.message
          : 'Failed to login. Please try again.';
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Register a new user
   */
  const register = async (
    name: string,
    phoneNumber: string,
    password: string,
    timeZone: string,
  ) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await authApi.register({
        name,
        phone_number: phoneNumber,
        password,
        time_zone: timeZone,
      });
      setToken(response.token);
      setUser(response.user);
      await saveAuthState(response.token, response.user);
    } catch (err) {
      const errorMessage =
        err instanceof ApiError
          ? err.message
          : 'Failed to register. Please try again.';
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Logout the current user
   */
  const logout = async () => {
    setUser(null);
    setToken(null);
    await clearAuthState();
  };

  /**
   * Clear the current error
   */
  const clearError = () => {
    setError(null);
  };

  const value = {
    user,
    token,
    isLoading,
    error,
    login,
    register,
    logout,
    clearError,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

/**
 * Hook to use authentication context
 */
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
