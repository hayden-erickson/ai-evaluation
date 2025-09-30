import { useState, useEffect } from 'react';
import {
    Box,
    Button,
    Container,
    Flex,
    Grid,
    Heading,
    Stack,
    Text,
    useDisclosure,
} from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';
import { getHabits } from '../services/api';
import type { Habit } from '../types';
import { HabitForm } from '../components/HabitForm';

export const MainPage = () => {
    const [habits, setHabits] = useState<Habit[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const { isOpen, onOpen, onClose } = useDisclosure();
    const navigate = useNavigate();

    const fetchHabits = async () => {
        try {
            const fetchedHabits = await getHabits();
            setHabits(fetchedHabits);
        } catch (error) {
            console.error('Failed to fetch habits:', error);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchHabits();
    }, []);

    const handleHabitClick = (habitId: number) => {
        navigate(`/habits/${habitId}`);
    };

    const handleHabitAdded = () => {
        fetchHabits();
        onClose();
    };

    return (
        <Container maxW="container.xl" py={8}>
            <Flex justify="space-between" align="center" mb={8}>
                <Heading size="lg">My Habits</Heading>
                <Flex gap={4}>
                    <Button colorScheme="blue" onClick={() => navigate('/calendar')}>
                        Calendar View
                    </Button>
                    <Button colorScheme="teal" onClick={onOpen}>
                        Add New Habit
                    </Button>
                </Flex>
            </Flex>

            {isLoading ? (
                <Text>Loading habits...</Text>
            ) : habits.length === 0 ? (
                <Box textAlign="center" py={8}>
                    <Text fontSize="lg" color="gray.600">
                        No habits yet. Create your first habit to get started!
                    </Text>
                </Box>
            ) : (
                <Grid templateColumns="repeat(auto-fill, minmax(300px, 1fr))" gap={6}>
                    {habits.map((habit) => (
                        <Box
                            key={habit.id}
                            p={5}
                            shadow="md"
                            borderWidth="1px"
                            borderRadius="lg"
                            cursor="pointer"
                            onClick={() => handleHabitClick(habit.id)}
                            _hover={{ shadow: 'lg', transform: 'translateY(-2px)' }}
                            transition="all 0.2s"
                        >
                            <Stack spacing={4}>
                                <Heading size="md">{habit.name}</Heading>
                                <Text>{habit.description}</Text>
                                <Text color="gray.600">Every {habit.frequency} days</Text>
                            </Stack>
                        </Box>
                    ))}
                </Grid>
            )}

            <HabitForm isOpen={isOpen} onClose={onClose} onSuccess={handleHabitAdded} />
        </Container>
    );
};