import { useState, useEffect } from 'react';
import type { CustomStatus } from '../../types/project';
import { createStatusApi, updateStatusApi, deleteStatusApi } from '../../api/projects';
import { useNotifications } from '../../contexts/useNotifications';
import modalStyles from '../Dashboard/Modal.module.css';
import styles from './EditStatusesModal.module.css';

interface Props { projectId: string; statuses: CustomStatus[]; onClose: () => void; onSave: () => void; }

interface Item extends CustomStatus { _new?: boolean; _dirty?: boolean; }

let tempId = -1;

export default function EditStatusesModal({ projectId, statuses, onClose, onSave }: Props) {
    const { addNotification } = useNotifications();
    const [items, setItems] = useState<Item[]>(() => statuses.map(s => ({ ...s, _dirty: false })));
    const [deletedIds, setDeletedIds] = useState<number[]>([]);
    const [saving, setSaving] = useState(false);

    useEffect(() => { const esc = (e: KeyboardEvent) => e.key === 'Escape' && onClose(); document.addEventListener('keydown', esc); return () => document.removeEventListener('keydown', esc); }, [onClose]);

    const update = (id: number, patch: Partial<Item>) => setItems(prev => prev.map(s => s.id === id ? { ...s, ...patch, _dirty: true } : s));

    const remove = (id: number) => {
        if (id > 0) setDeletedIds(prev => [...prev, id]);
        setItems(prev => prev.filter(s => s.id !== id));
    };

    const add = () => {
        const id = tempId--;
        setItems(prev => [...prev, { id, name: '', color: '#808080', order_index: prev.length, _new: true }]);
    };

    const moveUp = (i: number) => {
        if (i === 0) return;
        setItems(prev => {
            const next = [...prev];
            [next[i - 1], next[i]] = [next[i], next[i - 1]];
            return next.map((s, idx) => ({ ...s, order_index: idx, _dirty: true }));
        });
    };

    const moveDown = (i: number) => {
        setItems(prev => {
            if (i === prev.length - 1) return prev;
            const next = [...prev];
            [next[i], next[i + 1]] = [next[i + 1], next[i]];
            return next.map((s, idx) => ({ ...s, order_index: idx, _dirty: true }));
        });
    };

    const handleSave = async () => {
        setSaving(true);
        try {
            const deletes = deletedIds.map(id => deleteStatusApi(projectId, id));
            await Promise.all(deletes);

            for (const s of items) {
                if (s._new) {
                    await createStatusApi(projectId, { name: s.name, color: s.color, order_index: s.order_index });
                } else if (s._dirty) {
                    await updateStatusApi(projectId, s.id, { name: s.name, color: s.color, order_index: s.order_index });
                }
            }

            addNotification('success', 'Статусы сохранены');
            onSave();
        } catch {
            addNotification('error', 'Ошибка при сохранении статусов');
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className={modalStyles.overlay} onClick={onClose}>
            <div className={modalStyles.modal} onClick={e => e.stopPropagation()}>
                <div className={modalStyles.header}><h2>Редактировать статусы</h2></div>
                <div className={modalStyles.body}>
                    {items.map((s, i) => (
                        <div key={s.id} className={styles.row}>
                            <span className={styles.dot} style={{ backgroundColor: s.color }} />
                            <input className={styles.input} value={s.name} onChange={e => update(s.id, { name: e.target.value })} placeholder="Название статуса" />
                            <input className={styles.colorInput} type="color" value={s.color} onChange={e => update(s.id, { color: e.target.value })} />
                            <button className={styles.btn} onClick={() => moveUp(i)} disabled={i === 0} title="Вверх"><svg className={styles.btnIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7"/></svg></button>
                            <button className={styles.btn} onClick={() => moveDown(i)} disabled={i === items.length - 1} title="Вниз"><svg className={styles.btnIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7"/></svg></button>
                            <button className={`${styles.btn} ${styles.btnDelete}`} onClick={() => remove(s.id)} title="Удалить"><svg className={styles.btnIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12"/></svg></button>
                        </div>
                    ))}
                    <button className={styles.addBtn} onClick={add}>+ Добавить статус</button>
                </div>
                <div className={modalStyles.footer}>
                    <button className={modalStyles.cancel} onClick={onClose} disabled={saving}>Отмена</button>
                    <button className={modalStyles.save} onClick={handleSave} disabled={saving}>{saving ? 'Сохранение...' : 'Сохранить'}</button>
                </div>
            </div>
        </div>
    );
}
