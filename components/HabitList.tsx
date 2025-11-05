import React, { useState, useEffect } from 'react';
import { View, FlatList, StyleSheet, Text, Alert, TouchableOpacity } from 'react-native';
import Habit from './Habit';
import HabitDetailsModal from './HabitDetailsModal';
import { getHabits, createHabit, updateHabit, deleteHabit } from '../services/api';

interface HabitType {
  id: string;
  name: string;
  description: string;
}

const HabitList = () => {
  const [habits, setHabits] = useState<HabitType[]>([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState<HabitType | undefined>(undefined);

  // Fetches habits from the API
  const fetchHabits = async () => {
    try {
      const habitsData = await getHabits();
      setHabits(habitsData);
    } catch (error) {
      Alert.alert('Error', 'Failed to fetch habits');
    }
  };

  useEffect(() => {
    fetchHabits();
  }, []);

  // Opens the modal for editing a habit
  const handleEdit = (habit: HabitType) => {
    setSelectedHabit(habit);
    setModalVisible(true);
  };

  // Opens the modal for adding a new habit
  const handleAddNew = () => {
    setSelectedHabit(undefined);
    setModalVisible(true);
  };

  // Saves a new or edited habit
  const handleSave = async (name: string, description: string) => {
    try {
      if (selectedHabit) {
        await updateHabit(selectedHabit.id, { name, description });
      } else {
        await createHabit({ name, description });
      }
      fetchHabits();
    } catch (error) {
      Alert.alert('Error', 'Failed to save habit');
    }
    setModalVisible(false);
  };

  // Deletes a habit
  const handleDelete = async (id: string) => {
    try {
      await deleteHabit(id);
      fetchHabits();
    } catch (error) {
      Alert.alert('Error', 'Failed to delete habit');
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>My Habits</Text>
      <FlatList
        data={habits}
        renderItem={({ item }) => <Habit habit={item} onEdit={handleEdit} onDelete={handleDelete} />}
        keyExtractor={item => item.id}
        contentContainerStyle={styles.list}
      />
      <TouchableOpacity style={styles.addButton} onPress={handleAddNew}>
        <Text style={styles.addButtonText}>Add New Habit</Text>
      </TouchableOpacity>
      <HabitDetailsModal
        visible={modalVisible}
        onClose={() => setModalVisible(false)}
        onSave={handleSave}
        habit={selectedHabit}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    marginBottom: 20,
    color: '#333',
    textAlign: 'center',
  },
  list: {
    paddingBottom: 80,
  },
  addButton: {
    backgroundColor: '#add8e6',
    padding: 15,
    borderRadius: 25,
    alignItems: 'center',
    justifyContent: 'center',
    position: 'absolute',
    bottom: 20,
    left: 20,
    right: 20,
  },
  addButtonText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: 'bold',
  },
});

export default HabitList;
