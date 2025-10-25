import React, { useState } from 'react';
import { register, login } from '../../services/api';
import './Auth.css';

/**
 * Register component for new user registration
 * @param {function} onRegisterSuccess - Callback function called after successful registration
 * @param {function} onSwitchToLogin - Callback to switch to login view
 */
const Register = ({ onRegisterSuccess, onSwitchToLogin }) => {
  const [formData, setFormData] = useState({
    name: '',
    phoneNumber: '',
    timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
    password: '',
    confirmPassword: '',
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
  };

  /**
   * Validate form inputs before submission
   * @returns {boolean} True if inputs are valid
   */
  const validateInputs = () => {
    if (!formData.name.trim()) {
      setError('Name is required');
      return false;
    }
    if (!formData.phoneNumber.trim()) {
      setError('Phone number is required');
      return false;
    }
    if (!formData.password) {
      setError('Password is required');
      return false;
    }
    if (formData.password.length < 8) {
      setError('Password must be at least 8 characters long');
      return false;
    }
    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match');
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
      // Register the user
      await register({
        name: formData.name,
        phone_number: formData.phoneNumber,
        time_zone: formData.timeZone,
        password: formData.password,
      });

      // Automatically log in after successful registration
      const loginResponse = await login(formData.phoneNumber, formData.password);
      
      // Call success callback with user data
      onRegisterSuccess(loginResponse.user);
    } catch (err) {
      // Display error message to user
      setError(err.message || 'Registration failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h1 className="auth-title">Create Account</h1>
        <p className="auth-subtitle">Start tracking your habits today</p>

        {/* Display error message if exists */}
        {error && <div className="error-message">{error}</div>}

        <form onSubmit={handleSubmit} className="auth-form">
          <div className="form-group">
            <label htmlFor="name">Name</label>
            <input
              id="name"
              name="name"
              type="text"
              placeholder="Enter your name"
              value={formData.name}
              onChange={handleChange}
              disabled={loading}
              autoComplete="name"
            />
          </div>

          <div className="form-group">
            <label htmlFor="phoneNumber">Phone Number</label>
            <input
              id="phoneNumber"
              name="phoneNumber"
              type="tel"
              placeholder="Enter your phone number"
              value={formData.phoneNumber}
              onChange={handleChange}
              disabled={loading}
              autoComplete="tel"
            />
          </div>

          <div className="form-group">
            <label htmlFor="timeZone">Time Zone</label>
            <input
              id="timeZone"
              name="timeZone"
              type="text"
              placeholder="Time zone"
              value={formData.timeZone}
              onChange={handleChange}
              disabled={loading}
            />
            <small className="text-secondary">Auto-detected from your browser</small>
          </div>

          <div className="form-group">
            <label htmlFor="password">Password</label>
            <input
              id="password"
              name="password"
              type="password"
              placeholder="Create a password (min 8 characters)"
              value={formData.password}
              onChange={handleChange}
              disabled={loading}
              autoComplete="new-password"
            />
          </div>

          <div className="form-group">
            <label htmlFor="confirmPassword">Confirm Password</label>
            <input
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              placeholder="Confirm your password"
              value={formData.confirmPassword}
              onChange={handleChange}
              disabled={loading}
              autoComplete="new-password"
            />
          </div>

          <button 
            type="submit" 
            className="primary auth-button"
            disabled={loading}
          >
            {loading ? 'Creating Account...' : 'Register'}
          </button>
        </form>

        <p className="auth-switch">
          Already have an account?{' '}
          <button 
            type="button"
            className="link-button"
            onClick={onSwitchToLogin}
            disabled={loading}
          >
            Login
          </button>
        </p>
      </div>
    </div>
  );
};

export default Register;

