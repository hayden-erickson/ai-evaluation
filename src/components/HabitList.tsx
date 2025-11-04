/**
 * HabitList component - Displays a list of all user habits
 * Manages the list of habits and allows adding new ones
 */

import React from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  ScrollView,
  StyleSheet,
  ActivityIndicator,
  RefreshControl,
} from 'react-native';
import { Habit as HabitType } from '../types';
import { Colors, Spacing, BorderRadius, Typography, CommonStyles } from '../styles/theme';
import Habit from './Habit';

interface HabitListProps {
  habits: HabitType[];
  isLoading: boolean;
  onRefresh: () => Promise<void>;
  onAddNew: () => void;
  onEditHabit: (habit: HabitType) => void;
  onDeleteHabit: (habit: HabitType) => void;
}

export default function HabitList({
  habits,
  isLoading,
  onRefresh,
  onAddNew,
  onEditHabit,
  onDeleteHabit,
}: HabitListProps) {
  const [refreshing, setRefreshing] = React.useState(false);

  /**
   * Handle pull-to-refresh
   */
  const handleRefresh = async () => {
    setRefreshing(true);
    await onRefresh();
    setRefreshing(false);
  };

  // Show loading state on initial load
  if (isLoading && habits.length === 0) {
    return (
      <View style={[CommonStyles.container, styles.centered]}>
        <ActivityIndicator size="large" color={Colors.primary} />
        <Text style={styles.loadingText}>Loading habits...</Text>
      </View>
    );
  }

  // Show empty state if no habits
  if (!isLoading && habits.length === 0) {
    return (
      <View style={[CommonStyles.container, styles.container]}>
        <ScrollView
          contentContainerStyle={styles.emptyContainer}
          refreshControl={
            <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />
          }
        >
          <Text style={styles.emptyEmoji}>📝</Text>
          <Text style={styles.emptyTitle}>No habits yet</Text>
          <Text style={styles.emptyText}>
            Start building positive habits by adding your first one!
          </Text>
          <TouchableOpacity
            style={[CommonStyles.button, styles.addButton]}
            onPress={onAddNew}
          >
            <Text style={CommonStyles.buttonText}>Add Your First Habit</Text>
          </TouchableOpacity>
        </ScrollView>
      </View>
    );
  }

  // Show list of habits
  return (
    <View style={[CommonStyles.container, styles.container]}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />
        }
      >
        {/* Header */}
        <View style={styles.header}>
          <Text style={styles.headerTitle}>My Habits</Text>
          <Text style={styles.headerSubtitle}>
            {habits.length} {habits.length === 1 ? 'habit' : 'habits'}
          </Text>
        </View>

        {/* Habit list */}
        {habits.map(habit => (
          <Habit
            key={habit.id}
            habit={habit}
            onEdit={() => onEditHabit(habit)}
            onDelete={() => onDeleteHabit(habit)}
            onUpdate={onRefresh}
          />
        ))}

        {/* Add some bottom padding */}
        <View style={styles.bottomPadding} />
      </ScrollView>

      {/* Floating add button */}
      <TouchableOpacity
        style={styles.floatingButton}
        onPress={onAddNew}
      >
        <Text style={styles.floatingButtonText}>+</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  
  centered: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  
  scrollContent: {
    flexGrow: 1,
    paddingBottom: 80,  // Make room for floating button
  },
  
  header: {
    padding: Spacing.md,
    paddingTop: Spacing.lg,
  },
  
  headerTitle: {
    fontSize: Typography.h2,
    fontWeight: Typography.bold,
    color: Colors.textPrimary,
    marginBottom: Spacing.xs,
  },
  
  headerSubtitle: {
    fontSize: Typography.body,
    color: Colors.textSecondary,
  },
  
  loadingText: {
    marginTop: Spacing.md,
    fontSize: Typography.body,
    color: Colors.textSecondary,
  },
  
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: Spacing.xl,
  },
  
  emptyEmoji: {
    fontSize: 64,
    marginBottom: Spacing.lg,
  },
  
  emptyTitle: {
    fontSize: Typography.h2,
    fontWeight: Typography.bold,
    color: Colors.textPrimary,
    marginBottom: Spacing.sm,
  },
  
  emptyText: {
    fontSize: Typography.body,
    color: Colors.textSecondary,
    textAlign: 'center',
    marginBottom: Spacing.xl,
  },
  
  addButton: {
    paddingHorizontal: Spacing.xl,
  },
  
  bottomPadding: {
    height: Spacing.xl,
  },
  
  floatingButton: {
    position: 'absolute',
    bottom: Spacing.lg,
    right: Spacing.lg,
    width: 60,
    height: 60,
    borderRadius: BorderRadius.full,
    backgroundColor: Colors.primary,
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
    elevation: 8,
  },
  
  floatingButtonText: {
    fontSize: 32,
    color: Colors.textPrimary,
    fontWeight: Typography.bold,
    lineHeight: 32,
  },
});
