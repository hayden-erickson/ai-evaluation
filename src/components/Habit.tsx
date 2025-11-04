/**
 * Habit component displaying habit details, streak, and logs
 */

import React, {useState, useEffect, useCallback} from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Alert,
  ScrollView,
} from 'react-native';
import {Habit as HabitType, Log} from '../types';
import {theme} from '../styles/theme';
import {
  calculateStreak,
  getLastNDays,
  hasLogForDate,
  getLogsForDate,
  formatDayShort,
  formatDateShort,
  isToday,
} from '../utils/streakUtils';
import * as api from '../services/api';
import {LogDetailsModal} from './LogDetailsModal';

interface HabitComponentProps {
  habit: HabitType;
  onEdit: () => void;
  onDelete: () => void;
  onLogCreated: () => void;
}

/**
 * Habit component with streak tracking and log management
 */
export const HabitComponent: React.FC<HabitComponentProps> = ({
  habit,
  onEdit,
  onDelete,
  onLogCreated,
}) => {
  const [logs, setLogs] = useState<Log[]>([]);
  const [showLogs, setShowLogs] = useState(false);
  const [selectedLog, setSelectedLog] = useState<Log | null>(null);
  const [showLogModal, setShowLogModal] = useState(false);
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);

  /**
   * Load logs for the habit
   */
  const loadLogs = useCallback(async () => {
    try {
      const habitLogs = await api.getHabitLogs(habit.id);
      setLogs(habitLogs);
    } catch (error) {
      console.error('Error loading logs:', error);
      Alert.alert('Error', 'Failed to load habit logs');
    }
  }, [habit.id]);

  // Load logs when component mounts
  useEffect(() => {
    loadLogs();
  }, [loadLogs]);

  /**
   * Create a new log for today
   */
  const handleCreateLog = async () => {
    const today = new Date();
    const todayLogs = getLogsForDate(logs, today);

    if (todayLogs.length > 0) {
      // Edit existing log
      setSelectedLog(todayLogs[0]);
      setSelectedDate(today);
      setShowLogModal(true);
    } else {
      // Create new log
      setSelectedLog(null);
      setSelectedDate(today);
      setShowLogModal(true);
    }
  };

  /**
   * Handle log save
   */
  const handleSaveLog = async (data: {notes: string}) => {
    try {
      if (selectedLog) {
        // Update existing log
        await api.updateLog(selectedLog.id, data);
      } else {
        // Create new log
        await api.createLog(habit.id, data);
      }
      await loadLogs();
      onLogCreated();
      setShowLogModal(false);
    } catch (error) {
      console.error('Error saving log:', error);
      Alert.alert('Error', 'Failed to save log entry');
    }
  };

  /**
   * Handle log click
   */
  const handleLogClick = (date: Date) => {
    const dateLogs = getLogsForDate(logs, date);
    if (dateLogs.length > 0) {
      setSelectedLog(dateLogs[0]);
      setSelectedDate(date);
      setShowLogModal(true);
    }
  };

  /**
   * Handle delete confirmation
   */
  const handleDeleteConfirm = () => {
    Alert.alert(
      'Delete Habit',
      `Are you sure you want to delete "${habit.name}"? This will also delete all associated logs.`,
      [
        {text: 'Cancel', style: 'cancel'},
        {text: 'Delete', style: 'destructive', onPress: onDelete},
      ],
    );
  };

  const streak = calculateStreak(logs);
  const days = getLastNDays(14); // Show last 14 days

  return (
    <View style={styles.container}>
      {/* Habit header */}
      <View style={styles.header}>
        <View style={styles.headerLeft}>
          <Text style={styles.name}>{habit.name}</Text>
          {habit.description && (
            <Text style={styles.description}>{habit.description}</Text>
          )}
        </View>
        <View style={styles.actions}>
          <TouchableOpacity
            onPress={onEdit}
            style={styles.actionButton}
            hitSlop={{top: 10, bottom: 10, left: 10, right: 10}}>
            <Text style={styles.actionButtonText}>✏️</Text>
          </TouchableOpacity>
          <TouchableOpacity
            onPress={handleDeleteConfirm}
            style={styles.actionButton}
            hitSlop={{top: 10, bottom: 10, left: 10, right: 10}}>
            <Text style={styles.actionButtonText}>🗑️</Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Streak count */}
      <View style={styles.streakContainer}>
        <Text style={styles.streakNumber}>{streak}</Text>
        <Text style={styles.streakLabel}>day streak</Text>
      </View>

      {/* New log button */}
      <TouchableOpacity style={styles.newLogButton} onPress={handleCreateLog}>
        <Text style={styles.newLogButtonText}>
          {getLogsForDate(logs, new Date()).length > 0
            ? '✓ Logged Today'
            : '+ Log Today'}
        </Text>
      </TouchableOpacity>

      {/* Toggle streak list */}
      <TouchableOpacity
        style={styles.toggleButton}
        onPress={() => setShowLogs(!showLogs)}>
        <Text style={styles.toggleButtonText}>
          {showLogs ? 'Hide History' : 'Show History'}
        </Text>
      </TouchableOpacity>

      {/* Streak list */}
      {showLogs && (
        <ScrollView
          horizontal
          showsHorizontalScrollIndicator={false}
          style={styles.streakList}>
          {days.map((date, index) => {
            const hasLog = hasLogForDate(logs, date);
            const isTodayDate = isToday(date);

            return (
              <TouchableOpacity
                key={index}
                style={[
                  styles.dayContainer,
                  hasLog && styles.dayContainerActive,
                  isTodayDate && styles.dayContainerToday,
                ]}
                onPress={() => hasLog && handleLogClick(date)}>
                <Text style={styles.dayText}>{formatDayShort(date)}</Text>
                <Text style={styles.dateText}>{formatDateShort(date)}</Text>
                {hasLog && <View style={styles.logIndicator} />}
              </TouchableOpacity>
            );
          })}
        </ScrollView>
      )}

      {/* Log modal would be rendered here if we import it */}
      <LogDetailsModal
        visible={showLogModal}
        log={selectedLog}
        date={selectedDate || undefined}
        onClose={() => {
          setShowLogModal(false);
          setSelectedLog(null);
          setSelectedDate(null);
        }}
        onSave={handleSaveLog}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: theme.colors.surface,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.md,
    marginBottom: theme.spacing.md,
    ...theme.shadows.medium,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: theme.spacing.md,
  },
  headerLeft: {
    flex: 1,
    marginRight: theme.spacing.md,
  },
  name: {
    fontSize: theme.fontSize.lg,
    fontWeight: 'bold',
    color: theme.colors.text,
    marginBottom: theme.spacing.xs,
  },
  description: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textLight,
  },
  actions: {
    flexDirection: 'row',
    gap: theme.spacing.sm,
  },
  actionButton: {
    padding: theme.spacing.xs,
  },
  actionButtonText: {
    fontSize: theme.fontSize.lg,
  },
  streakContainer: {
    alignItems: 'center',
    paddingVertical: theme.spacing.md,
    backgroundColor: theme.colors.primaryLight,
    borderRadius: theme.borderRadius.md,
    marginBottom: theme.spacing.md,
  },
  streakNumber: {
    fontSize: theme.fontSize.xxl,
    fontWeight: 'bold',
    color: theme.colors.primary,
  },
  streakLabel: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textLight,
  },
  newLogButton: {
    backgroundColor: theme.colors.primary,
    paddingVertical: theme.spacing.sm,
    paddingHorizontal: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
    alignItems: 'center',
    marginBottom: theme.spacing.sm,
  },
  newLogButtonText: {
    color: theme.colors.surface,
    fontSize: theme.fontSize.md,
    fontWeight: '600',
  },
  toggleButton: {
    paddingVertical: theme.spacing.sm,
    alignItems: 'center',
  },
  toggleButtonText: {
    color: theme.colors.primary,
    fontSize: theme.fontSize.sm,
    fontWeight: '600',
  },
  streakList: {
    marginTop: theme.spacing.md,
  },
  dayContainer: {
    width: 60,
    height: 70,
    backgroundColor: theme.colors.streakEmpty,
    borderRadius: theme.borderRadius.md,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: theme.spacing.sm,
    padding: theme.spacing.xs,
  },
  dayContainerActive: {
    backgroundColor: theme.colors.streakActive,
  },
  dayContainerToday: {
    borderWidth: 2,
    borderColor: theme.colors.primary,
  },
  dayText: {
    fontSize: theme.fontSize.xs,
    color: theme.colors.textLight,
    marginBottom: theme.spacing.xs,
  },
  dateText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.text,
    fontWeight: '600',
  },
  logIndicator: {
    width: 6,
    height: 6,
    borderRadius: 3,
    backgroundColor: theme.colors.success,
    marginTop: theme.spacing.xs,
  },
});
