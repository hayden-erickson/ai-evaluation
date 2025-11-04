import React, { useState, useEffect } from 'react';
import { View, StyleSheet } from 'react-native';
import Modal from 'react-native-modal';
import { Input, Button, Text } from '@rneui/themed';
import { Habit } from '../types/types';

interface HabitDetailsModalProps {
  isVisible: boolean;
  onClose: () => void;
  onSave: (habitData: { name: string; description: string }) => void;
  habit: Habit | null;
}

const HabitDetailsModal = ({ isVisible, onClose, onSave, habit }: HabitDetailsModalProps) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  useEffect(() => {
    if (habit) {
      setName(habit.name);
      setDescription(habit.description);
    } else {
      setName('');
      setDescription('');
    }
  }, [habit]);

  const handleSave = () => {
    onSave({ name, description });
  };

  return (
    <Modal isVisible={isVisible} onBackdropPress={onClose}>
      <View style={styles.content}>
        <Text h4>{habit ? 'Edit Habit' : 'Add Habit'}</Text>
        <Input
          placeholder="Habit Name"
          value={name}
          onChangeText={setName}
        />
        <Input
          placeholder="Description"
          value={description}
          onChangeText={setDescription}
          multiline
        />
        <Button title="Save" onPress={handleSave} />
        <Button title="Cancel" onPress={onClose} type="clear" />
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  content: {
    backgroundColor: 'white',
    padding: 22,
    justifyContent: 'center',
    alignItems: 'center',
    borderRadius: 4,
    borderColor: 'rgba(0, 0, 0, 0.1)',
  },
});

export default HabitDetailsModal;
