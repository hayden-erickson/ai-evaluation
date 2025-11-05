import React, { useEffect, useState } from 'react';
import { Modal, StyleSheet, Text, TextInput, TouchableOpacity, View } from 'react-native';
import { colors, radius, spacing } from '../theme';
import { apiFetch } from '../api/client';
import type { CreateLogRequest, LogEntry, UpdateLogRequest } from '../types';

interface Props {
  visible: boolean;
  onClose: () => void;
  onSaved: (log: LogEntry) => void;
  habitId: number;
  log?: LogEntry | null;
}

/** Modal to create a new log or edit an existing log for a habit. */
export default function LogDetailsModal({ visible, onClose, onSaved, habitId, log }: Props) {
  const isEdit = !!log;
  const [notes, setNotes] = useState('');
  const [duration, setDuration] = useState<string>('');
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    setNotes(log?.notes || '');
    setDuration(log?.duration_seconds != null ? String(log.duration_seconds) : '');
  }, [log]);

  const save = async () => {
    if (busy) return;
    setBusy(true);
    try {
      if (isEdit) {
        const payload: UpdateLogRequest = {
          notes: notes.trim() || undefined,
          duration_seconds: duration ? Number(duration) : undefined,
        };
        const updated = await apiFetch<LogEntry>(`/logs/${log!.id}`, { method: 'PUT', body: JSON.stringify(payload) });
        onSaved(updated);
      } else {
        const payload: CreateLogRequest = {
          notes: notes.trim() || undefined,
          duration_seconds: duration ? Number(duration) : undefined,
        };
        const created = await apiFetch<LogEntry>(`/habits/${habitId}/logs`, { method: 'POST', body: JSON.stringify(payload) });
        onSaved(created);
      }
      onClose();
    } catch (e) {
      // Parent handles error messaging; keep modal UX simple
    } finally {
      setBusy(false);
    }
  };

  return (
    <Modal visible={visible} transparent animationType="slide" onRequestClose={onClose}>
      <View style={styles.overlay}>
        <View style={styles.card}>
          <Text style={styles.title}>{isEdit ? 'Edit Log' : 'New Log'}</Text>
          <TextInput placeholder="Notes" style={[styles.input, styles.multiline]} value={notes} onChangeText={setNotes} multiline />
          <TextInput placeholder="Duration (seconds)" style={styles.input} value={duration} onChangeText={setDuration} keyboardType="numeric" />

          <View style={styles.row}>
            <TouchableOpacity onPress={onClose} style={[styles.button, styles.cancel]}>
              <Text style={styles.cancelText}>Cancel</Text>
            </TouchableOpacity>
            <TouchableOpacity onPress={save} style={styles.button} disabled={busy}>
              <Text style={styles.saveText}>{busy ? 'Saving...' : 'Save'}</Text>
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  overlay: { flex: 1, backgroundColor: 'rgba(0,0,0,0.2)', alignItems: 'center', justifyContent: 'center' },
  card: { width: '90%', backgroundColor: colors.surface, borderRadius: radius.lg, padding: spacing.lg },
  title: { fontSize: 18, fontWeight: '700', color: colors.text, marginBottom: spacing.md },
  input: { backgroundColor: colors.muted, borderRadius: radius.md, padding: spacing.md, marginBottom: spacing.md, borderWidth: 1, borderColor: colors.border },
  multiline: { minHeight: 80, textAlignVertical: 'top' },
  row: { flexDirection: 'row', justifyContent: 'flex-end', gap: spacing.sm },
  button: { backgroundColor: colors.accent, padding: spacing.md, borderRadius: radius.md, minWidth: 96, alignItems: 'center' },
  saveText: { color: '#3730A3', fontWeight: '700' },
  cancel: { backgroundColor: colors.muted },
  cancelText: { color: colors.subtext, fontWeight: '600' },
});
