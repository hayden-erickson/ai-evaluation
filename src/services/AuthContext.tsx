/**
 * Authentication context for managing user state
 */

import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from 'react';
import {User, LoginRequest, RegisterRequest} from '../types';
import * as api from '../services/api';
import {ApiError} from '../services/api';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  error: string | null;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

/**
 * Authentication provider component
 * Manages user authentication state and provides auth methods
 */
export const AuthProvider: React.FC<AuthProviderProps> = ({children}) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Check for existing token on mount
  useEffect(() => {
    checkAuth();
  }, []);

  /**
   * Check if user is already authenticated
   * Note: In a production app, you should validate the token with the backend
   * or decode the JWT to check expiration
   */
  const checkAuth = async () => {
    try {
      const token = await api.getToken();
      if (token) {
        // Token exists - in a real app, validate it with the backend or decode JWT
        // For now, we'll require re-login on app restart for security
        setIsLoading(false);
      } else {
        setIsLoading(false);
      }
    } catch (err) {
      console.error('Auth check error:', err);
      setIsLoading(false);
    }
  };

  /**
   * Login with phone number and password
   */
  const login = async (credentials: LoginRequest) => {
    try {
      setError(null);
      setIsLoading(true);
      const response = await api.login(credentials);
      await api.setToken(response.token);
      setUser(response.user);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('An unexpected error occurred. Please try again.');
      }
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Register a new user
   */
  const register = async (data: RegisterRequest) => {
    try {
      setError(null);
      setIsLoading(true);
      // Register returns just the user, not a token
      await api.register(data);
      // Automatically login after successful registration
      await login({
        phone_number: data.phone_number,
        password: data.password,
      });
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('An unexpected error occurred. Please try again.');
      }
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Logout the current user
   */
  const logout = async () => {
    try {
      await api.removeToken();
      setUser(null);
    } catch (err) {
      console.error('Logout error:', err);
    }
  };

  /**
   * Clear any error messages
   */
  const clearError = () => {
    setError(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: user !== null,
        login,
        register,
        logout,
        error,
        clearError,
      }}>
      {children}
    </AuthContext.Provider>
  );
};

/**
 * Hook to use authentication context
 */
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
