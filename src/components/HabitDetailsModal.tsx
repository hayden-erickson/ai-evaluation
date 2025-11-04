import React, { useEffect, useState } from 'react';
import { Alert, Modal, StyleSheet, Text, TextInput, TouchableOpacity, View } from 'react-native';
import { colors, radius, spacing, typography } from '../theme';
import { CreateHabitRequest, Habit, UpdateHabitRequest } from '../api/client';

type Props = {
	visible: boolean;
	onClose: () => void;
	onSave: (input: CreateHabitRequest | UpdateHabitRequest) => Promise<void> | void;
	initial?: Habit | null;
};

export const HabitDetailsModal: React.FC<Props> = ({ visible, onClose, onSave, initial }) => {
	const [name, setName] = useState('');
	const [description, setDescription] = useState('');
	const [duration, setDuration] = useState<string>('');
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		if (initial) {
			setName(initial.name ?? '');
			setDescription(initial.description ?? '');
			setDuration(initial.duration_seconds != null ? String(initial.duration_seconds) : '');
		} else {
			setName('');
			setDescription('');
			setDuration('');
		}
	}, [initial, visible]);

	async function handleSave() {
		if (!name.trim()) {
			Alert.alert('Validation', 'Name is required');
			return;
		}
		const durationNumber = duration.trim() === '' ? null : Number(duration);
		if (durationNumber != null && (isNaN(durationNumber) || durationNumber < 0)) {
			Alert.alert('Validation', 'Duration must be a non-negative number');
			return;
		}
		setLoading(true);
		try {
			if (initial) {
				const payload: UpdateHabitRequest = {
					name: name,
					description: description,
					duration_seconds: durationNumber,
				};
				await onSave(payload);
			} else {
				const payload: CreateHabitRequest = {
					name: name.trim(),
					description: description.trim() || undefined,
					duration_seconds: durationNumber,
				};
				await onSave(payload);
			}
			onClose();
		} catch (e: any) {
			Alert.alert('Error', e?.message ?? 'Failed to save habit');
		} finally {
			setLoading(false);
		}
	}

	return (
		<Modal visible={visible} transparent animationType="slide" onRequestClose={onClose}>
			<View style={styles.backdrop}>
				<View style={styles.card}>
					<Text style={styles.title}>{initial ? 'Edit Habit' : 'New Habit'}</Text>
					<View style={styles.group}>
						<Text style={styles.label}>Name</Text>
						<TextInput style={styles.input} value={name} onChangeText={setName} placeholder="Habit name" placeholderTextColor={colors.subtext} />
					</View>
					<View style={styles.group}>
						<Text style={styles.label}>Description</Text>
						<TextInput
							style={styles.input}
							value={description}
							onChangeText={setDescription}
							placeholder="Optional"
							placeholderTextColor={colors.subtext}
						/>
					</View>
					<View style={styles.group}>
						<Text style={styles.label}>Duration (seconds)</Text>
						<TextInput
							style={styles.input}
							value={duration}
							onChangeText={setDuration}
							keyboardType="number-pad"
							placeholder="Optional"
							placeholderTextColor={colors.subtext}
						/>
					</View>
					<View style={styles.actions}>
						<TouchableOpacity onPress={onClose} style={[styles.button, styles.ghost]} disabled={loading}>
							<Text style={[styles.buttonText, { color: colors.text }]}>Cancel</Text>
						</TouchableOpacity>
						<TouchableOpacity onPress={handleSave} style={styles.button} disabled={loading}>
							<Text style={styles.buttonText}>{initial ? 'Save' : 'Create'}</Text>
						</TouchableOpacity>
					</View>
				</View>
			</View>
		</Modal>
	);
};

const styles = StyleSheet.create({
	backdrop: { flex: 1, backgroundColor: 'rgba(15,23,42,0.3)', justifyContent: 'flex-end' },
	card: { backgroundColor: colors.card, padding: spacing.xl, borderTopLeftRadius: radius.lg, borderTopRightRadius: radius.lg, gap: spacing.md },
	title: { fontSize: typography.subtitle, fontWeight: '600', color: colors.text, textAlign: 'center' },
	group: { gap: spacing.xs },
	label: { color: colors.subtext, fontSize: typography.small },
	input: {
		backgroundColor: colors.background,
		borderColor: colors.border,
		borderWidth: 1,
		borderRadius: radius.md,
		paddingHorizontal: spacing.lg,
		paddingVertical: spacing.md,
		color: colors.text,
		fontSize: typography.body,
	},
	actions: { flexDirection: 'row', gap: spacing.md, justifyContent: 'flex-end', marginTop: spacing.sm },
	button: { backgroundColor: colors.primary, paddingVertical: spacing.md, paddingHorizontal: spacing.xl, borderRadius: radius.md },
	ghost: { backgroundColor: 'transparent', borderWidth: 1, borderColor: colors.border },
	buttonText: { color: '#fff', fontWeight: '600' },
});


