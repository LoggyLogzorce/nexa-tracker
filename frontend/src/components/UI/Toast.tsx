import type { NotificationType } from '../../contexts/NotificationContext';
import styles from './Toast.module.css';

interface Notification {
    id: string;
    type: NotificationType;
    message: string;
    leaving: boolean;
}

interface Props {
    notifications: Notification[];
    onRemove: (id: string) => void;
}

export default function Toast({ notifications, onRemove }: Props) {
    if (notifications.length === 0) return null;

    return (
        <div className={styles.container}>
            {notifications.map(n => (
                <div key={n.id} className={`${styles.toast} ${styles[n.type]} ${n.leaving ? styles.leaving : ''}`} onClick={() => onRemove(n.id)}>
                    <svg className={styles.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        {n.type === 'success' ? (
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        ) : (
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        )}
                    </svg>
                    <span className={styles.message}>{n.message}</span>
                </div>
            ))}
        </div>
    );
}
