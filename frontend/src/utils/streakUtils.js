/**
 * Utility functions for calculating habit streaks
 */

/**
 * Check if two dates are on the same day
 * @param {Date} date1 - First date
 * @param {Date} date2 - Second date
 * @returns {boolean} True if dates are on the same day
 */
export const isSameDay = (date1, date2) => {
  return (
    date1.getFullYear() === date2.getFullYear() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getDate() === date2.getDate()
  );
};

/**
 * Get the start of day for a given date
 * @param {Date} date - The date
 * @returns {Date} Start of day
 */
export const getStartOfDay = (date) => {
  const newDate = new Date(date);
  newDate.setHours(0, 0, 0, 0);
  return newDate;
};

/**
 * Get the number of days between two dates
 * @param {Date} date1 - First date
 * @param {Date} date2 - Second date
 * @returns {number} Number of days between dates
 */
export const daysBetween = (date1, date2) => {
  const start = getStartOfDay(date1);
  const end = getStartOfDay(date2);
  const diffTime = end - start;
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
  return Math.abs(diffDays);
};

/**
 * Calculate the current streak for a habit
 * User is allowed to skip one day before streak resets
 * @param {Array} logs - Array of log objects with created_at timestamps
 * @returns {number} Current streak count
 */
export const calculateStreak = (logs) => {
  if (!logs || logs.length === 0) {
    return 0;
  }

  // Sort logs by date (newest first)
  const sortedLogs = [...logs].sort((a, b) => 
    new Date(b.created_at) - new Date(a.created_at)
  );

  // Get unique days (in case multiple logs per day)
  const uniqueDays = [];
  const seenDays = new Set();

  for (const log of sortedLogs) {
    const logDate = getStartOfDay(new Date(log.created_at));
    const dayKey = logDate.toISOString();
    
    if (!seenDays.has(dayKey)) {
      seenDays.add(dayKey);
      uniqueDays.push(logDate);
    }
  }

  const today = getStartOfDay(new Date());
  const mostRecentLog = uniqueDays[0];

  // Check if most recent log is today or yesterday or day before yesterday (with skip day allowance)
  const daysSinceLastLog = daysBetween(today, mostRecentLog);
  
  if (daysSinceLastLog > 2) {
    // Streak is broken if more than 2 days have passed
    return 0;
  }

  // Count consecutive days (allowing one skip day)
  let streak = 0;
  let currentDate = today;
  let skippedDays = 0;

  for (let i = 0; i < 365; i++) { // Max 365 days lookback
    const hasLogToday = uniqueDays.some(logDate => 
      isSameDay(logDate, currentDate)
    );

    if (hasLogToday) {
      streak++;
      skippedDays = 0; // Reset skip counter when log is found
    } else {
      skippedDays++;
      // Allow one skip day
      if (skippedDays > 1) {
        break; // Streak ends after missing more than one day
      }
    }

    // Move to previous day
    currentDate = new Date(currentDate);
    currentDate.setDate(currentDate.getDate() - 1);
  }

  return streak;
};

/**
 * Generate an array of day objects for the last N days
 * @param {number} days - Number of days to generate
 * @param {Array} logs - Array of logs
 * @returns {Array} Array of day objects with date and log info
 */
export const generateDayList = (days, logs) => {
  const dayList = [];
  const today = getStartOfDay(new Date());

  // Create a map of dates to logs for quick lookup
  const logsByDay = new Map();
  
  if (logs) {
    logs.forEach(log => {
      const logDate = getStartOfDay(new Date(log.created_at));
      const dayKey = logDate.toISOString();
      
      // Keep only the most recent log for each day
      if (!logsByDay.has(dayKey)) {
        logsByDay.set(dayKey, log);
      }
    });
  }

  // Generate day objects from most recent to oldest
  for (let i = 0; i < days; i++) {
    const date = new Date(today);
    date.setDate(date.getDate() - i);
    const dayKey = date.toISOString();
    
    dayList.push({
      date: date,
      log: logsByDay.get(dayKey) || null,
      isToday: i === 0,
    });
  }

  return dayList;
};

/**
 * Format a date for display
 * @param {Date} date - The date to format
 * @returns {string} Formatted date string
 */
export const formatDate = (date) => {
  const today = getStartOfDay(new Date());
  const targetDate = getStartOfDay(date);
  
  if (isSameDay(today, targetDate)) {
    return 'Today';
  }
  
  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);
  
  if (isSameDay(yesterday, targetDate)) {
    return 'Yesterday';
  }
  
  return date.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric' 
  });
};

