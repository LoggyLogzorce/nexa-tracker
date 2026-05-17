import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../contexts/useAuth';
import Avatar from '../UI/Avatar';
import styles from './Header.module.css';

interface Props { onToggleSidebar?: () => void; }

export default function Header({ onToggleSidebar }: Props) {
    const { user } = useAuth();
    const [notifOpen, setNotifOpen] = useState(false);

    return (
        <header className={styles.header}>
            <button className={styles.mobileBtn} onClick={onToggleSidebar}>
                <svg className={styles.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16"/></svg>
            </button>
            <div className={styles.search}>
                <svg className={styles.searchIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
                <input type="text" placeholder="Поиск задач, проектов..." />
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