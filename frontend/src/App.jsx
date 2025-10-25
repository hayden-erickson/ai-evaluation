import React, { useState, useEffect } from 'react';
import { isAuthenticated, logout } from './services/api';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import HabitList from './components/Habits/HabitList';
import './App.css';

/**
 * Main App component
 * Handles authentication state and routing between views
 */
function App() {
  const [user, setUser] = useState(null);
  const [view, setView] = useState('login'); // 'login' or 'register'
  const [loading, setLoading] = useState(true);

  /**
   * Check authentication status on mount
   */
  useEffect(() => {
    // Check if user is already authenticated
    if (isAuthenticated()) {
      // Note: We don't have user data stored, so we'll just set a placeholder
      // In a production app, you might want to fetch user data from an API
      setUser({ authenticated: true });
    }
    setLoading(false);
  }, []);

  /**
   * Handle successful login
   * @param {object} userData - User data from login response
   */
  const handleLoginSuccess = (userData) => {
    setUser(userData);
  };

  /**
   * Handle successful registration
   * @param {object} userData - User data from registration response
   */
  const handleRegisterSuccess = (userData) => {
    setUser(userData);
  };

  /**
   * Handle logout
   */
  const handleLogout = () => {
    logout();
    setUser(null);
    setView('login');
  };

  /**
   * Switch to register view
   */
  const handleSwitchToRegister = () => {
    setView('register');
  };

  /**
   * Switch to login view
   */
  const handleSwitchToLogin = () => {
    setView('login');
  };

  // Show loading state while checking authentication
  if (loading) {
    return (
      <div className="app-loading">
        <div className="spinner"></div>
      </div>
    );
  }

  // Show authenticated app if user is logged in
  if (user) {
    return (
      <div className="app">
        <nav className="app-nav">
          <div className="nav-content">
            <div className="nav-brand">
              <span className="nav-logo">🎯</span>
              <h2>Habit Tracker</h2>
            </div>
            <div className="nav-user">
              <span className="user-name">
                {user.name || 'User'}
              </span>
              <button 
                className="secondary small"
                onClick={handleLogout}
              >
                Logout
              </button>
            </div>
          </div>
        </nav>
        <main className="app-main">
          <HabitList user={user} />
        </main>
      </div>
    );
  }

  // Show authentication views if user is not logged in
  return (
    <div className="app">
      {view === 'login' ? (
        <Login
          onLoginSuccess={handleLoginSuccess}
          onSwitchToRegister={handleSwitchToRegister}
        />
      ) : (
        <Register
          onRegisterSuccess={handleRegisterSuccess}
          onSwitchToLogin={handleSwitchToLogin}
        />
      )}
    </div>
  );
}

export default App;

