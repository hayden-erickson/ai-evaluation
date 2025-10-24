import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import HabitList from './HabitList';
import '../styles/Dashboard.css';

/**
 * Dashboard Component
 * Main dashboard view for authenticated users
 */
const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <div className="header-content">
          <h1 className="logo">Habit Tracker</h1>
          <div className="user-info">
            <span className="user-name">{user?.name}</span>
            <button className="btn-logout" onClick={logout}>
              Logout
            </button>
          </div>
        </div>
      </header>

      <main className="dashboard-main">
        <HabitList />
      </main>
    </div>
  );
};

export default Dashboard;
