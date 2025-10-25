import React, { useState, useEffect } from 'react';
import { createLog, updateLog, deleteLog } from '../../services/api';
import './Modals.css';

/**
 * LogDetailsModal component for creating or editing logs
 * @param {number} habitId - Habit ID for creating new log
 * @param {object} log - Log object (null for new log)
 * @param {function} onSave - Callback to notify parent of save
 * @param {function} onClose - Callback to close modal
 */
const LogDetailsModal = ({ habitId, log, onSave, onClose }) => {
  const [formData, setFormData] = useState({
    notes: '',
    duration_seconds: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  /**
   * Initialize form data when modal opens or log changes
   */
  useEffect(() => {
    if (log) {
      setFormData({
        notes: log.notes || '',
        duration_seconds: log.duration_seconds || '',
      });
    }
  }, [log]);

  /**
   * Handle input changes
   */
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  /**
   * Validate form inputs
   * @returns {boolean} True if inputs are valid
   */
  const validateInputs = () => {
    // Validate duration if provided
    if (formData.duration_seconds && formData.duration_seconds < 0) {
      setError('Duration cannot be negative');
      return false;
    }

    return true;
  };

  /**
   * Handle form submission
   */
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    // Validate inputs
    if (!validateInputs()) {
      return;
    }

    setLoading(true);

    try {
      // Prepare data for API
      const dataToSave = {};
      
      // Include notes if provided
      if (formData.notes.trim()) {
        dataToSave.notes = formData.notes.trim();
      }
      
      // Include duration if provided
      if (formData.duration_seconds !== '') {
        dataToSave.duration_seconds = parseInt(formData.duration_seconds, 10);
      }

      if (log) {
        // Update existing log
        await updateLog(log.id, dataToSave);
      } else {
        // Create new log
        await createLog(habitId, dataToSave);
      }

      // Notify parent component
      onSave();
    } catch (err) {
      setError(err.message || 'Failed to save log');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Handle deleting a log
   */
  const handleDelete = async () => {
    // Confirm deletion
    if (!window.confirm('Are you sure you want to delete this log?')) {
      return;
    }

    setLoading(true);

    try {
      await deleteLog(log.id);
      // Notify parent component
      onSave();
    } catch (err) {
      setError(err.message || 'Failed to delete log');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Handle background click to close modal
   */
  const handleBackgroundClick = (e) => {
    if (e.target.classList.contains('modal-overlay')) {
      onClose();
    }
  };

  return (
    <div className="modal-overlay" onClick={handleBackgroundClick}>
      <div className="modal">
        <div className="modal-header">
          <h2>{log ? 'Edit Log' : 'New Log Entry'}</h2>
          <button 
            className="modal-close-button"
            onClick={onClose}
            disabled={loading}
          >
            ✕
          </button>
        </div>

        {/* Display error message if exists */}
        {error && <div className="error-message">{error}</div>}

        <form onSubmit={handleSubmit} className="modal-form">
          <div className="form-group">
            <label htmlFor="notes">Notes</label>
            <textarea
              id="notes"
              name="notes"
              placeholder="How did it go? (optional)"
              value={formData.notes}
              onChange={handleChange}
              disabled={loading}
              rows={4}
              autoFocus
            />
          </div>

          <div className="form-group">
            <label htmlFor="duration_seconds">Duration (seconds)</label>
            <input
              id="duration_seconds"
              name="duration_seconds"
              type="number"
              placeholder="e.g., 1800 for 30 minutes"
              value={formData.duration_seconds}
              onChange={handleChange}
              disabled={loading}
              min="0"
            />
            <small className="text-secondary">
              Optional: How long did you spend?
            </small>
          </div>

          <div className="modal-actions">
            {/* Show delete button only when editing existing log */}
            {log && (
              <button
                type="button"
                className="danger"
                onClick={handleDelete}
                disabled={loading}
              >
                Delete
              </button>
            )}
            <div className="modal-actions-right">
              <button
                type="button"
                className="secondary"
                onClick={onClose}
                disabled={loading}
              >
                Cancel
              </button>
              <button
                type="submit"
                className="primary"
                disabled={loading}
              >
                {loading ? 'Saving...' : 'Save Log'}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default LogDetailsModal;

