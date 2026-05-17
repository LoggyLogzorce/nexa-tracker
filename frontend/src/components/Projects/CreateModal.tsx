import React, { useState, useEffect } from 'react';
import type { ProjectStatus, Priority } from '../../types/project';
import modalStyles from '../Dashboard/Modal.module.css';

interface Props { onClose: () => void; onSave: (data: { title: string; description: string; status: ProjectStatus; priority: Priority }) => Promise<void>; }

export default function CreateModal({ onClose, onSave }: Props) {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [status, setStatus] = useState<ProjectStatus>('Планирование');
    const [priority, setPriority] = useState<Priority>('Средний');

    useEffect(() => { const esc = (e: KeyboardEvent) => e.key === 'Escape' && onClose(); document.addEventListener('keydown', esc); return () => document.removeEventListener('keydown', esc); }, [onClose]);

    const handleSubmit = async (e: React.FormEvent) => { e.preventDefault(); try { await onSave({ title, description, status, priority }); onClose(); } catch { /* ignored */ } };

    return (
        <div className={modalStyles.overlay} onClick={onClose}>
            <div className={modalStyles.modal} onClick={e => e.stopPropagation()}>
                <div className={modalStyles.header}><h2>Новый проект</h2></div>
                <form onSubmit={handleSubmit} className={modalStyles.body}>
                    <label className={modalStyles.label}>Название
                        <input className={modalStyles.input} value={title} onChange={e => setTitle(e.target.value)} required placeholder="Название проекта" />
                    </label>
                    <label className={modalStyles.label}>Описание
                        <textarea className={modalStyles.textarea} rows={3} value={description} onChange={e => setDescription(e.target.value)} placeholder="Описание проекта" />
                    </label>
                    <div className={modalStyles.row}>
                        <label className={modalStyles.label}>Статус
                            <select className={modalStyles.select} value={status} onChange={e => setStatus(e.target.value as ProjectStatus)}>
                                <option value="В работе">В работе</option><option value="Планирование">Планирование</option><option value="Завершен">Завершен</option>
                            </select>
                        </label>
                        <label className={modalStyles.label}>Приоритет
                            <select className={modalStyles.select} value={priority} onChange={e => setPriority(e.target.value as Priority)}>
                                <option value="Низкий">Низкий</option><option value="Средний">Средний</option><option value="Высокий">Высокий</option>
                            </select>
                        </label>
                    </div>
                    <div className={modalStyles.footer}>
                        <button type="button" className={modalStyles.cancel} onClick={onClose}>Отмена</button>
                        <button type="submit" className={modalStyles.save}>Создать</button>
                    </div>
                </form>
            </div>
        </div>
    );
}
