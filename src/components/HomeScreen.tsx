/**
 * HomeScreen Component
 * Main screen showing habits list with add and logout functionality
 */

import React, {useState, useEffect} from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Alert,
  SafeAreaView,
} from 'react-native';
import {useAuth} from '../contexts/AuthContext';
import {
  Habit as HabitType,
  CreateHabitRequest,
  UpdateHabitRequest,
} from '../types';
import {HabitList} from './HabitList';
import {HabitDetailsModal} from './HabitDetailsModal';
import {habitApi, ApiError} from '../services/api';
import {colors, spacing, borderRadius, fontSize, shadow} from '../styles/theme';

/**
 * Home screen with habits list and management
 */
export function HomeScreen() {
  const {user, token, logout} = useAuth();

  const [habits, setHabits] = useState<HabitType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [habitModalVisible, setHabitModalVisible] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState<HabitType | null>(null);

  // Load habits when component mounts
  useEffect(() => {
    loadHabits();
  }, []);

  /**
   * Loads all habits for the current user
   */
  const loadHabits = async () => {
    if (!token) {
      return;
    }

    setIsLoading(true);
    try {
      const userHabits = await habitApi.getHabits(token);
      setHabits(userHabits);
    } catch (err) {
      console.error('Failed to load habits:', err);
      Alert.alert('Error', 'Failed to load habits. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Opens modal to create a new habit
   */
  const handleAddNewHabit = () => {
    setSelectedHabit(null);
    setHabitModalVisible(true);
  };

  /**
   * Opens modal to edit an existing habit
   */
  const handleEditHabit = (habit: HabitType) => {
    setSelectedHabit(habit);
    setHabitModalVisible(true);
  };

  /**
   * Saves a habit (create or update)
   */
  const handleSaveHabit = async (
    data: CreateHabitRequest | UpdateHabitRequest,
  ) => {
    if (!token) {
      return;
    }

    try {
      if (selectedHabit) {
        // Update existing habit
        await habitApi.updateHabit(
          selectedHabit.id,
          data as UpdateHabitRequest,
          token,
        );
      } else {
        // Create new habit
        await habitApi.createHabit(data as CreateHabitRequest, token);
      }

      // Reload habits
      await loadHabits();
    } catch (err) {
      throw err; // Let modal handle the error display
    }
  };

  /**
   * Deletes a habit with confirmation
   */
  const handleDeleteHabit = async (habitId: number) => {
    if (!token) {
      return;
    }

    try {
      await habitApi.deleteHabit(habitId, token);
      await loadHabits();
    } catch (err) {
      Alert.alert('Error', 'Failed to delete habit. Please try again.');
    }
  };

  /**
   * Logs out the current user
   */
  const handleLogout = () => {
    Alert.alert('Logout', 'Are you sure you want to logout?', [
      {text: 'Cancel', style: 'cancel'},
      {
        text: 'Logout',
        style: 'destructive',
        onPress: logout,
      },
    ]);
  };

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <View>
          <Text style={styles.greeting}>Hello, {user?.name}!</Text>
          <Text style={styles.subtitle}>Track your daily habits</Text>
        </View>
        <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
          <Text style={styles.logoutText}>Logout</Text>
        </TouchableOpacity>
      </View>

      {/* Habits List */}
      <HabitList
        habits={habits}
        token={token || ''}
        isLoading={isLoading}
        onEditHabit={handleEditHabit}
        onDeleteHabit={handleDeleteHabit}
        onRefresh={loadHabits}
        onAddNew={handleAddNewHabit}
      />

      {/* Floating Add Button */}
      {habits.length > 0 && (
        <TouchableOpacity
          style={styles.fab}
          onPress={handleAddNewHabit}
          accessibilityLabel="Add new habit">
          <Text style={styles.fabText}>+</Text>
        </TouchableOpacity>
      )}

      {/* Habit Details Modal */}
      <HabitDetailsModal
        visible={habitModalVisible}
        habit={selectedHabit}
        onClose={() => {
          setHabitModalVisible(false);
          setSelectedHabit(null);
        }}
        onSave={handleSaveHabit}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    backgroundColor: colors.surfaceLight,
    borderBottomWidth: 1,
    borderBottomColor: colors.borderLight,
  },
  greeting: {
    fontSize: fontSize.xl,
    fontWeight: '600',
    color: colors.text,
  },
  subtitle: {
    fontSize: fontSize.sm,
    color: colors.textLight,
    marginTop: spacing.xs,
  },
  logoutButton: {
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: borderRadius.sm,
    borderWidth: 1,
    borderColor: colors.border,
  },
  logoutText: {
    fontSize: fontSize.sm,
    color: colors.textLight,
  },
  fab: {
    position: 'absolute',
    right: spacing.lg,
    bottom: spacing.lg,
    width: 56,
    height: 56,
    borderRadius: borderRadius.round,
    backgroundColor: colors.accent,
    justifyContent: 'center',
    alignItems: 'center',
    ...shadow.lg,
  },
  fabText: {
    fontSize: 28,
    color: colors.surfaceLight,
    fontWeight: 'bold',
  },
});
