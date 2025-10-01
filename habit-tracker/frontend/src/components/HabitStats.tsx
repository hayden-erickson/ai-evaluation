import { useEffect, useState } from 'react';
import {
    Box,
    Grid,
    Stat,
    StatLabel,
    StatNumber,
    StatHelpText,
    Text,
    useColorModeValue,
} from '@chakra-ui/react';
import type { Habit, Log } from '../types';
import { getLogs } from '../services/api';

interface HabitStatsProps {
    habit: Habit;
}

interface Stats {
    currentStreak: number;
    longestStreak: number;
    completionRate: number;
    totalCompletions: number;
    lastCompleted: string | null;
}

export const HabitStats = ({ habit }: HabitStatsProps) => {
    const [stats, setStats] = useState<Stats>({
        currentStreak: 0,
        longestStreak: 0,
        completionRate: 0,
        totalCompletions: 0,
        lastCompleted: null,
    });
    const [isLoading, setIsLoading] = useState(true);

    const bgColor = useColorModeValue('white', 'gray.800');
    const borderColor = useColorModeValue('gray.200', 'gray.700');

    useEffect(() => {
        const fetchStats = async () => {
            try {
                // Fetch all logs for this habit
                const logs = await getLogs(habit.id);
                
                // Sort logs by date
                const sortedLogs = logs.sort((a, b) => {
                    const dateA = new Date(a.completedAt);
                    const dateB = new Date(b.completedAt);
                    if (isNaN(dateA.getTime()) || isNaN(dateB.getTime())) return 0;
                    return dateA.getTime() - dateB.getTime();
                });

                // Calculate completion rate
                const createdDate = new Date(habit.createdAt);
                const daysActive = !isNaN(createdDate.getTime()) ? Math.max(1, Math.floor(
                    (new Date().getTime() - createdDate.getTime()) / (1000 * 60 * 60 * 24)
                )) : 1;
                const completionRate = (logs.length / daysActive) * 100;

                // Find the last completion date
                const lastCompleted = sortedLogs.length > 0 ? (() => {
                    const date = new Date(sortedLogs[sortedLogs.length - 1].completedAt);
                    return isNaN(date.getTime()) ? null : date.toLocaleDateString();
                })() : null;

                // Calculate streaks
                let currentStreak = 0;
                let longestStreak = 0;
                let tempStreak = 0;
                
                // Convert logs to a Set of date strings for easier lookup
                const completedDates = new Set(
                    logs.map(log => {
                        const date = new Date(log.completedAt);
                        return !isNaN(date.getTime()) ? date.toISOString().split('T')[0] : null;
                    }).filter((date): date is string => date !== null)
                );

                // Check the last 30 days for current streak
                const today = new Date();
                for (let i = 0; i < 30; i++) {
                    const date = new Date(today);
                    date.setDate(today.getDate() - i);
                    const dateStr = date.toISOString().split('T')[0];
                    
                    if (completedDates.has(dateStr)) {
                        currentStreak = i + 1;
                    } else {
                        break;
                    }
                }

                // Calculate longest streak
                const dateArray = Array.from(completedDates).sort();
                dateArray.forEach((date, i) => {
                    if (i === 0) {
                        tempStreak = 1;
                    } else {
                        const curr = new Date(date);
                        const prev = new Date(dateArray[i - 1]);
                        const diffDays = Math.floor((curr.getTime() - prev.getTime()) / (1000 * 60 * 60 * 24));
                        
                        if (diffDays === 1) {
                            tempStreak++;
                        } else {
                            if (tempStreak > longestStreak) {
                                longestStreak = tempStreak;
                            }
                            tempStreak = 1;
                        }
                    }
                });

                // Check final streak
                if (tempStreak > longestStreak) {
                    longestStreak = tempStreak;
                }

                setStats({
                    currentStreak,
                    longestStreak,
                    completionRate: Math.round(completionRate * 10) / 10,
                    totalCompletions: logs.length,
                    lastCompleted,
                });
            } catch (error) {
                console.error('Failed to fetch habit stats:', error);
            } finally {
                setIsLoading(false);
            }
        };

        fetchStats();
    }, [habit]);

    return (
        <Grid
            templateColumns={{
                base: 'repeat(1, 1fr)',
                sm: 'repeat(2, 1fr)',
                md: 'repeat(3, 1fr)',
            }}
            gap={4}
        >
            <Box p={4} borderRadius="lg" bg={bgColor} borderWidth="1px" borderColor={borderColor}>
                <Stat>
                    <StatLabel>Current Streak</StatLabel>
                    <StatNumber>{stats.currentStreak}</StatNumber>
                    <StatHelpText>days</StatHelpText>
                </Stat>
            </Box>

            <Box p={4} borderRadius="lg" bg={bgColor} borderWidth="1px" borderColor={borderColor}>
                <Stat>
                    <StatLabel>Longest Streak</StatLabel>
                    <StatNumber>{stats.longestStreak}</StatNumber>
                    <StatHelpText>days</StatHelpText>
                </Stat>
            </Box>

            <Box p={4} borderRadius="lg" bg={bgColor} borderWidth="1px" borderColor={borderColor}>
                <Stat>
                    <StatLabel>Completion Rate</StatLabel>
                    <StatNumber>{stats.completionRate}%</StatNumber>
                    <StatHelpText>since starting</StatHelpText>
                </Stat>
            </Box>

            <Box p={4} borderRadius="lg" bg={bgColor} borderWidth="1px" borderColor={borderColor}>
                <Stat>
                    <StatLabel>Total Completions</StatLabel>
                    <StatNumber>{stats.totalCompletions}</StatNumber>
                    <StatHelpText>times completed</StatHelpText>
                </Stat>
            </Box>

            <Box p={4} borderRadius="lg" bg={bgColor} borderWidth="1px" borderColor={borderColor}>
                <Stat>
                    <StatLabel>Last Completed</StatLabel>
                    <StatNumber>
                        <Text fontSize="lg">
                            {stats.lastCompleted || 'Never'}
                        </Text>
                    </StatNumber>
                    <StatHelpText>most recent</StatHelpText>
                </Stat>
            </Box>
        </Grid>
    );
};