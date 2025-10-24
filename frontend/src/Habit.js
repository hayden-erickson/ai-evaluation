import React, { useState, useEffect } from 'react';
import LogDetailsModal from './LogDetailsModal';
import { logsAPI } from './api';

/**
 * Calculate streak with 1-day grace period
 * Returns the current streak count
 */
function calculateStreak(logs) {
  if (!logs || logs.length === 0) {
    return 0;
  }

  // Sort logs by date descending (most recent first)
  const sortedLogs = [...logs].sort((a, b) => 
    new Date(b.created_at) - new Date(a.created_at)
  );

  let streak = 0;
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  // Check if there's a log today or yesterday (grace period)
  const mostRecentLog = new Date(sortedLogs[0].created_at);
  mostRecentLog.setHours(0, 0, 0, 0);
  
  const daysDiff = Math.floor((today - mostRecentLog) / (1000 * 60 * 60 * 24));
  
  // If most recent log is more than 1 day old, streak is broken
  if (daysDiff > 1) {
    return 0;
  }

  // Count consecutive days (with grace period)
  let currentDate = new Date(today);
  let lastLogDate = null;

  for (const log of sortedLogs) {
    const logDate = new Date(log.created_at);
    logDate.setHours(0, 0, 0, 0);

    if (lastLogDate === null) {
      // First log
      lastLogDate = logDate;
      streak = 1;
    } else {
      // Calculate days between this log and the last one
      const daysBetween = Math.floor((lastLogDate - logDate) / (1000 * 60 * 60 * 24));
      
      // If gap is more than 1 day (accounting for grace period), streak is broken
      if (daysBetween > 2) {
        break;
      }
      
      // Only increment streak if it's a different day
      if (daysBetween >= 1) {
        streak++;
        lastLogDate = logDate;
      }
    }
  }

  return streak;
}

/**
 * Get last 14 days with log status
 */
function getLast14Days(logs) {
  const days = [];
  const today = new Date();
  
  for (let i = 13; i >= 0; i--) {
    const date = new Date(today);
    date.setDate(date.getDate() - i);
    date.setHours(0, 0, 0, 0);
    
    // Check if there's a log for this day
    const hasLog = logs.some(log => {
      const logDate = new Date(log.created_at);
      logDate.setHours(0, 0, 0, 0);
      return logDate.getTime() === date.getTime();
    });
    
    // Find the log for this day
    const dayLog = logs.find(log => {
      const logDate = new Date(log.created_at);
      logDate.setHours(0, 0, 0, 0);
      return logDate.getTime() === date.getTime();
    });
    
    days.push({
      date,
      hasLog,
      log: dayLog,
    });
  }
  
  return days;
}

/**
 * Format date to short format (e.g., "12/25")
 */
function formatShortDate(date) {
  return `${date.getMonth() + 1}/${date.getDate()}`;
}

/**
 * Check if date is today
 */
function isToday(date) {
  const today = new Date();
  return date.getDate() === today.getDate() &&
         date.getMonth() === today.getMonth() &&
         date.getFullYear() === today.getFullYear();
}

/**
 * Habit component displaying a single habit with streak and logs
 */
function Habit({ habit, onDelete, onEdit, onRefresh }) {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showLogModal, setShowLogModal] = useState(false);
  const [editingLog, setEditingLog] = useState(null);

  /**
   * Load logs for this habit
   */
  useEffect(() => {
    loadLogs();
  }, [habit.id]);

  const loadLogs = async () => {
    try {
      setLoading(true);
      const data = await logsAPI.getByHabit(habit.id);
      setLogs(data || []);
      setError('');
    } catch (err) {
      setError('Failed to load logs');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  /**
   * Create new log for today
   */
  const handleNewLog = () => {
    setEditingLog(null);
    setShowLogModal(true);
  };

  /**
   * Edit existing log
   */
  const handleEditLog = (log) => {
    setEditingLog(log);
    setShowLogModal(true);
  };

  /**
   * Save log (create or update)
   */
  const handleSaveLog = async (logData) => {
    try {
      if (editingLog) {
        // Update existing log
        await logsAPI.update(editingLog.id, logData);
      } else {
        // Create new log
        await logsAPI.create(habit.id, logData);
      }
      await loadLogs();
      onRefresh();
    } catch (err) {
      throw err;
    }
  };

  const streak = calculateStreak(logs);
  const last14Days = getLast14Days(logs);

  return (
    <div className="habit-card">
      <div className="habit-header">
        <div className="habit-info">
          <h3>{habit.name}</h3>
          {habit.description && <p>{habit.description}</p>}
        </div>
        <div className="habit-actions">
          <button 
            className="btn btn-small" 
            onClick={() => onEdit(habit)}
            title="Edit habit"
          >
            Edit
          </button>
          <button 
            className="btn btn-danger btn-small" 
            onClick={() => onDelete(habit.id)}
            title="Delete habit"
          >
            Delete
          </button>
        </div>
      </div>

      <div className="streak-counter">
        <div className="streak-number">{streak}</div>
        <div className="streak-label">Day Streak 🔥</div>
      </div>

      <div className="streak-list">
        {last14Days.map((day, index) => (
          <div
            key={index}
            className={`day-container ${day.hasLog ? 'has-log' : ''}`}
            onClick={() => day.log && handleEditLog(day.log)}
            title={day.hasLog ? 'Click to edit' : 'No log for this day'}
          >
            <div className="day-date">{formatShortDate(day.date)}</div>
            <div className="day-status">
              {day.hasLog ? '✓' : '○'}
            </div>
          </div>
        ))}
      </div>

      <button 
        className="btn btn-success" 
        onClick={handleNewLog}
        style={{ width: '100%' }}
      >
        + Log Today
      </button>

      {showLogModal && (
        <LogDetailsModal
          log={editingLog}
          onClose={() => setShowLogModal(false)}
          onSave={handleSaveLog}
        />
      )}
    </div>
  );
}

export default Habit;
