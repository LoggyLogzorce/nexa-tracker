import { useState } from 'react';
import { useDraggable } from '@dnd-kit/core';
import { useNavigate } from 'react-router-dom';
import type { Task } from '../../types/task';
import type { Project } from '../../types/project';
import Avatar from '../UI/Avatar';
import styles from './KanbanCard.module.css';

interface Props { task: Task; project: Project; onEdit?: (task: Task) => void; onArchive?: (task: Task) => void; onDelete?: (task: Task) => void; }

export default function KanbanCard({ task, project, onEdit, onArchive, onDelete }: Props) {
    const navigate = useNavigate();
    const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({ id: String(task.id) });
    const [menuOpen, setMenuOpen] = useState(false);

    const style = transform ? { transform: `translate(${transform.x}px, ${transform.y}px)` } : undefined;

    const formatDate = (dateStr: string) => {
        const d = new Date(dateStr);
        return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
    };

    const handleCardClick = () => {
        navigate(`/projects/${project.id}/tasks/${task.id}${task.is_archive ? '?archived=true' : ''}`, { state: { project } });
    };

    const handleDotsClick = (e: React.MouseEvent) => {
        e.stopPropagation();
        setMenuOpen(prev => !prev);
    };

    const handleEditClick = (e: React.MouseEvent) => {
        e.stopPropagation();
        setMenuOpen(false);
        onEdit?.(task);
    };

    const handleArchiveClick = (e: React.MouseEvent) => {
        e.stopPropagation();
        setMenuOpen(false);
        onArchive?.(task);
    };

    const handleDeleteClick = (e: React.MouseEvent) => {
        e.stopPropagation();
        setMenuOpen(false);
        onDelete?.(task);
    };

    return (
        <div ref={setNodeRef} {...listeners} {...attributes} className={`${styles.card} ${isDragging ? styles.dragging : ''}`} style={style} onClick={handleCardClick}>
            <div className={styles.body}>
                <div className={styles.topRow}>
                    <span className={styles.priorityBadge} style={{ backgroundColor: task.priority.color + '20', color: task.priority.color }}>
                        {task.priority.title}
                    </span>
                    <div className={styles.idWrap}>
                        <span className={styles.id}>#{task.id}</span>
                        <button className={styles.dotsBtn} onClick={handleDotsClick}>&#8942;</button>
                        {menuOpen && (
                            <div className={styles.dropdown} onClick={e => e.stopPropagation()}>
                                <button className={styles.dropdownItem} onClick={handleEditClick}>
                                    <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                                    Редактировать
                                </button>
                                <button className={styles.dropdownItem} onClick={handleArchiveClick}>
                                    <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"/></svg>
                                    {task.is_archive ? 'Восстановить' : 'В архив'}
                                </button>
                                <button className={`${styles.dropdownItem} ${styles.dropdownItemDanger}`} onClick={handleDeleteClick}>
                                    <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                                    Удалить
                                </button>
                            </div>
                        )}
                    </div>
                </div>
                <h4 className={styles.title}>{task.title}</h4>
                <div className={styles.bottomRow}>
                    {task.assignee ? (
                        <div className={styles.assignee}>
                            <Avatar name={task.assignee.name} avatarUrl={task.assignee.avatar_url} size={24} />
                        </div>
                    ) : (
                        <span className={styles.noAssignee}>Не назначен</span>
                    )}
                    {task.deadline && (
                        <span className={styles.deadline}>
                            <svg className={styles.deadlineIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
                            {formatDate(task.deadline)}
                        </span>
                    )}
                </div>
            </div>
        </div>
    );
}
