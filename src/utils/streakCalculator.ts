import {
  startOfDay,
  differenceInDays,
  parseISO,
  format,
  subDays,
  isToday,
  isYesterday,
} from 'date-fns';
import {Log} from '../services/api';

/**
 * Calculate the current streak for a habit based on its logs
 * Users are allowed to skip one day before their streak resets
 * @param logs - Array of log entries for the habit
 * @returns Current streak count
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

  // Get unique days (in case there are multiple logs per day)
  const uniqueDays = new Set(
    sortedLogs.map(log => format(parseISO(log.created_at), 'yyyy-MM-dd')),
  );

  const uniqueDaysArray = Array.from(uniqueDays).sort((a, b) =>
    b.localeCompare(a),
  );

  if (uniqueDaysArray.length === 0) {
    return 0;
  }

  const today = startOfDay(new Date());
  const mostRecentLogDate = startOfDay(parseISO(uniqueDaysArray[0]));

  // Check if the most recent log is today or yesterday (or the day before - one skip allowed)
  const daysSinceLastLog = differenceInDays(today, mostRecentLogDate);

  // If more than 2 days have passed (allowing 1 skip day), streak is broken
  if (daysSinceLastLog > 2) {
    return 0;
  }

  // Count consecutive days (allowing one skip)
  let streak = 0;
  let expectedDate = today;
  let skipsUsed = 0;

  for (const dayStr of uniqueDaysArray) {
    const logDate = startOfDay(parseISO(dayStr));
    const daysDiff = differenceInDays(expectedDate, logDate);

    if (daysDiff === 0) {
      // Log is on the expected date
      streak++;
      expectedDate = subDays(expectedDate, 1);
    } else if (daysDiff === 1 && skipsUsed === 0) {
      // One day was skipped, but we allow one skip
      skipsUsed++;
      streak++;
      expectedDate = subDays(logDate, 1);
    } else if (daysDiff === 1 && skipsUsed > 0) {
      // Already used the skip, streak ends
      break;
    } else if (daysDiff > 1) {
      // More than one day gap, streak ends
      break;
    }
  }

  return streak;
};

/**
 * Get the last 30 days with their log status
 * @param logs - Array of log entries for the habit
 * @returns Array of day objects with date and log information
 */
export const getLast30Days = (
  logs: Log[],
): Array<{date: Date; dateStr: string; log?: Log; hasLog: boolean}> => {
  const days: Array<{
    date: Date;
    dateStr: string;
    log?: Log;
    hasLog: boolean;
  }> = [];
  const today = startOfDay(new Date());

  // Create a map of logs by date
  const logsByDate = new Map<string, Log>();
  logs.forEach(log => {
    const dateStr = format(parseISO(log.created_at), 'yyyy-MM-dd');
    // Keep the first log for each day (they're already sorted)
    if (!logsByDate.has(dateStr)) {
      logsByDate.set(dateStr, log);
    }
  });

  // Generate last 30 days
  for (let i = 0; i < 30; i++) {
    const date = subDays(today, i);
    const dateStr = format(date, 'yyyy-MM-dd');
    const log = logsByDate.get(dateStr);

    days.push({
      date,
      dateStr,
      log,
      hasLog: !!log,
    });
  }

  return days;
};

/**
 * Check if a habit has been logged today
 * @param logs - Array of log entries for the habit
 * @returns True if logged today, false otherwise
 */
export const hasLoggedToday = (logs: Log[]): boolean => {
  return logs.some(log => isToday(parseISO(log.created_at)));
};

/**
 * Check if a habit was logged yesterday
 * @param logs - Array of log entries for the habit
 * @returns True if logged yesterday, false otherwise
 */
export const hasLoggedYesterday = (logs: Log[]): boolean => {
  return logs.some(log => isYesterday(parseISO(log.created_at)));
};

/**
 * Get the log for today if it exists
 * @param logs - Array of log entries for the habit
 * @returns Today's log or undefined
 */
export const getTodayLog = (logs: Log[]): Log | undefined => {
  return logs.find(log => isToday(parseISO(log.created_at)));
};

/**
 * Format a date for display
 * @param date - Date to format
 * @returns Formatted date string
 */
export const formatDisplayDate = (date: Date | string): string => {
  const dateObj = typeof date === 'string' ? parseISO(date) : date;

  if (isToday(dateObj)) {
    return 'Today';
  }
  if (isYesterday(dateObj)) {
    return 'Yesterday';
  }

  return format(dateObj, 'MMM d');
};
