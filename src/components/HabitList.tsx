import React, { useEffect, useState } from 'react';
import { Alert, FlatList, RefreshControl, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { useAuth } from '../auth/AuthContext';
import { CreateHabitRequest, Habit } from '../api/client';
import { colors, radius, spacing, typography } from '../theme';
import { HabitItem } from './HabitItem';
import { HabitDetailsModal } from './HabitDetailsModal';

export const HabitList: React.FC = () => {
	const { client, logout } = useAuth();
	const [habits, setHabits] = useState<Habit[]>([]);
	const [loading, setLoading] = useState(false);
	const [habitModalVisible, setHabitModalVisible] = useState(false);

	async function load() {
		try {
			setLoading(true);
			const data = await client.getHabits();
			setHabits(data);
		} catch (e: any) {
			Alert.alert('Error', e?.message ?? 'Failed to load habits');
		} finally {
			setLoading(false);
		}
	}

	useEffect(() => {
		load();
	}, []);

	async function createHabit(input: CreateHabitRequest) {
		const created = await client.createHabit(input);
		setHabits((prev) => [created, ...prev]);
	}

	return (
		<View style={styles.container}>
			<View style={styles.headerRow}>
				<Text style={styles.title}>Your Habits</Text>
				<TouchableOpacity onPress={logout} style={styles.logout}>
					<Text style={styles.logoutText}>Logout</Text>
				</TouchableOpacity>
			</View>
			<FlatList
				data={habits}
				keyExtractor={(h) => String(h.id)}
				renderItem={({ item }) => (
					<HabitItem
						habit={item}
						onDeleted={() => setHabits((prev) => prev.filter((h) => h.id !== item.id))}
						onUpdated={(updated) => setHabits((prev) => prev.map((h) => (h.id === updated.id ? updated : h)))}
					/>
				)}
				ItemSeparatorComponent={() => <View style={{ height: spacing.md }} />}
				refreshControl={<RefreshControl refreshing={loading} onRefresh={load} />}
				ListEmptyComponent={!loading ? <Text style={styles.empty}>No habits yet. Add one!</Text> : null}
				contentContainerStyle={{ paddingBottom: spacing.xxl }}
			/>
			<TouchableOpacity style={styles.fab} onPress={() => setHabitModalVisible(true)}>
				<Text style={styles.fabText}>＋</Text>
			</TouchableOpacity>

			<HabitDetailsModal visible={habitModalVisible} onClose={() => setHabitModalVisible(false)} onSave={createHabit} />
		</View>
	);
};

const styles = StyleSheet.create({
	container: { flex: 1, padding: spacing.lg, gap: spacing.md },
	headerRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: spacing.md },
	title: { fontSize: typography.title, color: colors.text, fontWeight: '700' },
	logout: { paddingVertical: spacing.sm, paddingHorizontal: spacing.md, borderRadius: radius.md, backgroundColor: colors.accent },
	logoutText: { color: colors.text, fontWeight: '600' },
	empty: { textAlign: 'center', color: colors.subtext, marginTop: spacing.xl },
	fab: {
		position: 'absolute',
		bottom: spacing.xl,
		right: spacing.xl,
		backgroundColor: colors.primary,
		width: 56,
		height: 56,
		borderRadius: 28,
		alignItems: 'center',
		justifyContent: 'center',
		elevation: 3,
		shadowColor: '#000',
		shadowOpacity: 0.15,
		shadowRadius: 8,
		shadowOffset: { width: 0, height: 2 },
	},
	fabText: { color: '#fff', fontSize: 28, lineHeight: 28, marginTop: -2 },
});


