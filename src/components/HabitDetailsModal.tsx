/**
 * HabitDetailsModal - Modal for creating or editing habits
 * Allows users to set habit name and description
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
  ScrollView,
} from 'react-native';
import { CreateHabitRequest, UpdateHabitRequest, Habit } from '../types';
import { Colors, Spacing, BorderRadius, Typography, CommonStyles } from '../styles/theme';

interface HabitDetailsModalProps {
  visible: boolean;
  onClose: () => void;
  onSave: (data: CreateHabitRequest | UpdateHabitRequest) => Promise<void>;
  habit?: Habit | null;  // If provided, we're editing; otherwise creating
}

export default function HabitDetailsModal({
  visible,
  onClose,
  onSave,
  habit,
}: HabitDetailsModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  // Update form when habit changes
  useEffect(() => {
    if (habit) {
      setName(habit.name);
      setDescription(habit.description || '');
    } else {
      setName('');
      setDescription('');
    }
    setError('');
  }, [habit, visible]);

  /**
   * Validate form inputs
   * Returns error message if validation fails, empty string otherwise
   */
  const validateForm = (): string => {
    if (!name.trim()) {
      return 'Habit name is required';
    }
    if (name.trim().length < 2) {
      return 'Habit name must be at least 2 characters';
    }
    if (name.trim().length > 100) {
      return 'Habit name must be less than 100 characters';
    }
    return '';
  };

  /**
   * Handle save button press
   * Validates input and calls onSave callback
   */
  const handleSave = async () => {
    try {
      setError('');
      
      // Validate form
      const validationError = validateForm();
      if (validationError) {
        setError(validationError);
        return;
      }

      setIsLoading(true);

      const data: CreateHabitRequest | UpdateHabitRequest = {
        name: name.trim(),
        description: description.trim() || undefined,
      };

      await onSave(data);
      
      // Reset form and close modal on success
      setName('');
      setDescription('');
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save habit');
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Handle cancel button press
   * Resets form and closes modal
   */
  const handleCancel = () => {
    setName('');
    setDescription('');
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
          <ScrollView showsVerticalScrollIndicator={false}>
            <Text style={styles.title}>
              {habit ? 'Edit Habit' : 'New Habit'}
            </Text>

            {/* Name input field */}
            <View style={styles.inputContainer}>
              <Text style={styles.label}>Name *</Text>
              <TextInput
                style={CommonStyles.input}
                placeholder="e.g., Morning Exercise"
                placeholderTextColor={Colors.textSecondary}
                value={name}
                onChangeText={setName}
                maxLength={100}
                autoFocus={!habit}
              />
            </View>

            {/* Description input field */}
            <View style={styles.inputContainer}>
              <Text style={styles.label}>Description (optional)</Text>
              <TextInput
                style={[CommonStyles.input, styles.textArea]}
                placeholder="e.g., 30 minutes of cardio"
                placeholderTextColor={Colors.textSecondary}
                value={description}
                onChangeText={setDescription}
                multiline
                numberOfLines={4}
                textAlignVertical="top"
                maxLength={500}
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
          </ScrollView>
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
    gap: Spacing.md,
  },
  
  button: {
    flex: 1,
    paddingVertical: Spacing.md,
    borderRadius: BorderRadius.md,
    alignItems: 'center',
    justifyContent: 'center',
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
