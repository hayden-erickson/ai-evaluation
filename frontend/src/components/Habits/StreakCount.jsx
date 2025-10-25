import React from 'react';
import './Habits.css';

/**
 * StreakCount component displays the current streak with a visual indicator
 * @param {number} streak - Current streak count
 */
const StreakCount = ({ streak }) => {
  /**
   * Get emoji based on streak count
   */
  const getStreakEmoji = () => {
    if (streak === 0) return '💤';
    if (streak < 7) return '🔥';
    if (streak < 30) return '⚡';
    if (streak < 100) return '🌟';
    return '👑';
  };

  /**
   * Get color class based on streak count
   */
  const getStreakClass = () => {
    if (streak === 0) return 'streak-none';
    if (streak < 7) return 'streak-low';
    if (streak < 30) return 'streak-medium';
    return 'streak-high';
  };

  return (
    <div className={`streak-count ${getStreakClass()}`}>
      <span className="streak-emoji">{getStreakEmoji()}</span>
      <div className="streak-info">
        <span className="streak-number">{streak}</span>
        <span className="streak-label">day streak</span>
      </div>
      {streak > 0 && (
        <span className="streak-note text-secondary">
          {streak === 1 ? 'Keep it going!' : 'Amazing! Keep it up!'}
        </span>
      )}
    </div>
  );
};

export default StreakCount;

