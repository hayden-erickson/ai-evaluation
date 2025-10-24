import React, { useState, useEffect } from 'react';
import Auth from './Auth';
import HabitList from './HabitList';
import { authAPI } from './api';

/**
 * Main App component
 * Handles authentication state and routing
 */
function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  /**
   * Check authentication status on mount
   */
  useEffect(() => {
    const checkAuth = async () => {
      if (authAPI.isAuthenticated()) {
        // User has token, consider them authenticated
        // In a production app, you might want to validate the token with the backend
        setUser({ authenticated: true });
      }
      setLoading(false);
    };

    checkAuth();
  }, []);

  /**
   * Handle successful authentication
   */
  const handleAuthenticated = (userData) => {
    setUser(userData);
  };

  /**
   * Handle logout
   */
  const handleLogout = () => {
    authAPI.logout();
    setUser(null);
  };

  // Show loading spinner while checking auth
  if (loading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  // Show Auth component if not authenticated
  if (!user) {
    return <Auth onAuthenticated={handleAuthenticated} />;
  }

  // Show HabitList if authenticated
  return <HabitList user={user} onLogout={handleLogout} />;
}

export default App;
