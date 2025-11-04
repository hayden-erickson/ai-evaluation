/**
 * Login screen component
 */

import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
  ActivityIndicator,
} from 'react-native';
import { useAuth } from '../contexts/AuthContext';
import { theme } from '../styles/theme';
import { commonStyles } from '../styles/commonStyles';

interface LoginScreenProps {
  onNavigateToRegister: () => void;
}

/**
 * Login screen for user authentication
 */
export const LoginScreen: React.FC<LoginScreenProps> = ({ onNavigateToRegister }) => {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [localError, setLocalError] = useState('');
  
  const { login, isLoading, error } = useAuth();

  /**
   * Validate form inputs
   */
  const validateForm = (): boolean => {
    if (!phoneNumber.trim()) {
      setLocalError('Phone number is required');
      return false;
    }
    if (!password.trim()) {
      setLocalError('Password is required');
      return false;
    }
    if (password.length < 8) {
      setLocalError('Password must be at least 8 characters');
      return false;
    }
    return true;
  };

  /**
   * Handle login submission
   */
  const handleLogin = async () => {
    setLocalError('');
    
    if (!validateForm()) {
      return;
    }

    try {
      await login({
        phone_number: phoneNumber.trim(),
        password: password,
      });
    } catch (err) {
      // Error is handled by AuthContext
    }
  };

  const displayError = localError || error;

  return (
    <SafeAreaView style={commonStyles.safeArea}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.container}
      >
        <ScrollView
          contentContainerStyle={styles.scrollContent}
          keyboardShouldPersistTaps="handled"
        >
          <View style={styles.content}>
            <Text style={styles.title}>Welcome Back</Text>
            <Text style={styles.subtitle}>Track your habits, build your streaks</Text>

            <View style={styles.form}>
              <View style={styles.inputGroup}>
                <Text style={commonStyles.label}>Phone Number</Text>
                <TextInput
                  style={commonStyles.input}
                  placeholder="Enter your phone number"
                  placeholderTextColor={theme.colors.textSecondary}
                  value={phoneNumber}
                  onChangeText={(text) => {
                    setPhoneNumber(text);
                    setLocalError('');
                  }}
                  keyboardType="phone-pad"
                  autoCapitalize="none"
                  editable={!isLoading}
                />
              </View>

              <View style={styles.inputGroup}>
                <Text style={commonStyles.label}>Password</Text>
                <TextInput
                  style={commonStyles.input}
                  placeholder="Enter your password"
                  placeholderTextColor={theme.colors.textSecondary}
                  value={password}
                  onChangeText={(text) => {
                    setPassword(text);
                    setLocalError('');
                  }}
                  secureTextEntry
                  editable={!isLoading}
                />
              </View>

              {displayError ? (
                <Text style={commonStyles.errorText}>{displayError}</Text>
              ) : null}

              <TouchableOpacity
                style={[commonStyles.button, styles.loginButton, isLoading && styles.buttonDisabled]}
                onPress={handleLogin}
                disabled={isLoading}
              >
                {isLoading ? (
                  <ActivityIndicator color={theme.colors.surface} />
                ) : (
                  <Text style={commonStyles.buttonText}>Login</Text>
                )}
              </TouchableOpacity>

              <View style={styles.registerContainer}>
                <Text style={styles.registerText}>Don't have an account? </Text>
                <TouchableOpacity onPress={onNavigateToRegister} disabled={isLoading}>
                  <Text style={styles.registerLink}>Register</Text>
                </TouchableOpacity>
              </View>
            </View>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
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
    textAlign: 'center',
    marginBottom: theme.spacing.sm,
  },
  
  subtitle: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textSecondary,
    textAlign: 'center',
    marginBottom: theme.spacing.xl,
  },
  
  form: {
    width: '100%',
  },
  
  inputGroup: {
    marginBottom: theme.spacing.md,
  },
  
  loginButton: {
    marginTop: theme.spacing.md,
  },
  
  buttonDisabled: {
    opacity: 0.6,
  },
  
  registerContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: theme.spacing.lg,
  },
  
  registerText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textSecondary,
  },
  
  registerLink: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.primary,
    fontWeight: '600',
  },
});

