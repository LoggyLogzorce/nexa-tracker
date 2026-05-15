import React from 'react';
import styles from './AuthLayout.module.css';

interface AuthLayoutProps {
    children: React.ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
    return (
        <div className={styles.wrapper}>
            <div className={styles.card}>
                <h1 className={styles.logo}>NexaFlow</h1>
                {children}
            </div>
        </div>
    );
}