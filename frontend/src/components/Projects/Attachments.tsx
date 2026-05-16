import type { Attachment } from '../../types/task';
import styles from './Attachments.module.css';

interface Props { attachments: Attachment[]; }

function formatFileSize(bytes: number): string {
    if (bytes >= 1000000) return (bytes / 1000000).toFixed(1) + ' MB';
    if (bytes >= 1000) return (bytes / 1000).toFixed(0) + ' KB';
    return bytes + ' B';
}

function getFileConfig(mime: string, filename: string) {
    if (mime.includes('pdf')) return { icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z', color: '#dc2626', bg: '#fef2f2' };
    if (mime.includes('image') || filename.endsWith('.fig')) return { icon: 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z', color: '#9333ea', bg: '#f3e8ff' };
    return { icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z', color: '#64748b', bg: '#f1f5f9' };
}

export default function Attachments({ attachments }: Props) {
    const formatDate = (dateStr: string) => {
        const d = new Date(dateStr);
        return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
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
                        <a key={a.id} href="#" className={styles.item}>
                            <div className={styles.iconWrap} style={{ backgroundColor: cfg.bg }}>
                                <svg className={styles.fileIcon} fill="none" stroke={cfg.color} viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={cfg.icon} /></svg>
                            </div>
                            <div className={styles.fileInfo}>
                                <p className={styles.fileName}>{a.filename}</p>
                                <p className={styles.fileMeta}>{formatFileSize(a.file_size)} · {formatDate(a.created_at)}</p>
                            </div>
                        </a>
                    );
                })}
            </div>
            )}
        </div>
    );
}
