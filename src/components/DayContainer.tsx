import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { colors, radius, spacing } from '../theme';

interface Props {
  label: string;
  hasLog: boolean;
  onPress?: () => void;
}

/** Visual cell for a single day in the streak view. */
export default function DayContainer({ label, hasLog, onPress }: Props) {
  const content = (
    <View style={[styles.box, hasLog ? styles.filled : styles.empty]}>
      <Text style={[styles.label, hasLog ? styles.filledText : styles.emptyText]}>{label}</Text>
    </View>
  );

  if (onPress && hasLog) {
    return (
      <Pressable onPress={onPress} style={styles.wrap} accessibilityRole="button" accessibilityLabel={`Edit log for ${label}`}>
        {content}
      </Pressable>
    );
  }
  return <View style={styles.wrap}>{content}</View>;
}

const styles = StyleSheet.create({
  wrap: { marginRight: spacing.sm },
  box: { width: 48, height: 56, borderRadius: radius.md, borderWidth: 1, borderColor: colors.border, alignItems: 'center', justifyContent: 'center' },
  filled: { backgroundColor: colors.primary },
  filledText: { color: colors.primaryText, fontWeight: '700' },
  empty: { backgroundColor: colors.surface },
  emptyText: { color: colors.subtext },
  label: { fontSize: 12 },
});
