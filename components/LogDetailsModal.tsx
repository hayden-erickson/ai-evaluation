import React, { useState, useEffect } from 'react';
import { Modal, View, Text, TextInput, StyleSheet, TouchableOpacity, Alert } from 'react-native';

interface LogDetailsModalProps {
  visible: boolean;
  onClose: () => void;
  onSave: (notes: string) => void;
  log?: { notes: string };
}

const LogDetailsModal: React.FC<LogDetailsModalProps> = ({ visible, onClose, onSave, log }) => {
  const [notes, setNotes] = useState('');

  useEffect(() => {
    if (visible) {
      if (log) {
        setNotes(log.notes);
      } else {
        setNotes('');
      }
    }
  }, [visible, log]);

  const handleSave = () => {
    if (!notes.trim()) {
      Alert.alert('Validation Error', 'Notes cannot be empty.');
      return;
    }
    onSave(notes);
    onClose();
  };

  return (
    <Modal visible={visible} animationType="slide" transparent={true} onRequestClose={onClose}>
      <View style={styles.modalContainer}>
        <View style={styles.modalContent}>
          <Text style={styles.title}>{log ? 'Edit Log' : 'Add Log'}</Text>
          <TextInput
            style={styles.input}
            placeholder="Notes..."
            value={notes}
            onChangeText={setNotes}
            multiline
          />
          <View style={styles.buttonsContainer}>
            <TouchableOpacity style={[styles.button, styles.cancelButton]} onPress={onClose}>
              <Text style={styles.buttonText}>Cancel</Text>
            </TouchableOpacity>
            <TouchableOpacity
              style={[styles.button, styles.saveButton, !notes.trim() && styles.disabledButton]}
              onPress={handleSave}
              disabled={!notes.trim()}
            >
              <Text style={styles.buttonText}>Save</Text>
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  modalContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
  },
  modalContent: {
    backgroundColor: '#fff',
    padding: 25,
    borderRadius: 15,
    width: '90%',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 5,
  },
  title: {
    fontSize: 22,
    fontWeight: 'bold',
    marginBottom: 20,
    color: '#333',
    textAlign: 'center',
  },
  input: {
    backgroundColor: '#f0f0f0',
    borderWidth: 0,
    padding: 15,
    borderRadius: 10,
    marginBottom: 15,
    fontSize: 16,
    minHeight: 120,
    textAlignVertical: 'top',
  },
  buttonsContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginTop: 10,
  },
  button: {
    flex: 1,
    padding: 15,
    borderRadius: 25,
    alignItems: 'center',
    justifyContent: 'center',
    marginHorizontal: 5,
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
  saveButton: {
    backgroundColor: '#90ee90',
  },
  cancelButton: {
    backgroundColor: '#f08080',
  },
  disabledButton: {
    backgroundColor: '#ccc',
  },
});

export default LogDetailsModal;
