import React from 'react';
import { SafeAreaView, StatusBar, StyleSheet, View } from 'react-native';
import { colors } from '../theme';
import { HabitList } from '../components/HabitList';

export const HabitsScreen: React.FC = () => {
	return (
		<SafeAreaView style={styles.safe}>
			<StatusBar barStyle={'dark-content'} />
			<View style={{ flex: 1 }}>
				<HabitList />
			</View>
		</SafeAreaView>
	);
};

const styles = StyleSheet.create({
	safe: { flex: 1, backgroundColor: colors.background },
});


