import React, { useMemo, useState } from 'react';
import { Alert, KeyboardAvoidingView, Platform, SafeAreaView, ScrollView, StyleSheet, Text, TextInput, TouchableOpacity, View } from 'react-native';
import { useAuth } from '../auth/AuthContext';
import { colors, radius, spacing, typography } from '../theme';

export const AuthScreen: React.FC = () => {
	const { login, register } = useAuth();
	const [mode, setMode] = useState<'login' | 'register'>('login');
	const [loading, setLoading] = useState(false);
	const [name, setName] = useState('');
	const [phone, setPhone] = useState('');
	const [password, setPassword] = useState('');
	const [timeZone, setTimeZone] = useState<string>(
		(() => {
			try {
				return Intl.DateTimeFormat().resolvedOptions().timeZone ?? '';
			} catch {
				return '';
			}
		})()
	);

	const title = useMemo(() => (mode === 'login' ? 'Welcome back' : 'Create account'), [mode]);

	async function onSubmit() {
		if (loading) return;
		try {
			setLoading(true);
			if (mode === 'login') {
				if (!phone || !password) {
					Alert.alert('Missing fields', 'Phone number and password are required');
					return;
				}
				await login(phone, password);
			} else {
				if (!name || !phone || !password || !timeZone) {
					Alert.alert('Missing fields', 'Name, phone, password and time zone are required');
					return;
				}
				if (password.length < 8) {
					Alert.alert('Weak password', 'Password must be at least 8 characters');
					return;
				}
				await register({ name, phone_number: phone, password, time_zone: timeZone });
			}
		} catch (e: any) {
			const message = e?.message ?? 'Something went wrong';
			Alert.alert('Error', message);
		} finally {
			setLoading(false);
		}
	}

	return (
		<SafeAreaView style={styles.safe}>
			<KeyboardAvoidingView style={{ flex: 1 }} behavior={Platform.OS === 'ios' ? 'padding' : undefined}>
				<ScrollView contentContainerStyle={styles.container} keyboardShouldPersistTaps="handled">
					<Text style={styles.title}>{title}</Text>
					{mode === 'register' && (
						<View style={styles.fieldGroup}>
							<Text style={styles.label}>Name</Text>
							<TextInput value={name} onChangeText={setName} placeholder="Jane Doe" style={styles.input} placeholderTextColor={colors.subtext} />
						</View>
					)}
					<View style={styles.fieldGroup}>
						<Text style={styles.label}>Phone Number</Text>
						<TextInput
							value={phone}
							onChangeText={setPhone}
							keyboardType="phone-pad"
							placeholder="+1234567890"
							style={styles.input}
							placeholderTextColor={colors.subtext}
							autoCapitalize="none"
						/>
					</View>
					<View style={styles.fieldGroup}>
						<Text style={styles.label}>Password</Text>
						<TextInput
							value={password}
							onChangeText={setPassword}
							secureTextEntry
							placeholder="••••••••"
							style={styles.input}
							placeholderTextColor={colors.subtext}
						/>
					</View>
					{mode === 'register' && (
						<View style={styles.fieldGroup}>
							<Text style={styles.label}>Time Zone</Text>
							<TextInput value={timeZone} onChangeText={setTimeZone} placeholder="America/Los_Angeles" style={styles.input} placeholderTextColor={colors.subtext} />
						</View>
					)}

					<TouchableOpacity style={[styles.button, loading && { opacity: 0.6 }]} onPress={onSubmit} disabled={loading}>
						<Text style={styles.buttonText}>{mode === 'login' ? 'Login' : 'Create Account'}</Text>
					</TouchableOpacity>

					<TouchableOpacity
						style={styles.secondaryButton}
						onPress={() => setMode((m) => (m === 'login' ? 'register' : 'login'))}
						disabled={loading}
					>
						<Text style={styles.secondaryText}>
							{mode === 'login' ? "Don't have an account? Register" : 'Already have an account? Login'}
						</Text>
					</TouchableOpacity>
				</ScrollView>
			</KeyboardAvoidingView>
		</SafeAreaView>
	);
};

const styles = StyleSheet.create({
	safe: { flex: 1, backgroundColor: colors.background },
	container: { flexGrow: 1, padding: spacing.xl, gap: spacing.lg, justifyContent: 'center' },
	title: { fontSize: typography.title, color: colors.text, fontWeight: '600', textAlign: 'center', marginBottom: spacing.lg },
	fieldGroup: { gap: spacing.sm },
	label: { color: colors.subtext, fontSize: typography.small },
	input: {
		backgroundColor: colors.card,
		borderColor: colors.border,
		borderWidth: 1,
		borderRadius: radius.md,
		paddingHorizontal: spacing.lg,
		paddingVertical: spacing.md,
		color: colors.text,
		fontSize: typography.body,
	},
	button: {
		backgroundColor: colors.primary,
		paddingVertical: spacing.md,
		borderRadius: radius.lg,
		alignItems: 'center',
	},
	buttonText: { color: '#fff', fontWeight: '600', fontSize: typography.body },
	secondaryButton: { alignItems: 'center' },
	secondaryText: { color: colors.subtext },
});


