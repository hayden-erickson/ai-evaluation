import React, { useEffect, useState } from 'react';
import { Habit } from '../../types';
import { HabitService } from '../../services/HabitService';
import HabitComponent from './Habit';
import HabitDetailsModal from './HabitDetailsModal';
import './HabitList.css';
import Notification from './Notification';

const HabitList: React.FC = () => {
  const [habits, setHabits] = useState<Habit[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [notification, setNotification] = useState<{ message: string, type: 'success' | 'error' } | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState<Habit | null>(null);

  useEffect(() => {
    const fetchHabits = async () => {
      try {
        const data = await HabitService.getHabits();
        setHabits(data);
      } catch (err) {
        setNotification({ message: 'Failed to fetch habits.', type: 'error' });
      }
    };

    fetchHabits();
  }, []);

  const handleSaveHabit = (savedHabit: Habit) => {
    if (selectedHabit) {
      setHabits(habits.map((h) => (h.id === savedHabit.id ? savedHabit : h)));
    } else {
      setHabits([...habits, savedHabit]);
    }
    setNotification({ message: 'Habit saved successfully!', type: 'success' });
    setIsModalOpen(false);
    setSelectedHabit(null);
  };

  const openModal = (habit: Habit | null) => {
    setSelectedHabit(habit);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setSelectedHabit(null);
    setIsModalOpen(false);
  };

  const handleDeleteHabit = async (id: number) => {
    try {
      await HabitService.deleteHabit(id);
      setHabits(habits.filter((h) => h.id !== id));
      setNotification({ message: 'Habit deleted successfully!', type: 'success' });
    } catch (error) {
      setNotification({ message: 'Failed to delete habit.', type: 'error' });
    }
  };

  return (
    <div className="habit-list-container">
      {notification && <Notification message={notification.message} type={notification.type} />}
      <h2>My Habits</h2>
      <button onClick={() => openModal(null)}>Add New Habit</button>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <div>
        {habits.map((habit) => (
          <HabitComponent
            key={habit.id}
            habit={habit}
            onEdit={() => openModal(habit)}
            onDelete={() => handleDeleteHabit(habit.id)}
          />
        ))}
      </div>
      {isModalOpen && (
        <HabitDetailsModal
          habit={selectedHabit}
          onClose={closeModal}
          onSave={handleSaveHabit}
        />
      )}
    </div>
  );
};

export default HabitList;
