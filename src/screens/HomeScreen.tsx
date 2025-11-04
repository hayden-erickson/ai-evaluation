/**
 * HomeScreen - Main screen showing user's habits
 * Manages habit list and user logout
 */

import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Alert,
} from 'react-native';
import { Habit, CreateHabitRequest, UpdateHabitRequest, User } from '../types';
import { apiService } from '../services/api';
import { Colors, Spacing, Typography, CommonStyles } from '../styles/theme';
import HabitList from '../components/HabitList';
import HabitDetailsModal from '../components/HabitDetailsModal';

interface HomeScreenProps {
  user: User;
  onLogout: () => void;
}

export default function HomeScreen({ user, onLogout }: HomeScreenProps) {
  const [habits, setHabits] = useState<Habit[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showHabitModal, setShowHabitModal] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState<Habit | null>(null);

  // Load habits on component mount
  useEffect(() => {
    loadHabits();
  }, []);

  /**
   * Load all habits from the API
   */
  const loadHabits = async () => {
    try {
      setIsLoading(true);
      const fetchedHabits = await apiService.getHabits();
      setHabits(fetchedHabits);
    } catch (error) {
      Alert.alert('Error', 'Failed to load habits. Please try again.');
      console.error('Failed to load habits:', error);
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Handle opening modal to add a new habit
   */
  const handleAddNewHabit = () => {
    setSelectedHabit(null);
    setShowHabitModal(true);
  };

  /**
   * Handle opening modal to edit an existing habit
   */
  const handleEditHabit = (habit: Habit) => {
    setSelectedHabit(habit);
    setShowHabitModal(true);
  };

  /**
   * Handle saving a habit (create or update)
   */
  const handleSaveHabit = async (data: CreateHabitRequest | UpdateHabitRequest) => {
    try {
      if (selectedHabit) {
        // Update existing habit
        await apiService.updateHabit(selectedHabit.id, data as UpdateHabitRequest);
      } else {
        // Create new habit
        await apiService.createHabit(data as CreateHabitRequest);
      }
      
      // Reload habits to reflect changes
      await loadHabits();
    } catch (error) {
      throw new Error(error instanceof Error ? error.message : 'Failed to save habit');
    }
  };

  /**
   * Handle deleting a habit
   */
  const handleDeleteHabit = (habit: Habit) => {
    Alert.alert(
      'Delete Habit',
      `Are you sure you want to delete "${habit.name}"? This will also delete all associated logs.`,
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Delete',
          style: 'destructive',
          onPress: async () => {
            try {
              await apiService.deleteHabit(habit.id);
              await loadHabits();
            } catch (error) {
              Alert.alert('Error', 'Failed to delete habit');
            }
          },
        },
      ]
    );
  };

  /**
   * Handle user logout
   */
  const handleLogout = () => {
    Alert.alert(
      'Logout',
      'Are you sure you want to logout?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Logout',
          onPress: () => {
            apiService.logout();
            onLogout();
          },
        },
      ]
    );
  };

  return (
    <View style={CommonStyles.container}>
      {/* Header with user info and logout */}
      <View style={styles.headerBar}>
        <View style={styles.headerLeft}>
          <Text style={styles.greeting}>Hello, {user.name}! 👋</Text>
        </View>
        <TouchableOpacity onPress={handleLogout} style={styles.logoutButton}>
          <Text style={styles.logoutText}>Logout</Text>
        </TouchableOpacity>
      </View>

      {/* Habit list */}
      <HabitList
        habits={habits}
        isLoading={isLoading}
        onRefresh={loadHabits}
        onAddNew={handleAddNewHabit}
        onEditHabit={handleEditHabit}
        onDeleteHabit={handleDeleteHabit}
      />

      {/* Habit details modal */}
      <HabitDetailsModal
        visible={showHabitModal}
        onClose={() => {
          setShowHabitModal(false);
          setSelectedHabit(null);
        }}
        onSave={handleSaveHabit}
        habit={selectedHabit}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  headerBar: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: Spacing.md,
    backgroundColor: Colors.cardBackground,
    borderBottomWidth: 1,
    borderBottomColor: Colors.border,
  },
  
  headerLeft: {
    flex: 1,
  },
  
  greeting: {
    fontSize: Typography.h4,
    fontWeight: Typography.semibold,
    color: Colors.textPrimary,
  },
  
  logoutButton: {
    paddingHorizontal: Spacing.md,
    paddingVertical: Spacing.sm,
  },
  
  logoutText: {
    fontSize: Typography.body,
    color: Colors.accent,
    fontWeight: Typography.medium,
  },
});
