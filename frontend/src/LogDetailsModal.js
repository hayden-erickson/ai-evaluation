import React, { useState, useEffect } from 'react';

/**
 * Modal for creating or editing a habit log
 */
function LogDetailsModal({ log, onClose, onSave }) {
  const [formData, setFormData] = useState({
    notes: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // Populate form if editing existing log
  useEffect(() => {
    if (log) {
      setFormData({
        notes: log.notes || '',
      });
    }
  }, [log]);

  /**
   * Handle input changes
   */
  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
    // Clear error when user starts typing
    if (error) setError('');
  };

  /**
   * Handle form submission
   */
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    setLoading(true);
    setError('');

    try {
      await onSave({
        notes: formData.notes.trim(),
      });
      onClose();
    } catch (err) {
      setError(err.message || 'Failed to save log. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Handle modal overlay click
   */
  const handleOverlayClick = (e) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <div className="modal-overlay" onClick={handleOverlayClick}>
      <div className="modal">
        <h2>{log ? 'Edit Log' : 'Add Log Entry'}</h2>
        
        {error && <div className="error-message">{error}</div>}
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="notes">Notes</label>
            <textarea
              id="notes"
              name="notes"
              value={formData.notes}
              onChange={handleChange}
              placeholder="How did it go? Add any notes here..."
              disabled={loading}
              autoFocus
            />
          </div>
          
          <div className="form-actions">
            <button 
              type="button" 
              className="btn" 
              onClick={onClose}
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
}

export default LogDetailsModal;
