/**
 * LogDetailsModal - Modal for creating or editing habit logs
 * Allows users to add notes when logging a habit completion
 */

import React, { useState, useEffect } from 'react';
import {
  Modal,
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
} from 'react-native';
import { Log, CreateLogRequest, UpdateLogRequest } from '../types';
import { Colors, Spacing, BorderRadius, Typography, CommonStyles } from '../styles/theme';

interface LogDetailsModalProps {
  visible: boolean;
  onClose: () => void;
  onSave: (data: CreateLogRequest | UpdateLogRequest) => Promise<void>;
  log?: Log | null;  // If provided, we're editing; otherwise creating
}

export default function LogDetailsModal({
  visible,
  onClose,
  onSave,
  log,
}: LogDetailsModalProps) {
  const [notes, setNotes] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  // Update form when log changes
  useEffect(() => {
    if (log) {
      setNotes(log.notes || '');
    } else {
      setNotes('');
    }
    setError('');
  }, [log, visible]);

  /**
   * Handle save button press
   * Validates input and calls onSave callback
   */
  const handleSave = async () => {
    try {
      setError('');
      setIsLoading(true);

      const data: CreateLogRequest | UpdateLogRequest = {
        notes: notes.trim() || undefined,
      };

      await onSave(data);
      
      // Reset form and close modal on success
      setNotes('');
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save log');
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Handle cancel button press
   * Resets form and closes modal
   */
  const handleCancel = () => {
    setNotes('');
    setError('');
    onClose();
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={handleCancel}
    >
      <View style={CommonStyles.modalContainer}>
        <View style={[CommonStyles.modalContent, styles.modalContent]}>
          <Text style={styles.title}>
            {log ? 'Edit Log' : 'New Log'}
          </Text>

          {/* Notes input field */}
          <View style={styles.inputContainer}>
            <Text style={styles.label}>Notes (optional)</Text>
            <TextInput
              style={[CommonStyles.input, styles.textArea]}
              placeholder="Add notes about this habit completion..."
              placeholderTextColor={Colors.textSecondary}
              value={notes}
              onChangeText={setNotes}
              multiline
              numberOfLines={4}
              textAlignVertical="top"
            />
          </View>

          {/* Error message */}
          {error ? <Text style={CommonStyles.errorText}>{error}</Text> : null}

          {/* Action buttons */}
          <View style={styles.buttonContainer}>
            <TouchableOpacity
              style={[styles.button, styles.cancelButton]}
              onPress={handleCancel}
              disabled={isLoading}
            >
              <Text style={styles.cancelButtonText}>Cancel</Text>
            </TouchableOpacity>

            <TouchableOpacity
              style={[styles.button, styles.saveButton]}
              onPress={handleSave}
              disabled={isLoading}
            >
              {isLoading ? (
                <ActivityIndicator color={Colors.textPrimary} />
              ) : (
                <Text style={styles.saveButtonText}>Save</Text>
              )}
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  modalContent: {
    maxHeight: '80%',
  },
  
  title: {
    fontSize: Typography.h3,
    fontWeight: Typography.bold,
    color: Colors.textPrimary,
    marginBottom: Spacing.lg,
  },
  
  inputContainer: {
    marginBottom: Spacing.md,
  },
  
  label: {
    fontSize: Typography.small,
    fontWeight: Typography.medium,
    color: Colors.textPrimary,
    marginBottom: Spacing.xs,
  },
  
  textArea: {
    height: 100,
    paddingTop: Spacing.md,
  },
  
  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginTop: Spacing.lg,
  },
  
  button: {
    flex: 1,
    paddingVertical: Spacing.md,
    borderRadius: BorderRadius.md,
    alignItems: 'center',
    justifyContent: 'center',
    marginHorizontal: Spacing.xs,
  },
  
  cancelButton: {
    backgroundColor: Colors.inputBackground,
    borderWidth: 1,
    borderColor: Colors.border,
  },
  
  cancelButtonText: {
    color: Colors.textPrimary,
    fontSize: Typography.body,
    fontWeight: Typography.medium,
  },
  
  saveButton: {
    backgroundColor: Colors.primary,
  },
  
  saveButtonText: {
    color: Colors.textPrimary,
    fontSize: Typography.body,
    fontWeight: Typography.semibold,
  },
});
