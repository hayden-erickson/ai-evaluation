import { useState } from 'react';
import {
    Button,
    FormControl,
    FormLabel,
    Input,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    NumberInput,
    NumberInputField,
    Textarea,
    useToast,
} from '@chakra-ui/react';
import { createHabit } from '../services/api';

interface HabitFormProps {
    isOpen: boolean;
    onClose: () => void;
    onSuccess: () => void;
}

export const HabitForm = ({ isOpen, onClose, onSuccess }: HabitFormProps) => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [frequency, setFrequency] = useState('1');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const toast = useToast();

    const handleSubmit = async () => {
        if (!name) {
            toast({
                title: 'Error',
                description: 'Habit name is required',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
            return;
        }

        setIsSubmitting(true);
        try {
            await createHabit({
                name,
                description,
                frequency: parseInt(frequency),
            });
            toast({
                title: 'Success',
                description: 'Habit created successfully',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
            onSuccess();
            handleClose();
        } catch (error) {
            console.error('Failed to create habit:', error);
            toast({
                title: 'Error',
                description: 'Failed to create habit',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleClose = () => {
        setName('');
        setDescription('');
        setFrequency('1');
        onClose();
    };

    return (
        <Modal isOpen={isOpen} onClose={handleClose}>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Create New Habit</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <FormControl isRequired mb={4}>
                        <FormLabel>Name</FormLabel>
                        <Input
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="Enter habit name"
                        />
                    </FormControl>
                    <FormControl mb={4}>
                        <FormLabel>Description</FormLabel>
                        <Textarea
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            placeholder="Enter habit description"
                        />
                    </FormControl>
                    <FormControl isRequired>
                        <FormLabel>Frequency (days)</FormLabel>
                        <NumberInput min={1} value={frequency} onChange={setFrequency}>
                            <NumberInputField />
                        </NumberInput>
                    </FormControl>
                </ModalBody>
                <ModalFooter>
                    <Button variant="ghost" mr={3} onClick={handleClose}>
                        Cancel
                    </Button>
                    <Button
                        colorScheme="teal"
                        onClick={handleSubmit}
                        isLoading={isSubmitting}
                    >
                        Create
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};