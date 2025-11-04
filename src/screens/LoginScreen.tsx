/**
 * LoginScreen - User authentication screen
 * Allows users to login or navigate to registration
 */

import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  Alert,
} from 'react-native';
import { Colors, Spacing, BorderRadius, Typography, CommonStyles } from '../styles/theme';
import { apiService } from '../services/api';
import { LoginResponse, User } from '../types';

interface LoginScreenProps {
  onLoginSuccess: (user: User, token: string) => void;
  onNavigateToRegister: () => void;
}

export default function LoginScreen({ onLoginSuccess, onNavigateToRegister }: LoginScreenProps) {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  /**
   * Validate form inputs
   * Returns error message if validation fails, empty string otherwise
   */
  const validateForm = (): string => {
    if (!phoneNumber.trim()) {
      return 'Phone number is required';
    }
    if (!password) {
      return 'Password is required';
    }
    return '';
  };

  /**
   * Handle login button press
   * Validates inputs and calls API
   */
  const handleLogin = async () => {
    try {
      setError('');
      
      // Validate form
      const validationError = validateForm();
      if (validationError) {
        setError(validationError);
        return;
      }

      setIsLoading(true);

      const response: LoginResponse = await apiService.login({
        phone_number: phoneNumber.trim(),
        password: password,
      });

      // Call success callback with user and token
      onLoginSuccess(response.user, response.token);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Login failed';
      setError(errorMessage);
      Alert.alert('Login Failed', errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      style={[CommonStyles.container, styles.container]}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
    >
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        keyboardShouldPersistTaps="handled"
      >
        {/* Header */}
        <View style={styles.header}>
          <Text style={styles.emoji}>📊</Text>
          <Text style={styles.title}>Habit Tracker</Text>
          <Text style={styles.subtitle}>Build better habits, one day at a time</Text>
        </View>

        {/* Login form */}
        <View style={styles.form}>
          <View style={styles.inputContainer}>
            <Text style={styles.label}>Phone Number</Text>
            <TextInput
              style={CommonStyles.input}
              placeholder="+1234567890"
              placeholderTextColor={Colors.textSecondary}
              value={phoneNumber}
              onChangeText={setPhoneNumber}
              keyboardType="phone-pad"
              autoCapitalize="none"
              autoCorrect={false}
              editable={!isLoading}
            />
          </View>

          <View style={styles.inputContainer}>
            <Text style={styles.label}>Password</Text>
            <TextInput
              style={CommonStyles.input}
              placeholder="Enter your password"
              placeholderTextColor={Colors.textSecondary}
              value={password}
              onChangeText={setPassword}
              secureTextEntry
              autoCapitalize="none"
              autoCorrect={false}
              editable={!isLoading}
            />
          </View>

          {/* Error message */}
          {error ? <Text style={CommonStyles.errorText}>{error}</Text> : null}

          {/* Login button */}
          <TouchableOpacity
            style={[CommonStyles.button, styles.loginButton]}
            onPress={handleLogin}
            disabled={isLoading}
          >
            {isLoading ? (
              <ActivityIndicator color={Colors.textPrimary} />
            ) : (
              <Text style={CommonStyles.buttonText}>Login</Text>
            )}
          </TouchableOpacity>

          {/* Register link */}
          <View style={styles.registerContainer}>
            <Text style={styles.registerText}>Don't have an account? </Text>
            <TouchableOpacity
              onPress={onNavigateToRegister}
              disabled={isLoading}
            >
              <Text style={styles.registerLink}>Register</Text>
            </TouchableOpacity>
          </View>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: Colors.background,
  },
  
  scrollContent: {
    flexGrow: 1,
    justifyContent: 'center',
    padding: Spacing.xl,
  },
  
  header: {
    alignItems: 'center',
    marginBottom: Spacing.xxl,
  },
  
  emoji: {
    fontSize: 64,
    marginBottom: Spacing.md,
  },
  
  title: {
    fontSize: Typography.h1,
    fontWeight: Typography.bold,
    color: Colors.textPrimary,
    marginBottom: Spacing.xs,
  },
  
  subtitle: {
    fontSize: Typography.body,
    color: Colors.textSecondary,
    textAlign: 'center',
  },
  
  form: {
    width: '100%',
  },
  
  inputContainer: {
    marginBottom: Spacing.md,
  },
  
  label: {
    fontSize: Typography.small,
    fontWeight: Typography.medium,
    color: Colors.textPrimary,
    marginBottom: Spacing.xs,
  },
  
  loginButton: {
    marginTop: Spacing.lg,
  },
  
  registerContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    marginTop: Spacing.lg,
  },
  
  registerText: {
    fontSize: Typography.body,
    color: Colors.textSecondary,
  },
  
  registerLink: {
    fontSize: Typography.body,
    color: Colors.primary,
    fontWeight: Typography.semibold,
  },
});
