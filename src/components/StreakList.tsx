/**
 * Streak list component showing recent days and logs
 */

import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { Log } from '../types/models';
import { getRecentDays } from '../utils/streakCalculator';
import { theme } from '../styles/theme';

interface StreakListProps {
  logs: Log[];
  onLogPress: (log: Log) => void;
}

/**
 * Display a list of recent days with their logs
 */
export const StreakList: React.FC<StreakListProps> = ({ logs, onLogPress }) => {
  // Get the last 14 days
  const recentDays = getRecentDays(logs, 14);

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Recent Activity</Text>
      <View style={styles.daysContainer}>
        {recentDays.map((day, index) => (
          <DayContainer
            key={index}
            date={day.date}
            log={day.log}
            onPress={day.log ? () => onLogPress(day.log!) : undefined}
          />
        ))}
      </View>
    </View>
  );
};

interface DayContainerProps {
  date: Date;
  log: Log | null;
  onPress?: () => void;
}

/**
 * Display a single day with or without a log
 */
const DayContainer: React.FC<DayContainerProps> = ({ date, log, onPress }) => {
  const dayOfWeek = date.toLocaleDateString('en-US', { weekday: 'short' });
  const dayOfMonth = date.getDate();
  const isToday = new Date().toDateString() === date.toDateString();

  const containerStyle = [
    styles.dayContainer,
    log ? styles.dayWithLog : styles.dayWithoutLog,
    isToday && styles.todayContainer,
  ];

  const content = (
    <View style={containerStyle}>
      <Text style={[styles.dayOfWeek, log && styles.dayOfWeekActive]}>
        {dayOfWeek}
      </Text>
      <Text style={[styles.dayOfMonth, log && styles.dayOfMonthActive]}>
        {dayOfMonth}
      </Text>
      {log && <View style={styles.logIndicator} />}
    </View>
  );

  if (onPress) {
    return (
      <TouchableOpacity onPress={onPress} activeOpacity={0.7}>
        {content}
      </TouchableOpacity>
    );
  }

  return content;
};

const styles = StyleSheet.create({
  container: {
    marginTop: theme.spacing.md,
  },

  title: {
    fontSize: theme.fontSize.sm,
    fontWeight: '600',
    color: theme.colors.textSecondary,
    marginBottom: theme.spacing.sm,
  },

  daysContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: theme.spacing.xs,
  },

  dayContainer: {
    width: 44,
    height: 56,
    borderRadius: theme.borderRadius.sm,
    padding: theme.spacing.xs,
    alignItems: 'center',
    justifyContent: 'center',
    position: 'relative',
  },

  dayWithoutLog: {
    backgroundColor: theme.colors.background,
    borderWidth: 1,
    borderColor: theme.colors.border,
  },

  dayWithLog: {
    backgroundColor: theme.colors.success,
    borderWidth: 1,
    borderColor: theme.colors.success,
  },

  todayContainer: {
    borderWidth: 2,
    borderColor: theme.colors.primary,
  },

  dayOfWeek: {
    fontSize: 10,
    fontWeight: '600',
    color: theme.colors.textSecondary,
    marginBottom: 2,
  },

  dayOfWeekActive: {
    color: theme.colors.surface,
  },

  dayOfMonth: {
    fontSize: theme.fontSize.md,
    fontWeight: 'bold',
    color: theme.colors.text,
  },

  dayOfMonthActive: {
    color: theme.colors.surface,
  },

  logIndicator: {
    position: 'absolute',
    bottom: 4,
    width: 4,
    height: 4,
    borderRadius: 2,
    backgroundColor: theme.colors.surface,
  },
});

