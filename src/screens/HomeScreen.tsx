/**
 * Home screen displaying habit list
 */

import React, { useState, useEffect, useCallback } from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  SafeAreaView,
  ScrollView,
  RefreshControl,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { useAuth } from '../contexts/AuthContext';
import {
  Habit as HabitModel,
  Log,
  CreateHabitRequest,
  UpdateHabitRequest,
  CreateLogRequest,
  UpdateLogRequest,
} from '../types/models';
import { habitAPI, logAPI } from '../services/api';
import { theme } from '../styles/theme';
import { commonStyles } from '../styles/commonStyles';
import { Habit } from '../components/Habit';
import { HabitDetailsModal } from '../components/HabitDetailsModal';

interface HabitWithLogs {
  habit: HabitModel;
  logs: Log[];
}

/**
 * Home screen showing list of habits
 */
export const HomeScreen: React.FC = () => {
  const { user, logout } = useAuth();
  const [habitsWithLogs, setHabitsWithLogs] = useState<HabitWithLogs[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [isHabitModalVisible, setIsHabitModalVisible] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState<HabitModel | null>(null);
  const [error, setError] = useState<string | null>(null);

  /**
   * Fetch habits and their logs
   */
  const fetchHabitsAndLogs = async (showLoading = true) => {
    try {
      if (showLoading) {
        setIsLoading(true);
      }
      setError(null);

      // Fetch all habits
      const habits = await habitAPI.getUserHabits();

      // Fetch logs for each habit
      const habitsWithLogsData = await Promise.all(
        habits.map(async (habit) => {
          try {
            const logs = await logAPI.getHabitLogs(habit.id);
            return { habit, logs };
          } catch (err) {
            console.error(`Failed to fetch logs for habit ${habit.id}:`, err);
            return { habit, logs: [] };
          }
        })
      );

      setHabitsWithLogs(habitsWithLogsData);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch habits';
      setError(errorMessage);
      Alert.alert('Error', errorMessage);
    } finally {
      setIsLoading(false);
      setIsRefreshing(false);
    }
  };

  /**
   * Load habits on mount
   */
  useEffect(() => {
    fetchHabitsAndLogs();
  }, []);

  /**
   * Handle refresh
   */
  const handleRefresh = useCallback(() => {
    setIsRefreshing(true);
    fetchHabitsAndLogs(false);
  }, []);

  /**
   * Handle create new habit button
   */
  const handleNewHabitPress = () => {
    setSelectedHabit(null);
    setIsHabitModalVisible(true);
  };

  /**
   * Handle edit habit
   */
  const handleEditHabit = (habit: HabitModel) => {
    setSelectedHabit(habit);
    setIsHabitModalVisible(true);
  };

  /**
   * Handle save habit (create or update)
   */
  const handleSaveHabit = async (data: CreateHabitRequest | UpdateHabitRequest) => {
    try {
      if (selectedHabit) {
        // Update existing habit
        await habitAPI.updateHabit(selectedHabit.id, data as UpdateHabitRequest);
      } else {
        // Create new habit
        await habitAPI.createHabit(data as CreateHabitRequest);
      }
      await fetchHabitsAndLogs(false);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to save habit';
      throw new Error(errorMessage);
    }
  };

  /**
   * Handle delete habit
   */
  const handleDeleteHabit = async (habitId: number) => {
    try {
      await habitAPI.deleteHabit(habitId);
      await fetchHabitsAndLogs(false);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete habit';
      Alert.alert('Error', errorMessage);
    }
  };

  /**
   * Handle create log
   */
  const handleCreateLog = async (habitId: number, data: CreateLogRequest) => {
    try {
      await logAPI.createLog(habitId, data);
      await fetchHabitsAndLogs(false);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create log';
      throw new Error(errorMessage);
    }
  };

  /**
   * Handle update log
   */
  const handleUpdateLog = async (logId: number, data: UpdateLogRequest) => {
    try {
      await logAPI.updateLog(logId, data);
      await fetchHabitsAndLogs(false);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update log';
      throw new Error(errorMessage);
    }
  };

  /**
   * Handle logout
   */
  const handleLogout = () => {
    Alert.alert('Logout', 'Are you sure you want to logout?', [
      { text: 'Cancel', style: 'cancel' },
      { text: 'Logout', style: 'destructive', onPress: logout },
    ]);
  };

  if (isLoading) {
    return (
      <SafeAreaView style={commonStyles.safeArea}>
        <View style={commonStyles.centered}>
          <ActivityIndicator size="large" color={theme.colors.primary} />
          <Text style={styles.loadingText}>Loading your habits...</Text>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={commonStyles.safeArea}>
      <View style={styles.container}>
        {/* Header */}
        <View style={styles.header}>
          <View>
            <Text style={styles.greeting}>Hello, {user?.name}! 👋</Text>
            <Text style={styles.subtitle}>Keep up your streaks!</Text>
          </View>
          <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
            <Text style={styles.logoutButtonText}>Logout</Text>
          </TouchableOpacity>
        </View>

        {/* Habit List */}
        <ScrollView
          style={styles.scrollView}
          contentContainerStyle={styles.scrollContent}
          refreshControl={
            <RefreshControl
              refreshing={isRefreshing}
              onRefresh={handleRefresh}
              tintColor={theme.colors.primary}
            />
          }
        >
          {habitsWithLogs.length === 0 ? (
            <View style={styles.emptyState}>
              <Text style={styles.emptyStateEmoji}>🎯</Text>
              <Text style={styles.emptyStateTitle}>No habits yet</Text>
              <Text style={styles.emptyStateText}>
                Start building your streak by creating your first habit!
              </Text>
            </View>
          ) : (
            habitsWithLogs.map((item) => (
              <Habit
                key={item.habit.id}
                habit={item.habit}
                logs={item.logs}
                onEdit={() => handleEditHabit(item.habit)}
                onDelete={() => handleDeleteHabit(item.habit.id)}
                onCreateLog={(data) => handleCreateLog(item.habit.id, data)}
                onUpdateLog={handleUpdateLog}
              />
            ))
          )}
        </ScrollView>

        {/* Add New Habit Button */}
        <TouchableOpacity style={styles.fab} onPress={handleNewHabitPress}>
          <Text style={styles.fabText}>+</Text>
        </TouchableOpacity>

        {/* Habit Details Modal */}
        <HabitDetailsModal
          visible={isHabitModalVisible}
          habit={selectedHabit}
          onClose={() => {
            setIsHabitModalVisible(false);
            setSelectedHabit(null);
          }}
          onSave={handleSaveHabit}
        />
      </View>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },

  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: theme.spacing.lg,
    backgroundColor: theme.colors.surface,
    borderBottomWidth: 1,
    borderBottomColor: theme.colors.border,
  },

  greeting: {
    fontSize: theme.fontSize.xl,
    fontWeight: 'bold',
    color: theme.colors.text,
  },

  subtitle: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textSecondary,
    marginTop: 2,
  },

  logoutButton: {
    backgroundColor: theme.colors.border,
    borderRadius: theme.borderRadius.sm,
    paddingVertical: theme.spacing.xs,
    paddingHorizontal: theme.spacing.md,
  },

  logoutButtonText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.text,
    fontWeight: '600',
  },

  scrollView: {
    flex: 1,
  },

  scrollContent: {
    padding: theme.spacing.lg,
  },

  loadingText: {
    marginTop: theme.spacing.md,
    fontSize: theme.fontSize.md,
    color: theme.colors.textSecondary,
  },

  emptyState: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: theme.spacing.xxl,
  },

  emptyStateEmoji: {
    fontSize: 64,
    marginBottom: theme.spacing.md,
  },

  emptyStateTitle: {
    fontSize: theme.fontSize.xl,
    fontWeight: 'bold',
    color: theme.colors.text,
    marginBottom: theme.spacing.sm,
  },

  emptyStateText: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textSecondary,
    textAlign: 'center',
    paddingHorizontal: theme.spacing.xl,
  },

  fab: {
    position: 'absolute',
    right: theme.spacing.lg,
    bottom: theme.spacing.lg,
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: theme.colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
    ...theme.shadows.lg,
  },

  fabText: {
    fontSize: 32,
    color: theme.colors.surface,
    fontWeight: '300',
  },
});

