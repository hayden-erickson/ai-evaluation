import React, { useState, useEffect } from 'react';
import { habitsAPI, logsAPI } from '../../services/api';
import { calculateStreak, getLast30Days, formatDate, isToday } from '../../utils/streakCalculator';
import LogDetailsModal from '../Modals/LogDetailsModal';

/**
 * Habit Component
 * Displays a single habit with streak tracking and log management
 */
const Habit = ({ habit, onEdit, onDelete, onUpdate }) => {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isLogModalOpen, setIsLogModalOpen] = useState(false);
  const [selectedLog, setSelectedLog] = useState(null);
  const [selectedDate, setSelectedDate] = useState(null);

  // Load logs when component mounts
  useEffect(() => {
    loadLogs();
  }, [habit.id]);

  /**
   * Load habit logs from API
   */
  const loadLogs = async () => {
    try {
      setLoading(true);
      const fetchedLogs = await habitsAPI.getLogs(habit.id);
      setLogs(fetchedLogs || []);
      setError('');
    } catch (err) {
      setError('Failed to load logs');
      console.error('Error loading logs:', err);
    } finally {
      setLoading(false);
    }
  };

  /**
   * Handle creating a new log for today
   */
  const handleNewLog = () => {
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    setSelectedDate(today.toISOString());
    setSelectedLog(null);
    setIsLogModalOpen(true);
  };

  /**
   * Handle clicking on a day in the streak list
   */
  const handleDayClick = (day) => {
    if (day.hasLog && day.logs.length > 0) {
      // Edit existing log
      setSelectedLog(day.logs[0]);
      setSelectedDate(day.dateStr);
      setIsLogModalOpen(true);
    } else {
      // Create new log for this date
      setSelectedLog(null);
      setSelectedDate(day.dateStr);
      setIsLogModalOpen(true);
    }
  };

  /**
   * Handle saving a log (create or update)
   */
  const handleSaveLog = async (logData) => {
    try {
      if (selectedLog) {
        // Update existing log
        await logsAPI.update(selectedLog.id, logData);
      } else {
        // Create new log
        await logsAPI.create(habit.id, logData);
      }
      
      // Reload logs to update UI
      await loadLogs();
      
      // Notify parent component
      if (onUpdate) {
        onUpdate();
      }
    } catch (err) {
      throw new Error(err.message || 'Failed to save log');
    }
  };

  /**
   * Handle deleting the habit
   */
  const handleDelete = async () => {
    if (window.confirm(`Are you sure you want to delete "${habit.name}"? This will also delete all associated logs.`)) {
      try {
        await habitsAPI.delete(habit.id);
        if (onDelete) {
          onDelete(habit.id);
        }
      } catch (err) {
        alert('Failed to delete habit: ' + err.message);
      }
    }
  };

  // Calculate streak and get last 30 days
  const streak = calculateStreak(logs);
  const days = getLast30Days(logs);

  return (
    <div className="habit-card">
      <div className="habit-header">
        <div className="habit-info">
          <h3>{habit.name}</h3>
          {habit.description && <p>{habit.description}</p>}
        </div>
        <div className="habit-actions">
          <button 
            className="btn btn-secondary btn-small" 
            onClick={() => onEdit(habit)}
            title="Edit habit"
          >
            ✏️
          </button>
          <button 
            className="btn btn-danger btn-small" 
            onClick={handleDelete}
            title="Delete habit"
          >
            🗑️
          </button>
        </div>
      </div>

      <div className="streak-badge">
        <span className="streak-number">{streak}</span>
        <span>day streak 🔥</span>
      </div>

      <button 
        className="btn btn-success" 
        onClick={handleNewLog}
        style={{ marginBottom: '16px' }}
      >
        ✅ Log Today
      </button>

      {error && <div className="error-message">{error}</div>}

      {loading ? (
        <div className="loading">Loading logs...</div>
      ) : (
        <div className="days-grid">
          {days.map((day, index) => (
            <div
              key={index}
              className={`day-container ${day.hasLog ? 'has-log' : ''} ${isToday(day.date) ? 'is-today' : ''}`}
              onClick={() => handleDayClick(day)}
              title={`${formatDate(day.date)}${day.hasLog ? ' - Logged' : ' - No log'}`}
            >
              <span className="day-date">{day.date.getDate()}</span>
            </div>
          ))}
        </div>
      )}

      <LogDetailsModal
        isOpen={isLogModalOpen}
        onClose={() => setIsLogModalOpen(false)}
        onSave={handleSaveLog}
        log={selectedLog}
        date={selectedDate}
      />
    </div>
  );
};

export default Habit;
