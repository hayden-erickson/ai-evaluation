import React, { useState } from 'react';
import { authAPI } from '../../services/api';

/**
 * Register Component
 * Handles new user registration
 */
const Register = ({ onRegister, onSwitchToLogin }) => {
  const [formData, setFormData] = useState({
    name: '',
    phone_number: '',
    time_zone: 'America/New_York',
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
      if (!formData.name || !formData.phone_number || !formData.password) {
        throw new Error('Please fill in all required fields');
      }

      if (formData.password !== formData.confirmPassword) {
        throw new Error('Passwords do not match');
      }

      if (formData.password.length < 8) {
        throw new Error('Password must be at least 8 characters long');
      }

      // Prepare registration data
      const registrationData = {
        name: formData.name,
        phone_number: formData.phone_number,
        time_zone: formData.time_zone,
        password: formData.password,
      };

      // Attempt registration
      const response = await authAPI.register(registrationData);
      
      // Call parent callback on success
      if (onRegister) {
        onRegister(response.user);
      }
    } catch (err) {
      setError(err.message || 'Registration failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h1 className="auth-title">Create Account</h1>
        
        {error && <div className="error-message">{error}</div>}
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label" htmlFor="name">
              Name *
            </label>
            <input
              type="text"
              id="name"
              name="name"
              className="form-input"
              value={formData.name}
              onChange={handleChange}
              placeholder="Enter your name"
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="phone_number">
              Phone Number *
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
            <label className="form-label" htmlFor="time_zone">
              Time Zone
            </label>
            <input
              type="text"
              id="time_zone"
              name="time_zone"
              className="form-input"
              value={formData.time_zone}
              onChange={handleChange}
              placeholder="e.g., America/New_York"
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="password">
              Password *
            </label>
            <input
              type="password"
              id="password"
              name="password"
              className="form-input"
              value={formData.password}
              onChange={handleChange}
              placeholder="At least 8 characters"
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="confirmPassword">
              Confirm Password *
            </label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              className="form-input"
              value={formData.confirmPassword}
              onChange={handleChange}
              placeholder="Re-enter your password"
              disabled={loading}
            />
          </div>

          <button 
            type="submit" 
            className="btn btn-primary" 
            style={{ width: '100%' }}
            disabled={loading}
          >
            {loading ? 'Creating Account...' : 'Register'}
          </button>
        </form>

        <div className="auth-link">
          Already have an account?{' '}
          <a href="#" onClick={(e) => { e.preventDefault(); onSwitchToLogin(); }}>
            Login
          </a>
        </div>
      </div>
    </div>
  );
};

export default Register;
