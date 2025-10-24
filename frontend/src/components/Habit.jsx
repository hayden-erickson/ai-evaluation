import { useState } from 'react';
import { calculateStreak, generateDateRange, getLogForDate, formatDate, isToday } from '../utils/streakUtils';
import LogDetailsModal from './LogDetailsModal';

/**
 * Individual habit component with streak tracking and log management
 */
const Habit = ({ habit, logs, onUpdate, onDelete, onCreateLog, onUpdateLog, onDeleteLog }) => {
  const [showHabitModal, setShowHabitModal] = useState(false);
  const [showLogModal, setShowLogModal] = useState(false);
  const [selectedLog, setSelectedLog] = useState(null);
  const [selectedDate, setSelectedDate] = useState(null);
  const [error, setError] = useState('');

  const streak = calculateStreak(logs);
  const dateRange = generateDateRange(30);

  /**
   * Handle creating a log for today
   */
  const handleNewLog = () => {
    const today = new Date();
    const existingLog = getLogForDate(logs, today);
    
    if (existingLog) {
      // Open log for editing if it exists
      setSelectedLog(existingLog);
      setSelectedDate(formatDate(today));
    } else {
      // Create new log for today
      setSelectedLog(null);
      setSelectedDate(formatDate(today));
    }
    setShowLogModal(true);
  };

  /**
   * Handle clicking on a day in the streak calendar
   */
  const handleDayClick = (date) => {
    const log = getLogForDate(logs, date);
    setSelectedLog(log);
    setSelectedDate(formatDate(date));
    setShowLogModal(true);
  };

  /**
   * Handle saving a log
   */
  const handleSaveLog = async (logData) => {
    setError('');
    try {
      if (selectedLog) {
        // Update existing log
        await onUpdateLog(selectedLog.id, logData);
      } else {
        // Create new log
        await onCreateLog(habit.id, logData);
      }
      setShowLogModal(false);
      setSelectedLog(null);
      setSelectedDate(null);
    } catch (err) {
      throw err; // Let the modal handle the error display
    }
  };

  /**
   * Handle deleting the habit
   */
  const handleDelete = async () => {
    if (window.confirm(`Are you sure you want to delete "${habit.name}"? This will also delete all associated logs.`)) {
      setError('');
      try {
        await onDelete(habit.id);
      } catch (err) {
        setError(err.message || 'Failed to delete habit');
      }
    }
  };

  return (
    <div className="habit-card">
      {error && (
        <div className="error-message">
          <span>⚠️</span>
          <span>{error}</span>
        </div>
      )}

      <div className="habit-header">
        <div className="habit-title-section">
          <h3 className="habit-name">{habit.name}</h3>
          {habit.description && (
            <p className="habit-description">{habit.description}</p>
          )}
        </div>
        <div className="habit-actions">
          <button 
            className="btn btn-secondary btn-sm"
            onClick={() => setShowHabitModal(true)}
            aria-label="Edit habit"
          >
            ✏️ Edit
          </button>
          <button 
            className="btn btn-danger btn-sm"
            onClick={handleDelete}
            aria-label="Delete habit"
          >
            🗑️ Delete
          </button>
        </div>
      </div>

      <div className="streak-section">
        <div className="streak-count">{streak}</div>
        <div className="streak-label">Day Streak 🔥</div>
        <button 
          className="btn btn-success mt-md"
          onClick={handleNewLog}
        >
          ➕ Log Today
        </button>
      </div>

      <div className="streak-list">
        {dateRange.map((date) => {
          const log = getLogForDate(logs, date);
          const isTodayDate = isToday(date);
          
          return (
            <div
              key={date.toISOString()}
              className={`day-container ${log ? 'with-log' : 'empty'} ${isTodayDate ? 'today' : ''}`}
              onClick={() => handleDayClick(date)}
              title={log ? `${formatDate(date)}: ${log.notes || 'No notes'}` : formatDate(date)}
            >
              <div>{log ? '✓' : ''}</div>
              <div className="day-date">{formatDate(date)}</div>
            </div>
          );
        })}
      </div>

      {/* Modals are rendered by parent HabitList component to avoid duplication */}
      <LogDetailsModal
        isOpen={showLogModal}
        onClose={() => {
          setShowLogModal(false);
          setSelectedLog(null);
          setSelectedDate(null);
        }}
        onSave={handleSaveLog}
        log={selectedLog}
        date={selectedDate}
      />
    </div>
  );
};

export default Habit;
