import { useState, useEffect } from 'react';
import {
    Box,
    Button,
    HStack,
    Input,
    Tag,
    TagLabel,
    TagCloseButton,
    Text,
    VStack,
    useToast,
    Wrap,
    WrapItem,
} from '@chakra-ui/react';
import { AddIcon } from '@chakra-ui/icons';
import type { Habit, Tag as TagType } from '../types';
import { getTags, createTag, deleteTag } from '../services/api';

interface TagManagementProps {
    habit: Habit;
}

export const TagManagement = ({ habit }: TagManagementProps) => {
    const [tags, setTags] = useState<TagType[]>([]);
    const [newTagValue, setNewTagValue] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const toast = useToast();

    const fetchTags = async () => {
        try {
            const fetchedTags = await getTags(habit.id);
            setTags(fetchedTags);
        } catch (error) {
            console.error('Failed to fetch tags:', error);
            toast({
                title: 'Error',
                description: 'Failed to load tags',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchTags();
    }, [habit.id]);

    const handleAddTag = async () => {
        if (!newTagValue.trim()) return;

        try {
            await createTag(habit.id, newTagValue.trim());
            await fetchTags();
            setNewTagValue('');
            toast({
                title: 'Success',
                description: 'Tag added successfully',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
        } catch (error) {
            console.error('Failed to add tag:', error);
            toast({
                title: 'Error',
                description: 'Failed to add tag',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    const handleDeleteTag = async (tagId: number) => {
        try {
            await deleteTag(habit.id, tagId);
            await fetchTags();
            toast({
                title: 'Success',
                description: 'Tag removed successfully',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });
        } catch (error) {
            console.error('Failed to remove tag:', error);
            toast({
                title: 'Error',
                description: 'Failed to remove tag',
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    const handleKeyPress = (event: React.KeyboardEvent) => {
        if (event.key === 'Enter') {
            handleAddTag();
        }
    };

    return (
        <Box>
            <VStack spacing={6} align="stretch">
                <HStack>
                    <Input
                        placeholder="Add a new tag..."
                        value={newTagValue}
                        onChange={(e) => setNewTagValue(e.target.value)}
                        onKeyPress={handleKeyPress}
                    />
                    <Button
                        leftIcon={<AddIcon />}
                        colorScheme="blue"
                        onClick={handleAddTag}
                        isDisabled={!newTagValue.trim()}
                    >
                        Add
                    </Button>
                </HStack>

                <Box>
                    <Text mb={4} fontWeight="medium">Current Tags</Text>
                    {tags.length === 0 ? (
                        <Text color="gray.500">No tags added yet</Text>
                    ) : (
                        <Wrap spacing={2}>
                            {tags.map(tag => (
                                <WrapItem key={tag.id}>
                                    <Tag
                                        size="md"
                                        borderRadius="full"
                                        variant="solid"
                                        colorScheme="blue"
                                    >
                                        <TagLabel>{tag.value}</TagLabel>
                                        <TagCloseButton
                                            onClick={() => handleDeleteTag(tag.id)}
                                        />
                                    </Tag>
                                </WrapItem>
                            ))}
                        </Wrap>
                    )}
                </Box>
            </VStack>
        </Box>
    );
};