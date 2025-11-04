/**
 * LogDetailsModal Component
 * Modal for creating and editing habit logs
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
} from 'react-native';
import {Log, CreateLogRequest, UpdateLogRequest} from '../types';
import {colors, spacing, borderRadius, fontSize, shadow} from '../styles/theme';
import {formatDate} from '../utils/helpers';

interface LogDetailsModalProps {
  visible: boolean;
  log?: Log | null; // If provided, we're editing; otherwise creating
  logDate?: Date; // For new logs, the date the log is for
  onClose: () => void;
  onSave: (data: CreateLogRequest | UpdateLogRequest) => Promise<void>;
}

/**
 * Modal for creating or editing a habit log
 */
export function LogDetailsModal({
  visible,
  log,
  logDate,
  onClose,
  onSave,
}: LogDetailsModalProps) {
  const [notes, setNotes] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isEditing = !!log;

  // Initialize form with log data when editing
  useEffect(() => {
    if (log) {
      setNotes(log.notes || '');
    } else {
      setNotes('');
    }
    setError(null);
  }, [log, visible]);

  /**
   * Validates and saves the log
   */
  const handleSave = async () => {
    setError(null);
    setIsLoading(true);

    try {
      if (isEditing) {
        // Update log
        const updateData: UpdateLogRequest = {};
        if (notes !== (log?.notes || '')) {
          updateData.notes = notes.trim() || undefined;
        }

        await onSave(updateData);
      } else {
        // Create new log
        const createData: CreateLogRequest = {
          notes: notes.trim() || undefined,
        };

        await onSave(createData);
      }

      // Reset form and close
      setNotes('');
      setError(null);
      onClose();
    } catch (err: any) {
      setError(err.message || 'Failed to save log. Please try again.');
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

  // Determine the display date
  const displayDate = log
    ? new Date(log.created_at)
    : logDate || new Date();

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
              {isEditing ? 'Edit Log' : 'New Log'}
            </Text>

            <Text style={styles.dateText}>{formatDate(displayDate)}</Text>

            <ScrollView
              style={styles.formContainer}
              showsVerticalScrollIndicator={false}>
              {/* Notes Field */}
              <View style={styles.fieldContainer}>
                <Text style={styles.fieldLabel}>Notes (optional)</Text>
                <TextInput
                  style={[styles.input, styles.textArea]}
                  placeholder="Add notes about this log entry..."
                  placeholderTextColor={colors.textMuted}
                  value={notes}
                  onChangeText={setNotes}
                  editable={!isLoading}
                  multiline={true}
                  numberOfLines={6}
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
    maxHeight: '70%',
    ...shadow.lg,
  },
  modalTitle: {
    fontSize: fontSize.xl,
    fontWeight: '600',
    color: colors.text,
    marginBottom: spacing.sm,
  },
  dateText: {
    fontSize: fontSize.md,
    color: colors.textLight,
    marginBottom: spacing.lg,
  },
  formContainer: {
    flex: 1,
  },
  fieldContainer: {
    marginBottom: spacing.md,
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
    minHeight: 120,
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
