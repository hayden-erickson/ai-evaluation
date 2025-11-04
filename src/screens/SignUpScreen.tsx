import React, { useState, useContext } from 'react';
import { View, StyleSheet } from 'react-native';
import { Input, Button, Text } from '@rneui/themed';
import { AuthContext } from '../contexts/AuthContext';
import { colors } from '../styles/colors';

const SignUpScreen = ({ navigation }) => {
  const [name, setName] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [timeZone, setTimeZone] = useState('');
  const [error, setError] = useState('');
  const { register } = useContext(AuthContext);

  const handleSignUp = async () => {
    if (!name || !phoneNumber || !password || !timeZone) {
      setError('All fields are required.');
      return;
    }
    try {
      await register(name, phoneNumber, password, timeZone);
      // navigation is handled by AppNavigator after successful registration and login
    } catch (e) {
      setError('Failed to sign up. Please try again.');
    }
  };

  return (
    <View style={styles.container}>
      <Text h3>Sign Up</Text>
      {error ? <Text style={styles.error}>{error}</Text> : null}
      <Input
        placeholder="Name"
        value={name}
        onChangeText={setName}
      />
      <Input
        placeholder="Phone Number"
        value={phoneNumber}
        onChangeText={setPhoneNumber}
        keyboardType="phone-pad"
        autoCapitalize="none"
      />
      <Input
        placeholder="Time Zone (e.g., America/New_York)"
        value={timeZone}
        onChangeText={setTimeZone}
        autoCapitalize="none"
      />
      <Input
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
      />
      <Button title="Sign Up" onPress={handleSignUp} />
      <Button
        title="Already have an account? Login"
        onPress={() => navigation.navigate('Login')}
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

export default SignUpScreen;
