import client, { extractData } from './client';
import type { LoginPayload, RegisterPayload, LoginResponse, UserResponse, ApiResponse } from '../types/auth';

let refreshPromise: Promise<LoginResponse> | null = null;

export const login = async (payload: LoginPayload): Promise<LoginResponse> => {
    const response = await client.post<ApiResponse<LoginResponse>>('/auth/login', payload);
    return extractData(response.data);
};

export const register = async (payload: RegisterPayload): Promise<void> => {
    const response = await client.post<ApiResponse>('/auth/register', payload);
    extractData(response.data);
};

export const refresh = async (): Promise<LoginResponse> => {
    if (refreshPromise) return refreshPromise;

    refreshPromise = (async () => {
        const response = await client.post<ApiResponse<LoginResponse>>('/auth/refresh');
        return extractData(response.data);
    })();

    try {
        return await refreshPromise;
    } finally {
        refreshPromise = null;
    }
};

export const logout = async (): Promise<void> => {
    await client.post('/auth/logout');
};

export const getMe = async (): Promise<UserResponse> => {
    const response = await client.get<ApiResponse<UserResponse>>('/users/me');
    return extractData(response.data);
};
