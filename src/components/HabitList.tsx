/**
 * HabitList Component
 * Displays a scrollable list of habits
 */

import React from 'react';
import {
  View,
  Text,
  FlatList,
  StyleSheet,
  ActivityIndicator,
  TouchableOpacity,
} from 'react-native';
import {Habit as HabitType} from '../types';
import {Habit} from './Habit';
import {colors, spacing, fontSize} from '../styles/theme';

interface HabitListProps {
  habits: HabitType[];
  token: string;
  isLoading: boolean;
  onEditHabit: (habit: HabitType) => void;
  onDeleteHabit: (habitId: number) => void;
  onRefresh: () => void;
  onAddNew: () => void;
}

/**
 * Component that displays a list of habits
 */
export function HabitList({
  habits,
  token,
  isLoading,
  onEditHabit,
  onDeleteHabit,
  onRefresh,
  onAddNew,
}: HabitListProps) {
  if (isLoading) {
    return (
      <View style={styles.centerContainer}>
        <ActivityIndicator size="large" color={colors.primary} />
        <Text style={styles.loadingText}>Loading habits...</Text>
      </View>
    );
  }

  if (habits.length === 0) {
    return (
      <View style={styles.centerContainer}>
        <Text style={styles.emptyTitle}>No habits yet</Text>
        <Text style={styles.emptySubtitle}>
          Start tracking your habits by creating one!
        </Text>
        <TouchableOpacity style={styles.emptyButton} onPress={onAddNew}>
          <Text style={styles.emptyButtonText}>+ Create Your First Habit</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <FlatList
      data={habits}
      keyExtractor={item => item.id.toString()}
      renderItem={({item}) => (
        <Habit
          habit={item}
          token={token}
          onEdit={() => onEditHabit(item)}
          onDelete={() => onDeleteHabit(item.id)}
          onUpdate={onRefresh}
        />
      )}
      contentContainerStyle={styles.listContainer}
      showsVerticalScrollIndicator={false}
      refreshing={isLoading}
      onRefresh={onRefresh}
    />
  );
}

const styles = StyleSheet.create({
  listContainer: {
    padding: spacing.lg,
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: spacing.xl,
  },
  loadingText: {
    fontSize: fontSize.md,
    color: colors.textLight,
    marginTop: spacing.md,
  },
  emptyTitle: {
    fontSize: fontSize.xl,
    fontWeight: '600',
    color: colors.text,
    marginBottom: spacing.sm,
  },
  emptySubtitle: {
    fontSize: fontSize.md,
    color: colors.textLight,
    textAlign: 'center',
    marginBottom: spacing.lg,
  },
  emptyButton: {
    backgroundColor: colors.primary,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderRadius: 12,
  },
  emptyButtonText: {
    fontSize: fontSize.md,
    fontWeight: '600',
    color: colors.text,
  },
});
