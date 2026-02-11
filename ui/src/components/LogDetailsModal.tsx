import React, { useState, useEffect } from 'react';
import { Log } from '../types';
import { LogService } from '../services/LogService';

interface LogDetailsModalProps {
  log: Log | null;
  habitId: number;
  onClose: () => void;
  onSave: (log: Log) => void;
}

const LogDetailsModal: React.FC<LogDetailsModalProps> = ({ log, habitId, onClose, onSave }) => {
  const [notes, setNotes] = useState('');
  const [duration, setDuration] = useState(0);

  useEffect(() => {
    if (log) {
      setNotes(log.notes);
      setDuration(log.duration);
    } else {
      setNotes('');
      setDuration(0);
    }
  }, [log]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const logData = { date: new Date().toISOString(), notes, duration };
    try {
      if (log) {
        const updatedLog = await LogService.updateLog({ ...log, ...logData });
        onSave(updatedLog);
      } else {
        const newLog = await LogService.createLog(habitId, logData);
        onSave(newLog);
      }
      onClose();
    } catch (error) {
      console.error('Failed to save log', error);
      // TODO: Show error to user
    }
  };

  return (
    <div className="modal">
      <div className="modal-content">
        <h2>{log ? 'Edit Log' : 'Add Log'}</h2>
        <form onSubmit={handleSubmit}>
          <div>
            <label htmlFor="notes">Notes</label>
            <textarea
              id="notes"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
            />
          </div>
          <div>
            <label htmlFor="duration">Duration (minutes)</label>
            <input
              type="number"
              id="duration"
              value={duration}
              onChange={(e) => setDuration(parseInt(e.target.value, 10))}
              required
            />
          </div>
          <button type="submit">Save</button>
          <button type="button" onClick={onClose}>Cancel</button>
        </form>
      </div>
    </div>
  );
};

export default LogDetailsModal;
