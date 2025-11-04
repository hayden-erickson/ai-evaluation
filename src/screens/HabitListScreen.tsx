import React, { useEffect, useState, useContext, useCallback } from 'react';
import { View, StyleSheet, ActivityIndicator } from 'react-native';
import { Text, Button } from '@rneui/themed';
import { createHabit, getHabits, updateHabit, deleteHabit } from '../api/habitService';
import { Habit } from '../types/types';
import { AuthContext } from '../contexts/AuthContext';
import HabitList from '../components/HabitList';
import HabitDetailsModal from '../components/HabitDetailsModal';
import { colors } from '../styles/colors';

const HabitListScreen = () => {
  const [habits, setHabits] = useState<Habit[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isModalVisible, setModalVisible] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState<Habit | null>(null);
  const { logout } = useContext(AuthContext);

  const fetchHabits = useCallback(async () => {
    try {
      setLoading(true);
      const fetchedHabits = await getHabits();
      setHabits(fetchedHabits);
      setError('');
    } catch (e) {
      setError('Failed to fetch habits.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchHabits();
  }, [fetchHabits]);

  const handleAddHabit = () => {
    setSelectedHabit(null);
    setModalVisible(true);
  };

  const handleEditHabit = (habit: Habit) => {
    setSelectedHabit(habit);
    setModalVisible(true);
  };

  const handleDeleteHabit = async (id: number) => {
    try {
      await deleteHabit(id);
      fetchHabits();
    } catch (e) {
      setError('Failed to delete habit.');
    }
  };

  const handleSaveHabit = async (habitData: { name: string; description: string }) => {
    try {
      if (selectedHabit) {
        await updateHabit(selectedHabit.id, habitData);
      } else {
        await createHabit(habitData);
      }
      fetchHabits();
    } catch (e) {
      setError('Failed to save habit.');
    } finally {
      setModalVisible(false);
      setSelectedHabit(null);
    }
  };

  const handleLogCreated = () => {
    // For now, just refetch all habits to update streak info etc.
    // A more optimized approach might be to just refetch logs for the specific habit
    fetchHabits();
  }


  if (loading) {
    return <ActivityIndicator style={styles.centered} size="large" />;
  }

  if (error) {
    return <Text style={styles.centered}>{error}</Text>;
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Button title="Logout" onPress={logout} buttonStyle={styles.logoutButton}/>
        <Button title="Add New Habit" onPress={handleAddHabit} />
      </View>
      <HabitList habits={habits} onEdit={handleEditHabit} onDelete={handleDeleteHabit} onLogCreated={handleLogCreated} />
      <HabitDetailsModal
        isVisible={isModalVisible}
        onClose={() => setModalVisible(false)}
        onSave={handleSaveHabit}
        habit={selectedHabit}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 10,
    backgroundColor: colors.background,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 10,
  },
  logoutButton: {
    backgroundColor: colors.secondary,
  },
  centered: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    color: colors.text,
  },
});

export default HabitListScreen;
