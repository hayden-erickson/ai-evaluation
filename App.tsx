import React, { useState } from 'react';
import {
  SafeAreaView,
  StatusBar,
  StyleSheet,
  useColorScheme,
  View,
} from 'react-native';
import HabitList from './components/HabitList';
import Login from './components/Login';
import SignUp from './components/SignUp';
import { setAuthToken } from './services/api';

type Screen = 'Login' | 'SignUp' | 'HabitList';

function App() {
  const isDarkMode = useColorScheme() === 'dark';
  const [screen, setScreen] = useState<Screen>('Login');
  const [token, setToken] = useState<string | null>(null);

  const handleLogin = (newToken: string) => {
    setToken(newToken);
    setAuthToken(newToken);
    setScreen('HabitList');
  };

  const handleSignUp = (newToken: string) => {
    setToken(newToken);
    setAuthToken(newToken);
    setScreen('HabitList');
  };

  const backgroundStyle = {
    backgroundColor: isDarkMode ? '#121212' : '#f0f8ff',
    flex: 1,
  };

  return (
    <SafeAreaView style={backgroundStyle}>
      <StatusBar barStyle={isDarkMode ? 'light-content' : 'dark-content'} />
      <View style={styles.container}>
        {screen === 'Login' && <Login onLogin={handleLogin} onNavigateToSignUp={() => setScreen('SignUp')} />}
        {screen === 'SignUp' && <SignUp onSignUp={handleSignUp} onNavigateToLogin={() => setScreen('Login')} />}
        {screen === 'HabitList' && <HabitList />}
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 10,
    fontFamily: 'sans-serif',
  },
});

export default App;
