import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import type { Task } from '../../types/task';
import { getTasksByUserApi } from '../../api/projects';
import Avatar from '../../components/UI/Avatar';
import styles from './MyTasksPage.module.css';

type Tab = 'assigned' | 'reported';

export default function MyTasksPage() {
    const [tab, setTab] = useState<Tab>('assigned');
    const [tasks, setTasks] = useState<Task[]>([]);
    const [loading, setLoading] = useState(true);

    const handleTabChange = (newTab: Tab) => {
        setTab(newTab);
        setLoading(true);
    };

    useEffect(() => {
        getTasksByUserApi(tab)
            .then(setTasks)
            .finally(() => setLoading(false));
    }, [tab]);

    const fmtDate = (iso: string) => {
        const d = new Date(iso);
        return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
    };

    return (
        <div className={styles.page}>
            <div className={styles.header}>
                <h1 className={styles.title}>Мои задачи</h1>
            </div>

            <div className={styles.tabs}>
                <button className={`${styles.tab} ${tab === 'assigned' ? styles.active : ''}`} onClick={() => handleTabChange('assigned')}>
                    На меня назначены
                </button>
                <button className={`${styles.tab} ${tab === 'reported' ? styles.active : ''}`} onClick={() => handleTabChange('reported')}>
                    Созданные мной
                </button>
            </div>

            {loading ? (
                <p className={styles.loading}>Загрузка...</p>
            ) : tasks.length === 0 ? (
                <div className={styles.empty}>
                    <p>Нет задач</p>
                </div>
            ) : (
                <div className={styles.list}>
                    {tasks.map(task => (
                        <Link key={task.id} to={`/projects/${task.project_id}/tasks/${task.id}`} className={styles.card}>
                            <div className={styles.cardTop}>
                                <h3 className={styles.cardTitle}>{task.title}</h3>
                                <span className={styles.projectId}>#{task.id}</span>
                            </div>
                            {task.description && (
                                <p className={styles.cardDesc}>{task.description.substring(0, 100)}{task.description.length > 100 ? '...' : ''}</p>
                            )}
                            <div className={styles.cardMeta}>
                                <span className={styles.chip} style={{ backgroundColor: task.status.color + '18', color: task.status.color }}>
                                    <span className={styles.dot} style={{ background: task.status.color }} />
                                    {task.status.name}
                                </span>
                                <span className={styles.chip} style={{ backgroundColor: task.priority.color + '18', color: task.priority.color }}>
                                    <span className={styles.dot} style={{ background: task.priority.color }} />
                                    {task.priority.title}
                                </span>
                                {task.assignee && (
                                    <span className={styles.assignee}>
                                        <Avatar name={task.assignee.name} avatarUrl={task.assignee.avatar_url} size={18} />
                                        {task.assignee.name}
                                    </span>
                                )}
                                {task.deadline && (
                                    <span className={`${styles.deadline} ${new Date(task.deadline) < new Date() ? styles.overdue : ''}`}>
                                        {fmtDate(task.deadline)}
                                    </span>
                                )}
                            </div>
                        </Link>
                    ))}
                </div>
            )}
        </div>
    );
}