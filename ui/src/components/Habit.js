import React, { useState, useEffect } from 'react';
import StreakList from './StreakList';
import { useAuth } from '../App';

// This component will display a single habit.
function Habit({ habit, onOpenLogModal, onOpenHabitModal, onDeleteHabit }) {
  const [logs, setLogs] = useState([]);
  const { token } = useAuth();

  useEffect(() => {
    const fetchLogs = async () => {
      const response = await fetch(`/habits/${habit.id}/logs`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (response.ok) {
        const data = await response.json();
        setLogs(data || []);
      } else {
        // Handle error
        console.error('Failed to fetch logs');
      }
    };

    if (token) {
      fetchLogs();
    }
  }, [token, habit.id]);

  return (
    <li className="p-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-medium text-gray-900">{habit.name}</h3>
          <p className="text-sm text-gray-500">{habit.description}</p>
        </div>
        <div className="flex items-center">
          <span className="text-sm font-semibold text-gray-900 mr-4">Streak: {habit.streak}</span>
          <button onClick={() => onOpenLogModal(habit)} className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded mr-2">Log</button>
          <button onClick={() => onOpenHabitModal(habit)} className="bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded mr-2">Edit</button>
          <button onClick={() => onDeleteHabit(habit.id)} className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">Delete</button>
        </div>
      </div>
      <StreakList logs={logs} />
    </li>
  );
}

export default Habit;
