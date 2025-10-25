import React, { useState, useEffect } from 'react';

/**
 * LogDetailsModal Component
 * Modal for creating or editing a habit log
 */
const LogDetailsModal = ({ isOpen, onClose, onSave, log, date }) => {
  const [formData, setFormData] = useState({
    notes: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // Populate form when editing existing log
  useEffect(() => {
    if (log) {
      setFormData({
        notes: log.notes || '',
      });
    } else {
      setFormData({
        notes: '',
      });
    }
    setError('');
  }, [log, isOpen]);

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
      // Call parent save handler
      await onSave(formData);
      
      // Reset form and close modal
      setFormData({ notes: '' });
      onClose();
    } catch (err) {
      setError(err.message || 'Failed to save log');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Handle modal close
   */
  const handleClose = () => {
    if (!loading) {
      setFormData({ notes: '' });
      setError('');
      onClose();
    }
  };

  /**
   * Format date for display
   */
  const formatDate = (dateStr) => {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', { 
      weekday: 'long', 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    });
  };

  // Don't render if modal is not open
  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={handleClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2 className="modal-title">
            {log ? 'Edit Log' : 'New Log'}
          </h2>
          <button 
            className="modal-close" 
            onClick={handleClose}
            disabled={loading}
          >
            ×
          </button>
        </div>

        {date && (
          <p style={{ marginBottom: '16px', color: 'var(--text-secondary)' }}>
            {formatDate(date)}
          </p>
        )}

        {error && <div className="error-message">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label" htmlFor="notes">
              Notes
            </label>
            <textarea
              id="notes"
              name="notes"
              className="form-textarea"
              value={formData.notes}
              onChange={handleChange}
              placeholder="How did it go? Any thoughts or reflections?"
              disabled={loading}
              autoFocus
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
              className="btn btn-success"
              disabled={loading}
            >
              {loading ? 'Saving...' : 'Save Log'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default LogDetailsModal;
