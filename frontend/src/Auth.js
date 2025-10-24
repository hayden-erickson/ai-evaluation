import React, { useState } from 'react';
import { authAPI } from './api';

/**
 * Auth component for user login and registration
 */
function Auth({ onAuthenticated }) {
  const [isLogin, setIsLogin] = useState(true);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // Form state
  const [formData, setFormData] = useState({
    name: '',
    phone_number: '',
    password: '',
    time_zone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'America/New_York',
  });

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
   * Validate form inputs
   */
  const validateForm = () => {
    if (!formData.phone_number.trim()) {
      setError('Phone number is required');
      return false;
    }
    if (!formData.password.trim()) {
      setError('Password is required');
      return false;
    }
    if (!isLogin && !formData.name.trim()) {
      setError('Name is required');
      return false;
    }
    // Basic phone number validation
    if (!/^\+?[\d\s-()]+$/.test(formData.phone_number)) {
      setError('Please enter a valid phone number');
      return false;
    }
    // Password length validation
    if (formData.password.length < 6) {
      setError('Password must be at least 6 characters');
      return false;
    }
    return true;
  };

  /**
   * Handle form submission
   */
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setLoading(true);
    setError('');

    try {
      if (isLogin) {
        // Login
        const response = await authAPI.login({
          phone_number: formData.phone_number,
          password: formData.password,
        });
        onAuthenticated(response.user);
      } else {
        // Register
        const response = await authAPI.register({
          name: formData.name,
          phone_number: formData.phone_number,
          password: formData.password,
          time_zone: formData.time_zone,
        });
        onAuthenticated(response.user);
      }
    } catch (err) {
      setError(err.message || 'Authentication failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Toggle between login and register
   */
  const toggleMode = () => {
    setIsLogin(!isLogin);
    setError('');
    setFormData({
      ...formData,
      name: '',
    });
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h2>{isLogin ? 'Login' : 'Sign Up'}</h2>
        
        {error && <div className="error-message">{error}</div>}
        
        <form onSubmit={handleSubmit}>
          {!isLogin && (
            <div className="form-group">
              <label htmlFor="name">Name</label>
              <input
                type="text"
                id="name"
                name="name"
                value={formData.name}
                onChange={handleChange}
                placeholder="Enter your name"
                disabled={loading}
              />
            </div>
          )}
          
          <div className="form-group">
            <label htmlFor="phone_number">Phone Number</label>
            <input
              type="tel"
              id="phone_number"
              name="phone_number"
              value={formData.phone_number}
              onChange={handleChange}
              placeholder="+1234567890"
              disabled={loading}
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="password">Password</label>
            <input
              type="password"
              id="password"
              name="password"
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
            {loading ? 'Please wait...' : (isLogin ? 'Login' : 'Sign Up')}
          </button>
        </form>
        
        <div className="auth-toggle">
          {isLogin ? "Don't have an account? " : 'Already have an account? '}
          <button onClick={toggleMode} disabled={loading}>
            {isLogin ? 'Sign Up' : 'Login'}
          </button>
        </div>
      </div>
    </div>
  );
}

export default Auth;
