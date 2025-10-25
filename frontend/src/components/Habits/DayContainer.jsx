import React from 'react';
import { formatDate } from '../../utils/streakUtils';
import './Habits.css';

/**
 * DayContainer component displays a single day with or without a log
 * @param {object} day - Day object with date and log info
 * @param {function} onLogClick - Callback when log is clicked
 */
const DayContainer = ({ day, onLogClick }) => {
  const { date, log, isToday } = day;

  /**
   * Handle click on day container
   */
  const handleClick = () => {
    // Only allow clicking if there's a log
    if (log && onLogClick) {
      onLogClick(log);
    }
  };

  return (
    <div 
      className={`day-container ${log ? 'has-log' : 'no-log'} ${isToday ? 'is-today' : ''}`}
      onClick={handleClick}
      title={log ? 'Click to edit' : 'No log for this day'}
    >
      <div className="day-date">{formatDate(date)}</div>
      <div className="day-indicator">
        {log ? (
          <div className="log-marker">✓</div>
        ) : (
          <div className="empty-marker">·</div>
        )}
      </div>
      {log && log.notes && (
        <div className="day-notes">{log.notes.substring(0, 30)}{log.notes.length > 30 ? '...' : ''}</div>
      )}
    </div>
  );
};

export default DayContainer;

