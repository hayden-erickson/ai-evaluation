import React, { useState, useEffect } from 'react';
import { Log } from '../types';
import '../styles/Modal.css';

interface LogDetailsModalProps {
  isOpen: boolean;
  log: Log | null;
  onClose: () => void;
  onSave: (notes: string) => Promise<void>;
}

/**
 * LogDetailsModal Component
 * Modal for creating or editing habit logs
 */
const LogDetailsModal: React.FC<LogDetailsModalProps> = ({
  isOpen,
  log,
  onClose,
  onSave,
}) => {
  const [notes, setNotes] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  // Reset form when modal opens with a different log
  useEffect(() => {
    if (isOpen) {
      setNotes(log?.notes || '');
      setError('');
    }
  }, [isOpen, log]);

  /**
   * Handle form submission
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    setIsLoading(true);
    try {
      await onSave(notes.trim());
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save log');
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
          <h2>{log ? 'Edit Log' : 'New Log'}</h2>
          <button className="modal-close" onClick={handleClose} disabled={isLoading}>
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit} className="modal-form">
          <div className="form-group">
            <label htmlFor="logNotes">Notes</label>
            <textarea
              id="logNotes"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Add notes about today's progress..."
              rows={6}
              disabled={isLoading}
              autoFocus
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

export default LogDetailsModal;
