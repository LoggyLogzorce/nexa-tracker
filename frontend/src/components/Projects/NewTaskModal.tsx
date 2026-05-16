import React, { useState, useEffect } from 'react';
import type { TaskStatus, TaskPriority, ProjectMember } from '../../types/task';
import modalStyles from '../Dashboard/Modal.module.css';

interface Props {
    statuses: TaskStatus[];
    priorities: TaskPriority[];
    members: ProjectMember[];
    onClose: () => void;
    onSave: (data: { title: string; description: string; status: string; priority: string; deadline: string; assignee: string }) => void;
}

export default function NewTaskModal({ statuses, priorities, members, onClose, onSave }: Props) {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [status, setStatus] = useState(statuses[0]?.name || '');
    const [priority, setPriority] = useState(priorities[0]?.title || '');
    const [deadline, setDeadline] = useState('');
    const [assignee, setAssignee] = useState('');

    useEffect(() => {
        const esc = (e: KeyboardEvent) => e.key === 'Escape' && onClose();
        document.addEventListener('keydown', esc);
        return () => document.removeEventListener('keydown', esc);
    }, [onClose]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!title) return;
        onSave({ title, description, status, priority, deadline, assignee });
    };

    return (
        <div className={modalStyles.overlay} onClick={onClose}>
            <div className={modalStyles.modal} onClick={e => e.stopPropagation()}>
                <div className={modalStyles.header}>
                    <h2>Новая задача</h2>
                </div>
                <form onSubmit={handleSubmit} className={modalStyles.body}>
                    <label className={modalStyles.label}>Название задачи
                        <input className={modalStyles.input} value={title} onChange={e => setTitle(e.target.value)} required placeholder="Введите название" />
                    </label>
                    <label className={modalStyles.label}>Описание
                        <textarea className={modalStyles.textarea} rows={3} value={description} onChange={e => setDescription(e.target.value)} placeholder="Описание задачи" />
                    </label>
                    <div className={modalStyles.row}>
                        <label className={modalStyles.label}>Статус
                            <select className={modalStyles.select} value={status} onChange={e => setStatus(e.target.value)}>
                                {statuses.map(s => <option key={s.id} value={s.name}>{s.name}</option>)}
                            </select>
                        </label>
                        <label className={modalStyles.label}>Приоритет
                            <select className={modalStyles.select} value={priority} onChange={e => setPriority(e.target.value)}>
                                {priorities.map(p => <option key={p.id} value={p.title}>{p.title}</option>)}
                            </select>
                        </label>
                    </div>
                    <label className={modalStyles.label}>Исполнитель
                        <select className={modalStyles.select} value={assignee} onChange={e => setAssignee(e.target.value)}>
                            <option value="">Не назначен</option>
                            {members.filter(m => m.User.name !== 'Deleted User' && (m.role === 'owner' || m.role === 'member')).map(m => (
                                <option key={m.User.user_id} value={m.User.user_id}>{m.User.name}</option>
                            ))}
                        </select>
                    </label>
                    <label className={modalStyles.label}>Дедлайн
                        <input className={modalStyles.input} type="date" value={deadline} onChange={e => setDeadline(e.target.value)} />
                    </label>
                    <div className={modalStyles.footer}>
                        <button type="button" className={modalStyles.cancel} onClick={onClose}>Отмена</button>
                        <button type="submit" className={modalStyles.save}>Создать</button>
                    </div>
                </form>
            </div>
        </div>
    );
}
