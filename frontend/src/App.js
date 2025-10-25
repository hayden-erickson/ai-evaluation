import React, { useState, useEffect } from 'react';
import { authAPI } from './services/api';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import HabitList from './components/Habits/HabitList';

/**
 * Main App Component
 * Handles authentication state and routing
 */
function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [showRegister, setShowRegister] = useState(false);
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // Check authentication status on mount
  useEffect(() => {
    const checkAuth = () => {
      const authenticated = authAPI.isAuthenticated();
      setIsAuthenticated(authenticated);
      setLoading(false);
    };

    checkAuth();
  }, []);

  /**
   * Handle successful login
   */
  const handleLogin = (userData) => {
    setUser(userData);
    setIsAuthenticated(true);
  };

  /**
   * Handle successful registration
   */
  const handleRegister = (userData) => {
    setUser(userData);
    setIsAuthenticated(true);
  };

  /**
   * Handle logout
   */
  const handleLogout = () => {
    authAPI.logout();
    setUser(null);
    setIsAuthenticated(false);
  };

  // Show loading state
  if (loading) {
    return (
      <div className="loading" style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        Loading...
      </div>
    );
  }

  // Show authentication screens if not authenticated
  if (!isAuthenticated) {
    if (showRegister) {
      return (
        <Register
          onRegister={handleRegister}
          onSwitchToLogin={() => setShowRegister(false)}
        />
      );
    }
    return (
      <Login
        onLogin={handleLogin}
        onSwitchToRegister={() => setShowRegister(true)}
      />
    );
  }

  // Show main app if authenticated
  return (
    <div className="container">
      <div className="header">
        <h1>🎯 Habit Tracker</h1>
        <p style={{ color: 'var(--text-secondary)', marginBottom: '8px' }}>
          {user?.name ? `Welcome back, ${user.name}!` : 'Keep your streaks alive!'}
        </p>
        <div className="header-actions">
          <button className="btn btn-secondary" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </div>

      <HabitList />
    </div>
  );
}

export default App;
