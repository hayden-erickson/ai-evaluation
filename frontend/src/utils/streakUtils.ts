import { Log, StreakData } from '../types';

/**
 * Helper function to check if two dates are on the same day
 */
export const isSameDay = (date1: Date, date2: Date): boolean => {
  return (
    date1.getFullYear() === date2.getFullYear() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getDate() === date2.getDate()
  );
};

/**
 * Helper function to get the difference in days between two dates
 */
export const getDaysDifference = (date1: Date, date2: Date): number => {
  const msPerDay = 1000 * 60 * 60 * 24;
  const utc1 = Date.UTC(date1.getFullYear(), date1.getMonth(), date1.getDate());
  const utc2 = Date.UTC(date2.getFullYear(), date2.getMonth(), date2.getDate());
  return Math.floor((utc2 - utc1) / msPerDay);
};

/**
 * Calculate streak data from logs
 * A streak continues if there's a log for consecutive days, allowing one day skip
 */
export const calculateStreak = (logs: Log[]): StreakData => {
  if (logs.length === 0) {
    return {
      currentStreak: 0,
      longestStreak: 0,
      lastLogDate: null,
    };
  }

  // Sort logs by date (most recent first)
  const sortedLogs = [...logs].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );

  // Get unique dates (one log per day maximum)
  const uniqueDates = sortedLogs
    .map((log) => {
      const date = new Date(log.created_at);
      return new Date(date.getFullYear(), date.getMonth(), date.getDate()).getTime();
    })
    .filter((value, index, self) => self.indexOf(value) === index)
    .sort((a, b) => b - a); // Most recent first

  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const todayTime = today.getTime();

  // Calculate current streak
  let currentStreak = 0;
  let lastDate = todayTime;
  let allowedSkip = true;

  for (const dateTime of uniqueDates) {
    const daysDiff = getDaysDifference(new Date(dateTime), new Date(lastDate));

    // Check if current streak continues
    if (daysDiff === 0) {
      // Same day (can happen if lastDate is today but no log for today yet)
      currentStreak = 1;
      lastDate = dateTime;
    } else if (daysDiff === 1) {
      // Consecutive day
      currentStreak++;
      lastDate = dateTime;
      allowedSkip = true;
    } else if (daysDiff === 2 && allowedSkip) {
      // One day skip allowed
      currentStreak++;
      lastDate = dateTime;
      allowedSkip = false;
    } else {
      // Streak broken
      break;
    }
  }

  // Calculate longest streak
  let longestStreak = 0;
  let tempStreak = 0;
  let tempLastDate = uniqueDates[0];
  let tempAllowedSkip = true;

  for (let i = 0; i < uniqueDates.length; i++) {
    const currentDate = uniqueDates[i];
    
    if (i === 0) {
      tempStreak = 1;
      tempLastDate = currentDate;
      continue;
    }

    const daysDiff = getDaysDifference(new Date(currentDate), new Date(tempLastDate));

    if (daysDiff === 1) {
      // Consecutive day
      tempStreak++;
      tempAllowedSkip = true;
    } else if (daysDiff === 2 && tempAllowedSkip) {
      // One day skip allowed
      tempStreak++;
      tempAllowedSkip = false;
    } else {
      // Streak broken
      longestStreak = Math.max(longestStreak, tempStreak);
      tempStreak = 1;
      tempAllowedSkip = true;
    }

    tempLastDate = currentDate;
  }
  longestStreak = Math.max(longestStreak, tempStreak);

  return {
    currentStreak,
    longestStreak,
    lastLogDate: sortedLogs[0].created_at,
  };
};

/**
 * Get last N days for displaying in the streak list
 */
export const getLastNDays = (n: number): Date[] => {
  const days: Date[] = [];
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  for (let i = 0; i < n; i++) {
    const date = new Date(today);
    date.setDate(today.getDate() - i);
    days.push(date);
  }

  return days;
};

/**
 * Check if there's a log for a specific date
 */
export const getLogForDate = (logs: Log[], date: Date): Log | undefined => {
  return logs.find((log) => {
    const logDate = new Date(log.created_at);
    return isSameDay(logDate, date);
  });
};

/**
 * Format date for display
 */
export const formatDate = (date: Date): string => {
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
  });
};

/**
 * Format date for API (ISO string)
 */
export const formatDateForAPI = (date: Date): string => {
  return date.toISOString();
};
