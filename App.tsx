/**
 * Habit Tracker App
 * Main application entry point
 */

import React, { useState } from 'react';
import { StatusBar } from 'react-native';
import { AuthProvider, useAuth } from './src/contexts/AuthContext';
import { LoginScreen } from './src/screens/LoginScreen';
import { RegisterScreen } from './src/screens/RegisterScreen';
import { HomeScreen } from './src/screens/HomeScreen';
import { theme } from './src/styles/theme';

type Screen = 'login' | 'register' | 'home';

/**
 * Main app navigation component
 */
const AppNavigator: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const [currentScreen, setCurrentScreen] = useState<Screen>('login');

  // Show home screen if authenticated
  if (isAuthenticated) {
    return <HomeScreen />;
  }

  // Show authentication screens
  if (currentScreen === 'register') {
    return <RegisterScreen onNavigateToLogin={() => setCurrentScreen('login')} />;
  }

  return <LoginScreen onNavigateToRegister={() => setCurrentScreen('register')} />;
};

/**
 * Root app component with providers
 */
function App(): React.JSX.Element {
  return (
    <AuthProvider>
      <StatusBar
        barStyle="dark-content"
        backgroundColor={theme.colors.background}
      />
      <AppNavigator />
    </AuthProvider>
  );
}

export default App;
