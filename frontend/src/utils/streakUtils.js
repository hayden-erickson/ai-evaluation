/**
 * Calculate the current streak for a habit based on its logs
 * A user can skip one day before the streak resets
 * @param {Array} logs - Array of log objects with created_at timestamps
 * @returns {number} Current streak count
 */
export const calculateStreak = (logs) => {
  if (!logs || logs.length === 0) {
    return 0;
  }

  // Sort logs by date (most recent first)
  const sortedLogs = [...logs].sort((a, b) => 
    new Date(b.created_at) - new Date(a.created_at)
  );

  // Get dates only (no time component)
  const logDates = sortedLogs.map(log => {
    const date = new Date(log.created_at);
    return new Date(date.getFullYear(), date.getMonth(), date.getDate());
  });

  // Remove duplicate dates (multiple logs on same day)
  const uniqueDates = [...new Set(logDates.map(d => d.getTime()))].map(t => new Date(t));

  const today = new Date();
  today.setHours(0, 0, 0, 0);
  
  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);

  const twoDaysAgo = new Date(today);
  twoDaysAgo.setDate(twoDaysAgo.getDate() - 2);

  // Check if most recent log is today, yesterday, or two days ago (allowing 1 skip)
  const mostRecentDate = uniqueDates[0];
  
  if (mostRecentDate < twoDaysAgo) {
    // More than 2 days ago - streak is broken
    return 0;
  }

  let streak = 1;
  let currentDate = new Date(mostRecentDate);
  
  // Walk backwards through dates counting consecutive days (allowing 1 skip)
  for (let i = 1; i < uniqueDates.length; i++) {
    const nextDate = uniqueDates[i];
    const expectedDate = new Date(currentDate);
    expectedDate.setDate(expectedDate.getDate() - 1);
    
    const allowedSkipDate = new Date(currentDate);
    allowedSkipDate.setDate(allowedSkipDate.getDate() - 2);

    if (nextDate.getTime() === expectedDate.getTime()) {
      // Consecutive day
      streak++;
      currentDate = nextDate;
    } else if (nextDate.getTime() === allowedSkipDate.getTime()) {
      // One day skipped (allowed)
      streak += 2; // Count both days
      currentDate = nextDate;
    } else {
      // Gap is too large - streak broken
      break;
    }
  }

  return streak;
};

/**
 * Generate an array of dates for displaying the streak calendar
 * @param {number} days - Number of days to generate (default 30)
 * @returns {Array} Array of date objects
 */
export const generateDateRange = (days = 30) => {
  const dates = [];
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  for (let i = days - 1; i >= 0; i--) {
    const date = new Date(today);
    date.setDate(date.getDate() - i);
    dates.push(date);
  }

  return dates;
};

/**
 * Check if a log exists for a specific date
 * @param {Array} logs - Array of log objects
 * @param {Date} date - Date to check
 * @returns {Object|null} Log object if exists, null otherwise
 */
export const getLogForDate = (logs, date) => {
  if (!logs || logs.length === 0) {
    return null;
  }

  const targetDate = new Date(date);
  targetDate.setHours(0, 0, 0, 0);

  return logs.find(log => {
    const logDate = new Date(log.created_at);
    logDate.setHours(0, 0, 0, 0);
    return logDate.getTime() === targetDate.getTime();
  }) || null;
};

/**
 * Format date to display string
 * @param {Date} date - Date to format
 * @returns {string} Formatted date string
 */
export const formatDate = (date) => {
  const options = { month: 'short', day: 'numeric' };
  return new Date(date).toLocaleDateString('en-US', options);
};

/**
 * Check if a date is today
 * @param {Date} date - Date to check
 * @returns {boolean} True if date is today
 */
export const isToday = (date) => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const checkDate = new Date(date);
  checkDate.setHours(0, 0, 0, 0);
  return today.getTime() === checkDate.getTime();
};
