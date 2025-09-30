import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Box,
    Container,
    Flex,
    Heading,
    Spinner,
    Stack,
    HStack,
    IconButton,
    useToast,
    Tab,
    TabList,
    TabPanel,
    TabPanels,
    Tabs,
    Tag,
    TagLabel,
    Wrap,
    WrapItem,
} from '@chakra-ui/react';
import { ArrowBackIcon } from '@chakra-ui/icons';
import type { Habit } from '../types';
import { getHabit } from '../services/api';
import { HabitStats } from '../components/HabitStats';
import { HabitLogList } from '../components/HabitLogList';
import { TagManagement } from '../components/TagManagement';
import { HabitSettings } from '../components/HabitSettings';

export const HabitDetailPage = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const toast = useToast();
    const [habit, setHabit] = useState<Habit | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [selectedTags, setSelectedTags] = useState<string[]>([]);

    useEffect(() => {
        const fetchHabit = async () => {
            try {
                if (!id) throw new Error('No habit ID provided');
                const habitId = parseInt(id);
                const fetchedHabit = await getHabit(habitId);
                setHabit(fetchedHabit);
            } catch (error) {
                console.error('Failed to fetch habit:', error);
                toast({
                    title: 'Error',
                    description: 'Failed to load habit details',
                    status: 'error',
                    duration: 5000,
                    isClosable: true,
                });
                navigate('/');
            } finally {
                setIsLoading(false);
            }
        };

        fetchHabit();
    }, [id, navigate, toast]);

    if (isLoading) {
        return (
            <Container maxW="container.xl" py={8}>
                <Flex justify="center" align="center" h="200px">
                    <Spinner size="xl" />
                </Flex>
            </Container>
        );
    }

    if (!habit) return null;

    return (
        <Container maxW="container.xl" py={8}>
            <Stack spacing={8}>
                <HStack spacing={4}>
                    <IconButton
                        aria-label="Go back"
                        icon={<ArrowBackIcon />}
                        onClick={() => navigate('/')}
                    />
                    <Heading size="lg">{habit.name}</Heading>
                </HStack>

                <Tabs variant="enclosed">
                    <TabList>
                        <Tab>Overview</Tab>
                        <Tab>Logs</Tab>
                        <Tab>Tags</Tab>
                        <Tab>Settings</Tab>
                    </TabList>

                    <TabPanels>
                        <TabPanel>
                            <HabitStats habit={habit} />
                        </TabPanel>
                        <TabPanel>
                            <HabitLogList habit={habit} selectedTags={selectedTags} />
                        </TabPanel>
                        <TabPanel>
                            <Box mb={6}>
                                <TagManagement habit={habit} />
                            </Box>
                            <Heading size="md" mb={4}>Filter Logs by Tags</Heading>
                            <Wrap spacing={2}>
                                {habit.tags?.map(tag => (
                                    <WrapItem key={tag.id}>
                                        <Tag
                                            size="md"
                                            borderRadius="full"
                                            variant={selectedTags.includes(tag.value) ? "solid" : "outline"}
                                            colorScheme="blue"
                                            cursor="pointer"
                                            onClick={() => {
                                                setSelectedTags(prev =>
                                                    prev.includes(tag.value)
                                                        ? prev.filter(t => t !== tag.value)
                                                        : [...prev, tag.value]
                                                );
                                            }}
                                        >
                                            <TagLabel>{tag.value}</TagLabel>
                                        </Tag>
                                    </WrapItem>
                                ))}
                            </Wrap>
                        </TabPanel>
                        <TabPanel>
                            <HabitSettings
                                habit={habit}
                                onUpdate={(updatedHabit) => setHabit(updatedHabit)}
                            />
                        </TabPanel>
                    </TabPanels>
                </Tabs>
            </Stack>
        </Container>
    );
};