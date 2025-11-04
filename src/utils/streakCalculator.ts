/**
 * Utilities for calculating habit streaks
 */

import { Log } from '../types/models';

/**
 * Check if two dates are on the same day (ignoring time)
 */
export const isSameDay = (date1: Date, date2: Date): boolean => {
  return (
    date1.getFullYear() === date2.getFullYear() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getDate() === date2.getDate()
  );
};

/**
 * Get the number of days between two dates
 */
export const getDaysDifference = (date1: Date, date2: Date): number => {
  const msPerDay = 24 * 60 * 60 * 1000;
  const utc1 = Date.UTC(date1.getFullYear(), date1.getMonth(), date1.getDate());
  const utc2 = Date.UTC(date2.getFullYear(), date2.getMonth(), date2.getDate());
  return Math.floor((utc2 - utc1) / msPerDay);
};

/**
 * Calculate the current streak for a habit based on its logs
 * Users are allowed to skip one day before the streak resets
 */
export const calculateStreak = (logs: Log[]): number => {
  if (logs.length === 0) {
    return 0;
  }

  // Sort logs by date (most recent first)
  const sortedLogs = [...logs].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );

  const today = new Date();
  const mostRecentLog = new Date(sortedLogs[0].created_at);
  
  // Check if the most recent log is today or within the allowed skip period (1 day)
  const daysSinceLastLog = getDaysDifference(mostRecentLog, today);
  
  if (daysSinceLastLog > 2) {
    // More than 2 days (allowing 1 skip day) - streak is broken
    return 0;
  }

  let streak = 0;
  let currentDate = new Date(mostRecentLog);
  
  // Group logs by date
  const logsByDate = new Map<string, Log[]>();
  sortedLogs.forEach(log => {
    const logDate = new Date(log.created_at);
    const dateKey = `${logDate.getFullYear()}-${logDate.getMonth()}-${logDate.getDate()}`;
    if (!logsByDate.has(dateKey)) {
      logsByDate.set(dateKey, []);
    }
    logsByDate.get(dateKey)!.push(log);
  });

  let skippedDays = 0;
  let checkDate = new Date(currentDate);
  
  // Count backwards from the most recent log
  while (true) {
    const dateKey = `${checkDate.getFullYear()}-${checkDate.getMonth()}-${checkDate.getDate()}`;
    
    if (logsByDate.has(dateKey)) {
      // Found a log for this date
      streak++;
      skippedDays = 0; // Reset skip counter
    } else {
      // No log for this date
      skippedDays++;
      if (skippedDays > 1) {
        // More than 1 day skipped - streak is broken
        break;
      }
    }
    
    // Move to previous day
    checkDate = new Date(checkDate);
    checkDate.setDate(checkDate.getDate() - 1);
    
    // Stop if we've gone too far back
    if (streak > 0 && getDaysDifference(checkDate, mostRecentLog) > 365) {
      break;
    }
  }

  return streak;
};

/**
 * Get logs for the last N days
 */
export const getRecentDays = (logs: Log[], numDays: number): { date: Date; log: Log | null }[] => {
  const days: { date: Date; log: Log | null }[] = [];
  const today = new Date();
  
  // Group logs by date
  const logsByDate = new Map<string, Log>();
  logs.forEach(log => {
    const logDate = new Date(log.created_at);
    const dateKey = `${logDate.getFullYear()}-${logDate.getMonth()}-${logDate.getDate()}`;
    // Keep the most recent log for each date if there are multiple
    if (!logsByDate.has(dateKey)) {
      logsByDate.set(dateKey, log);
    }
  });
  
  // Create array of days with their logs
  for (let i = 0; i < numDays; i++) {
    const date = new Date(today);
    date.setDate(date.getDate() - i);
    
    const dateKey = `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}`;
    const log = logsByDate.get(dateKey) || null;
    
    days.push({ date, log });
  }
  
  return days;
};

/**
 * Check if a log exists for today
 */
export const hasLoggedToday = (logs: Log[]): boolean => {
  const today = new Date();
  return logs.some(log => isSameDay(new Date(log.created_at), today));
};

/**
 * Get the date of the last log
 */
export const getLastLogDate = (logs: Log[]): Date | null => {
  if (logs.length === 0) {
    return null;
  }
  
  const sortedLogs = [...logs].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );
  
  return new Date(sortedLogs[0].created_at);
};

