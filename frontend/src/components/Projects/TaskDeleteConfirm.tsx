import { useEffect } from 'react';
import type { Task } from '../../types/task';
import modalStyles from '../Dashboard/Modal.module.css';

interface Props { task: Task; onClose: () => void; onConfirm: (task: Task) => void; }

export default function TaskDeleteConfirm({ task, onClose, onConfirm }: Props) {
    useEffect(() => { const esc = (e: KeyboardEvent) => e.key === 'Escape' && onClose(); document.addEventListener('keydown', esc); return () => document.removeEventListener('keydown', esc); }, [onClose]);

    return (
        <div className={modalStyles.overlay} onClick={onClose}>
            <div className={`${modalStyles.modal} ${modalStyles.confirm}`} onClick={e => e.stopPropagation()}>
                <div className={modalStyles.iconWrap}>
                    <svg className={modalStyles.warnIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>
                </div>
                <h3 className={modalStyles.confirmTitle}>Удалить задачу?</h3>
                <p className={modalStyles.confirmDesc}>Вы уверены, что хотите удалить «{task.title}»? Это действие невозможно отменить.</p>
                <div className={modalStyles.footer}>
                    <button className={modalStyles.cancel} onClick={onClose}>Отмена</button>
                    <button className={modalStyles.delete} onClick={() => onConfirm(task)}>Удалить</button>
                </div>
            </div>
        </div>
    );
}
