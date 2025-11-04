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

  const today = new Date();
  today.setHours(0, 0, 0, 0);
  
  // Check if the most recent log is within allowed range (today or yesterday)
  const mostRecentLog = new Date(sortedLogs[0].created_at);
  mostRecentLog.setHours(0, 0, 0, 0);
  
  const daysSinceLastLog = Math.floor(
    (today.getTime() - mostRecentLog.getTime()) / (1000 * 60 * 60 * 24)
  );

  // If more than 1 day has passed since last log, streak is broken
  if (daysSinceLastLog > 1) {
    return 0;
  }

  let streak = 1; // Start with the most recent log
  let previousLogDate = mostRecentLog;

  // Check each subsequent log
  for (let i = 1; i < sortedLogs.length; i++) {
    const currentLogDate = new Date(sortedLogs[i].created_at);
    currentLogDate.setHours(0, 0, 0, 0);

    const daysBetween = Math.floor(
      (previousLogDate.getTime() - currentLogDate.getTime()) / (1000 * 60 * 60 * 24)
    );

    // If logs are consecutive days (1 day apart) or allow 1-day skip (2 days apart)
    if (daysBetween === 1 || daysBetween === 2) {
      streak++;
      previousLogDate = currentLogDate;
    } 
    // If gap is more than 2 days, streak is broken
    else if (daysBetween > 2) {
      break;
    }
    // If daysBetween === 0, it's the same day, skip this log
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
