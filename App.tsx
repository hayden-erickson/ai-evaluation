/**
 * Main App component for the Habit Tracker application
 * @format
 */

import React, {useState} from 'react';
import {StatusBar, StyleSheet, View, ActivityIndicator} from 'react-native';
import {SafeAreaProvider} from 'react-native-safe-area-context';
import {AuthProvider, useAuth} from './src/services/AuthContext';
import {LoginScreen} from './src/screens/LoginScreen';
import {RegisterScreen} from './src/screens/RegisterScreen';
import {HomeScreen} from './src/screens/HomeScreen';
import {theme} from './src/styles/theme';

/**
 * Main navigation component
 * Handles authentication flow and screen navigation
 */
function AppNavigator() {
  const {isAuthenticated, isLoading} = useAuth();
  const [showRegister, setShowRegister] = useState(false);

  // Show loading spinner while checking auth state
  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={theme.colors.primary} />
      </View>
    );
  }

  // Show appropriate screen based on auth state
  if (!isAuthenticated) {
    return showRegister ? (
      <RegisterScreen onNavigateToLogin={() => setShowRegister(false)} />
    ) : (
      <LoginScreen onNavigateToRegister={() => setShowRegister(true)} />
    );
  }

  return <HomeScreen />;
}

/**
 * Root App component
 */
function App() {
  return (
    <SafeAreaProvider>
      <AuthProvider>
        <StatusBar barStyle="dark-content" backgroundColor={theme.colors.background} />
        <View style={styles.container}>
          <AppNavigator />
        </View>
      </AuthProvider>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: theme.colors.background,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: theme.colors.background,
  },
});

export default App;
