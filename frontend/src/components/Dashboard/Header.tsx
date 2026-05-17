import { useState, useEffect, useRef, useCallback } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/useAuth';
import { searchTasksApi, searchProjectsApi } from '../../api/projects';
import type { Task } from '../../types/task';
import type { Project } from '../../types/project';
import Avatar from '../UI/Avatar';
import styles from './Header.module.css';

interface Props { onToggleSidebar?: () => void; }

export default function Header({ onToggleSidebar }: Props) {
    const { user } = useAuth();
    const navigate = useNavigate();
    const [notifOpen, setNotifOpen] = useState(false);
    const [query, setQuery] = useState('');
    const [tasks, setTasks] = useState<Task[]>([]);
    const [projects, setProjects] = useState<Project[]>([]);
    const [open, setOpen] = useState(false);
    const [loading, setLoading] = useState(false);
    const searchRef = useRef<HTMLDivElement>(null);
    const inputRef = useRef<HTMLInputElement>(null);

    const doSearch = useCallback(async (q: string) => {
        if (q.length < 2) {
            setTasks([]);
            setProjects([]);
            setOpen(false);
            return;
        }
        setLoading(true);
        try {
            const [t, p] = await Promise.all([
                searchTasksApi(q),
                searchProjectsApi(q),
            ]);
            setTasks(t);
            setProjects(p);
            setOpen(t.length > 0 || p.length > 0);
        } catch {
            setTasks([]);
            setProjects([]);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        const timer = setTimeout(() => doSearch(query), 300);
        return () => clearTimeout(timer);
    }, [query, doSearch]);

    useEffect(() => {
        const handleClick = (e: MouseEvent) => {
            if (searchRef.current && !searchRef.current.contains(e.target as Node)) {
                setOpen(false);
            }
        };
        const handleEsc = (e: KeyboardEvent) => {
            if (e.key === 'Escape') {
                setOpen(false);
                inputRef.current?.blur();
            }
        };
        document.addEventListener('mousedown', handleClick);
        document.addEventListener('keydown', handleEsc);
        return () => {
            document.removeEventListener('mousedown', handleClick);
            document.removeEventListener('keydown', handleEsc);
        };
    }, []);

    const handleSelect = () => {
        setQuery('');
        setTasks([]);
        setProjects([]);
        setOpen(false);
    };

    return (
        <header className={styles.header}>
            <button className={styles.mobileBtn} onClick={onToggleSidebar}>
                <svg className={styles.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16"/></svg>
            </button>
            <div className={styles.search} ref={searchRef}>
                <svg className={styles.searchIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
                <input
                    ref={inputRef}
                    type="text"
                    placeholder="Поиск задач, проектов..."
                    value={query}
                    onChange={e => setQuery(e.target.value)}
                    onFocus={() => { if (tasks.length > 0 || projects.length > 0) setOpen(true); }}
                />
                {loading && <span className={styles.searchSpinner} />}
                {open && (
                    <div className={styles.dropdown}>
                        {projects.length > 0 && (
                            <div className={styles.section}>
                                <div className={styles.sectionTitle}>Проекты</div>
                                {projects.map(p => (
                                    <button key={p.id} className={styles.resultItem} onClick={() => { navigate(`/projects/${p.id}`); handleSelect(); }}>
                                        <svg className={styles.resultIcon} width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
                                        <span>{p.title}</span>
                                    </button>
                                ))}
                            </div>
                        )}
                        {tasks.length > 0 && (
                            <div className={styles.section}>
                                <div className={styles.sectionTitle}>Задачи</div>
                                {tasks.map(t => (
                                    <button key={t.id} className={styles.resultItem} onClick={() => { navigate(`/projects/${t.project_id}/tasks/${t.id}`); handleSelect(); }}>
                                        <span className={styles.taskId}>#{t.id}</span>
                                        <div className={styles.taskInfo}>
                                            <span className={styles.taskTitle}>{t.title}</span>
                                            {t.project_title && <span className={styles.taskProject}>{t.project_title}</span>}
                                        </div>
                                    </button>
                                ))}
                            </div>
                        )}
                    </div>
                )}
            </div>
            <div className={styles.actions}>
                <div className={styles.notifWrap}>
                    <button className={styles.notifBtn} onClick={() => setNotifOpen(!notifOpen)}>
                        <span className={styles.badge}></span>
                        <svg className={styles.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/></svg>
                    </button>
                    {notifOpen && (
                        <div className={styles.notifPanel}>
                            <div className={styles.notifHeader}><h3>Уведомления</h3></div>
                            <div className={styles.notifList}>
                                <div className={styles.notifItem}>
                                    <p className={styles.notifTitle}>Новая задача назначена</p>
                                    <p className={styles.notifDesc}>В проекте "Редизайн мобильного приложения"</p>
                                    <span className={styles.notifTime}>5 минут назад</span>
                                </div>
                            </div>
                            <div className={styles.notifFooter}><button>Показать все</button></div>
                        </div>
                    )}
                </div>
                    <Link to="/profile" className={styles.profile}>
                        <div className={styles.userInfo}>
                            <p className={styles.userName}>{user?.name || 'Пользователь'}</p>
                            {user?.role && <p className={styles.userRole}>{user.role}</p>}
                        </div>
                        <Avatar name={user?.name || '?'} avatarUrl={user?.avatar_url} size={40} />
                    </Link>
            </div>
        </header>
    );
}
