import { useState, useRef } from 'react';
import { useAuth } from '../../contexts/useAuth';
import {updateUserMeApi, uploadAvatarApi, changePasswordApi, deleteUserMeApi} from '../../api/auth';
import { useNotifications } from '../../contexts/useNotifications';
import Avatar from '../../components/UI/Avatar';
import modalStyles from '../../components/Dashboard/Modal.module.css';
import styles from './ProfilePage.module.css';

export default function ProfilePage() {
    const { user, logout } = useAuth();
    const { addNotification } = useNotifications();
    const fileInputRef = useRef<HTMLInputElement>(null);

    const [showEdit, setShowEdit] = useState(false);
    const [editName, setEditName] = useState(user?.name || '');
    const [saving, setSaving] = useState(false);

    const [showPassword, setShowPassword] = useState(false);
    const [currentPassword, setCurrentPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [savingPassword, setSavingPassword] = useState(false);

    const [showDelete, setShowDelete] = useState(false);
    const [deletePassword, setDeletePassword] = useState('');
    const [deletingAccount, setDeletingAccount] = useState(false);

    const handleSave = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        try {
            const updated = await updateUserMeApi({ name: editName });
            Object.assign(user!, updated);
            setShowEdit(false);
            addNotification('success', 'Профиль обновлён');
        } catch {
            addNotification('error', 'Ошибка при обновлении профиля');
        } finally {
            setSaving(false);
        }
    };

    const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;
        try {
            const updated = await uploadAvatarApi(file);
            Object.assign(user!, updated);
            addNotification('success', 'Аватар обновлён');
        } catch {
            addNotification('error', 'Ошибка при загрузке аватара');
        }
        if (fileInputRef.current) fileInputRef.current.value = '';
    };

    const handleChangePassword = async (e: React.FormEvent) => {
        e.preventDefault();
        if (newPassword !== confirmPassword) {
            addNotification('error', 'Пароли не совпадают');
            return;
        }
        setSavingPassword(true);
        try {
            await changePasswordApi({ current_password: currentPassword, new_password: newPassword });
            setShowPassword(false);
            setCurrentPassword('');
            setNewPassword('');
            setConfirmPassword('');
            addNotification('success', 'Пароль изменён');
        } catch {
            addNotification('error', 'Ошибка при смене пароля');
        } finally {
            setSavingPassword(false);
        }
    };


    const handleDeleteAccount = async (e: React.FormEvent) => {
        e.preventDefault();
        setDeletingAccount(true);
        try {
            await deleteUserMeApi(deletePassword);
            await logout();
        } catch {
            addNotification('error', 'Ошибка при удалении аккаунта');
        } finally {
            setDeletingAccount(false);
        }
    };

    const handleLogout = async () => {
        await logout();
    };

    const fmtDate = (iso: string) => new Date(iso).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' });

    return (
        <div className={styles.page}>
            <div className={styles.breadcrumb}>
                <span>Профиль</span>
            </div>

            <div className={styles.card}>
                <div className={styles.avatarSection}>
                    <div className={styles.avatarWrap} onClick={() => fileInputRef.current?.click()}>
                        <Avatar name={user?.name || '?'} avatarUrl={user?.avatar_url} size={80} />
                        <div className={styles.avatarOverlay}>
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M23 19a2 2 0 01-2 2H3a2 2 0 01-2-2V8a2 2 0 012-2h4l2-3h6l2 3h4a2 2 0 012 2z"/><circle cx="12" cy="13" r="4"/></svg>
                        </div>
                    </div>
                    <input ref={fileInputRef} type="file" accept="image/*" onChange={handleAvatarChange} hidden />
                    <h1 className={styles.name}>{user?.name || 'Пользователь'}</h1>
                    <span className={styles.role}>{user?.role}</span>
                </div>

                <div className={styles.infoGrid}>
                    <div className={styles.infoItem}>
                        <span className={styles.label}>Email</span>
                        <span className={styles.value}>{user?.email || '—'}</span>
                    </div>
                    <div className={styles.infoItem}>
                        <span className={styles.label}>Роль</span>
                        <span className={styles.value}>{user?.role || '—'}</span>
                    </div>
                    <div className={styles.infoItem}>
                        <span className={styles.label}>Зарегистрирован</span>
                        <span className={styles.value}>{user?.created_at ? fmtDate(user.created_at) : '—'}</span>
                    </div>
                </div>

                <div className={styles.actions}>
                    <button className={styles.editBtn} onClick={() => {
                        setEditName(user?.name || '');
                        setShowEdit(true);
                    }}>
                        Редактировать профиль
                    </button>
                    <button className={styles.editBtn} onClick={() => setShowPassword(true)}>
                        Сменить пароль
                    </button>
                    <button className={styles.logoutBtn} onClick={handleLogout}>
                        Выйти
                    </button>
                    <button className={styles.deleteBtn} onClick={() => setShowDelete(true)}>
                        Удалить аккаунт
                    </button>
                </div>
            </div>

            {showEdit && (
                <div className={modalStyles.overlay} onClick={() => setShowEdit(false)}>
                    <div className={modalStyles.modal} onClick={e => e.stopPropagation()}>
                        <div className={modalStyles.header}><h2>Редактировать профиль</h2></div>
                        <form onSubmit={handleSave} className={modalStyles.body}>
                            <label className={modalStyles.label}>
                                Имя
                                <input className={modalStyles.input} value={editName} onChange={e => setEditName(e.target.value)} required />
                            </label>
                            <div className={modalStyles.footer}>
                                <button type="button" className={modalStyles.cancel} onClick={() => setShowEdit(false)}>Отмена</button>
                                <button type="submit" className={modalStyles.save} disabled={saving}>{saving ? 'Сохранение...' : 'Сохранить'}</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {showPassword && (
                <div className={modalStyles.overlay} onClick={() => setShowPassword(false)}>
                    <div className={modalStyles.modal} onClick={e => e.stopPropagation()}>
                        <div className={modalStyles.header}><h2>Сменить пароль</h2></div>
                        <form onSubmit={handleChangePassword} className={modalStyles.body}>
                            <label className={modalStyles.label}>
                                Текущий пароль
                                <input className={modalStyles.input} type="password" value={currentPassword} onChange={e => setCurrentPassword(e.target.value)} required />
                            </label>
                            <label className={modalStyles.label}>
                                Новый пароль
                                <input className={modalStyles.input} type="password" value={newPassword} onChange={e => setNewPassword(e.target.value)} required minLength={8} />
                            </label>
                            <label className={modalStyles.label}>
                                Повторите новый пароль
                                <input className={modalStyles.input} type="password" value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} required minLength={8} />
                            </label>
                            <div className={modalStyles.footer}>
                                <button type="button" className={modalStyles.cancel} onClick={() => setShowPassword(false)}>Отмена</button>
                                <button type="submit" className={modalStyles.save} disabled={savingPassword}>{savingPassword ? 'Сохранение...' : 'Сохранить'}</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {showDelete && (
                <div className={modalStyles.overlay} onClick={() => setShowDelete(false)}>
                    <div className={modalStyles.modal} onClick={e => e.stopPropagation()}>
                        <div className={modalStyles.header}><h2>Удалить аккаунт</h2></div>
                        <form onSubmit={handleDeleteAccount} className={modalStyles.body}>
                            <p>Это действие необратимо. Введите пароль для подтверждения.</p>
                            <label className={modalStyles.label}>
                                Пароль
                                <input className={modalStyles.input} type="password" value={deletePassword} onChange={e => setDeletePassword(e.target.value)} required />
                            </label>
                            <div className={modalStyles.footer}>
                                <button type="button" className={modalStyles.cancel} onClick={() => setShowDelete(false)}>Отмена</button>
                                <button type="submit" className={modalStyles.save} disabled={deletingAccount}>{deletingAccount ? 'Удаление...' : 'Удалить'}</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}