/**
 * Register screen component
 */

import React, {useState} from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
  Alert,
} from 'react-native';
import {useAuth} from '../services/AuthContext';
import {Button} from '../components/Button';
import {Input} from '../components/Input';
import {theme} from '../styles/theme';

interface RegisterScreenProps {
  onNavigateToLogin: () => void;
}

/**
 * Registration screen for new users
 */
export const RegisterScreen: React.FC<RegisterScreenProps> = ({
  onNavigateToLogin,
}) => {
  const {register, isLoading, error, clearError} = useAuth();
  const [name, setName] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errors, setErrors] = useState<{
    name?: string;
    phoneNumber?: string;
    password?: string;
    confirmPassword?: string;
  }>({});

  /**
   * Validate form inputs
   */
  const validate = (): boolean => {
    const newErrors: typeof errors = {};

    if (!name.trim()) {
      newErrors.name = 'Name is required';
    }

    if (!phoneNumber.trim()) {
      newErrors.phoneNumber = 'Phone number is required';
    }

    if (!password) {
      newErrors.password = 'Password is required';
    } else if (password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters';
    }

    if (!confirmPassword) {
      newErrors.confirmPassword = 'Please confirm your password';
    } else if (password !== confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  /**
   * Handle registration form submission
   */
  const handleRegister = async () => {
    clearError();

    if (!validate()) {
      return;
    }

    try {
      await register({
        name: name.trim(),
        phone_number: phoneNumber.trim(),
        password,
        time_zone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
      });
    } catch {
      // Error is already handled by AuthContext
      Alert.alert('Registration Failed', error || 'Please try again.');
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={styles.container}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        keyboardShouldPersistTaps="handled">
        <View style={styles.content}>
          <Text style={styles.title}>Create Account</Text>
          <Text style={styles.subtitle}>
            Start your journey to better habits
          </Text>

          <View style={styles.form}>
            <Input
              label="Name"
              value={name}
              onChangeText={text => {
                setName(text);
                setErrors({...errors, name: undefined});
              }}
              placeholder="Enter your name"
              autoCapitalize="words"
              error={errors.name}
            />

            <Input
              label="Phone Number"
              value={phoneNumber}
              onChangeText={text => {
                setPhoneNumber(text);
                setErrors({...errors, phoneNumber: undefined});
              }}
              placeholder="+1234567890"
              keyboardType="phone-pad"
              autoCapitalize="none"
              error={errors.phoneNumber}
            />

            <Input
              label="Password"
              value={password}
              onChangeText={text => {
                setPassword(text);
                setErrors({...errors, password: undefined});
              }}
              placeholder="At least 8 characters"
              secureTextEntry
              error={errors.password}
            />

            <Input
              label="Confirm Password"
              value={confirmPassword}
              onChangeText={text => {
                setConfirmPassword(text);
                setErrors({...errors, confirmPassword: undefined});
              }}
              placeholder="Re-enter your password"
              secureTextEntry
              error={errors.confirmPassword}
            />

            <Button
              title="Sign Up"
              onPress={handleRegister}
              loading={isLoading}
              style={styles.registerButton}
            />

            <View style={styles.loginContainer}>
              <Text style={styles.loginText}>Already have an account? </Text>
              <Button
                title="Sign In"
                onPress={onNavigateToLogin}
                variant="secondary"
              />
            </View>
          </View>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: theme.colors.background,
  },
  scrollContent: {
    flexGrow: 1,
    justifyContent: 'center',
  },
  content: {
    padding: theme.spacing.lg,
  },
  title: {
    fontSize: theme.fontSize.xxl,
    fontWeight: 'bold',
    color: theme.colors.text,
    marginBottom: theme.spacing.sm,
    textAlign: 'center',
  },
  subtitle: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textLight,
    marginBottom: theme.spacing.xl,
    textAlign: 'center',
  },
  form: {
    marginTop: theme.spacing.lg,
  },
  registerButton: {
    marginTop: theme.spacing.md,
  },
  loginContainer: {
    marginTop: theme.spacing.xl,
    alignItems: 'center',
  },
  loginText: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textLight,
    marginBottom: theme.spacing.sm,
  },
});
