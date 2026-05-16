import { useContext } from 'react';
import { NotificationContext } from './NotificationContext';
import type { NotificationContextType } from './NotificationContext';

export function useNotifications(): NotificationContextType {
    const ctx = useContext(NotificationContext);
    if (!ctx) throw new Error('useNotifications must be used within NotificationProvider');
    return ctx;
}
