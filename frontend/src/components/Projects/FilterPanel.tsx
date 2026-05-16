import { useState, useEffect, useRef } from 'react';
import type { TaskStatus, TaskPriority, ProjectMember } from '../../types/task';
import styles from './FilterPanel.module.css';

interface Filters {
    status: string;
    priority: string;
    assignee: string;
}

interface Props {
    statuses: TaskStatus[];
    priorities: TaskPriority[];
    members: ProjectMember[];
    filters: Filters;
    onFilterChange: (filters: Filters) => void;
    onApply: () => void;
    onClear: () => void;
}

export default function FilterPanel({ statuses, priorities, members, filters, onFilterChange, onApply, onClear }: Props) {
    const [open, setOpen] = useState(false);
    const ref = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handler = (e: MouseEvent) => { if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false); };
        document.addEventListener('mousedown', handler);
        return () => document.removeEventListener('mousedown', handler);
    }, []);

    const set = (key: keyof Filters) => (e: React.ChangeEvent<HTMLSelectElement>) => {
        onFilterChange({ ...filters, [key]: e.target.value });
    };

    const handleApply = () => { onApply(); setOpen(false); };
    const handleClear = () => { onClear(); setOpen(false); };

    return (
        <div className={styles.wrap} ref={ref}>
            <button className={styles.trigger} onClick={() => setOpen(!open)}>
                <svg className={styles.triggerIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"/></svg>
                Фильтры
            </button>
            {open && (
                <div className={styles.dropdown}>
                    <h4 className={styles.title}>Фильтры</h4>
                    <div className={styles.group}>
                        <label className={styles.label}>Статус</label>
                        <select className={styles.select} value={filters.status} onChange={set('status')}>
                            <option value="">Все статусы</option>
                            {statuses.map(s => <option key={s.id} value={s.name}>{s.name}</option>)}
                        </select>
                    </div>
                    <div className={styles.group}>
                        <label className={styles.label}>Приоритет</label>
                        <select className={styles.select} value={filters.priority} onChange={set('priority')}>
                            <option value="">Все приоритеты</option>
                            {priorities.map(p => <option key={p.id} value={p.title}>{p.title}</option>)}
                        </select>
                    </div>
                    <div className={styles.group}>
                        <label className={styles.label}>Исполнитель</label>
                        <select className={styles.select} value={filters.assignee} onChange={set('assignee')}>
                            <option value="">Все исполнители</option>
                            {members.filter(m => m.User.name !== 'Deleted User').map(m => (
                                <option key={m.User.user_id} value={m.User.user_id}>{m.User.name}</option>
                            ))}
                        </select>
                    </div>
                    <div className={styles.footer}>
                        <button className={styles.clearBtn} onClick={handleClear}>Сбросить</button>
                        <button className={styles.applyBtn} onClick={handleApply}>Применить</button>
                    </div>
                </div>
            )}
        </div>
    );
}
