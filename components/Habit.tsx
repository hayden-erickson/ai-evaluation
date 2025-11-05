import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, Alert, TouchableOpacity } from 'react-native';
import LogDetailsModal from './LogDetailsModal';
import StreakList from './StreakList';
import { getLogs, createLog } from '../services/api';

interface Log {
  id: string;
  habit_id: string;
  date: string;
  notes: string;
}

interface HabitProps {
  habit: {
    id: string;
    name: string;
    description: string;
  };
  onEdit: (habit: HabitProps['habit']) => void;
  onDelete: (id: string) => void;
}

const Habit: React.FC<HabitProps> = ({ habit, onEdit, onDelete }) => {
  const [logModalVisible, setLogModalVisible] = useState(false);
  const [logs, setLogs] = useState<Log[]>([]);

  const fetchLogs = async () => {
    try {
      const logsData = await getLogs(habit.id);
      setLogs(logsData.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()));
    } catch (error) {
      Alert.alert('Error', 'Failed to fetch logs');
    }
  };

  useEffect(() => {
    fetchLogs();
  }, []);

  const calculateStreak = () => {
    if (logs.length === 0) return 0;

    const logDates = new Set(logs.map(log => {
        const d = new Date(log.date);
        d.setUTCHours(0, 0, 0, 0);
        return d.getTime();
    }));

    let streak = 0;
    let skippedDays = 0;
    const today = new Date();
    today.setUTCHours(0, 0, 0, 0);

    for (let i = 0; i < 60; i++) { // Check up to 60 days back
        const dateToCheck = new Date(today);
        dateToCheck.setDate(today.getDate() - i);

        if (logDates.has(dateToCheck.getTime())) {
            streak++;
            skippedDays = 0;
        } else {
            skippedDays++;
            if (skippedDays > 1) {
                break;
            }
        }
    }
    return streak;
  };

  const handleSaveLog = async (notes: string) => {
    try {
      await createLog(habit.id, { date: new Date().toISOString().split('T')[0], notes });
      fetchLogs();
    } catch (error) {
      Alert.alert('Error', 'Failed to save log');
    }
    setLogModalVisible(false);
  };

  return (
    <View style={styles.container}>
      <Text style={styles.name}>{habit.name}</Text>
      <Text style={styles.description}>{habit.description}</Text>
      <View style={styles.streakContainer}>
        <Text style={styles.streakText}>Streak: {calculateStreak()}</Text>
        <StreakList logs={logs} />
      </View>
      <View style={styles.buttonsContainer}>
        <TouchableOpacity style={[styles.button, styles.logButton]} onPress={() => setLogModalVisible(true)}>
          <Text style={styles.buttonText}>Log</Text>
        </TouchableOpacity>
        <TouchableOpacity style={[styles.button, styles.editButton]} onPress={() => onEdit(habit)}>
          <Text style={styles.buttonText}>Edit</Text>
        </TouchableOpacity>
        <TouchableOpacity style={[styles.button, styles.deleteButton]} onPress={() => onDelete(habit.id)}>
          <Text style={styles.buttonText}>Delete</Text>
        </TouchableOpacity>
      </View>
      <LogDetailsModal
        visible={logModalVisible}
        onClose={() => setLogModalVisible(false)}
        onSave={handleSaveLog}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: '#ffffff',
    padding: 20,
    borderRadius: 15,
    marginBottom: 15,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  name: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#333',
  },
  description: {
    fontSize: 16,
    color: '#666',
    marginTop: 5,
    marginBottom: 15,
  },
  streakContainer: {
    marginBottom: 15,
  },
  streakText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#444',
    marginBottom: 10,
  },
  buttonsContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginTop: 10,
  },
  button: {
    paddingVertical: 10,
    paddingHorizontal: 20,
    borderRadius: 20,
    alignItems: 'center',
    justifyContent: 'center',
  },
  buttonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: 'bold',
  },
  logButton: {
    backgroundColor: '#90ee90',
  },
  editButton: {
    backgroundColor: '#87ceeb',
  },
  deleteButton: {
    backgroundColor: '#f08080',
  },
});

export default Habit;
