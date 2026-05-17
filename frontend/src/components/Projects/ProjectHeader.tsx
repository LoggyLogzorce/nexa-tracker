import { useState } from 'react';
import type { Project, ProjectStatus, Priority } from '../../types/project';
import styles from './ProjectHeader.module.css';

const statusConfig: Record<ProjectStatus, string> = {
    'В работе': styles.badgeBlue,
    'Планирование': styles.badgeAmber,
    'Завершен': styles.badgeGreen,
};

const priorityConfig: Record<Priority, string> = {
    'Высокий': '#ef4444',
    'Средний': '#f59e0b',
    'Низкий': '#22c55e',
};

interface Props { project: Project; onEdit?: () => void; onDelete?: () => void; onEditStatuses?: () => void; onEditPriorities?: () => void; }

export default function ProjectHeader({ project, onEdit, onDelete, onEditStatuses, onEditPriorities }: Props) {
    const [menuOpen, setMenuOpen] = useState(false);
    const createdDate = new Date(project.createdAt).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' });

    return (
        <div className={styles.header}>
            <div className={styles.left}>
                <div className={styles.titleRow}>
                    <h1 className={styles.title}>{project.title}</h1>
                    <span className={`${styles.statusBadge} ${statusConfig[project.status]}`}>{project.status}</span>
                    <span className={styles.priorityBadge} style={{ color: priorityConfig[project.priority] }}>
                        <svg className={styles.priorityDot} width="8" height="8" viewBox="0 0 8 8"><circle cx="4" cy="4" r="3" fill="currentColor"/></svg>
                        {project.priority}
                    </span>
                </div>
                <p className={styles.description}>{project.description}</p>
                <div className={styles.meta}>
                    <span className={styles.metaItem}>
                        <svg className={styles.metaIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
                        Создан: {createdDate}
                    </span>
                    <span className={styles.metaItem}>
                        <svg className={styles.metaIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>
                        Автор: {project.owner.name || project.owner.email}
                    </span>
                </div>
            </div>
            {project.role === 'owner' && (
                <div className={styles.menuWrap}>
                    <button className={styles.menuBtn} onClick={() => setMenuOpen(!menuOpen)}>
                        <svg className={styles.actionIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"/></svg>
                    </button>
                    {menuOpen && (
                        <div className={styles.dropdown}>
                            <button className={styles.dropdownItem} onClick={() => { onEdit?.(); setMenuOpen(false); }}>
                                <svg className={styles.dropdownIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                                Редактировать
                            </button>
                            <button className={`${styles.dropdownItem} ${styles.dropdownDelete}`} onClick={() => { onDelete?.(); setMenuOpen(false); }}>
                                <svg className={styles.dropdownIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                                Удалить
                            </button>
                            <div className={styles.dropdownDivider} />
                            <button className={styles.dropdownItem} onClick={() => { onEditStatuses?.(); setMenuOpen(false); }}>
                                <svg className={styles.dropdownIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>
                                Редактировать статусы
                            </button>
                            <button className={styles.dropdownItem} onClick={() => { onEditPriorities?.(); setMenuOpen(false); }}>
                                <svg className={styles.dropdownIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4h4v4H3V4zm7 0h4v4h-4V4zm7 0h4v4h-4V4zM3 10h4v4H3v-4zm7 0h4v4h-4v-4zm7 0h4v4h-4v-4zM3 16h4v4H3v-4zm7 0h4v4h-4v-4zm7 0h4v4h-4v-4z"/></svg>
                                Редактировать приоритеты
                            </button>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}
