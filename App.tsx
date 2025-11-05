import React from 'react';
import { StatusBar, StyleSheet, View } from 'react-native';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { AuthProvider, useAuth } from './src/context/AuthContext';
import AuthScreen from './src/components/AuthScreen';
import HabitList from './src/components/HabitList';
import { colors } from './src/theme';

function Root() {
  const { token, loading } = useAuth();

  return (
    <View style={styles.container}>
      {loading ? null : token ? <HabitList /> : <AuthScreen />}
    </View>
  );
}

function App() {
  return (
    <SafeAreaProvider>
      <StatusBar barStyle={'dark-content'} />
      <AuthProvider>
        <Root />
      </AuthProvider>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
});

export default App;
