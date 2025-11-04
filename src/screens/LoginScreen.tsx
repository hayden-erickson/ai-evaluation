/**
 * Login screen component
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

interface LoginScreenProps {
  onNavigateToRegister: () => void;
}

/**
 * Login screen with phone number and password authentication
 */
export const LoginScreen: React.FC<LoginScreenProps> = ({
  onNavigateToRegister,
}) => {
  const {login, isLoading, error, clearError} = useAuth();
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<{
    phoneNumber?: string;
    password?: string;
  }>({});

  /**
   * Validate form inputs
   */
  const validate = (): boolean => {
    const newErrors: typeof errors = {};

    if (!phoneNumber.trim()) {
      newErrors.phoneNumber = 'Phone number is required';
    }

    if (!password) {
      newErrors.password = 'Password is required';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  /**
   * Handle login form submission
   */
  const handleLogin = async () => {
    clearError();
    
    if (!validate()) {
      return;
    }

    try {
      await login({
        phone_number: phoneNumber.trim(),
        password,
      });
    } catch (err) {
      // Error is already handled by AuthContext
      Alert.alert('Login Failed', error || 'Please check your credentials and try again.');
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
          <Text style={styles.title}>Welcome Back</Text>
          <Text style={styles.subtitle}>Sign in to continue tracking your habits</Text>

          <View style={styles.form}>
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
              placeholder="Enter your password"
              secureTextEntry
              error={errors.password}
            />

            <Button
              title="Sign In"
              onPress={handleLogin}
              loading={isLoading}
              style={styles.loginButton}
            />

            <View style={styles.registerContainer}>
              <Text style={styles.registerText}>Don't have an account? </Text>
              <Button
                title="Sign Up"
                onPress={onNavigateToRegister}
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
  loginButton: {
    marginTop: theme.spacing.md,
  },
  registerContainer: {
    marginTop: theme.spacing.xl,
    alignItems: 'center',
  },
  registerText: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textLight,
    marginBottom: theme.spacing.sm,
  },
});
