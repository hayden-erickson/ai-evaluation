/**
 * Modal for creating/editing log details
 */

import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Modal,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  ActivityIndicator,
} from 'react-native';
import { Log, CreateLogRequest, UpdateLogRequest } from '../types/models';
import { theme } from '../styles/theme';
import { commonStyles } from '../styles/commonStyles';

interface LogDetailsModalProps {
  visible: boolean;
  log: Log | null; // null for creating new log
  onClose: () => void;
  onSave: (data: CreateLogRequest | UpdateLogRequest) => Promise<void>;
}

/**
 * Modal for creating or editing log details
 */
export const LogDetailsModal: React.FC<LogDetailsModalProps> = ({
  visible,
  log,
  onClose,
  onSave,
}) => {
  const [notes, setNotes] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const isEditing = log !== null;

  /**
   * Reset form when modal opens/closes or log changes
   */
  useEffect(() => {
    if (visible) {
      if (log) {
        setNotes(log.notes || '');
      } else {
        setNotes('');
      }
      setError('');
    }
  }, [visible, log]);

  /**
   * Handle save button press
   */
  const handleSave = async () => {
    setError('');
    setIsLoading(true);

    try {
      if (isEditing) {
        // Update existing log
        const updateData: UpdateLogRequest = {
          notes: notes.trim() || undefined,
        };
        await onSave(updateData);
      } else {
        // Create new log
        const createData: CreateLogRequest = {
          notes: notes.trim() || undefined,
        };
        await onSave(createData);
      }
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save log');
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Handle cancel button press
   */
  const handleCancel = () => {
    setError('');
    onClose();
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent={true}
      onRequestClose={handleCancel}
    >
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.modalContainer}
      >
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <ScrollView
              keyboardShouldPersistTaps="handled"
              showsVerticalScrollIndicator={false}
            >
              <Text style={styles.modalTitle}>
                {isEditing ? 'Edit Log' : 'New Log'}
              </Text>

              <View style={styles.form}>
                <View style={styles.inputGroup}>
                  <Text style={commonStyles.label}>Notes</Text>
                  <TextInput
                    style={[commonStyles.input, styles.textArea]}
                    placeholder="How did it go today?"
                    placeholderTextColor={theme.colors.textSecondary}
                    value={notes}
                    onChangeText={(text) => {
                      setNotes(text);
                      setError('');
                    }}
                    multiline
                    numberOfLines={5}
                    editable={!isLoading}
                  />
                </View>

                {error ? (
                  <Text style={commonStyles.errorText}>{error}</Text>
                ) : null}
              </View>

              <View style={styles.buttonContainer}>
                <TouchableOpacity
                  style={[styles.button, styles.cancelButton]}
                  onPress={handleCancel}
                  disabled={isLoading}
                >
                  <Text style={styles.cancelButtonText}>Cancel</Text>
                </TouchableOpacity>

                <TouchableOpacity
                  style={[commonStyles.button, styles.saveButton, isLoading && styles.buttonDisabled]}
                  onPress={handleSave}
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <ActivityIndicator color={theme.colors.surface} size="small" />
                  ) : (
                    <Text style={commonStyles.buttonText}>
                      {isEditing ? 'Update' : 'Save'}
                    </Text>
                  )}
                </TouchableOpacity>
              </View>
            </ScrollView>
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
};

const styles = StyleSheet.create({
  modalContainer: {
    flex: 1,
  },

  modalOverlay: {
    flex: 1,
    backgroundColor: theme.colors.overlay,
    justifyContent: 'center',
    alignItems: 'center',
    padding: theme.spacing.lg,
  },

  modalContent: {
    backgroundColor: theme.colors.surface,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.lg,
    width: '100%',
    maxWidth: 500,
    maxHeight: '80%',
    ...theme.shadows.lg,
  },

  modalTitle: {
    fontSize: theme.fontSize.xl,
    fontWeight: 'bold',
    color: theme.colors.text,
    marginBottom: theme.spacing.lg,
    textAlign: 'center',
  },

  form: {
    marginBottom: theme.spacing.lg,
  },

  inputGroup: {
    marginBottom: theme.spacing.md,
  },

  textArea: {
    height: 120,
    textAlignVertical: 'top',
    paddingTop: theme.spacing.sm,
  },

  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: theme.spacing.md,
  },

  button: {
    flex: 1,
    borderRadius: theme.borderRadius.md,
    paddingVertical: theme.spacing.md,
    alignItems: 'center',
    justifyContent: 'center',
  },

  cancelButton: {
    backgroundColor: theme.colors.border,
  },

  cancelButtonText: {
    color: theme.colors.text,
    fontSize: theme.fontSize.md,
    fontWeight: '600',
  },

  saveButton: {
    flex: 1,
  },

  buttonDisabled: {
    opacity: 0.6,
  },
});

