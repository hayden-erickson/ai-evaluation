import React, { useEffect, useState } from 'react';
import { ActivityIndicator, FlatList, RefreshControl, SafeAreaView, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { colors, radius, spacing } from '../theme';
import type { Habit } from '../types';
import { apiFetch } from '../api/client';
import HabitItem from './HabitItem';
import HabitDetailsModal from './HabitDetailsModal';
import { useAuth } from '../context/AuthContext';

/** Fetches and renders the user's habit list with create/edit/delete. */
export default function HabitList() {
  const { logout } = useAuth();
  const [habits, setHabits] = useState<Habit[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshing, setRefreshing] = useState(false);
  const [newHabitOpen, setNewHabitOpen] = useState(false);

  const load = async () => {
    setError(null);
    try {
      const data = await apiFetch<Habit[]>('/habits');
      setHabits(data);
    } catch (e: any) {
      setError(e.message || 'Failed to fetch habits');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  /** Adds a new habit to local state when created. */
  const handleSavedHabit = (h: Habit) => {
    setHabits((prev) => [h, ...prev]);
  };

  /** Updates habit in local list. */
  const handleUpdatedHabit = (h: Habit) => {
    setHabits((prev) => prev.map((x) => (x.id === h.id ? h : x)));
  };

  /** Removes habit from local list by id. */
  const handleDeletedHabit = (id: number) => {
    setHabits((prev) => prev.filter((x) => x.id !== id));
  };

  if (loading) {
    return (
      <SafeAreaView style={styles.center}>
        <ActivityIndicator />
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Your Habits</Text>
        <TouchableOpacity onPress={() => setNewHabitOpen(true)} style={styles.addBtn}>
          <Text style={styles.addText}>+ Add</Text>
        </TouchableOpacity>
      </View>

      {!!error && <Text style={styles.error}>{error}</Text>}

      <FlatList
        data={habits}
        keyExtractor={(item) => String(item.id)}
        contentContainerStyle={{ padding: spacing.lg }}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={() => { setRefreshing(true); load(); }} />}
        renderItem={({ item }) => (
          <HabitItem habit={item} onDeleted={handleDeletedHabit} onUpdated={handleUpdatedHabit} />
        )}
        ListEmptyComponent={<Text style={styles.empty}>No habits yet. Create your first one!</Text>}
      />

      <View style={styles.footer}>
        <TouchableOpacity onPress={logout}>
          <Text style={styles.logout}>Logout</Text>
        </TouchableOpacity>
      </View>

      <HabitDetailsModal visible={newHabitOpen} onClose={() => setNewHabitOpen(false)} onSaved={handleSavedHabit} />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  center: { flex: 1, alignItems: 'center', justifyContent: 'center', backgroundColor: colors.background },
  header: { padding: spacing.lg, flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between' },
  title: { fontSize: 24, fontWeight: '700', color: colors.text },
  addBtn: { backgroundColor: colors.primary, paddingHorizontal: spacing.lg, paddingVertical: spacing.sm, borderRadius: radius.lg },
  addText: { color: colors.primaryText, fontWeight: '700' },
  empty: { color: colors.subtext, textAlign: 'center', marginTop: spacing.xl },
  footer: { padding: spacing.lg, alignItems: 'center' },
  logout: { color: colors.subtext, textDecorationLine: 'underline' },
  error: { color: '#B91C1C', paddingHorizontal: spacing.lg },
});
