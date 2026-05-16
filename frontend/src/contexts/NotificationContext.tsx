import { createContext, useState, useCallback, type ReactNode } from 'react';
import Toast from '../components/UI/Toast';

export type NotificationType = 'success' | 'error';

interface Notification {
    id: string;
    type: NotificationType;
    message: string;
    leaving: boolean;
}

export interface NotificationContextType {
    addNotification: (type: NotificationType, message: string) => void;
}

// eslint-disable-next-line react-refresh/only-export-components
export const NotificationContext = createContext<NotificationContextType | null>(null);

let nextId = 1;

export function NotificationProvider({ children }: { children: ReactNode }) {
    const [notifications, setNotifications] = useState<Notification[]>([]);

    const removeNotification = useCallback((id: string) => {
        setNotifications(prev => prev.filter(n => n.id !== id));
    }, []);

    const addNotification = useCallback((type: NotificationType, message: string) => {
        const id = String(nextId++);
        setNotifications(prev => [...prev, { id, type, message, leaving: false }]);
        setTimeout(() => setNotifications(prev => prev.map(n => n.id === id ? { ...n, leaving: true } : n)), 2500);
        setTimeout(() => removeNotification(id), 3000);
    }, [removeNotification]);

    return (
        <NotificationContext.Provider value={{ addNotification }}>
            {children}
            <Toast notifications={notifications} onRemove={removeNotification} />
        </NotificationContext.Provider>
    );
}
