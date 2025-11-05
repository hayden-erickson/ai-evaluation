import React, { useEffect, useState } from 'react';
import { Modal, StyleSheet, Text, TextInput, TouchableOpacity, View } from 'react-native';
import { colors, radius, spacing } from '../theme';
import { apiFetch } from '../api/client';
import type { CreateHabitRequest, Habit, UpdateHabitRequest } from '../types';

interface Props {
  visible: boolean;
  onClose: () => void;
  onSaved: (habit: Habit) => void;
  habit?: Habit | null;
}

/** Modal to create or update a habit. */
export default function HabitDetailsModal({ visible, onClose, onSaved, habit }: Props) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [busy, setBusy] = useState(false);
  const isEdit = !!habit;

  useEffect(() => {
    setName(habit?.name || '');
    setDescription(habit?.description || '');
  }, [habit]);

  const save = async () => {
    if (busy) return;
    setBusy(true);
    try {
      if (isEdit) {
        const payload: UpdateHabitRequest = { name: name.trim() || undefined, description: description.trim() || undefined };
        const updated = await apiFetch<Habit>(`/habits/${habit!.id}`, { method: 'PUT', body: JSON.stringify(payload) });
        onSaved(updated);
      } else {
        const payload: CreateHabitRequest = { name: name.trim(), description: description.trim() || undefined };
        const created = await apiFetch<Habit>('/habits', { method: 'POST', body: JSON.stringify(payload) });
        onSaved(created);
      }
      onClose();
    } catch (e) {
      // Errors are shown through parent or could add local message
    } finally {
      setBusy(false);
    }
  };

  return (
    <Modal visible={visible} transparent animationType="slide" onRequestClose={onClose}>
      <View style={styles.overlay}>
        <View style={styles.card}>
          <Text style={styles.title}>{isEdit ? 'Edit Habit' : 'New Habit'}</Text>
          <TextInput placeholder="Name" style={styles.input} value={name} onChangeText={setName} />
          <TextInput placeholder="Description" style={styles.input} value={description} onChangeText={setDescription} />

          <View style={styles.row}>
            <TouchableOpacity onPress={onClose} style={[styles.button, styles.cancel]}> 
              <Text style={styles.cancelText}>Cancel</Text>
            </TouchableOpacity>
            <TouchableOpacity onPress={save} style={styles.button} disabled={busy || !name.trim()}>
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
  row: { flexDirection: 'row', justifyContent: 'flex-end', gap: spacing.sm },
  button: { backgroundColor: colors.primary, padding: spacing.md, borderRadius: radius.md, minWidth: 96, alignItems: 'center' },
  saveText: { color: colors.primaryText, fontWeight: '700' },
  cancel: { backgroundColor: colors.muted },
  cancelText: { color: colors.subtext, fontWeight: '600' },
});
