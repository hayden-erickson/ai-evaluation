import { useState } from 'react';
import {
    Box,
    Button,
    FormControl,
    FormLabel,
    Input,
    NumberInput,
    NumberInputField,
    NumberInputStepper,
    NumberIncrementStepper,
    NumberDecrementStepper,
    Stack,
    Textarea,
    useToast,
    FormHelperText,
} from '@chakra-ui/react';
import type { Habit } from '../types';
import { updateHabit } from '../services/api';

interface HabitSettingsProps {
    habit: Habit;
    onUpdate: (updatedHabit: Habit) => void;
}

interface HabitFormData {
    name: string;
    description: string;
    frequency: number;
}

export const HabitSettings = ({ habit, onUpdate }: HabitSettingsProps) => {
    const [formData, setFormData] = useState<HabitFormData>({
        name: habit.name,
        description: habit.description,
        frequency: habit.frequency,
    });
    const [isSubmitting, setIsSubmitting] = useState(false);
    const toast = useToast();

    const handleSubmit = async () => {
        if (!formData.name.trim()) {
            toast({
                title: 'Error',
                description: 'Habit name is required',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
            return;
        }

        try {
            setIsSubmitting(true);
            const updatedHabit = await updateHabit(habit.id, formData);
            onUpdate(updatedHabit);
            toast({
                title: 'Success',
                description: 'Habit updated successfully',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
        } catch (error) {
            console.error('Failed to update habit:', error);
            toast({
                title: 'Error',
                description: 'Failed to update habit',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <Box maxW="xl">
            <Stack spacing={6}>
                <FormControl isRequired>
                    <FormLabel>Name</FormLabel>
                    <Input
                        value={formData.name}
                        onChange={(e) => setFormData(prev => ({
                            ...prev,
                            name: e.target.value
                        }))}
                        placeholder="Enter habit name"
                    />
                </FormControl>

                <FormControl>
                    <FormLabel>Description</FormLabel>
                    <Textarea
                        value={formData.description}
                        onChange={(e) => setFormData(prev => ({
                            ...prev,
                            description: e.target.value
                        }))}
                        placeholder="Enter habit description"
                    />
                    <FormHelperText>
                        Optional: Add details about your habit
                    </FormHelperText>
                </FormControl>

                <FormControl>
                    <FormLabel>Target Frequency (days per week)</FormLabel>
                    <NumberInput
                        value={formData.frequency}
                        onChange={(_, value) => setFormData(prev => ({
                            ...prev,
                            frequency: value
                        }))}
                        min={1}
                        max={7}
                    >
                        <NumberInputField />
                        <NumberInputStepper>
                            <NumberIncrementStepper />
                            <NumberDecrementStepper />
                        </NumberInputStepper>
                    </NumberInput>
                    <FormHelperText>
                        How many days per week you want to complete this habit
                    </FormHelperText>
                </FormControl>

                <Button
                    colorScheme="blue"
                    onClick={handleSubmit}
                    isLoading={isSubmitting}
                >
                    Save Changes
                </Button>
            </Stack>
        </Box>
    );
};