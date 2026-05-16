import { useState, useEffect } from 'react';
import type { CustomPriority } from '../../types/project';
import { createPriorityApi, updatePriorityApi, deletePriorityApi } from '../../api/projects';
import { useNotifications } from '../../contexts/useNotifications';
import modalStyles from '../Dashboard/Modal.module.css';
import styles from './EditStatusesModal.module.css';

interface Props { projectId: string; priorities: CustomPriority[]; onClose: () => void; onSave: () => void; }

interface Item extends CustomPriority { _new?: boolean; _dirty?: boolean; }

let tempId = -1;

export default function EditPrioritiesModal({ projectId, priorities, onClose, onSave }: Props) {
    const { addNotification } = useNotifications();
    const [items, setItems] = useState<Item[]>(() => priorities.map(s => ({ ...s, _dirty: false })));
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
        setItems(prev => [...prev, { id, title: '', color: '#808080', _new: true }]);
    };

    const handleSave = async () => {
        setSaving(true);
        try {
            const deletes = deletedIds.map(id => deletePriorityApi(projectId, id));
            await Promise.all(deletes);

            for (const s of items) {
                if (s._new) {
                    await createPriorityApi(projectId, { title: s.title, color: s.color });
                } else if (s._dirty) {
                    await updatePriorityApi(projectId, s.id, { title: s.title, color: s.color });
                }
            }

            addNotification('success', 'Приоритеты сохранены');
            onSave();
        } catch {
            addNotification('error', 'Ошибка при сохранении приоритетов');
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className={modalStyles.overlay} onClick={onClose}>
            <div className={modalStyles.modal} onClick={e => e.stopPropagation()}>
                <div className={modalStyles.header}><h2>Редактировать приоритеты</h2></div>
                <div className={modalStyles.body}>
                    {items.map(s => (
                        <div key={s.id} className={styles.row}>
                            <span className={styles.dot} style={{ backgroundColor: s.color }} />
                            <input className={styles.input} value={s.title} onChange={e => update(s.id, { title: e.target.value })} placeholder="Название приоритета" />
                            <input className={styles.colorInput} type="color" value={s.color} onChange={e => update(s.id, { color: e.target.value })} />
                            <button className={`${styles.btn} ${styles.btnDelete}`} onClick={() => remove(s.id)} title="Удалить">
                                <svg className={styles.btnIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12"/></svg>
                            </button>
                        </div>
                    ))}
                    <button className={styles.addBtn} onClick={add}>+ Добавить приоритет</button>
                </div>
                <div className={modalStyles.footer}>
                    <button className={modalStyles.cancel} onClick={onClose} disabled={saving}>Отмена</button>
                    <button className={modalStyles.save} onClick={handleSave} disabled={saving}>{saving ? 'Сохранение...' : 'Сохранить'}</button>
                </div>
            </div>
        </div>
    );
}
