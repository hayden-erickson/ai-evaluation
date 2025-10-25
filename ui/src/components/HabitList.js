import React from 'react';
import Habit from './Habit';

// This component will display the list of habits.
function HabitList({ habits, onOpenLogModal, onOpenHabitModal, onDeleteHabit }) {
  return (
    <div className="bg-white shadow overflow-hidden sm:rounded-md">
      <ul className="divide-y divide-gray-200">
        {habits.map((habit) => (
          <Habit key={habit.id} habit={habit} onOpenLogModal={onOpenLogModal} onOpenHabitModal={onOpenHabitModal} onDeleteHabit={onDeleteHabit} />
        ))}
      </ul>
    </div>
  );
}

export default HabitList;
