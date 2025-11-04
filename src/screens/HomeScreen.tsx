import React, {useState, useEffect, useCallback} from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  FlatList,
  Alert,
  RefreshControl,
  ActivityIndicator,
} from 'react-native';
import {Habit as HabitType, getHabits, logout, ApiError} from '../services/api';
import {Habit} from '../components/Habit';
import {HabitDetailsModal} from '../components/HabitDetailsModal';
import {colors, spacing, borderRadius, typography, shadows} from '../styles/theme';

interface HomeScreenProps {
  onLogout: () => void;
}

/**
 * Home screen component
 * Displays the list of user's habits and allows creating new ones
 */
export const HomeScreen: React.FC<HomeScreenProps> = ({onLogout}) => {
  const [habits, setHabits] = useState<HabitType[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [showNewHabitModal, setShowNewHabitModal] = useState(false);

  /**
   * Load habits from API
   */
  const loadHabits = async (isRefreshing = false) => {
    try {
      if (!isRefreshing) {
        setLoading(true);
      }
      const userHabits = await getHabits();
      setHabits(userHabits);
    } catch (error) {
      if (error instanceof ApiError) {
        // If unauthorized, logout
        if (error.status === 401) {
          Alert.alert('Session Expired', 'Please login again.');
          handleLogout();
        } else {
          Alert.alert('Error', `Failed to load habits: ${error.message}`);
        }
      }
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  // Load habits when component mounts
  useEffect(() => {
    loadHabits();
  }, []);

  /**
   * Handle pull-to-refresh
   */
  const onRefresh = useCallback(() => {
    setRefreshing(true);
    loadHabits(true);
  }, []);

  /**
   * Handle logout
   */
  const handleLogout = async () => {
    try {
      await logout();
      onLogout();
    } catch (error) {
      Alert.alert('Error', 'Failed to logout. Please try again.');
    }
  };

  /**
   * Confirm logout with user
   */
  const confirmLogout = () => {
    Alert.alert(
      'Logout',
      'Are you sure you want to logout?',
      [
        {text: 'Cancel', style: 'cancel'},
        {text: 'Logout', style: 'destructive', onPress: handleLogout},
      ],
    );
  };

  /**
   * Render empty state when no habits exist
   */
  const renderEmptyState = () => (
    <View style={styles.emptyState}>
      <Text style={styles.emptyStateIcon}>🎯</Text>
      <Text style={styles.emptyStateTitle}>No Habits Yet</Text>
      <Text style={styles.emptyStateText}>
        Start building better habits by creating your first one!
      </Text>
    </View>
  );

  /**
   * Render a single habit item
   */
  const renderHabit = ({item}: {item: HabitType}) => (
    <Habit habit={item} onUpdate={() => loadHabits()} />
  );

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={colors.primary} />
        <Text style={styles.loadingText}>Loading your habits...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <View>
          <Text style={styles.headerTitle}>My Habits</Text>
          <Text style={styles.headerSubtitle}>
            {habits.length} {habits.length === 1 ? 'habit' : 'habits'} tracked
          </Text>
        </View>
        <TouchableOpacity style={styles.logoutButton} onPress={confirmLogout}>
          <Text style={styles.logoutButtonText}>Logout</Text>
        </TouchableOpacity>
      </View>

      {/* Habit List */}
      <FlatList
        data={habits}
        renderItem={renderHabit}
        keyExtractor={item => item.id.toString()}
        contentContainerStyle={[
          styles.listContent,
          habits.length === 0 && styles.listContentEmpty,
        ]}
        ListEmptyComponent={renderEmptyState}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            tintColor={colors.primary}
            colors={[colors.primary]}
          />
        }
      />

      {/* Add New Habit Button */}
      <TouchableOpacity
        style={styles.fab}
        onPress={() => setShowNewHabitModal(true)}>
        <Text style={styles.fabText}>+</Text>
      </TouchableOpacity>

      {/* New Habit Modal */}
      <HabitDetailsModal
        visible={showNewHabitModal}
        onClose={() => setShowNewHabitModal(false)}
        onSave={() => loadHabits()}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.background,
  },
  loadingText: {
    marginTop: spacing.md,
    fontSize: typography.base,
    color: colors.textSecondary,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: spacing.lg,
    backgroundColor: colors.surfaceLight,
    borderBottomWidth: 1,
    borderBottomColor: colors.borderLight,
    ...shadows.sm,
  },
  headerTitle: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.text,
  },
  headerSubtitle: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  logoutButton: {
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    backgroundColor: colors.surface,
    borderRadius: borderRadius.md,
    borderWidth: 1,
    borderColor: colors.border,
  },
  logoutButtonText: {
    fontSize: typography.sm,
    fontWeight: typography.medium,
    color: colors.text,
  },
  listContent: {
    padding: spacing.lg,
  },
  listContentEmpty: {
    flexGrow: 1,
  },
  emptyState: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: spacing.xl,
  },
  emptyStateIcon: {
    fontSize: 64,
    marginBottom: spacing.lg,
  },
  emptyStateTitle: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.text,
    marginBottom: spacing.sm,
    textAlign: 'center',
  },
  emptyStateText: {
    fontSize: typography.base,
    color: colors.textSecondary,
    textAlign: 'center',
    lineHeight: 24,
  },
  fab: {
    position: 'absolute',
    right: spacing.lg,
    bottom: spacing.lg,
    width: 60,
    height: 60,
    borderRadius: borderRadius.full,
    backgroundColor: colors.secondary,
    justifyContent: 'center',
    alignItems: 'center',
    ...shadows.lg,
  },
  fabText: {
    fontSize: 32,
    color: colors.text,
    fontWeight: typography.bold,
  },
});
