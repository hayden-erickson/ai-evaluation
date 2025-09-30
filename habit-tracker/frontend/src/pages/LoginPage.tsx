import { Box, Center, Heading, Stack, Text } from '@chakra-ui/react';
import { GoogleLogin } from '@react-oauth/google';
import type { CredentialResponse } from '@react-oauth/google';
import { useAuth } from '../hooks/useAuth';
import { loginWithGoogle } from '../services/api';

export const LoginPage = () => {
    const { login } = useAuth();

    const handleGoogleSuccess = async (credentialResponse: CredentialResponse) => {
        try {
            if (!credentialResponse.credential) {
                throw new Error('No credential received from Google');
            }
            const response = await loginWithGoogle(credentialResponse.credential);
            login(response.token, response.user);
        } catch (error: any) {
            console.error('Login failed:', error);
            // Show error toast or message to user
            if (error?.response?.data?.error) {
                alert(error.response.data.error);
            } else {
                alert(error.message || 'Failed to sign in with Google. Please try again.');
            }
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