/**
 * Main App Component - Habit Tracker
 * Manages authentication state and navigation between screens
 */

import React, { useState } from 'react';
import { StatusBar, StyleSheet, View } from 'react-native';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { User } from './src/types';
import { apiService } from './src/services/api';
import LoginScreen from './src/screens/LoginScreen';
import RegisterScreen from './src/screens/RegisterScreen';
import HomeScreen from './src/screens/HomeScreen';
import { Colors } from './src/styles/theme';

type Screen = 'login' | 'register' | 'home';

function App() {
  const [currentScreen, setCurrentScreen] = useState<Screen>('login');
  const [currentUser, setCurrentUser] = useState<User | null>(null);

  /**
   * Handle successful login
   * Stores user and token, navigates to home screen
   */
  const handleLoginSuccess = (user: User, token: string) => {
    setCurrentUser(user);
    apiService.setToken(token);
    setCurrentScreen('home');
  };

  /**
   * Handle successful registration
   * Navigates back to login screen
   */
  const handleRegisterSuccess = () => {
    setCurrentScreen('login');
  };

  /**
   * Handle user logout
   * Clears user data and navigates to login screen
   */
  const handleLogout = () => {
    setCurrentUser(null);
    apiService.setToken(null);
    setCurrentScreen('login');
  };

  return (
    <SafeAreaProvider>
      <View style={styles.container}>
        <StatusBar barStyle="dark-content" backgroundColor={Colors.background} />
        
        {/* Render current screen based on state */}
        {currentScreen === 'login' && (
          <LoginScreen
            onLoginSuccess={handleLoginSuccess}
            onNavigateToRegister={() => setCurrentScreen('register')}
          />
        )}
        
        {currentScreen === 'register' && (
          <RegisterScreen
            onRegisterSuccess={handleRegisterSuccess}
            onNavigateToLogin={() => setCurrentScreen('login')}
          />
        )}
        
        {currentScreen === 'home' && currentUser && (
          <HomeScreen
            user={currentUser}
            onLogout={handleLogout}
          />
        )}
      </View>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
});

export default App;
