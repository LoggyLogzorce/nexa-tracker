import { NavLink } from 'react-router-dom';
import { useAuth } from '../../contexts/useAuth';
import { useVersion } from '../../contexts/useVersion';
import styles from './Sidebar.module.css';

const navItems = [
    { to: '/dashboard', label: 'Дашборд', icon: 'M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z' },
    { to: '/tasks', label: 'Задачи', icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01' },
    { to: '/projects', label: 'Проекты', icon: 'M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10' },
    { to: '/changelog', label: 'Что нового', icon: 'M11.48 3.499a.562.562 0 0 1 1.04 0l2.125 5.111a.563.563 0 0 0 .475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 0 0-.182.557l1.285 5.385a.562.562 0 0 1-.84.61l-4.725-2.885a.562.562 0 0 0-.586 0L6.982 20.54a.562.562 0 0 1-.84-.61l1.285-5.386a.562.562 0 0 0-.182-.557l-4.204-3.602a.563.563 0 0 1 .321-.988l5.518-.442a.563.563 0 0 0 .475-.345L11.48 3.5Z' },
];

interface Props { isOpen?: boolean; onClose?: () => void; collapsed?: boolean; onToggleCollapse?: () => void; }

export default function Sidebar({ isOpen, onClose, collapsed, onToggleCollapse }: Props) {
    const { logout } = useAuth();
    const { hasNewVersion } = useVersion();
    return (
        <>
            {isOpen && <div className={styles.overlay} onClick={onClose} />}
            <aside className={`${styles.sidebar} ${isOpen ? styles.open : ''} ${collapsed ? styles.collapsed : ''}`}>
            <div>
                <div className={styles.logo}>
                    <img src="/logo.png" alt="NexaFlow" className={styles.logoImg} />
                    <span>NexaFlow</span>
                </div>
                <nav className={styles.nav}>
                    {navItems.map(item => (
                        <NavLink key={item.to} to={item.to} className={({ isActive }) =>
                            `${styles.navItem} ${isActive ? styles.active : ''}`}>
                            <span className={styles.iconWrap}>
                                <svg className={styles.navIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={item.icon} />
                                </svg>
                                {item.to === '/changelog' && hasNewVersion && <span className={styles.badge} />}
                            </span>
                            <span>{item.label}</span>
                        </NavLink>
                    ))}
                </nav>
            </div>
            <div className={styles.bottom}>
                <button className={`${styles.navItem} ${styles.logout}`} onClick={logout}>
                    <svg className={styles.navIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
                    </svg>
                    <span>Выход</span>
                </button>
                <button className={styles.collapseBtn} onClick={onToggleCollapse}>
                    <svg className={styles.navIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={collapsed ? 'M9 5l7 7-7 7' : 'M15 19l-7-7 7-7'} />
                    </svg>
                </button>
            </div>
        </aside>
        </>
    );
}