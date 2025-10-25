export type LoginRequest = { phone_number: string; password: string };
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
	duration_seconds?: number;
	created_at: string;
};

export type Log = {
	id: number;
	habit_id: number;
	notes?: string;
	duration_seconds?: number;
	created_at: string;
};

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8080";

function getAuthToken(): string | null {
	return localStorage.getItem("auth_token");
}

function setAuthToken(token: string | null) {
	if (token) localStorage.setItem("auth_token", token);
	else localStorage.removeItem("auth_token");
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
	const headers: HeadersInit = {
		"Content-Type": "application/json",
		...(options.headers || {}),
	};
	const token = getAuthToken();
	if (token) {
		(headers as Record<string, string>)["Authorization"] = `Bearer ${token}`;
	}

	const res = await fetch(`${API_BASE}${path}`, { ...options, headers });
	if (!res.ok) {
		let message = `Request failed with ${res.status}`;
		try {
			const text = await res.text();
			message = text || message;
		} catch {}
		throw new Error(message);
	}
	if (res.status === 204) return undefined as unknown as T;
	return (await res.json()) as T;
}

export const api = {
	setAuthToken,
	login: (body: LoginRequest) => request<LoginResponse>("/users/login", { method: "POST", body: JSON.stringify(body) }),
	register: (body: {
		name: string;
		time_zone: string;
		phone_number: string;
		password: string;
		profile_image_url?: string;
	}) => request("/users/register", { method: "POST", body: JSON.stringify(body) }),
	getHabits: () => request<Habit[]>("/habits"),
	createHabit: (body: { name: string; description?: string; duration_seconds?: number }) =>
		request<Habit>("/habits", { method: "POST", body: JSON.stringify(body) }),
	updateHabit: (id: number, body: { name?: string; description?: string; duration_seconds?: number }) =>
		request<Habit>(`/habits/${id}`, { method: "PUT", body: JSON.stringify(body) }),
	deleteHabit: (id: number) => request<void>(`/habits/${id}`, { method: "DELETE" }),
	getHabitLogs: (habitId: number) => request<Log[]>(`/habits/${habitId}/logs`),
	createLog: (habitId: number, body: { notes?: string; duration_seconds?: number }) =>
		request<Log>(`/habits/${habitId}/logs`, { method: "POST", body: JSON.stringify(body) }),
	getLog: (id: number) => request<Log>(`/logs/${id}`),
	updateLog: (id: number, body: { notes?: string; duration_seconds?: number }) =>
		request<Log>(`/logs/${id}`, { method: "PUT", body: JSON.stringify(body) }),
	deleteLog: (id: number) => request<void>(`/logs/${id}`, { method: "DELETE" }),
};


