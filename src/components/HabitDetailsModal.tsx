/**
 * Modal for creating and editing habits
 */

import React, {useState, useEffect} from 'react';
import {
  View,
  Text,
  Modal,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  TouchableOpacity,
} from 'react-native';
import {Button} from './Button';
import {Input} from './Input';
import {theme} from '../styles/theme';
import {Habit} from '../types';

interface HabitDetailsModalProps {
  visible: boolean;
  habit?: Habit | null;
  onClose: () => void;
  onSave: (data: {name: string; description: string}) => Promise<void>;
}

/**
 * Modal component for creating and editing habits
 */
export const HabitDetailsModal: React.FC<HabitDetailsModalProps> = ({
  visible,
  habit,
  onClose,
  onSave,
}) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [errors, setErrors] = useState<{name?: string}>({});

  // Update form when habit changes
  useEffect(() => {
    if (habit) {
      setName(habit.name);
      setDescription(habit.description || '');
    } else {
      setName('');
      setDescription('');
    }
    setErrors({});
  }, [habit, visible]);

  /**
   * Validate form inputs
   */
  const validate = (): boolean => {
    const newErrors: typeof errors = {};

    if (!name.trim()) {
      newErrors.name = 'Habit name is required';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  /**
   * Handle save button press
   */
  const handleSave = async () => {
    if (!validate()) {
      return;
    }

    try {
      setIsSaving(true);
      await onSave({
        name: name.trim(),
        description: description.trim(),
      });
      onClose();
    } catch (error) {
      console.error('Error saving habit:', error);
    } finally {
      setIsSaving(false);
    }
  };

  /**
   * Handle cancel button press
   */
  const handleCancel = () => {
    setName('');
    setDescription('');
    setErrors({});
    onClose();
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent={true}
      onRequestClose={handleCancel}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.overlay}>
        <TouchableOpacity
          style={styles.backdrop}
          activeOpacity={1}
          onPress={handleCancel}
        />
        <View style={styles.modalContainer}>
          <ScrollView
            contentContainerStyle={styles.scrollContent}
            keyboardShouldPersistTaps="handled">
            <View style={styles.content}>
              <Text style={styles.title}>
                {habit ? 'Edit Habit' : 'New Habit'}
              </Text>

              <Input
                label="Name"
                value={name}
                onChangeText={text => {
                  setName(text);
                  setErrors({...errors, name: undefined});
                }}
                placeholder="e.g., Morning Exercise"
                autoCapitalize="words"
                error={errors.name}
              />

              <Input
                label="Description (Optional)"
                value={description}
                onChangeText={setDescription}
                placeholder="Add details about your habit"
                multiline
                numberOfLines={4}
                style={styles.textArea}
              />

              <View style={styles.buttons}>
                <Button
                  title="Cancel"
                  onPress={handleCancel}
                  variant="secondary"
                  style={styles.button}
                />
                <Button
                  title="Save"
                  onPress={handleSave}
                  loading={isSaving}
                  style={styles.button}
                />
              </View>
            </View>
          </ScrollView>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
};

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    justifyContent: 'flex-end',
  },
  backdrop: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
  },
  modalContainer: {
    backgroundColor: theme.colors.surface,
    borderTopLeftRadius: theme.borderRadius.lg,
    borderTopRightRadius: theme.borderRadius.lg,
    maxHeight: '80%',
    ...theme.shadows.large,
  },
  scrollContent: {
    flexGrow: 1,
  },
  content: {
    padding: theme.spacing.lg,
  },
  title: {
    fontSize: theme.fontSize.xl,
    fontWeight: 'bold',
    color: theme.colors.text,
    marginBottom: theme.spacing.lg,
    textAlign: 'center',
  },
  textArea: {
    minHeight: 100,
    textAlignVertical: 'top',
  },
  buttons: {
    flexDirection: 'row',
    gap: theme.spacing.md,
    marginTop: theme.spacing.lg,
  },
  button: {
    flex: 1,
  },
});
