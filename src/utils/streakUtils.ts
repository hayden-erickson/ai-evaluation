/**
 * Utility functions for habit streak calculations
 */

import {Log} from '../types';

/**
 * Calculate the current streak for a habit
 * Users are allowed to skip one day before their streak resets
 * 
 * @param logs - Array of log entries for a habit (will be sorted internally)
 * @returns The current streak count
 */
export const calculateStreak = (logs: Log[]): number => {
  if (logs.length === 0) {
    return 0;
  }

  // Sort logs by date (most recent first)
  const sortedLogs = [...logs].sort(
    (a, b) =>
      new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  );

  let streak = 0;
  let lastLogDate: Date | null = null;
  let skippedDay = false;

  for (const log of sortedLogs) {
    const logDate = new Date(log.created_at);
    // Normalize to start of day for comparison
    logDate.setHours(0, 0, 0, 0);

    if (lastLogDate === null) {
      // First log - check if it's today or yesterday (or day before if skipped)
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      
      const daysDiff = Math.floor(
        (today.getTime() - logDate.getTime()) / (1000 * 60 * 60 * 24),
      );

      if (daysDiff === 0 || daysDiff === 1) {
        // Today or yesterday
        streak = 1;
        lastLogDate = logDate;
      } else if (daysDiff === 2) {
        // Day before yesterday (allowed skip)
        streak = 1;
        lastLogDate = logDate;
        skippedDay = true;
      } else {
        // Too old, streak is broken
        break;
      }
    } else {
      // Calculate days between this log and the last one
      const daysDiff = Math.floor(
        (lastLogDate.getTime() - logDate.getTime()) / (1000 * 60 * 60 * 24),
      );

      if (daysDiff === 1) {
        // Consecutive day
        streak++;
        lastLogDate = logDate;
      } else if (daysDiff === 2 && !skippedDay) {
        // One day skipped (allowed once)
        streak++;
        lastLogDate = logDate;
        skippedDay = true;
      } else if (daysDiff === 0) {
        // Multiple logs on the same day - don't count as extra streak
        continue;
      } else {
        // Gap is too large, streak is broken
        break;
      }
    }
  }

  return streak;
};

/**
 * Check if a log exists for a specific date
 * 
 * @param logs - Array of log entries
 * @param date - Date to check
 * @returns True if a log exists for the date
 */
export const hasLogForDate = (logs: Log[], date: Date): boolean => {
  const targetDate = new Date(date);
  targetDate.setHours(0, 0, 0, 0);

  return logs.some(log => {
    const logDate = new Date(log.created_at);
    logDate.setHours(0, 0, 0, 0);
    return logDate.getTime() === targetDate.getTime();
  });
};

/**
 * Get logs for a specific date
 * 
 * @param logs - Array of log entries
 * @param date - Date to get logs for
 * @returns Array of logs for the date
 */
export const getLogsForDate = (logs: Log[], date: Date): Log[] => {
  const targetDate = new Date(date);
  targetDate.setHours(0, 0, 0, 0);

  return logs.filter(log => {
    const logDate = new Date(log.created_at);
    logDate.setHours(0, 0, 0, 0);
    return logDate.getTime() === targetDate.getTime();
  });
};

/**
 * Get the last N days as an array of dates
 * 
 * @param days - Number of days to get
 * @returns Array of dates
 */
export const getLastNDays = (days: number): Date[] => {
  const dates: Date[] = [];
  const today = new Date();

  for (let i = 0; i < days; i++) {
    const date = new Date(today);
    date.setDate(today.getDate() - i);
    date.setHours(0, 0, 0, 0);
    dates.push(date);
  }

  return dates;
};

/**
 * Format a date as a short string (e.g., "Mon", "Tue")
 * 
 * @param date - Date to format
 * @returns Formatted date string
 */
export const formatDayShort = (date: Date): string => {
  const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  return days[date.getDay()];
};

/**
 * Format a date as a short date string (e.g., "1/15")
 * 
 * @param date - Date to format
 * @returns Formatted date string
 */
export const formatDateShort = (date: Date): string => {
  return `${date.getMonth() + 1}/${date.getDate()}`;
};

/**
 * Check if a date is today
 * 
 * @param date - Date to check
 * @returns True if the date is today
 */
export const isToday = (date: Date): boolean => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const checkDate = new Date(date);
  checkDate.setHours(0, 0, 0, 0);
  return today.getTime() === checkDate.getTime();
};
