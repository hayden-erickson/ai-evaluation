/**
 * Utility functions for the Habit Tracker application
 */

import {Log, StreakInfo} from '../types';

/**
 * Calculates the current streak for a habit based on its logs
 * Users are allowed to skip one day before their streak resets
 * @param logs - Array of logs for the habit, sorted by created_at
 * @returns StreakInfo object containing streak details
 */
export function calculateStreak(logs: Log[]): StreakInfo {
  if (!logs || logs.length === 0) {
    return {
      currentStreak: 0,
      lastLogDate: null,
      canContinue: true,
    };
  }

  // Sort logs by date (newest first)
  const sortedLogs = [...logs].sort(
    (a, b) =>
      new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  );

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const lastLog = sortedLogs[0];
  const lastLogDate = new Date(lastLog.created_at);
  lastLogDate.setHours(0, 0, 0, 0);

  // Calculate days since last log
  const daysSinceLastLog = Math.floor(
    (today.getTime() - lastLogDate.getTime()) / (1000 * 60 * 60 * 24),
  );

  // If more than 1 day has passed, streak is broken
  if (daysSinceLastLog > 1) {
    return {
      currentStreak: 0,
      lastLogDate: lastLogDate,
      canContinue: false,
    };
  }

  // Calculate the actual streak by checking consecutive days
  let streak = 0;
  let currentDate = new Date(lastLogDate);

  for (let i = 0; i < sortedLogs.length; i++) {
    const logDate = new Date(sortedLogs[i].created_at);
    logDate.setHours(0, 0, 0, 0);

    const daysDiff = Math.floor(
      (currentDate.getTime() - logDate.getTime()) / (1000 * 60 * 60 * 24),
    );

    // If this log is on the expected day or we're at the start
    if (daysDiff === 0) {
      streak++;
      currentDate.setDate(currentDate.getDate() - 1);
    }
    // If we skipped one day (allowed), continue but move current date back
    else if (daysDiff === 1) {
      currentDate = new Date(logDate);
      currentDate.setDate(currentDate.getDate() - 1);
    }
    // If gap is more than 1 day, streak ends
    else {
      break;
    }
  }

  return {
    currentStreak: streak,
    lastLogDate: lastLogDate,
    canContinue: daysSinceLastLog <= 1,
  };
}

/**
 * Formats a date to a readable string (e.g., "Jan 15, 2024")
 */
export function formatDate(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });
}

/**
 * Formats a date to just the day (e.g., "15")
 */
export function formatDay(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return d.getDate().toString();
}

/**
 * Formats a date to the month name (e.g., "January")
 */
export function formatMonth(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return d.toLocaleDateString('en-US', {month: 'long'});
}

/**
 * Checks if two dates are on the same day
 */
export function isSameDay(date1: Date, date2: Date): boolean {
  return (
    date1.getFullYear() === date2.getFullYear() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getDate() === date2.getDate()
  );
}

/**
 * Gets the last N days including today
 */
export function getLastNDays(n: number): Date[] {
  const days: Date[] = [];
  const today = new Date();

  for (let i = n - 1; i >= 0; i--) {
    const date = new Date(today);
    date.setDate(date.getDate() - i);
    date.setHours(0, 0, 0, 0);
    days.push(date);
  }

  return days;
}

/**
 * Validates phone number format (basic validation)
 */
export function validatePhoneNumber(phone: string): boolean {
  // Basic validation: should start with + and have at least 10 digits
  const phoneRegex = /^\+\d{10,}$/;
  return phoneRegex.test(phone);
}

/**
 * Validates password strength
 */
export function validatePassword(password: string): {
  isValid: boolean;
  message?: string;
} {
  if (password.length < 8) {
    return {
      isValid: false,
      message: 'Password must be at least 8 characters long',
    };
  }

  return {isValid: true};
}

/**
 * Validates habit name
 */
export function validateHabitName(name: string): {
  isValid: boolean;
  message?: string;
} {
  if (!name || name.trim().length === 0) {
    return {
      isValid: false,
      message: 'Habit name is required',
    };
  }

  if (name.length > 100) {
    return {
      isValid: false,
      message: 'Habit name must be less than 100 characters',
    };
  }

  return {isValid: true};
}
