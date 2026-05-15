import { useState } from 'react';
import type {Project, ProjectStatus, Priority} from '../../types/project';
import styles from './ProjectCard.module.css';

const statusConfig: Record<ProjectStatus, { badge: string; bar: string }> = {
    'В работе': { badge: styles.badgeBlue, bar: styles.barBlue },
    'Планирование': { badge: styles.badgeAmber, bar: styles.barAmber },
    'Завершен': { badge: styles.badgeGreen, bar: styles.barGreen },
};
const priorityConfig: Record<Priority, { color: string; icon: string }> = {
    'Высокий': { color: styles.priorityHigh, icon: 'M12 2L2 22h20L12 2zm0 4l6.5 13h-13L12 6z' },
    'Средний': { color: styles.priorityMedium, icon: 'M12 2L2 22h20L12 2zm0 4l6.5 13h-13L12 6z' },
    'Низкий': { color: styles.priorityLow, icon: 'M12 2L2 22h20L12 2zm0 4l6.5 13h-13L12 6z' },
};

interface Props { project: Project; onEdit: (p: Project) => void; onDelete: (p: Project) => void; }

export default function ProjectCard({ project, onEdit, onDelete }: Props) {
    const [menuOpen, setMenuOpen] = useState(false);
    const s = statusConfig[project.status];
    const p = priorityConfig[project.priority];

    return (
        <div className={styles.card}>
            <div className={styles.content}>
                <div className={styles.header}>
                    <span className={`${styles.badge} ${s.badge}`}>{project.status}</span>
                    <div className={styles.menu}>
                        {project.role === 'owner' && (
                            <>
                                <button className={styles.menuBtn} onClick={(e) => { e.stopPropagation(); setMenuOpen(!menuOpen); }}>
                                    <svg className={styles.dots} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"/></svg>
                                </button>
                                {menuOpen && (
                                    <div className={styles.dropdown}>
                                        <button className={styles.dropdownItem} onClick={(e) => { e.stopPropagation(); onEdit(project); setMenuOpen(false); }}>
                                            <svg className={styles.smallIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                                            Редактировать
                                        </button>
                                        <button className={`${styles.dropdownItem} ${styles.delete}`} onClick={(e) => { e.stopPropagation(); onDelete(project); setMenuOpen(false); }}>
                                            <svg className={styles.smallIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                                            Удалить
                                        </button>
                                    </div>
                                )}
                            </>
                        )}
                    </div>
                </div>
                <h3 className={styles.projectTitle}>{project.title}</h3>
                <p className={styles.description}>{project.description}</p>
                <div className={styles.footer}>
                    <div className={styles.priority}>
                        <svg className={`${styles.priorityIcon} ${p.color}`} fill="currentColor" viewBox="0 0 24 24"><path d={p.icon}/></svg>
                        <span className={styles.priorityText}>{project.priority}</span>
                    </div>
                    <span className={styles.date}>Создан: {new Date(project.createdAt).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })}</span>
                </div>
            </div>
            <div className={`${styles.statusBar} ${s.bar}`}></div>
        </div>
    );
}