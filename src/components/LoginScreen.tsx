/**
 * LoginScreen Component
 * Screen for user login and registration
 */

import React, {useState} from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  Alert,
} from 'react-native';
import {useAuth} from '../contexts/AuthContext';
import {colors, spacing, borderRadius, fontSize, shadow} from '../styles/theme';
import {
  validatePhoneNumber,
  validatePassword,
  validateHabitName,
} from '../utils/helpers';

/**
 * Login and Registration Screen
 */
export function LoginScreen() {
  const {login, register, error, clearError} = useAuth();

  const [isRegistering, setIsRegistering] = useState(false);
  const [name, setName] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  /**
   * Handles user login
   */
  const handleLogin = async () => {
    clearError();

    // Validate inputs
    if (!validatePhoneNumber(phoneNumber)) {
      Alert.alert(
        'Invalid Phone Number',
        'Please enter a valid phone number starting with + (e.g., +1234567890)',
      );
      return;
    }

    const passwordValidation = validatePassword(password);
    if (!passwordValidation.isValid) {
      Alert.alert('Invalid Password', passwordValidation.message || '');
      return;
    }

    setIsLoading(true);
    try {
      await login(phoneNumber, password);
    } catch (err) {
      // Error is handled by context and displayed below
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Handles user registration
   */
  const handleRegister = async () => {
    clearError();

    // Validate inputs
    const nameValidation = validateHabitName(name);
    if (!nameValidation.isValid) {
      Alert.alert('Invalid Name', 'Please enter your name');
      return;
    }

    if (!validatePhoneNumber(phoneNumber)) {
      Alert.alert(
        'Invalid Phone Number',
        'Please enter a valid phone number starting with + (e.g., +1234567890)',
      );
      return;
    }

    const passwordValidation = validatePassword(password);
    if (!passwordValidation.isValid) {
      Alert.alert('Invalid Password', passwordValidation.message || '');
      return;
    }

    setIsLoading(true);
    try {
      // Get timezone (simplified - using UTC for now)
      const timeZone = 'America/New_York';
      await register(name, phoneNumber, password, timeZone);
    } catch (err) {
      // Error is handled by context and displayed below
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Toggles between login and registration mode
   */
  const toggleMode = () => {
    setIsRegistering(!isRegistering);
    clearError();
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={styles.container}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        keyboardShouldPersistTaps="handled">
        <View style={styles.content}>
          {/* Title */}
          <Text style={styles.title}>Habit Tracker</Text>
          <Text style={styles.subtitle}>
            {isRegistering
              ? 'Create your account to get started'
              : 'Sign in to continue'}
          </Text>

          {/* Form */}
          <View style={styles.form}>
            {isRegistering && (
              <View style={styles.fieldContainer}>
                <Text style={styles.label}>Name</Text>
                <TextInput
                  style={styles.input}
                  placeholder="Enter your name"
                  placeholderTextColor={colors.textMuted}
                  value={name}
                  onChangeText={setName}
                  editable={!isLoading}
                  autoCapitalize="words"
                />
              </View>
            )}

            <View style={styles.fieldContainer}>
              <Text style={styles.label}>Phone Number</Text>
              <TextInput
                style={styles.input}
                placeholder="+1234567890"
                placeholderTextColor={colors.textMuted}
                value={phoneNumber}
                onChangeText={setPhoneNumber}
                editable={!isLoading}
                keyboardType="phone-pad"
                autoCapitalize="none"
              />
            </View>

            <View style={styles.fieldContainer}>
              <Text style={styles.label}>Password</Text>
              <TextInput
                style={styles.input}
                placeholder="Enter your password"
                placeholderTextColor={colors.textMuted}
                value={password}
                onChangeText={setPassword}
                editable={!isLoading}
                secureTextEntry
                autoCapitalize="none"
              />
            </View>

            {/* Error Message */}
            {error && (
              <View style={styles.errorContainer}>
                <Text style={styles.errorText}>{error}</Text>
              </View>
            )}

            {/* Submit Button */}
            <TouchableOpacity
              style={[styles.button, isLoading && styles.buttonDisabled]}
              onPress={isRegistering ? handleRegister : handleLogin}
              disabled={isLoading}>
              <Text style={styles.buttonText}>
                {isLoading
                  ? 'Loading...'
                  : isRegistering
                  ? 'Create Account'
                  : 'Sign In'}
              </Text>
            </TouchableOpacity>

            {/* Toggle Mode */}
            <TouchableOpacity
              style={styles.toggleButton}
              onPress={toggleMode}
              disabled={isLoading}>
              <Text style={styles.toggleText}>
                {isRegistering
                  ? 'Already have an account? Sign in'
                  : "Don't have an account? Register"}
              </Text>
            </TouchableOpacity>
          </View>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  scrollContent: {
    flexGrow: 1,
  },
  content: {
    flex: 1,
    justifyContent: 'center',
    padding: spacing.xl,
  },
  title: {
    fontSize: fontSize.xxxl,
    fontWeight: 'bold',
    color: colors.text,
    textAlign: 'center',
    marginBottom: spacing.sm,
  },
  subtitle: {
    fontSize: fontSize.md,
    color: colors.textLight,
    textAlign: 'center',
    marginBottom: spacing.xl,
  },
  form: {
    width: '100%',
  },
  fieldContainer: {
    marginBottom: spacing.lg,
  },
  label: {
    fontSize: fontSize.sm,
    fontWeight: '500',
    color: colors.text,
    marginBottom: spacing.sm,
  },
  input: {
    backgroundColor: colors.surfaceLight,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    fontSize: fontSize.md,
    color: colors.text,
    borderWidth: 1,
    borderColor: colors.border,
  },
  errorContainer: {
    backgroundColor: colors.error,
    borderRadius: borderRadius.sm,
    padding: spacing.md,
    marginBottom: spacing.lg,
  },
  errorText: {
    color: colors.text,
    fontSize: fontSize.sm,
    textAlign: 'center',
  },
  button: {
    backgroundColor: colors.primary,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    alignItems: 'center',
    ...shadow.sm,
  },
  buttonDisabled: {
    opacity: 0.5,
  },
  buttonText: {
    fontSize: fontSize.md,
    fontWeight: '600',
    color: colors.text,
  },
  toggleButton: {
    marginTop: spacing.lg,
    alignItems: 'center',
  },
  toggleText: {
    fontSize: fontSize.sm,
    color: colors.textLight,
  },
});
