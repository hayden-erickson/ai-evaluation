import React, { useState, useEffect } from 'react';
import { View, StyleSheet } from 'react-native';
import Modal from 'react-native-modal';
import { Input, Button, Text } from '@rneui/themed';
import { Log } from '../types/types';

interface LogDetailsModalProps {
  isVisible: boolean;
  onClose: () => void;
  onSave: (logData: { notes: string }) => void;
  log: Log | null;
}

const LogDetailsModal = ({ isVisible, onClose, onSave, log }: LogDetailsModalProps) => {
  const [notes, setNotes] = useState('');

  useEffect(() => {
    if (log) {
      setNotes(log.notes);
    } else {
      setNotes('');
    }
  }, [log]);

  const handleSave = () => {
    onSave({ notes });
  };

  return (
    <Modal isVisible={isVisible} onBackdropPress={onClose}>
      <View style={styles.content}>
        <Text h4>{log ? 'Edit Log' : 'Add Log'}</Text>
        <Input
          placeholder="Notes"
          value={notes}
          onChangeText={setNotes}
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

export default LogDetailsModal;
