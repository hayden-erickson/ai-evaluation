import React from 'react';
import { FlatList, StyleSheet } from 'react-native';
import { Text } from '@rneui/themed';
import { Habit } from '../types/types';
import HabitListItem from './HabitListItem';
import { colors } from '../styles/colors';

interface HabitListProps {
  habits: Habit[];
  onEdit: (habit: Habit) => void;
  onDelete: (id: number) => void;
  onLogCreated: () => void;
}

const HabitList = ({ habits, onEdit, onDelete, onLogCreated }: HabitListProps) => {
  return (
    <FlatList
      data={habits}
      keyExtractor={(item) => item.id.toString()}
      renderItem={({ item }) => (
        <HabitListItem
          habit={item}
          onEdit={() => onEdit(item)}
          onDelete={() => onDelete(item.id)}
          onLogCreated={onLogCreated}
        />
      )}
      ListEmptyComponent={<Text style={styles.centered}>No habits yet. Add one!</Text>}
    />
  );
};

const styles = StyleSheet.create({
  centered: {
    flex: 1,
    textAlign: 'center',
    marginTop: 20,
    color: colors.text,
  },
});

export default HabitList;
