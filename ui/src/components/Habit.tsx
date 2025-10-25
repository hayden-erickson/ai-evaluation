import React, { useState, useEffect } from 'react';
import { Habit, Log } from '../types';
import LogDetailsModal from './LogDetailsModal';
import { LogService } from '../services/LogService';
import { calculateStreak } from '../utils/streak';
import LogList from './LogList';
import './Habit.css';

interface HabitProps {
  habit: Habit;
  onEdit: () => void;
  onDelete: () => void;
}

const HabitComponent: React.FC<HabitProps> = ({ habit, onEdit, onDelete }) => {
  const [isLogModalOpen, setIsLogModalOpen] = useState(false);
  const [logs, setLogs] = useState<Log[]>([]);
  const [selectedLog, setSelectedLog] = useState<Log | null>(null);
  const [streak, setStreak] = useState(0);

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        const data = await LogService.getLogs(habit.id);
        setLogs(data);
        setStreak(calculateStreak(data));
      } catch (error) {
        console.error('Failed to fetch logs', error);
      }
    };
    fetchLogs();
  }, [habit.id]);

  const handleSaveLog = (savedLog: Log) => {
    let updatedLogs;
    if (selectedLog) {
      updatedLogs = logs.map((l) => (l.id === savedLog.id ? savedLog : l));
    } else {
      updatedLogs = [...logs, savedLog];
    }
    setLogs(updatedLogs);
    setStreak(calculateStreak(updatedLogs));
    closeLogModal();
  };

  const openLogModal = (log: Log | null) => {
    setSelectedLog(log);
    setIsLogModalOpen(true);
  };

  const closeLogModal = () => {
    setSelectedLog(null);
    setIsLogModalOpen(false);
  };

  return (
    <div className="habit-container">
      <h3>{habit.name}</h3>
      <p>{habit.description}</p>
      <p>Streak: {streak}</p>
      <button onClick={onEdit}>Edit</button>
      <button onClick={onDelete}>Delete</button>
      <button onClick={() => openLogModal(null)}>Add Log</button>
      <LogList logs={logs} onLogClick={openLogModal} />
      {isLogModalOpen && (
        <LogDetailsModal
          log={selectedLog}
          habitId={habit.id}
          onClose={closeLogModal}
          onSave={handleSaveLog}
        />
      )}
    </div>
  );
};

export default HabitComponent;
