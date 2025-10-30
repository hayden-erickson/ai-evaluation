import React, { useState, createContext, useContext, useEffect } from 'react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from 'react-router-dom';
import HabitList from './components/HabitList';
import AddNewHabitButton from './components/AddNewHabitButton';
import HabitDetailsModal from './components/HabitDetailsModal';
import LogDetailsModal from './components/LogDetailsModal';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';

const AuthContext = createContext(null);

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'));

  const setAuthToken = (newToken) => {
    if (newToken) {
      localStorage.setItem('token', newToken);
    } else {
      localStorage.removeItem('token');
    }
    setToken(newToken);
  };

  return (
    <AuthContext.Provider value={{ token, setAuthToken }}>
      <Router>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route
            path="/"
            element={
              token ? <MainApp /> : <Navigate to="/login" />
            }
          />
        </Routes>
      </Router>
    </AuthContext.Provider>
  );
}

function MainApp() {
  const [habits, setHabits] = useState([]);
  const [isHabitModalOpen, setIsHabitModalOpen] = useState(false);
  const [isLogModalOpen, setIsLogModalOpen] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState(null);
  const [selectedLog, setSelectedLog] = useState(null);
  const { token, setAuthToken } = useAuth();

  const fetchHabits = async () => {
    const response = await fetch('/habits', {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (response.ok) {
      const data = await response.json();
      setHabits(data || []);
    }
  };

  useEffect(() => {
    fetchHabits();
  }, [token]);

  const handleSaveHabit = async (habit) => {
    const method = habit.id ? 'PUT' : 'POST';
    const url = habit.id ? `/habits/${habit.id}` : '/habits';
    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(habit),
    });

    if (response.ok) {
      fetchHabits();
    }
  };

  const handleDeleteHabit = async (habitId) => {
    const response = await fetch(`/habits/${habitId}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    });

    if (response.ok) {
      fetchHabits();
    }
  };

  const handleSaveLog = async (log) => {
    const habitId = selectedHabit.id;
    const method = log.id ? 'PUT' : 'POST';
    const url = log.id ? `/logs/${log.id}` : `/habits/${habitId}/logs`;

    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(log),
    });

    if (response.ok) {
      fetchHabits();
    }
  };

  const openLogModal = (habit, log = null) => {
    setSelectedHabit(habit);
    setSelectedLog(log);
    setIsLogModalOpen(true);
  };

  return (
    <div className="bg-gray-100 min-h-screen">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8 flex justify-between items-center">
          <h1 className="text-3xl font-bold text-gray-900">Habit Tracker</h1>
          <div>
            <AddNewHabitButton onOpen={() => { setSelectedHabit(null); setIsHabitModalOpen(true); }} />
            <button onClick={() => setAuthToken(null)} className="ml-4 bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">
              Logout
            </button>
          </div>
        </div>
      </header>
      <main>
        <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          <HabitList
            habits={habits}
            onOpenLogModal={openLogModal}
            onOpenHabitModal={(habit) => { setSelectedHabit(habit); setIsHabitModalOpen(true); }}
            onDeleteHabit={handleDeleteHabit}
          />
        </div>
      </main>
      <HabitDetailsModal isOpen={isHabitModalOpen} onClose={() => setIsHabitModalOpen(false)} onSave={handleSaveHabit} habit={selectedHabit} />
      <LogDetailsModal isOpen={isLogModalOpen} onClose={() => setIsLogModalOpen(false)} onSave={handleSaveLog} log={selectedLog} />
    </div>
  );
}

export const useAuth = () => useContext(AuthContext);

export default App;

