import type { Attachment } from '../../types/task';
import { downloadAttachmentApi } from '../../api/projects';
import styles from './Attachments.module.css';

interface Props { attachments: Attachment[]; projectId: string; }

function formatFileSize(bytes: number): string {
    if (bytes >= 1000000) return (bytes / 1000000).toFixed(1) + ' MB';
    if (bytes >= 1000) return (bytes / 1000).toFixed(0) + ' KB';
    return bytes + ' B';
}

function getFileConfig(mime: string, filename: string) {
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
}

export default function Attachments({ attachments, projectId }: Props) {
    const formatDate = (dateStr: string) => {
        const d = new Date(dateStr);
        return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
    };

    const handleDownload = async (a: Attachment) => {
        try {
            await downloadAttachmentApi(projectId, a.task_id, a.id, a.filename);
        } catch { /* ignored */ }
    };

    return (
        <div className={styles.wrap}>
            <h3 className={styles.title}>Вложения</h3>
            {attachments.length === 0 ? (
                <p className={styles.empty}>Нет вложений</p>
            ) : (
            <div className={styles.list}>
                {attachments.map(a => {
                    const cfg = getFileConfig(a.mime_type, a.filename);
                    return (
                        <div key={a.id} className={styles.item} onClick={() => handleDownload(a)}>
                            <div className={styles.iconWrap} style={{ backgroundColor: cfg.bg }}>
                                <svg className={styles.fileIcon} fill="none" stroke={cfg.color} viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={cfg.icon} /></svg>
                            </div>
                            <div className={styles.fileInfo}>
                                <p className={styles.fileName}>{a.filename}</p>
                                <p className={styles.fileMeta}>{formatFileSize(a.file_size)} · {formatDate(a.created_at)}</p>
                            </div>
                        </div>
                    );
                })}
            </div>
            )}
        </div>
    );
}
