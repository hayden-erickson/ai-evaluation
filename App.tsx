/**
 * Habit Tracker App
 * Main entry point for the React Native application
 *
 * @format
 */

import React from 'react';
import {StatusBar, StyleSheet, View, ActivityIndicator} from 'react-native';
import {SafeAreaProvider} from 'react-native-safe-area-context';
import {AuthProvider, useAuth} from './src/contexts/AuthContext';
import {LoginScreen} from './src/components/LoginScreen';
import {HomeScreen} from './src/components/HomeScreen';
import {colors} from './src/styles/theme';

/**
 * Main App Component with Authentication Flow
 */
function AppContent() {
  const {user, isLoading} = useAuth();

  // Show loading indicator while checking auth state
  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  // Show appropriate screen based on auth state
  return user ? <HomeScreen /> : <LoginScreen />;
}

/**
 * Root App Component
 */
function App() {
  return (
    <SafeAreaProvider>
      <AuthProvider>
        <StatusBar barStyle="dark-content" backgroundColor={colors.background} />
        <AppContent />
      </AuthProvider>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.background,
  },
});

export default App;
