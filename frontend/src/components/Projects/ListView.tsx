import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import type { Task } from '../../types/task';
import type { Project } from '../../types/project';
import Avatar from '../UI/Avatar';
import styles from './ListView.module.css';

interface Props { tasks: Task[]; project: Project; onEdit?: (task: Task) => void; onArchive?: (task: Task) => void; onDelete?: (task: Task) => void; }

export default function ListView({ tasks, project, onEdit, onArchive, onDelete }: Props) {
    const navigate = useNavigate();
    const [menuTaskId, setMenuTaskId] = useState<number | null>(null);

    const formatDate = (dateStr: string) => {
        const d = new Date(dateStr);
        return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
    };

    const handleRowClick = (task: Task) => {
        navigate(`/projects/${project.id}/tasks/${task.id}${task.is_archive ? '?archived=true' : ''}`, { state: { project } });
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
                        <th className={styles.thThin}></th>
                    </tr>
                </thead>
                <tbody className={styles.tbody}>
                    {tasks.map(task => {
                        const isOverdue = task.deadline && new Date(task.deadline) < new Date();
                        return (
                            <tr key={task.id} className={styles.tr} onClick={() => handleRowClick(task)}>
                                <td className={styles.td}>
                                    <div className={styles.titleRow}>
                                        <span className={styles.taskId}>#{task.id}</span>
                                        <div>
                                            <p className={styles.taskTitle}>{task.title}</p>
                                            {task.description && <p className={styles.taskDesc}>{task.description.substring(0, 50)}{task.description.length > 50 ? '...' : ''}</p>}
                                        </div>
                                    </div>
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
                                            <Avatar name={task.assignee.name} avatarUrl={task.assignee.avatar_url} size={24} />
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
                                <td className={styles.tdActions}>
                                    <div className={styles.dotsWrap} onClick={e => e.stopPropagation()}>
                                        <button className={styles.dotsBtn} onClick={() => setMenuTaskId(menuTaskId === task.id ? null : task.id)}>&#8942;</button>
                                        {menuTaskId === task.id && (
                                            <div className={styles.dropdown}>
                                                <button className={styles.dropdownItem} onClick={() => { onEdit?.(task); setMenuTaskId(null); }}>
                                                    <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                                                    Редактировать
                                                </button>
                                                <button className={styles.dropdownItem} onClick={() => { onArchive?.(task); setMenuTaskId(null); }}>
                                                    <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"/></svg>
                                                    {task.is_archive ? 'Восстановить' : 'В архив'}
                                                </button>
                                                <button className={`${styles.dropdownItem} ${styles.dropdownItemDanger}`} onClick={() => { onDelete?.(task); setMenuTaskId(null); }}>
                                                    <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                                                    Удалить
                                                </button>
                                            </div>
                                        )}
                                    </div>
                                </td>
                            </tr>
                        );
                    })}
                </tbody>
            </table>
        </div>
    );
}
