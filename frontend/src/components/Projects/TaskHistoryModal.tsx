import { useState, useEffect } from 'react';
import type { TaskHistoryEntry } from '../../api/projects';
import { getTaskHistoryApi } from '../../api/projects';
import Avatar from '../UI/Avatar';
import modalStyles from '../Dashboard/Modal.module.css';
import styles from './TaskHistoryModal.module.css';

interface Props {
    projectId: string;
    taskId: number;
    archived: boolean;
    onClose: () => void;
}

const fieldLabels: Record<string, string> = {
    title: 'Название',
    description: 'Описание',
    status: 'Статус',
    priority: 'Приоритет',
    assignee: 'Исполнитель',
    deadline: 'Дедлайн',
    is_archive: 'Архив',
};

const fmtDate = (iso: string) =>
    new Date(iso).toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });

export default function TaskHistoryModal({ projectId, taskId, archived, onClose }: Props) {
    const [entries, setEntries] = useState<TaskHistoryEntry[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const esc = (e: KeyboardEvent) => e.key === 'Escape' && onClose();
        document.addEventListener('keydown', esc);
        return () => document.removeEventListener('keydown', esc);
    }, [onClose]);

    useEffect(() => {
        let cancelled = false;
        const load = async () => {
            try {
                setLoading(true);
                setError(null);
                const data = await getTaskHistoryApi(projectId, taskId, archived);
                if (!cancelled) setEntries(data);
            } catch (err) {
                if (!cancelled) setError(err instanceof Error ? err.message : 'Ошибка загрузки');
            } finally {
                if (!cancelled) setLoading(false);
            }
        };
        load();
        return () => { cancelled = true; };
    }, [projectId, taskId, archived]);

    return (
        <div className={modalStyles.overlay} onClick={onClose}>
            <div
                className={`${modalStyles.modal} ${styles.modal}`}
                onClick={e => e.stopPropagation()}
            >
                <div className={modalStyles.header}>
                    <h2>История изменений</h2>
                </div>
                <div className={modalStyles.body}>
                    {loading && <div className={styles.loader}>Загрузка истории...</div>}
                    {error && <div className={styles.error}>{error}</div>}
                    {!loading && !error && entries.length === 0 && (
                        <div className={styles.empty}>История изменений пуста</div>
                    )}
                    {!loading && !error && entries.length > 0 && (
                        <div className={styles.entries}>
                            {entries.map(entry => (
                                <div key={entry.id} className={styles.entry}>
                                    <div className={styles.entryHeader}>
                                        <Avatar
                                            name={entry.user.name}
                                            avatarUrl={entry.user.avatar_url}
                                            size={24}
                                        />
                                        <span className={styles.userName}>{entry.user.name}</span>
                                        <span className={styles.entryDate}>
                                            {fmtDate(entry.created_at)}
                                        </span>
                                    </div>
                                    <div className={styles.changes}>
                                        {entry.changes.map((change, idx) => {
                                            const label = fieldLabels[change.field] || change.field;
                                            const oldVal = change.old_name || change.old_value;
                                            const newVal = change.new_name || change.new_value;
                                            return (
                                                <div key={idx} className={styles.changeItem}>
                                                    <span className={styles.changeField}>{label}:</span>
                                                    <span className={styles.changeOld}>{String(oldVal ?? '—')}</span>
                                                    <span className={styles.changeArrow}>→</span>
                                                    <span className={styles.changeNew}>{String(newVal ?? '—')}</span>
                                                </div>
                                            );
                                        })}
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
                <div className={modalStyles.footer}>
                    <button className={modalStyles.cancel} onClick={onClose}>Закрыть</button>
                </div>
            </div>
        </div>
    );
}
