import { ChakraProvider } from '@chakra-ui/react';
import { GoogleOAuthProvider } from '@react-oauth/google';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './hooks/useAuth';
import { LoginPage } from './pages/LoginPage';
import { MainPage } from './pages/MainPage';
import { CalendarPage } from './pages/CalendarPage';
import { HabitDetailPage } from './pages/HabitDetailPage';

// Protected route wrapper
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
    const { isAuthenticated } = useAuth();
    return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
};

// App routes wrapper that requires AuthProvider
const AppRoutes = () => (
    <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route
            path="/"
            element={
                <ProtectedRoute>
                    <MainPage />
                </ProtectedRoute>
            }
        />
        <Route path="/calendar" element={
                <ProtectedRoute>
                    <CalendarPage />
                </ProtectedRoute>
            } />
        <Route path="/habits/:id" element={
                <ProtectedRoute>
                    <HabitDetailPage />
                </ProtectedRoute>
            } />
    </Routes>
);

function App() {
    const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;

    return (
        <Router>
            <ChakraProvider>
                <GoogleOAuthProvider clientId={clientId}>
                    <AuthProvider>
                        <AppRoutes />
                    </AuthProvider>
                </GoogleOAuthProvider>
            </ChakraProvider>
        </Router>
    );
}

export default App;
