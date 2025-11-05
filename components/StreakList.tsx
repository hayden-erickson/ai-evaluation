import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

interface Log {
  id: string;
  habit_id: string;
  date: string;
  notes: string;
}

interface StreakListProps {
  logs: Log[];
}

const StreakList: React.FC<StreakListProps> = ({ logs }) => {
  const today = new Date();
  const days = Array.from({ length: 30 }, (_, i) => {
    const date = new Date(today);
    date.setDate(today.getDate() - i);
    return date.toISOString().split('T')[0];
  }).reverse();

  const loggedDates = new Set(logs.map(log => log.date));

  return (
    <View style={styles.container}>
      {days.map(day => (
        <View
          key={day}
          style={[styles.dayContainer, loggedDates.has(day) && styles.loggedDay]}
        />
      ))}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginTop: 10,
  },
  dayContainer: {
    width: 20,
    height: 20,
    backgroundColor: '#e0e0e0',
    margin: 2,
    borderRadius: 5,
  },
  loggedDay: {
    backgroundColor: '#a0e0a0',
  },
});

export default StreakList;
