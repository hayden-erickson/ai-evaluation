import React, {useState, useEffect} from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  Modal,
  StyleSheet,
  Alert,
  KeyboardAvoidingView,
  Platform,
  ActivityIndicator,
} from 'react-native';
import {Log, createLog, updateLog, ApiError} from '../services/api';
import {colors, spacing, borderRadius, typography, shadows} from '../styles/theme';

interface LogDetailsModalProps {
  visible: boolean;
  habitId: number;
  log?: Log;
  onClose: () => void;
  onSave: () => void;
}

/**
 * Modal for creating or editing a habit log
 * Shows input field for notes
 */
export const LogDetailsModal: React.FC<LogDetailsModalProps> = ({
  visible,
  habitId,
  log,
  onClose,
  onSave,
}) => {
  const [notes, setNotes] = useState('');
  const [loading, setLoading] = useState(false);

  const isEditing = !!log;

  // Initialize form with log data when editing
  useEffect(() => {
    if (log) {
      setNotes(log.notes || '');
    } else {
      setNotes('');
    }
  }, [log, visible]);

  /**
   * Handle form submission
   * Creates a new log or updates an existing one
   */
  const handleSave = async () => {
    setLoading(true);

    try {
      if (isEditing && log) {
        // Update existing log
        await updateLog(log.id, {
          notes: notes.trim() || undefined,
        });
      } else {
        // Create new log
        await createLog(habitId, {
          notes: notes.trim() || undefined,
        });
      }

      onSave();
      onClose();
    } catch (error) {
      if (error instanceof ApiError) {
        Alert.alert('Error', error.message);
      } else {
        Alert.alert('Error', 'An unexpected error occurred. Please try again.');
      }
    } finally {
      setLoading(false);
    }
  };

  /**
   * Handle modal close
   * Resets form state
   */
  const handleClose = () => {
    if (!loading) {
      onClose();
    }
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={handleClose}>
      <KeyboardAvoidingView
        style={styles.overlay}
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}>
        <TouchableOpacity
          style={styles.backdrop}
          activeOpacity={1}
          onPress={handleClose}
        />
        <View style={styles.modalContainer}>
          <View style={styles.modal}>
            <Text style={styles.title}>
              {isEditing ? 'Edit Log' : 'Log Habit'}
            </Text>

            <View style={styles.form}>
              <View style={styles.inputContainer}>
                <Text style={styles.label}>Notes (Optional)</Text>
                <TextInput
                  style={[styles.input, styles.textArea]}
                  placeholder="How did it go? Any thoughts?"
                  placeholderTextColor={colors.textLight}
                  value={notes}
                  onChangeText={setNotes}
                  multiline
                  numberOfLines={6}
                  textAlignVertical="top"
                  editable={!loading}
                  autoFocus
                />
              </View>
            </View>

            <View style={styles.buttons}>
              <TouchableOpacity
                style={[styles.button, styles.cancelButton]}
                onPress={handleClose}
                disabled={loading}>
                <Text style={styles.cancelButtonText}>Cancel</Text>
              </TouchableOpacity>

              <TouchableOpacity
                style={[styles.button, styles.saveButton, loading && styles.buttonDisabled]}
                onPress={handleSave}
                disabled={loading}>
                {loading ? (
                  <ActivityIndicator color={colors.textInverse} size="small" />
                ) : (
                  <Text style={styles.saveButtonText}>Save</Text>
                )}
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
};

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  backdrop: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
  },
  modalContainer: {
    width: '90%',
    maxWidth: 400,
  },
  modal: {
    backgroundColor: colors.surfaceLight,
    borderRadius: borderRadius.lg,
    padding: spacing.lg,
    ...shadows.lg,
  },
  title: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.text,
    marginBottom: spacing.lg,
    textAlign: 'center',
  },
  form: {
    marginBottom: spacing.lg,
  },
  inputContainer: {
    marginBottom: spacing.md,
  },
  label: {
    fontSize: typography.sm,
    fontWeight: typography.medium,
    color: colors.text,
    marginBottom: spacing.sm,
  },
  input: {
    backgroundColor: colors.background,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    fontSize: typography.base,
    color: colors.text,
  },
  textArea: {
    minHeight: 120,
  },
  buttons: {
    flexDirection: 'row',
    gap: spacing.md,
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
    fontSize: typography.base,
    fontWeight: typography.semibold,
  },
  saveButton: {
    backgroundColor: colors.success,
    ...shadows.sm,
  },
  saveButtonText: {
    color: colors.text,
    fontSize: typography.base,
    fontWeight: typography.semibold,
  },
  buttonDisabled: {
    opacity: 0.6,
  },
});
