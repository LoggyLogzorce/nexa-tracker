import { API_ORIGIN } from '../../api/client';
import styles from './Avatar.module.css';

const avatarColors = ['#0ea5e9', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#ef4444', '#14b8a6', '#f97316', '#6366f1', '#84cc16'];

function getAvatarColor(name: string): string {
    let hash = 0;
    for (let i = 0; i < name.length; i++) {
        hash = name.charCodeAt(i) + ((hash << 5) - hash);
    }
    return avatarColors[Math.abs(hash) % avatarColors.length];
}

interface Props {
    name: string;
    avatarUrl?: string | null;
    size?: number;
}

export default function Avatar({ name, avatarUrl, size = 32 }: Props) {
    const initials = name ? name.charAt(0).toUpperCase() : '?';

    const avatarSrc = avatarUrl
        ? avatarUrl.startsWith('http')
            ? avatarUrl
            : `${API_ORIGIN}/${avatarUrl.replace(/^\//, '')}`
        : null;

    return (
        <div className={styles.wrap} style={{ width: size, height: size, ...(!avatarSrc ? { background: getAvatarColor(name) } : {}) }}>
            {avatarSrc ? (
                <img className={styles.img} src={avatarSrc} alt={name} />
            ) : (
                <span className={styles.text} style={{ fontSize: Math.round(size * 0.44) }}>{initials}</span>
            )}
        </div>
    );
}
