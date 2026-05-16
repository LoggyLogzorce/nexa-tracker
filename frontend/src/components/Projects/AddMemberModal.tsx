import { useState, useEffect, useRef } from 'react';
import { searchUsersApi, addMemberApi } from '../../api/projects';
import { useNotifications } from '../../contexts/useNotifications';
import type { UserSearchResult } from '../../api/projects';
import modalStyles from '../Dashboard/Modal.module.css';
import styles from './AddMemberModal.module.css';

interface Props { projectId: string; onClose: () => void; onSave: () => void; }

export default function AddMemberModal({ projectId, onClose, onSave }: Props) {
    const { addNotification } = useNotifications();
    const [email, setEmail] = useState('');
    const [results, setResults] = useState<UserSearchResult[]>([]);
    const [selected, setSelected] = useState<UserSearchResult | null>(null);
    const [role, setRole] = useState<'member' | 'read_only'>('member');
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');
    const inputRef = useRef<HTMLInputElement>(null);

    useEffect(() => { inputRef.current?.focus(); }, []);

    useEffect(() => {
        const esc = (e: KeyboardEvent) => e.key === 'Escape' && onClose();
        document.addEventListener('keydown', esc);
        return () => document.removeEventListener('keydown', esc);
    }, [onClose]);

    useEffect(() => {
        const timer = setTimeout(() => {
            if (selected || !email.trim()) { setResults([]); setError(''); return; }
            searchUsersApi(email.trim())
                .then(data => { setResults(data); if (data.length === 0) setError('Пользователи не найдены'); else setError(''); })
                .catch(() => setError('Ошибка поиска'));
        }, 450);
        return () => clearTimeout(timer);
    }, [email, selected]);

    const handleAdd = async () => {
        if (!selected) return;
        setSaving(true);
        try {
            await addMemberApi(projectId, selected.id, role);
            addNotification('success', `Участник ${selected.name} добавлен в проект`);
            onSave();
        } catch {
            addNotification('error', 'Ошибка при добавлении участника');
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className={modalStyles.overlay} onClick={onClose}>
            <div className={modalStyles.modal} onClick={e => e.stopPropagation()}>
                <div className={modalStyles.header}><h2>Добавить участника</h2></div>
                <div className={modalStyles.body}>
                    <label className={modalStyles.label}>Email пользователя
                        <input ref={inputRef} className={modalStyles.input} type="email" value={email} onChange={e => { setEmail(e.target.value); setSelected(null); }} placeholder="Введите email" autoComplete="off" />
                    </label>

                    {selected ? (
                        <>
                        <div className={styles.selected}>
                            <div className={styles.selectedInfo}>
                                <span className={styles.selectedName}>{selected.name}</span>
                                <span className={styles.selectedEmail}>{selected.email}</span>
                            </div>
                            <button className={styles.selectedRemove} onClick={() => setSelected(null)}>
                                <svg className={styles.selectedRemoveIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12"/></svg>
                            </button>
                        </div>
                        <div className={styles.roleSelector}>
                            <p className={styles.roleLabel}>Роль</p>
                            <div className={styles.roleOptions}>
                                <label className={`${styles.roleOption} ${role === 'member' ? styles.roleActive : ''}`}>
                                    <input type="radio" name="role" value="member" checked={role === 'member'} onChange={() => setRole('member')} />
                                    Участник
                                </label>
                                <label className={`${styles.roleOption} ${role === 'read_only' ? styles.roleActive : ''}`}>
                                    <input type="radio" name="role" value="read_only" checked={role === 'read_only'} onChange={() => setRole('read_only')} />
                                    Только чтение
                                </label>
                            </div>
                        </div>
                        </>
                    ) : error ? (
                        <p className={styles.hintError}>{error}</p>
                    ) : results.length > 0 ? (
                        <div className={styles.results}>
                            {results.map(u => (
                                <button key={u.id} className={styles.resultItem} onClick={() => { setSelected(u); setResults([]); }}>
                                    <div className={styles.resultInfo}>
                                        <span className={styles.resultName}>{u.name}</span>
                                        <span className={styles.resultEmail}>{u.email}</span>
                                    </div>
                                </button>
                            ))}
                        </div>
                    ) : null}

                    {email.trim() && !selected && !error && results.length === 0 && <p className={styles.hint}>Введите email для поиска</p>}
                </div>
                <div className={modalStyles.footer}>
                    <button className={modalStyles.cancel} onClick={onClose} disabled={saving}>Отмена</button>
                    <button className={modalStyles.save} onClick={handleAdd} disabled={!selected || saving}>{saving ? 'Добавление...' : 'Добавить'}</button>
                </div>
            </div>
        </div>
    );
}
