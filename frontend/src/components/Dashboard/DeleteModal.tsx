import { useEffect } from 'react';
import type {Project} from '../../types/project';
import styles from './Modal.module.css';

interface Props { project: Project; onClose: () => void; onConfirm: (id: string) => void; }

export default function DeleteModal({ project, onClose, onConfirm }: Props) {
    useEffect(() => { const esc = (e: KeyboardEvent) => e.key === 'Escape' && onClose(); document.addEventListener('keydown', esc); return () => document.removeEventListener('keydown', esc); }, [onClose]);

    return (
        <div className={styles.overlay} onClick={onClose}>
            <div className={`${styles.modal} ${styles.confirm}`} onClick={e => e.stopPropagation()}>
                <div className={styles.iconWrap}>
                    <svg className={styles.warnIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>
                </div>
                <h3 className={styles.confirmTitle}>Удалить проект?</h3>
                <p className={styles.confirmDesc}>Вы уверены, что хотите удалить «{project.title}»? Это действие невозможно отменить.</p>
                <div className={styles.footer}>
                    <button className={styles.cancel} onClick={onClose}>Отмена</button>
                    <button className={styles.delete} onClick={() => onConfirm(project.id)}>Удалить</button>
                </div>
            </div>
        </div>
    );
}