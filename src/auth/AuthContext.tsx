import React, { createContext, useCallback, useEffect, useMemo, useState } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { ApiClient, LoginResponse } from '../api/client';

type AuthState = {
	token: string | null;
	user: LoginResponse['user'] | null;
};

type AuthContextValue = {
	state: AuthState;
	client: ApiClient;
	login: (phone: string, password: string) => Promise<void>;
	logout: () => Promise<void>;
	register: (input: {
		name: string;
		phone_number: string;
		password: string;
		time_zone: string;
		profile_image_url?: string;
	}) => Promise<void>;
};

const STORAGE_KEY = 'auth_state_v1';

export const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	const [state, setState] = useState<AuthState>({ token: null, user: null });

	const client = useMemo(
		() => new ApiClient({ getToken: () => state.token }),
		[state.token]
	);

	useEffect(() => {
		(async () => {
			try {
				const saved = await AsyncStorage.getItem(STORAGE_KEY);
				if (saved) {
					setState(JSON.parse(saved));
				}
			} catch {}
		})();
	}, []);

	useEffect(() => {
		(async () => {
			try {
				await AsyncStorage.setItem(STORAGE_KEY, JSON.stringify(state));
			} catch {}
		})();
	}, [state]);

	const login = useCallback(async (phone: string, password: string) => {
		const res = await client.login({ phone_number: phone, password });
		setState({ token: res.token, user: res.user });
	}, [client]);

	const register = useCallback(
		async (input: { name: string; phone_number: string; password: string; time_zone: string; profile_image_url?: string }) => {
			await client.register(input);
			// Auto-login after register
			const res = await client.login({ phone_number: input.phone_number, password: input.password });
			setState({ token: res.token, user: res.user });
		},
		[client]
	);

	const logout = useCallback(async () => {
		setState({ token: null, user: null });
		await AsyncStorage.removeItem(STORAGE_KEY);
	}, []);

	const value: AuthContextValue = useMemo(
		() => ({ state, client, login, logout, register }),
		[state, client, login, logout, register]
	);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export function useAuth(): AuthContextValue {
	const ctx = React.useContext(AuthContext);
	if (!ctx) throw new Error('useAuth must be used within AuthProvider');
	return ctx;
}


