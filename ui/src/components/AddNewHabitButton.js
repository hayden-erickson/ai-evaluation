import React from 'react';

// This component will be a button to add a new habit.
function AddNewHabitButton({ onOpen }) {
  return (
    <button onClick={onOpen} className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
      Add New Habit
    </button>
  );
}

export default AddNewHabitButton;
