import { useState, useEffect } from 'react';
import { habitsAPI, logsAPI } from '../services/api';
import Habit from './Habit';
import HabitDetailsModal from './HabitDetailsModal';

/**
 * Main component for displaying and managing the list of habits
 */
const HabitList = () => {
  const [habits, setHabits] = useState([]);
  const [habitLogs, setHabitLogs] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showHabitModal, setShowHabitModal] = useState(false);
  const [editingHabit, setEditingHabit] = useState(null);

  /**
   * Fetch all habits and their logs on component mount
   */
  useEffect(() => {
    fetchHabits();
  }, []);

  /**
   * Fetch all habits from the API
   */
  const fetchHabits = async () => {
    setLoading(true);
    setError('');
    try {
      const fetchedHabits = await habitsAPI.getAll();
      setHabits(fetchedHabits || []);
      
      // Fetch logs for each habit
      if (fetchedHabits && fetchedHabits.length > 0) {
        await fetchAllLogs(fetchedHabits);
      }
    } catch (err) {
      setError(err.message || 'Failed to load habits');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Fetch logs for all habits
   */
  const fetchAllLogs = async (habitsList) => {
    const logsPromises = habitsList.map(async (habit) => {
      try {
        const logs = await logsAPI.getByHabit(habit.id);
        return { habitId: habit.id, logs: logs || [] };
      } catch (err) {
        console.error(`Failed to fetch logs for habit ${habit.id}:`, err);
        return { habitId: habit.id, logs: [] };
      }
    });

    const logsResults = await Promise.all(logsPromises);
    const logsMap = {};
    logsResults.forEach(({ habitId, logs }) => {
      logsMap[habitId] = logs;
    });
    setHabitLogs(logsMap);
  };

  /**
   * Handle creating a new habit
   */
  const handleCreateHabit = async (habitData) => {
    const newHabit = await habitsAPI.create(habitData);
    setHabits([...habits, newHabit]);
    setHabitLogs({ ...habitLogs, [newHabit.id]: [] });
  };

  /**
   * Handle updating a habit
   */
  const handleUpdateHabit = async (habitId, habitData) => {
    const updatedHabit = await habitsAPI.update(habitId, habitData);
    setHabits(habits.map(h => h.id === habitId ? updatedHabit : h));
  };

  /**
   * Handle deleting a habit
   */
  const handleDeleteHabit = async (habitId) => {
    await habitsAPI.delete(habitId);
    setHabits(habits.filter(h => h.id !== habitId));
    const newHabitLogs = { ...habitLogs };
    delete newHabitLogs[habitId];
    setHabitLogs(newHabitLogs);
  };

  /**
   * Handle creating a new log
   */
  const handleCreateLog = async (habitId, logData) => {
    const newLog = await logsAPI.create(habitId, logData);
    setHabitLogs({
      ...habitLogs,
      [habitId]: [...(habitLogs[habitId] || []), newLog],
    });
  };

  /**
   * Handle updating a log
   */
  const handleUpdateLog = async (logId, logData) => {
    const updatedLog = await logsAPI.update(logId, logData);
    
    // Update the log in the appropriate habit's logs array
    const newHabitLogs = { ...habitLogs };
    Object.keys(newHabitLogs).forEach(habitId => {
      newHabitLogs[habitId] = newHabitLogs[habitId].map(log =>
        log.id === logId ? updatedLog : log
      );
    });
    setHabitLogs(newHabitLogs);
  };

  /**
   * Handle deleting a log
   */
  const handleDeleteLog = async (logId) => {
    await logsAPI.delete(logId);
    
    // Remove the log from the appropriate habit's logs array
    const newHabitLogs = { ...habitLogs };
    Object.keys(newHabitLogs).forEach(habitId => {
      newHabitLogs[habitId] = newHabitLogs[habitId].filter(log => log.id !== logId);
    });
    setHabitLogs(newHabitLogs);
  };

  /**
   * Handle saving a habit (create or update)
   */
  const handleSaveHabit = async (habitData) => {
    if (editingHabit) {
      await handleUpdateHabit(editingHabit.id, habitData);
    } else {
      await handleCreateHabit(habitData);
    }
    setShowHabitModal(false);
    setEditingHabit(null);
  };

  /**
   * Open modal to create a new habit
   */
  const handleNewHabit = () => {
    setEditingHabit(null);
    setShowHabitModal(true);
  };

  if (loading) {
    return (
      <div className="loading">
        <p>Loading your habits...</p>
      </div>
    );
  }

  return (
    <div className="container">
      {error && (
        <div className="error-message">
          <span>⚠️</span>
          <span>{error}</span>
        </div>
      )}

      <div className="flex-between mb-xl">
        <h2 style={{ fontSize: '1.75rem', fontWeight: '600' }}>My Habits</h2>
        <button className="btn btn-primary" onClick={handleNewHabit}>
          ➕ Add New Habit
        </button>
      </div>

      {habits.length === 0 ? (
        <div className="empty-state">
          <div className="empty-state-icon">📝</div>
          <h3 className="empty-state-title">No habits yet</h3>
          <p className="empty-state-description">
            Start tracking your habits to build consistent streaks and achieve your goals!
          </p>
          <button className="btn btn-primary btn-lg" onClick={handleNewHabit}>
            Create Your First Habit
          </button>
        </div>
      ) : (
        <div className="habits-container">
          {habits.map((habit) => (
            <Habit
              key={habit.id}
              habit={habit}
              logs={habitLogs[habit.id] || []}
              onUpdate={handleUpdateHabit}
              onDelete={handleDeleteHabit}
              onCreateLog={handleCreateLog}
              onUpdateLog={handleUpdateLog}
              onDeleteLog={handleDeleteLog}
            />
          ))}
        </div>
      )}

      <HabitDetailsModal
        isOpen={showHabitModal}
        onClose={() => {
          setShowHabitModal(false);
          setEditingHabit(null);
        }}
        onSave={handleSaveHabit}
        habit={editingHabit}
      />
    </div>
  );
};

export default HabitList;
