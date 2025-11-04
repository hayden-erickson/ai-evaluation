/**
 * Habit component - Displays a single habit with its streak and logs
 * Shows name, description, streak count, and recent activity
 */

import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { Habit as HabitType, Log, CreateLogRequest, UpdateLogRequest } from '../types';
import { apiService } from '../services/api';
import { calculateStreak, getLogsGroupedByDate, formatDateForDisplay, getTodayStart } from '../utils/streakUtils';
import { Colors, Spacing, BorderRadius, Typography, CommonStyles } from '../styles/theme';
import LogDetailsModal from './LogDetailsModal';

interface HabitProps {
  habit: HabitType;
  onEdit: () => void;
  onDelete: () => void;
  onUpdate: () => void;  // Callback to refresh habit list
}

export default function Habit({ habit, onEdit, onDelete, onUpdate }: HabitProps) {
  const [logs, setLogs] = useState<Log[]>([]);
  const [isLoadingLogs, setIsLoadingLogs] = useState(true);
  const [isCreatingLog, setIsCreatingLog] = useState(false);
  const [showLogModal, setShowLogModal] = useState(false);
  const [selectedLog, setSelectedLog] = useState<Log | null>(null);
  const [streak, setStreak] = useState(0);

  // Load logs when component mounts or habit changes
  useEffect(() => {
    loadLogs();
  }, [habit.id]);

  // Recalculate streak when logs change
  useEffect(() => {
    setStreak(calculateStreak(logs));
  }, [logs]);

  /**
   * Load all logs for this habit from the API
   */
  const loadLogs = async () => {
    try {
      setIsLoadingLogs(true);
      const fetchedLogs = await apiService.getHabitLogs(habit.id);
      setLogs(fetchedLogs);
    } catch (error) {
      Alert.alert('Error', 'Failed to load habit logs');
      console.error('Failed to load logs:', error);
    } finally {
      setIsLoadingLogs(false);
    }
  };

  /**
   * Handle creating a new log for today
   */
  const handleNewLog = () => {
    // Check if there's already a log for today
    const today = getTodayStart();
    const hasLogToday = logs.some(log => {
      const logDate = new Date(log.created_at);
      logDate.setHours(0, 0, 0, 0);
      return logDate.getTime() === today.getTime();
    });

    if (hasLogToday) {
      Alert.alert('Already Logged', 'You have already logged this habit for today.');
      return;
    }

    setSelectedLog(null);
    setShowLogModal(true);
  };

  /**
   * Handle editing an existing log
   */
  const handleEditLog = (log: Log) => {
    setSelectedLog(log);
    setShowLogModal(true);
  };

  /**
   * Handle saving a log (create or update)
   */
  const handleSaveLog = async (data: CreateLogRequest | UpdateLogRequest) => {
    try {
      if (selectedLog) {
        // Update existing log
        await apiService.updateLog(selectedLog.id, data as UpdateLogRequest);
      } else {
        // Create new log
        await apiService.createLog(habit.id, data as CreateLogRequest);
      }
      
      // Reload logs to reflect changes
      await loadLogs();
      onUpdate();  // Notify parent to refresh
    } catch (error) {
      throw new Error(error instanceof Error ? error.message : 'Failed to save log');
    }
  };

  /**
   * Handle deleting a log
   */
  const handleDeleteLog = (log: Log) => {
    Alert.alert(
      'Delete Log',
      'Are you sure you want to delete this log?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Delete',
          style: 'destructive',
          onPress: async () => {
            try {
              await apiService.deleteLog(log.id);
              await loadLogs();
              onUpdate();
            } catch (error) {
              Alert.alert('Error', 'Failed to delete log');
            }
          },
        },
      ]
    );
  };

  // Get last 7 days of logs grouped by date
  const logsGroupedByDate = getLogsGroupedByDate(logs, 7);
  const dates = Object.keys(logsGroupedByDate).sort((a, b) => 
    new Date(b).getTime() - new Date(a).getTime()
  );

  return (
    <View style={[CommonStyles.card, styles.habitCard]}>
      {/* Habit header with name and actions */}
      <View style={styles.header}>
        <View style={styles.headerLeft}>
          <Text style={styles.habitName}>{habit.name}</Text>
          {habit.description ? (
            <Text style={styles.habitDescription}>{habit.description}</Text>
          ) : null}
        </View>
        
        <View style={styles.headerActions}>
          <TouchableOpacity onPress={onEdit} style={styles.iconButton}>
            <Text style={styles.iconText}>✏️</Text>
          </TouchableOpacity>
          <TouchableOpacity onPress={onDelete} style={styles.iconButton}>
            <Text style={styles.iconText}>🗑️</Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Streak display */}
      <View style={styles.streakContainer}>
        <Text style={styles.streakLabel}>Current Streak:</Text>
        <View style={styles.streakBadge}>
          <Text style={styles.streakCount}>{streak}</Text>
          <Text style={styles.streakUnit}>{streak === 1 ? 'day' : 'days'}</Text>
        </View>
      </View>

      {/* New log button */}
      <TouchableOpacity
        style={[CommonStyles.button, styles.newLogButton]}
        onPress={handleNewLog}
        disabled={isCreatingLog}
      >
        {isCreatingLog ? (
          <ActivityIndicator color={Colors.textPrimary} />
        ) : (
          <Text style={CommonStyles.buttonText}>Log Today</Text>
        )}
      </TouchableOpacity>

      {/* Streak list - last 7 days */}
      {isLoadingLogs ? (
        <View style={styles.loadingContainer}>
          <ActivityIndicator color={Colors.primary} />
        </View>
      ) : (
        <View style={styles.streakList}>
          <Text style={styles.streakListTitle}>Recent Activity</Text>
          {dates.map(dateKey => {
            const log = logsGroupedByDate[dateKey];
            const date = new Date(dateKey);
            const isToday = date.toDateString() === new Date().toDateString();
            
            return (
              <TouchableOpacity
                key={dateKey}
                style={[
                  styles.dayContainer,
                  log ? styles.dayContainerActive : styles.dayContainerInactive,
                ]}
                onPress={() => log && handleEditLog(log)}
                onLongPress={() => log && handleDeleteLog(log)}
                disabled={!log}
              >
                <View style={styles.dayLeft}>
                  <Text style={styles.dayDate}>
                    {formatDateForDisplay(date)}
                    {isToday ? ' (Today)' : ''}
                  </Text>
                  {log && log.notes ? (
                    <Text style={styles.dayNotes} numberOfLines={1}>
                      {log.notes}
                    </Text>
                  ) : null}
                </View>
                <View style={styles.dayRight}>
                  {log ? (
                    <Text style={styles.dayCheck}>✓</Text>
                  ) : (
                    <Text style={styles.dayEmpty}>−</Text>
                  )}
                </View>
              </TouchableOpacity>
            );
          })}
        </View>
      )}

      {/* Log details modal */}
      <LogDetailsModal
        visible={showLogModal}
        onClose={() => {
          setShowLogModal(false);
          setSelectedLog(null);
        }}
        onSave={handleSaveLog}
        log={selectedLog}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  habitCard: {
    marginHorizontal: Spacing.md,
  },
  
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: Spacing.md,
  },
  
  headerLeft: {
    flex: 1,
  },
  
  headerActions: {
    flexDirection: 'row',
    gap: Spacing.sm,
  },
  
  iconButton: {
    padding: Spacing.xs,
  },
  
  iconText: {
    fontSize: Typography.h4,
  },
  
  habitName: {
    fontSize: Typography.h3,
    fontWeight: Typography.bold,
    color: Colors.textPrimary,
    marginBottom: Spacing.xs,
  },
  
  habitDescription: {
    fontSize: Typography.small,
    color: Colors.textSecondary,
  },
  
  streakContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: Spacing.md,
    gap: Spacing.sm,
  },
  
  streakLabel: {
    fontSize: Typography.body,
    color: Colors.textSecondary,
  },
  
  streakBadge: {
    backgroundColor: Colors.streakActive,
    borderRadius: BorderRadius.full,
    paddingHorizontal: Spacing.md,
    paddingVertical: Spacing.xs,
    flexDirection: 'row',
    alignItems: 'center',
    gap: Spacing.xs,
  },
  
  streakCount: {
    fontSize: Typography.h3,
    fontWeight: Typography.bold,
    color: Colors.textPrimary,
  },
  
  streakUnit: {
    fontSize: Typography.small,
    color: Colors.textPrimary,
  },
  
  newLogButton: {
    marginBottom: Spacing.md,
  },
  
  loadingContainer: {
    padding: Spacing.lg,
    alignItems: 'center',
  },
  
  streakList: {
    marginTop: Spacing.sm,
  },
  
  streakListTitle: {
    fontSize: Typography.body,
    fontWeight: Typography.semibold,
    color: Colors.textPrimary,
    marginBottom: Spacing.sm,
  },
  
  dayContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: Spacing.md,
    borderRadius: BorderRadius.sm,
    marginBottom: Spacing.xs,
  },
  
  dayContainerActive: {
    backgroundColor: Colors.streakActive,
  },
  
  dayContainerInactive: {
    backgroundColor: Colors.streakInactive,
  },
  
  dayLeft: {
    flex: 1,
  },
  
  dayDate: {
    fontSize: Typography.small,
    fontWeight: Typography.medium,
    color: Colors.textPrimary,
  },
  
  dayNotes: {
    fontSize: Typography.small,
    color: Colors.textSecondary,
    marginTop: Spacing.xs,
  },
  
  dayRight: {
    marginLeft: Spacing.sm,
  },
  
  dayCheck: {
    fontSize: Typography.h3,
    color: Colors.textPrimary,
  },
  
  dayEmpty: {
    fontSize: Typography.h3,
    color: Colors.textSecondary,
  },
});
