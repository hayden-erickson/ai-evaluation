/**
 * Habit Tracker App
 * A React Native app for tracking daily habits and building streaks
 *
 * @format
 */

import React, {useState, useEffect} from 'react';
import {StatusBar, StyleSheet, View, ActivityIndicator} from 'react-native';
import {SafeAreaProvider} from 'react-native-safe-area-context';
import {NavigationContainer} from '@react-navigation/native';
import {
  createNativeStackNavigator,
  NativeStackScreenProps,
} from '@react-navigation/native-stack';
import {getToken} from './src/services/api';
import {LoginScreen} from './src/screens/LoginScreen';
import {RegisterScreen} from './src/screens/RegisterScreen';
import {HomeScreen} from './src/screens/HomeScreen';
import {colors} from './src/styles/theme';

type RootStackParamList = {
  Login: undefined;
  Register: undefined;
  Home: undefined;
};

type LoginScreenProps = NativeStackScreenProps<RootStackParamList, 'Login'>;
type RegisterScreenProps = NativeStackScreenProps<
  RootStackParamList,
  'Register'
>;
type HomeScreenProps = NativeStackScreenProps<RootStackParamList, 'Home'>;

const Stack = createNativeStackNavigator<RootStackParamList>();

/**
 * Main App component
 * Handles authentication state and navigation
 */
function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  /**
   * Check if user is already authenticated on app start
   */
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const token = await getToken();
        setIsAuthenticated(!!token);
      } catch (error) {
        console.error('Error checking authentication:', error);
        setIsAuthenticated(false);
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

  /**
   * Handle successful login/registration
   */
  const handleAuthSuccess = () => {
    setIsAuthenticated(true);
  };

  /**
   * Handle logout
   */
  const handleLogout = () => {
    setIsAuthenticated(false);
  };

  // Show loading screen while checking authentication
  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  return (
    <SafeAreaProvider>
      <StatusBar barStyle="dark-content" backgroundColor={colors.background} />
      <NavigationContainer>
        <Stack.Navigator
          screenOptions={{
            headerShown: false,
            animation: 'fade',
          }}>
          {isAuthenticated ? (
            // Authenticated screens
            <Stack.Screen name="Home">
              {(props: HomeScreenProps) => (
                <HomeScreen {...props} onLogout={handleLogout} />
              )}
            </Stack.Screen>
          ) : (
            // Authentication screens
            <>
              <Stack.Screen name="Login">
                {(props: LoginScreenProps) => (
                  <LoginScreen
                    {...props}
                    onLoginSuccess={handleAuthSuccess}
                    onNavigateToRegister={() =>
                      props.navigation.navigate('Register')
                    }
                  />
                )}
              </Stack.Screen>
              <Stack.Screen name="Register">
                {(props: RegisterScreenProps) => (
                  <RegisterScreen
                    {...props}
                    onRegisterSuccess={handleAuthSuccess}
                    onNavigateToLogin={() => props.navigation.navigate('Login')}
                  />
                )}
              </Stack.Screen>
            </>
          )}
        </Stack.Navigator>
      </NavigationContainer>
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
