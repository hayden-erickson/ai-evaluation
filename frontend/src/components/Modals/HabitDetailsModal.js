import React, { useState, useEffect } from 'react';

/**
 * HabitDetailsModal Component
 * Modal for creating or editing a habit
 */
const HabitDetailsModal = ({ isOpen, onClose, onSave, habit }) => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // Populate form when editing existing habit
  useEffect(() => {
    if (habit) {
      setFormData({
        name: habit.name || '',
        description: habit.description || '',
      });
    } else {
      setFormData({
        name: '',
        description: '',
      });
    }
    setError('');
  }, [habit, isOpen]);

  /**
   * Handle input changes
   */
  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
    setError('');
  };

  /**
   * Handle form submission
   */
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      // Validate inputs
      if (!formData.name.trim()) {
        throw new Error('Habit name is required');
      }

      // Call parent save handler
      await onSave(formData);
      
      // Reset form and close modal
      setFormData({ name: '', description: '' });
      onClose();
    } catch (err) {
      setError(err.message || 'Failed to save habit');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Handle modal close
   */
  const handleClose = () => {
    if (!loading) {
      setFormData({ name: '', description: '' });
      setError('');
      onClose();
    }
  };

  // Don't render if modal is not open
  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={handleClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">
            {habit ? 'Edit Habit' : 'New Habit'}
          </h2>
          <button 
            className="modal-close" 
            onClick={handleClose}
            disabled={loading}
          >
            ×
          </button>
        </div>

        {error && <div className="error-message">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label" htmlFor="name">
              Habit Name *
            </label>
            <input
              type="text"
              id="name"
              name="name"
              className="form-input"
              value={formData.name}
              onChange={handleChange}
              placeholder="e.g., Morning Exercise"
              disabled={loading}
              autoFocus
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="description">
              Description
            </label>
            <textarea
              id="description"
              name="description"
              className="form-textarea"
              value={formData.description}
              onChange={handleChange}
              placeholder="Optional description of your habit"
              disabled={loading}
            />
          </div>

          <div className="modal-actions">
            <button
              type="button"
              className="btn btn-secondary"
              onClick={handleClose}
              disabled={loading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
            >
              {loading ? 'Saving...' : 'Save'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default HabitDetailsModal;
