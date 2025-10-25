import React, { useState, useEffect } from 'react';
import './Modals.css';

/**
 * HabitDetailsModal component for creating or editing habits
 * @param {object} habit - Habit object (null for new habit)
 * @param {function} onSave - Callback to save habit
 * @param {function} onClose - Callback to close modal
 */
const HabitDetailsModal = ({ habit, onSave, onClose }) => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    duration_seconds: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  /**
   * Initialize form data when modal opens or habit changes
   */
  useEffect(() => {
    if (habit) {
      setFormData({
        name: habit.name || '',
        description: habit.description || '',
        duration_seconds: habit.duration_seconds || '',
      });
    }
  }, [habit]);

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
    if (!formData.name.trim()) {
      setError('Name is required');
      return false;
    }
    
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
      
      // Only include fields that have values
      if (formData.name.trim()) {
        dataToSave.name = formData.name.trim();
      }
      if (formData.description.trim()) {
        dataToSave.description = formData.description.trim();
      }
      if (formData.duration_seconds !== '') {
        dataToSave.duration_seconds = parseInt(formData.duration_seconds, 10);
      }

      await onSave(dataToSave);
      // Modal will be closed by parent component on success
    } catch (err) {
      setError(err.message || 'Failed to save habit');
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
          <h2>{habit ? 'Edit Habit' : 'Create New Habit'}</h2>
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
            <label htmlFor="name">Habit Name *</label>
            <input
              id="name"
              name="name"
              type="text"
              placeholder="e.g., Morning Exercise"
              value={formData.name}
              onChange={handleChange}
              disabled={loading}
              autoFocus
            />
          </div>

          <div className="form-group">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              name="description"
              placeholder="Describe your habit (optional)"
              value={formData.description}
              onChange={handleChange}
              disabled={loading}
              rows={3}
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
              Optional: Set a target duration in seconds
            </small>
          </div>

          <div className="modal-actions">
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
              {loading ? 'Saving...' : 'Save Habit'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default HabitDetailsModal;

