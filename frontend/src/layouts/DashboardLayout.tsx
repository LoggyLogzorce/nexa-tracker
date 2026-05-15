import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from '../components/Dashboard/Sidebar';
import Header from '../components/Dashboard/Header';
import styles from './DashboardLayout.module.css';

export default function DashboardLayout() {
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

    return (
        <div className={styles.wrapper}>
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} collapsed={sidebarCollapsed} onToggleCollapse={() => setSidebarCollapsed(prev => !prev)} />
            <div className={styles.main}>
                <Header onToggleSidebar={() => setSidebarOpen(prev => !prev)} />
                <main className={styles.content}><Outlet /></main>
            </div>
        </div>
    );
}