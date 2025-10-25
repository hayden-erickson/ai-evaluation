/**
 * Streak Calculator Utility
 * Calculates habit streaks based on logs with 1-day skip allowance
 */

/**
 * Calculate the current streak for a habit based on its logs
 * Users are allowed to skip one day before their streak resets
 * @param {Array} logs - Array of log objects with created_at timestamps
 * @returns {number} - Current streak count
 */
export const calculateStreak = (logs) => {
  if (!logs || logs.length === 0) {
    return 0;
  }

  // Sort logs by date (most recent first)
  const sortedLogs = [...logs].sort((a, b) => 
    new Date(b.created_at) - new Date(a.created_at)
  );

  // Get dates only (without time) for comparison
  const logDates = sortedLogs.map(log => {
    const date = new Date(log.created_at);
    return new Date(date.getFullYear(), date.getMonth(), date.getDate());
  });

  // Remove duplicate dates (multiple logs on same day)
  const uniqueDates = [];
  const dateStrings = new Set();
  for (const date of logDates) {
    const dateStr = date.toISOString().split('T')[0];
    if (!dateStrings.has(dateStr)) {
      dateStrings.add(dateStr);
      uniqueDates.push(date);
    }
  }

  if (uniqueDates.length === 0) {
    return 0;
  }

  const today = new Date();
  today.setHours(0, 0, 0, 0);
  
  const mostRecentLog = uniqueDates[0];
  const daysSinceLastLog = Math.floor((today - mostRecentLog) / (1000 * 60 * 60 * 24));

  // If last log is more than 2 days ago (allowing 1 skip day), streak is broken
  if (daysSinceLastLog > 2) {
    return 0;
  }

  // Count consecutive days (allowing 1 skip)
  let streak = 1;
  let skipsUsed = 0;
  
  for (let i = 0; i < uniqueDates.length - 1; i++) {
    const currentDate = uniqueDates[i];
    const nextDate = uniqueDates[i + 1];
    const daysDiff = Math.floor((currentDate - nextDate) / (1000 * 60 * 60 * 24));

    if (daysDiff === 1) {
      // Consecutive day
      streak++;
    } else if (daysDiff === 2 && skipsUsed === 0) {
      // One day skipped, allowed once
      streak++;
      skipsUsed++;
    } else {
      // Streak broken
      break;
    }
  }

  return streak;
};

/**
 * Get the last 30 days with log status
 * @param {Array} logs - Array of log objects
 * @returns {Array} - Array of day objects with date and logs
 */
export const getLast30Days = (logs) => {
  const days = [];
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  // Create a map of dates to logs
  const logsByDate = new Map();
  if (logs) {
    logs.forEach(log => {
      const logDate = new Date(log.created_at);
      logDate.setHours(0, 0, 0, 0);
      const dateStr = logDate.toISOString().split('T')[0];
      
      if (!logsByDate.has(dateStr)) {
        logsByDate.set(dateStr, []);
      }
      logsByDate.get(dateStr).push(log);
    });
  }

  // Generate last 30 days
  for (let i = 29; i >= 0; i--) {
    const date = new Date(today);
    date.setDate(date.getDate() - i);
    const dateStr = date.toISOString().split('T')[0];
    
    days.push({
      date: new Date(date),
      dateStr,
      logs: logsByDate.get(dateStr) || [],
      hasLog: logsByDate.has(dateStr),
    });
  }

  return days;
};

/**
 * Format date for display
 * @param {Date} date - Date object
 * @returns {string} - Formatted date string
 */
export const formatDate = (date) => {
  const options = { month: 'short', day: 'numeric' };
  return date.toLocaleDateString('en-US', options);
};

/**
 * Check if a date is today
 * @param {Date} date - Date to check
 * @returns {boolean} - True if date is today
 */
export const isToday = (date) => {
  const today = new Date();
  return date.getDate() === today.getDate() &&
         date.getMonth() === today.getMonth() &&
         date.getFullYear() === today.getFullYear();
};
