/**
 * RegisterScreen - User registration screen
 * Allows new users to create an account
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
import { User } from '../types';

interface RegisterScreenProps {
  onRegisterSuccess: () => void;
  onNavigateToLogin: () => void;
}

export default function RegisterScreen({ onRegisterSuccess, onNavigateToLogin }: RegisterScreenProps) {
  const [name, setName] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [timeZone, setTimeZone] = useState('America/New_York');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  /**
   * Validate form inputs
   * Returns error message if validation fails, empty string otherwise
   */
  const validateForm = (): string => {
    if (!name.trim()) {
      return 'Name is required';
    }
    if (name.trim().length < 2) {
      return 'Name must be at least 2 characters';
    }
    if (!phoneNumber.trim()) {
      return 'Phone number is required';
    }
    if (!password) {
      return 'Password is required';
    }
    if (password.length < 8) {
      return 'Password must be at least 8 characters';
    }
    if (password !== confirmPassword) {
      return 'Passwords do not match';
    }
    if (!timeZone.trim()) {
      return 'Time zone is required';
    }
    return '';
  };

  /**
   * Handle register button press
   * Validates inputs and calls API
   */
  const handleRegister = async () => {
    try {
      setError('');
      
      // Validate form
      const validationError = validateForm();
      if (validationError) {
        setError(validationError);
        return;
      }

      setIsLoading(true);

      await apiService.register({
        name: name.trim(),
        phone_number: phoneNumber.trim(),
        password: password,
        time_zone: timeZone.trim(),
      });

      Alert.alert(
        'Success',
        'Account created successfully! Please login.',
        [{ text: 'OK', onPress: onRegisterSuccess }]
      );
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Registration failed';
      setError(errorMessage);
      Alert.alert('Registration Failed', errorMessage);
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
          <Text style={styles.emoji}>✨</Text>
          <Text style={styles.title}>Create Account</Text>
          <Text style={styles.subtitle}>Start your habit tracking journey</Text>
        </View>

        {/* Registration form */}
        <View style={styles.form}>
          <View style={styles.inputContainer}>
            <Text style={styles.label}>Name</Text>
            <TextInput
              style={CommonStyles.input}
              placeholder="John Doe"
              placeholderTextColor={Colors.textSecondary}
              value={name}
              onChangeText={setName}
              autoCapitalize="words"
              autoCorrect={false}
              editable={!isLoading}
            />
          </View>

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
            <Text style={styles.label}>Password (min 8 characters)</Text>
            <TextInput
              style={CommonStyles.input}
              placeholder="Enter password"
              placeholderTextColor={Colors.textSecondary}
              value={password}
              onChangeText={setPassword}
              secureTextEntry
              autoCapitalize="none"
              autoCorrect={false}
              editable={!isLoading}
            />
          </View>

          <View style={styles.inputContainer}>
            <Text style={styles.label}>Confirm Password</Text>
            <TextInput
              style={CommonStyles.input}
              placeholder="Re-enter password"
              placeholderTextColor={Colors.textSecondary}
              value={confirmPassword}
              onChangeText={setConfirmPassword}
              secureTextEntry
              autoCapitalize="none"
              autoCorrect={false}
              editable={!isLoading}
            />
          </View>

          <View style={styles.inputContainer}>
            <Text style={styles.label}>Time Zone</Text>
            <TextInput
              style={CommonStyles.input}
              placeholder="America/New_York"
              placeholderTextColor={Colors.textSecondary}
              value={timeZone}
              onChangeText={setTimeZone}
              autoCapitalize="none"
              autoCorrect={false}
              editable={!isLoading}
            />
          </View>

          {/* Error message */}
          {error ? <Text style={CommonStyles.errorText}>{error}</Text> : null}

          {/* Register button */}
          <TouchableOpacity
            style={[CommonStyles.button, styles.registerButton]}
            onPress={handleRegister}
            disabled={isLoading}
          >
            {isLoading ? (
              <ActivityIndicator color={Colors.textPrimary} />
            ) : (
              <Text style={CommonStyles.buttonText}>Register</Text>
            )}
          </TouchableOpacity>

          {/* Login link */}
          <View style={styles.loginContainer}>
            <Text style={styles.loginText}>Already have an account? </Text>
            <TouchableOpacity
              onPress={onNavigateToLogin}
              disabled={isLoading}
            >
              <Text style={styles.loginLink}>Login</Text>
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
    paddingTop: Spacing.xxl,
    paddingBottom: Spacing.xxl,
  },
  
  header: {
    alignItems: 'center',
    marginBottom: Spacing.xl,
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
  
  registerButton: {
    marginTop: Spacing.lg,
  },
  
  loginContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    marginTop: Spacing.lg,
  },
  
  loginText: {
    fontSize: Typography.body,
    color: Colors.textSecondary,
  },
  
  loginLink: {
    fontSize: Typography.body,
    color: Colors.primary,
    fontWeight: Typography.semibold,
  },
});
