import { useState, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { getChangelog } from '../../api/version';
import styles from './ChangelogPage.module.css';

export default function ChangelogPage() {
    const [md, setMd] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        let cancelled = false;
        getChangelog()
            .then(content => { if (!cancelled) setMd(content); })
            .catch(() => { if (!cancelled) setError('Не удалось загрузить список изменений'); })
            .finally(() => { if (!cancelled) setLoading(false); });
        return () => { cancelled = true; };
    }, []);

    if (loading) return <div className={styles.root}><p className={styles.loader}>Загрузка...</p></div>;
    if (error) return <div className={styles.root}><p className={styles.error}>{error}</p></div>;

    return (
        <div className={styles.root}>
            <div className={styles.md}>
                <ReactMarkdown remarkPlugins={[remarkGfm]}>
                    {md}
                </ReactMarkdown>
            </div>
        </div>
    );
}
