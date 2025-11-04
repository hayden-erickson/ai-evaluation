/**
 * Habit Component
 * Displays a single habit with streak information and log history
 */

import React, {useState, useEffect} from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Alert,
  ActivityIndicator,
} from 'react-native';
import {
  Habit as HabitType,
  Log,
  CreateLogRequest,
  UpdateLogRequest,
} from '../types';
import {colors, spacing, borderRadius, fontSize, shadow} from '../styles/theme';
import {
  calculateStreak,
  formatDay,
  getLastNDays,
  isSameDay,
} from '../utils/helpers';
import {LogDetailsModal} from './LogDetailsModal';
import {logApi} from '../services/api';

interface HabitProps {
  habit: HabitType;
  token: string;
  onEdit: () => void;
  onDelete: () => void;
  onUpdate: () => void; // Callback to refresh habit data
}

/**
 * Component that displays a habit with its streak and logs
 */
export function Habit({habit, token, onEdit, onDelete, onUpdate}: HabitProps) {
  const [logs, setLogs] = useState<Log[]>([]);
  const [isLoadingLogs, setIsLoadingLogs] = useState(true);
  const [isExpanded, setIsExpanded] = useState(false);
  const [logModalVisible, setLogModalVisible] = useState(false);
  const [selectedLog, setSelectedLog] = useState<Log | null>(null);
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);

  // Load logs when component mounts
  useEffect(() => {
    loadLogs();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  /**
   * Loads all logs for this habit
   */
  const loadLogs = async () => {
    setIsLoadingLogs(true);
    try {
      const habitLogs = await logApi.getHabitLogs(habit.id, token);
      setLogs(habitLogs);
    } catch (err) {
      console.error('Failed to load logs:', err);
      Alert.alert('Error', 'Failed to load habit logs');
    } finally {
      setIsLoadingLogs(false);
    }
  };

  /**
   * Creates a new log for today
   */
  const handleNewLog = () => {
    setSelectedLog(null);
    setSelectedDate(new Date());
    setLogModalVisible(true);
  };

  /**
   * Opens modal to edit an existing log
   */
  const handleEditLog = (log: Log) => {
    setSelectedLog(log);
    setSelectedDate(null);
    setLogModalVisible(true);
  };

  /**
   * Saves a log (create or update)
   */
  const handleSaveLog = async (data: CreateLogRequest | UpdateLogRequest) => {
    try {
      if (selectedLog) {
        // Update existing log
        await logApi.updateLog(selectedLog.id, data as UpdateLogRequest, token);
      } else {
        // Create new log
        await logApi.createLog(habit.id, data as CreateLogRequest, token);
      }

      // Reload logs and notify parent
      await loadLogs();
      onUpdate();
    } catch (error) {
      throw error; // Let modal handle the error display
    }
  };

  /**
   * Deletes a log with confirmation
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
            await logApi.deleteLog(log.id, token);
            await loadLogs();
            onUpdate();
          },
        },
      ],
    );
  };

  /**
   * Confirms and deletes the habit
   */
  const confirmDelete = () => {
    Alert.alert(
      'Delete Habit',
      `Are you sure you want to delete "${habit.name}"? This will also delete all associated logs.`,
      [
        {text: 'Cancel', style: 'cancel'},
        {
          text: 'Delete',
          style: 'destructive',
          onPress: onDelete,
        },
      ],
    );
  };

  // Calculate streak
  const streakInfo = calculateStreak(logs);
  const lastNDays = getLastNDays(7);

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <View style={styles.headerLeft}>
          <Text style={styles.name}>{habit.name}</Text>
          {habit.description && (
            <Text style={styles.description}>{habit.description}</Text>
          )}
        </View>

        <View style={styles.headerRight}>
          <TouchableOpacity
            style={styles.iconButton}
            onPress={onEdit}
            accessibilityLabel="Edit habit">
            <Text style={styles.iconText}>✏️</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={styles.iconButton}
            onPress={confirmDelete}
            accessibilityLabel="Delete habit">
            <Text style={styles.iconText}>🗑️</Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Streak Counter */}
      <View style={styles.streakContainer}>
        <Text style={styles.streakNumber}>{streakInfo.currentStreak}</Text>
        <Text style={styles.streakLabel}>
          Day Streak {streakInfo.canContinue ? '🔥' : ''}
        </Text>
      </View>

      {/* New Log Button */}
      <TouchableOpacity style={styles.newLogButton} onPress={handleNewLog}>
        <Text style={styles.newLogButtonText}>+ Log Today</Text>
      </TouchableOpacity>

      {/* Expand/Collapse Button */}
      <TouchableOpacity
        style={styles.expandButton}
        onPress={() => setIsExpanded(!isExpanded)}>
        <Text style={styles.expandButtonText}>
          {isExpanded ? 'Hide History ▲' : 'Show History ▼'}
        </Text>
      </TouchableOpacity>

      {/* Streak List (Last 7 days) */}
      {isExpanded && (
        <View style={styles.streakList}>
          {isLoadingLogs ? (
            <ActivityIndicator size="small" color={colors.primary} />
          ) : (
            lastNDays.map(day => {
              const dayLog = logs.find(log =>
                isSameDay(new Date(log.created_at), day),
              );
              const isToday = isSameDay(day, new Date());

              return (
                <View key={day.toISOString()} style={styles.dayContainer}>
                  <Text style={styles.dayLabel}>
                    {formatDay(day)} {isToday && '(Today)'}
                  </Text>
                  {dayLog ? (
                    <TouchableOpacity
                      style={styles.logItem}
                      onPress={() => handleEditLog(dayLog)}
                      onLongPress={() => handleDeleteLog(dayLog)}>
                      <Text style={styles.logItemText}>✓</Text>
                      {dayLog.notes && (
                        <Text style={styles.logNotes} numberOfLines={1}>
                          {dayLog.notes}
                        </Text>
                      )}
                    </TouchableOpacity>
                  ) : (
                    <View style={styles.emptyLog}>
                      <Text style={styles.emptyLogText}>-</Text>
                    </View>
                  )}
                </View>
              );
            })
          )}
        </View>
      )}

      {/* Log Details Modal */}
      <LogDetailsModal
        visible={logModalVisible}
        log={selectedLog}
        logDate={selectedDate || undefined}
        onClose={() => {
          setLogModalVisible(false);
          setSelectedLog(null);
          setSelectedDate(null);
        }}
        onSave={handleSaveLog}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.surfaceLight,
    borderRadius: borderRadius.lg,
    padding: spacing.lg,
    marginBottom: spacing.md,
    ...shadow.sm,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: spacing.md,
  },
  headerLeft: {
    flex: 1,
  },
  headerRight: {
    flexDirection: 'row',
    gap: spacing.sm,
  },
  name: {
    fontSize: fontSize.lg,
    fontWeight: '600',
    color: colors.text,
    marginBottom: spacing.xs,
  },
  description: {
    fontSize: fontSize.sm,
    color: colors.textLight,
  },
  iconButton: {
    padding: spacing.xs,
  },
  iconText: {
    fontSize: fontSize.lg,
  },
  streakContainer: {
    alignItems: 'center',
    paddingVertical: spacing.md,
    backgroundColor: colors.surface,
    borderRadius: borderRadius.md,
    marginBottom: spacing.md,
  },
  streakNumber: {
    fontSize: fontSize.xxxl,
    fontWeight: 'bold',
    color: colors.text,
  },
  streakLabel: {
    fontSize: fontSize.sm,
    color: colors.textLight,
    marginTop: spacing.xs,
  },
  newLogButton: {
    backgroundColor: colors.primary,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    alignItems: 'center',
    marginBottom: spacing.sm,
  },
  newLogButtonText: {
    fontSize: fontSize.md,
    fontWeight: '600',
    color: colors.text,
  },
  expandButton: {
    alignItems: 'center',
    padding: spacing.sm,
  },
  expandButtonText: {
    fontSize: fontSize.sm,
    color: colors.textLight,
  },
  streakList: {
    marginTop: spacing.md,
    gap: spacing.sm,
  },
  dayContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: spacing.sm,
    borderBottomWidth: 1,
    borderBottomColor: colors.borderLight,
  },
  dayLabel: {
    fontSize: fontSize.sm,
    color: colors.textLight,
    flex: 1,
  },
  logItem: {
    backgroundColor: colors.streakActive,
    borderRadius: borderRadius.sm,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    flex: 2,
  },
  logItemText: {
    fontSize: fontSize.md,
    fontWeight: '600',
  },
  logNotes: {
    fontSize: fontSize.sm,
    color: colors.textLight,
    flex: 1,
  },
  emptyLog: {
    backgroundColor: colors.streakInactive,
    borderRadius: borderRadius.sm,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    alignItems: 'center',
    flex: 2,
  },
  emptyLogText: {
    fontSize: fontSize.md,
    color: colors.textMuted,
  },
});
