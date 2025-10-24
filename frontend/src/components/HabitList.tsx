import React, { useState, useEffect } from 'react';
import { Habit as HabitType } from '../types';
import { habitAPI } from '../services/api';
import Habit from './Habit';
import HabitDetailsModal from './HabitDetailsModal';
import '../styles/HabitList.css';

/**
 * HabitList Component
 * Displays list of all user habits with add functionality
 */
const HabitList: React.FC = () => {
  const [habits, setHabits] = useState<HabitType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState<HabitType | null>(null);

  // Fetch habits when component mounts
  useEffect(() => {
    fetchHabits();
  }, []);

  /**
   * Fetch all habits for the authenticated user
   */
  const fetchHabits = async () => {
    setIsLoading(true);
    setError('');
    try {
      const fetchedHabits = await habitAPI.getUserHabits();
      setHabits(fetchedHabits);
    } catch (err) {
      setError('Failed to load habits. Please try again.');
      console.error('Error fetching habits:', err);
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Open modal to create a new habit
   */
  const handleAddHabit = () => {
    setSelectedHabit(null);
    setIsModalOpen(true);
  };

  /**
   * Open modal to edit an existing habit
   */
  const handleEditHabit = (habit: HabitType) => {
    setSelectedHabit(habit);
    setIsModalOpen(true);
  };

  /**
   * Save habit (create or update)
   */
  const handleSaveHabit = async (name: string, description: string) => {
    try {
      if (selectedHabit) {
        // Update existing habit
        await habitAPI.updateHabit(selectedHabit.id, { name, description });
      } else {
        // Create new habit
        await habitAPI.createHabit({ name, description });
      }
      await fetchHabits();
    } catch (err) {
      throw err;
    }
  };

  /**
   * Delete a habit
   */
  const handleDeleteHabit = async (habitId: number) => {
    if (!window.confirm('Are you sure you want to delete this habit? All logs will be deleted as well.')) {
      return;
    }

    try {
      await habitAPI.deleteHabit(habitId);
      await fetchHabits();
    } catch (err) {
      setError('Failed to delete habit');
      console.error('Error deleting habit:', err);
    }
  };

  if (isLoading) {
    return (
      <div className="habit-list-container">
        <div className="loading-container">
          <div className="loading-spinner"></div>
          <p>Loading your habits...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="habit-list-container">
      <div className="habit-list-header">
        <h2>My Habits</h2>
        <button className="btn-add-habit" onClick={handleAddHabit}>
          + Add Habit
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {habits.length === 0 ? (
        <div className="empty-state">
          <p>No habits yet. Start by creating your first habit!</p>
          <button className="btn-primary" onClick={handleAddHabit}>
            Create Your First Habit
          </button>
        </div>
      ) : (
        <div className="habit-list">
          {habits.map((habit) => (
            <Habit
              key={habit.id}
              habit={habit}
              onEdit={() => handleEditHabit(habit)}
              onDelete={() => handleDeleteHabit(habit.id)}
              onRefresh={fetchHabits}
            />
          ))}
        </div>
      )}

      <HabitDetailsModal
        isOpen={isModalOpen}
        habit={selectedHabit}
        onClose={() => setIsModalOpen(false)}
        onSave={handleSaveHabit}
      />
    </div>
  );
};

export default HabitList;
