import type { Task } from '../../types/task';
import styles from './ListView.module.css';

const avatarImages = [
    'https://image.qwenlm.ai/public_source/87cc46bb-251e-456e-a7c0-8d3b41f35211/1281ab55e-b3ab-4e9e-be91-85c89891a3c3.png',
];

interface Props { tasks: Task[]; }

export default function ListView({ tasks }: Props) {
    const formatDate = (dateStr: string) => {
        const d = new Date(dateStr);
        return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
    };

    return (
        <div className={styles.wrap}>
            <table className={styles.table}>
                <thead className={styles.thead}>
                    <tr>
                        <th className={styles.th}>Задача</th>
                        <th className={styles.th}>Статус</th>
                        <th className={styles.th}>Приоритет</th>
                        <th className={styles.th}>Исполнитель</th>
                        <th className={styles.th}>Дедлайн</th>
                    </tr>
                </thead>
                <tbody className={styles.tbody}>
                    {tasks.map(task => {
                        const isOverdue = task.deadline && new Date(task.deadline) < new Date();
                        return (
                            <tr key={task.id} className={styles.tr}>
                                <td className={styles.td}>
                                    <p className={styles.taskTitle}>{task.title}</p>
                                    {task.description && <p className={styles.taskDesc}>{task.description.substring(0, 50)}{task.description.length > 50 ? '...' : ''}</p>}
                                </td>
                                <td className={styles.td}>
                                    <span className={styles.statusBadge} style={{ backgroundColor: task.status.color + '20', color: task.status.color }}>
                                        {task.status.name}
                                    </span>
                                </td>
                                <td className={styles.td}>
                                    <span className={styles.priorityBadge} style={{ backgroundColor: task.priority.color + '20', color: task.priority.color }}>
                                        {task.priority.title}
                                    </span>
                                </td>
                                <td className={styles.td}>
                                    {task.assignee ? (
                                        <div className={styles.assignee}>
                                            <div className={styles.avatar}><img src={avatarImages[0]} alt="" /></div>
                                            <span className={styles.assigneeName}>{task.assignee.name}</span>
                                        </div>
                                    ) : (
                                        <span className={styles.noAssignee}>Не назначен</span>
                                    )}
                                </td>
                                <td className={styles.td}>
                                    <span className={`${styles.deadline} ${isOverdue ? styles.overdue : ''}`}>
                                        {task.deadline ? formatDate(task.deadline) : '-'}
                                    </span>
                                </td>
                            </tr>
                        );
                    })}
                </tbody>
            </table>
        </div>
    );
}
