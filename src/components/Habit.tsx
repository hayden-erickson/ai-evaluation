/**
 * Habit component displaying habit details and logs
 */

import React, { useState } from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { Habit as HabitModel, Log, CreateLogRequest, UpdateLogRequest } from '../types/models';
import { calculateStreak, hasLoggedToday } from '../utils/streakCalculator';
import { theme } from '../styles/theme';
import { commonStyles } from '../styles/commonStyles';
import { StreakList } from './StreakList';
import { LogDetailsModal } from './LogDetailsModal';

interface HabitProps {
  habit: HabitModel;
  logs: Log[];
  onEdit: () => void;
  onDelete: () => void;
  onCreateLog: (data: CreateLogRequest) => Promise<void>;
  onUpdateLog: (logId: number, data: UpdateLogRequest) => Promise<void>;
}

/**
 * Display a single habit with its details, streak, and logs
 */
export const Habit: React.FC<HabitProps> = ({
  habit,
  logs,
  onEdit,
  onDelete,
  onCreateLog,
  onUpdateLog,
}) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [isCreatingLog, setIsCreatingLog] = useState(false);
  const [isLogModalVisible, setIsLogModalVisible] = useState(false);
  const [selectedLog, setSelectedLog] = useState<Log | null>(null);

  const streak = calculateStreak(logs);
  const loggedToday = hasLoggedToday(logs);

  /**
   * Handle new log button press
   */
  const handleNewLogPress = () => {
    setSelectedLog(null);
    setIsLogModalVisible(true);
  };

  /**
   * Handle edit log press
   */
  const handleEditLog = (log: Log) => {
    setSelectedLog(log);
    setIsLogModalVisible(true);
  };

  /**
   * Handle save log
   */
  const handleSaveLog = async (data: CreateLogRequest | UpdateLogRequest) => {
    if (selectedLog) {
      // Update existing log
      await onUpdateLog(selectedLog.id, data as UpdateLogRequest);
    } else {
      // Create new log
      await onCreateLog(data as CreateLogRequest);
    }
  };

  /**
   * Handle delete button with confirmation
   */
  const handleDeletePress = () => {
    Alert.alert(
      'Delete Habit',
      `Are you sure you want to delete "${habit.name}"? This will also delete all associated logs.`,
      [
        { text: 'Cancel', style: 'cancel' },
        { text: 'Delete', style: 'destructive', onPress: onDelete },
      ]
    );
  };

  return (
    <View style={commonStyles.card}>
      {/* Header */}
      <TouchableOpacity onPress={() => setIsExpanded(!isExpanded)} activeOpacity={0.7}>
        <View style={styles.header}>
          <View style={styles.headerLeft}>
            <Text style={styles.habitName}>{habit.name}</Text>
            {habit.description ? (
              <Text style={styles.habitDescription} numberOfLines={isExpanded ? undefined : 2}>
                {habit.description}
              </Text>
            ) : null}
          </View>
          <View style={styles.streakBadge}>
            <Text style={styles.streakNumber}>{streak}</Text>
            <Text style={styles.streakLabel}>🔥</Text>
          </View>
        </View>
      </TouchableOpacity>

      {/* Action Buttons */}
      <View style={styles.actions}>
        <TouchableOpacity style={styles.actionButton} onPress={onEdit}>
          <Text style={styles.actionButtonText}>✏️ Edit</Text>
        </TouchableOpacity>

        <TouchableOpacity style={styles.actionButton} onPress={handleDeletePress}>
          <Text style={[styles.actionButtonText, styles.deleteText]}>🗑️ Delete</Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={[
            styles.newLogButton,
            loggedToday && styles.newLogButtonDisabled,
          ]}
          onPress={handleNewLogPress}
          disabled={loggedToday}
        >
          <Text style={styles.newLogButtonText}>
            {loggedToday ? '✓ Logged Today' : '+ Log Today'}
          </Text>
        </TouchableOpacity>
      </View>

      {/* Streak List (only show when expanded) */}
      {isExpanded && (
        <StreakList logs={logs} onLogPress={handleEditLog} />
      )}

      {/* Log Details Modal */}
      <LogDetailsModal
        visible={isLogModalVisible}
        log={selectedLog}
        onClose={() => {
          setIsLogModalVisible(false);
          setSelectedLog(null);
        }}
        onSave={handleSaveLog}
      />
    </View>
  );
};

const styles = StyleSheet.create({
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

  habitName: {
    fontSize: theme.fontSize.lg,
    fontWeight: 'bold',
    color: theme.colors.text,
    marginBottom: theme.spacing.xs,
  },

  habitDescription: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textSecondary,
    lineHeight: 18,
  },

  streakBadge: {
    backgroundColor: theme.colors.accent,
    borderRadius: theme.borderRadius.md,
    paddingVertical: theme.spacing.xs,
    paddingHorizontal: theme.spacing.sm,
    alignItems: 'center',
    minWidth: 60,
  },

  streakNumber: {
    fontSize: theme.fontSize.xl,
    fontWeight: 'bold',
    color: theme.colors.text,
  },

  streakLabel: {
    fontSize: 16,
    marginTop: 2,
  },

  actions: {
    flexDirection: 'row',
    gap: theme.spacing.sm,
    flexWrap: 'wrap',
  },

  actionButton: {
    backgroundColor: theme.colors.background,
    borderRadius: theme.borderRadius.sm,
    paddingVertical: theme.spacing.xs,
    paddingHorizontal: theme.spacing.sm,
    borderWidth: 1,
    borderColor: theme.colors.border,
  },

  actionButtonText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.text,
    fontWeight: '600',
  },

  deleteText: {
    color: theme.colors.error,
  },

  newLogButton: {
    backgroundColor: theme.colors.primary,
    borderRadius: theme.borderRadius.sm,
    paddingVertical: theme.spacing.xs,
    paddingHorizontal: theme.spacing.md,
    flex: 1,
    alignItems: 'center',
  },

  newLogButtonDisabled: {
    backgroundColor: theme.colors.success,
  },

  newLogButtonText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.surface,
    fontWeight: '600',
  },
});

