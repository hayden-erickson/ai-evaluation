/**
 * Registration screen component
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

interface RegisterScreenProps {
  onNavigateToLogin: () => void;
}

/**
 * Registration screen for new user signup
 */
export const RegisterScreen: React.FC<RegisterScreenProps> = ({ onNavigateToLogin }) => {
  const [name, setName] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [timeZone, setTimeZone] = useState('America/New_York');
  const [localError, setLocalError] = useState('');
  
  const { register, isLoading, error } = useAuth();

  /**
   * Validate form inputs
   */
  const validateForm = (): boolean => {
    if (!name.trim()) {
      setLocalError('Name is required');
      return false;
    }
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
    if (password !== confirmPassword) {
      setLocalError('Passwords do not match');
      return false;
    }
    if (!timeZone.trim()) {
      setLocalError('Time zone is required');
      return false;
    }
    return true;
  };

  /**
   * Handle registration submission
   */
  const handleRegister = async () => {
    setLocalError('');
    
    if (!validateForm()) {
      return;
    }

    try {
      await register({
        name: name.trim(),
        phone_number: phoneNumber.trim(),
        password: password,
        time_zone: timeZone.trim(),
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
            <Text style={styles.title}>Create Account</Text>
            <Text style={styles.subtitle}>Start tracking your habits today</Text>

            <View style={styles.form}>
              <View style={styles.inputGroup}>
                <Text style={commonStyles.label}>Name</Text>
                <TextInput
                  style={commonStyles.input}
                  placeholder="Enter your name"
                  placeholderTextColor={theme.colors.textSecondary}
                  value={name}
                  onChangeText={(text) => {
                    setName(text);
                    setLocalError('');
                  }}
                  editable={!isLoading}
                />
              </View>

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
                <Text style={commonStyles.label}>Time Zone</Text>
                <TextInput
                  style={commonStyles.input}
                  placeholder="e.g., America/New_York"
                  placeholderTextColor={theme.colors.textSecondary}
                  value={timeZone}
                  onChangeText={(text) => {
                    setTimeZone(text);
                    setLocalError('');
                  }}
                  editable={!isLoading}
                />
              </View>

              <View style={styles.inputGroup}>
                <Text style={commonStyles.label}>Password</Text>
                <TextInput
                  style={commonStyles.input}
                  placeholder="At least 8 characters"
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

              <View style={styles.inputGroup}>
                <Text style={commonStyles.label}>Confirm Password</Text>
                <TextInput
                  style={commonStyles.input}
                  placeholder="Re-enter your password"
                  placeholderTextColor={theme.colors.textSecondary}
                  value={confirmPassword}
                  onChangeText={(text) => {
                    setConfirmPassword(text);
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
                style={[commonStyles.button, styles.registerButton, isLoading && styles.buttonDisabled]}
                onPress={handleRegister}
                disabled={isLoading}
              >
                {isLoading ? (
                  <ActivityIndicator color={theme.colors.surface} />
                ) : (
                  <Text style={commonStyles.buttonText}>Register</Text>
                )}
              </TouchableOpacity>

              <View style={styles.loginContainer}>
                <Text style={styles.loginText}>Already have an account? </Text>
                <TouchableOpacity onPress={onNavigateToLogin} disabled={isLoading}>
                  <Text style={styles.loginLink}>Login</Text>
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
  
  registerButton: {
    marginTop: theme.spacing.md,
  },
  
  buttonDisabled: {
    opacity: 0.6,
  },
  
  loginContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: theme.spacing.lg,
  },
  
  loginText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textSecondary,
  },
  
  loginLink: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.primary,
    fontWeight: '600',
  },
});

