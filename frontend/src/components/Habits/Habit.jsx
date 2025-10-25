import React, { useState, useEffect } from 'react';
import { getHabitLogs, createLog } from '../../services/api';
import { calculateStreak, generateDayList } from '../../utils/streakUtils';
import StreakCount from './StreakCount';
import DayContainer from './DayContainer';
import LogDetailsModal from '../Modals/LogDetailsModal';
import './Habits.css';

/**
 * Habit component displays a single habit with its logs and streak
 * @param {object} habit - Habit object
 * @param {function} onEdit - Callback to edit habit
 * @param {function} onDelete - Callback to delete habit
 */
const Habit = ({ habit, onEdit, onDelete }) => {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showLogModal, setShowLogModal] = useState(false);
  const [editingLog, setEditingLog] = useState(null);
  const [isExpanded, setIsExpanded] = useState(false);

  /**
   * Fetch logs when component mounts or habit changes
   */
  useEffect(() => {
    fetchLogs();
  }, [habit.id]);

  /**
   * Fetch logs for this habit from API
   */
  const fetchLogs = async () => {
    try {
      setLoading(false);
      const data = await getHabitLogs(habit.id);
      setLogs(data || []);
      setError('');
    } catch (err) {
      setError(err.message || 'Failed to load logs');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Handle creating a new log for today
   */
  const handleNewLog = () => {
    setEditingLog(null);
    setShowLogModal(true);
  };

  /**
   * Handle editing an existing log
   * @param {object} log - Log to edit
   */
  const handleEditLog = (log) => {
    setEditingLog(log);
    setShowLogModal(true);
  };

  /**
   * Handle log modal save (handled by modal component)
   * Refresh logs after save
   */
  const handleLogSaved = () => {
    setShowLogModal(false);
    setEditingLog(null);
    fetchLogs(); // Refresh logs
  };

  /**
   * Handle closing modal
   */
  const handleCloseModal = () => {
    setShowLogModal(false);
    setEditingLog(null);
  };

  // Calculate current streak
  const currentStreak = calculateStreak(logs);

  // Generate day list for the last 30 days
  const dayList = generateDayList(30, logs);

  return (
    <div className="habit-card card">
      <div className="habit-header">
        <div className="habit-info">
          <h3 className="habit-name">{habit.name}</h3>
          {habit.description && (
            <p className="habit-description text-secondary">{habit.description}</p>
          )}
        </div>
        
        <div className="habit-actions">
          <button
            className="small secondary"
            onClick={() => onEdit(habit)}
            title="Edit habit"
          >
            ✏️ Edit
          </button>
          <button
            className="small danger"
            onClick={() => onDelete(habit.id)}
            title="Delete habit"
          >
            🗑️ Delete
          </button>
        </div>
      </div>

      {/* Streak Count */}
      <StreakCount streak={currentStreak} />

      {/* New Log Button */}
      <button
        className="success new-log-button"
        onClick={handleNewLog}
      >
        ✓ Log for Today
      </button>

      {/* Error message */}
      {error && <div className="error-message mt-2">{error}</div>}

      {/* Streak List Toggle */}
      <button
        className="small secondary expand-button"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        {isExpanded ? '▲ Hide History' : '▼ Show History'}
      </button>

      {/* Streak List - shows last 30 days */}
      {isExpanded && (
        <div className="streak-list">
          <h4 className="streak-list-title">Last 30 Days</h4>
          {loading ? (
            <div className="spinner"></div>
          ) : (
            <div className="day-grid">
              {dayList.map((day, index) => (
                <DayContainer
                  key={index}
                  day={day}
                  onLogClick={handleEditLog}
                />
              ))}
            </div>
          )}
        </div>
      )}

      {/* Log Details Modal */}
      {showLogModal && (
        <LogDetailsModal
          habitId={habit.id}
          log={editingLog}
          onSave={handleLogSaved}
          onClose={handleCloseModal}
        />
      )}
    </div>
  );
};

export default Habit;

