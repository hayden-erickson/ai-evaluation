import React, { useState, useEffect } from 'react';
import { getHabits, createHabit, deleteHabit, updateHabit } from '../../services/api';
import Habit from './Habit';
import HabitDetailsModal from '../Modals/HabitDetailsModal';
import './Habits.css';

/**
 * HabitList component displays all habits for the authenticated user
 * @param {object} user - Current user object
 */
const HabitList = ({ user }) => {
  const [habits, setHabits] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showHabitModal, setShowHabitModal] = useState(false);
  const [editingHabit, setEditingHabit] = useState(null);

  /**
   * Fetch all habits on component mount
   */
  useEffect(() => {
    fetchHabits();
  }, []);

  /**
   * Fetch habits from API
   */
  const fetchHabits = async () => {
    try {
      setLoading(true);
      const data = await getHabits();
      setHabits(data || []);
      setError('');
    } catch (err) {
      setError(err.message || 'Failed to load habits');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Handle creating a new habit
   */
  const handleAddNewHabit = () => {
    setEditingHabit(null);
    setShowHabitModal(true);
  };

  /**
   * Handle editing an existing habit
   * @param {object} habit - Habit to edit
   */
  const handleEditHabit = (habit) => {
    setEditingHabit(habit);
    setShowHabitModal(true);
  };

  /**
   * Handle saving habit (create or update)
   * @param {object} habitData - Habit data to save
   */
  const handleSaveHabit = async (habitData) => {
    try {
      if (editingHabit) {
        // Update existing habit
        const updated = await updateHabit(editingHabit.id, habitData);
        setHabits(habits.map(h => h.id === updated.id ? updated : h));
      } else {
        // Create new habit
        const newHabit = await createHabit(habitData);
        setHabits([...habits, newHabit]);
      }
      setShowHabitModal(false);
      setEditingHabit(null);
    } catch (err) {
      throw err; // Let modal handle the error display
    }
  };

  /**
   * Handle deleting a habit
   * @param {number} habitId - ID of habit to delete
   */
  const handleDeleteHabit = async (habitId) => {
    // Confirm deletion
    if (!window.confirm('Are you sure you want to delete this habit? This action cannot be undone.')) {
      return;
    }

    try {
      await deleteHabit(habitId);
      setHabits(habits.filter(h => h.id !== habitId));
    } catch (err) {
      setError(err.message || 'Failed to delete habit');
    }
  };

  /**
   * Handle closing modal
   */
  const handleCloseModal = () => {
    setShowHabitModal(false);
    setEditingHabit(null);
  };

  // Show loading spinner while fetching
  if (loading) {
    return (
      <div className="habits-container">
        <div className="spinner"></div>
      </div>
    );
  }

  return (
    <div className="habits-container">
      <div className="habits-header">
        <div>
          <h1>My Habits</h1>
          <p className="text-secondary">Track your daily habits and build lasting streaks</p>
        </div>
        <button 
          className="primary"
          onClick={handleAddNewHabit}
        >
          + Add New Habit
        </button>
      </div>

      {/* Display error message if exists */}
      {error && <div className="error-message">{error}</div>}

      {/* Display habits or empty state */}
      {habits.length === 0 ? (
        <div className="empty-state">
          <h2>No habits yet</h2>
          <p className="text-secondary">
            Start building better routines by creating your first habit
          </p>
          <button 
            className="primary mt-2"
            onClick={handleAddNewHabit}
          >
            Create Your First Habit
          </button>
        </div>
      ) : (
        <div className="habits-grid">
          {habits.map(habit => (
            <Habit
              key={habit.id}
              habit={habit}
              onEdit={handleEditHabit}
              onDelete={handleDeleteHabit}
            />
          ))}
        </div>
      )}

      {/* Habit Details Modal */}
      {showHabitModal && (
        <HabitDetailsModal
          habit={editingHabit}
          onSave={handleSaveHabit}
          onClose={handleCloseModal}
        />
      )}
    </div>
  );
};

export default HabitList;

