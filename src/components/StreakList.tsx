import React from 'react';
import { ScrollView, StyleSheet, View } from 'react-native';
import DayContainer from './DayContainer';
import { indexLogsByDay, dayKey } from '../utils/streak';
import type { LogEntry } from '../types';

interface Props {
  logs: LogEntry[];
  days?: number;
  onEditLogById?: (id: number) => void;
}

/** Horizontal list of day cells indicating which days have logs. */
export default function StreakList({ logs, days = 14, onEditLogById }: Props) {
  const map = indexLogsByDay(logs);
  const today = new Date();
  const cells: { label: string; key: string; logId?: number }[] = [];

  for (let i = 0; i < days; i++) {
    const d = new Date(today);
    d.setDate(today.getDate() - i);
    const key = dayKey(d);
    const label = d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
    cells.push({ label, key, logId: map[key]?.id });
  }

  return (
    <ScrollView horizontal showsHorizontalScrollIndicator={false} contentContainerStyle={styles.row}>
      {cells.map((c) => (
        <DayContainer key={c.key} label={c.label} hasLog={!!c.logId} onPress={c.logId && onEditLogById ? () => onEditLogById(c.logId!) : undefined} />
      ))}
      <View style={{ width: 4 }} />
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  row: { paddingVertical: 4, paddingHorizontal: 2 },
});
