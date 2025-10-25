import React, { useState, useEffect } from 'react';
import { Habit } from '../types';
import '../styles/Modal.css';

interface HabitDetailsModalProps {
  isOpen: boolean;
  habit: Habit | null;
  onClose: () => void;
  onSave: (name: string, description: string) => Promise<void>;
}

/**
 * HabitDetailsModal Component
 * Modal for creating or editing habits
 */
const HabitDetailsModal: React.FC<HabitDetailsModalProps> = ({
  isOpen,
  habit,
  onClose,
  onSave,
}) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  // Reset form when modal opens with a different habit
  useEffect(() => {
    if (isOpen) {
      setName(habit?.name || '');
      setDescription(habit?.description || '');
      setError('');
    }
  }, [isOpen, habit]);

  /**
   * Handle form submission
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // Validate input
    if (!name.trim()) {
      setError('Habit name is required');
      return;
    }

    setIsLoading(true);
    try {
      await onSave(name.trim(), description.trim());
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save habit');
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Handle modal close
   */
  const handleClose = () => {
    if (!isLoading) {
      onClose();
    }
  };

  // Don't render if modal is closed
  if (!isOpen) {
    return null;
  }

  return (
    <div className="modal-overlay" onClick={handleClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>{habit ? 'Edit Habit' : 'New Habit'}</h2>
          <button className="modal-close" onClick={handleClose} disabled={isLoading}>
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit} className="modal-form">
          <div className="form-group">
            <label htmlFor="habitName">Habit Name *</label>
            <input
              id="habitName"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g., Morning Exercise"
              disabled={isLoading}
              autoFocus
            />
          </div>

          <div className="form-group">
            <label htmlFor="habitDescription">Description</label>
            <textarea
              id="habitDescription"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Add details about your habit..."
              rows={4}
              disabled={isLoading}
            />
          </div>

          {error && <div className="error-message">{error}</div>}

          <div className="modal-actions">
            <button
              type="button"
              onClick={handleClose}
              className="btn-secondary"
              disabled={isLoading}
            >
              Cancel
            </button>
            <button type="submit" className="btn-primary" disabled={isLoading}>
              {isLoading ? 'Saving...' : 'Save'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default HabitDetailsModal;
