import React, { useState } from 'react';
import { login } from '../../services/api';
import './Auth.css';

/**
 * Login component for user authentication
 * @param {function} onLoginSuccess - Callback function called after successful login
 * @param {function} onSwitchToRegister - Callback to switch to register view
 */
const Login = ({ onLoginSuccess, onSwitchToRegister }) => {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  /**
   * Validate form inputs before submission
   * @returns {boolean} True if inputs are valid
   */
  const validateInputs = () => {
    if (!phoneNumber.trim()) {
      setError('Phone number is required');
      return false;
    }
    if (!password) {
      setError('Password is required');
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
      const response = await login(phoneNumber, password);
      // Call success callback with user data
      onLoginSuccess(response.user);
    } catch (err) {
      // Display error message to user
      setError(err.message || 'Login failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h1 className="auth-title">Welcome Back</h1>
        <p className="auth-subtitle">Track your habits and build better routines</p>

        {/* Display error message if exists */}
        {error && <div className="error-message">{error}</div>}

        <form onSubmit={handleSubmit} className="auth-form">
          <div className="form-group">
            <label htmlFor="phoneNumber">Phone Number</label>
            <input
              id="phoneNumber"
              type="tel"
              placeholder="Enter your phone number"
              value={phoneNumber}
              onChange={(e) => setPhoneNumber(e.target.value)}
              disabled={loading}
              autoComplete="tel"
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              placeholder="Enter your password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              disabled={loading}
              autoComplete="current-password"
            />
          </div>

          <button 
            type="submit" 
            className="primary auth-button"
            disabled={loading}
          >
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </form>

        <p className="auth-switch">
          Don't have an account?{' '}
          <button 
            type="button"
            className="link-button"
            onClick={onSwitchToRegister}
            disabled={loading}
          >
            Register
          </button>
        </p>
      </div>
    </div>
  );
};

export default Login;

