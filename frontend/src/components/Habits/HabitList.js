import React, { useState, useEffect } from 'react';
import { habitsAPI } from '../../services/api';
import Habit from './Habit';
import HabitDetailsModal from '../Modals/HabitDetailsModal';

/**
 * HabitList Component
 * Displays all habits for the current user
 */
const HabitList = () => {
  const [habits, setHabits] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState(null);

  // Load habits when component mounts
  useEffect(() => {
    loadHabits();
  }, []);

  /**
   * Load all habits from API
   */
  const loadHabits = async () => {
    try {
      setLoading(true);
      const fetchedHabits = await habitsAPI.getAll();
      setHabits(fetchedHabits || []);
      setError('');
    } catch (err) {
      setError('Failed to load habits. Please try again.');
      console.error('Error loading habits:', err);
    } finally {
      setLoading(false);
    }
  };

  /**
   * Handle opening modal for new habit
   */
  const handleNewHabit = () => {
    setSelectedHabit(null);
    setIsModalOpen(true);
  };

  /**
   * Handle opening modal for editing habit
   */
  const handleEditHabit = (habit) => {
    setSelectedHabit(habit);
    setIsModalOpen(true);
  };

  /**
   * Handle saving a habit (create or update)
   */
  const handleSaveHabit = async (habitData) => {
    try {
      if (selectedHabit) {
        // Update existing habit
        await habitsAPI.update(selectedHabit.id, habitData);
      } else {
        // Create new habit
        await habitsAPI.create(habitData);
      }
      
      // Reload habits to update UI
      await loadHabits();
    } catch (err) {
      throw new Error(err.message || 'Failed to save habit');
    }
  };

  /**
   * Handle deleting a habit
   */
  const handleDeleteHabit = (habitId) => {
    setHabits(habits.filter(h => h.id !== habitId));
  };

  return (
    <div>
      <div style={{ marginBottom: '24px' }}>
        <button className="btn btn-primary" onClick={handleNewHabit}>
          ➕ Add New Habit
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {loading ? (
        <div className="loading">Loading your habits...</div>
      ) : habits.length === 0 ? (
        <div className="empty-state">
          <h3>No habits yet</h3>
          <p>Start building better habits by creating your first one!</p>
        </div>
      ) : (
        <div className="habits-list">
          {habits.map((habit) => (
            <Habit
              key={habit.id}
              habit={habit}
              onEdit={handleEditHabit}
              onDelete={handleDeleteHabit}
              onUpdate={loadHabits}
            />
          ))}
        </div>
      )}

      <HabitDetailsModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSave={handleSaveHabit}
        habit={selectedHabit}
      />
    </div>
  );
};

export default HabitList;
