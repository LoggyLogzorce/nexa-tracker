import { useDraggable } from '@dnd-kit/core';
import type { Task } from '../../types/task';
import styles from './KanbanCard.module.css';

const avatarImages = [
    'https://image.qwenlm.ai/public_source/87cc46bb-251e-456e-a7c0-8d3b41f35211/1281ab55e-b3ab-4e9e-be91-85c89891a3c3.png',
    'https://image.qwenlm.ai/public_source/87cc46bb-251e-456e-a7c0-8d3b41f35211/15745d149-3792-4261-a5c2-f7913b692184.png',
    'https://image.qwenlm.ai/public_source/87cc46bb-251e-456e-a7c0-8d3b41f35211/1818d7289-0020-4e2c-bba4-19c272b7c3ba.png',
    'https://image.qwenlm.ai/public_source/87cc46bb-251e-456e-a7c0-8d3b41f35211/1df69439c-4074-406c-b81d-7d3faea2d322.png'
];

interface Props { task: Task; }

export default function KanbanCard({ task }: Props) {
    const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({ id: String(task.id) });

    const style = transform ? { transform: `translate(${transform.x}px, ${transform.y}px)` } : undefined;

    const formatDate = (dateStr: string) => {
        const d = new Date(dateStr);
        return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
    };

    return (
        <div ref={setNodeRef} {...listeners} {...attributes} className={`${styles.card} ${isDragging ? styles.dragging : ''}`} style={style}>
            <div className={styles.body}>
                <div className={styles.topRow}>
                    <span className={styles.priorityBadge} style={{ backgroundColor: task.priority.color + '20', color: task.priority.color }}>
                        {task.priority.title}
                    </span>
                    <span className={styles.id}>#{task.id}</span>
                </div>
                <h4 className={styles.title}>{task.title}</h4>
                <div className={styles.bottomRow}>
                    {task.assignee ? (
                        <div className={styles.assignee}>
                            <div className={styles.avatar}>
                                <img src={avatarImages[0]} alt={task.assignee.name} />
                            </div>
                        </div>
                    ) : (
                        <span className={styles.noAssignee}>Не назначен</span>
                    )}
                    {task.deadline && (
                        <span className={styles.deadline}>
                            <svg className={styles.deadlineIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
                            {formatDate(task.deadline)}
                        </span>
                    )}
                </div>
            </div>
        </div>
    );
}
