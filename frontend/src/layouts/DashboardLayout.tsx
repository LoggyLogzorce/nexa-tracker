import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from '../components/Dashboard/Sidebar';
import Header from '../components/Dashboard/Header';
import NewVersionModal from '../components/Changelog/NewVersionModal';
import { useVersion } from '../contexts/useVersion';
import styles from './DashboardLayout.module.css';

export default function DashboardLayout() {
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
    const { hasNewVersion, currentVersion, markSeen } = useVersion();

    return (
        <div className={styles.wrapper}>
            {hasNewVersion && <NewVersionModal version={currentVersion} onClose={markSeen} />}
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} collapsed={sidebarCollapsed} onToggleCollapse={() => setSidebarCollapsed(prev => !prev)} />
            <div className={styles.main}>
                <Header onToggleSidebar={() => setSidebarOpen(prev => !prev)} />
                <main className={styles.content}><Outlet /></main>
            </div>
        </div>
    );
}