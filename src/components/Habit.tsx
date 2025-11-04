import React, {useState, useEffect} from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Alert,
  ScrollView,
  ActivityIndicator,
} from 'react-native';
import {
  Habit as HabitType,
  Log,
  getHabitLogs,
  deleteHabit,
  deleteLog,
  ApiError,
} from '../services/api';
import {
  calculateStreak,
  getLast30Days,
  hasLoggedToday,
  formatDisplayDate,
} from '../utils/streakCalculator';
import {colors, spacing, borderRadius, typography, shadows} from '../styles/theme';
import {HabitDetailsModal} from './HabitDetailsModal';
import {LogDetailsModal} from './LogDetailsModal';

interface HabitProps {
  habit: HabitType;
  onUpdate: () => void;
}

/**
 * Habit component displays a single habit with its streak and logs
 * Includes buttons for editing, deleting, and logging
 */
export const Habit: React.FC<HabitProps> = ({habit, onUpdate}) => {
  const [logs, setLogs] = useState<Log[]>([]);
  const [loading, setLoading] = useState(true);
  const [showHabitModal, setShowHabitModal] = useState(false);
  const [showLogModal, setShowLogModal] = useState(false);
  const [selectedLog, setSelectedLog] = useState<Log | undefined>();
  const [expanded, setExpanded] = useState(false);

  const streak = calculateStreak(logs);
  const loggedToday = hasLoggedToday(logs);
  const last30Days = getLast30Days(logs);

  /**
   * Load habit logs from API
   */
  const loadLogs = async () => {
    try {
      setLoading(true);
      const habitLogs = await getHabitLogs(habit.id);
      setLogs(habitLogs);
    } catch (error) {
      if (error instanceof ApiError) {
        Alert.alert('Error', `Failed to load logs: ${error.message}`);
      }
    } finally {
      setLoading(false);
    }
  };

  // Load logs when component mounts
  useEffect(() => {
    loadLogs();
  }, [habit.id]);

  /**
   * Handle habit deletion
   * Confirms with user before deleting
   */
  const handleDelete = () => {
    Alert.alert(
      'Delete Habit',
      `Are you sure you want to delete "${habit.name}"? This action cannot be undone.`,
      [
        {text: 'Cancel', style: 'cancel'},
        {
          text: 'Delete',
          style: 'destructive',
          onPress: async () => {
            try {
              await deleteHabit(habit.id);
              onUpdate();
            } catch (error) {
              if (error instanceof ApiError) {
                Alert.alert('Error', `Failed to delete habit: ${error.message}`);
              }
            }
          },
        },
      ],
    );
  };

  /**
   * Handle log deletion
   * Confirms with user before deleting
   */
  const handleDeleteLog = (log: Log) => {
    Alert.alert(
      'Delete Log',
      'Are you sure you want to delete this log entry?',
      [
        {text: 'Cancel', style: 'cancel'},
        {
          text: 'Delete',
          style: 'destructive',
          onPress: async () => {
            try {
              await deleteLog(log.id);
              await loadLogs();
              onUpdate();
            } catch (error) {
              if (error instanceof ApiError) {
                Alert.alert('Error', `Failed to delete log: ${error.message}`);
              }
            }
          },
        },
      ],
    );
  };

  /**
   * Open modal to create a new log
   */
  const handleNewLog = () => {
    setSelectedLog(undefined);
    setShowLogModal(true);
  };

  /**
   * Open modal to edit an existing log
   */
  const handleEditLog = (log: Log) => {
    setSelectedLog(log);
    setShowLogModal(true);
  };

  /**
   * Handle successful save from modals
   */
  const handleModalSave = async () => {
    await loadLogs();
    onUpdate();
  };

  return (
    <View style={styles.container}>
      {/* Header Section */}
      <View style={styles.header}>
        <View style={styles.headerLeft}>
          <Text style={styles.name}>{habit.name}</Text>
          {habit.description ? (
            <Text style={styles.description}>{habit.description}</Text>
          ) : null}
        </View>

        <View style={styles.headerRight}>
          <TouchableOpacity
            style={styles.iconButton}
            onPress={() => setShowHabitModal(true)}>
            <Text style={styles.iconText}>✏️</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.iconButton} onPress={handleDelete}>
            <Text style={styles.iconText}>🗑️</Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Streak Section */}
      <View style={styles.streakSection}>
        <View style={styles.streakBadge}>
          <Text style={styles.streakNumber}>{streak}</Text>
          <Text style={styles.streakLabel}>day streak</Text>
        </View>

        <TouchableOpacity
          style={[
            styles.logButton,
            loggedToday && styles.logButtonDisabled,
          ]}
          onPress={handleNewLog}
          disabled={loggedToday}>
          <Text style={styles.logButtonText}>
            {loggedToday ? '✓ Logged Today' : '+ Log Today'}
          </Text>
        </TouchableOpacity>
      </View>

      {/* Toggle Streak List */}
      <TouchableOpacity
        style={styles.toggleButton}
        onPress={() => setExpanded(!expanded)}>
        <Text style={styles.toggleText}>
          {expanded ? 'Hide History' : 'Show History'}
        </Text>
        <Text style={styles.toggleIcon}>{expanded ? '▲' : '▼'}</Text>
      </TouchableOpacity>

      {/* Streak List (Last 30 Days) */}
      {expanded && (
        <View style={styles.streakList}>
          {loading ? (
            <ActivityIndicator color={colors.primary} />
          ) : (
            <ScrollView
              horizontal
              showsHorizontalScrollIndicator={false}
              contentContainerStyle={styles.streakListContent}>
              {last30Days.map((day, index) => (
                <TouchableOpacity
                  key={day.dateStr}
                  style={[
                    styles.dayContainer,
                    day.hasLog && styles.dayContainerActive,
                  ]}
                  onPress={() => day.log && handleEditLog(day.log)}
                  disabled={!day.log}>
                  <Text style={styles.dayLabel}>
                    {formatDisplayDate(day.date)}
                  </Text>
                  {day.hasLog ? (
                    <View style={styles.logIndicator}>
                      <Text style={styles.logIndicatorText}>✓</Text>
                    </View>
                  ) : (
                    <View style={styles.emptyIndicator} />
                  )}
                </TouchableOpacity>
              ))}
            </ScrollView>
          )}
        </View>
      )}

      {/* Modals */}
      <HabitDetailsModal
        visible={showHabitModal}
        habit={habit}
        onClose={() => setShowHabitModal(false)}
        onSave={handleModalSave}
      />

      <LogDetailsModal
        visible={showLogModal}
        habitId={habit.id}
        log={selectedLog}
        onClose={() => {
          setShowLogModal(false);
          setSelectedLog(undefined);
        }}
        onSave={handleModalSave}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.surfaceLight,
    borderRadius: borderRadius.lg,
    padding: spacing.md,
    marginBottom: spacing.md,
    ...shadows.md,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: spacing.md,
  },
  headerLeft: {
    flex: 1,
    marginRight: spacing.md,
  },
  headerRight: {
    flexDirection: 'row',
    gap: spacing.sm,
  },
  name: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.text,
    marginBottom: spacing.xs,
  },
  description: {
    fontSize: typography.sm,
    color: colors.textSecondary,
  },
  iconButton: {
    padding: spacing.sm,
  },
  iconText: {
    fontSize: typography.lg,
  },
  streakSection: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: spacing.md,
  },
  streakBadge: {
    backgroundColor: colors.streakActive,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    alignItems: 'center',
    minWidth: 80,
  },
  streakNumber: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.text,
  },
  streakLabel: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  logButton: {
    flex: 1,
    backgroundColor: colors.primary,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    alignItems: 'center',
    marginLeft: spacing.md,
    ...shadows.sm,
  },
  logButtonDisabled: {
    backgroundColor: colors.success,
  },
  logButtonText: {
    fontSize: typography.base,
    fontWeight: typography.semibold,
    color: colors.text,
  },
  toggleButton: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: spacing.sm,
    gap: spacing.sm,
  },
  toggleText: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    fontWeight: typography.medium,
  },
  toggleIcon: {
    fontSize: typography.sm,
    color: colors.textSecondary,
  },
  streakList: {
    marginTop: spacing.sm,
    paddingTop: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.borderLight,
  },
  streakListContent: {
    gap: spacing.sm,
    paddingVertical: spacing.sm,
  },
  dayContainer: {
    width: 60,
    alignItems: 'center',
    padding: spacing.sm,
    backgroundColor: colors.surface,
    borderRadius: borderRadius.sm,
    borderWidth: 1,
    borderColor: colors.borderLight,
  },
  dayContainerActive: {
    backgroundColor: colors.streakActive,
    borderColor: colors.success,
  },
  dayLabel: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    marginBottom: spacing.xs,
    textAlign: 'center',
  },
  logIndicator: {
    width: 32,
    height: 32,
    borderRadius: borderRadius.full,
    backgroundColor: colors.success,
    justifyContent: 'center',
    alignItems: 'center',
  },
  logIndicatorText: {
    fontSize: typography.base,
    color: colors.text,
  },
  emptyIndicator: {
    width: 32,
    height: 32,
    borderRadius: borderRadius.full,
    backgroundColor: colors.streakInactive,
    borderWidth: 1,
    borderColor: colors.border,
  },
});
