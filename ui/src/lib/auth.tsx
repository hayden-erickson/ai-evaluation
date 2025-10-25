import React, { createContext, useContext, useEffect, useMemo, useState } from "react";
import { api, LoginRequest, LoginResponse } from "./api";

type AuthContextType = {
	user: LoginResponse["user"] | null;
	token: string | null;
	login: (req: LoginRequest) => Promise<void>;
	logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
	const [user, setUser] = useState<LoginResponse["user"] | null>(null);
	const [token, setToken] = useState<string | null>(localStorage.getItem("auth_token"));

	useEffect(() => {
		if (token) api.setAuthToken(token);
	}, [token]);

	const login = async (req: LoginRequest) => {
		const res = await api.login(req);
		setUser(res.user);
		setToken(res.token);
		api.setAuthToken(res.token);
	};

	const logout = () => {
		setUser(null);
		setToken(null);
		api.setAuthToken(null);
	};

	const value = useMemo(() => ({ user, token, login, logout }), [user, token]);
	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
	const ctx = useContext(AuthContext);
	if (!ctx) throw new Error("useAuth must be used within AuthProvider");
	return ctx;
}


