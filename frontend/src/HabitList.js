import React, { useState, useEffect } from 'react';
import Habit from './Habit';
import HabitDetailsModal from './HabitDetailsModal';
import { habitsAPI } from './api';

/**
 * HabitList component displaying all user habits
 */
function HabitList({ user, onLogout }) {
  const [habits, setHabits] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showHabitModal, setShowHabitModal] = useState(false);
  const [editingHabit, setEditingHabit] = useState(null);

  /**
   * Load all habits on component mount
   */
  useEffect(() => {
    loadHabits();
  }, []);

  /**
   * Load habits from API
   */
  const loadHabits = async () => {
    try {
      setLoading(true);
      const data = await habitsAPI.getAll();
      setHabits(data || []);
      setError('');
    } catch (err) {
      setError('Failed to load habits. Please try again.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  /**
   * Open modal to create new habit
   */
  const handleAddHabit = () => {
    setEditingHabit(null);
    setShowHabitModal(true);
  };

  /**
   * Open modal to edit existing habit
   */
  const handleEditHabit = (habit) => {
    setEditingHabit(habit);
    setShowHabitModal(true);
  };

  /**
   * Save habit (create or update)
   */
  const handleSaveHabit = async (habitData) => {
    try {
      if (editingHabit) {
        // Update existing habit
        await habitsAPI.update(editingHabit.id, habitData);
      } else {
        // Create new habit
        await habitsAPI.create(habitData);
      }
      await loadHabits();
    } catch (err) {
      throw err;
    }
  };

  /**
   * Delete a habit
   */
  const handleDeleteHabit = async (habitId) => {
    // Confirm deletion
    if (!window.confirm('Are you sure you want to delete this habit? All logs will be deleted as well.')) {
      return;
    }

    try {
      await habitsAPI.delete(habitId);
      await loadHabits();
    } catch (err) {
      alert('Failed to delete habit. Please try again.');
      console.error(err);
    }
  };

  return (
    <div className="container">
      <div className="header">
        <h1>Habit Tracker 📝</h1>
        <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
          <span>Welcome, {user?.name || 'User'}!</span>
          <button className="btn btn-small" onClick={onLogout}>
            Logout
          </button>
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}

      <div className="add-habit-container">
        <button className="btn btn-primary add-habit-btn" onClick={handleAddHabit}>
          + Add New Habit
        </button>
      </div>

      {loading ? (
        <div className="loading">
          <div className="spinner"></div>
        </div>
      ) : habits.length === 0 ? (
        <div className="empty-state">
          <h3>No habits yet</h3>
          <p>Start tracking your first habit by clicking the button above!</p>
        </div>
      ) : (
        <div className="habits-container">
          {habits.map((habit) => (
            <Habit
              key={habit.id}
              habit={habit}
              onDelete={handleDeleteHabit}
              onEdit={handleEditHabit}
              onRefresh={loadHabits}
            />
          ))}
        </div>
      )}

      {showHabitModal && (
        <HabitDetailsModal
          habit={editingHabit}
          onClose={() => setShowHabitModal(false)}
          onSave={handleSaveHabit}
        />
      )}
    </div>
  );
}

export default HabitList;
