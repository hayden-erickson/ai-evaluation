/**
 * Modal for creating and editing log entries
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
import {Log} from '../types';

interface LogDetailsModalProps {
  visible: boolean;
  log?: Log | null;
  date?: Date;
  onClose: () => void;
  onSave: (data: {notes: string}) => Promise<void>;
}

/**
 * Modal component for creating and editing log entries
 */
export const LogDetailsModal: React.FC<LogDetailsModalProps> = ({
  visible,
  log,
  date,
  onClose,
  onSave,
}) => {
  const [notes, setNotes] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  // Update form when log changes
  useEffect(() => {
    if (log) {
      setNotes(log.notes || '');
    } else {
      setNotes('');
    }
  }, [log, visible]);

  /**
   * Handle save button press
   */
  const handleSave = async () => {
    try {
      setIsSaving(true);
      await onSave({
        notes: notes.trim(),
      });
      onClose();
    } catch (error) {
      console.error('Error saving log:', error);
    } finally {
      setIsSaving(false);
    }
  };

  /**
   * Handle cancel button press
   */
  const handleCancel = () => {
    setNotes('');
    onClose();
  };

  /**
   * Format date for display
   */
  const formatDate = (d: Date): string => {
    const months = [
      'Jan',
      'Feb',
      'Mar',
      'Apr',
      'May',
      'Jun',
      'Jul',
      'Aug',
      'Sep',
      'Oct',
      'Nov',
      'Dec',
    ];
    return `${months[d.getMonth()]} ${d.getDate()}, ${d.getFullYear()}`;
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
                {log ? 'Edit Log Entry' : 'New Log Entry'}
              </Text>
              
              {date && (
                <Text style={styles.dateText}>{formatDate(date)}</Text>
              )}

              <Input
                label="Notes (Optional)"
                value={notes}
                onChangeText={setNotes}
                placeholder="How did it go? Any thoughts?"
                multiline
                numberOfLines={6}
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
    marginBottom: theme.spacing.sm,
    textAlign: 'center',
  },
  dateText: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textLight,
    marginBottom: theme.spacing.lg,
    textAlign: 'center',
  },
  textArea: {
    minHeight: 120,
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
