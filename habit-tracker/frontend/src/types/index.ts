export interface User {
    id: number;
    profileImageUrl: string | null;
    name: string;
    timeZone: string;
    phone: string | null;
    email: string;
    createdAt: string;
}

export interface Habit {
    id: number;
    userId: number;
    name: string;
    description: string;
    frequency: number;
    createdAt: string;
    logs?: Log[];
    tags?: Tag[];
}

export interface HabitWithStreak extends Habit {
    currentStreak: number;
    longestStreak: number;
    lastLogDate: string | null;
}

export interface LogFormData {
    notes: string;
    completedAt: string;
}

export interface Log {
    id: number;
    habitId: number;
    notes: string;
    createdAt: string;
    completedAt: string;
    tags?: Tag[];
}

export interface Tag {
    id: number;
    habitId: number;
    value: string;
    createdAt: string;
}

export interface AuthResponse {
    token: string;
    user: User;
}