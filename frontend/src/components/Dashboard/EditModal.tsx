import React, { useState, useEffect } from 'react';
import type {Project, ProjectStatus, Priority} from '../../types/project';
import styles from './Modal.module.css';

interface Props { project: Project; onClose: () => void; onSave: (p: Project) => void; }

export default function EditModal({ project, onClose, onSave }: Props) {
    const [form, setForm] = useState(project);

    useEffect(() => { const esc = (e: KeyboardEvent) => e.key === 'Escape' && onClose(); document.addEventListener('keydown', esc); return () => document.removeEventListener('keydown', esc); }, [onClose]);

    const handleSubmit = (e: React.FormEvent) => { e.preventDefault(); onSave(form); };

    return (
        <div className={styles.overlay} onClick={onClose}>
            <div className={styles.modal} onClick={e => e.stopPropagation()}>
                <div className={styles.header}><h2>Редактировать проект</h2></div>
                <form onSubmit={handleSubmit} className={styles.body}>
                    <label className={styles.label}>Название<input className={styles.input} value={form.title} onChange={e => setForm({ ...form, title: e.target.value })} required /></label>
                    <label className={styles.label}>Описание<textarea className={styles.textarea} rows={3} value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} required /></label>
                    <div className={styles.row}>
                        <label className={styles.label}>Статус
                            <select className={styles.select} value={form.status} onChange={e => setForm({ ...form, status: e.target.value as ProjectStatus })}>
                                <option value="В работе">В работе</option><option value="Планирование">Планирование</option><option value="Завершен">Завершен</option>
                            </select>
                        </label>
                        <label className={styles.label}>Приоритет
                            <select className={styles.select} value={form.priority} onChange={e => setForm({ ...form, priority: e.target.value as Priority })}>
                                <option value="Низкий">Низкий</option><option value="Средний">Средний</option><option value="Высокий">Высокий</option>
                            </select>
                        </label>
                    </div>
                    <div className={styles.footer}>
                        <button type="button" className={styles.cancel} onClick={onClose}>Отмена</button>
                        <button type="submit" className={styles.save}>Сохранить</button>
                    </div>
                </form>
            </div>
        </div>
    );
}