/**
 * Modal for creating/editing habit details
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
import { Habit, CreateHabitRequest, UpdateHabitRequest } from '../types/models';
import { theme } from '../styles/theme';
import { commonStyles } from '../styles/commonStyles';

interface HabitDetailsModalProps {
  visible: boolean;
  habit: Habit | null; // null for creating new habit
  onClose: () => void;
  onSave: (data: CreateHabitRequest | UpdateHabitRequest) => Promise<void>;
}

/**
 * Modal for creating or editing habit details
 */
export const HabitDetailsModal: React.FC<HabitDetailsModalProps> = ({
  visible,
  habit,
  onClose,
  onSave,
}) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const isEditing = habit !== null;

  /**
   * Reset form when modal opens/closes or habit changes
   */
  useEffect(() => {
    if (visible) {
      if (habit) {
        setName(habit.name);
        setDescription(habit.description || '');
      } else {
        setName('');
        setDescription('');
      }
      setError('');
    }
  }, [visible, habit]);

  /**
   * Validate form inputs
   */
  const validateForm = (): boolean => {
    if (!name.trim()) {
      setError('Habit name is required');
      return false;
    }
    return true;
  };

  /**
   * Handle save button press
   */
  const handleSave = async () => {
    setError('');

    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    try {
      if (isEditing) {
        // Update existing habit
        const updateData: UpdateHabitRequest = {
          name: name.trim(),
          description: description.trim() || undefined,
        };
        await onSave(updateData);
      } else {
        // Create new habit
        const createData: CreateHabitRequest = {
          name: name.trim(),
          description: description.trim() || undefined,
        };
        await onSave(createData);
      }
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save habit');
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
                {isEditing ? 'Edit Habit' : 'New Habit'}
              </Text>

              <View style={styles.form}>
                <View style={styles.inputGroup}>
                  <Text style={commonStyles.label}>Name *</Text>
                  <TextInput
                    style={commonStyles.input}
                    placeholder="e.g., Morning meditation"
                    placeholderTextColor={theme.colors.textSecondary}
                    value={name}
                    onChangeText={(text) => {
                      setName(text);
                      setError('');
                    }}
                    editable={!isLoading}
                  />
                </View>

                <View style={styles.inputGroup}>
                  <Text style={commonStyles.label}>Description</Text>
                  <TextInput
                    style={[commonStyles.input, styles.textArea]}
                    placeholder="Optional description"
                    placeholderTextColor={theme.colors.textSecondary}
                    value={description}
                    onChangeText={setDescription}
                    multiline
                    numberOfLines={3}
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
                      {isEditing ? 'Update' : 'Create'}
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
    height: 80,
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

