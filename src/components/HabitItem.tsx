import React, { useEffect, useState } from 'react';
import { Alert, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { colors, radius, spacing } from '../theme';
import type { Habit, LogEntry } from '../types';
import { apiFetch } from '../api/client';
import StreakList from './StreakList';
import LogDetailsModal from './LogDetailsModal';
import HabitDetailsModal from './HabitDetailsModal';
import { computeStreak } from '../utils/streak';

interface Props {
  habit: Habit;
  onDeleted: (id: number) => void;
  onUpdated: (habit: Habit) => void;
}

/** Renders a habit row with actions and streak visuals. */
export default function HabitItem({ habit, onDeleted, onUpdated }: Props) {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [logModalOpen, setLogModalOpen] = useState(false);
  const [editHabitOpen, setEditHabitOpen] = useState(false);
  const [editingLog, setEditingLog] = useState<LogEntry | null>(null);

  // Load logs for this habit
  useEffect(() => {
    (async () => {
      try {
        const data = await apiFetch<LogEntry[]>(`/habits/${habit.id}/logs`);
        setLogs(data);
      } catch (e: any) {
        setError(e.message || 'Failed to load logs');
      }
    })();
  }, [habit.id]);

  /** Creates a new log by opening the log modal. */
  const createNewLog = () => {
    setEditingLog(null);
    setLogModalOpen(true);
  };

  /** Deletes the habit using the API and notifies parent. */
  const deleteHabit = async () => {
    Alert.alert('Delete Habit', 'Are you sure?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Delete',
        style: 'destructive',
        onPress: async () => {
          try {
            await apiFetch(`/habits/${habit.id}`, { method: 'DELETE' });
            onDeleted(habit.id);
          } catch (e: any) {
            setError(e.message || 'Failed to delete habit');
          }
        },
      },
    ]);
  };

  /** Saves a created or updated log and refreshes local state. */
  const handleSavedLog = (saved: LogEntry) => {
    setLogs((prev) => {
      const idx = prev.findIndex((l) => l.id === saved.id);
      if (idx >= 0) {
        const next = prev.slice();
        next[idx] = saved;
        return next;
      }
      return [saved, ...prev];
    });
  };

  /** Opens edit log modal by id. */
  const editLogById = async (id: number) => {
    try {
      const log = await apiFetch<LogEntry>(`/logs/${id}`);
      setEditingLog(log);
      setLogModalOpen(true);
    } catch (e: any) {
      setError(e.message || 'Failed to load log');
    }
  };

  const streak = computeStreak(logs);

  return (
    <View style={styles.card}>
      <View style={styles.row}>
        <View style={{ flex: 1 }}>
          <Text style={styles.name}>{habit.name}</Text>
          {!!habit.description && <Text style={styles.desc}>{habit.description}</Text>}
          <Text style={styles.streak}>Streak: {streak}</Text>
        </View>
        <View style={styles.actions}>
          <TouchableOpacity style={[styles.iconBtn, { backgroundColor: colors.accent }]} onPress={() => setEditHabitOpen(true)}>
            <Text style={styles.iconText}>Edit</Text>
          </TouchableOpacity>
          <TouchableOpacity style={[styles.iconBtn, { backgroundColor: colors.danger }]} onPress={deleteHabit}>
            <Text style={styles.iconText}>Delete</Text>
          </TouchableOpacity>
        </View>
      </View>

      <StreakList logs={logs} onEditLogById={editLogById} />

      <View style={styles.footer}>
        <TouchableOpacity style={styles.newLogBtn} onPress={createNewLog}>
          <Text style={styles.newLogText}>New Log</Text>
        </TouchableOpacity>
      </View>

      {!!error && <Text style={styles.error}>{error}</Text>}

      <LogDetailsModal
        visible={logModalOpen}
        onClose={() => setLogModalOpen(false)}
        onSaved={handleSavedLog}
        habitId={habit.id}
        log={editingLog}
      />

      <HabitDetailsModal
        visible={editHabitOpen}
        onClose={() => setEditHabitOpen(false)}
        onSaved={(h) => onUpdated(h)}
        habit={habit}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  card: { backgroundColor: colors.surface, borderRadius: radius.lg, padding: spacing.lg, marginBottom: spacing.lg, borderWidth: 1, borderColor: colors.border },
  row: { flexDirection: 'row', alignItems: 'center', marginBottom: spacing.sm },
  name: { fontSize: 18, fontWeight: '700', color: colors.text },
  desc: { color: colors.subtext, marginTop: 2 },
  streak: { marginTop: spacing.sm, color: colors.primaryText },
  actions: { flexDirection: 'row', gap: spacing.sm },
  iconBtn: { paddingHorizontal: spacing.md, paddingVertical: spacing.sm, borderRadius: radius.md },
  iconText: { color: '#111827', fontWeight: '700' },
  footer: { marginTop: spacing.sm },
  newLogBtn: { backgroundColor: colors.primary, borderRadius: radius.md, padding: spacing.md, alignItems: 'center' },
  newLogText: { color: colors.primaryText, fontWeight: '700' },
  error: { color: '#B91C1C', marginTop: spacing.sm },
});
