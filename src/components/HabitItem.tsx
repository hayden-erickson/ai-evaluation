import React, { useEffect, useMemo, useState } from 'react';
import { Alert, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { Habit, Log } from '../api/client';
import { useAuth } from '../auth/AuthContext';
import { colors, radius, spacing, typography } from '../theme';
import { LogDetailsModal } from './LogDetailsModal';
import { HabitDetailsModal } from './HabitDetailsModal';

type Props = {
	habit: Habit;
	onDeleted: () => void;
	onUpdated: (habit: Habit) => void;
};

export const HabitItem: React.FC<Props> = ({ habit, onDeleted, onUpdated }) => {
	const { client } = useAuth();
	const [logs, setLogs] = useState<Log[] | null>(null);
	const [logModalVisible, setLogModalVisible] = useState(false);
	const [editLog, setEditLog] = useState<Log | null>(null);
	const [habitModalVisible, setHabitModalVisible] = useState(false);
	const [expanded, setExpanded] = useState(true);

	useEffect(() => {
		let mounted = true;
		(async () => {
			try {
				const data = await client.getHabitLogs(habit.id);
				if (mounted) setLogs(data);
			} catch (e: any) {
				Alert.alert('Error', e?.message ?? 'Failed to load logs');
			}
		})();
		return () => {
			mounted = false;
		};
	}, [client, habit.id]);

	const days = useMemo(() => buildRecentDays(14), []);
	const logsByDay = useMemo(() => groupLogsByDay(logs ?? []), [logs]);
	const streak = useMemo(() => computeStreakWithOneSkip(logs ?? []), [logs]);

	async function handleCreateLog(input: { notes?: string; duration_seconds?: number | null }) {
		const created = await client.createLog(habit.id, input);
		setLogs((prev) => [created, ...(prev ?? [])]);
	}

	async function handleUpdateLog(input: { notes?: string; duration_seconds?: number | null }) {
		if (!editLog) return;
		const updated = await client.updateLog(editLog.id, input);
		setLogs((prev) => (prev ?? []).map((l) => (l.id === updated.id ? updated : l)));
	}

	async function handleDeleteLog(id: number) {
		await client.deleteLog(id);
		setLogs((prev) => (prev ?? []).filter((l) => l.id !== id));
	}

	async function handleDeleteHabit() {
		Alert.alert('Delete Habit', 'Are you sure you want to delete this habit?', [
			{ text: 'Cancel', style: 'cancel' },
			{
				text: 'Delete',
				style: 'destructive',
				onPress: async () => {
					await client.deleteHabit(habit.id);
					onDeleted();
				},
			},
		]);
	}

	async function handleSaveHabit(input: { name?: string; description?: string; duration_seconds?: number | null }) {
		if (habit) {
			const updated = await client.updateHabit(habit.id, input);
			onUpdated(updated);
		}
	}

	return (
		<View style={styles.card}>
			<View style={styles.header}>
				<View style={{ flex: 1 }}>
					<Text style={styles.name}>{habit.name}</Text>
					{!!habit.description && <Text style={styles.description}>{habit.description}</Text>}
				</View>
				<View style={styles.headerActions}>
					<TouchableOpacity onPress={() => setHabitModalVisible(true)} style={[styles.smallButton, { backgroundColor: colors.accent }]}>
						<Text style={styles.smallButtonText}>Edit</Text>
					</TouchableOpacity>
					<TouchableOpacity onPress={handleDeleteHabit} style={[styles.smallButton, { backgroundColor: colors.error }]}>
						<Text style={styles.smallButtonText}>Delete</Text>
					</TouchableOpacity>
				</View>
			</View>
			<View style={styles.streakRow}>
				<Text style={styles.streak}>🔥 {streak} day streak</Text>
				<TouchableOpacity onPress={() => setLogModalVisible(true)} style={styles.newLogButton}>
					<Text style={styles.newLogText}>New Log</Text>
				</TouchableOpacity>
			</View>
			<TouchableOpacity onPress={() => setExpanded((e) => !e)}>
				<Text style={styles.sectionTitle}>{expanded ? 'Hide' : 'Show'} recent days</Text>
			</TouchableOpacity>
			{expanded && (
				<View style={styles.daysRow}>
					{days.map((d) => {
						const key = d.key;
						const log = logsByDay.get(key)?.[0];
						return (
							<TouchableOpacity
								key={key}
								style={[styles.day, log ? styles.dayFilled : styles.dayEmpty]}
								onPress={() => {
									if (log) {
										setEditLog(log);
										setLogModalVisible(true);
									} else {
										setEditLog(null);
										setLogModalVisible(true);
									}
								}}
							>
								<Text style={styles.dayText}>{d.label}</Text>
							</TouchableOpacity>
						);
					})}
				</View>
			)}

			<LogDetailsModal
				visible={logModalVisible}
				onClose={() => setLogModalVisible(false)}
				onSave={async (input) => {
					if (editLog) {
						await handleUpdateLog(input);
					} else {
						await handleCreateLog(input);
					}
				}}
				initial={editLog}
			/>

			<HabitDetailsModal
				visible={habitModalVisible}
				onClose={() => setHabitModalVisible(false)}
				onSave={handleSaveHabit}
				initial={habit}
			/>
		</View>
	);
};

function buildRecentDays(count: number): { key: string; label: string; date: Date }[] {
	const today = new Date();
	const out: { key: string; label: string; date: Date }[] = [];
	for (let i = 0; i < count; i++) {
		const d = new Date(today);
		d.setDate(today.getDate() - i);
		const key = dateKey(d);
		out.push({ key, label: d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' }), date: d });
	}
	return out.reverse();
}

function groupLogsByDay(logs: Log[]): Map<string, Log[]> {
	const map = new Map<string, Log[]>();
	for (const log of logs) {
		const d = new Date(log.created_at);
		const key = dateKey(d);
		const arr = map.get(key) ?? [];
		arr.push(log);
		map.set(key, arr);
	}
	return map;
}

function dateKey(d: Date): string {
	return `${d.getFullYear()}-${d.getMonth() + 1}-${d.getDate()}`;
}

// Computes current streak up to today allowing one missed day between logs
function computeStreakWithOneSkip(logs: Log[]): number {
	if (!logs.length) return 0;
	const days = Array.from(new Set(logs.map((l) => dateKey(new Date(l.created_at)))));
	days.sort((a, b) => (a > b ? -1 : 1));

	let streak = 0;
	let skipsRemaining = 1;
	let currentDate = new Date();
	let cursorKey = dateKey(currentDate);

	while (true) {
		if (days.includes(cursorKey)) {
			streak += 1;
			currentDate.setDate(currentDate.getDate() - 1);
			cursorKey = dateKey(currentDate);
			continue;
		}
		if (skipsRemaining > 0) {
			skipsRemaining -= 1;
			currentDate.setDate(currentDate.getDate() - 1);
			cursorKey = dateKey(currentDate);
			continue;
		}
		break;
	}
	return streak;
}

const styles = StyleSheet.create({
	card: { backgroundColor: colors.card, borderRadius: radius.lg, padding: spacing.lg, gap: spacing.md, borderWidth: 1, borderColor: colors.border },
	header: { flexDirection: 'row', alignItems: 'center', gap: spacing.md },
	headerActions: { flexDirection: 'row', gap: spacing.sm },
	name: { fontSize: typography.subtitle, color: colors.text, fontWeight: '600' },
	description: { color: colors.subtext },
	streakRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between' },
	streak: { color: colors.text, fontWeight: '600' },
	newLogButton: { backgroundColor: colors.secondary, paddingVertical: spacing.sm, paddingHorizontal: spacing.lg, borderRadius: radius.md },
	newLogText: { color: colors.text, fontWeight: '600' },
	sectionTitle: { color: colors.subtext },
	daysRow: { flexDirection: 'row', flexWrap: 'wrap', gap: spacing.sm },
	day: { borderRadius: radius.md, borderWidth: 1, borderColor: colors.border, paddingVertical: spacing.xs, paddingHorizontal: spacing.sm },
	dayFilled: { backgroundColor: colors.success, borderColor: colors.success },
	dayEmpty: { backgroundColor: colors.background },
	dayText: { color: colors.text, fontSize: typography.small },
});


