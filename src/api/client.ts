import { Platform } from 'react-native';

export type LoginResponse = {
	token: string;
	user: {
		id: number;
		name: string;
		phone_number: string;
		time_zone: string;
		profile_image_url?: string;
		created_at: string;
	};
};

export type Habit = {
	id: number;
	user_id: number;
	name: string;
	description?: string;
	duration_seconds?: number | null;
	created_at: string;
};

export type Log = {
	id: number;
	habit_id: number;
	notes?: string;
	duration_seconds?: number | null;
	created_at: string;
};

export type CreateHabitRequest = {
	name: string;
	description?: string;
	duration_seconds?: number | null;
};

export type UpdateHabitRequest = Partial<CreateHabitRequest>;

export type CreateLogRequest = {
	notes?: string;
	duration_seconds?: number | null;
};

export type UpdateLogRequest = Partial<CreateLogRequest>;

export type ApiError = {
	status: number;
	message: string;
};

type FetchOptions = RequestInit & { token?: string };

function getDefaultBaseUrl(): string {
	// Use device localhost mapping for Android emulator
	const host = Platform.OS === 'android' ? '10.0.2.2' : 'localhost';
	return `http://${host}:8080`;
}

export class ApiClient {
	private readonly baseUrl: string;
	private readonly getToken: () => Promise<string | null> | string | null;

	constructor(params?: { baseUrl?: string; getToken?: () => Promise<string | null> | string | null }) {
		this.baseUrl = params?.baseUrl ?? getDefaultBaseUrl();
		this.getToken = params?.getToken ?? (() => null);
	}

	private async request<T>(path: string, options: FetchOptions = {}): Promise<T> {
		const token = options.token ?? (await Promise.resolve(this.getToken()));
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
			...(token ? { Authorization: `Bearer ${token}` } : {}),
			...(options.headers ?? {}),
		};

		const res = await fetch(`${this.baseUrl}${path}`, {
			...options,
			headers,
		});

		const text = await res.text();
		const maybeJson = text ? safeJsonParse(text) : null;

		if (!res.ok) {
			const message = typeof maybeJson === 'string' ? maybeJson : (maybeJson?.message ?? text ?? 'Request failed');
			throw <ApiError>{ status: res.status, message };
		}

		return (maybeJson as T) ?? ({} as T);
	}

	// Auth
	async register(payload: {
		name: string;
		phone_number: string;
		password: string;
		time_zone: string;
		profile_image_url?: string;
	}): Promise<{ id: number } & Omit<LoginResponse['user'], 'id'>> {
		return this.request(`/users/register`, { method: 'POST', body: JSON.stringify(payload) });
	}

	async login(payload: { phone_number: string; password: string }): Promise<LoginResponse> {
		return this.request(`/users/login`, { method: 'POST', body: JSON.stringify(payload) });
	}

	// Habits
	async getHabits(): Promise<Habit[]> {
		return this.request(`/habits`);
	}

	async createHabit(payload: CreateHabitRequest): Promise<Habit> {
		return this.request(`/habits`, { method: 'POST', body: JSON.stringify(payload) });
	}

	async updateHabit(habitId: number, payload: UpdateHabitRequest): Promise<Habit> {
		return this.request(`/habits/${habitId}`, { method: 'PUT', body: JSON.stringify(payload) });
	}

	async deleteHabit(habitId: number): Promise<void> {
		await this.request(`/habits/${habitId}`, { method: 'DELETE' });
	}

	// Logs
	async getHabitLogs(habitId: number): Promise<Log[]> {
		return this.request(`/habits/${habitId}/logs`);
	}

	async createLog(habitId: number, payload: CreateLogRequest): Promise<Log> {
		return this.request(`/habits/${habitId}/logs`, { method: 'POST', body: JSON.stringify(payload) });
	}

	async updateLog(logId: number, payload: UpdateLogRequest): Promise<Log> {
		return this.request(`/logs/${logId}`, { method: 'PUT', body: JSON.stringify(payload) });
	}

	async deleteLog(logId: number): Promise<void> {
		await this.request(`/logs/${logId}`, { method: 'DELETE' });
	}
}

function safeJsonParse(text: string): unknown {
	try {
		return JSON.parse(text);
	} catch {
		return text;
	}
}


