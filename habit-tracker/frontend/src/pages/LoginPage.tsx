import { Box, Center, Heading, Stack, Text } from '@chakra-ui/react';
import { GoogleLogin } from '@react-oauth/google';
import { useAuth } from '../hooks/useAuth';
import { loginWithGoogle } from '../services/api';

export const LoginPage = () => {
    const { login } = useAuth();

    const handleGoogleSuccess = async (credentialResponse: { credential: string }) => {
        try {
            const response = await loginWithGoogle(credentialResponse.credential);
            login(response.token, response.user);
        } catch (error) {
            console.error('Login failed:', error);
        }
    };

    return (
        <Center minH="100vh" bg="gray.50">
            <Box bg="white" p={8} rounded="lg" shadow="lg" maxW="md" w="full">
                <Stack spacing={6} align="center">
                    <Heading size="xl">Habit Tracker</Heading>
                    <Text color="gray.600" textAlign="center">
                        Track your habits and build better routines
                    </Text>
                    <Box w="full">
                        <GoogleLogin
                            onSuccess={handleGoogleSuccess}
                            onError={() => {
                                console.error('Login failed');
                            }}
                        />
                    </Box>
                </Stack>
            </Box>
        </Center>
    );
};