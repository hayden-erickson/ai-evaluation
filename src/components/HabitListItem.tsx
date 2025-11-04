import React, { useState } from 'react';
import { View, StyleSheet } from 'react-native';
import { Text, Button, Icon } from '@rneui/themed';
import { Habit, Log } from '../types/types';
import LogDetailsModal from './LogDetailsModal';
import { createLog } from '../api/logService';
import { colors } from '../styles/colors';

interface HabitListItemProps {
  habit: Habit;
  onEdit: () => void;
  onDelete: () => void;
  onLogCreated: () => void;
}

const HabitListItem = ({ habit, onEdit, onDelete, onLogCreated }: HabitListItemProps) => {
  const [isLogModalVisible, setLogModalVisible] = useState(false);

  const handleSaveLog = async (logData: { notes: string }) => {
    try {
      await createLog(habit.id, logData);
      onLogCreated();
    } catch (error) {
      console.error('Failed to save log', error);
    } finally {
      setLogModalVisible(false);
    }
  };

  return (
    <View style={styles.container}>
      <View style={styles.habitInfo}>
        <Text h4>{habit.name}</Text>
        <Text>{habit.description}</Text>
      </View>
      <View style={styles.buttons}>
        <Button
          title="Log"
          onPress={() => setLogModalVisible(true)}
        />
        <Button
          icon={<Icon name="edit" size={25} />}
          type="clear"
          onPress={onEdit}
        />
        <Button
          icon={<Icon name="delete" size={25} color="red" />}
          type="clear"
          onPress={onDelete}
        />
      </View>
      <LogDetailsModal
        isVisible={isLogModalVisible}
        onClose={() => setLogModalVisible(false)}
        onSave={handleSaveLog}
        log={null}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 15,
    borderBottomWidth: 1,
    borderBottomColor: colors.gray,
    backgroundColor: colors.white,
    borderRadius: 8,
    marginBottom: 10,
  },
  habitInfo: {
    flex: 1,
  },
  buttons: {
    flexDirection: 'row',
    alignItems: 'center',
  },
});

export default HabitListItem;
