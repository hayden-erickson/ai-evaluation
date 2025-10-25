import { useState, useEffect } from 'react';

/**
 * Modal for creating or editing a log entry
 */
const LogDetailsModal = ({ isOpen, onClose, onSave, log = null, date = null }) => {
  const [notes, setNotes] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // Populate form when editing an existing log
  useEffect(() => {
    if (log) {
      setNotes(log.notes || '');
    } else {
      setNotes('');
    }
    setError('');
  }, [log, isOpen]);

  /**
   * Handle form submission
   */
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await onSave({ notes });
      onClose();
    } catch (err) {
      setError(err.message || 'Failed to save log');
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  const modalTitle = log ? 'Edit Log' : date ? `Log for ${date}` : 'Create Log';

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>{modalTitle}</h2>
          <button className="modal-close" onClick={onClose} aria-label="Close">
            ×
          </button>
        </div>

        {error && (
          <div className="error-message">
            <span>⚠️</span>
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label" htmlFor="notes">
              Notes
            </label>
            <textarea
              id="notes"
              className="form-textarea"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Add notes about your progress..."
              disabled={loading}
              autoFocus
            />
          </div>

          <div className="modal-footer">
            <button 
              type="button" 
              className="btn btn-outline" 
              onClick={onClose}
              disabled={loading}
            >
              Cancel
            </button>
            <button 
              type="submit" 
              className="btn btn-success"
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

export default LogDetailsModal;
