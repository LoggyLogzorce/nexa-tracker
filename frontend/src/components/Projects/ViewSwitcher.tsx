import styles from './ViewSwitcher.module.css';

type View = 'kanban' | 'list';

interface Props { active: View; onChange: (view: View) => void; }

export default function ViewSwitcher({ active, onChange }: Props) {
    return (
        <div className={styles.switcher}>
            <button className={`${styles.btn} ${active === 'kanban' ? styles.active : ''}`} onClick={() => onChange('kanban')}>
                <svg className={styles.btnIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2"/></svg>
                Канбан
            </button>
            <button className={`${styles.btn} ${active === 'list' ? styles.active : ''}`} onClick={() => onChange('list')}>
                <svg className={styles.btnIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16"/></svg>
                Список
            </button>
        </div>
    );
}
