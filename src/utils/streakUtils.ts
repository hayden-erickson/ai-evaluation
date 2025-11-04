/**
 * Utility functions for calculating habit streaks
 */

import { Log } from '../types';
import { STREAK_DISPLAY_DAYS } from '../config';

/**
 * Calculate the current streak for a habit based on its logs
 * The streak continues if the user logs daily or skips at most 1 day
 * 
 * @param logs - Array of logs sorted by created_at descending (newest first)
 * @returns The current streak count
 */
export function calculateStreak(logs: Log[]): number {
  if (!logs || logs.length === 0) {
    return 0;
  }

  // Sort logs by date descending (newest first)
  const sortedLogs = [...logs].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );

  let streak = 0;
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  // Check each log starting from most recent
  for (let i = 0; i < sortedLogs.length; i++) {
    const logDate = new Date(sortedLogs[i].created_at);
    logDate.setHours(0, 0, 0, 0);

    // Calculate expected date for this position in the streak
    // The first log should be today or yesterday, subsequent logs continue the pattern
    const expectedDate = new Date(today);
    expectedDate.setDate(today.getDate() - streak);

    const daysDifference = Math.floor(
      (expectedDate.getTime() - logDate.getTime()) / (1000 * 60 * 60 * 24)
    );

    // If log is on the expected day, increment streak
    if (daysDifference === 0) {
      streak++;
    } 
    // If log is 1 day before expected (user skipped a day), still count it
    else if (daysDifference === 1) {
      streak++;
    } 
    // If gap is more than 1 day, streak is broken
    else if (daysDifference > 1) {
      break;
    }
  }

  // Check if the streak is still active (last log should be recent)
  const mostRecentLog = new Date(sortedLogs[0].created_at);
  mostRecentLog.setHours(0, 0, 0, 0);
  
  const daysSinceLastLog = Math.floor(
    (today.getTime() - mostRecentLog.getTime()) / (1000 * 60 * 60 * 24)
  );

  // If more than 1 day has passed since last log, streak is broken
  if (daysSinceLastLog > 1) {
    return 0;
  }

  return streak;
}

/**
 * Get logs grouped by date for display in the streak list
 * 
 * @param logs - Array of logs
 * @param days - Number of days to display (default from config)
 * @returns Object mapping date strings to logs
 */
export function getLogsGroupedByDate(
  logs: Log[],
  days: number = STREAK_DISPLAY_DAYS
): Record<string, Log | null> {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const result: Record<string, Log | null> = {};

  // Create entries for the last N days
  for (let i = 0; i < days; i++) {
    const date = new Date(today);
    date.setDate(today.getDate() - i);
    const dateKey = formatDate(date);
    result[dateKey] = null;
  }

  // Fill in logs for dates that have them
  logs.forEach(log => {
    const logDate = new Date(log.created_at);
    logDate.setHours(0, 0, 0, 0);
    const dateKey = formatDate(logDate);
    
    if (dateKey in result) {
      // If multiple logs on same day, keep the most recent
      if (!result[dateKey] || 
          new Date(log.created_at) > new Date(result[dateKey]!.created_at)) {
        result[dateKey] = log;
      }
    }
  });

  return result;
}

/**
 * Format a date to YYYY-MM-DD string
 */
export function formatDate(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

/**
 * Format a date for display (e.g., "Mon, Dec 25")
 */
export function formatDateForDisplay(date: Date): string {
  const options: Intl.DateTimeFormatOptions = {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
  };
  return date.toLocaleDateString('en-US', options);
}

/**
 * Check if a date is today
 */
export function isToday(date: Date): boolean {
  const today = new Date();
  return (
    date.getDate() === today.getDate() &&
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()
  );
}

/**
 * Get the start of today (00:00:00)
 */
export function getTodayStart(): Date {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  return today;
}
