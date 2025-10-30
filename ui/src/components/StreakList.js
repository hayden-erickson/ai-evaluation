import React from 'react';

// This component will display the streak of a habit.
function StreakList({ logs }) {
  // Create a map of dates to logs for easy lookup
  const logMap = new Map(logs.map(log => [new Date(log.date).toDateString(), log]));

  // Get the last 30 days
  const days = Array.from({ length: 30 }, (_, i) => {
    const date = new Date();
    date.setDate(date.getDate() - i);
    return date;
  }).reverse();

  return (
    <div className="flex space-x-1 mt-2">
      {days.map(day => {
        const log = logMap.get(day.toDateString());
        return (
          <div
            key={day.toISOString()}
            className={`w-4 h-4 rounded-sm ${log ? 'bg-green-500' : 'bg-gray-200'}`}
            title={day.toDateString()}
          ></div>
        );
      })}
    </div>
  );
}

export default StreakList;
