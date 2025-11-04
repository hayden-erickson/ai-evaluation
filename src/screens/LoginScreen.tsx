import React, { useState, useContext } from 'react';
import { View, StyleSheet } from 'react-native';
import { Input, Button, Text } from '@rneui/themed';
import { AuthContext } from '../contexts/AuthContext';
import { colors } from '../styles/colors';

const LoginScreen = ({ navigation }) => {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login } = useContext(AuthContext);

  const handleLogin = async () => {
    if (!phoneNumber || !password) {
      setError('Phone number and password are required.');
      return;
    }
    try {
      await login(phoneNumber, password);
      // navigation is handled by AppNavigator
    } catch (e) {
      setError('Failed to log in. Please check your credentials.');
    }
  };

  return (
    <View style={styles.container}>
      <Text h3>Login</Text>
      {error ? <Text style={styles.error}>{error}</Text> : null}
      <Input
        placeholder="Phone Number"
        value={phoneNumber}
        onChangeText={setPhoneNumber}
        keyboardType="phone-pad"
        autoCapitalize="none"
      />
      <Input
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
      />
      <Button title="Login" onPress={handleLogin} />
      <Button
        title="Don't have an account? Sign Up"
        onPress={() => navigation.navigate('SignUp')}
        type="clear"
        titleStyle={{ color: colors.primary }}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
    backgroundColor: colors.background,
  },
  error: {
    color: colors.error,
    marginBottom: 10,
  },
});

export default LoginScreen;
