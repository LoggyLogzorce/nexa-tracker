import { createContext, useState, useEffect, useCallback, type ReactNode } from 'react';
import * as authApi from '../api/auth';
import { setCurrentAccessToken, setOnLogout } from '../api/client';
import type { UserResponse, LoginPayload, RegisterPayload } from '../types/auth';

export interface AuthContextType {
    user: UserResponse | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (payload: LoginPayload) => Promise<void>;
    register: (payload: RegisterPayload) => Promise<void>;
    logout: () => Promise<void>;
    loadUser: () => Promise<void>;
}

// eslint-disable-next-line react-refresh/only-export-components
export const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<UserResponse | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [sessionValid, setSessionValid] = useState(false);

    const logout = useCallback(async () => {
        try {
            await authApi.logout();
        } catch {
            // ignore
        }
        setCurrentAccessToken(null);
        setUser(null);
        setSessionValid(false);
    }, []);

    const loadUser = useCallback(async () => {
        const userData = await authApi.getMe();
        setUser(userData);
    }, []);

    useEffect(() => {
        setOnLogout(() => {
            setCurrentAccessToken(null);
            setUser(null);
            setSessionValid(false);
        });
    }, []);

    useEffect(() => {
        const init = async () => {
            try {
                const { access_token } = await authApi.refresh();
                setCurrentAccessToken(access_token);
                setSessionValid(true);
                await loadUser();
            } catch {
                setCurrentAccessToken(null);
            } finally {
                setIsLoading(false);
            }
        };
        init();
    }, [loadUser]);

    const login = async (payload: LoginPayload) => {
        const { access_token } = await authApi.login(payload);
        setCurrentAccessToken(access_token);
        setSessionValid(true);
        await loadUser();
    };

    const register = async (payload: RegisterPayload) => {
        await authApi.register(payload);
    };

    return (
        <AuthContext.Provider value={{
            user,
            isAuthenticated: sessionValid,
            isLoading,
            login,
            register,
            logout,
            loadUser,
        }}>
            {children}
        </AuthContext.Provider>
    );
}


