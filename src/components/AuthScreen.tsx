import React, { useState } from 'react';
import { Alert, KeyboardAvoidingView, Platform, StyleSheet, Text, TextInput, TouchableOpacity, View } from 'react-native';
import { useAuth } from '../context/AuthContext';
import { colors, radius, spacing } from '../theme';

/** Minimal auth screen that toggles between Login and Register flows. */
export default function AuthScreen() {
  const { login, register, error } = useAuth();
  const [mode, setMode] = useState<'login' | 'register'>('login');
  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [busy, setBusy] = useState(false);

  const submit = async () => {
    if (busy) return;
    setBusy(true);
    try {
      if (mode === 'login') {
        await login(phone.trim(), password);
      } else {
        if (!name.trim()) throw new Error('Name is required');
        await register(name.trim(), phone.trim(), password, Intl.DateTimeFormat().resolvedOptions().timeZone);
      }
    } catch (e: any) {
      Alert.alert('Authentication Failed', e.message || 'Please try again.');
    } finally {
      setBusy(false);
    }
  };

  return (
    <KeyboardAvoidingView behavior={Platform.select({ ios: 'padding' })} style={styles.container}>
      <View style={styles.card}>
        <Text style={styles.title}>{mode === 'login' ? 'Welcome back' : 'Create account'}</Text>

        {mode === 'register' && (
          <TextInput
            placeholder="Name"
            style={styles.input}
            value={name}
            onChangeText={setName}
            autoCapitalize="words"
          />
        )}
        <TextInput
          placeholder="Phone Number"
          style={styles.input}
          value={phone}
          onChangeText={setPhone}
          autoCapitalize="none"
          keyboardType="phone-pad"
        />
        <TextInput
          placeholder="Password"
          style={styles.input}
          value={password}
          onChangeText={setPassword}
          autoCapitalize="none"
          secureTextEntry
        />

        {!!error && <Text style={styles.error}>{error}</Text>}

        <TouchableOpacity onPress={submit} style={[styles.button, busy && { opacity: 0.7 }]} disabled={busy}>
          <Text style={styles.buttonText}>{busy ? 'Please wait...' : mode === 'login' ? 'Login' : 'Register'}</Text>
        </TouchableOpacity>

        <TouchableOpacity onPress={() => setMode(mode === 'login' ? 'register' : 'login')}>
          <Text style={styles.link}>
            {mode === 'login' ? "Don't have an account? Register" : 'Already have an account? Login'}
          </Text>
        </TouchableOpacity>
      </View>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background, alignItems: 'center', justifyContent: 'center' },
  card: { width: '90%', backgroundColor: colors.surface, borderRadius: radius.lg, padding: spacing.xl, shadowColor: '#000', shadowOpacity: 0.05, shadowRadius: 10 },
  title: { fontSize: 24, fontWeight: '700', color: colors.text, marginBottom: spacing.lg },
  input: { backgroundColor: colors.muted, borderRadius: radius.md, padding: spacing.md, marginBottom: spacing.md, borderColor: colors.border, borderWidth: 1 },
  button: { backgroundColor: colors.primary, borderRadius: radius.lg, padding: spacing.md, alignItems: 'center', marginTop: spacing.sm },
  buttonText: { color: colors.primaryText, fontWeight: '700' },
  link: { color: colors.subtext, marginTop: spacing.md, textAlign: 'center' },
  error: { color: '#B91C1C', marginBottom: spacing.sm },
});
