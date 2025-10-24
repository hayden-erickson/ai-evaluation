import React, { useState, useEffect } from 'react';
import { Habit as HabitType, Log } from '../types';
import { logAPI } from '../services/api';
import { calculateStreak, getLastNDays, getLogForDate, formatDate } from '../utils/streakUtils';
import LogDetailsModal from './LogDetailsModal';
import '../styles/Habit.css';

interface HabitProps {
  habit: HabitType;
  onEdit: () => void;
  onDelete: () => void;
  onRefresh: () => void;
}

/**
 * Habit Component
 * Displays a single habit with streak tracking and log management
 */
const Habit: React.FC<HabitProps> = ({ habit, onEdit, onDelete, onRefresh }) => {
  const [logs, setLogs] = useState<Log[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [isLogModalOpen, setIsLogModalOpen] = useState(false);
  const [selectedLog, setSelectedLog] = useState<Log | null>(null);
  const [isExpanded, setIsExpanded] = useState(false);

  // Fetch logs when component mounts or habit changes
  useEffect(() => {
    fetchLogs();
  }, [habit.id]);

  /**
   * Fetch all logs for this habit
   */
  const fetchLogs = async () => {
    setIsLoading(true);
    setError('');
    try {
      const fetchedLogs = await logAPI.getHabitLogs(habit.id);
      setLogs(fetchedLogs);
    } catch (err) {
      setError('Failed to load logs');
      console.error('Error fetching logs:', err);
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Create a new log for today
   */
  const handleNewLog = () => {
    setSelectedLog(null);
    setIsLogModalOpen(true);
  };

  /**
   * Edit an existing log
   */
  const handleEditLog = (log: Log) => {
    setSelectedLog(log);
    setIsLogModalOpen(true);
  };

  /**
   * Save log (create or update)
   */
  const handleSaveLog = async (notes: string) => {
    try {
      if (selectedLog) {
        // Update existing log
        await logAPI.updateLog(selectedLog.id, { notes });
      } else {
        // Create new log
        await logAPI.createLog(habit.id, { notes });
      }
      await fetchLogs();
      onRefresh();
    } catch (err) {
      throw err;
    }
  };

  /**
   * Delete a log
   */
  const handleDeleteLog = async (logId: number) => {
    if (!window.confirm('Are you sure you want to delete this log?')) {
      return;
    }

    try {
      await logAPI.deleteLog(logId);
      await fetchLogs();
      onRefresh();
    } catch (err) {
      setError('Failed to delete log');
      console.error('Error deleting log:', err);
    }
  };

  // Calculate streak data
  const streakData = calculateStreak(logs);
  const last14Days = getLastNDays(14);

  return (
    <div className="habit-card">
      <div className="habit-header">
        <div className="habit-info">
          <h3>{habit.name}</h3>
          {habit.description && (
            <p className="habit-description">{habit.description}</p>
          )}
        </div>
        <div className="habit-actions">
          <button
            className="btn-icon"
            onClick={onEdit}
            title="Edit habit"
            aria-label="Edit habit"
          >
            ✏️
          </button>
          <button
            className="btn-icon"
            onClick={onDelete}
            title="Delete habit"
            aria-label="Delete habit"
          >
            🗑️
          </button>
        </div>
      </div>

      <div className="streak-info">
        <div className="streak-stat">
          <span className="streak-number">{streakData.currentStreak}</span>
          <span className="streak-label">Current Streak</span>
        </div>
        <div className="streak-stat">
          <span className="streak-number">{streakData.longestStreak}</span>
          <span className="streak-label">Longest Streak</span>
        </div>
      </div>

      <button className="btn-new-log" onClick={handleNewLog}>
        + Log Today
      </button>

      {error && <div className="error-message">{error}</div>}

      <div className="streak-list-container">
        <button
          className="streak-list-toggle"
          onClick={() => setIsExpanded(!isExpanded)}
        >
          {isExpanded ? '▼' : '▶'} Last 14 Days
        </button>

        {isExpanded && (
          <div className="streak-list">
            {isLoading ? (
              <div className="loading">Loading logs...</div>
            ) : (
              last14Days.map((date, index) => {
                const log = getLogForDate(logs, date);
                const isToday = index === 0;

                return (
                  <div
                    key={date.toISOString()}
                    className={`day-container ${log ? 'has-log' : ''} ${
                      isToday ? 'today' : ''
                    }`}
                  >
                    <div className="day-date">{formatDate(date)}</div>
                    {log ? (
                      <div
                        className="log-entry"
                        onClick={() => handleEditLog(log)}
                        title="Click to edit"
                      >
                        <div className="log-indicator">✓</div>
                        {log.notes && (
                          <div className="log-notes">{log.notes}</div>
                        )}
                        <button
                          className="btn-delete-log"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDeleteLog(log.id);
                          }}
                          title="Delete log"
                          aria-label="Delete log"
                        >
                          ×
                        </button>
                      </div>
                    ) : (
                      <div className="empty-log">-</div>
                    )}
                  </div>
                );
              })
            )}
          </div>
        )}
      </div>

      <LogDetailsModal
        isOpen={isLogModalOpen}
        log={selectedLog}
        onClose={() => setIsLogModalOpen(false)}
        onSave={handleSaveLog}
      />
    </div>
  );
};

export default Habit;
