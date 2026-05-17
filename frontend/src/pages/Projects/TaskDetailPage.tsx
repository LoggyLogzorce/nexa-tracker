import { useState, useEffect, useRef, useMemo } from 'react';
import { useParams, useSearchParams, useLocation, Link, useNavigate } from 'react-router-dom';
import type { Project } from '../../types/project';
import type { Task, Attachment, ProjectMember } from '../../types/task';
import { getProjectById, getTaskById, getTaskAttachments, deleteTaskApi, uploadTaskAttachmentApi, deleteTaskAttachmentApi, downloadAttachmentApi, getProjectMembers, updateTaskApi } from '../../api/projects';
import { useNotifications } from '../../contexts/useNotifications';
import TaskComments from '../../components/Projects/TaskComments';
import EditTaskModal from '../../components/Projects/EditTaskModal';
import Avatar from '../../components/UI/Avatar';
import modalStyles from '../../components/Dashboard/Modal.module.css';
import styles from './TaskDetailPage.module.css';

const fmtDate = (iso: string) => new Date(iso).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' });
const fmtDateShort = (iso: string) => new Date(iso).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
const fmtSize = (b: number) => b < 1048576 ? Math.round(b / 1024) + ' KB' : (b / 1048576).toFixed(1) + ' MB';
const getFileConfig = (mime: string, filename: string) => {
    const name = filename.toLowerCase();
    if (mime.includes('pdf')) return { icon: 'M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8zM10 12h4M10 16h4', color: '#dc2626', bg: '#fef2f2' };
    if (mime.includes('image')) return { icon: 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z', color: '#9333ea', bg: '#f3e8ff' };
    if (mime.includes('word') || name.endsWith('.doc') || name.endsWith('.docx')) return { icon: 'M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8zM10 12h4M10 16h4', color: '#2563eb', bg: '#eff6ff' };
    if (mime.includes('sheet') || mime.includes('excel') || name.endsWith('.xls') || name.endsWith('.xlsx')) return { icon: 'M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8zM9 12h6M9 16h6M12 12v4', color: '#16a34a', bg: '#f0fdf4' };
    if (mime.includes('presentation') || name.endsWith('.ppt') || name.endsWith('.pptx') || name.endsWith('.pptm')) return { icon: 'M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8zM12 12l-3 2v-4l3 2z', color: '#ea580c', bg: '#fff7ed' };
    if (mime.includes('zip') || mime.includes('rar') || mime.includes('gzip') || mime.includes('tar') || mime.includes('7z') || mime.includes('compress') || name.endsWith('.zip') || name.endsWith('.rar') || name.endsWith('.7z') || name.endsWith('.gz') || name.endsWith('.tar')) return { icon: 'M2 7a2 2 0 012-2h3l2 2h5a2 2 0 012 2v1H4V7zM4 12h16l-1 6a2 2 0 01-2 2H7a2 2 0 01-2-2l-1-6z', color: '#d97706', bg: '#fffbeb' };
    if (mime.includes('text') || name.endsWith('.txt')) return { icon: 'M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8zM10 12h4', color: '#64748b', bg: '#f1f5f9' };
    if (mime.includes('audio') || mime.includes('mp3') || mime.includes('wav') || mime.includes('aac') || mime.includes('ogg') || name.endsWith('.mp3') || name.endsWith('.wav') || name.endsWith('.aac') || name.endsWith('.ogg') || name.endsWith('.flac')) return { icon: 'M9 18V5l12-2v13M9 18a3 3 0 01-6 0 3 3 0 016 0zM21 16a3 3 0 01-6 0 3 3 0 016 0z', color: '#0891b2', bg: '#ecfeff' };
    if (mime.includes('video') || mime.includes('mp4') || mime.includes('avi') || mime.includes('mkv') || mime.includes('mov') || mime.includes('webm') || name.endsWith('.mp4') || name.endsWith('.avi') || name.endsWith('.mkv') || name.endsWith('.mov') || name.endsWith('.webm')) return { icon: 'M23 7l-7 5 7 5V7zM15 7H3a2 2 0 00-2 2v6a2 2 0 002 2h12a2 2 0 002-2V9a2 2 0 00-2-2z', color: '#db2777', bg: '#fdf2f8' };
    return { icon: 'M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8zM10 12h4M10 16h4', color: '#64748b', bg: '#f1f5f9' };
};

export default function TaskDetailPage() {
    const { id: projectId, taskId } = useParams<{ id: string; taskId: string }>();
    const location = useLocation();
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const archivedParam = searchParams.get('archived') === 'true';
    const projectFromState = (location.state as { project?: Project })?.project;
    const { addNotification } = useNotifications();
    const fileInputRef = useRef<HTMLInputElement>(null);

    const [project, setProject] = useState<Project | null>(projectFromState ?? null);
    const [task, setTask] = useState<Task | null>(null);
    const [attachments, setAttachments] = useState<Attachment[]>([]);
    const [members, setMembers] = useState<ProjectMember[]>([]);
    const [loading, setLoading] = useState(true);
    const [uploading, setUploading] = useState(false);
    const [deleteAttach, setDeleteAttach] = useState<Attachment | null>(null);
    const [showEditTask, setShowEditTask] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showDeleteTask, setShowDeleteTask] = useState(false);
    const [showEditDesc, setShowEditDesc] = useState(false);
    const [editDescText, setEditDescText] = useState('');
    const [commentsOpen, setCommentsOpen] = useState(true);

    const allMembers = useMemo(() => {
        if (!project) return members;
        const ownerMember: ProjectMember = {
            project_id: project.id,
            role: 'owner',
            User: {
                user_id: project.owner.id,
                name: project.owner.name,
                email: project.owner.email,
                avatar_url: project.owner.avatar_url,
            },
        };
        return [ownerMember, ...members];
    }, [project, members]);

    useEffect(() => {
        if (!projectId || !taskId) return;

        const loadData = async () => {
            try {
                const [taskData, attachData, membersData, projectData] = await Promise.all([
                    getTaskById(projectId, Number(taskId), archivedParam),
                    getTaskAttachments(projectId, Number(taskId), archivedParam),
                    getProjectMembers(projectId),
                    getProjectById(projectId),
                ]);
                setTask(taskData);
                setAttachments(attachData);
                setMembers(membersData);
                setProject(projectData);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Ошибка загрузки');
            } finally {
                setLoading(false);
            }
        };

        loadData();
    }, [projectId, taskId, archivedParam]);

    const handleDeleteTask = () => {
        setShowDeleteTask(true);
    };

    const handleConfirmDelete = async () => {
        if (!projectId || !task) return;
        try {
            await deleteTaskApi(projectId, task.id, task.is_archive);
            addNotification('success', 'Задача удалена');
            navigate(`/projects/${projectId}`);
        } catch {
            addNotification('error', 'Ошибка при удалении задачи');
        }
    };

    const handleUploadAttachment = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file || !projectId || !task) return;
        setUploading(true);
        try {
            const attachment = await uploadTaskAttachmentApi(projectId, task.id, file, task.is_archive);
            setAttachments(prev => [...prev, attachment]);
            addNotification('success', 'Файл загружен');
        } catch {
            addNotification('error', 'Ошибка при загрузке файла');
        } finally {
            setUploading(false);
            if (fileInputRef.current) fileInputRef.current.value = '';
        }
    };

    const handleDeleteAttachment = async (attachment: Attachment) => {
        if (!projectId || !task) return;
        try {
            await deleteTaskAttachmentApi(projectId, task.id, attachment.id, task.is_archive);
            setAttachments(prev => prev.filter(a => a.id !== attachment.id));
            addNotification('success', 'Вложение удалено');
        } catch {
            addNotification('error', 'Ошибка при удалении вложения');
        }
    };

    const handleEditTask = async (data: { title: string; description: string; status: string; priority: string; deadline: string; assignee: string }) => {
        if (!projectId || !task || !project) return;
        const statusObj = project.statuses.find(s => s.name === data.status) || project.statuses[0];
        const priorityObj = project.priorities.find(p => p.title === data.priority) || project.priorities[0];
        const assigneeUser = data.assignee ? allMembers.find(m => m.User.user_id === data.assignee) : null;
        try {
            const updated = await updateTaskApi(projectId, task.id, {
                title: data.title,
                description: data.description,
                deadline: data.deadline ? `${data.deadline}T00:00:00Z` : null,
                status_id: statusObj.id,
                priority_id: priorityObj.id,
                assignee_id: assigneeUser?.User.user_id ?? null,
            }, task.is_archive);
            setTask(updated);
            setShowEditTask(false);
            addNotification('success', 'Задача обновлена');
        } catch {
            addNotification('error', 'Ошибка при обновлении задачи');
        }
    };

    const handleDownloadAttachment = async (attachment: Attachment) => {
        if (!projectId || !task) return;
        try {
            await downloadAttachmentApi(projectId, task.id, attachment.id, attachment.filename, task.is_archive);
        } catch {
            addNotification('error', 'Ошибка при скачивании файла');
        }
    };

    const handleEditDescription = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!projectId || !task) return;
        try {
            const updated = await updateTaskApi(projectId, task.id, { description: editDescText }, task.is_archive);
            setTask(updated);
            setShowEditDesc(false);
            addNotification('success', 'Описание обновлено');
        } catch {
            addNotification('error', 'Ошибка при обновлении описания');
        }
    };

    const handleArchiveToggle = async () => {
        if (!projectId || !task) return;
        const newState = !task.is_archive;
        try {
            const updated = await updateTaskApi(projectId, task.id, { is_archive: newState }, task.is_archive);
            setTask(updated);
            addNotification('success', newState ? 'Задача архивирована' : 'Задача восстановлена');
        } catch {
            addNotification('error', 'Ошибка при изменении статуса задачи');
        }
    };

    if (loading) {
        return <div className={styles.page}><div className={styles.centerMsg}>Загрузка задачи...</div></div>;
    }

    if (error || !task || !project) {
        return (
            <div className={styles.page}>
                <div className={styles.centerMsg}>{error || 'Задача не найдена'}</div>
            </div>
        );
    }

    const deadline = task.deadline
        ? new Date(task.deadline + 'T00:00:00').toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
        : '—';

    const assigneeContent = task.assignee ? (
        <span className={styles.userChip}>
            <Avatar name={task.assignee.name} avatarUrl={task.assignee.avatar_url} size={22} />
            {task.assignee.name}
        </span>
    ) : (
        <span style={{ color: '#7b7f9e' }}>—</span>
    );

    const reporterContent = (
        <span className={styles.userChip}>
            <Avatar name={task.reporter.name} avatarUrl={task.reporter.avatar_url} size={22} />
            {task.reporter.name}
        </span>
    );

    return (
        <div className={`${styles.page} ${commentsOpen ? styles.hasComments : ''}`}>
            <div className={styles.taskScroll}>
                <div className={styles.breadcrumb}>
                    <Link to="/projects">Проекты</Link>
                    <span>›</span>
                    <Link to={`/projects/${project.id}`}>{project.title}</Link>
                    <span>›</span>
                    <span>Задача #{task.id}</span>
                </div>

                {/* Title + Meta */}
                <div className={styles.card}>
                    <div className={styles.titleBlock}>
                        <h1 className={styles.taskTitle}>{task.title}</h1>
                        <div className={styles.titleActions}>
                            <button className={task.is_archive ? styles.unarchiveBtn : styles.archiveBtn} onClick={handleArchiveToggle} title={task.is_archive ? 'Восстановить задачу' : 'Архивировать задачу'}>
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    {task.is_archive ? (
                                        <>
                                            <path d="M4 4h16v2H4V4z"/><path d="M4 8h16l-2 12H6L4 8z"/><path d="M12 12v4"/><path d="M10 14l2 2 2-2"/>
                                        </>
                                    ) : (
                                        <>
                                            <path d="M4 4h16v2H4V4z"/><path d="M4 8h16l-2 12H6L4 8z"/><path d="M12 16v-4"/><path d="M10 14l2-2 2 2"/>
                                        </>
                                    )}
                                </svg>
                            </button>
                            <button className={styles.editTaskBtn} onClick={() => setShowEditTask(true)} title="Редактировать задачу">
                                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                            </button>
                            <button className={styles.deleteTaskBtn} onClick={handleDeleteTask} title="Удалить задачу">
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
                            </button>
                        </div>
                    </div>
                    <div className={styles.metaGrid}>
                        <div className={styles.metaCell}>
                            <div className={styles.metaKey}>Статус</div>
                            <div className={styles.metaVal}>
                                <span className={styles.chip} style={{ backgroundColor: task.status.color + '18', color: task.status.color }}>
                                    <span className={styles.dot} style={{ background: task.status.color }} />
                                    {task.status.name}
                                </span>
                            </div>
                        </div>
                        <div className={styles.metaCell}>
                            <div className={styles.metaKey}>Приоритет</div>
                            <div className={styles.metaVal}>
                                <span className={styles.chip} style={{ backgroundColor: task.priority.color + '18', color: task.priority.color }}>
                                    <span className={styles.dot} style={{ background: task.priority.color }} />
                                    {task.priority.title}
                                </span>
                            </div>
                        </div>
                        <div className={styles.metaCell}>
                            <div className={styles.metaKey}>Исполнитель</div>
                            <div className={styles.metaVal}>{assigneeContent}</div>
                        </div>
                        <div className={styles.metaCell}>
                            <div className={styles.metaKey}>Автор</div>
                            <div className={styles.metaVal}>{reporterContent}</div>
                        </div>
                        <div className={styles.metaCell}>
                            <div className={styles.metaKey}>Дедлайн</div>
                            <div className={styles.metaVal}>{deadline}</div>
                        </div>
                        <div className={styles.metaCell}>
                            <div className={styles.metaKey}>Создана</div>
                            <div className={styles.metaVal} style={{ color: '#7b7f9e' }}>{fmtDate(task.created_at)}</div>
                        </div>
                    </div>
                </div>

                {/* Description */}
                <div className={styles.card}>
                    <div className={styles.cardHeader}>
                        <span className={styles.cardLabel}>Описание</span>
                        <button className={styles.editDescBtn} onClick={() => { setEditDescText(task.description || ''); setShowEditDesc(true); }} title="Редактировать описание">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                        </button>
                    </div>
                    <div className={styles.descBody}>
                        {task.description || <span className={styles.descEmpty}>Описание не указано</span>}
                    </div>
                </div>

                {/* Attachments */}
                <div className={styles.card}>
                    <div className={styles.cardHeader}>
                        <span className={styles.cardLabel}>Вложения</span>
                        <span style={{ fontSize: '0.75rem', color: '#7b7f9e' }}>{attachments.length}</span>
                    </div>
                    <div className={styles.attachList}>
                        {attachments.length === 0 ? (
                            <div style={{ padding: '0.75rem 1.25rem', color: '#7b7f9e', fontSize: '0.8125rem' }}>Вложений нет</div>
                        ) : (
                            attachments.map(a => (
                                <div key={a.id} className={styles.attachItem}>
                                    <div className={styles.attachIcon} style={{ backgroundColor: getFileConfig(a.mime_type, a.filename).bg }}>
                                        <svg width="16" height="16" fill="none" stroke={getFileConfig(a.mime_type, a.filename).color} viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={getFileConfig(a.mime_type, a.filename).icon} /></svg>
                                    </div>
                                    <div className={styles.attachInfo} onClick={() => handleDownloadAttachment(a)} style={{ cursor: 'pointer' }}>
                                        <div className={styles.attachName}>{a.filename}</div>
                                        <div className={styles.attachMeta}>{fmtSize(a.file_size)} · {fmtDateShort(a.created_at)}</div>
                                    </div>
                                    <button className={styles.deleteAttachBtn} onClick={() => setDeleteAttach(a)} title="Удалить вложение">
                                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#f87171" strokeWidth="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
                                    </button>
                                </div>
                            ))
                        )}
                        <div className={styles.attachUpload}>
                            <input ref={fileInputRef} type="file" onChange={handleUploadAttachment} hidden />
                            <button className={styles.uploadBtn} onClick={() => fileInputRef.current?.click()} disabled={uploading}>
                                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                                {uploading ? 'Загрузка...' : 'Загрузить файл'}
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <button className={`${styles.commentsToggle} ${commentsOpen ? styles.commentsToggleOpen : ''}`} onClick={() => setCommentsOpen(!commentsOpen)} title={commentsOpen ? 'Скрыть чат' : 'Показать чат'}>
                <svg width="18" height="18" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                    <path d={commentsOpen ? 'M9 5l7 7-7 7' : 'M15 19l-7-7 7-7'} />
                </svg>
            </button>
            <div className={`${styles.sidebarWrap} ${commentsOpen ? styles.open : ''}`}>
                <TaskComments projectId={project.id} taskId={task.id} archived={task.is_archive} />
            </div>

            {showEditDesc && (
                <div className={modalStyles.overlay} onClick={() => setShowEditDesc(false)}>
                    <div className={modalStyles.modal} onClick={e => e.stopPropagation()}>
                        <div className={modalStyles.header}><h2>Редактировать описание</h2></div>
                        <form onSubmit={handleEditDescription} className={modalStyles.body}>
                            <textarea className={modalStyles.textarea} rows={6} value={editDescText} onChange={e => setEditDescText(e.target.value)} placeholder="Описание задачи" />
                            <div className={modalStyles.footer}>
                                <button type="button" className={modalStyles.cancel} onClick={() => setShowEditDesc(false)}>Отмена</button>
                                <button type="submit" className={modalStyles.save}>Сохранить</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {showEditTask && project && (
                <EditTaskModal task={task} statuses={project.statuses} priorities={project.priorities} members={allMembers} onClose={() => setShowEditTask(false)} onSave={handleEditTask} />
            )}

            {showDeleteTask && (
                <div className={modalStyles.overlay} onClick={() => setShowDeleteTask(false)}>
                    <div className={`${modalStyles.modal} ${modalStyles.confirm}`} onClick={e => e.stopPropagation()}>
                        <div className={modalStyles.iconWrap}>
                            <svg className={modalStyles.warnIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>
                        </div>
                        <h3 className={modalStyles.confirmTitle}>Удалить задачу?</h3>
                        <p className={modalStyles.confirmDesc}>Задача будет безвозвратно удалена. Вы можете архивировать её вместо удаления.</p>
                        <div className={modalStyles.footer}>
                            <button className={modalStyles.cancel} onClick={() => setShowDeleteTask(false)}>Отмена</button>
                            <button className={modalStyles.save} onClick={() => { handleArchiveToggle(); setShowDeleteTask(false); }} style={{ background: '#d97706' }}>Архивировать</button>
                            <button className={modalStyles.delete} onClick={() => { handleConfirmDelete(); setShowDeleteTask(false); }}>Удалить</button>
                        </div>
                    </div>
                </div>
            )}

            {deleteAttach && (
                <div className={modalStyles.overlay} onClick={() => setDeleteAttach(null)}>
                    <div className={`${modalStyles.modal} ${modalStyles.confirm}`} onClick={e => e.stopPropagation()}>
                        <div className={modalStyles.iconWrap}>
                            <svg className={modalStyles.warnIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>
                        </div>
                        <h3 className={modalStyles.confirmTitle}>Удалить вложение?</h3>
                        <p className={modalStyles.confirmDesc}>Вы уверены, что хотите удалить «{deleteAttach.filename}»? Это действие невозможно отменить.</p>
                        <div className={modalStyles.footer}>
                            <button className={modalStyles.cancel} onClick={() => setDeleteAttach(null)}>Отмена</button>
                            <button className={modalStyles.delete} onClick={() => { handleDeleteAttachment(deleteAttach); setDeleteAttach(null); }}>Удалить</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
