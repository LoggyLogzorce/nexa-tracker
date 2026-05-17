import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios';
import type { ApiResponse } from '../types/auth';

// export const API_ORIGIN = 'http://10.147.18.142:8080';
export const API_ORIGIN = 'http://192.168.3.17:8080';

const client = axios.create({
    baseURL: `${API_ORIGIN}/api/v1`,
    withCredentials: true,
    headers: { 'Content-Type': 'application/json' },
});

let currentAccessToken: string | null = null;

export function setCurrentAccessToken(token: string | null) {
    currentAccessToken = token;
}

let onLogout: () => void = () => {};

export function setOnLogout(handler: () => void) {
    onLogout = handler;
}

client.interceptors.request.use((config: InternalAxiosRequestConfig) => {
    if (currentAccessToken) {
        config.headers.Authorization = `Bearer ${currentAccessToken}`;
    }
    return config;
});

let isRefreshing = false;
let failedQueue: Array<{
    resolve: (value: unknown) => void;
    reject: (reason: unknown) => void;
}> = [];

function processQueue(error: unknown) {
    failedQueue.forEach(prom => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(undefined);
        }
    });
    failedQueue = [];
}

client.interceptors.response.use(
    response => response,
    async (error: AxiosError<ApiResponse>) => {
        const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

        if (error.response?.status === 401 && !originalRequest._retry && !originalRequest.url?.includes('/auth/refresh')) {
            if (isRefreshing) {
                return new Promise((resolve, reject) => {
                    failedQueue.push({ resolve, reject });
                }).then(() => {
                    originalRequest._retry = true;
                    if (currentAccessToken) {
                        originalRequest.headers.Authorization = `Bearer ${currentAccessToken}`;
                    }
                    return client(originalRequest);
                });
            }

            originalRequest._retry = true;
            isRefreshing = true;

            try {
                const { data } = await client.post('/auth/refresh');
                const newToken = (data as ApiResponse<{ access_token: string }>).data?.access_token;
                if (newToken) {
                    setCurrentAccessToken(newToken);
                }
                processQueue(null);
                if (newToken) {
                    originalRequest.headers.Authorization = `Bearer ${newToken}`;
                }
                return client(originalRequest);
            } catch (refreshError) {
                processQueue(refreshError);
                onLogout();
                return Promise.reject(refreshError);
            } finally {
                isRefreshing = false;
            }
        }

        return Promise.reject(error);
    },
);

export function extractData<T>(apiResponse: ApiResponse<T>): T {
    if (!apiResponse.success) {
        throw new Error(apiResponse.error || 'Unknown error');
    }
    return apiResponse.data as T;
}

export default client;
