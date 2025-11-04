/**
 * HabitDetailsModal Component
 * Modal for creating and editing habits
 */

import React, {useState, useEffect} from 'react';
import {
  Modal,
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  Alert,
} from 'react-native';
import {Habit, CreateHabitRequest, UpdateHabitRequest} from '../types';
import {colors, spacing, borderRadius, fontSize, shadow} from '../styles/theme';
import {validateHabitName} from '../utils/helpers';

interface HabitDetailsModalProps {
  visible: boolean;
  habit?: Habit | null; // If provided, we're editing; otherwise creating
  onClose: () => void;
  onSave: (data: CreateHabitRequest | UpdateHabitRequest) => Promise<void>;
}

/**
 * Modal for creating or editing a habit
 */
export function HabitDetailsModal({
  visible,
  habit,
  onClose,
  onSave,
}: HabitDetailsModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isEditing = !!habit;

  // Initialize form with habit data when editing
  useEffect(() => {
    if (habit) {
      setName(habit.name);
      setDescription(habit.description || '');
    } else {
      setName('');
      setDescription('');
    }
    setError(null);
  }, [habit, visible]);

  /**
   * Validates and saves the habit
   */
  const handleSave = async () => {
    // Clear previous errors
    setError(null);

    // Validate name
    const nameValidation = validateHabitName(name);
    if (!nameValidation.isValid) {
      setError(nameValidation.message || 'Invalid habit name');
      return;
    }

    setIsLoading(true);

    try {
      if (isEditing) {
        // Update habit - only include changed fields
        const updateData: UpdateHabitRequest = {};
        if (name !== habit?.name) {
          updateData.name = name;
        }
        if (description !== (habit?.description || '')) {
          updateData.description = description;
        }

        await onSave(updateData);
      } else {
        // Create new habit
        const createData: CreateHabitRequest = {
          name: name.trim(),
          description: description.trim() || undefined,
        };

        await onSave(createData);
      }

      // Reset form and close
      setName('');
      setDescription('');
      setError(null);
      onClose();
    } catch (err: any) {
      setError(err.message || 'Failed to save habit. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Handles modal close
   */
  const handleClose = () => {
    if (!isLoading) {
      setError(null);
      onClose();
    }
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent={true}
      onRequestClose={handleClose}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.modalOverlay}>
        <View style={styles.modalContainer}>
          <View style={styles.modalContent}>
            <Text style={styles.modalTitle}>
              {isEditing ? 'Edit Habit' : 'New Habit'}
            </Text>

            <ScrollView
              style={styles.formContainer}
              showsVerticalScrollIndicator={false}>
              {/* Name Field */}
              <View style={styles.fieldContainer}>
                <Text style={styles.fieldLabel}>Name *</Text>
                <TextInput
                  style={styles.input}
                  placeholder="Enter habit name"
                  placeholderTextColor={colors.textMuted}
                  value={name}
                  onChangeText={setName}
                  editable={!isLoading}
                  maxLength={100}
                />
              </View>

              {/* Description Field */}
              <View style={styles.fieldContainer}>
                <Text style={styles.fieldLabel}>Description (optional)</Text>
                <TextInput
                  style={[styles.input, styles.textArea]}
                  placeholder="Add a description"
                  placeholderTextColor={colors.textMuted}
                  value={description}
                  onChangeText={setDescription}
                  editable={!isLoading}
                  multiline={true}
                  numberOfLines={4}
                  textAlignVertical="top"
                />
              </View>

              {/* Error message */}
              {error && (
                <View style={styles.errorContainer}>
                  <Text style={styles.errorText}>{error}</Text>
                </View>
              )}
            </ScrollView>

            {/* Action Buttons */}
            <View style={styles.buttonContainer}>
              <TouchableOpacity
                style={[styles.button, styles.cancelButton]}
                onPress={handleClose}
                disabled={isLoading}>
                <Text style={styles.cancelButtonText}>Cancel</Text>
              </TouchableOpacity>

              <TouchableOpacity
                style={[
                  styles.button,
                  styles.saveButton,
                  isLoading && styles.disabledButton,
                ]}
                onPress={handleSave}
                disabled={isLoading}>
                <Text style={styles.saveButtonText}>
                  {isLoading ? 'Saving...' : 'Save'}
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
}

const styles = StyleSheet.create({
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'flex-end',
  },
  modalContainer: {
    flex: 1,
    justifyContent: 'flex-end',
  },
  modalContent: {
    backgroundColor: colors.background,
    borderTopLeftRadius: borderRadius.lg,
    borderTopRightRadius: borderRadius.lg,
    paddingTop: spacing.lg,
    paddingBottom: Platform.OS === 'ios' ? spacing.xl : spacing.lg,
    paddingHorizontal: spacing.lg,
    maxHeight: '80%',
    ...shadow.lg,
  },
  modalTitle: {
    fontSize: fontSize.xl,
    fontWeight: '600',
    color: colors.text,
    marginBottom: spacing.lg,
  },
  formContainer: {
    flex: 1,
  },
  fieldContainer: {
    marginBottom: spacing.lg,
  },
  fieldLabel: {
    fontSize: fontSize.sm,
    fontWeight: '500',
    color: colors.text,
    marginBottom: spacing.sm,
  },
  input: {
    backgroundColor: colors.surface,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    fontSize: fontSize.md,
    color: colors.text,
    borderWidth: 1,
    borderColor: colors.border,
  },
  textArea: {
    minHeight: 100,
    paddingTop: spacing.md,
  },
  errorContainer: {
    backgroundColor: colors.error,
    borderRadius: borderRadius.sm,
    padding: spacing.md,
    marginTop: spacing.sm,
  },
  errorText: {
    color: colors.text,
    fontSize: fontSize.sm,
  },
  buttonContainer: {
    flexDirection: 'row',
    gap: spacing.md,
    marginTop: spacing.lg,
  },
  button: {
    flex: 1,
    padding: spacing.md,
    borderRadius: borderRadius.md,
    alignItems: 'center',
  },
  cancelButton: {
    backgroundColor: colors.surface,
    borderWidth: 1,
    borderColor: colors.border,
  },
  cancelButtonText: {
    color: colors.text,
    fontSize: fontSize.md,
    fontWeight: '500',
  },
  saveButton: {
    backgroundColor: colors.primary,
  },
  saveButtonText: {
    color: colors.text,
    fontSize: fontSize.md,
    fontWeight: '600',
  },
  disabledButton: {
    opacity: 0.5,
  },
});
