import { useState, useEffect, useRef } from 'react';
import type { Comment } from '../../types/comment';
import { useAuth } from '../../contexts/useAuth';
import { getTaskComments, createTaskComment, updateTaskCommentApi, deleteTaskCommentApi } from '../../api/projects';
import { useNotifications } from '../../contexts/useNotifications';
import Avatar from '../UI/Avatar';
import styles from './TaskComments.module.css';

interface Props { projectId: string; taskId: number; archived?: boolean; }

export default function TaskComments({ projectId, taskId, archived }: Props) {
    const { user } = useAuth();
    const { addNotification } = useNotifications();
    const [comments, setComments] = useState<Comment[]>([]);
    const [text, setText] = useState('');
    const [loading, setLoading] = useState(true);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [editText, setEditText] = useState('');
    const listRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        getTaskComments(projectId, taskId, archived)
            .then(setComments)
            .finally(() => setLoading(false));
    }, [projectId, taskId, archived]);

    const scrollToBottom = () => {
        requestAnimationFrame(() => {
            if (listRef.current) listRef.current.scrollTop = listRef.current.scrollHeight;
        });
    };

    useEffect(() => { if (!loading) scrollToBottom(); }, [loading, comments.length]);

    const handleSend = async () => {
        const trimmed = text.trim();
        if (!trimmed) return;
        setText('');
        try {
            const comment = await createTaskComment(projectId, taskId, trimmed, archived);
            setComments(prev => [...prev, comment]);
        } catch { /* ignored */ }
    };

    const handleKey = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleSend(); }
    };

    const handleEditKey = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleSaveEdit(); }
        if (e.key === 'Escape') { setEditingId(null); }
    };

    const handleStartEdit = (c: Comment) => {
        setEditingId(c.id);
        setEditText(c.content);
    };

    const handleSaveEdit = async () => {
        if (editingId === null) return;
        const trimmed = editText.trim();
        if (!trimmed) return;
        try {
            const updated = await updateTaskCommentApi(projectId, taskId, editingId, trimmed, archived);
            setComments(prev => prev.map(c => c.id === editingId ? updated : c));
            setEditingId(null);
            addNotification('success', 'Комментарий обновлён');
        } catch {
            addNotification('error', 'Ошибка при обновлении комментария');
        }
    };

    const handleDelete = async (commentId: number) => {
        if (!confirm('Удалить комментарий?')) return;
        try {
            await deleteTaskCommentApi(projectId, taskId, commentId, archived);
            setComments(prev => prev.filter(c => c.id !== commentId));
            addNotification('success', 'Комментарий удалён');
        } catch {
            addNotification('error', 'Ошибка при удалении комментария');
        }
    };

    const sameDay = (a: string, b: string) => new Date(a).toDateString() === new Date(b).toDateString();

    const fmtDate = (iso: string) => new Date(iso).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
    const fmtTime = (iso: string) => new Date(iso).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });

    const currentUserId = user?.id;

    return (
        <div className={styles.panel}>
            <div className={styles.header}>
                <span className={styles.title}>Комментарии</span>
                <span className={styles.count}>{comments.length}</span>
            </div>

            <div className={styles.messages} ref={listRef}>
                {loading ? (
                    <div className={styles.empty}>Загрузка...</div>
                ) : comments.length === 0 ? (
                    <div className={styles.empty}>
                        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z"/></svg>
                        <span>Комментариев пока нет</span>
                    </div>
                ) : (
                    comments.map((c, i) => (
                        <div key={c.id}>
                            {(i === 0 || !sameDay(comments[i - 1].created_at, c.created_at)) && (
                                <div className={styles.dateDivider}>{fmtDate(c.created_at)}</div>
                            )}
                            <div className={`${styles.msg} ${c.user.id === currentUserId ? styles.own : ''}`}>
                                {c.user.id !== currentUserId && (
                                    <div className={styles.avatarWrap}><Avatar name={c.user.name} avatarUrl={c.user.avatar_url} size={28} /></div>
                                )}
                                <div className={styles.msgContent}>
                                    {c.user.id !== currentUserId && (
                                        <div className={styles.author}>{c.user.name}</div>
                                    )}
                                    {editingId === c.id ? (
                                        <div className={styles.editWrap}>
                                            <textarea className={styles.editArea} value={editText} onChange={e => setEditText(e.target.value)} onKeyDown={handleEditKey} autoFocus />
                                            <div className={styles.editActions}>
                                                <button className={styles.editBtn} onClick={handleSaveEdit}>Сохранить</button>
                                                <button className={styles.editCancel} onClick={() => setEditingId(null)}>Отмена</button>
                                            </div>
                                        </div>
                                    ) : (
                                        <div className={`${styles.bubble} ${c.user.id === currentUserId ? styles.ownBubble : styles.otherBubble}`}>
                                            {c.content}
                                        </div>
                                    )}
                                    {editingId !== c.id && (
                                        <div className={styles.time}>
                                            {c.user.id === currentUserId && (
                                                <span className={styles.commentActions}>
                                                    <button className={styles.actionBtn} onClick={() => handleStartEdit(c)} title="Редактировать">
                                                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                                                    </button>
                                                    <button className={styles.actionBtn} onClick={() => handleDelete(c.id)} title="Удалить">
                                                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
                                                    </button>
                                                </span>
                                            )}
                                            {fmtTime(c.created_at)}
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))
                )}
            </div>

            <div className={styles.inputArea}>
                <div className={styles.inputRow}>
                    <textarea
                        className={styles.textarea}
                        value={text}
                        onChange={e => setText(e.target.value)}
                        onKeyDown={handleKey}
                        placeholder="Написать комментарий..."
                        rows={1}
                    />
                    <button className={styles.sendBtn} onClick={handleSend} disabled={!text.trim()}>
                        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2.5"><line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/></svg>
                    </button>
                </div>
            </div>
        </div>
    );
}
