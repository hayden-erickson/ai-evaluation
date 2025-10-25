import React, { useState, useEffect } from 'react';
import { Habit } from '../../types';
import { HabitService } from '../../services/HabitService';

interface HabitDetailsModalProps {
  habit: Habit | null;
  onClose: () => void;
  onSave: (habit: Habit) => void;
}

const HabitDetailsModal: React.FC<HabitDetailsModalProps> = ({ habit, onClose, onSave }) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  useEffect(() => {
    if (habit) {
      setName(habit.name);
      setDescription(habit.description);
    } else {
      setName('');
      setDescription('');
    }
  }, [habit]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (habit) {
        const updatedHabit = await HabitService.updateHabit({ ...habit, name, description });
        onSave(updatedHabit);
      } else {
        const newHabit = await HabitService.createHabit({ name, description });
        onSave(newHabit);
      }
      onClose();
    } catch (error) {
      console.error('Failed to save habit', error);
      // TODO: Show error to user
    }
  };

  return (
    <div className="modal">
      <div className="modal-content">
        <h2>{habit ? 'Edit Habit' : 'Add Habit'}</h2>
        <form onSubmit={handleSubmit}>
          <div>
            <label htmlFor="name">Name</label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>
          <div>
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>
          <button type="submit">Save</button>
          <button type="button" onClick={onClose}>Cancel</button>
        </form>
      </div>
    </div>
  );
};

export default HabitDetailsModal;
