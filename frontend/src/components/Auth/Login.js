import React, { useState } from 'react';
import { authAPI } from '../../services/api';

/**
 * Login Component
 * Handles user authentication
 */
const Login = ({ onLogin, onSwitchToRegister }) => {
  const [formData, setFormData] = useState({
    phone_number: '',
    password: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

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
      if (!formData.phone_number || !formData.password) {
        throw new Error('Please fill in all fields');
      }

      // Attempt login
      const response = await authAPI.login(formData);
      
      // Call parent callback on success
      if (onLogin) {
        onLogin(response.user);
      }
    } catch (err) {
      setError(err.message || 'Login failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h1 className="auth-title">Habit Tracker</h1>
        
        {error && <div className="error-message">{error}</div>}
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label" htmlFor="phone_number">
              Phone Number
            </label>
            <input
              type="tel"
              id="phone_number"
              name="phone_number"
              className="form-input"
              value={formData.phone_number}
              onChange={handleChange}
              placeholder="Enter your phone number"
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="password">
              Password
            </label>
            <input
              type="password"
              id="password"
              name="password"
              className="form-input"
              value={formData.password}
              onChange={handleChange}
              placeholder="Enter your password"
              disabled={loading}
            />
          </div>

          <button 
            type="submit" 
            className="btn btn-primary" 
            style={{ width: '100%' }}
            disabled={loading}
          >
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </form>

        <div className="auth-link">
          Don't have an account?{' '}
          <a href="#" onClick={(e) => { e.preventDefault(); onSwitchToRegister(); }}>
            Register
          </a>
        </div>
      </div>
    </div>
  );
};

export default Login;
