import { useState, useEffect } from 'react';
import {
    Box,
    Button,
    Table,
    Thead,
    Tbody,
    Tr,
    Th,
    Td,
    IconButton,
    useDisclosure,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalBody,
    ModalFooter,
    FormControl,
    FormLabel,
    Input,
    Textarea,
    VStack,
    Flex,
    Text,
    useToast,
} from '@chakra-ui/react';
import { AddIcon, DeleteIcon, EditIcon } from '@chakra-ui/icons';
import { format } from 'date-fns';
import type { Habit, Log } from '../types';
import { getLogs, createLog, updateLog, deleteLog } from '../services/api';

interface HabitLogListProps {
    habit: Habit;
    selectedTags?: string[];
}

interface LogFormData {
    id?: number;
    notes: string;
    completedAt: string;
}

export const HabitLogList = ({ habit }: HabitLogListProps) => {
    const [logs, setLogs] = useState<Log[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [selectedLog, setSelectedLog] = useState<Log | null>(null);
    const { isOpen, onOpen, onClose } = useDisclosure();
    const toast = useToast();

    const [formData, setFormData] = useState<LogFormData>({
        notes: '',
        completedAt: format(new Date(), 'yyyy-MM-dd'),
    });

    const fetchLogs = async () => {
        try {
            const fetchedLogs = await getLogs(habit.id);
            const sortedLogs = fetchedLogs.sort((a, b) => 
                new Date(b.completedAt).getTime() - new Date(a.completedAt).getTime()
            );
            
            // Filter logs by selected tags if any
            const filteredLogs = selectedTags?.length
                ? sortedLogs.filter(log => 
                    log.tags?.some(tag => selectedTags.includes(tag.value))
                )
                : sortedLogs;
            
            setLogs(filteredLogs);
        } catch (error) {
            console.error('Failed to fetch logs:', error);
            toast({
                title: 'Error',
                description: 'Failed to load habit logs',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchLogs();
    }, [habit.id]);

    const handleOpenModal = (log?: Log) => {
        if (log) {
            setSelectedLog(log);
            setFormData({
                id: log.id,
                notes: log.notes,
                completedAt: format(new Date(log.completedAt), 'yyyy-MM-dd'),
            });
        } else {
            setSelectedLog(null);
            setFormData({
                notes: '',
                completedAt: format(new Date(), 'yyyy-MM-dd'),
            });
        }
        onOpen();
    };

    const handleSubmit = async () => {
        try {
            if (selectedLog) {
                await updateLog(habit.id, selectedLog.id, formData);
            } else {
                await createLog(habit.id, formData);
            }
            await fetchLogs();
            onClose();
            toast({
                title: 'Success',
                description: selectedLog ? 'Log updated' : 'Log created',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
        } catch (error) {
            console.error('Failed to save log:', error);
            toast({
                title: 'Error',
                description: 'Failed to save log',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    const handleDelete = async (log: Log) => {
        try {
            await deleteLog(habit.id, log.id);
            await fetchLogs();
            toast({
                title: 'Success',
                description: 'Log deleted',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
        } catch (error) {
            console.error('Failed to delete log:', error);
            toast({
                title: 'Error',
                description: 'Failed to delete log',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    return (
        <Box>
            <Flex justify="space-between" align="center" mb={4}>
                <Text fontSize="lg" fontWeight="medium">Habit Logs</Text>
                <Button
                    leftIcon={<AddIcon />}
                    colorScheme="blue"
                    onClick={() => handleOpenModal()}
                >
                    Add Log
                </Button>
            </Flex>

            <Table variant="simple">
                <Thead>
                    <Tr>
                        <Th>Date</Th>
                        <Th>Notes</Th>
                        <Th>Actions</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {logs.map(log => (
                        <Tr key={log.id}>
                            <Td>{format(new Date(log.completedAt), 'MMM d, yyyy')}</Td>
                            <Td>{log.notes || '-'}</Td>
                            <Td>
                                <IconButton
                                    icon={<EditIcon />}
                                    aria-label="Edit log"
                                    size="sm"
                                    mr={2}
                                    onClick={() => handleOpenModal(log)}
                                />
                                <IconButton
                                    icon={<DeleteIcon />}
                                    aria-label="Delete log"
                                    size="sm"
                                    colorScheme="red"
                                    onClick={() => handleDelete(log)}
                                />
                            </Td>
                        </Tr>
                    ))}
                    {logs.length === 0 && (
                        <Tr>
                            <Td colSpan={3} textAlign="center">
                                <Text color="gray.500">No logs found</Text>
                            </Td>
                        </Tr>
                    )}
                </Tbody>
            </Table>

            <Modal isOpen={isOpen} onClose={onClose}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>
                        {selectedLog ? 'Edit Log' : 'Add Log'}
                    </ModalHeader>
                    <ModalBody>
                        <VStack spacing={4}>
                            <FormControl>
                                <FormLabel>Date</FormLabel>
                                <Input
                                    type="date"
                                    value={formData.completedAt}
                                    onChange={e => setFormData(prev => ({
                                        ...prev,
                                        completedAt: e.target.value
                                    }))}
                                />
                            </FormControl>
                            <FormControl>
                                <FormLabel>Notes</FormLabel>
                                <Textarea
                                    value={formData.notes}
                                    onChange={e => setFormData(prev => ({
                                        ...prev,
                                        notes: e.target.value
                                    }))}
                                    placeholder="Add any notes about this completion..."
                                />
                            </FormControl>
                        </VStack>
                    </ModalBody>
                    <ModalFooter>
                        <Button variant="ghost" mr={3} onClick={onClose}>
                            Cancel
                        </Button>
                        <Button colorScheme="blue" onClick={handleSubmit}>
                            {selectedLog ? 'Save' : 'Add'}
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </Box>
    );
};