import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';

/**
 * Application header with logout functionality
 */
const Header = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  /**
   * Handle user logout
   */
  const handleLogout = () => {
    if (window.confirm('Are you sure you want to logout?')) {
      logout();
      navigate('/login');
    }
  };

  return (
    <header className="app-header">
      <div className="app-header-content">
        <h1>🎯 Habit Tracker</h1>
        <div className="flex gap-md" style={{ alignItems: 'center' }}>
          {user && (
            <span style={{ color: 'var(--text-secondary)' }}>
              Welcome, {user.name}!
            </span>
          )}
          <button className="btn btn-outline btn-sm" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </div>
    </header>
  );
};

export default Header;
