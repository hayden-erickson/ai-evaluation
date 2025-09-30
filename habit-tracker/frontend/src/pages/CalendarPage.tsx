import { useState, useEffect } from 'react';
import {
    Box,
    Checkbox,
    Container,
    Flex,
    FormControl,
    FormLabel,
    Heading,
    Stack,
    VStack,
    Text,
    Spinner,
} from '@chakra-ui/react';
import type { Habit } from '../types';
import { HabitCalendar } from '../components/HabitCalendar';
import { getHabits } from '../services/api';

export const CalendarPage = () => {
    const [habits, setHabits] = useState<Habit[]>([]);
    const [selectedHabitIds, setSelectedHabitIds] = useState<number[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchHabits = async () => {
            try {
                const fetchedHabits = await getHabits();
                setHabits(fetchedHabits);
                // By default, select all habits
                setSelectedHabitIds(fetchedHabits.map(h => h.id));
            } catch (error) {
                console.error('Failed to fetch habits:', error);
            } finally {
                setIsLoading(false);
            }
        };
        
        fetchHabits();
    }, []);

    const handleHabitToggle = (habitId: number) => {
        setSelectedHabitIds(prev =>
            prev.includes(habitId)
                ? prev.filter(id => id !== habitId)
                : [...prev, habitId]
        );
    };

    return (
        <Container maxW="container.xl" py={8}>
            <Stack spacing={8}>
                <Heading size="lg">Habit Calendar</Heading>
                {isLoading ? (
                    <Flex justify="center" align="center" h="200px">
                        <Spinner size="xl" />
                    </Flex>
                ) : (
                    <Flex gap={8} direction={{ base: 'column', md: 'row' }}>
                        <Box flex="1" maxW={{ base: '100%', md: '300px' }}>
                            <VStack align="stretch" spacing={4}>
                                <FormControl>
                                    <FormLabel>Filter Habits</FormLabel>
                                    {habits.length === 0 ? (
                                        <Text color="gray.500">No habits found</Text>
                                    ) : (
                                        habits.map(habit => (
                                            <Checkbox
                                                key={habit.id}
                                                isChecked={selectedHabitIds.includes(habit.id)}
                                                onChange={() => handleHabitToggle(habit.id)}
                                            >
                                                {habit.name}
                                            </Checkbox>
                                        ))
                                    )}
                                </FormControl>
                            </VStack>
                        </Box>
                        <Box flex="3">
                            <HabitCalendar
                                habits={habits}
                                selectedHabitIds={selectedHabitIds}
                            />
                        </Box>
                    </Flex>
                )}
            </Stack>
        </Container>
    );
};