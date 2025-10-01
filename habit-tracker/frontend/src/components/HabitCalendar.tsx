import { useState, useEffect } from 'react';
import {
    Box,
    Grid,
    Text,
    VStack,
} from '@chakra-ui/react';
import Calendar from 'react-calendar';
import { format } from 'date-fns';
import type { Habit, Log } from '../types';
import { getLogs } from '../services/api';
import 'react-calendar/dist/Calendar.css';

interface HabitCalendarProps {
    habits: Habit[];
    selectedHabitIds: number[];
}

interface LogsByDate {
    [date: string]: Log[];
}

export const HabitCalendar = ({ habits, selectedHabitIds }: HabitCalendarProps) => {
    const [date, setDate] = useState<Date>(new Date());
    const [logsByDate, setLogsByDate] = useState<LogsByDate>({});

    // Color generation for habits
    const habitColors = habits.reduce((acc, habit) => {
        const hue = (habit.id * 137.508) % 360; // Golden angle approximation
        acc[habit.id] = `hsl(${hue}, 70%, 50%)`;
        return acc;
    }, {} as { [id: number]: string });

    // Fetch logs for selected habits
    useEffect(() => {
        const fetchLogs = async () => {
            if (selectedHabitIds.length === 0) {
                setLogsByDate({});
                return;
            }
            
            try {
                // Fetch logs for the selected habits
                const promises = selectedHabitIds.map(habitId => getLogs(habitId));
                const allLogs = (await Promise.all(promises)).flat();
                
                // Group logs by date
                const logsByDateMap = allLogs.reduce((acc: LogsByDate, log: Log) => {
                    const date = new Date(log.completedAt);
                    if (!isNaN(date.getTime())) {
                        const dateStr = format(date, 'yyyy-MM-dd');
                        if (!acc[dateStr]) {
                            acc[dateStr] = [];
                        }
                        acc[dateStr].push(log);
                    }
                    return acc;
                }, {});
                
                setLogsByDate(logsByDateMap);
            } catch (error) {
                console.error('Failed to fetch logs:', error);
            }
        };

        fetchLogs();
    }, [selectedHabitIds]);

    // Render habit completion dots for each date
    const renderTileContent = ({ date }: { date: Date }) => {
        if (!date || isNaN(date.getTime())) return null;
        
        const dateStr = format(date, 'yyyy-MM-dd');
        const logs = logsByDate[dateStr] || [];
        const filteredLogs = logs.filter(log => selectedHabitIds.includes(log.habitId));

        if (filteredLogs.length === 0) return null;

        return (
            <Grid
                templateColumns={`repeat(${Math.min(filteredLogs.length, 3)}, 1fr)`}
                gap={1}
                mt={1}
            >
                {filteredLogs.slice(0, 3).map((log, idx) => (
                    <Box
                        key={idx}
                        w="6px"
                        h="6px"
                        borderRadius="full"
                        bg={habitColors[log.habitId]}
                    />
                ))}
                {filteredLogs.length > 3 && (
                    <Text fontSize="xs" color="gray.500">
                        +{filteredLogs.length - 3}
                    </Text>
                )}
            </Grid>
        );
    };

    // Custom tile className function
    const getTileClassName = ({ date, view }: { date: Date; view: string }) => {
        if (view !== 'month' || !date || isNaN(date.getTime())) return '';
        
        const dateStr = format(date, 'yyyy-MM-dd');
        const logs = logsByDate[dateStr] || [];
        const filteredLogs = logs.filter(log => selectedHabitIds.includes(log.habitId));
        
        return filteredLogs.length > 0 ? 'has-logs' : '';
    };

    return (
        <Box>
            <VStack spacing={4} align="stretch">
                <Calendar
                    onChange={(value) => {
                        if (value instanceof Date && !isNaN(value.getTime())) {
                            setDate(value);
                        }
                    }}
                    value={date}
                    tileContent={renderTileContent}
                    tileClassName={getTileClassName}
                    className="react-calendar"
                />
                {selectedHabitIds.length === 0 && (
                    <Text color="gray.500" textAlign="center">
                        Select one or more habits to view their completion history
                    </Text>
                )}
            </VStack>
        </Box>
    );
};