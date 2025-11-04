import React, { createContext, useState, useEffect } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';
import api from '../api/api';
import { User } from '../types/types';

interface AuthContextType {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  login: (phoneNumber, password) => Promise<void>;
  logout: () => void;
  register: (name, phoneNumber, password, timeZone) => Promise<void>;
}

export const AuthContext = createContext<AuthContextType>({
  user: null,
  token: null,
  isLoading: true,
  login: async () => {},
  logout: () => {},
  register: async () => {},
});

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const bootstrapAsync = async () => {
      let userToken;
      try {
        userToken = await AsyncStorage.getItem('userToken');
        if (userToken) {
          // You might want to verify the token with your backend here
          const userResponse = await api.get('/users/me'); // A protected route to get user info
          setUser(userResponse.data);
          setToken(userToken);
        }
      } catch (e) {
        // Restoring token failed
        console.error('Restoring token failed', e);
      }
      setIsLoading(false);
    };

    bootstrapAsync();
  }, []);

  const authContext = {
    login: async (phoneNumber, password) => {
        try {
            const response = await api.post('/users/login', { phone_number: phoneNumber, password });
            const { token, user } = response.data;
            setToken(token);
            setUser(user);
            await AsyncStorage.setItem('userToken', token);
            api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        } catch (error) {
            console.error('Login failed', error);
            // Handle error, maybe show a message to the user
            throw error;
        }
    },
    logout: async () => {
      setToken(null);
      setUser(null);
      await AsyncStorage.removeItem('userToken');
      delete api.defaults.headers.common['Authorization'];
    },
    register: async (name, phoneNumber, password, timeZone) => {
        try {
            await api.post('/users/register', { name, phone_number: phoneNumber, password, time_zone: timeZone });
            // Optionally login the user automatically after registration
            await authContext.login(phoneNumber, password);
        } catch (error) {
            console.error('Registration failed', error);
            // Handle error
            throw error;
        }
    },
    user,
    token,
    isLoading,
  };

  return (
    <AuthContext.Provider value={authContext}>
      {children}
    </AuthContext.Provider>
  );
};
