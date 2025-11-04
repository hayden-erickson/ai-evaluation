/**
 * Home screen displaying the list of habits
 */

import React, {useState, useEffect, useCallback} from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  RefreshControl,
  Alert,
  TouchableOpacity,
} from 'react-native';
import {SafeAreaView} from 'react-native-safe-area-context';
import {useAuth} from '../services/AuthContext';
import {HabitComponent} from '../components/Habit';
import {HabitDetailsModal} from '../components/HabitDetailsModal';
import {Button} from '../components/Button';
import {theme} from '../styles/theme';
import {Habit} from '../types';
import * as api from '../services/api';

/**
 * Home screen with habit list and management
 */
export const HomeScreen: React.FC = () => {
  const {user, logout} = useAuth();
  const [habits, setHabits] = useState<Habit[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [showHabitModal, setShowHabitModal] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState<Habit | null>(null);

  // Load habits on mount
  useEffect(() => {
    loadHabits();
  }, []);

  /**
   * Load all habits for the user
   */
  const loadHabits = async () => {
    try {
      setIsLoading(true);
      const userHabits = await api.getHabits();
      setHabits(userHabits);
    } catch (error) {
      console.error('Error loading habits:', error);
      Alert.alert('Error', 'Failed to load habits');
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Refresh habits
   */
  const handleRefresh = useCallback(async () => {
    setIsRefreshing(true);
    await loadHabits();
    setIsRefreshing(false);
  }, []);

  /**
   * Handle creating a new habit
   */
  const handleCreateHabit = () => {
    setSelectedHabit(null);
    setShowHabitModal(true);
  };

  /**
   * Handle editing a habit
   */
  const handleEditHabit = (habit: Habit) => {
    setSelectedHabit(habit);
    setShowHabitModal(true);
  };

  /**
   * Handle deleting a habit
   */
  const handleDeleteHabit = async (habitId: number) => {
    try {
      await api.deleteHabit(habitId);
      await loadHabits();
    } catch (error) {
      console.error('Error deleting habit:', error);
      Alert.alert('Error', 'Failed to delete habit');
    }
  };

  /**
   * Handle saving habit (create or update)
   */
  const handleSaveHabit = async (data: {name: string; description: string}) => {
    try {
      if (selectedHabit) {
        // Update existing habit
        await api.updateHabit(selectedHabit.id, data);
      } else {
        // Create new habit
        await api.createHabit(data);
      }
      await loadHabits();
      setShowHabitModal(false);
    } catch (error) {
      console.error('Error saving habit:', error);
      Alert.alert('Error', 'Failed to save habit');
      throw error;
    }
  };

  /**
   * Handle logout confirmation
   */
  const handleLogout = () => {
    Alert.alert('Logout', 'Are you sure you want to logout?', [
      {text: 'Cancel', style: 'cancel'},
      {text: 'Logout', style: 'destructive', onPress: logout},
    ]);
  };

  /**
   * Render empty state
   */
  const renderEmpty = () => (
    <View style={styles.emptyContainer}>
      <Text style={styles.emptyTitle}>No Habits Yet</Text>
      <Text style={styles.emptyText}>
        Start building better habits by adding your first one!
      </Text>
    </View>
  );

  /**
   * Render habit item
   */
  const renderHabit = ({item}: {item: Habit}) => (
    <HabitComponent
      habit={item}
      onEdit={() => handleEditHabit(item)}
      onDelete={() => handleDeleteHabit(item.id)}
      onLogCreated={loadHabits}
    />
  );

  return (
    <SafeAreaView style={styles.container} edges={['top', 'left', 'right']}>
      {/* Header */}
      <View style={styles.header}>
        <View>
          <Text style={styles.greeting}>Hello, {user?.name || 'there'}!</Text>
          <Text style={styles.subtitle}>Keep up your streaks 🔥</Text>
        </View>
        <TouchableOpacity onPress={handleLogout} style={styles.logoutButton}>
          <Text style={styles.logoutText}>Logout</Text>
        </TouchableOpacity>
      </View>

      {/* Habit list */}
      <FlatList
        data={habits}
        renderItem={renderHabit}
        keyExtractor={item => item.id.toString()}
        contentContainerStyle={[
          styles.listContent,
          habits.length === 0 && styles.listContentEmpty,
        ]}
        ListEmptyComponent={!isLoading ? renderEmpty : null}
        refreshControl={
          <RefreshControl
            refreshing={isRefreshing}
            onRefresh={handleRefresh}
            tintColor={theme.colors.primary}
          />
        }
      />

      {/* Add new habit button */}
      <View style={styles.addButtonContainer}>
        <Button
          title="+ Add New Habit"
          onPress={handleCreateHabit}
          style={styles.addButton}
        />
      </View>

      {/* Habit details modal */}
      <HabitDetailsModal
        visible={showHabitModal}
        habit={selectedHabit}
        onClose={() => setShowHabitModal(false)}
        onSave={handleSaveHabit}
      />
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: theme.colors.background,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: theme.spacing.lg,
    backgroundColor: theme.colors.surface,
    ...theme.shadows.small,
  },
  greeting: {
    fontSize: theme.fontSize.xl,
    fontWeight: 'bold',
    color: theme.colors.text,
  },
  subtitle: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textLight,
    marginTop: theme.spacing.xs,
  },
  logoutButton: {
    paddingVertical: theme.spacing.sm,
    paddingHorizontal: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    backgroundColor: theme.colors.surfaceLight,
  },
  logoutText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textLight,
    fontWeight: '600',
  },
  listContent: {
    padding: theme.spacing.lg,
    paddingBottom: 100,
  },
  listContentEmpty: {
    flexGrow: 1,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: theme.spacing.xl,
  },
  emptyTitle: {
    fontSize: theme.fontSize.xl,
    fontWeight: 'bold',
    color: theme.colors.text,
    marginBottom: theme.spacing.sm,
  },
  emptyText: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textLight,
    textAlign: 'center',
  },
  addButtonContainer: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    padding: theme.spacing.lg,
    backgroundColor: theme.colors.background,
    ...theme.shadows.large,
  },
  addButton: {
    backgroundColor: theme.colors.accent,
  },
});
